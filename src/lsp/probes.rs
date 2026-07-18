use super::snapshot::DocumentSnapshot;
use std::collections::{HashMap, HashSet};
use std::sync::Arc;
use std::time::{Duration, Instant};
use tokio::process::Command;
use tokio::sync::Mutex;
use tokio::time::timeout;
use tower_lsp::lsp_types::{
    CompletionItem, CompletionItemKind, CompletionList, CompletionTextEdit, Documentation,
    Position, Range, TextEdit,
};
use tree_sitter::Node;

const QUERY_TIMEOUT: Duration = Duration::from_secs(3);
const SUCCESS_TTL: Duration = Duration::from_secs(30);
const FAILURE_TTL: Duration = Duration::from_secs(5);
const MAX_OUTPUT_BYTES: usize = 16 * 1024 * 1024;
const MAX_ITEMS: usize = 200;

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
enum QueryKind {
    Targets,
    Fields,
}

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
struct QueryKey {
    kind: QueryKind,
    pattern: String,
}

#[derive(Debug, Clone)]
enum QueryResult {
    Success(Arc<str>),
    Failure(Arc<str>),
}

#[derive(Debug, Clone)]
struct CacheEntry {
    created: Instant,
    result: QueryResult,
}

impl CacheEntry {
    fn is_fresh(&self) -> bool {
        let ttl = match self.result {
            QueryResult::Success(_) => SUCCESS_TTL,
            QueryResult::Failure(_) => FAILURE_TTL,
        };
        self.created.elapsed() < ttl
    }
}

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

impl CompletionRequest {
    fn query_key(&self) -> QueryKey {
        match self {
            Self::ProbeTargets {
                provider,
                target_prefix,
                ..
            } => QueryKey {
                kind: QueryKind::Targets,
                pattern: format!("{provider}:{target_prefix}*"),
            },
            Self::ArgsFields { probes, .. } => QueryKey {
                kind: QueryKind::Fields,
                pattern: probes.clone(),
            },
        }
    }
}

#[derive(Debug, Default)]
pub(super) struct ProbeCatalog {
    cache: Mutex<HashMap<QueryKey, CacheEntry>>,
    reported_errors: Mutex<HashSet<Arc<str>>>,
}

#[derive(Debug)]
pub(super) struct CompletionOutcome {
    pub list: CompletionList,
    pub warning: Option<Arc<str>>,
}

impl ProbeCatalog {
    pub(super) async fn complete(&self, request: CompletionRequest) -> CompletionOutcome {
        let key = request.query_key();
        match self.query(key).await {
            QueryResult::Success(output) => CompletionOutcome {
                list: completion_from_output(&request, &output),
                warning: None,
            },
            QueryResult::Failure(error) => {
                let warning = self
                    .reported_errors
                    .lock()
                    .await
                    .insert(Arc::clone(&error))
                    .then_some(error);
                CompletionOutcome {
                    list: CompletionList {
                        is_incomplete: false,
                        items: Vec::new(),
                    },
                    warning,
                }
            }
        }
    }

    async fn query(&self, key: QueryKey) -> QueryResult {
        if let Some(result) = self
            .cache
            .lock()
            .await
            .get(&key)
            .filter(|entry| entry.is_fresh())
            .map(|entry| entry.result.clone())
        {
            return result;
        }

        let flag = match key.kind {
            QueryKind::Targets => "-l",
            QueryKind::Fields => "-lv",
        };
        let mut command = Command::new("bpftrace");
        command.arg(flag).arg(&key.pattern).kill_on_drop(true);
        let result = match timeout(QUERY_TIMEOUT, command.output()).await {
            Err(_) => QueryResult::Failure(Arc::from(format!(
                "bpftrace query timed out after {} seconds",
                QUERY_TIMEOUT.as_secs()
            ))),
            Ok(Err(error)) => QueryResult::Failure(Arc::from(format!(
                "cannot run bpftrace for completion: {error}"
            ))),
            Ok(Ok(output)) if !output.status.success() => {
                let stderr = String::from_utf8_lossy(&output.stderr);
                QueryResult::Failure(Arc::from(format!(
                    "bpftrace completion query failed: {}",
                    stderr.trim()
                )))
            }
            Ok(Ok(output)) if output.stdout.len() > MAX_OUTPUT_BYTES => QueryResult::Failure(
                Arc::from("bpftrace completion query returned too much data"),
            ),
            Ok(Ok(output)) => QueryResult::Success(Arc::from(
                String::from_utf8_lossy(&output.stdout).into_owned(),
            )),
        };
        self.cache.lock().await.insert(
            key,
            CacheEntry {
                created: Instant::now(),
                result: result.clone(),
            },
        );
        result
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
    if !is_probe_provider(provider)
        || target_prefix.chars().any(char::is_whitespace)
        || (!allows_empty_target(provider) && target_prefix.len() < 2)
    {
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
    let mut node = tree
        .root_node()
        .descendant_for_byte_range(start, offset.max(start + 1).min(text.len()))?;
    loop {
        if node.kind() == "action_block" {
            let probes = named_child(node, "probes_list")?;
            return Some(text[probes.byte_range()].trim().to_string());
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

fn allows_empty_target(provider: &str) -> bool {
    matches!(
        provider,
        "tracepoint" | "rawtracepoint" | "hardware" | "software"
    )
}

fn is_probe_target_byte(byte: u8) -> bool {
    byte.is_ascii_alphanumeric()
        || matches!(byte, b'_' | b'-' | b'.' | b'/' | b':' | b'*' | b'?' | b'+')
}

fn is_identifier_byte(byte: u8) -> bool {
    byte.is_ascii_alphanumeric() || byte == b'_'
}

fn completion_from_output(request: &CompletionRequest, output: &str) -> CompletionList {
    match request {
        CompletionRequest::ProbeTargets {
            provider,
            target_prefix,
            range,
        } => target_completions(provider, target_prefix, *range, output),
        CompletionRequest::ArgsFields {
            field_prefix,
            range,
            ..
        } => field_completions(field_prefix, *range, output),
    }
}

fn target_completions(provider: &str, prefix: &str, range: Range, output: &str) -> CompletionList {
    let expected = format!("{provider}:");
    let mut targets: Vec<_> = output
        .lines()
        .map(str::trim)
        .filter_map(|line| line.strip_prefix(&expected))
        .filter(|target| target.starts_with(prefix))
        .collect();
    targets.sort_unstable();
    targets.dedup();
    let is_incomplete = targets.len() > MAX_ITEMS;
    targets.truncate(MAX_ITEMS);
    let items = targets
        .into_iter()
        .enumerate()
        .map(|(index, target)| CompletionItem {
            label: target.to_string(),
            kind: Some(CompletionItemKind::EVENT),
            detail: Some(format!("{provider}:{target}")),
            filter_text: Some(target.to_string()),
            sort_text: Some(format!("{index:06}")),
            text_edit: Some(CompletionTextEdit::Edit(TextEdit {
                range,
                new_text: target.to_string(),
            })),
            ..CompletionItem::default()
        })
        .collect();
    CompletionList {
        is_incomplete,
        items,
    }
}

#[derive(Debug, PartialEq, Eq)]
struct ProbeField<'a> {
    probe: &'a str,
    declaration: &'a str,
    name: &'a str,
}

fn field_completions(prefix: &str, range: Range, output: &str) -> CompletionList {
    let mut fields = parse_fields(output);
    fields.retain(|field| field.name.starts_with(prefix));
    fields.sort_by_key(|field| field.name);
    fields.dedup_by_key(|field| field.name);
    let is_incomplete = fields.len() > MAX_ITEMS;
    fields.truncate(MAX_ITEMS);
    let items = fields
        .into_iter()
        .enumerate()
        .map(|(index, field)| CompletionItem {
            label: field.name.to_string(),
            kind: Some(CompletionItemKind::FIELD),
            detail: Some(field.declaration.to_string()),
            documentation: Some(Documentation::String(format!("Probe: `{}`", field.probe))),
            sort_text: Some(format!("{index:06}")),
            text_edit: Some(CompletionTextEdit::Edit(TextEdit {
                range,
                new_text: field.name.to_string(),
            })),
            ..CompletionItem::default()
        })
        .collect();
    CompletionList {
        is_incomplete,
        items,
    }
}

fn parse_fields(output: &str) -> Vec<ProbeField<'_>> {
    let mut probe = "";
    let mut fields = Vec::new();
    for line in output.lines() {
        let declaration = line.trim();
        if declaration.is_empty() {
            continue;
        }
        if !line.starts_with(char::is_whitespace) {
            probe = declaration;
            continue;
        }
        let Some(name) = declaration_name(declaration) else {
            continue;
        };
        fields.push(ProbeField {
            probe,
            declaration,
            name,
        });
    }
    fields
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
    fn parses_target_and_field_output() {
        let range = Range::default();
        let targets = target_completions(
            "tracepoint",
            "sysc",
            range,
            "tracepoint:sched:sched_switch\ntracepoint:syscalls:sys_enter_openat\n",
        );
        assert_eq!(targets.items.len(), 1);
        assert_eq!(targets.items[0].label, "syscalls:sys_enter_openat");

        let output = "tracepoint:syscalls:sys_enter_openat\n    int __syscall_nr\n    const char * filename\n    char comm[16]\n";
        let fields = parse_fields(output);
        assert_eq!(
            fields,
            vec![
                ProbeField {
                    probe: "tracepoint:syscalls:sys_enter_openat",
                    declaration: "int __syscall_nr",
                    name: "__syscall_nr",
                },
                ProbeField {
                    probe: "tracepoint:syscalls:sys_enter_openat",
                    declaration: "const char * filename",
                    name: "filename",
                },
                ProbeField {
                    probe: "tracepoint:syscalls:sys_enter_openat",
                    declaration: "char comm[16]",
                    name: "comm",
                },
            ]
        );
    }
}
