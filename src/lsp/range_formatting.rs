use super::snapshot::DocumentSnapshot;
use crate::{config::Config, format::format_source, parse};
use anyhow::{bail, Result};
use std::ops::Range as ByteRange;
use tower_lsp::lsp_types::{Range, TextEdit};
use tree_sitter::Node;

pub(super) fn format_range(
    doc: &DocumentSnapshot,
    range: Range,
    config: &Config,
) -> Result<Vec<TextEdit>> {
    let source = &doc.text;
    let start = doc.line_index.offset_for_position(source, range.start);
    let end = doc.line_index.offset_for_position(source, range.end);
    if range.start > range.end
        || doc.line_index.position_for_offset(source, start) != range.start
        || doc.line_index.position_for_offset(source, end) != range.end
    {
        bail!("range must use valid UTF-16 character boundaries");
    }
    let selected = &source[start..end];
    if selected.trim().is_empty() {
        return Ok(Vec::new());
    }
    let requested = start + selected.len() - selected.trim_start().len()
        ..end - (selected.len() - selected.trim_end().len());
    let Some(tree) = &doc.tree else {
        return Ok(Vec::new());
    };
    let root = tree.root_node();
    if root.has_error() {
        return Ok(Vec::new());
    }
    let container = selection_container(root, &requested).unwrap_or(root);
    let mut cursor = container.walk();
    let children: Vec<_> = container.children(&mut cursor).collect();
    let Some(first) = children
        .iter()
        .position(|node| overlaps(node.byte_range(), &requested))
    else {
        return Ok(Vec::new());
    };
    let mut first = first;
    let mut last = children
        .iter()
        .rposition(|node| overlaps(node.byte_range(), &requested))
        .unwrap();
    if children[first].kind() == ";" && first > 0 {
        first -= 1;
    }
    if children
        .get(last + 1)
        .is_some_and(|node| node.kind() == ";")
    {
        last += 1;
    }
    // A trailing comment belongs to the statement, including when the formatter
    // puts it on a separate line according to comments.preserve_inline.
    while let Some(next) = children.get(last + 1) {
        if !is_comment(*next) || source[children[last].end_byte()..next.start_byte()].contains('\n')
        {
            break;
        }
        last += 1;
    }
    let nodes = &children[first..=last];
    if nodes.iter().all(|node| is_comment(*node))
        || nodes.iter().any(|node| contains_preprocessor(*node))
    {
        return Ok(Vec::new());
    }
    let effective = line_edges(
        source,
        nodes[0].start_byte()..nodes.last().unwrap().end_byte(),
    );
    if suppression_intersects(source, root, &effective) {
        return Ok(Vec::new());
    }

    // Formatting the actual document preserves all context (predicates, nested
    // blocks, GNU braces, tabs and comments) without a second fragment formatter.
    let formatted = format_source(source, config)?;
    let formatted_tree = parse::ensure_valid(&formatted)?;
    let formatted_root = formatted_tree.root_node();
    let Some(mapped) = nodes
        .iter()
        .map(|node| corresponding_node(*node, formatted_root))
        .collect::<Option<Vec<_>>>()
    else {
        return Ok(Vec::new());
    };
    if nodes
        .iter()
        .zip(&mapped)
        .any(|(before, after)| !same_syntax(*before, source, *after, &formatted))
    {
        return Ok(Vec::new());
    }
    let output = line_edges(
        &formatted,
        mapped[0].start_byte()..mapped.last().unwrap().end_byte(),
    );
    let mut replacement = formatted[output.clone()].to_string();
    // Splitting a one-line block must introduce a line break before the selected
    // statement, while keeping the preceding brace/statement untouched.
    if effective.start > 0
        && !source[..effective.start].ends_with('\n')
        && (output.start == 0 || formatted[..output.start].ends_with('\n'))
    {
        replacement.insert(0, '\n');
    }
    if replacement == source[effective.clone()] {
        return Ok(Vec::new());
    }
    let updated = format!(
        "{}{}{}",
        &source[..effective.start],
        replacement,
        &source[effective.end..]
    );
    if parse::ensure_valid(&updated).is_err() {
        return Ok(Vec::new());
    }
    Ok(vec![TextEdit {
        range: doc
            .line_index
            .range_for_offsets(source, effective.start, effective.end),
        new_text: replacement,
    }])
}

fn overlaps(left: ByteRange<usize>, right: &ByteRange<usize>) -> bool {
    left.start < right.end && right.start < left.end
}

fn selection_container<'a>(node: Node<'a>, range: &ByteRange<usize>) -> Option<Node<'a>> {
    let mut cursor = node.walk();
    for child in node.named_children(&mut cursor) {
        if child.start_byte() <= range.start && child.end_byte() >= range.end {
            if let Some(container) = selection_container(child, range) {
                return Some(container);
            }
        }
    }
    match node.kind() {
        "source_file" | "preamble" => Some(node),
        "action" | "block" if range.start > node.start_byte() && range.end < node.end_byte() => {
            Some(node)
        }
        _ => None,
    }
}

fn is_comment(node: Node<'_>) -> bool {
    matches!(node.kind(), "line_comment" | "block_comment")
}

fn contains_preprocessor(node: Node<'_>) -> bool {
    if matches!(
        node.kind(),
        "c_preproc" | "c_preproc_block" | "hashbang" | "c_struct"
    ) {
        return true;
    }
    let mut cursor = node.walk();
    let found = node.named_children(&mut cursor).any(contains_preprocessor);
    found
}

fn line_edges(text: &str, range: ByteRange<usize>) -> ByteRange<usize> {
    let mut start = range.start;
    let mut end = range.end;
    while start > 0 && matches!(text.as_bytes()[start - 1], b' ' | b'\t') {
        start -= 1;
    }
    while end < text.len() && matches!(text.as_bytes()[end], b' ' | b'\t' | b'\r') {
        end += 1;
    }
    if text.as_bytes().get(end) == Some(&b'\n') {
        end += 1;
    }
    start..end
}

fn corresponding_node<'a>(mut node: Node<'_>, mut root: Node<'a>) -> Option<Node<'a>> {
    let mut path = Vec::new();
    while let Some(parent) = node.parent() {
        let mut cursor = parent.walk();
        let index = parent
            .children(&mut cursor)
            .position(|child| child == node)?;
        path.push((index, parent.kind(), parent.child_count()));
        node = parent;
    }
    for (index, kind, count) in path.into_iter().rev() {
        if root.kind() != kind || root.child_count() != count {
            return None;
        }
        root = root.child(index)?;
    }
    Some(root)
}

fn same_syntax(before: Node<'_>, source: &str, after: Node<'_>, formatted: &str) -> bool {
    if before.kind() != after.kind() || before.child_count() != after.child_count() {
        return false;
    }
    if before.child_count() == 0 {
        return source[before.byte_range()].trim_end() == formatted[after.byte_range()].trim_end();
    }
    (0..before.child_count()).all(|index| {
        same_syntax(
            before.child(index).unwrap(),
            source,
            after.child(index).unwrap(),
            formatted,
        )
    })
}

// Temporary integration boundary for #43: replace this guard with the shared
// suppression-span helper once available. Only standalone parsed line comments
// count; strings and block comments containing directive text do not.
fn suppression_intersects(source: &str, root: Node<'_>, effective: &ByteRange<usize>) -> bool {
    fn collect(node: Node<'_>, source: &str, directives: &mut Vec<(ByteRange<usize>, bool)>) {
        if node.kind() == "line_comment" {
            let start = source[..node.start_byte()]
                .rfind('\n')
                .map_or(0, |offset| offset + 1);
            if source[start..node.start_byte()].trim().is_empty() {
                let text = source[node.byte_range()].trim_end();
                if matches!(text, "// btfmt: off" | "// btfmt: on") {
                    directives.push((start..node.end_byte(), text == "// btfmt: off"));
                }
            }
            return;
        }
        let mut cursor = node.walk();
        for child in node.named_children(&mut cursor) {
            collect(child, source, directives);
        }
    }
    let mut directives = Vec::new();
    collect(root, source, &mut directives);
    let mut disabled = None;
    for (range, off) in directives {
        if overlaps(range.clone(), effective) {
            return true;
        }
        if off {
            disabled.get_or_insert(range.start);
        } else if let Some(start) = disabled.take() {
            if overlaps(start..range.end, effective) {
                return true;
            }
        }
    }
    disabled.is_some_and(|start| overlaps(start..source.len(), effective))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::{config::BraceStyle, text::LineIndex};
    use tower_lsp::lsp_types::Position;

    fn edits(source: &str, selection: &str, config: &Config) -> Vec<TextEdit> {
        let start = source.find(selection).unwrap();
        let (doc, _) = DocumentSnapshot::analyze(source.to_string(), Some(1));
        format_range(
            &doc,
            doc.line_index
                .range_for_offsets(source, start, start + selection.len()),
            config,
        )
        .unwrap()
    }

    fn apply(source: &str, edits: &[TextEdit]) -> String {
        assert!(edits.len() <= 1);
        let Some(edit) = edits.first() else {
            return source.to_string();
        };
        let index = LineIndex::new(source);
        let start = index.offset_for_position(source, edit.range.start);
        let end = index.offset_for_position(source, edit.range.end);
        let result = format!("{}{}{}", &source[..start], edit.new_text, &source[end..]);
        assert!(result.starts_with(&source[..start]));
        assert!(result.ends_with(&source[end..]));
        parse::ensure_valid(&result).unwrap();
        result
    }

    #[test]
    fn partial_expression_expands_to_statement_and_preserves_other_bytes() {
        let source = "BEGIN\n{\n $x=1+2;\n  $y=  3;\n}\nEND{exit();}\n";
        let changes = edits(source, "1+", &Config::default());
        assert_eq!(
            changes[0].range,
            Range::new(Position::new(2, 0), Position::new(3, 0))
        );
        let result = apply(source, &changes);
        assert_eq!(
            result,
            "BEGIN\n{\n    $x = 1 + 2;\n  $y=  3;\n}\nEND{exit();}\n"
        );
        assert!(edits(&result, "$x = 1 + 2;", &Config::default()).is_empty());
    }

    #[test]
    fn nested_statements_match_full_formatting_for_all_indent_modes() {
        let source = "BEGIN { if(pid){ $x=1; $y=2; } else { $z=3; } }\n";
        for style in [BraceStyle::SameLine, BraceStyle::NextLine, BraceStyle::Gnu] {
            for spaces in [true, false] {
                for indent in [true, false] {
                    let mut config = Config::default();
                    config.blocks.brace_style = style.clone();
                    config.blocks.indent_statements = indent;
                    config.indent.use_spaces = spaces;
                    config.indent.size = 2;
                    let full = format_source(source, &config).unwrap();
                    let changes = edits(source, "$x=1; $y=2;", &config);
                    assert_eq!(changes.len(), 1, "{config:?}");
                    let result = apply(source, &changes);
                    let expected: String = full
                        .lines()
                        .filter(|line| line.contains("$x") || line.contains("$y"))
                        .map(|line| format!("{line}\n"))
                        .collect();
                    assert_eq!(changes[0].new_text, format!("\n{expected}"));
                    assert!(result.starts_with("BEGIN { if(pid){\n"));
                    assert!(result.ends_with("} else { $z=3; } }\n"));
                    assert!(edits(&result, "$x = 1;", &config).is_empty());
                }
            }
        }
    }

    #[test]
    fn complete_probes_and_control_flow_match_full_formatter() {
        let config = Config::default();
        for source in [
            "BEGIN{if(pid){$x=1;}else{$x=2;}while($x<3){$x++;}}\n",
            "kprobe:vfs_read,kprobe:vfs_write /pid==42/ { @x=count(); }\n",
            "macro plus(value) { value+1 }\n",
        ] {
            let changes = edits(source, source.trim(), &config);
            assert_eq!(
                apply(source, &changes),
                format_source(source, &config).unwrap()
            );
        }
        let source = "BEGIN { if(pid){$x=1;}else{$x=2;} $y=3; }\n";
        let result = apply(source, &edits(source, "pid", &config));
        assert!(result.contains("    if (pid)\n    {\n        $x = 1;\n    }\n    else"));
        assert!(result.ends_with("$y=3; }\n"));
    }

    #[test]
    fn multiple_probes_keep_outer_separators_and_half_open_selection_boundaries() {
        let config = Config::default();
        let source = "// unchanged\n\nBEGIN{$x=1;}\n\n\nEND{$y=2;}\n\n\nkprobe:f{$z=3;}\n";
        let changes = edits(source, "BEGIN{$x=1;}\n\n\nEND", &config);
        let result = apply(source, &changes);
        let expected = format_source("BEGIN{$x=1;}\nEND{$y=2;}\n", &config).unwrap();
        assert_eq!(
            result,
            format!("// unchanged\n\n{expected}\n\nkprobe:f{{$z=3;}}\n")
        );
        let source = "BEGIN\n{\n $x=1;\n $y=2;\n}\n";
        let changes = edits(source, " $x=1;\n", &config);
        assert_eq!(
            apply(source, &changes),
            "BEGIN\n{\n    $x = 1;\n $y=2;\n}\n"
        );
    }

    #[test]
    fn unicode_comments_and_crlf_are_preserved_outside_effective_range() {
        let source = "// 中文😀\r\nBEGIN { printf(\"😀\"); $x=1; // 注释😀\r\n $y= 2;\r\n}\r\n";
        for inline in [true, false] {
            let mut config = Config::default();
            config.comments.preserve_inline = inline;
            let changes = edits(source, "$x", &config);
            assert_eq!(changes[0].range.start, Position::new(1, 21));
            assert_eq!(changes[0].range.end, Position::new(2, 0));
            let result = apply(source, &changes);
            assert!(result.starts_with("// 中文😀\r\nBEGIN { printf(\"😀\");\n"));
            assert!(result.ends_with(" $y= 2;\r\n}\r\n"));
            assert!(result.contains("// 注释😀\n"));
            assert!(edits(&result, "$x = 1;", &config).is_empty());
        }
    }

    #[test]
    fn comments_between_statements_and_partial_block_comments_are_safe() {
        let config = Config::default();
        let source = "BEGIN { $x=1; /* 中文😀\n comment */ $y=2; // tail\n}\n";
        let result = apply(
            source,
            &edits(source, "$x=1; /* 中文😀\n comment */ $y", &config),
        );
        assert!(result.contains("/* 中文😀\n comment */"));
        assert!(result.contains("$y = 2; // tail"));
        assert!(edits(source, "中文", &config).is_empty());
        assert!(edits("// lone comment\nBEGIN{}\n", "comment", &config).is_empty());
    }

    #[test]
    fn empty_malformed_preprocessor_and_invalid_positions_are_predictable() {
        let config = Config::default();
        for (source, selection) in [
            ("BEGIN{$x=1;}\n", ""),
            ("BEGIN {\n  $x=1;\n}\n", "  "),
            ("BEGIN { $x=; }", "$x"),
            ("#define VALUE 1\nBEGIN{$x=VALUE;}\n", "VALUE 1"),
            ("#if X\n#define VALUE 1\n#endif\nBEGIN{}\n", "#if"),
        ] {
            assert!(edits(source, selection, &config).is_empty());
        }
        let source = "#define VALUE 1\nBEGIN{$x=VALUE;}\n";
        assert!(!edits(source, "$x", &config).is_empty());
        let (doc, _) = DocumentSnapshot::analyze("BEGIN{printf(\"😀\");}".into(), Some(1));
        for range in [
            Range::new(Position::new(0, 15), Position::new(0, 16)),
            Range::new(Position::new(4, 0), Position::new(4, 0)),
            Range::new(Position::new(0, 10), Position::new(0, 1)),
        ] {
            assert!(format_range(&doc, range, &config).is_err());
        }
    }

    #[test]
    fn suppression_guard_blocks_disabled_ranges_but_not_enabled_neighbors() {
        let config = Config::default();
        let source = "BEGIN {\n $a=1;\n // btfmt: off\n $b=2;\n // btfmt: on\n $c=3;\n}\n";
        for selection in ["$b", "btfmt: off", "BEGIN", "$a=1;\n // btfmt: off\n $b"] {
            assert!(edits(source, selection, &config).is_empty());
        }
        for selection in ["$a", "$c"] {
            assert!(!edits(source, selection, &config).is_empty());
        }
        assert!(edits("BEGIN{\n// btfmt: off\n$x=1;\n}", "$x", &config).is_empty());
        let source = "BEGIN{ /* // btfmt: off */ $x=1;}";
        assert!(!edits(source, "$x", &config).is_empty());
    }
}
