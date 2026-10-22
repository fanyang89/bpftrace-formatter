use super::{catalog, snapshot::DocumentSnapshot};
use tower_lsp::lsp_types::{
    Documentation, ParameterInformation, ParameterLabel, Position, SignatureHelp,
    SignatureInformation,
};

pub(super) fn signature_help(doc: &DocumentSnapshot, position: Position) -> Option<SignatureHelp> {
    let offset = doc.line_index.offset_for_position(&doc.text, position);
    if let Some(tree) = &doc.tree {
        let mut node = tree.root_node().descendant_for_byte_range(offset, offset)?;
        loop {
            if matches!(node.kind(), "c_preproc" | "c_preproc_block") {
                return None;
            }
            match node.parent() {
                Some(parent) => node = parent,
                None => break,
            }
        }
    }
    let (name, argument) = active_call(&doc.text[..offset])?;
    let entry = catalog::find(name)?;
    let signature = catalog::signature(name)?;
    let count = signature.parameters.len();
    // Do not highlight an unrelated parameter in an invalid, overfull call.
    if (count == 0 && argument > 0) || (count > 0 && argument >= count && !signature.variadic) {
        return None;
    }
    let active_parameter = (count > 0).then(|| argument.min(count - 1) as u32);
    let label = format!(
        "{}({})",
        name,
        signature
            .parameters
            .iter()
            .map(|(label, _)| *label)
            .collect::<Vec<_>>()
            .join(", ")
    );
    Some(SignatureHelp {
        signatures: vec![SignatureInformation {
            label,
            documentation: Some(Documentation::String(entry.documentation.into())),
            parameters: Some(
                signature
                    .parameters
                    .iter()
                    .map(|(label, doc)| ParameterInformation {
                        label: ParameterLabel::Simple((*label).into()),
                        documentation: Some(Documentation::String((*doc).into())),
                    })
                    .collect(),
            ),
            active_parameter,
        }],
        active_signature: Some(0),
        active_parameter,
    })
}

struct Delimiter<'a> {
    close: u8,
    callee: Option<&'a str>,
    argument: usize,
}

// Scan balanced syntax rather than relying on recovered call nodes: tree-sitter may
// not produce a call_expression at all until the closing parenthesis is entered.
fn active_call(text: &str) -> Option<(&str, usize)> {
    let bytes = text.as_bytes();
    let mut stack: Vec<Delimiter<'_>> = Vec::new();
    let mut previous = "";
    let mut callee = None;
    let mut i = 0;
    while i < bytes.len() {
        match bytes[i] {
            b if b.is_ascii_whitespace() => i += 1,
            b'/' if bytes.get(i + 1) == Some(&b'/') => {
                i += 2;
                while i < bytes.len() && bytes[i] != b'\n' {
                    i += 1;
                }
                if i == bytes.len() {
                    return None;
                }
            }
            b'/' if bytes.get(i + 1) == Some(&b'*') => {
                i += 2;
                while i + 1 < bytes.len() && &bytes[i..i + 2] != b"*/" {
                    i += 1;
                }
                if i + 1 >= bytes.len() {
                    return None;
                }
                i += 2;
            }
            b'#' => {
                // Preprocessor directives and shebangs cannot contain bpftrace calls.
                while i < bytes.len() {
                    let continuation = bytes[i] == b'\n' && i > 0 && bytes[i - 1] == b'\\';
                    if bytes[i] == b'\n' && !continuation {
                        break;
                    }
                    i += 1;
                }
                if i == bytes.len() {
                    return None;
                }
                previous = "";
                callee = None;
            }
            quote @ (b'"' | b'\'') => {
                i += 1;
                while i < bytes.len() {
                    if bytes[i] == b'\\' {
                        i = (i + 2).min(bytes.len());
                    } else if bytes[i] == quote {
                        i += 1;
                        break;
                    } else {
                        i += 1;
                    }
                }
                previous = "";
                callee = None;
            }
            b if b.is_ascii_alphanumeric() || matches!(b, b'_' | b'$' | b'@') => {
                let start = i;
                while i < bytes.len()
                    && (bytes[i].is_ascii_alphanumeric() || matches!(bytes[i], b'_' | b'$' | b'@'))
                {
                    i += 1;
                }
                let word = &text[start..i];
                callee = if matches!(previous, "." | ">" | "macro" | "fn") {
                    // Keep a non-builtin sentinel so nested member calls do not
                    // accidentally display the outer function's signature.
                    Some("")
                } else {
                    Some(word)
                };
                previous = word;
            }
            open @ (b'(' | b'[' | b'{') => {
                stack.push(Delimiter {
                    close: match open {
                        b'(' => b')',
                        b'[' => b']',
                        _ => b'}',
                    },
                    callee: if open == b'(' { callee } else { None },
                    argument: 0,
                });
                callee = None;
                previous = "";
                i += 1;
            }
            close @ (b')' | b']' | b'}') => {
                if stack.last().is_some_and(|frame| frame.close == close) {
                    stack.pop();
                } else {
                    stack.clear();
                }
                callee = None;
                previous = "";
                i += 1;
            }
            b',' => {
                if let Some(frame) = stack.last_mut() {
                    frame.argument += 1;
                }
                callee = None;
                previous = "";
                i += 1;
            }
            b';' => {
                // Recover at statement boundaries rather than retaining an old,
                // unfinished call into the next statement.
                while stack.last().is_some_and(|frame| frame.close != b'}') {
                    stack.pop();
                }
                callee = None;
                previous = "";
                i += 1;
            }
            byte => {
                previous = match byte {
                    b'.' => ".",
                    b'>' => ">",
                    _ => "",
                };
                callee = None;
                i += 1;
            }
        }
    }
    stack
        .iter()
        .rev()
        .find_map(|frame| frame.callee.map(|name| (name, frame.argument)))
}

#[cfg(test)]
mod tests {
    use super::*;

    fn help(source: &str) -> Option<SignatureHelp> {
        let offset = source.find('|').expect("cursor marker");
        let text = source.replacen('|', "", 1);
        let (doc, _) = DocumentSnapshot::analyze(text, Some(1));
        signature_help(&doc, doc.line_index.position_for_offset(&doc.text, offset))
    }

    #[test]
    fn resolves_nested_and_incomplete_calls() {
        for (source, name, parameter) in [
            ("BEGIN { printf(|", "printf", 0),
            ("BEGIN { @h = hist(42, |); }", "hist", 1),
            ("BEGIN { print(@m, 10, |); }", "print", 2),
            ("BEGIN { printf(\"%s\", str(arg0, |", "str", 1),
            ("BEGIN { printf(\"%s\", str(arg0), |); }", "printf", 1),
            ("BEGIN { lhist((1, 2), 0, 100, |); }", "lhist", 3),
            ("BEGIN { lhist(@m[1, 2], |); }", "lhist", 1),
            ("BEGIN { printf /* comment */ (\"x\", |", "printf", 1),
            ("BEGIN { printf(\"😀,\\\"x\", |); }", "printf", 1),
            ("BEGIN { printf(\"still, typing|", "printf", 0),
            ("BEGIN { lhist(1, /* , ) */ 0,\n |); }", "lhist", 2),
        ] {
            let result = help(source).unwrap_or_else(|| panic!("{source}"));
            assert!(
                result.signatures[0].label.starts_with(&format!("{name}(")),
                "{source}"
            );
            assert_eq!(result.active_parameter, Some(parameter), "{source}");
            assert!(result.signatures[0]
                .parameters
                .as_ref()
                .unwrap()
                .iter()
                .all(|p| p.documentation.is_some()));
        }
    }

    #[test]
    fn suppresses_non_calls_and_stale_context() {
        for source in [
            "BEGIN { | }",
            "BEGIN { printf(\"x\")|; }",
            "BEGIN { printf(\"x\", unknown(|)); }",
            "BEGIN { $printf(|",
            "BEGIN { @printf(|",
            "BEGIN { args.printf(|",
            "BEGIN { args->printf(|",
            "BEGIN { printf(\"x\", /* comment |",
            "BEGIN { printf( // comment |",
            "#define TEST printf(|",
            "// printf(|",
            "/* printf(|",
            "BEGIN { hist(1, 2, |",
            "BEGIN { printf(; |",
            "BEGIN { count(1, |",
            "macro printf(|",
            "fn printf(|",
        ] {
            assert!(help(source).is_none(), "{source}");
        }
    }

    #[test]
    fn tolerates_every_cursor_position_in_malformed_unicode_source() {
        let text = "BEGIN { printf(\"😀\", str(arg0, 8), @m[1, 2]);\r\n/* x */ lhist(1, (2, 3], ";
        let (doc, _) = DocumentSnapshot::analyze(text.into(), Some(1));
        for offset in text
            .char_indices()
            .map(|(offset, _)| offset)
            .chain([text.len()])
        {
            let _ = signature_help(&doc, doc.line_index.position_for_offset(text, offset));
        }
    }

    #[test]
    fn handles_zero_arguments_and_variadic_arguments() {
        let result = help("BEGIN { @c = count(|); }").unwrap();
        assert_eq!(result.signatures[0].label, "count()");
        assert_eq!(result.active_parameter, None);
        assert!(result.signatures[0].parameters.as_ref().unwrap().is_empty());
        let result = help("BEGIN { printf(\"%d %d %d\", 1, 2, |); }").unwrap();
        assert_eq!(result.active_parameter, Some(1));
    }
}
