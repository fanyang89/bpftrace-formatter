use super::snapshot::DocumentSnapshot;
use std::collections::{HashMap, HashSet};
use std::fs;
use std::path::{Path, PathBuf};
use std::sync::{Arc, Mutex};
use tokio::task;
use tower_lsp::lsp_types::{
    CompletionItem, CompletionItemKind, CompletionList, CompletionTextEdit, Documentation,
    Position, Range, TextEdit,
};
use tree_sitter::Node;

const MAX_ITEMS: usize = 200;

const STATIC_TARGETS: &[(&str, &str)] = &[
    ("hardware:branch-instructions", "hardware event"),
    ("hardware:branch-misses", "hardware event"),
    ("hardware:bus-cycles", "hardware event"),
    ("hardware:cache-misses", "hardware event"),
    ("hardware:cache-references", "hardware event"),
    ("hardware:cpu-cycles", "hardware event"),
    ("hardware:instructions", "hardware event"),
    ("hardware:ref-cycles", "hardware event"),
    ("hardware:stalled-cycles-backend", "hardware event"),
    ("hardware:stalled-cycles-frontend", "hardware event"),
    ("interval:ms:100", "interval probe"),
    ("interval:s:1", "interval probe"),
    ("interval:us:100", "interval probe"),
    ("profile:hz:99", "profile probe"),
    ("profile:ms:10", "profile probe"),
    ("profile:s:1", "profile probe"),
    ("software:alignment-faults", "software event"),
    ("software:bpf-output", "software event"),
    ("software:context-switches", "software event"),
    ("software:cpu-clock", "software event"),
    ("software:cpu-migrations", "software event"),
    ("software:dummy", "software event"),
    ("software:emulation-faults", "software event"),
    ("software:major-faults", "software event"),
    ("software:minor-faults", "software event"),
    ("software:page-faults", "software event"),
    ("software:task-clock", "software event"),
];

#[derive(Debug, Clone)]
pub(super) enum CompletionRequest {
    ProbeTargets {
        provider: String,
        target_prefix: String,
        range: Range,
    },
    ArgsFields {
        probes: String,
        field_prefix: String,
        range: Range,
    },
}

#[derive(Debug, Clone, PartialEq, Eq)]
struct TargetCandidate {
    full_name: String,
    detail: String,
}

#[derive(Debug, Clone, PartialEq, Eq)]
struct FieldCandidate {
    probe: String,
    declaration: String,
    name: String,
}

#[derive(Debug, Default)]
pub(super) struct WorkspaceMetadata {
    targets: Vec<TargetCandidate>,
    fields: HashMap<String, Vec<FieldCandidate>>,
}

impl WorkspaceMetadata {
    pub(super) fn add_snapshot(&mut self, snapshot: &DocumentSnapshot) {
        let Some(tree) = snapshot.tree.as_ref() else {
            return;
        };
        let text = snapshot.text.as_ref();
        walk(tree.root_node(), &mut |node| {
            if node.kind() == "probe" {
                let full_name = text[node.byte_range()].trim().to_string();
                if full_name.contains(':') {
                    self.targets.push(TargetCandidate {
                        full_name,
                        detail: "workspace probe".to_string(),
                    });
                }
                return;
            }
            if node.kind() != "field_expression"
                || node
                    .child_by_field_name("argument")
                    .is_none_or(|argument| argument.kind() != "args_keyword")
            {
                return;
            }
            let Some(field) = node.child_by_field_name("field") else {
                return;
            };
            let Some(action_block) = ancestor(node, "action_block") else {
                return;
            };
            let Some(probes_list) = named_child(action_block, "probes_list") else {
                return;
            };
            let name = text[field.byte_range()].to_string();
            let mut cursor = probes_list.walk();
            for probe in probes_list
                .named_children(&mut cursor)
                .filter(|child| child.kind() == "probe")
            {
                let probe = text[probe.byte_range()].trim().to_string();
                self.fields
                    .entry(probe.clone())
                    .or_default()
                    .push(FieldCandidate {
                        probe,
                        declaration: format!("{name} (observed in workspace)"),
                        name: name.clone(),
                    });
            }
        });
    }
}

#[derive(Debug, Default)]
struct KernelMetadataCache {
    targets: Option<Vec<TargetCandidate>>,
    fields: HashMap<String, Vec<FieldCandidate>>,
}

#[derive(Debug, Clone, Default)]
pub(super) struct ProbeCatalog {
    cache: Arc<Mutex<KernelMetadataCache>>,
}

impl ProbeCatalog {
    pub(super) async fn complete(
        &self,
        request: CompletionRequest,
        workspace: WorkspaceMetadata,
    ) -> CompletionList {
        match request {
            CompletionRequest::ProbeTargets {
                provider,
                target_prefix,
                range,
            } => {
                let mut candidates = workspace.targets;
                candidates.extend(self.kernel_targets().await);
                candidates.extend(STATIC_TARGETS.iter().map(|(full_name, detail)| {
                    TargetCandidate {
                        full_name: (*full_name).to_string(),
                        detail: (*detail).to_string(),
                    }
                }));
                target_completions(&provider, &target_prefix, range, candidates)
            }
            CompletionRequest::ArgsFields {
                probes,
                field_prefix,
                range,
            } => {
                let mut candidates = probes
                    .split(',')
                    .flat_map(|probe| {
                        workspace
                            .fields
                            .get(probe.trim())
                            .into_iter()
                            .flatten()
                            .cloned()
                    })
                    .collect::<Vec<_>>();
                candidates.extend(self.kernel_fields(&probes).await);
                field_completions(&field_prefix, range, candidates)
            }
        }
    }

    async fn kernel_targets(&self) -> Vec<TargetCandidate> {
        if let Some(targets) = self
            .cache
            .lock()
            .expect("probe metadata cache poisoned")
            .targets
            .clone()
        {
            return targets;
        }
        let targets = task::spawn_blocking(load_kernel_targets)
            .await
            .unwrap_or_default();
        self.cache
            .lock()
            .expect("probe metadata cache poisoned")
            .targets = Some(targets.clone());
        targets
    }

    async fn kernel_fields(&self, probes: &str) -> Vec<FieldCandidate> {
        let probe_names = probes
            .split(',')
            .map(str::trim)
            .filter(|probe| !probe.is_empty())
            .map(str::to_string)
            .collect::<Vec<_>>();
        let missing = {
            let cache = self.cache.lock().expect("probe metadata cache poisoned");
            probe_names
                .iter()
                .filter(|probe| !cache.fields.contains_key(*probe))
                .cloned()
                .collect::<Vec<_>>()
        };
        if !missing.is_empty() {
            let probes_to_load = missing.clone();
            let loaded = task::spawn_blocking(move || load_kernel_fields(&probes_to_load))
                .await
                .unwrap_or_default();
            let mut cache = self.cache.lock().expect("probe metadata cache poisoned");
            for probe in missing {
                cache.fields.insert(
                    probe.clone(),
                    loaded.get(&probe).cloned().unwrap_or_default(),
                );
            }
        }
        let cache = self.cache.lock().expect("probe metadata cache poisoned");
        probe_names
            .iter()
            .flat_map(|probe| cache.fields.get(probe).into_iter().flatten().cloned())
            .collect()
    }
}

pub(super) fn completion_request(
    snapshot: &DocumentSnapshot,
    position: Position,
) -> Option<CompletionRequest> {
    let text = snapshot.text.as_ref();
    let offset = snapshot.line_index.offset_for_position(text, position);
    args_fields_request(snapshot, offset).or_else(|| probe_targets_request(snapshot, offset))
}

fn probe_targets_request(snapshot: &DocumentSnapshot, offset: usize) -> Option<CompletionRequest> {
    let text = snapshot.text.as_ref();
    let offset = offset.min(text.len());
    let before = &text[..offset];
    let clause_start = before
        .rfind(['\n', ',', '{', '}'])
        .map_or(0, |index| index + 1);
    let whitespace = before[clause_start..]
        .len()
        .saturating_sub(before[clause_start..].trim_start().len());
    let start = clause_start + whitespace;
    let candidate = &text[start..offset];
    let (provider, target_prefix) = candidate.split_once(':')?;
    if !is_probe_provider(provider) || target_prefix.chars().any(char::is_whitespace) {
        return None;
    }

    let target_start = start + provider.len() + 1;
    let mut target_end = offset;
    while target_end < text.len() && is_probe_target_byte(text.as_bytes()[target_end]) {
        target_end += 1;
    }
    Some(CompletionRequest::ProbeTargets {
        provider: provider.to_string(),
        target_prefix: target_prefix.to_string(),
        range: snapshot
            .line_index
            .range_for_offsets(text, target_start, target_end),
    })
}

fn args_fields_request(snapshot: &DocumentSnapshot, offset: usize) -> Option<CompletionRequest> {
    let text = snapshot.text.as_ref();
    let offset = offset.min(text.len());
    let mut field_start = offset;
    while field_start > 0 && is_identifier_byte(text.as_bytes()[field_start - 1]) {
        field_start -= 1;
    }
    if field_start == 0 || text.as_bytes()[field_start - 1] != b'.' {
        return None;
    }
    let mut argument_start = field_start - 1;
    while argument_start > 0 && is_identifier_byte(text.as_bytes()[argument_start - 1]) {
        argument_start -= 1;
    }
    if &text[argument_start..field_start - 1] != "args" {
        return None;
    }
    let mut field_end = offset;
    while field_end < text.len() && is_identifier_byte(text.as_bytes()[field_end]) {
        field_end += 1;
    }
    let probes = enclosing_probes(snapshot, argument_start)?;
    Some(CompletionRequest::ArgsFields {
        probes,
        field_prefix: text[field_start..offset].to_string(),
        range: snapshot
            .line_index
            .range_for_offsets(text, field_start, field_end),
    })
}

fn enclosing_probes(snapshot: &DocumentSnapshot, offset: usize) -> Option<String> {
    let text = snapshot.text.as_ref();
    let tree = snapshot.tree.as_ref()?;
    let start = offset.saturating_sub(1);
    let node = tree
        .root_node()
        .descendant_for_byte_range(start, offset.max(start + 1).min(text.len()))?;
    let action_block = ancestor(node, "action_block")?;
    let probes = named_child(action_block, "probes_list")?;
    Some(text[probes.byte_range()].trim().to_string())
}

fn ancestor<'tree>(mut node: Node<'tree>, kind: &str) -> Option<Node<'tree>> {
    loop {
        if node.kind() == kind {
            return Some(node);
        }
        node = node.parent()?;
    }
}

fn named_child<'tree>(node: Node<'tree>, kind: &str) -> Option<Node<'tree>> {
    let mut cursor = node.walk();
    let child = node
        .named_children(&mut cursor)
        .find(|child| child.kind() == kind);
    child
}

fn walk(node: Node<'_>, visit: &mut impl FnMut(Node<'_>)) {
    visit(node);
    let mut cursor = node.walk();
    for child in node.named_children(&mut cursor) {
        walk(child, visit);
    }
}

fn is_probe_provider(provider: &str) -> bool {
    matches!(
        provider,
        "bench"
            | "fentry"
            | "fexit"
            | "hardware"
            | "interval"
            | "iter"
            | "kprobe"
            | "kretprobe"
            | "profile"
            | "rawtracepoint"
            | "software"
            | "test"
            | "tracepoint"
            | "uprobe"
            | "uretprobe"
            | "usdt"
            | "watchpoint"
    )
}

fn is_probe_target_byte(byte: u8) -> bool {
    byte.is_ascii_alphanumeric()
        || matches!(byte, b'_' | b'-' | b'.' | b'/' | b':' | b'*' | b'?' | b'+')
}

fn is_identifier_byte(byte: u8) -> bool {
    byte.is_ascii_alphanumeric() || byte == b'_'
}

fn target_completions(
    provider: &str,
    prefix: &str,
    range: Range,
    candidates: Vec<TargetCandidate>,
) -> CompletionList {
    let expected = format!("{provider}:");
    let mut seen = HashSet::new();
    let mut targets = candidates
        .into_iter()
        .filter_map(|candidate| {
            let target = candidate.full_name.strip_prefix(&expected)?.to_string();
            (target.starts_with(prefix) && seen.insert(target.clone()))
                .then_some((target, candidate.detail))
        })
        .collect::<Vec<_>>();
    targets.sort_by(|left, right| left.0.cmp(&right.0));
    let is_incomplete = targets.len() > MAX_ITEMS;
    targets.truncate(MAX_ITEMS);
    let items = targets
        .into_iter()
        .enumerate()
        .map(|(index, (target, detail))| CompletionItem {
            label: target.clone(),
            kind: Some(CompletionItemKind::EVENT),
            detail: Some(detail),
            filter_text: Some(target.clone()),
            sort_text: Some(format!("{index:06}")),
            text_edit: Some(CompletionTextEdit::Edit(TextEdit {
                range,
                new_text: target,
            })),
            ..CompletionItem::default()
        })
        .collect();
    CompletionList {
        is_incomplete,
        items,
    }
}

fn field_completions(
    prefix: &str,
    range: Range,
    candidates: Vec<FieldCandidate>,
) -> CompletionList {
    let mut seen = HashSet::new();
    let mut fields = candidates
        .into_iter()
        .filter(|field| field.name.starts_with(prefix) && seen.insert(field.name.clone()))
        .collect::<Vec<_>>();
    fields.sort_by(|left, right| left.name.cmp(&right.name));
    let is_incomplete = fields.len() > MAX_ITEMS;
    fields.truncate(MAX_ITEMS);
    let items = fields
        .into_iter()
        .enumerate()
        .map(|(index, field)| CompletionItem {
            label: field.name.clone(),
            kind: Some(CompletionItemKind::FIELD),
            detail: Some(field.declaration),
            documentation: Some(Documentation::String(format!("Probe: `{}`", field.probe))),
            sort_text: Some(format!("{index:06}")),
            text_edit: Some(CompletionTextEdit::Edit(TextEdit {
                range,
                new_text: field.name,
            })),
            ..CompletionItem::default()
        })
        .collect();
    CompletionList {
        is_incomplete,
        items,
    }
}

fn load_kernel_targets() -> Vec<TargetCandidate> {
    let mut targets = event_root()
        .map(|root| load_tracepoint_targets(&root))
        .unwrap_or_default();
    targets.extend(load_kprobe_targets());
    targets
}

fn event_root() -> Option<PathBuf> {
    [
        "/sys/kernel/events",
        "/sys/kernel/tracing/events",
        "/sys/kernel/debug/tracing/events",
    ]
    .into_iter()
    .map(PathBuf::from)
    .find(|path| fs::read_dir(path).is_ok())
}

fn load_tracepoint_targets(root: &Path) -> Vec<TargetCandidate> {
    let mut targets = Vec::new();
    let Ok(subsystems) = fs::read_dir(root) else {
        return targets;
    };
    for subsystem in subsystems.flatten() {
        let Ok(file_type) = subsystem.file_type() else {
            continue;
        };
        if !file_type.is_dir() {
            continue;
        }
        let subsystem_name = subsystem.file_name().to_string_lossy().into_owned();
        let Ok(events) = fs::read_dir(subsystem.path()) else {
            continue;
        };
        for event in events.flatten() {
            if !event.path().join("format").is_file() {
                continue;
            }
            let event_name = event.file_name().to_string_lossy().into_owned();
            targets.push(TargetCandidate {
                full_name: format!("tracepoint:{subsystem_name}:{event_name}"),
                detail: "kernel tracepoint".to_string(),
            });
        }
    }
    targets
}

fn load_kprobe_targets() -> Vec<TargetCandidate> {
    [
        "/sys/kernel/tracing/available_filter_functions",
        "/sys/kernel/debug/tracing/available_filter_functions",
    ]
    .into_iter()
    .find_map(|path| fs::read_to_string(path).ok())
    .map(|text| {
        text.lines()
            .filter_map(|line| line.split_whitespace().next())
            .map(|function| TargetCandidate {
                full_name: format!("kprobe:{function}"),
                detail: "kernel function".to_string(),
            })
            .collect()
    })
    .unwrap_or_default()
}

fn load_kernel_fields(probes: &[String]) -> HashMap<String, Vec<FieldCandidate>> {
    let Some(root) = event_root() else {
        return HashMap::new();
    };
    probes
        .iter()
        .filter_map(|probe| {
            let (_, target) = probe.split_once("tracepoint:")?;
            let (subsystem, event) = target.split_once(':')?;
            if !safe_component(subsystem) || !safe_component(event) {
                return None;
            }
            let format =
                fs::read_to_string(root.join(subsystem).join(event).join("format")).ok()?;
            Some((probe.clone(), parse_tracepoint_fields(probe, &format)))
        })
        .collect()
}

fn safe_component(component: &str) -> bool {
    !component.is_empty()
        && component
            .bytes()
            .all(|byte| byte.is_ascii_alphanumeric() || matches!(byte, b'_' | b'-' | b'.'))
}

fn parse_tracepoint_fields(probe: &str, format: &str) -> Vec<FieldCandidate> {
    format
        .lines()
        .filter_map(|line| line.trim().strip_prefix("field:"))
        .filter_map(|line| {
            line.split_once(';')
                .map(|(declaration, _)| declaration.trim())
        })
        .filter_map(|declaration| {
            let name = declaration_name(declaration)?;
            (!name.starts_with("common_")).then(|| FieldCandidate {
                probe: probe.to_string(),
                declaration: declaration.to_string(),
                name: name.to_string(),
            })
        })
        .collect()
}

fn declaration_name(declaration: &str) -> Option<&str> {
    let declaration = declaration.trim_end_matches(';').trim_end();
    let before_array = declaration
        .rfind('[')
        .map_or(declaration, |index| &declaration[..index]);
    let end = before_array
        .char_indices()
        .rev()
        .find(|(_, ch)| ch.is_ascii_alphanumeric() || *ch == '_')?
        .0
        + 1;
    let start = before_array[..end]
        .char_indices()
        .rev()
        .find(|(_, ch)| !ch.is_ascii_alphanumeric() && *ch != '_')
        .map_or(0, |(index, ch)| index + ch.len_utf8());
    (start < end).then_some(&before_array[start..end])
}

#[cfg(test)]
mod tests {
    use super::*;

    fn request(source: &str, marker: &str) -> CompletionRequest {
        let offset = source.find(marker).unwrap();
        let text = source.replacen(marker, "", 1);
        let (snapshot, _) = DocumentSnapshot::analyze(text, Some(1));
        let position = snapshot
            .line_index
            .position_for_offset(&snapshot.text, offset);
        completion_request(&snapshot, position).unwrap()
    }

    #[test]
    fn extracts_dynamic_completion_requests() {
        let target = request("tracepoint:sysc| { exit(); }", "|");
        assert!(matches!(
            target,
            CompletionRequest::ProbeTargets {
                provider,
                target_prefix,
                ..
            } if provider == "tracepoint" && target_prefix == "sysc"
        ));

        let fields = request(
            "tracepoint:syscalls:sys_enter_openat { print(args.file|); }",
            "|",
        );
        assert!(matches!(
            fields,
            CompletionRequest::ArgsFields {
                probes,
                field_prefix,
                ..
            } if probes == "tracepoint:syscalls:sys_enter_openat" && field_prefix == "file"
        ));
    }

    #[test]
    fn extracts_workspace_targets_and_fields() {
        let source = "tracepoint:syscalls:sys_enter_openat { print(args.filename); }";
        let (snapshot, _) = DocumentSnapshot::analyze(source.to_string(), Some(1));
        let mut metadata = WorkspaceMetadata::default();
        metadata.add_snapshot(&snapshot);
        assert!(metadata
            .targets
            .iter()
            .any(|target| { target.full_name == "tracepoint:syscalls:sys_enter_openat" }));
        assert_eq!(
            metadata.fields["tracepoint:syscalls:sys_enter_openat"][0].name,
            "filename"
        );
    }

    #[test]
    fn reads_tracepoint_metadata_without_external_commands() {
        let temp = tempfile::tempdir().unwrap();
        let event = temp.path().join("syscalls/sys_enter_openat");
        fs::create_dir_all(&event).unwrap();
        fs::write(
            event.join("format"),
            "format:\n\tfield:unsigned short common_type; offset:0;\n\tfield:const char * filename; offset:8;\n\tfield:int flags; offset:16;\n",
        )
        .unwrap();

        let targets = load_tracepoint_targets(temp.path());
        assert_eq!(targets[0].full_name, "tracepoint:syscalls:sys_enter_openat");
        let fields = parse_tracepoint_fields(
            "tracepoint:syscalls:sys_enter_openat",
            &fs::read_to_string(event.join("format")).unwrap(),
        );
        assert_eq!(
            fields
                .iter()
                .map(|field| field.name.as_str())
                .collect::<Vec<_>>(),
            ["filename", "flags"]
        );
    }

    #[test]
    fn static_targets_cover_unprivileged_environments() {
        let completions = target_completions(
            "hardware",
            "cpu",
            Range::default(),
            STATIC_TARGETS
                .iter()
                .map(|(full_name, detail)| TargetCandidate {
                    full_name: (*full_name).to_string(),
                    detail: (*detail).to_string(),
                })
                .collect(),
        );
        assert_eq!(completions.items[0].label, "cpu-cycles");
    }
}
