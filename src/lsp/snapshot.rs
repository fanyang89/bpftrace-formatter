use super::symbols::SymbolIndex;
use crate::parse;
use crate::text::LineIndex;
use std::sync::Arc;
use tower_lsp::lsp_types::{Diagnostic, DiagnosticSeverity};
use tree_sitter::Tree;

#[derive(Debug)]
pub(super) struct DocumentSnapshot {
    pub text: Arc<str>,
    pub version: Option<i32>,
    pub tree: Option<Tree>,
    pub line_index: LineIndex,
    pub symbols: Option<Arc<SymbolIndex>>,
    pub completion_symbols: Option<Arc<SymbolIndex>>,
}

impl DocumentSnapshot {
    pub(super) fn analyze(text: String, version: Option<i32>) -> (Arc<Self>, Vec<Diagnostic>) {
        let line_index = LineIndex::new(&text);
        let (tree, diagnostics) = match parse::parse(&text) {
            Ok(parsed) => (Some(parsed.tree), parsed.diagnostics),
            Err(err) => (
                None,
                vec![Diagnostic {
                    range: line_index.full_range(&text),
                    severity: Some(DiagnosticSeverity::ERROR),
                    source: Some("btfmt".to_string()),
                    message: format!("parse failed: {err:#}"),
                    ..Diagnostic::default()
                }],
            ),
        };
        let symbols = tree
            .as_ref()
            .and_then(|tree| SymbolIndex::new(&text, tree))
            .map(Arc::new);
        let completion_symbols = symbols.clone().or_else(|| {
            tree.as_ref()
                .map(|tree| Arc::new(SymbolIndex::new_for_completion(&text, tree)))
        });
        (
            Arc::new(Self {
                text: Arc::from(text),
                version,
                tree,
                line_index,
                symbols,
                completion_symbols,
            }),
            diagnostics,
        )
    }
}
