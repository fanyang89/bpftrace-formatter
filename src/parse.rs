use crate::text::range_for_offsets;
use anyhow::{Context, Result};
use tower_lsp::lsp_types::{Diagnostic, DiagnosticSeverity, Range};
use tree_sitter::{Node, Parser, Tree};

pub struct ParseResult {
    pub tree: Tree,
    pub diagnostics: Vec<Diagnostic>,
}

pub fn parse(source: &str) -> Result<ParseResult> {
    let mut parser = Parser::new();
    let language = tree_sitter_bpftrace::LANGUAGE;
    parser
        .set_language(&language.into())
        .context("loading tree-sitter-bpftrace grammar")?;
    let tree = parser
        .parse(source, None)
        .context("parsing bpftrace source")?;
    let diagnostics = diagnostics_for_tree(source, &tree);
    Ok(ParseResult { tree, diagnostics })
}

pub fn ensure_valid(source: &str) -> Result<Tree> {
    let result = parse(source)?;
    if !result.diagnostics.is_empty() {
        let first = &result.diagnostics[0];
        anyhow::bail!("parse failed: {}", first.message);
    }
    Ok(result.tree)
}

pub fn diagnostics_for_tree(source: &str, tree: &Tree) -> Vec<Diagnostic> {
    let mut diagnostics = Vec::new();
    collect_diagnostics(source, tree.root_node(), &mut diagnostics);
    diagnostics
}

fn collect_diagnostics(source: &str, node: Node<'_>, diagnostics: &mut Vec<Diagnostic>) {
    if node.is_error() || node.is_missing() {
        let range = if node.start_byte() <= node.end_byte() {
            range_for_offsets(
                source,
                node.start_byte(),
                node.end_byte().max(node.start_byte() + 1),
            )
        } else {
            Range::default()
        };
        let kind = if node.is_missing() {
            "missing"
        } else {
            "error"
        };
        diagnostics.push(Diagnostic {
            range,
            severity: Some(DiagnosticSeverity::ERROR),
            source: Some("btfmt".to_string()),
            message: format!("tree-sitter parse {kind} at {}", node.kind()),
            ..Diagnostic::default()
        });
    }

    let mut cursor = node.walk();
    for child in node.children(&mut cursor) {
        if child.has_error() || child.is_error() || child.is_missing() {
            collect_diagnostics(source, child, diagnostics);
        }
    }
}
