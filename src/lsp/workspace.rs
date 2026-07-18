use super::snapshot::DocumentSnapshot;
use ignore::WalkBuilder;
use std::collections::{HashMap, HashSet, VecDeque};
use std::fs;
use std::path::{Path, PathBuf};
use std::sync::Arc;
use tower_lsp::lsp_types::Url;
use tree_sitter::Node;

const MAX_FILES: usize = 5_000;
const MAX_FILE_BYTES: u64 = 2 * 1024 * 1024;
const MAX_TOTAL_BYTES: u64 = 50 * 1024 * 1024;

#[derive(Debug, Default)]
pub(super) struct WorkspaceIndex {
    roots: Vec<PathBuf>,
    disk: HashMap<Url, Arc<DocumentSnapshot>>,
    overlays: HashMap<Url, Arc<DocumentSnapshot>>,
    forward: HashMap<Url, Vec<Url>>,
    reverse: HashMap<Url, HashSet<Url>>,
    ready: bool,
    incomplete: bool,
}

impl WorkspaceIndex {
    pub(super) fn scan(&mut self, roots: &[PathBuf]) {
        self.roots = roots
            .iter()
            .filter_map(|root| fs::canonicalize(root).ok())
            .collect();
        self.disk.clear();
        self.incomplete = false;
        let mut total_bytes = 0u64;
        'roots: for root in &self.roots {
            let mut builder = WalkBuilder::new(root);
            builder
                .standard_filters(true)
                .hidden(true)
                .follow_links(false);
            for entry in builder.build().filter_map(Result::ok) {
                let path = entry.path();
                if path.extension().and_then(|ext| ext.to_str()) != Some("bt") {
                    continue;
                }
                if self.disk.len() >= MAX_FILES {
                    self.incomplete = true;
                    break 'roots;
                }
                let Ok(metadata) = entry.metadata() else {
                    continue;
                };
                if metadata.len() > MAX_FILE_BYTES
                    || total_bytes.saturating_add(metadata.len()) > MAX_TOTAL_BYTES
                {
                    self.incomplete = true;
                    continue;
                }
                let Ok(text) = fs::read_to_string(path) else {
                    continue;
                };
                let Ok(uri) = Url::from_file_path(path) else {
                    continue;
                };
                let (snapshot, _) = DocumentSnapshot::analyze(text, None);
                total_bytes += metadata.len();
                self.disk.insert(uri, snapshot);
            }
        }
        self.ready = true;
        self.rebuild_graph();
    }

    pub(super) fn set_overlay(&mut self, uri: Url, snapshot: Arc<DocumentSnapshot>) {
        self.overlays.insert(uri, snapshot);
        self.rebuild_graph();
    }

    pub(super) fn remove_overlay(&mut self, uri: &Url) {
        self.overlays.remove(uri);
        self.rebuild_graph();
    }

    pub(super) fn refresh(&mut self, uri: &Url) {
        let Ok(path) = uri.to_file_path() else {
            return;
        };
        if !self.path_is_in_workspace(&path) {
            return;
        }
        match fs::read_to_string(&path) {
            Ok(text) if path.extension().and_then(|ext| ext.to_str()) == Some("bt") => {
                let (snapshot, _) = DocumentSnapshot::analyze(text, None);
                self.disk.insert(uri.clone(), snapshot);
            }
            _ => {
                self.disk.remove(uri);
            }
        }
        self.rebuild_graph();
    }

    pub(super) fn remove_disk(&mut self, uri: &Url) {
        self.disk.remove(uri);
        self.rebuild_graph();
    }

    pub(super) fn snapshot(&self, uri: &Url) -> Option<Arc<DocumentSnapshot>> {
        self.overlays
            .get(uri)
            .or_else(|| self.disk.get(uri))
            .map(Arc::clone)
    }

    pub(super) fn program_files(&self, uri: &Url) -> Vec<Url> {
        if !self.ready {
            return Vec::new();
        }
        let roots = self.reverse_roots(uri);
        let starts = if roots.is_empty() {
            vec![uri.clone()]
        } else {
            roots
        };
        let mut files = HashSet::new();
        for start in starts {
            self.collect_forward(&start, &mut files);
        }
        let mut files: Vec<_> = files.into_iter().collect();
        files.sort_by(|left, right| left.as_str().cmp(right.as_str()));
        files
    }

    pub(super) fn ready_for_rename(&self) -> bool {
        self.ready && !self.incomplete
    }

    pub(super) fn is_incomplete(&self) -> bool {
        !self.ready || self.incomplete
    }

    fn rebuild_graph(&mut self) {
        self.forward.clear();
        self.reverse.clear();
        let uris: Vec<_> = self
            .disk
            .keys()
            .chain(self.overlays.keys())
            .cloned()
            .collect();
        for uri in uris {
            let Some(snapshot) = self.snapshot(&uri) else {
                continue;
            };
            let imports = self.resolve_imports(&uri, &snapshot);
            for target in &imports {
                self.reverse
                    .entry(target.clone())
                    .or_default()
                    .insert(uri.clone());
            }
            self.forward.insert(uri, imports);
        }
    }

    fn resolve_imports(&self, importer: &Url, snapshot: &DocumentSnapshot) -> Vec<Url> {
        let Ok(importer_path) = importer.to_file_path() else {
            return Vec::new();
        };
        let Some(parent) = importer_path.parent() else {
            return Vec::new();
        };
        let Some(tree) = snapshot.tree.as_ref() else {
            return Vec::new();
        };
        let mut specs = Vec::new();
        collect_imports(tree.root_node(), &snapshot.text, &mut specs);
        let mut targets = HashSet::new();
        for spec in specs {
            let target = parent.join(spec);
            let Ok(target) = fs::canonicalize(target) else {
                continue;
            };
            if !self.path_is_in_workspace(&target) {
                continue;
            }
            if target.is_dir() {
                for uri in self.disk.keys().chain(self.overlays.keys()) {
                    let Ok(path) = uri.to_file_path() else {
                        continue;
                    };
                    if path.parent() == Some(target.as_path()) {
                        targets.insert(uri.clone());
                    }
                }
            } else if target.extension().and_then(|ext| ext.to_str()) == Some("bt") {
                if let Ok(uri) = Url::from_file_path(target) {
                    if self.disk.contains_key(&uri) || self.overlays.contains_key(&uri) {
                        targets.insert(uri);
                    }
                }
            }
        }
        let mut targets: Vec<_> = targets.into_iter().collect();
        targets.sort_by(|left, right| left.as_str().cmp(right.as_str()));
        targets
    }

    fn reverse_roots(&self, uri: &Url) -> Vec<Url> {
        let mut queue = VecDeque::from([uri.clone()]);
        let mut seen = HashSet::new();
        let mut roots = HashSet::new();
        while let Some(current) = queue.pop_front() {
            if !seen.insert(current.clone()) {
                continue;
            }
            match self
                .reverse
                .get(&current)
                .filter(|parents| !parents.is_empty())
            {
                Some(parents) => queue.extend(parents.iter().cloned()),
                None => {
                    roots.insert(current);
                }
            }
        }
        roots.into_iter().collect()
    }

    fn collect_forward(&self, uri: &Url, files: &mut HashSet<Url>) {
        if !files.insert(uri.clone()) {
            return;
        }
        if let Some(imports) = self.forward.get(uri) {
            for target in imports {
                self.collect_forward(target, files);
            }
        }
    }

    fn path_is_in_workspace(&self, path: &Path) -> bool {
        fs::canonicalize(path)
            .ok()
            .is_some_and(|path| self.roots.iter().any(|root| path.starts_with(root)))
    }
}

fn collect_imports(node: Node<'_>, source: &str, imports: &mut Vec<PathBuf>) {
    if node.kind() == "import_statement" {
        let mut cursor = node.walk();
        if let Some(literal) = node
            .named_children(&mut cursor)
            .find(|child| child.kind() == "string_literal")
        {
            if let Some(path) = decode_string_literal(source, literal) {
                imports.push(PathBuf::from(path));
            }
        }
        return;
    }
    let mut cursor = node.walk();
    for child in node.named_children(&mut cursor) {
        collect_imports(child, source, imports);
    }
}

fn decode_string_literal(source: &str, node: Node<'_>) -> Option<String> {
    let text = source.get(node.byte_range())?;
    let text = text.strip_prefix('"')?.strip_suffix('"')?;
    let mut decoded = String::new();
    let mut chars = text.chars();
    while let Some(ch) = chars.next() {
        if ch != '\\' {
            decoded.push(ch);
            continue;
        }
        let escaped = chars.next()?;
        decoded.push(match escaped {
            'n' => '\n',
            'r' => '\r',
            't' => '\t',
            '\\' => '\\',
            '"' => '"',
            other => other,
        });
    }
    Some(decoded)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::lsp::symbols::GlobalSymbolKind;
    use tempfile::tempdir;

    #[test]
    fn workspace_graph_respects_import_closures_and_boundaries() {
        let dir = tempdir().unwrap();
        let root = dir.path();
        fs::create_dir(root.join("lib")).unwrap();
        fs::write(
            root.join("root.bt"),
            "import \"lib/helper.bt\";\nBEGIN { helper(); }\n",
        )
        .unwrap();
        fs::write(root.join("lib/helper.bt"), "macro helper() { 1 }\n").unwrap();
        fs::write(root.join("independent.bt"), "BEGIN { @same = 1; }\n").unwrap();

        let mut index = WorkspaceIndex::default();
        index.scan(&[root.to_path_buf()]);
        let root_uri = Url::from_file_path(root.join("root.bt")).unwrap();
        let helper_uri = Url::from_file_path(root.join("lib/helper.bt")).unwrap();
        let independent_uri = Url::from_file_path(root.join("independent.bt")).unwrap();

        assert_eq!(
            index.program_files(&root_uri),
            vec![helper_uri.clone(), root_uri.clone()]
        );
        assert_eq!(
            index.program_files(&helper_uri),
            vec![
                helper_uri,
                Url::from_file_path(root.join("root.bt")).unwrap()
            ]
        );
        assert_eq!(index.program_files(&independent_uri), vec![independent_uri]);
        assert!(index.ready_for_rename());

        let root_snapshot = index.snapshot(&root_uri).unwrap();
        let call_offset = root_snapshot.text.find("helper()").unwrap();
        let target = root_snapshot
            .symbols
            .as_ref()
            .unwrap()
            .global_target_at(call_offset)
            .unwrap();
        assert_eq!(target.kind, GlobalSymbolKind::Macro);
        let occurrences: Vec<_> = index
            .program_files(&root_uri)
            .into_iter()
            .flat_map(|uri| {
                index
                    .snapshot(&uri)
                    .unwrap()
                    .symbols
                    .as_ref()
                    .unwrap()
                    .global_occurrences(&target)
            })
            .collect();
        assert_eq!(occurrences.len(), 2);
        assert_eq!(occurrences.iter().filter(|item| item.definition).count(), 1);
    }
}
