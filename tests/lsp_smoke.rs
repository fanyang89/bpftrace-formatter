use assert_cmd::cargo::cargo_bin;
use std::process::Command;

#[test]
fn lsp_smoke_script_passes() {
    let root = env!("CARGO_MANIFEST_DIR");
    let btfmt = cargo_bin("btfmt");
    let status = Command::new("python3")
        .arg("scripts/lsp_smoke.py")
        .current_dir(root)
        .env("BTFMT_PATH", btfmt)
        .status()
        .expect("run scripts/lsp_smoke.py");
    assert!(status.success(), "lsp smoke failed");
}
