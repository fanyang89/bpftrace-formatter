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
fn cli_formats_explicit_stdin() {
    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-")
        .write_stdin("BEGIN{exit();}")
        .assert()
        .success()
        .stdout("BEGIN\n{\n    exit();\n}\n");

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-")
        .write_stdin("BEGIN{exit();")
        .assert()
        .failure()
        .stderr(predicate::str::contains("formatting <stdin>"));

    Command::cargo_bin("btfmt")
        .unwrap()
        .args(["-", "-"])
        .write_stdin("BEGIN{exit();}")
        .assert()
        .failure()
        .stderr(predicate::str::contains("stdin may only be specified once"));
}

#[test]
fn cli_check_reports_changed_files_without_writing() {
    let dir = tempdir().unwrap();
    let unchanged = dir.path().join("unchanged.bt");
    let changed = dir.path().join("changed.bt");
    let formatted = "BEGIN\n{\n    exit();\n}\n";
    fs::write(&unchanged, formatted).unwrap();
    fs::write(&changed, "END{exit();}").unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("--check")
        .arg(&unchanged)
        .assert()
        .success()
        .stdout("");

    Command::cargo_bin("btfmt")
        .unwrap()
        .args(["--check", "--verbose"])
        .arg(&unchanged)
        .arg(&changed)
        .assert()
        .failure()
        .stdout("")
        .stderr(
            predicate::str::contains("Unchanged:")
                .and(predicate::str::contains("Would reformat:"))
                .and(predicate::str::contains("format check failed")),
        );
    assert_eq!(fs::read_to_string(&changed).unwrap(), "END{exit();}");

    Command::cargo_bin("btfmt")
        .unwrap()
        .args(["--check", "-"])
        .write_stdin(formatted)
        .assert()
        .success()
        .stdout("");
    Command::cargo_bin("btfmt")
        .unwrap()
        .args(["--check", "-"])
        .write_stdin("BEGIN{exit();}")
        .assert()
        .failure()
        .stdout("")
        .stderr(predicate::str::contains("<stdin>"));
}

#[test]
fn cli_rejects_writing_stdin_and_check_write_conflicts() {
    Command::cargo_bin("btfmt")
        .unwrap()
        .args(["-w", "-"])
        .write_stdin("BEGIN{exit();}")
        .assert()
        .failure()
        .stderr(predicate::str::contains("cannot write stdin in place"));

    Command::cargo_bin("btfmt")
        .unwrap()
        .args(["--check", "-w", "input.bt"])
        .assert()
        .failure()
        .stderr(predicate::str::contains("cannot be used with"));
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
fn cli_write_preserves_invalid_and_unchanged_files() {
    let dir = tempdir().unwrap();
    let invalid = dir.path().join("invalid.bt");
    let invalid_source = "BEGIN{exit();";
    fs::write(&invalid, invalid_source).unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-w")
        .arg(&invalid)
        .assert()
        .failure();
    assert_eq!(fs::read_to_string(&invalid).unwrap(), invalid_source);

    let unchanged = dir.path().join("unchanged.bt");
    let formatted_source = "BEGIN\n{\n    exit();\n}\n";
    fs::write(&unchanged, formatted_source).unwrap();
    let modified = fs::metadata(&unchanged).unwrap().modified().unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .args(["-w", "-v"])
        .arg(&unchanged)
        .assert()
        .success()
        .stderr(predicate::str::contains("Unchanged:"));
    assert_eq!(fs::read_to_string(&unchanged).unwrap(), formatted_source);
    assert_eq!(
        fs::metadata(&unchanged).unwrap().modified().unwrap(),
        modified
    );
}

#[cfg(unix)]
#[test]
fn cli_write_preserves_file_permissions() {
    use std::os::unix::fs::PermissionsExt;

    let dir = tempdir().unwrap();
    let path = dir.path().join("input.bt");
    fs::write(&path, "BEGIN{exit();}").unwrap();
    fs::set_permissions(&path, fs::Permissions::from_mode(0o640)).unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-w")
        .arg(&path)
        .assert()
        .success();
    assert_eq!(
        fs::metadata(path).unwrap().permissions().mode() & 0o777,
        0o640
    );
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
