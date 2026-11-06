use btfmt::config::{BraceStyle, Config};
use btfmt::format::format_source;
use btfmt::suppression::suppression_ranges;

fn check(source: &str, config: &Config) -> String {
    let output = format_source(source, config).unwrap();
    assert_eq!(
        format_source(&output, config).unwrap(),
        output,
        "not idempotent"
    );
    let before = suppression_ranges(source).unwrap();
    let after = suppression_ranges(&output).unwrap();
    assert_eq!(
        before
            .iter()
            .map(|range| &source[range.clone()])
            .collect::<Vec<_>>(),
        after
            .iter()
            .map(|range| &output[range.clone()])
            .collect::<Vec<_>>()
    );
    output
}

#[test]
fn multiple_regions_preserve_every_byte_and_format_surroundings() {
    let source = concat!(
        "BEGIN{@before=1;\n",
        " \t// btfmt: off \t\r\n",
        "\t@x=  2;  \r\n\r\n",
        "  // btfmt: on\t\r\n",
        "@middle=3;\n",
        "// btfmt: off\n",
        "  @y=4; /* 保留 😀 */ \t\n",
        "// btfmt: on\n",
        "@after=5;}\n",
    );
    let output = check(source, &Config::default());
    assert!(output.starts_with("BEGIN\n{\n    @before = 1;\n"));
    assert!(output.contains("    @middle = 3;\n"));
    assert!(output.ends_with("    @after = 5;\n}\n"));
}

#[test]
fn unmatched_off_preserves_eof_whitespace_and_missing_newline() {
    for tail in ["", "\n", "  \t", "\n\n\t ", "\r\n\r\n"] {
        let source = format!("BEGIN{{exit();}}\n// btfmt: off\nEND{{ exit( ); }}{tail}");
        let output = check(&source, &Config::default());
        assert!(output.starts_with("BEGIN\n{\n    exit();\n}\n"));
        assert!(output.ends_with(&format!("// btfmt: off\nEND{{ exit( ); }}{tail}")));
    }
    for source in [
        "// btfmt: off",
        "// btfmt: off  \t",
        "// btfmt: off\n\n",
        "// btfmt: off\u{a0}",
        "\u{b}// btfmt: off",
    ] {
        assert_eq!(check(source, &Config::default()), source);
    }
}

#[test]
fn repeated_off_is_not_nested_and_unmatched_on_is_ignored() {
    let source = concat!(
        "// btfmt: on\n",
        "// btfmt: off\n",
        "BEGIN{exit();}\n",
        "// btfmt: off\n",
        "// btfmt: on\n",
        "END{exit();}\n",
        "// btfmt: on\n",
    );
    let ranges = suppression_ranges(source).unwrap();
    assert_eq!(ranges.len(), 1);
    assert_eq!(
        &source[ranges[0].clone()],
        "// btfmt: off\nBEGIN{exit();}\n// btfmt: off\n// btfmt: on\n"
    );
    assert!(check(source, &Config::default()).contains("END\n{\n    exit();\n}\n"));
}

#[test]
fn directive_lookalikes_do_not_trigger_even_after_formatting() {
    let source = concat!(
        "/*\n// btfmt: off\n*/\n",
        "// explanation: // btfmt: off\n",
        "// btfmt: off (not a directive)\n",
        "//btfmt: off\n",
        "// btfmt: OFF\n",
        "BEGIN{ // btfmt: off\n",
        "printf(\"text // btfmt: off\"); // btfmt: off\n",
        "printf(\"text // btfmt: on\");\n",
        "printf(\"\n// btfmt: off\n\");\n",
        "} // btfmt: off\n",
        "END{exit();}\n",
    );
    for preserve_inline in [true, false] {
        let mut config = Config::default();
        config.comments.preserve_inline = preserve_inline;
        assert!(suppression_ranges(source).unwrap().is_empty());
        let output = check(source, &config);
        assert!(output.contains("END\n{\n    exit();\n}\n"));
    }
}

#[test]
fn suppression_can_cross_block_and_expression_boundaries() {
    let source = concat!(
        "BEGIN{\n",
        "// btfmt: off\n",
        " if(1){\n",
        "// btfmt: on\n",
        "@x=1;\n",
        "// btfmt: off\n",
        " } }\nEND{printf(\n",
        "// btfmt: on\n",
        "\"x\",1);}\n",
    );
    for brace_style in [BraceStyle::SameLine, BraceStyle::NextLine, BraceStyle::Gnu] {
        for indent_statements in [true, false] {
            let mut config = Config::default();
            config.blocks.brace_style = brace_style.clone();
            config.blocks.indent_statements = indent_statements;
            let output = check(source, &config);
            assert!(output.contains("@x = 1;\n"));
            assert!(output.contains("\"x\", 1);\n"));
        }
    }
}

#[test]
fn adjacent_and_empty_regions_and_on_at_eof_are_preserved() {
    for ending in ["", "\n", " \t", "\r\n"] {
        let source = format!(
            "// btfmt: off\n// btfmt: on\n// btfmt: off\nBEGIN{{exit();}}\n// btfmt: on{ending}"
        );
        assert_eq!(suppression_ranges(&source).unwrap().len(), 2);
        assert_eq!(check(&source, &Config::default()), source);
    }
}

#[test]
fn preprocessor_and_shebang_handling_survives_suppression() {
    let source = concat!(
        "#!/usr/bin/env bpftrace\n",
        "#define MARKER \"// btfmt: off\"\n",
        "#ifdef TEST\n",
        "// btfmt: off\n",
        "#define VALUE 1\n",
        "#endif\n",
        "// btfmt: off\n",
        "#define SCALE(x) \\\n  ((x) * 2)\n",
        "BEGIN{ /* keep */ @x=1;}\n",
        "// btfmt: on\n",
        "END{exit();}\n",
    );
    assert_eq!(suppression_ranges(source).unwrap().len(), 1);
    let output = check(source, &Config::default());
    assert!(output.starts_with("#!/usr/bin/env bpftrace\n\n#define MARKER"));
    assert!(output.contains("#ifdef TEST\n// btfmt: off\n#define VALUE 1\n#endif"));
    assert!(output.contains("END\n{\n    exit();\n}\n"));
}

#[test]
fn suppressed_code_still_requires_valid_syntax() {
    let source = "// btfmt: off\nBEGIN{\n// btfmt: on\n";
    assert!(format_source(source, &Config::default()).is_err());
    assert!(suppression_ranges(source).is_err());
}

#[test]
fn range_offsets_are_utf8_bytes_and_include_complete_directive_lines() {
    let prefix = "// 😀\nBEGIN{}\n";
    let region = " \t// btfmt: off\t\r\nEND{}\r\n  // btfmt: on \r\n";
    let source = format!("{prefix}{region}interval:s:1{{exit();}}\n");
    assert_eq!(
        suppression_ranges(&source).unwrap(),
        vec![prefix.len()..prefix.len() + region.len()]
    );
    check(&source, &Config::default());
}

#[test]
fn cli_uses_shared_suppression_and_check_is_idempotent() {
    let source = "BEGIN{exit();}\n// btfmt: off\nEND{ exit( ); } \t";
    let expected = check(source, &Config::default());
    assert_cmd::Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-")
        .write_stdin(source)
        .assert()
        .success()
        .stdout(expected.clone());
    assert_cmd::Command::cargo_bin("btfmt")
        .unwrap()
        .args(["--check", "-"])
        .write_stdin(expected)
        .assert()
        .success()
        .stdout("");
}

#[test]
fn lsp_preserves_suppressed_bytes_and_matches_cli() {
    let status = std::process::Command::new("python3")
        .arg("-B")
        .arg("scripts/lsp_suppression.py")
        .current_dir(env!("CARGO_MANIFEST_DIR"))
        .env("BTFMT_PATH", assert_cmd::cargo::cargo_bin("btfmt"))
        .status()
        .expect("run scripts/lsp_suppression.py");
    assert!(status.success(), "LSP suppression regression failed");
}
