use assert_cmd::cargo::cargo_bin;
use std::process::Command;

#[test]
fn workspace_symbol_protocol_and_lifecycle() {
    let status = Command::new("python3")
        .arg("scripts/lsp_workspace_symbols.py")
        .current_dir(env!("CARGO_MANIFEST_DIR"))
        .env("BTFMT_PATH", cargo_bin("btfmt"))
        .env("PYTHONDONTWRITEBYTECODE", "1")
        .status()
        .expect("run workspace symbol protocol tests");
    assert!(status.success(), "workspace symbol protocol tests failed");
}
