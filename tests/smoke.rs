use btfmt::config::Config;
use btfmt::format::format_source;
use btfmt::parse::parse;
use std::fs;
use std::path::{Path, PathBuf};

#[test]
fn formats_testdata_fixtures() {
    let root = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    let config = Config::default();
    for path in bt_files(&root.join("tests/testdata")) {
        let source = fs::read_to_string(&path)
            .unwrap_or_else(|err| panic!("read {}: {err}", path.display()));
        let formatted = format_source(&source, &config)
            .unwrap_or_else(|err| panic!("format {}: {err:#}", path.display()));
        assert!(
            !formatted.trim().is_empty(),
            "empty output for {}",
            path.display()
        );
        assert!(
            parse(&formatted).unwrap().diagnostics.is_empty(),
            "invalid output for {}",
            path.display()
        );
        assert_eq!(
            format_source(&formatted, &config).unwrap(),
            formatted,
            "non-idempotent output for {}",
            path.display()
        );
    }
}

#[test]
fn formats_bpftrace_tools_corpus() {
    let root = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    let tools = root.join("tests/corpus/bpftrace-tools");

    let config = Config::default();
    for path in bt_files(&tools) {
        let source = fs::read_to_string(&path)
            .unwrap_or_else(|err| panic!("read {}: {err}", path.display()));
        let formatted = format_source(&source, &config)
            .unwrap_or_else(|err| panic!("format {}: {err:#}", path.display()));
        assert!(
            !formatted.trim().is_empty(),
            "empty output for {}",
            path.display()
        );
        assert!(
            parse(&formatted).unwrap().diagnostics.is_empty(),
            "invalid output for {}",
            path.display()
        );
        assert_eq!(
            format_source(&formatted, &config).unwrap(),
            formatted,
            "non-idempotent output for {}",
            path.display()
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
