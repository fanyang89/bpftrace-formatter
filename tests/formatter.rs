use btfmt::config::{BraceStyle, Config};
use btfmt::format::format_source;
use btfmt::parse::parse;
use std::fs;
use std::path::{Path, PathBuf};

fn fmt(input: &str, config: &Config) -> String {
    format_source(input, config).unwrap()
}

#[test]
fn parses_valid_input_and_reports_invalid_input() {
    let valid = parse("BEGIN { printf(\"x\"); }").unwrap();
    assert!(valid.diagnostics.is_empty());

    let invalid = parse("BEGIN { printf(\"x\");").unwrap();
    assert!(!invalid.diagnostics.is_empty());
    assert!(format_source("BEGIN { printf(\"x\");", &Config::default()).is_err());

    let invalid_unicode = parse("BEGIN { printf(\"😀\");").unwrap();
    assert!(!invalid_unicode.diagnostics.is_empty());
}

#[test]
fn formats_golden_fixtures_exactly() {
    let root = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    for path in bt_files(&root.join("tests/testdata")) {
        if path
            .components()
            .any(|component| component.as_os_str() == "golden")
        {
            continue;
        }
        let source = fs::read_to_string(&path).unwrap();
        let formatted = fmt(&source, &Config::default());
        let golden_path = root
            .join("tests/testdata/golden")
            .join(path.file_name().expect("fixture filename"));
        let expected = fs::read_to_string(&golden_path)
            .unwrap_or_else(|err| panic!("read {}: {err}", golden_path.display()));
        assert_eq!(
            formatted,
            expected,
            "golden mismatch for {}",
            path.display()
        );
    }
}

#[test]
fn brace_styles_are_honored() {
    let mut config = Config::default();
    config.blocks.brace_style = BraceStyle::SameLine;
    assert_eq!(fmt("BEGIN{exit();}", &config), "BEGIN {\n    exit();\n}\n");

    config.blocks.brace_style = BraceStyle::NextLine;
    assert_eq!(fmt("BEGIN{exit();}", &config), "BEGIN\n{\n    exit();\n}\n");

    config.blocks.brace_style = BraceStyle::Gnu;
    assert_eq!(
        fmt("BEGIN{exit();}", &config),
        "BEGIN\n    {\n        exit();\n}\n"
    );
}

#[test]
fn indentation_settings_are_honored() {
    let mut config = Config::default();
    config.indent.size = 2;
    assert_eq!(fmt("BEGIN{exit();}", &config), "BEGIN\n{\n  exit();\n}\n");

    config.indent.use_spaces = false;
    assert_eq!(fmt("BEGIN{exit();}", &config), "BEGIN\n{\n\texit();\n}\n");

    config.blocks.indent_statements = false;
    assert_eq!(fmt("BEGIN{exit();}", &config), "BEGIN\n{\nexit();\n}\n");
}

#[test]
fn spacing_settings_are_honored() {
    let mut config = Config::default();
    config.spacing.around_commas = false;
    config.spacing.around_parentheses = true;
    assert_eq!(
        fmt("BEGIN{printf(\"a\",1,2);}", &config),
        "BEGIN\n{\n    printf( \"a\",1,2 );\n}\n"
    );

    let mut config = Config::default();
    config.spacing.around_brackets = true;
    config.spacing.around_operators = false;
    assert_eq!(
        fmt("BEGIN{@[pid]=count();}", &config),
        "BEGIN\n{\n    @[ pid ]=count();\n}\n"
    );
}

#[test]
fn comments_shebang_preprocessor_and_probe_spacing_are_preserved() {
    let mut config = Config::default();
    config.line_breaks.empty_lines_after_shebang = 2;
    config.line_breaks.empty_lines_between_probes = 2;

    let output = fmt(
        "#!/usr/bin/env bpftrace\n#define X 1\nBEGIN{printf(\"x\"); // hello\n}END{exit();}",
        &config,
    );
    assert_eq!(
        output,
        concat!(
            "#!/usr/bin/env bpftrace\n\n\n",
            "#define X 1\n",
            "BEGIN\n",
            "{\n",
            "    printf(\"x\"); // hello\n",
            "}\n\n\n",
            "END\n",
            "{\n",
            "    exit();\n",
            "}\n",
        )
    );
}

#[test]
fn control_flow_empty_blocks_and_comments_are_exact() {
    assert_eq!(fmt("BEGIN{}", &Config::default()), "BEGIN\n{\n}\n");
    assert_eq!(
        fmt(
            "BEGIN{if(1){exit();}else{exit();} exit();}",
            &Config::default()
        ),
        concat!(
            "BEGIN\n",
            "{\n",
            "    if (1)\n",
            "    {\n",
            "        exit();\n",
            "    }\n",
            "    else\n",
            "    {\n",
            "        exit();\n",
            "    }\n",
            "    exit();\n",
            "}\n",
        )
    );
    assert_eq!(
        fmt(
            "BEGIN{exit();}\n// next probe\nEND{exit();}",
            &Config::default()
        ),
        concat!(
            "BEGIN\n",
            "{\n",
            "    exit();\n",
            "}\n\n",
            "// next probe\n",
            "END\n",
            "{\n",
            "    exit();\n",
            "}\n",
        )
    );
}

#[test]
fn supported_layout_settings_are_exact() {
    let mut config = Config::default();
    config.indent.size = 8;
    assert_eq!(
        fmt("BEGIN{exit();}", &config),
        "BEGIN\n{\n        exit();\n}\n"
    );

    config.blocks.brace_style = BraceStyle::SameLine;
    config.spacing.before_block_start = false;
    assert_eq!(
        fmt("BEGIN{exit();}", &config),
        "BEGIN{\n        exit();\n}\n"
    );
}

#[test]
fn empty_line_counts_are_exact() {
    for count in 0..=2 {
        let mut config = Config::default();
        config.line_breaks.empty_lines_between_probes = count;
        let expected = format!(
            "BEGIN\n{{\n    exit();\n}}{}END\n{{\n    exit();\n}}\n",
            "\n".repeat(count + 1)
        );
        assert_eq!(fmt("BEGIN{exit();}END{exit();}", &config), expected);

        config.line_breaks.empty_lines_after_shebang = count;
        let expected = format!(
            "#!/usr/bin/env bpftrace{}BEGIN\n{{\n    exit();\n}}\n",
            "\n".repeat(count + 1)
        );
        assert_eq!(
            fmt("#!/usr/bin/env bpftrace\nBEGIN{exit();}", &config),
            expected
        );
    }
}

#[test]
fn important_bpftrace_tokens_are_not_split_by_spacing() {
    let output = fmt(
        "tracepoint:syscalls:sys_enter_*/pid==1234/{printf(\"%s\", str(args->filename));}",
        &Config::default(),
    );
    assert!(output.contains("sys_enter_*"));
    assert!(output.contains("sys_enter_*\n/pid == 1234/\n"));
    assert!(output.contains("args->filename"));

    let output = fmt(
        "kprobe:vfs_read*,kprobe:vfs_write*{exit();}",
        &Config::default(),
    );
    assert!(output.contains("vfs_read*"));
    assert!(output.contains("vfs_write*"));
}

#[test]
fn reproduction_inputs_keep_constructs_parseable() {
    let cases = [
        "BEGIN{if(1){exit();}else{exit();}}",
        "BEGIN{// hello\nprintf(\"x\");}",
        "tracepoint:syscalls:sys_enter_openat/pid==1234/{exit();}",
        "tracepoint:syscalls:sys_enter_openat,tracepoint:syscalls:sys_enter_open{exit();}",
    ];

    for input in cases {
        let output = fmt(input, &Config::default());
        assert!(
            parse(&output).unwrap().diagnostics.is_empty(),
            "{input}: {output}"
        );
    }
}

#[test]
fn protected_comments_and_preprocessor_regions_are_preserved() {
    let source = concat!(
        "#define SCALE(x) \\\n",
        "  ((x) * 2)\n",
        "#ifdef __x86_64__\n",
        "#define REG ax\n",
        "#else\n",
        "#define REG r0\n",
        "#endif\n",
        "BEGIN { /* keep this block comment */ printf(\"%d\", SCALE(1)); }\n",
    );
    let formatted = fmt(source, &Config::default());

    assert!(formatted.contains("#define SCALE(x) \\\n  ((x) * 2)"));
    assert!(formatted.contains("#ifdef __x86_64__\n#define REG ax\n#else\n#define REG r0\n#endif"));
    assert!(formatted.contains("/* keep this block comment */"));
    assert!(parse(&formatted).unwrap().diagnostics.is_empty());
    assert_eq!(fmt(&formatted, &Config::default()), formatted);
}

fn bt_files(root: &Path) -> Vec<PathBuf> {
    let mut files = Vec::new();
    collect_bt_files(root, &mut files);
    files.sort();
    files
}

fn collect_bt_files(path: &Path, files: &mut Vec<PathBuf>) {
    let Ok(entries) = fs::read_dir(path) else {
        return;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            collect_bt_files(&path, files);
        } else if path.extension().is_some_and(|ext| ext == "bt") {
            files.push(path);
        }
    }
}
