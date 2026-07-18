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
}

#[test]
fn formats_golden_fixtures_exactly() {
    let root = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    for path in bt_files(&root.join("testdata")) {
        if path
            .components()
            .any(|component| component.as_os_str() == "golden")
        {
            continue;
        }
        let source = fs::read_to_string(&path).unwrap();
        let formatted = fmt(&source, &Config::default());
        let golden_path = root
            .join("testdata/golden")
            .join(path.file_name().expect("fixture filename"));
        if golden_path.exists() {
            let expected = fs::read_to_string(&golden_path).unwrap();
            assert_eq!(
                formatted,
                expected,
                "golden mismatch for {}",
                path.display()
            );
        }
    }
}

#[test]
fn formats_bpftrace_tools_tree() {
    let root = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    let tools = root.join("bpftrace/tools");
    if !tools.exists() {
        return;
    }
    for path in bt_files(&tools) {
        let source = fs::read_to_string(&path).unwrap();
        let formatted = fmt(&source, &Config::default());
        assert!(
            parse(&formatted).unwrap().diagnostics.is_empty(),
            "formatter produced invalid output for {}",
            path.display()
        );
        assert_eq!(
            fmt(&formatted, &Config::default()),
            formatted,
            "formatter is not idempotent for {}",
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
    assert!(output.starts_with("#!/usr/bin/env bpftrace\n\n#define X 1\n"));
    assert!(output.contains("printf(\"x\"); // hello\n"));
    assert!(output.contains("}\n\nEND"));
}

#[test]
fn important_bpftrace_tokens_are_not_split_by_spacing() {
    let output = fmt(
        "tracepoint:syscalls:sys_enter_*/pid==1234/{printf(\"%s\", str(args->filename));}",
        &Config::default(),
    );
    assert!(output.contains("sys_enter_*"));
    assert!(output.contains("/pid == 1234/"));
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
