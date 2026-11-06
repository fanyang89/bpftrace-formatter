use assert_cmd::cargo::cargo_bin;
use std::process::Command;

#[test]
fn range_formatting_protocol() {
    let status = Command::new("python3")
        .arg("-B")
        .arg("scripts/lsp_range_formatting.py")
        .current_dir(env!("CARGO_MANIFEST_DIR"))
        .env("BTFMT_PATH", cargo_bin("btfmt"))
        .status()
        .expect("run range formatting protocol test");
    assert!(status.success(), "range formatting protocol test failed");
}
