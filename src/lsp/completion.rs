use super::catalog;
use super::snapshot::DocumentSnapshot;
use super::symbols::{CompletionSymbol, CompletionSymbolKind};
use std::collections::HashSet;
use tower_lsp::lsp_types::{
    CompletionItem, CompletionItemKind, CompletionList, CompletionTextEdit, Position, TextEdit,
};
use tree_sitter::Node;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum CursorContext {
    TopLevel,
    ProbeProvider,
    ProbeTarget,
    Expression,
    Statement,
    MemberAccess,
    Suppressed,
    Unknown,
}

#[derive(Debug)]
struct CompletionPrefix {
    prefix: String,
    start: usize,
    end: usize,
}

#[derive(Debug)]
struct RankedItem {
    rank: usize,
    scope_distance: usize,
    item: CompletionItem,
}

pub(super) fn complete(
    snapshot: &DocumentSnapshot,
    position: Position,
    workspace_symbols: Vec<CompletionSymbol>,
    is_incomplete: bool,
) -> CompletionList {
    let text = snapshot.text.as_ref();
    let offset = snapshot.line_index.offset_for_position(text, position);
    let context = cursor_context(snapshot, offset);
    if context == CursorContext::Suppressed || context == CursorContext::MemberAccess {
        return CompletionList {
            is_incomplete,
            items: Vec::new(),
        };
    }
    let prefix = completion_prefix(text, offset);
    let range = snapshot
        .line_index
        .range_for_offsets(text, prefix.start, prefix.end);
    let mut ranked = Vec::new();
    if matches!(
        context,
        CursorContext::Expression | CursorContext::Statement | CursorContext::Unknown
    ) {
        if let Some(index) = snapshot.completion_symbols.as_ref() {
            for symbol in index.visible_symbols_at(offset) {
                push_symbol(&mut ranked, symbol, &prefix.prefix, range);
            }
        }
        for symbol in workspace_symbols {
            push_symbol(&mut ranked, symbol, &prefix.prefix, range);
        }
    }

    let context_mask = match context {
        CursorContext::TopLevel => catalog::TOP_LEVEL,
        CursorContext::ProbeProvider => catalog::PROBE_PROVIDER,
        CursorContext::Expression | CursorContext::Unknown => catalog::EXPRESSION,
        CursorContext::Statement => catalog::STATEMENT | catalog::EXPRESSION,
        CursorContext::ProbeTarget | CursorContext::MemberAccess | CursorContext::Suppressed => 0,
    };
    if !prefix.prefix.starts_with(['$', '@']) {
        for entry in catalog::entries()
            .iter()
            .filter(|entry| entry.contexts & context_mask != 0)
        {
            if !entry.label.starts_with(&prefix.prefix) {
                continue;
            }
            let mut item = entry.completion_item();
            item.text_edit = Some(CompletionTextEdit::Edit(TextEdit {
                range,
                new_text: entry.label.to_string(),
            }));
            ranked.push(RankedItem {
                rank: match entry.kind {
                    catalog::CatalogKind::Provider => 4,
                    catalog::CatalogKind::Value => 5,
                    catalog::CatalogKind::Function => 6,
                    catalog::CatalogKind::Keyword => 7,
                },
                scope_distance: usize::MAX,
                item,
            });
        }
    }

    ranked.sort_by(|left, right| {
        (left.rank, left.scope_distance, left.item.label.as_str()).cmp(&(
            right.rank,
            right.scope_distance,
            right.item.label.as_str(),
        ))
    });
    let mut seen = HashSet::new();
    let mut items = Vec::new();
    for (idx, mut ranked) in ranked.into_iter().enumerate() {
        if !seen.insert(ranked.item.label.clone()) {
            continue;
        }
        ranked.item.sort_text = Some(format!("{idx:06}"));
        items.push(ranked.item);
    }
    CompletionList {
        is_incomplete,
        items,
    }
}

fn push_symbol(
    ranked: &mut Vec<RankedItem>,
    symbol: CompletionSymbol,
    prefix: &str,
    range: tower_lsp::lsp_types::Range,
) {
    if !symbol.label.starts_with(prefix) {
        return;
    }
    if prefix.starts_with('$') && symbol.kind != CompletionSymbolKind::Scratch {
        return;
    }
    if prefix.starts_with('@') && symbol.kind != CompletionSymbolKind::Map {
        return;
    }
    if !prefix.starts_with(['$', '@'])
        && matches!(
            symbol.kind,
            CompletionSymbolKind::Scratch | CompletionSymbolKind::Map
        )
    {
        return;
    }
    let (kind, rank, detail) = match symbol.kind {
        CompletionSymbolKind::Scratch => (CompletionItemKind::VARIABLE, 0, "scratch variable"),
        CompletionSymbolKind::Parameter => (CompletionItemKind::VARIABLE, 1, "macro parameter"),
        CompletionSymbolKind::Map => (CompletionItemKind::VARIABLE, 2, "map"),
        CompletionSymbolKind::Macro => (CompletionItemKind::FUNCTION, 3, "macro"),
    };
    ranked.push(RankedItem {
        rank,
        scope_distance: symbol.scope_distance,
        item: CompletionItem {
            label: symbol.label.clone(),
            kind: Some(kind),
            detail: Some(detail.to_string()),
            text_edit: Some(CompletionTextEdit::Edit(TextEdit {
                range,
                new_text: symbol.label,
            })),
            ..CompletionItem::default()
        },
    });
}

fn completion_prefix(text: &str, offset: usize) -> CompletionPrefix {
    let bytes = text.as_bytes();
    let mut start = offset.min(bytes.len());
    while start > 0 && is_completion_byte(bytes[start - 1]) {
        start -= 1;
    }
    let mut end = offset.min(bytes.len());
    while end < bytes.len() && is_completion_byte(bytes[end]) {
        end += 1;
    }
    CompletionPrefix {
        prefix: text[start..offset.min(text.len())].to_string(),
        start,
        end,
    }
}

fn is_completion_byte(byte: u8) -> bool {
    byte.is_ascii_alphanumeric() || matches!(byte, b'_' | b'$' | b'@')
}

fn cursor_context(snapshot: &DocumentSnapshot, offset: usize) -> CursorContext {
    let text = snapshot.text.as_ref();
    let before = &text[..offset.min(text.len())];
    if before.ends_with('.') || before.ends_with("->") {
        return CursorContext::MemberAccess;
    }
    if let Some(tree) = snapshot.tree.as_ref() {
        let start = offset.saturating_sub(1);
        if let Some(node) = tree
            .root_node()
            .descendant_for_byte_range(start, offset.max(start + 1).min(text.len()))
        {
            if has_ancestor_kind(
                node,
                &[
                    "line_comment",
                    "block_comment",
                    "string_literal",
                    "c_preproc",
                    "c_preproc_block",
                ],
            ) {
                return CursorContext::Suppressed;
            }
            if has_ancestor_kind(node, &["predicate"]) {
                return CursorContext::Expression;
            }
            if has_ancestor_kind(
                node,
                &["action", "block", "block_expression", "macro_definition"],
            ) {
                let previous = before.trim_end().chars().last();
                return if matches!(previous, None | Some('{') | Some(';') | Some('}')) {
                    CursorContext::Statement
                } else {
                    CursorContext::Expression
                };
            }
            if has_ancestor_kind(node, &["probe", "probes_list"]) {
                let line = before.rsplit('\n').next().unwrap_or(before);
                return if line.contains(':') {
                    CursorContext::ProbeTarget
                } else {
                    CursorContext::ProbeProvider
                };
            }
            if has_ancestor_kind(node, &["source_file", "preamble"]) {
                return CursorContext::TopLevel;
            }
        }
    }
    if before
        .rfind('{')
        .is_some_and(|open| before.rfind('}').is_none_or(|close| open > close))
    {
        CursorContext::Statement
    } else if before
        .rsplit('\n')
        .next()
        .is_some_and(|line| !line.contains(':'))
    {
        CursorContext::TopLevel
    } else {
        CursorContext::Unknown
    }
}

fn has_ancestor_kind(mut node: Node<'_>, kinds: &[&str]) -> bool {
    loop {
        if kinds.contains(&node.kind()) {
            return true;
        }
        let Some(parent) = node.parent() else {
            return false;
        };
        node = parent;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn labels(source: &str, marker: &str) -> Vec<String> {
        let offset = source.find(marker).unwrap();
        let text = source.replacen(marker, "", 1);
        let (snapshot, _) = DocumentSnapshot::analyze(text, Some(1));
        let position = snapshot
            .line_index
            .position_for_offset(&snapshot.text, offset);
        complete(&snapshot, position, Vec::new(), false)
            .items
            .into_iter()
            .map(|item| item.label)
            .collect()
    }

    #[test]
    fn completion_filters_by_context_and_prefix() {
        let top = labels("kp|", "|");
        assert!(top.contains(&"kprobe".to_string()));
        assert!(!top.contains(&"printf".to_string()));

        let action = labels("BEGIN { pri| }", "|");
        assert!(action.contains(&"print".to_string()));
        assert!(action.contains(&"printf".to_string()));

        let scratch = labels("BEGIN { $value = 1; print($va|); }", "|");
        assert_eq!(scratch, vec!["$value"]);

        assert!(labels("BEGIN { \"pri|\"; }", "|").is_empty());
        assert!(labels("BEGIN { args.| }", "|").is_empty());
    }
}
