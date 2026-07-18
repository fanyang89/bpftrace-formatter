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
    let Ok(parsed) = parse::parse(text) else {
        return Vec::new();
    };
    let root = parsed.tree.root_node();
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
                    range_for_offsets(text, node.start_byte(), node.end_byte()),
                    range_for_offsets(text, probes.start_byte(), probes.end_byte()),
                ))
            }
            "macro_definition" => {
                let name_node = node.child_by_field_name("name")?;
                let name = format!("macro {}", text_for_node(text, name_node));
                Some(new_document_symbol(
                    name,
                    SymbolKind::FUNCTION,
                    range_for_offsets(text, node.start_byte(), node.end_byte()),
                    range_for_offsets(text, name_node.start_byte(), name_node.end_byte()),
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

#[cfg(test)]
mod tests {
    use super::*;

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

        let symbols = document_symbols(text);
        assert_eq!(symbols.len(), 2);
        assert_eq!(symbols[0].name, "BEGIN");
        assert_eq!(symbols[1].name, "kprobe:vfs_read*, kprobe:vfs_write*");
        assert_eq!(symbols[0].selection_range.start, Position::new(0, 0));
        assert_eq!(symbols[1].selection_range.start, Position::new(5, 0));
        assert_eq!(symbols[1].range.end.line, 9);
    }

    #[test]
    fn document_symbols_include_macros() {
        let text = "macro add_one(value) { return value + 1; }\n";
        let symbols = document_symbols(text);

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
}
