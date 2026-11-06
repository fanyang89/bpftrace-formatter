use assert_cmd::cargo::cargo_bin;
use std::process::Command;

#[test]
fn semantic_tokens_protocol() {
    let status = Command::new("python3")
        .arg("scripts/lsp_semantic_tokens.py")
        .current_dir(env!("CARGO_MANIFEST_DIR"))
        .env("BTFMT_PATH", cargo_bin("btfmt"))
        .env("PYTHONDONTWRITEBYTECODE", "1")
        .status()
        .expect("run semantic tokens protocol tests");
    assert!(status.success(), "semantic tokens protocol tests failed");
}
