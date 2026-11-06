use std::collections::HashMap;

use super::catalog::{self, CatalogKind};
use super::snapshot::DocumentSnapshot;
use super::symbols::CompletionSymbolKind;
use tower_lsp::lsp_types::{
    SemanticToken, SemanticTokenModifier, SemanticTokenType, SemanticTokens,
    SemanticTokensFullOptions, SemanticTokensLegend, SemanticTokensOptions,
};
use tree_sitter::Node;

#[derive(Clone, Copy)]
#[repr(u32)]
enum Kind {
    Variable,
    Map,
    Parameter,
    Macro,
    Function,
    Property,
    Type,
    Event,
    Comment,
    String,
    Number,
    Keyword,
}

const DECLARATION: u32 = 1;
const DEFAULT_LIBRARY: u32 = 1 << 1;

pub(super) fn options() -> SemanticTokensOptions {
    SemanticTokensOptions {
        legend: SemanticTokensLegend {
            token_types: vec![
                SemanticTokenType::VARIABLE,
                SemanticTokenType::new("map"),
                SemanticTokenType::PARAMETER,
                SemanticTokenType::MACRO,
                SemanticTokenType::FUNCTION,
                SemanticTokenType::PROPERTY,
                SemanticTokenType::TYPE,
                SemanticTokenType::EVENT,
                SemanticTokenType::COMMENT,
                SemanticTokenType::STRING,
                SemanticTokenType::NUMBER,
                SemanticTokenType::KEYWORD,
            ],
            token_modifiers: vec![
                SemanticTokenModifier::DECLARATION,
                SemanticTokenModifier::DEFAULT_LIBRARY,
            ],
        },
        full: Some(SemanticTokensFullOptions::Bool(true)),
        range: Some(false),
        ..SemanticTokensOptions::default()
    }
}

pub(super) fn full(doc: &DocumentSnapshot) -> SemanticTokens {
    let mut result = SemanticTokens::default();
    let Some(tree) = &doc.tree else {
        return result;
    };
    // The recovery index belongs to this exact snapshot, even while typing invalid syntax.
    let symbols: HashMap<_, _> = doc
        .completion_symbols
        .iter()
        .flat_map(|index| index.classified_occurrences())
        .map(|(range, kind, declaration)| {
            let kind = match kind {
                CompletionSymbolKind::Scratch => Kind::Variable,
                CompletionSymbolKind::Map => Kind::Map,
                CompletionSymbolKind::Parameter => Kind::Parameter,
                CompletionSymbolKind::Macro => Kind::Macro,
            };
            (
                range.start,
                (kind, if declaration { DECLARATION } else { 0 }),
            )
        })
        .collect();
    let mut spans = Vec::new();
    let mut pending = vec![tree.root_node()];
    while let Some(node) = pending.pop() {
        if node.is_missing() || node.start_byte() == node.end_byte() {
            continue;
        }
        let text = &doc.text[node.byte_range()];
        if node.kind() == "map_variable" {
            // A map node includes its indexes; only color the sigil and name here.
            let len = text
                .bytes()
                .take_while(|byte| byte.is_ascii_alphanumeric() || matches!(byte, b'@' | b'_'))
                .count();
            let (kind, modifiers) = symbols
                .get(&node.start_byte())
                .copied()
                .unwrap_or((Kind::Map, 0));
            spans.push((node.start_byte(), node.start_byte() + len, kind, modifiers));
        } else if let Some((kind, modifiers)) = classify(node, text, &symbols) {
            spans.push((node.start_byte(), node.end_byte(), kind, modifiers));
            continue;
        }
        let mut cursor = node.walk();
        pending.extend(node.children(&mut cursor));
    }
    spans.sort_by_key(|&(start, end, _, _)| (start, end));
    let mut previous_end = 0;
    let mut previous_line = 0;
    let mut previous_character = 0;
    for (start, end, kind, modifiers) in spans {
        if start < previous_end || start >= end {
            continue;
        }
        previous_end = end;
        // Never require multilineTokenSupport; exclude both CR and LF from each segment.
        let mut offset = start;
        for segment in doc.text[start..end].split_inclusive('\n') {
            let content = segment.trim_end_matches(['\r', '\n']);
            if !content.is_empty() {
                let position = doc.line_index.position_for_offset(&doc.text, offset);
                result.data.push(SemanticToken {
                    delta_line: position.line - previous_line,
                    delta_start: if position.line == previous_line {
                        position.character - previous_character
                    } else {
                        position.character
                    },
                    length: content.encode_utf16().count() as u32,
                    token_type: kind as u32,
                    token_modifiers_bitset: modifiers,
                });
                previous_line = position.line;
                previous_character = position.character;
            }
            offset += segment.len();
        }
    }
    result
}

fn classify(
    node: Node<'_>,
    text: &str,
    symbols: &HashMap<usize, (Kind, u32)>,
) -> Option<(Kind, u32)> {
    let kind = match node.kind() {
        "line_comment" | "block_comment" | "hashbang" => Kind::Comment,
        "string_literal" => Kind::String,
        "integer_literal" | "decimal_number" => Kind::Number,
        "boolean_literal" | "break_statement" | "continue_statement" | "return_statement" => {
            Kind::Keyword
        }
        "integer_type" | "string_type" => Kind::Type,
        "script_parameter" => return Some((Kind::Parameter, DEFAULT_LIBRARY)),
        "args_keyword" | "retval_identifier" => return Some((Kind::Variable, DEFAULT_LIBRARY)),
        "scratch_variable" => {
            return Some(
                symbols
                    .get(&node.start_byte())
                    .copied()
                    .unwrap_or((Kind::Variable, 0)),
            )
        }
        "identifier" => {
            let parent = node.parent()?;
            if field_is(parent, "field", node) || parent.kind() == "offsetof_expression" {
                return Some((Kind::Property, 0));
            }
            if matches!(
                parent.kind(),
                "struct_type" | "type_specifier" | "pointer_type"
            ) {
                return Some((Kind::Type, 0));
            }
            if matches!(parent.kind(), "probe" | "config_assignment") {
                return None;
            }
            if let Some(class) = symbols.get(&node.start_byte()) {
                return Some(*class);
            }
            let entry = catalog::find(text)?;
            return match entry.kind {
                CatalogKind::Function
                    if parent.kind() == "call_expression" && field_is(parent, "function", node) =>
                {
                    Some((Kind::Function, DEFAULT_LIBRARY))
                }
                CatalogKind::Value if parent.kind() != "call_expression" => {
                    Some((Kind::Variable, DEFAULT_LIBRARY))
                }
                _ => None,
            };
        }
        kind if kind.ends_with("_provider") => return Some((Kind::Event, DEFAULT_LIBRARY)),
        _ if !node.is_named()
            && matches!(
                text,
                "let"
                    | "macro"
                    | "if"
                    | "else"
                    | "while"
                    | "for"
                    | "unroll"
                    | "comptime"
                    | "import"
                    | "config"
                    | "struct"
                    | "sizeof"
                    | "offsetof"
                    | "typeof"
            ) =>
        {
            Kind::Keyword
        }
        _ => return None,
    };
    Some((kind, 0))
}

fn field_is(parent: Node<'_>, field: &str, node: Node<'_>) -> bool {
    parent.child_by_field_name(field) == Some(node)
}

#[cfg(test)]
mod tests {
    use super::*;
    use tower_lsp::lsp_types::Position;

    fn classified(text: &str) -> Vec<(String, String, u32, Position)> {
        let (doc, _) = DocumentSnapshot::analyze(text.into(), Some(1));
        let legend = options().legend;
        let mut position = Position::default();
        let mut previous_end = Position::default();
        full(&doc)
            .data
            .into_iter()
            .map(|token| {
                position.line += token.delta_line;
                position.character = if token.delta_line == 0 {
                    position.character + token.delta_start
                } else {
                    token.delta_start
                };
                assert!(position >= previous_end);
                assert!(token.length > 0);
                let end = Position::new(position.line, position.character + token.length);
                let start_offset = doc.line_index.offset_for_position(text, position);
                let end_offset = doc.line_index.offset_for_position(text, end);
                assert_eq!(doc.line_index.position_for_offset(text, end_offset), end);
                let value = &text[start_offset..end_offset];
                assert!(!value.contains(['\n', '\r']));
                assert_eq!(value.encode_utf16().count(), token.length as usize);
                assert!(token.token_modifiers_bitset < 4);
                previous_end = end;
                (
                    value.into(),
                    legend.token_types[token.token_type as usize]
                        .as_str()
                        .into(),
                    token.token_modifiers_bitset,
                    position,
                )
            })
            .collect()
    }

    #[test]
    fn classifies_bindings_macro_parameters_and_builtins() {
        let text = concat!(
            "macro apply(value, $local, @target) { @target[value] = $local; value }\n",
            "BEGIN { let $local = pid; @total[$local] = count(); ",
            "apply($local, $local, @total); printf(\"%d\", @total[$local]); }\n",
        );
        crate::parse::ensure_valid(text).unwrap();
        let tokens = classified(text);
        let classes = |name: &str| -> Vec<_> {
            tokens
                .iter()
                .filter(|token| token.0 == name)
                .map(|token| (token.1.as_str(), token.2))
                .collect()
        };
        assert_eq!(classes("apply"), [("macro", DECLARATION), ("macro", 0)]);
        assert_eq!(
            classes("value"),
            [
                ("parameter", DECLARATION),
                ("parameter", 0),
                ("parameter", 0)
            ]
        );
        assert_eq!(
            classes("@target"),
            [("parameter", DECLARATION), ("parameter", 0)]
        );
        assert_eq!(
            classes("$local"),
            [
                ("parameter", DECLARATION),
                ("parameter", 0),
                ("variable", DECLARATION),
                ("variable", 0),
                ("variable", 0),
                ("variable", 0),
                ("variable", 0),
            ]
        );
        assert_eq!(classes("@total"), [("map", 0); 3]);
        assert_eq!(classes("pid"), [("variable", DEFAULT_LIBRARY)]);
        assert_eq!(classes("count"), [("function", DEFAULT_LIBRARY)]);
        assert_eq!(classes("printf"), [("function", DEFAULT_LIBRARY)]);
    }

    #[test]
    fn uses_syntax_context_and_local_shadowing_not_identifier_spelling() {
        let text = concat!(
            "macro use(pid, printf) { printf(pid); $x = args.pid; $y = (struct pid *)0; pid }\n",
            "uprobe:/bin/test:printf { unknown(); $x = args.pid; printf(\"pid @map $x\"); }\n",
        );
        crate::parse::ensure_valid(text).unwrap();
        let tokens = classified(text);
        let classes = |name: &str| -> Vec<_> {
            tokens
                .iter()
                .filter(|token| token.0 == name)
                .map(|token| token.1.as_str())
                .collect()
        };
        assert_eq!(
            classes("pid"),
            [
                "parameter",
                "parameter",
                "property",
                "type",
                "parameter",
                "property"
            ]
        );
        assert_eq!(classes("printf"), ["parameter", "parameter", "function"]);
        assert!(classes("unknown").is_empty());
        assert!(classes("@map").is_empty());
    }

    #[test]
    fn encodes_unicode_crlf_and_multiline_comments_without_overlap() {
        let text = "/*😀 first\r\n\r\n終 */\r\nBEGIN { printf(\"😀中\"); @x[pid] = count(); }\r\n";
        crate::parse::ensure_valid(text).unwrap();
        let tokens = classified(text);
        assert_eq!(
            tokens[0],
            (
                "/*😀 first".into(),
                "comment".into(),
                0,
                Position::new(0, 0)
            )
        );
        assert_eq!(
            tokens[1],
            ("終 */".into(), "comment".into(), 0, Position::new(2, 0))
        );
        let map = tokens.iter().find(|token| token.0 == "@x").unwrap();
        assert_eq!(map.3, Position::new(3, 23));
        assert!(tokens
            .iter()
            .any(|token| token.0 == "\"😀中\"" && token.1 == "string"));
    }

    #[test]
    fn map_tokens_exclude_trivia_and_indexes_without_losing_parameter_bindings() {
        let text = concat!(
            "macro m(@target) { @target /*comment*/ [pid] = 1; }\n",
            "BEGIN { @ [pid] = count(); @named [tid] = 1; }\n",
        );
        crate::parse::ensure_valid(text).unwrap();
        let tokens = classified(text);
        let classes: Vec<_> = tokens
            .iter()
            .filter(|token| token.0.starts_with('@'))
            .map(|token| (token.0.as_str(), token.1.as_str(), token.2))
            .collect();
        assert_eq!(
            classes,
            [
                ("@target", "parameter", DECLARATION),
                ("@target", "parameter", 0),
                ("@", "map", 0),
                ("@named", "map", 0),
            ]
        );
        assert!(tokens
            .iter()
            .any(|token| token.0 == "/*comment*/" && token.1 == "comment"));
        assert!(tokens
            .iter()
            .any(|token| token.0 == "pid" && token.2 == DEFAULT_LIBRARY));
    }

    #[test]
    fn incomplete_edits_always_use_current_snapshot_and_valid_ranges() {
        let text = "macro m(value, @target) { @target[value] = count(); }\nBEGIN { printf(\"😀\"); m(pid, @x); }\n";
        for (end, _) in text
            .char_indices()
            .chain(std::iter::once((text.len(), '\0')))
        {
            classified(&text[..end]);
        }
        let (old, _) = DocumentSnapshot::analyze("BEGIN { @old = count(); }".into(), Some(1));
        assert!(!full(&old).data.is_empty());
        let (empty, _) = DocumentSnapshot::analyze(String::new(), Some(2));
        assert!(full(&empty).data.is_empty());
        let tokens = classified("BEGIN { @new = count(); $unfinished =");
        assert!(tokens.iter().any(|token| token.0 == "@new"));
        assert!(!tokens.iter().any(|token| token.0 == "@old"));
    }
}
