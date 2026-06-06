use assert_cmd::Command;
use predicates::prelude::*;
use std::fs;
use tempfile::tempdir;

#[test]
fn cli_writes_to_stdout_by_default() {
    let dir = tempdir().unwrap();
    let path = dir.path().join("input.bt");
    fs::write(&path, "BEGIN{printf(\"x\",1);}").unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg(&path)
        .assert()
        .success()
        .stdout(predicate::str::contains("printf"));

    assert_eq!(
        fs::read_to_string(&path).unwrap(),
        "BEGIN{printf(\"x\",1);}"
    );
}

#[test]
fn cli_write_and_in_place_update_files() {
    for flag in ["-w", "-i"] {
        let dir = tempdir().unwrap();
        let path = dir.path().join("input.bt");
        fs::write(&path, "BEGIN{printf(\"x\",1);}").unwrap();

        Command::cargo_bin("btfmt")
            .unwrap()
            .arg(flag)
            .arg(&path)
            .assert()
            .success();
        let updated = fs::read_to_string(&path).unwrap();
        assert!(updated.contains("printf"));
        assert!(updated.ends_with('\n'));
    }
}

#[test]
fn cli_formats_multiple_files() {
    let dir = tempdir().unwrap();
    let first = dir.path().join("first.bt");
    let second = dir.path().join("second.bt");
    fs::write(&first, "BEGIN{exit();}").unwrap();
    fs::write(&second, "END{exit();}").unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-w")
        .arg(&first)
        .arg(&second)
        .assert()
        .success();
    assert!(fs::read_to_string(first).unwrap().contains("BEGIN"));
    assert!(fs::read_to_string(second).unwrap().contains("END"));
}

#[test]
fn cli_generates_default_config_and_accepts_legacy_flag_names() {
    let dir = tempdir().unwrap();
    let path = dir.path().join("generated.json");

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-generate-config")
        .arg("-config-output")
        .arg(&path)
        .assert()
        .success()
        .stdout(predicate::str::contains("Generated default configuration"));

    let contents = fs::read_to_string(path).unwrap();
    assert!(contents.contains("indent"));
}

#[test]
fn cli_uses_explicit_config() {
    let dir = tempdir().unwrap();
    let path = dir.path().join("input.bt");
    let config = dir.path().join("config.json");
    fs::write(&path, "BEGIN{exit();}").unwrap();
    fs::write(&config, r#"{"indent":{"size":2}}"#).unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-config")
        .arg(&config)
        .arg(&path)
        .assert()
        .success()
        .stdout(predicate::str::contains("\n  exit();"));
}

#[test]
fn cli_reports_no_input_read_parse_and_config_errors() {
    Command::cargo_bin("btfmt")
        .unwrap()
        .assert()
        .failure()
        .stderr(predicate::str::contains("no input files specified"));

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("does-not-exist.bt")
        .assert()
        .failure()
        .stderr(predicate::str::contains("reading"));

    let dir = tempdir().unwrap();
    let input = dir.path().join("input.bt");
    let config = dir.path().join("bad.json");
    fs::write(&input, "BEGIN{exit();}").unwrap();
    fs::write(&config, "{").unwrap();
    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-c")
        .arg(&config)
        .arg(&input)
        .assert()
        .failure()
        .stderr(predicate::str::contains("parsing"));

    fs::write(&input, "BEGIN{exit();").unwrap();
    Command::cargo_bin("btfmt")
        .unwrap()
        .arg(&input)
        .assert()
        .failure()
        .stderr(predicate::str::contains("parse failed"));
}

#[test]
fn cli_help_version_and_lsp_help_are_available() {
    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-help")
        .assert()
        .success()
        .stdout(predicate::str::contains("Usage"));

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-version")
        .assert()
        .success()
        .stdout(predicate::str::contains("btfmt"));

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("lsp")
        .arg("--help")
        .assert()
        .success()
        .stdout(predicate::str::contains("Usage"));
}
