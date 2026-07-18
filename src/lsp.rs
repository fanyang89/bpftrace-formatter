use crate::config::{load_from_base, Config};
use crate::format::format_source;
use crate::parse;
use crate::text::{full_range, identifier_at_position, range_for_offsets};
use anyhow::Result as AnyResult;
use serde_json::Value;
use std::collections::HashMap;
use std::path::{Path, PathBuf};
use std::sync::{Arc, Mutex};
use tower_lsp::jsonrpc::Result;
use tower_lsp::lsp_types::*;
use tower_lsp::{async_trait, Client, LanguageServer, LspService, Server};

#[derive(Debug, Clone)]
struct Document {
    text: String,
}

#[derive(Debug, Default)]
struct ServerState {
    workspace_roots: Vec<PathBuf>,
    config_path: Option<PathBuf>,
}

#[derive(Debug)]
struct Backend {
    client: Client,
    docs: Arc<Mutex<HashMap<Url, Document>>>,
    state: Arc<Mutex<ServerState>>,
}

pub async fn run_server() -> AnyResult<()> {
    let stdin = tokio::io::stdin();
    let stdout = tokio::io::stdout();
    let (service, socket) = LspService::new(|client| Backend {
        client,
        docs: Arc::new(Mutex::new(HashMap::new())),
        state: Arc::new(Mutex::new(ServerState::default())),
    });
    Server::new(stdin, stdout, socket).serve(service).await;
    Ok(())
}

#[async_trait]
impl LanguageServer for Backend {
    async fn initialize(&self, params: InitializeParams) -> Result<InitializeResult> {
        self.configure_from_initialize(&params);
        Ok(InitializeResult {
            server_info: Some(ServerInfo {
                name: "btfmt".to_string(),
                version: Some(env!("CARGO_PKG_VERSION").to_string()),
            }),
            capabilities: ServerCapabilities {
                text_document_sync: Some(TextDocumentSyncCapability::Kind(
                    TextDocumentSyncKind::FULL,
                )),
                document_formatting_provider: Some(OneOf::Left(true)),
                hover_provider: Some(HoverProviderCapability::Simple(true)),
                completion_provider: Some(CompletionOptions::default()),
                document_symbol_provider: Some(OneOf::Left(true)),
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
        let version = Some(params.text_document.version);
        let text = params.text_document.text;
        self.upsert(uri, text, version).await;
    }

    async fn did_change(&self, params: DidChangeTextDocumentParams) {
        let Some(change) = params.content_changes.into_iter().last() else {
            return;
        };
        self.upsert(
            params.text_document.uri,
            change.text,
            Some(params.text_document.version),
        )
        .await;
    }

    async fn did_close(&self, params: DidCloseTextDocumentParams) {
        let uri = params.text_document.uri;
        self.docs
            .lock()
            .expect("document store poisoned")
            .remove(&uri);
        self.client.publish_diagnostics(uri, Vec::new(), None).await;
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
            Ok(text) => Ok(Some(vec![TextEdit {
                range: full_range(&doc.text),
                new_text: text,
            }])),
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
        Ok(Some(DocumentSymbolResponse::Nested(document_symbols(
            &doc.text,
        ))))
    }

    async fn hover(&self, params: HoverParams) -> Result<Option<Hover>> {
        let uri = params.text_document_position_params.text_document.uri;
        let position = params.text_document_position_params.position;
        let Some(doc) = self.doc(&uri) else {
            return Ok(None);
        };
        let Some((word, range)) = identifier_at_position(&doc.text, position) else {
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

    async fn completion(&self, _: CompletionParams) -> Result<Option<CompletionResponse>> {
        Ok(Some(CompletionResponse::Array(completion_items())))
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
    }

    fn config_for_uri(&self, uri: &Url) -> AnyResult<Config> {
        let state = self.state.lock().expect("server state poisoned");
        let document_dir = uri
            .to_file_path()
            .ok()
            .and_then(|path| path.parent().map(Path::to_path_buf));
        if let Some(config_path) = &state.config_path {
            if config_path.as_os_str().is_empty() {
                return Ok(Config::default());
            }
            let path = if config_path.is_absolute() {
                config_path.clone()
            } else if let Some(root) = state.workspace_roots.first() {
                root.join(config_path)
            } else if let Some(dir) = &document_dir {
                dir.join(config_path)
            } else {
                config_path.clone()
            };
            if path.exists() {
                return Config::load(&path);
            }
            return Ok(Config::default());
        }

        if let Some(dir) = document_dir {
            return load_from_base(&dir, None);
        }
        Ok(Config::default())
    }

    async fn upsert(&self, uri: Url, text: String, version: Option<i32>) {
        let diagnostics = parse::parse(&text)
            .map(|parsed| parsed.diagnostics)
            .unwrap_or_else(|err| {
                vec![Diagnostic {
                    range: full_range(&text),
                    severity: Some(DiagnosticSeverity::ERROR),
                    source: Some("btfmt".to_string()),
                    message: format!("parse failed: {err:#}"),
                    ..Diagnostic::default()
                }]
            });
        self.docs
            .lock()
            .expect("document store poisoned")
            .insert(uri.clone(), Document { text });
        self.client
            .publish_diagnostics(uri, diagnostics, version)
            .await;
    }

    fn doc(&self, uri: &Url) -> Option<Document> {
        self.docs
            .lock()
            .expect("document store poisoned")
            .get(uri)
            .cloned()
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

fn document_symbols(text: &str) -> Vec<DocumentSymbol> {
    text.lines()
        .enumerate()
        .filter_map(|(line_idx, line)| {
            let trimmed = line.trim();
            if !trimmed.contains('{')
                || trimmed.starts_with("if")
                || trimmed.starts_with("while")
                || trimmed.starts_with("for")
            {
                return None;
            }
            let name = trimmed.split('{').next().unwrap_or(trimmed).trim();
            if name.is_empty() {
                return None;
            }
            let start = text
                .lines()
                .take(line_idx)
                .map(|line| line.len() + 1)
                .sum::<usize>();
            let end = start + line.len();
            let range = range_for_offsets(text, start, end);
            #[allow(deprecated)]
            let symbol = DocumentSymbol {
                name: name.to_string(),
                detail: None,
                kind: SymbolKind::EVENT,
                tags: None,
                deprecated: None,
                range,
                selection_range: range,
                children: None,
            };
            Some(symbol)
        })
        .collect()
}

fn completion_items() -> Vec<CompletionItem> {
    [
        ("BEGIN", CompletionItemKind::KEYWORD),
        ("END", CompletionItemKind::KEYWORD),
        ("tracepoint", CompletionItemKind::EVENT),
        ("kprobe", CompletionItemKind::EVENT),
        ("kretprobe", CompletionItemKind::EVENT),
        ("uprobe", CompletionItemKind::EVENT),
        ("printf", CompletionItemKind::FUNCTION),
        ("print", CompletionItemKind::FUNCTION),
        ("count", CompletionItemKind::FUNCTION),
        ("sum", CompletionItemKind::FUNCTION),
        ("avg", CompletionItemKind::FUNCTION),
        ("hist", CompletionItemKind::FUNCTION),
        ("if", CompletionItemKind::KEYWORD),
        ("while", CompletionItemKind::KEYWORD),
        ("for", CompletionItemKind::KEYWORD),
    ]
    .into_iter()
    .map(|(label, kind)| CompletionItem {
        label: label.to_string(),
        kind: Some(kind),
        ..CompletionItem::default()
    })
    .collect()
}

fn hover_markdown(word: &str) -> Option<String> {
    let value = match word {
        "BEGIN" => "Run once when bpftrace starts.",
        "END" => "Run once when bpftrace exits.",
        "tracepoint" => "Kernel tracepoint probe.",
        "kprobe" => "Kernel function entry probe.",
        "kretprobe" => "Kernel function return probe.",
        "printf" => "Print formatted output.",
        "print" => "Print a value or map.",
        "count" => "Count occurrences in a map aggregation.",
        "sum" => "Sum values in a map aggregation.",
        "hist" => "Build a power-of-two histogram.",
        _ => return None,
    };
    Some(format!("**bpftrace** `{word}`\n\n{value}"))
}
