mod catalog;
mod completion;
mod probes;
mod snapshot;
mod symbols;
mod tracepoint_catalog;
mod workspace;

use self::snapshot::DocumentSnapshot;
use self::symbols::{
    AccessKind, ByteRange, GlobalSymbolKind, GlobalSymbolOccurrence, GlobalSymbolTarget,
    SymbolIndex,
};
use self::workspace::WorkspaceIndex;
use crate::config::{load_from_base, Config};
use crate::format::format_source;
use crate::text::identifier_at_position_with_index;
use anyhow::Result as AnyResult;
use serde_json::Value;
use std::collections::HashMap;
use std::path::{Path, PathBuf};
use std::sync::{Arc, Mutex};
use tower_lsp::jsonrpc::{Error, Result};
use tower_lsp::lsp_types::*;
use tower_lsp::{async_trait, Client, LanguageServer, LspService, Server};
use tree_sitter::Tree;

#[derive(Debug, Default)]
struct DocumentStore {
    documents: HashMap<Url, Arc<DocumentSnapshot>>,
}

impl DocumentStore {
    fn can_accept(&self, uri: &Url, version: i32, require_existing: bool) -> bool {
        match self.documents.get(uri) {
            Some(current) => current.version.is_some_and(|current| version > current),
            None => !require_existing,
        }
    }

    fn insert(
        &mut self,
        uri: Url,
        document: Arc<DocumentSnapshot>,
        require_existing: bool,
    ) -> bool {
        let Some(version) = document.version else {
            return false;
        };
        if !self.can_accept(&uri, version, require_existing) {
            return false;
        }
        self.documents.insert(uri, document);
        true
    }

    fn remove(&mut self, uri: &Url) {
        self.documents.remove(uri);
    }

    fn get(&self, uri: &Url) -> Option<Arc<DocumentSnapshot>> {
        self.documents.get(uri).map(Arc::clone)
    }
}

#[derive(Debug)]
struct ServerState {
    workspace_roots: Vec<PathBuf>,
    config_path: Option<PathBuf>,
    supports_document_changes: bool,
    trusted: bool,
}

impl Default for ServerState {
    fn default() -> Self {
        Self {
            workspace_roots: Vec::new(),
            config_path: None,
            supports_document_changes: false,
            trusted: true,
        }
    }
}

#[derive(Debug)]
struct Backend {
    client: Client,
    docs: Arc<Mutex<DocumentStore>>,
    state: Arc<Mutex<ServerState>>,
    workspace: Arc<Mutex<WorkspaceIndex>>,
    probes: probes::ProbeCatalog,
}

type PendingDocumentEdits = (
    Arc<DocumentSnapshot>,
    Vec<OneOf<TextEdit, AnnotatedTextEdit>>,
);

pub async fn run_server() -> AnyResult<()> {
    let stdin = tokio::io::stdin();
    let stdout = tokio::io::stdout();
    let (service, socket) = LspService::new(|client| Backend {
        client,
        docs: Arc::new(Mutex::new(DocumentStore::default())),
        state: Arc::new(Mutex::new(ServerState::default())),
        workspace: Arc::new(Mutex::new(WorkspaceIndex::default())),
        probes: probes::ProbeCatalog::default(),
    });
    Server::new(stdin, stdout, socket).serve(service).await;
    Ok(())
}

#[async_trait]
impl LanguageServer for Backend {
    async fn initialize(&self, params: InitializeParams) -> Result<InitializeResult> {
        let supports_rename = supports_document_changes(&params);
        self.configure_from_initialize(&params);
        let (roots, trusted) = {
            let state = self.state.lock().expect("server state poisoned");
            (state.workspace_roots.clone(), state.trusted)
        };
        self.workspace
            .lock()
            .expect("workspace index poisoned")
            .scan(if trusted { &roots } else { &[] });
        Ok(InitializeResult {
            server_info: Some(ServerInfo {
                name: "btfmt".to_string(),
                version: Some(env!("CARGO_PKG_VERSION").to_string()),
            }),
            capabilities: ServerCapabilities {
                text_document_sync: Some(TextDocumentSyncCapability::Options(
                    TextDocumentSyncOptions {
                        open_close: Some(true),
                        change: Some(TextDocumentSyncKind::FULL),
                        save: Some(TextDocumentSyncSaveOptions::Supported(true)),
                        ..TextDocumentSyncOptions::default()
                    },
                )),
                document_formatting_provider: Some(OneOf::Left(true)),
                hover_provider: Some(HoverProviderCapability::Simple(true)),
                completion_provider: Some(CompletionOptions {
                    trigger_characters: Some(vec![
                        "$".to_string(),
                        "@".to_string(),
                        ":".to_string(),
                        ".".to_string(),
                    ]),
                    ..CompletionOptions::default()
                }),
                document_symbol_provider: Some(OneOf::Left(true)),
                definition_provider: Some(OneOf::Left(true)),
                references_provider: Some(OneOf::Left(true)),
                document_highlight_provider: Some(OneOf::Left(true)),
                rename_provider: supports_rename.then(|| {
                    OneOf::Right(RenameOptions {
                        prepare_provider: Some(true),
                        work_done_progress_options: WorkDoneProgressOptions::default(),
                    })
                }),
                workspace: Some(WorkspaceServerCapabilities {
                    workspace_folders: Some(WorkspaceFoldersServerCapabilities {
                        supported: Some(true),
                        change_notifications: Some(OneOf::Left(true)),
                    }),
                    file_operations: None,
                }),
                ..ServerCapabilities::default()
            },
        })
    }

    async fn initialized(&self, _: InitializedParams) {
        self.client
            .log_message(MessageType::INFO, "btfmt server initialized")
            .await;
    }

    async fn shutdown(&self) -> Result<()> {
        Ok(())
    }

    async fn did_open(&self, params: DidOpenTextDocumentParams) {
        let uri = params.text_document.uri;
        let version = params.text_document.version;
        let text = params.text_document.text;
        self.upsert(uri, text, version, false).await;
    }

    async fn did_change(&self, params: DidChangeTextDocumentParams) {
        let Some(change) = params.content_changes.into_iter().last() else {
            return;
        };
        self.upsert(
            params.text_document.uri,
            change.text,
            params.text_document.version,
            true,
        )
        .await;
    }

    async fn did_close(&self, params: DidCloseTextDocumentParams) {
        let uri = params.text_document.uri;
        self.docs
            .lock()
            .expect("document store poisoned")
            .remove(&uri);
        self.workspace
            .lock()
            .expect("workspace index poisoned")
            .remove_overlay(&uri);
        self.client.publish_diagnostics(uri, Vec::new(), None).await;
    }

    async fn did_save(&self, params: DidSaveTextDocumentParams) {
        if !self.is_trusted() {
            return;
        }
        self.workspace
            .lock()
            .expect("workspace index poisoned")
            .refresh(&params.text_document.uri);
    }

    async fn did_change_watched_files(&self, params: DidChangeWatchedFilesParams) {
        if !self.is_trusted() {
            return;
        }
        let mut workspace = self.workspace.lock().expect("workspace index poisoned");
        for change in params.changes {
            match change.typ {
                FileChangeType::CREATED | FileChangeType::CHANGED => workspace.refresh(&change.uri),
                FileChangeType::DELETED => workspace.remove_disk(&change.uri),
                _ => {}
            }
        }
    }

    async fn did_change_workspace_folders(&self, params: DidChangeWorkspaceFoldersParams) {
        if !self.is_trusted() {
            return;
        }
        let mut state = self.state.lock().expect("server state poisoned");
        for removed in params.event.removed {
            if let Ok(path) = removed.uri.to_file_path() {
                state.workspace_roots.retain(|root| root != &path);
            }
        }
        for added in params.event.added {
            if let Ok(path) = added.uri.to_file_path() {
                if !state.workspace_roots.contains(&path) {
                    state.workspace_roots.push(path);
                }
            }
        }
        let roots = state.workspace_roots.clone();
        drop(state);
        self.workspace
            .lock()
            .expect("workspace index poisoned")
            .scan(&roots);
    }

    async fn did_change_configuration(&self, params: DidChangeConfigurationParams) {
        if let Some(path) = config_path_setting_from_value(&params.settings) {
            self.state
                .lock()
                .expect("server state poisoned")
                .config_path = path;
        }
    }

    async fn formatting(&self, params: DocumentFormattingParams) -> Result<Option<Vec<TextEdit>>> {
        let uri = params.text_document.uri;
        let Some(doc) = self.doc(&uri) else {
            return Ok(Some(Vec::new()));
        };
        let config = match self.config_for_uri(&uri) {
            Ok(config) => config,
            Err(err) => {
                self.client
                    .log_message(MessageType::ERROR, format!("config load failed: {err:#}"))
                    .await;
                return Ok(None);
            }
        };
        match format_source(&doc.text, &config) {
            Ok(text) => Ok(Some(formatting_edits(&doc, text))),
            Err(err) => {
                self.client
                    .log_message(MessageType::ERROR, format!("formatting failed: {err:#}"))
                    .await;
                Ok(None)
            }
        }
    }

    async fn document_symbol(
        &self,
        params: DocumentSymbolParams,
    ) -> Result<Option<DocumentSymbolResponse>> {
        let Some(doc) = self.doc(&params.text_document.uri) else {
            return Ok(Some(DocumentSymbolResponse::Nested(Vec::new())));
        };
        let Some(tree) = doc.tree.as_ref() else {
            return Ok(Some(DocumentSymbolResponse::Nested(Vec::new())));
        };
        Ok(Some(DocumentSymbolResponse::Nested(document_symbols(
            &doc.text,
            tree,
            &doc.line_index,
        ))))
    }

    async fn hover(&self, params: HoverParams) -> Result<Option<Hover>> {
        let uri = params.text_document_position_params.text_document.uri;
        let position = params.text_document_position_params.position;
        let Some(doc) = self.doc(&uri) else {
            return Ok(None);
        };
        let Some((word, range)) =
            identifier_at_position_with_index(&doc.text, &doc.line_index, position)
        else {
            return Ok(None);
        };
        let value = hover_markdown(&word).unwrap_or_else(|| format!("`{word}`"));
        Ok(Some(Hover {
            contents: HoverContents::Markup(MarkupContent {
                kind: MarkupKind::Markdown,
                value,
            }),
            range: Some(range),
        }))
    }

    async fn completion(&self, params: CompletionParams) -> Result<Option<CompletionResponse>> {
        let uri = params.text_document_position.text_document.uri;
        let Some(doc) = self.doc(&uri) else {
            return Ok(Some(CompletionResponse::List(CompletionList {
                is_incomplete: false,
                items: Vec::new(),
            })));
        };
        let workspace_symbols = self.workspace_completion_symbols(&uri);
        let incomplete = self
            .workspace
            .lock()
            .expect("workspace index poisoned")
            .is_incomplete();
        let mut list = completion::complete(
            &doc,
            params.text_document_position.position,
            workspace_symbols,
            incomplete,
        );
        if let Some(request) =
            probes::completion_request(&doc, params.text_document_position.position)
        {
            let dynamic = self
                .probes
                .complete(
                    request,
                    self.workspace_probe_metadata(&uri),
                    self.is_trusted(),
                )
                .await;
            list.is_incomplete |= dynamic.is_incomplete;
            list.items.extend(dynamic.items);
        }
        Ok(Some(CompletionResponse::List(list)))
    }

    async fn goto_definition(
        &self,
        params: GotoDefinitionParams,
    ) -> Result<Option<GotoDefinitionResponse>> {
        let uri = params.text_document_position_params.text_document.uri;
        let position = params.text_document_position_params.position;
        let Some(doc) = self.doc(&uri) else {
            return Ok(None);
        };
        let Some(index) = doc.symbols.as_ref() else {
            return Ok(None);
        };
        let offset = doc.line_index.offset_for_position(&doc.text, position);
        if let Some(target) = index.global_target_at(offset) {
            let locations: Vec<_> = self
                .global_occurrences(&uri, &target)
                .into_iter()
                .filter(|(_, _, occurrence)| occurrence.definition)
                .map(|(uri, snapshot, occurrence)| Location {
                    uri,
                    range: lsp_range(&snapshot, occurrence.range),
                })
                .collect();
            return Ok(match locations.len() {
                0 => None,
                1 => Some(GotoDefinitionResponse::Scalar(
                    locations.into_iter().next().unwrap(),
                )),
                _ => Some(GotoDefinitionResponse::Array(locations)),
            });
        }
        let ranges = index.definitions_at(offset);
        if ranges.is_empty() {
            return Ok(None);
        }
        let locations: Vec<_> = ranges
            .into_iter()
            .map(|range| Location {
                uri: uri.clone(),
                range: lsp_range(&doc, range),
            })
            .collect();
        Ok(Some(if locations.len() == 1 {
            GotoDefinitionResponse::Scalar(locations.into_iter().next().unwrap())
        } else {
            GotoDefinitionResponse::Array(locations)
        }))
    }

    async fn references(&self, params: ReferenceParams) -> Result<Option<Vec<Location>>> {
        let uri = params.text_document_position.text_document.uri;
        let position = params.text_document_position.position;
        let Some(doc) = self.doc(&uri) else {
            return Ok(Some(Vec::new()));
        };
        let Some(index) = doc.symbols.as_ref() else {
            return Ok(Some(Vec::new()));
        };
        let offset = doc.line_index.offset_for_position(&doc.text, position);
        if let Some(target) = index.global_target_at(offset) {
            let occurrences = self.global_occurrences(&uri, &target);
            if target.kind == GlobalSymbolKind::Macro
                && !occurrences
                    .iter()
                    .any(|(_, _, occurrence)| occurrence.definition)
            {
                return Ok(Some(Vec::new()));
            }
            let locations = occurrences
                .into_iter()
                .filter(|(_, _, occurrence)| {
                    params.context.include_declaration || !occurrence.definition
                })
                .map(|(uri, snapshot, occurrence)| Location {
                    uri,
                    range: lsp_range(&snapshot, occurrence.range),
                })
                .collect();
            return Ok(Some(locations));
        }
        let locations = index
            .references_at(offset, params.context.include_declaration)
            .into_iter()
            .map(|range| Location {
                uri: uri.clone(),
                range: lsp_range(&doc, range),
            })
            .collect();
        Ok(Some(locations))
    }

    async fn document_highlight(
        &self,
        params: DocumentHighlightParams,
    ) -> Result<Option<Vec<DocumentHighlight>>> {
        let uri = params.text_document_position_params.text_document.uri;
        let position = params.text_document_position_params.position;
        let Some(doc) = self.doc(&uri) else {
            return Ok(Some(Vec::new()));
        };
        let Some(index) = doc.symbols.as_ref() else {
            return Ok(Some(Vec::new()));
        };
        let offset = doc.line_index.offset_for_position(&doc.text, position);
        let highlights = index
            .highlights_at(offset)
            .into_iter()
            .map(|(range, access)| DocumentHighlight {
                range: lsp_range(&doc, range),
                kind: Some(match access {
                    AccessKind::Read => DocumentHighlightKind::READ,
                    AccessKind::Write | AccessKind::Declaration => DocumentHighlightKind::WRITE,
                }),
            })
            .collect();
        Ok(Some(highlights))
    }

    async fn prepare_rename(
        &self,
        params: TextDocumentPositionParams,
    ) -> Result<Option<PrepareRenameResponse>> {
        let Some(doc) = self.doc(&params.text_document.uri) else {
            return Ok(None);
        };
        let Some(index) = doc.symbols.as_ref() else {
            return Ok(None);
        };
        let offset = doc
            .line_index
            .offset_for_position(&doc.text, params.position);
        Ok(index
            .prepare_range(offset)
            .map(|range| PrepareRenameResponse::Range(lsp_range(&doc, range))))
    }

    async fn rename(&self, params: RenameParams) -> Result<Option<WorkspaceEdit>> {
        if !self
            .state
            .lock()
            .expect("server state poisoned")
            .supports_document_changes
        {
            return Err(Error::invalid_params(
                "client must support versioned documentChanges for rename",
            ));
        }
        let uri = params.text_document_position.text_document.uri;
        let position = params.text_document_position.position;
        let Some(doc) = self.doc(&uri) else {
            return Ok(None);
        };
        let Some(index) = doc.symbols.as_ref() else {
            return Ok(None);
        };
        let offset = doc.line_index.offset_for_position(&doc.text, position);
        if let Some(target) = index.global_target_at(offset) {
            return self
                .rename_global(&uri, &target, &params.new_name)
                .map_err(Error::invalid_params);
        }
        let Some(rename) = index
            .rename_at(offset, &params.new_name)
            .map_err(Error::invalid_params)?
        else {
            return Ok(None);
        };
        let edits = rename
            .ranges
            .into_iter()
            .map(|range| {
                OneOf::Left(TextEdit {
                    range: lsp_range(&doc, range),
                    new_text: rename.replacement.clone(),
                })
            })
            .collect();
        Ok(Some(WorkspaceEdit {
            changes: None,
            document_changes: Some(DocumentChanges::Edits(vec![TextDocumentEdit {
                text_document: OptionalVersionedTextDocumentIdentifier {
                    uri,
                    version: doc.version,
                },
                edits,
            }])),
            change_annotations: None,
        }))
    }
}

impl Backend {
    fn configure_from_initialize(&self, params: &InitializeParams) {
        let mut state = self.state.lock().expect("server state poisoned");
        state.workspace_roots = workspace_roots(params);
        state.config_path = params
            .initialization_options
            .as_ref()
            .and_then(config_path_setting_from_value)
            .flatten();
        state.supports_document_changes = supports_document_changes(params);
        state.trusted = params
            .initialization_options
            .as_ref()
            .and_then(trusted_setting_from_value)
            .unwrap_or(true);
    }

    fn is_trusted(&self) -> bool {
        self.state.lock().expect("server state poisoned").trusted
    }

    fn config_for_uri(&self, uri: &Url) -> AnyResult<Config> {
        let (workspace_roots, config_path) = {
            let state = self.state.lock().expect("server state poisoned");
            (state.workspace_roots.clone(), state.config_path.clone())
        };
        let document_path = uri.to_file_path().ok();
        let document_dir = document_path
            .as_deref()
            .and_then(Path::parent)
            .map(Path::to_path_buf);
        if let Some(config_path) = config_path {
            let base_dir = document_path
                .as_deref()
                .and_then(|path| workspace_root_for_document(&workspace_roots, path))
                .map(Path::to_path_buf)
                .or_else(|| document_dir.clone())
                .unwrap_or_else(|| PathBuf::from("."));
            return load_from_base(&base_dir, Some(&config_path));
        }

        if let Some(dir) = document_dir {
            return load_from_base(&dir, None);
        }
        Ok(Config::default())
    }

    async fn upsert(&self, uri: Url, text: String, version: i32, require_existing: bool) {
        if !self
            .docs
            .lock()
            .expect("document store poisoned")
            .can_accept(&uri, version, require_existing)
        {
            return;
        }

        let (document, diagnostics) = DocumentSnapshot::analyze(text, Some(version));
        let inserted = self.docs.lock().expect("document store poisoned").insert(
            uri.clone(),
            Arc::clone(&document),
            require_existing,
        );
        if !inserted {
            return;
        }
        self.workspace
            .lock()
            .expect("workspace index poisoned")
            .set_overlay(uri.clone(), document);
        if self
            .doc(&uri)
            .is_some_and(|doc| doc.version == Some(version))
        {
            self.client
                .publish_diagnostics(uri, diagnostics, Some(version))
                .await;
        }
    }

    fn doc(&self, uri: &Url) -> Option<Arc<DocumentSnapshot>> {
        self.docs.lock().expect("document store poisoned").get(uri)
    }

    fn effective_snapshot(&self, uri: &Url) -> Option<Arc<DocumentSnapshot>> {
        self.doc(uri).or_else(|| {
            self.workspace
                .lock()
                .expect("workspace index poisoned")
                .snapshot(uri)
        })
    }

    fn global_occurrences(
        &self,
        origin: &Url,
        target: &GlobalSymbolTarget,
    ) -> Vec<(Url, Arc<DocumentSnapshot>, GlobalSymbolOccurrence)> {
        let files = self
            .workspace
            .lock()
            .expect("workspace index poisoned")
            .program_files(origin);
        let mut seen = std::collections::HashSet::new();
        let mut occurrences = Vec::new();
        for uri in files {
            let Some(snapshot) = self.effective_snapshot(&uri) else {
                continue;
            };
            let Some(index) = snapshot.symbols.as_ref() else {
                continue;
            };
            for occurrence in index.global_occurrences(target) {
                if seen.insert((uri.clone(), occurrence.range.start, occurrence.range.end)) {
                    occurrences.push((uri.clone(), Arc::clone(&snapshot), occurrence));
                }
            }
        }
        occurrences
    }

    fn workspace_completion_symbols(&self, origin: &Url) -> Vec<self::symbols::CompletionSymbol> {
        let files = self
            .workspace
            .lock()
            .expect("workspace index poisoned")
            .program_files(origin);
        let mut symbols = Vec::new();
        for uri in files {
            let Some(snapshot) = self.effective_snapshot(&uri) else {
                continue;
            };
            let Some(index) = snapshot.completion_symbols.as_ref() else {
                continue;
            };
            symbols.extend(index.global_completion_symbols());
        }
        symbols
    }

    fn workspace_probe_metadata(&self, origin: &Url) -> probes::WorkspaceMetadata {
        let mut files = self
            .workspace
            .lock()
            .expect("workspace index poisoned")
            .program_files(origin);
        if !files.contains(origin) {
            files.push(origin.clone());
        }
        let mut metadata = probes::WorkspaceMetadata::default();
        for uri in files {
            if let Some(snapshot) = self.effective_snapshot(&uri) {
                metadata.add_snapshot(&snapshot);
            }
        }
        metadata
    }

    fn rename_global(
        &self,
        origin: &Url,
        target: &GlobalSymbolTarget,
        new_name: &str,
    ) -> std::result::Result<Option<WorkspaceEdit>, String> {
        if !self
            .workspace
            .lock()
            .expect("workspace index poisoned")
            .ready_for_rename()
        {
            return Err("workspace index is incomplete; cross-file rename is unavailable".into());
        }
        let name = SymbolIndex::normalize_global_rename(target, new_name)?;
        if name != target.name {
            let collision = GlobalSymbolTarget {
                kind: target.kind,
                name: name.clone(),
            };
            if !self.global_occurrences(origin, &collision).is_empty() {
                return Err(format!("rename target {new_name:?} already exists"));
            }
        }
        let replacement = match target.kind {
            GlobalSymbolKind::Map => format!("@{name}"),
            GlobalSymbolKind::Macro => name,
        };
        let occurrences = self.global_occurrences(origin, target);
        if occurrences.is_empty()
            || (target.kind == GlobalSymbolKind::Macro
                && !occurrences
                    .iter()
                    .any(|(_, _, occurrence)| occurrence.definition))
        {
            return Ok(None);
        }
        let mut by_document: HashMap<Url, PendingDocumentEdits> = HashMap::new();
        for (uri, snapshot, occurrence) in occurrences {
            by_document
                .entry(uri)
                .or_insert_with(|| (Arc::clone(&snapshot), Vec::new()))
                .1
                .push(OneOf::Left(TextEdit {
                    range: lsp_range(&snapshot, occurrence.range),
                    new_text: replacement.clone(),
                }));
        }
        let mut edits: Vec<_> = by_document
            .into_iter()
            .map(|(uri, (snapshot, edits))| TextDocumentEdit {
                text_document: OptionalVersionedTextDocumentIdentifier {
                    uri,
                    version: snapshot.version,
                },
                edits,
            })
            .collect();
        edits.sort_by(|left, right| {
            left.text_document
                .uri
                .as_str()
                .cmp(right.text_document.uri.as_str())
        });
        Ok(Some(WorkspaceEdit {
            changes: None,
            document_changes: Some(DocumentChanges::Edits(edits)),
            change_annotations: None,
        }))
    }
}

fn workspace_root_for_document<'a>(
    workspace_roots: &'a [PathBuf],
    document_path: &Path,
) -> Option<&'a Path> {
    workspace_roots
        .iter()
        .filter(|root| document_path.starts_with(root))
        .max_by_key(|root| root.components().count())
        .map(PathBuf::as_path)
}

fn supports_document_changes(params: &InitializeParams) -> bool {
    params
        .capabilities
        .workspace
        .as_ref()
        .and_then(|workspace| workspace.workspace_edit.as_ref())
        .and_then(|workspace_edit| workspace_edit.document_changes)
        .unwrap_or(false)
}

fn lsp_range(doc: &DocumentSnapshot, range: ByteRange) -> Range {
    doc.line_index
        .range_for_offsets(&doc.text, range.start, range.end)
}

fn formatting_edits(doc: &DocumentSnapshot, formatted: String) -> Vec<TextEdit> {
    if formatted == doc.text.as_ref() {
        Vec::new()
    } else {
        vec![TextEdit {
            range: doc.line_index.full_range(&doc.text),
            new_text: formatted,
        }]
    }
}

fn workspace_roots(params: &InitializeParams) -> Vec<PathBuf> {
    if let Some(folders) = &params.workspace_folders {
        let roots: Vec<PathBuf> = folders
            .iter()
            .filter_map(|folder| folder.uri.to_file_path().ok())
            .collect();
        if !roots.is_empty() {
            return roots;
        }
    }
    params
        .root_uri
        .as_ref()
        .and_then(|uri| uri.to_file_path().ok())
        .into_iter()
        .collect()
}

fn config_path_setting_from_value(value: &Value) -> Option<Option<PathBuf>> {
    let path = value
        .get("btfmt")
        .and_then(|btfmt| btfmt.get("configPath"))
        .or_else(|| value.get("configPath"))
        .and_then(Value::as_str)?;
    if path.is_empty() {
        Some(None)
    } else {
        Some(Some(PathBuf::from(path)))
    }
}

fn trusted_setting_from_value(value: &Value) -> Option<bool> {
    value
        .get("btfmt")
        .and_then(|btfmt| btfmt.get("trusted"))
        .or_else(|| value.get("trusted"))
        .and_then(Value::as_bool)
}

fn document_symbols(
    text: &str,
    tree: &Tree,
    line_index: &crate::text::LineIndex,
) -> Vec<DocumentSymbol> {
    let root = tree.root_node();
    let mut cursor = root.walk();
    root.named_children(&mut cursor)
        .filter_map(|node| match node.kind() {
            "action_block" => {
                let mut child_cursor = node.walk();
                let probes = node
                    .named_children(&mut child_cursor)
                    .find(|child| child.kind() == "probes_list")?;
                let name = text_for_node(text, probes);
                Some(new_document_symbol(
                    name,
                    SymbolKind::EVENT,
                    line_index.range_for_offsets(text, node.start_byte(), node.end_byte()),
                    line_index.range_for_offsets(text, probes.start_byte(), probes.end_byte()),
                ))
            }
            "macro_definition" => {
                let name_node = node.child_by_field_name("name")?;
                let name = format!("macro {}", text_for_node(text, name_node));
                Some(new_document_symbol(
                    name,
                    SymbolKind::FUNCTION,
                    line_index.range_for_offsets(text, node.start_byte(), node.end_byte()),
                    line_index.range_for_offsets(
                        text,
                        name_node.start_byte(),
                        name_node.end_byte(),
                    ),
                ))
            }
            _ => None,
        })
        .collect()
}

fn text_for_node(text: &str, node: tree_sitter::Node<'_>) -> String {
    text[node.byte_range()]
        .split_whitespace()
        .collect::<Vec<_>>()
        .join(" ")
}

fn new_document_symbol(
    name: String,
    kind: SymbolKind,
    range: Range,
    selection_range: Range,
) -> DocumentSymbol {
    #[allow(deprecated)]
    DocumentSymbol {
        name,
        detail: None,
        kind,
        tags: None,
        deprecated: None,
        range,
        selection_range,
        children: None,
    }
}

fn hover_markdown(word: &str) -> Option<String> {
    let entry = catalog::find(word)?;
    Some(format!(
        "**bpftrace** `{word}`\n\n{}\n\n_{}_",
        entry.documentation,
        catalog::CATALOG_VERSION
    ))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::parse;

    #[test]
    fn document_symbols_follow_action_block_syntax() {
        let text = concat!(
            "BEGIN\r\n",
            "{\r\n",
            "    exit();\r\n",
            "}\r\n",
            "\r\n",
            "kprobe:vfs_read*,\r\n",
            "kprobe:vfs_write* /pid == 1/\r\n",
            "{\r\n",
            "    exit();\r\n",
            "}\r\n",
        );

        let tree = parse::ensure_valid(text).unwrap();
        let line_index = crate::text::LineIndex::new(text);
        let symbols = document_symbols(text, &tree, &line_index);
        assert_eq!(symbols.len(), 2);
        assert_eq!(symbols[0].name, "BEGIN");
        assert_eq!(symbols[1].name, "kprobe:vfs_read*, kprobe:vfs_write*");
        assert_eq!(symbols[0].selection_range.start, Position::new(0, 0));
        assert_eq!(symbols[1].selection_range.start, Position::new(5, 0));
        assert_eq!(symbols[1].range.end.line, 9);
    }

    #[test]
    fn document_symbols_include_macros() {
        let text = "macro add_one(value) { value + 1 }\n";
        let tree = parse::ensure_valid(text).unwrap();
        let line_index = crate::text::LineIndex::new(text);
        let symbols = document_symbols(text, &tree, &line_index);

        assert_eq!(symbols.len(), 1);
        assert_eq!(symbols[0].name, "macro add_one");
        assert_eq!(symbols[0].kind, SymbolKind::FUNCTION);
    }

    #[test]
    fn workspace_root_uses_longest_matching_path() {
        let roots = vec![
            PathBuf::from("/workspace"),
            PathBuf::from("/workspace/nested"),
        ];

        assert_eq!(
            workspace_root_for_document(&roots, Path::new("/workspace/nested/script.bt")),
            Some(Path::new("/workspace/nested"))
        );
        assert_eq!(
            workspace_root_for_document(&roots, Path::new("/workspace/other.bt")),
            Some(Path::new("/workspace"))
        );
        assert_eq!(
            workspace_root_for_document(&roots, Path::new("/outside/script.bt")),
            None
        );
    }

    #[test]
    fn document_store_rejects_stale_and_closed_changes() {
        let uri = Url::parse("file:///workspace/script.bt").unwrap();
        let mut store = DocumentStore::default();
        let snapshot =
            |text: &str, version| DocumentSnapshot::analyze(text.to_string(), Some(version)).0;
        assert!(store.insert(uri.clone(), snapshot("v1", 1), false,));
        assert!(store.insert(uri.clone(), snapshot("v3", 3), true,));
        for version in [2, 3] {
            assert!(!store.insert(
                uri.clone(),
                snapshot(&format!("stale-{version}"), version),
                true,
            ));
        }
        let current = store.get(&uri).unwrap();
        assert_eq!(current.version, Some(3));
        assert_eq!(current.text.as_ref(), "v3");
        assert!(Arc::ptr_eq(&current, &store.get(&uri).unwrap()));

        store.remove(&uri);
        assert!(!store.insert(uri.clone(), snapshot("closed", 4), true,));
        assert!(store.get(&uri).is_none());
    }

    #[test]
    fn rename_requires_versioned_document_change_support() {
        let mut params = InitializeParams::default();
        assert!(!supports_document_changes(&params));

        params.capabilities.workspace = Some(WorkspaceClientCapabilities {
            workspace_edit: Some(WorkspaceEditClientCapabilities {
                document_changes: Some(true),
                ..WorkspaceEditClientCapabilities::default()
            }),
            ..WorkspaceClientCapabilities::default()
        });
        assert!(supports_document_changes(&params));
    }

    #[test]
    fn trusted_setting_defaults_and_parses_explicit_values() {
        assert_eq!(trusted_setting_from_value(&serde_json::json!({})), None);
        assert_eq!(
            trusted_setting_from_value(&serde_json::json!({"btfmt": {"trusted": false}})),
            Some(false)
        );
    }

    #[test]
    fn unchanged_formatting_returns_no_edits() {
        let (snapshot, _) = DocumentSnapshot::analyze("BEGIN\n{\n}\n".to_string(), Some(1));
        assert!(formatting_edits(&snapshot, snapshot.text.to_string()).is_empty());
        assert_eq!(
            formatting_edits(&snapshot, "BEGIN {}\n".to_string()).len(),
            1
        );
    }
}
