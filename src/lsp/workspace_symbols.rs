use super::symbols::GlobalSymbolKind;
use super::workspace::WorkspaceIndex;
use tower_lsp::lsp_types::{Location, SymbolInformation, SymbolKind};

pub(super) fn search(
    workspace: &WorkspaceIndex,
    query: &str,
    trusted: bool,
) -> Vec<SymbolInformation> {
    // Trust is checked even if the index still contains snapshots from a trusted session.
    if !trusted {
        return Vec::new();
    }
    let query = query.trim().to_lowercase();
    let mut results = Vec::new();
    for (uri, snapshot) in workspace.symbol_snapshots() {
        let Some(index) = snapshot.symbols.as_ref() else {
            continue;
        };
        let container = Some(uri.path().to_string());
        let mut push = |name: String, kind, range| {
            if name.to_lowercase().contains(&query) {
                #[allow(deprecated)]
                results.push(SymbolInformation {
                    name,
                    kind,
                    tags: None,
                    deprecated: None,
                    location: Location::new(uri.clone(), range),
                    container_name: container.clone(),
                });
            }
        };
        for (target, range) in index.global_definitions() {
            let (name, kind) = match target.kind {
                GlobalSymbolKind::Map => (format!("@{}", target.name), SymbolKind::VARIABLE),
                GlobalSymbolKind::Macro => (target.name, SymbolKind::FUNCTION),
            };
            push(name, kind, super::lsp_range(&snapshot, range));
        }
        if let Some(tree) = snapshot.tree.as_ref() {
            for symbol in super::document_symbols(&snapshot.text, tree, &snapshot.line_index) {
                if symbol.kind == SymbolKind::EVENT {
                    push(symbol.name, symbol.kind, symbol.selection_range);
                }
            }
        }
    }
    results.sort_by(|left, right| {
        left.name
            .cmp(&right.name)
            .then_with(|| left.location.uri.as_str().cmp(right.location.uri.as_str()))
            .then_with(|| left.location.range.start.cmp(&right.location.range.start))
            .then_with(|| left.location.range.end.cmp(&right.location.range.end))
    });
    results
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::lsp::snapshot::DocumentSnapshot;
    use std::fs;
    use std::path::Path;
    use tempfile::tempdir;
    use tower_lsp::lsp_types::{Position, Url};

    fn uri(path: &Path) -> Url {
        Url::from_file_path(path).unwrap()
    }

    fn overlay(index: &mut WorkspaceIndex, path: &Path, text: &str) {
        index.set_overlay(uri(path), DocumentSnapshot::analyze(text.into(), Some(1)).0);
    }

    fn names(index: &WorkspaceIndex, query: &str) -> Vec<String> {
        search(index, query, true)
            .into_iter()
            .map(|symbol| symbol.name)
            .collect()
    }

    #[test]
    fn queries_return_definitions_with_stable_kinds_and_utf16_locations() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("symbols.bt");
        let text = concat!(
            "macro Measure(value, @local) { @private = value; value }\r\n",
            "BEGIN { printf(\"😀\"); @Count = 1; @Count++; print(@Count); ",
            "print(@read_only); $scratch = Measure(1, @Count); missing(); @ = 1; }\r\n",
            "END { print(@Count); }\r\n",
        );
        fs::write(&path, text).unwrap();
        let mut index = WorkspaceIndex::default();
        index.scan(&[dir.path().into()]);
        assert_eq!(
            names(&index, ""),
            ["@", "@Count", "BEGIN", "END", "Measure"]
        );
        assert_eq!(names(&index, "  coUNT  "), ["@Count"]);
        assert_eq!(names(&index, "@c"), ["@Count"]);
        assert_eq!(names(&index, "MEAS"), ["Measure"]);
        assert!(names(&index, "not_present").is_empty());
        let all = search(&index, "", true);
        assert_eq!(all, search(&index, "   ", true));
        for (name, kind) in [
            ("@Count", SymbolKind::VARIABLE),
            ("Measure", SymbolKind::FUNCTION),
            ("BEGIN", SymbolKind::EVENT),
        ] {
            let symbol = all.iter().find(|symbol| symbol.name == name).unwrap();
            assert_eq!(symbol.kind, kind);
            assert_eq!(symbol.location.uri, uri(&path));
            let snapshot = index.snapshot(&symbol.location.uri).unwrap();
            let start = snapshot
                .line_index
                .offset_for_position(text, symbol.location.range.start);
            let end = snapshot
                .line_index
                .offset_for_position(text, symbol.location.range.end);
            assert_eq!(&text[start..end], name);
        }
        assert_eq!(all[1].location.range.start, Position::new(1, 22));
        assert!(
            search(&index, "", false).is_empty(),
            "cached symbols require trust"
        );
    }

    #[test]
    fn overlays_disk_updates_and_root_changes_share_the_index_lifecycle() {
        let dir = tempdir().unwrap();
        let first = dir.path().join("first");
        let second = dir.path().join("second");
        fs::create_dir(&first).unwrap();
        fs::create_dir(&second).unwrap();
        let a = first.join("a.bt");
        let b = second.join("b.bt");
        fs::write(&a, "BEGIN { @same = 1; @same++; }").unwrap();
        fs::write(&b, "BEGIN { @same = 2; }").unwrap();
        let mut index = WorkspaceIndex::default();
        index.scan(&[first.clone(), second.clone(), first.clone()]);
        let results = search(&index, "@same", true);
        assert_eq!(results.len(), 2);
        assert_eq!(results[0].location.uri, uri(&a));
        assert_eq!(results[1].location.uri, uri(&b));
        overlay(&mut index, &a, "BEGIN { @unsaved = 1; }");
        fs::write(&a, "BEGIN { @saved = 1; }").unwrap();
        index.refresh(&uri(&a));
        assert_eq!(names(&index, "@"), ["@same", "@unsaved"]);
        index.remove_overlay(&uri(&a));
        assert_eq!(names(&index, "@"), ["@same", "@saved"]);
        fs::remove_file(&a).unwrap();
        index.remove_disk(&uri(&a));
        assert_eq!(names(&index, "@"), ["@same"]);
        overlay(&mut index, &a, "BEGIN { @new = 1; }");
        overlay(
            &mut index,
            &dir.path().join("outside.bt"),
            "BEGIN { @outside = 1; }",
        );
        overlay(&mut index, &b, "BEGIN { @second = 1; }");
        index.scan(&[first]);
        assert_eq!(names(&index, "@"), ["@new"]);
        index.scan(&[second]);
        assert_eq!(names(&index, "@"), ["@second"]);
        index.scan(&[]);
        assert!(names(&index, "").is_empty());
    }

    #[test]
    fn ignore_rules_apply_to_scans_refreshes_and_new_buffers() {
        let dir = tempdir().unwrap();
        let root = dir.path();
        fs::create_dir(root.join(".git")).unwrap();
        fs::create_dir(root.join("excluded")).unwrap();
        fs::write(root.join(".gitignore"), "ignored*.bt\nexcluded/\n").unwrap();
        fs::write(root.join(".ignore"), "local.bt\n!ignored_keep.bt\n").unwrap();
        let mut index = WorkspaceIndex::default();
        for file in [
            "ignored.bt",
            "local.bt",
            ".hidden.bt",
            "excluded/a.bt",
            "ignored_keep.bt",
            "keep.bt",
        ] {
            fs::write(root.join(file), "BEGIN { @symbol = 1; }").unwrap();
        }
        index.scan(&[root.into()]);
        assert_eq!(search(&index, "@symbol", true).len(), 2);
        for file in [
            "ignored.bt",
            "local.bt",
            ".hidden.bt",
            "excluded/a.bt",
            "ignored_new.bt",
        ] {
            index.refresh(&uri(&root.join(file)));
            overlay(&mut index, &root.join(file), "BEGIN { @leak = 1; }");
        }
        assert!(names(&index, "leak").is_empty());
        overlay(&mut index, &root.join("new.bt"), "BEGIN { @new = 1; }");
        assert_eq!(names(&index, "@new"), ["@new"]);
        fs::write(root.join(".ignore"), "*.bt\n").unwrap();
        index.refresh(&uri(&root.join(".ignore")));
        assert!(names(&index, "").is_empty());
        fs::remove_file(root.join(".ignore")).unwrap();
        index.remove_disk(&uri(&root.join(".ignore")));
        assert_eq!(names(&index, "@new"), ["@new"]);
    }

    #[test]
    fn invalid_overlay_hides_stale_definitions_until_close() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("main.bt");
        fs::write(&path, "BEGIN { @valid = 1; }").unwrap();
        let mut index = WorkspaceIndex::default();
        index.scan(&[dir.path().into()]);
        overlay(&mut index, &path, "BEGIN { @invalid =");
        assert!(names(&index, "").is_empty());
        index.remove_overlay(&uri(&path));
        assert_eq!(names(&index, "@valid"), ["@valid"]);
    }

    #[tokio::test]
    async fn restricted_initialization_never_indexes_workspace_files() {
        use crate::lsp::{Backend, DocumentStore, ServerState};
        use std::sync::{Arc, Mutex};
        use tower_lsp::lsp_types::{InitializeParams, WorkspaceSymbolParams};
        use tower_lsp::{LanguageServer, LspService};

        let dir = tempdir().unwrap();
        let path = dir.path().join("secret.bt");
        fs::write(&path, "BEGIN { @secret = 1; }").unwrap();
        let workspace = Arc::new(Mutex::new(WorkspaceIndex::default()));
        let (service, _socket) = LspService::new(|client| Backend {
            client,
            docs: Arc::new(Mutex::new(DocumentStore::default())),
            state: Arc::new(Mutex::new(ServerState::default())),
            workspace: Arc::clone(&workspace),
            probes: Default::default(),
        });
        service
            .inner()
            .initialize(InitializeParams {
                root_uri: Some(uri(dir.path())),
                initialization_options: Some(serde_json::json!({"btfmt": {"trusted": false}})),
                ..Default::default()
            })
            .await
            .unwrap();
        assert!(workspace.lock().unwrap().snapshot(&uri(&path)).is_none());
        // Exercise the handler's trust check with a populated cache as well.
        workspace.lock().unwrap().scan(&[dir.path().into()]);
        assert!(workspace.lock().unwrap().snapshot(&uri(&path)).is_some());
        let result = service
            .inner()
            .symbol(WorkspaceSymbolParams {
                query: String::new(),
                ..Default::default()
            })
            .await
            .unwrap()
            .unwrap();
        assert!(result.is_empty());
    }

    #[test]
    fn repeated_macro_definitions_deduplicate_but_distinct_probe_blocks_remain() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("main.bt");
        fs::write(
            &path,
            concat!(
                "macro helper() { 1 }\nmacro helper() { 2 }\n",
                "BEGIN { helper(); }\nBEGIN { helper(); }\n",
            ),
        )
        .unwrap();
        let mut index = WorkspaceIndex::default();
        index.scan(&[dir.path().into()]);
        let macros = search(&index, "helper", true);
        assert_eq!(macros.len(), 1);
        assert_eq!(macros[0].location.range.start, Position::new(0, 6));
        let probes = search(&index, "BEGIN", true);
        assert_eq!(probes.len(), 2);
        assert_eq!(probes[0].location.range.start.line, 2);
        assert_eq!(probes[1].location.range.start.line, 3);
    }

    #[cfg(unix)]
    #[test]
    fn symlinks_cannot_expose_files_outside_the_workspace() {
        let dir = tempdir().unwrap();
        let root = dir.path().join("root");
        fs::create_dir(&root).unwrap();
        let outside = dir.path().join("outside.bt");
        fs::write(&outside, "BEGIN { @secret = 1; }").unwrap();
        let link = root.join("link.bt");
        std::os::unix::fs::symlink(outside, &link).unwrap();
        let mut index = WorkspaceIndex::default();
        index.scan(&[root]);
        index.refresh(&uri(&link));
        overlay(&mut index, &link, "BEGIN { @secret = 2; }");
        assert!(names(&index, "").is_empty());
    }
}
