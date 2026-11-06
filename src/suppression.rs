//! Detection of whole-line formatting suppression directives.

use crate::parse;
use anyhow::Result;
use std::ops::Range;

/// Return sorted, disjoint, half-open UTF-8 byte ranges to preserve verbatim.
///
/// Ranges include the complete `// btfmt: off` and `// btfmt: on` lines,
/// including indentation, trailing whitespace and line endings. An unmatched
/// `off` extends to EOF; repeated `off` directives do not nest, and unmatched
/// `on` directives are ignored. Only standalone line comments count, excluding
/// comments inside preprocessor nodes. The entire source must parse successfully.
///
/// Range-formatting callers must detect these ranges on the original document,
/// before slicing it: an `off` may precede the requested selection.
pub fn suppression_ranges(source: &str) -> Result<Vec<Range<usize>>> {
    let tree = parse::ensure_valid(source)?;
    Ok(ranges_from_tree(source, tree.root_node()))
}

pub(crate) fn ranges_from_tree(source: &str, root: tree_sitter::Node<'_>) -> Vec<Range<usize>> {
    fn collect(
        source: &str,
        node: tree_sitter::Node<'_>,
        start: &mut Option<usize>,
        ranges: &mut Vec<Range<usize>>,
    ) {
        match node.kind() {
            "line_comment" => {
                let line_start = line_start(source, node.start_byte());
                if !is_standalone(source, node.start_byte()) {
                    return;
                }
                let line_end = source[node.end_byte()..]
                    .find('\n')
                    .map_or(source.len(), |offset| node.end_byte() + offset + 1);
                match directive(&source[node.start_byte()..line_end]) {
                    Some(Directive::Off) => {
                        start.get_or_insert(line_start);
                    }
                    Some(Directive::On) => {
                        if let Some(begin) = start.take() {
                            ranges.push(begin..line_end);
                        }
                    }
                    None => {}
                }
                return;
            }
            "string_literal" | "block_comment" | "c_preproc" | "c_preproc_block" => return,
            _ => {}
        }
        let mut cursor = node.walk();
        for child in node.named_children(&mut cursor) {
            collect(source, child, start, ranges);
        }
    }

    let mut ranges = Vec::new();
    let mut start = None;
    collect(source, root, &mut start, &mut ranges);
    if let Some(start) = start {
        ranges.push(start..source.len());
    }
    ranges
}

#[derive(Clone, Copy)]
pub(crate) enum Directive {
    Off,
    On,
}

pub(crate) fn directive(comment: &str) -> Option<Directive> {
    match comment.trim_end() {
        "// btfmt: off" => Some(Directive::Off),
        "// btfmt: on" => Some(Directive::On),
        _ => None,
    }
}

pub(crate) fn is_standalone(source: &str, comment_start: usize) -> bool {
    source[line_start(source, comment_start)..comment_start]
        .chars()
        .all(char::is_whitespace)
}

fn line_start(source: &str, offset: usize) -> usize {
    source[..offset].rfind('\n').map_or(0, |pos| pos + 1)
}
