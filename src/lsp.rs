use crate::config::Config;
use crate::format::format_source;
use crate::parse;
use crate::text::{
    all_identifier_occurrences, full_range, identifier_at_position, range_for_offsets,
};
use anyhow::Result as AnyResult;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use tower_lsp::jsonrpc::Result;
use tower_lsp::lsp_types::*;
use tower_lsp::{async_trait, Client, LanguageServer, LspService, Server};

#[derive(Debug, Clone)]
struct Document {
    text: String,
}

#[derive(Debug)]
struct Backend {
    client: Client,
    docs: Arc<Mutex<HashMap<Url, Document>>>,
}

pub async fn run_server() -> AnyResult<()> {
    let stdin = tokio::io::stdin();
    let stdout = tokio::io::stdout();
    let (service, socket) = LspService::new(|client| Backend {
        client,
        docs: Arc::new(Mutex::new(HashMap::new())),
    });
    Server::new(stdin, stdout, socket).serve(service).await;
    Ok(())
}

#[async_trait]
impl LanguageServer for Backend {
    async fn initialize(&self, _: InitializeParams) -> Result<InitializeResult> {
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
                definition_provider: Some(OneOf::Left(true)),
                references_provider: Some(OneOf::Left(true)),
                document_highlight_provider: Some(OneOf::Left(true)),
                rename_provider: Some(OneOf::Right(RenameOptions {
                    prepare_provider: Some(true),
                    work_done_progress_options: WorkDoneProgressOptions::default(),
                })),
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

    async fn formatting(&self, params: DocumentFormattingParams) -> Result<Option<Vec<TextEdit>>> {
        let Some(doc) = self.doc(&params.text_document.uri) else {
            return Ok(Some(Vec::new()));
        };
        match format_source(&doc.text, &Config::default()) {
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

    async fn goto_definition(
        &self,
        params: GotoDefinitionParams,
    ) -> Result<Option<GotoDefinitionResponse>> {
        let uri = params.text_document_position_params.text_document.uri;
        let position = params.text_document_position_params.position;
        let Some(doc) = self.doc(&uri) else {
            return Ok(None);
        };
        let Some((word, _)) = identifier_at_position(&doc.text, position) else {
            return Ok(None);
        };
        let Some(range) = all_identifier_occurrences(&doc.text, &word)
            .into_iter()
            .next()
        else {
            return Ok(None);
        };
        Ok(Some(GotoDefinitionResponse::Array(vec![Location {
            uri,
            range,
        }])))
    }

    async fn references(&self, params: ReferenceParams) -> Result<Option<Vec<Location>>> {
        let uri = params.text_document_position.text_document.uri;
        let position = params.text_document_position.position;
        let Some(doc) = self.doc(&uri) else {
            return Ok(Some(Vec::new()));
        };
        let Some((word, _)) = identifier_at_position(&doc.text, position) else {
            return Ok(Some(Vec::new()));
        };
        let locations = all_identifier_occurrences(&doc.text, &word)
            .into_iter()
            .map(|range| Location {
                uri: uri.clone(),
                range,
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
        let Some((word, _)) = identifier_at_position(&doc.text, position) else {
            return Ok(Some(Vec::new()));
        };
        let highlights = all_identifier_occurrences(&doc.text, &word)
            .into_iter()
            .map(|range| DocumentHighlight {
                range,
                kind: Some(DocumentHighlightKind::TEXT),
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
        Ok(identifier_at_position(&doc.text, params.position)
            .map(|(_, range)| PrepareRenameResponse::Range(range)))
    }

    async fn rename(&self, params: RenameParams) -> Result<Option<WorkspaceEdit>> {
        let uri = params.text_document_position.text_document.uri;
        let position = params.text_document_position.position;
        let Some(doc) = self.doc(&uri) else {
            return Ok(None);
        };
        let Some((word, _)) = identifier_at_position(&doc.text, position) else {
            return Ok(None);
        };
        let edits: Vec<TextEdit> = all_identifier_occurrences(&doc.text, &word)
            .into_iter()
            .map(|range| TextEdit {
                range,
                new_text: params.new_name.clone(),
            })
            .collect();
        let mut changes = HashMap::new();
        changes.insert(uri, edits);
        Ok(Some(WorkspaceEdit {
            changes: Some(changes),
            document_changes: None,
            change_annotations: None,
        }))
    }
}

impl Backend {
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
