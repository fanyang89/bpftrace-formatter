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
fn cli_preflights_all_inputs_before_output_or_write() {
    let dir = tempdir().unwrap();
    let valid = dir.path().join("valid.bt");
    let invalid = dir.path().join("invalid.bt");
    let valid_source = "BEGIN{exit();}";
    let invalid_source = "END{exit();";
    fs::write(&valid, valid_source).unwrap();
    fs::write(&invalid, invalid_source).unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-w")
        .arg(&valid)
        .arg(&invalid)
        .assert()
        .failure();
    assert_eq!(fs::read_to_string(&valid).unwrap(), valid_source);
    assert_eq!(fs::read_to_string(&invalid).unwrap(), invalid_source);

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg(&valid)
        .arg(&invalid)
        .assert()
        .failure()
        .stdout("");
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
fn cli_generate_config_requires_force_to_overwrite() {
    let dir = tempdir().unwrap();
    let path = dir.path().join("generated.json");
    fs::write(&path, "keep me\n").unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("--generate-config")
        .arg("--config-output")
        .arg(&path)
        .assert()
        .failure()
        .stderr(
            predicate::str::contains("already exists").and(predicate::str::contains("--force")),
        );
    assert_eq!(fs::read_to_string(&path).unwrap(), "keep me\n");

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("-generate-config")
        .arg("-force")
        .arg("-config-output")
        .arg(&path)
        .assert()
        .success();
    assert!(fs::read_to_string(&path).unwrap().contains("indent"));

    Command::cargo_bin("btfmt")
        .unwrap()
        .arg("--force")
        .arg("input.bt")
        .assert()
        .failure()
        .stderr(predicate::str::contains("required arguments"));
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

#[test]
fn cli_discovers_directories_in_sorted_order_and_respects_ignore_rules() {
    let dir = tempdir().unwrap();
    fs::create_dir(dir.path().join("nested")).unwrap();
    fs::create_dir(dir.path().join(".hidden")).unwrap();
    fs::write(dir.path().join(".gitignore"), "ignored.bt\n").unwrap();
    fs::write(dir.path().join(".ignore"), "*.skip.bt\n").unwrap();
    fs::write(dir.path().join("nested/.gitignore"), "*.bt\n!keep.bt\n").unwrap();
    for name in [
        "ignored.bt",
        "bad.skip.bt",
        ".hidden.bt",
        ".hidden/bad.bt",
        "nested/bad.bt",
        "other.txt",
        "uppercase.BT",
    ] {
        fs::write(dir.path().join(name), "invalid bpftrace").unwrap();
    }
    fs::write(dir.path().join("z.bt"), "END{exit();}").unwrap();
    fs::write(dir.path().join("a.bt"), "BEGIN{exit();}").unwrap();
    fs::write(dir.path().join("nested/keep.bt"), "interval:s:1{exit();}").unwrap();

    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .arg(".")
        .assert()
        .success()
        .stdout(
            "BEGIN\n{\n    exit();\n}\ninterval:s:1\n{\n    exit();\n}\nEND\n{\n    exit();\n}\n",
        );
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--check", "."])
        .assert()
        .failure()
        .stdout("")
        .stderr(predicate::str::contains("format check failed"));
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--write", "."])
        .assert()
        .success()
        .stdout("");
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--check", "."])
        .assert()
        .success()
        .stdout("");
    assert_eq!(
        fs::read_to_string(dir.path().join("ignored.bt")).unwrap(),
        "invalid bpftrace"
    );
    assert_eq!(
        fs::read_to_string(dir.path().join("nested/bad.bt")).unwrap(),
        "invalid bpftrace"
    );
}

#[test]
fn cli_preserves_argument_order_and_deduplicates_overlapping_inputs() {
    let dir = tempdir().unwrap();
    fs::create_dir(dir.path().join("nested")).unwrap();
    fs::write(dir.path().join("a.bt"), "BEGIN{exit();}").unwrap();
    fs::write(dir.path().join("nested/z.bt"), "END{exit();}").unwrap();
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["nested/z.bt", ".", "./a.bt", "nested", "-"])
        .write_stdin("interval:s:1{exit();}")
        .assert()
        .success()
        .stdout(
            "END\n{\n    exit();\n}\nBEGIN\n{\n    exit();\n}\ninterval:s:1\n{\n    exit();\n}\n",
        );
}

#[test]
fn cli_explicit_files_bypass_discovery_filters() {
    let dir = tempdir().unwrap();
    fs::write(dir.path().join(".gitignore"), "*.bt\n").unwrap();
    for name in ["ignored.bt", "script.txt", ".hidden.bt"] {
        fs::write(dir.path().join(name), "BEGIN{exit();}").unwrap();
        Command::cargo_bin("btfmt")
            .unwrap()
            .current_dir(dir.path())
            .arg(name)
            .assert()
            .success()
            .stdout("BEGIN\n{\n    exit();\n}\n");
    }
}

#[test]
fn cli_empty_directories_are_successful_no_ops() {
    let dir = tempdir().unwrap();
    for args in [
        vec!["."],
        vec!["--check", "."],
        vec!["--diff", "."],
        vec!["--write", "."],
    ] {
        Command::cargo_bin("btfmt")
            .unwrap()
            .current_dir(dir.path())
            .args(args)
            .assert()
            .success()
            .stdout("");
    }
}

#[test]
fn cli_directory_failures_do_not_write_or_emit_partial_diffs() {
    let dir = tempdir().unwrap();
    fs::write(dir.path().join("a.bt"), "BEGIN{exit();}").unwrap();
    fs::write(dir.path().join("z.bt"), "END{exit();").unwrap();
    for mode in ["--write", "--diff", "--check"] {
        Command::cargo_bin("btfmt")
            .unwrap()
            .current_dir(dir.path())
            .args([mode, "."])
            .assert()
            .failure()
            .stdout("");
        assert_eq!(
            fs::read_to_string(dir.path().join("a.bt")).unwrap(),
            "BEGIN{exit();}"
        );
    }
    fs::remove_file(dir.path().join("z.bt")).unwrap();
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--diff", ".", "missing.bt"])
        .assert()
        .failure()
        .stdout("")
        .stderr(predicate::str::contains("reading missing.bt"));
}

#[test]
fn cli_diff_emits_unified_diff_and_preserves_input() {
    let dir = tempdir().unwrap();
    let source = "BEGIN{exit();}";
    fs::write(dir.path().join("script.bt"), source).unwrap();
    let expected = "--- script.bt\n+++ script.bt\n@@ -1 +1,4 @@\n-BEGIN{exit();}\n\\ No newline at end of file\n+BEGIN\n+{\n+    exit();\n+}\n";
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--diff", "script.bt"])
        .assert()
        .success()
        .stdout(expected);
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--check", "--diff", "."])
        .assert()
        .failure()
        .stdout(
            predicate::str::contains("-BEGIN{exit();}\n")
                .and(predicate::str::contains("+    exit();\n")),
        )
        .stderr(predicate::str::contains("format check failed"));
    assert_eq!(
        fs::read_to_string(dir.path().join("script.bt")).unwrap(),
        source
    );
}

#[test]
fn cli_diff_omits_unchanged_files_and_supports_stdin() {
    let dir = tempdir().unwrap();
    fs::write(
        dir.path().join("unchanged.bt"),
        "BEGIN\n{\n    exit();\n}\n",
    )
    .unwrap();
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--check", "--diff", "."])
        .assert()
        .success()
        .stdout("");
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--diff", ".", "-"])
        .write_stdin("BEGIN{printf(\"😀\");}")
        .assert()
        .success()
        .stdout(
            predicate::str::contains("--- <stdin>\n+++ <stdin>\n")
                .and(predicate::str::contains("+    printf(\"😀\");"))
                .and(predicate::str::contains("unchanged.bt").not()),
        );
    Command::cargo_bin("btfmt")
        .unwrap()
        .args(["--check", "--diff", "-"])
        .write_stdin("BEGIN{exit();}")
        .assert()
        .failure()
        .stdout(predicate::str::contains("--- <stdin>"));
    Command::cargo_bin("btfmt")
        .unwrap()
        .args(["--diff", "-", "-"])
        .assert()
        .failure()
        .stderr(predicate::str::contains("stdin may only be specified once"));
}

#[test]
fn cli_diff_rejects_mutating_options() {
    for flag in ["--write", "--in-place", "--generate-config"] {
        Command::cargo_bin("btfmt")
            .unwrap()
            .args(["--diff", flag, "input.bt"])
            .assert()
            .failure()
            .stderr(predicate::str::contains("cannot be used with"));
    }
}

#[cfg(unix)]
#[test]
fn cli_directory_discovery_skips_symlinks_but_explicit_links_work() {
    use std::os::unix::fs::symlink;
    let dir = tempdir().unwrap();
    let outside = tempdir().unwrap();
    fs::write(outside.path().join("bad.bt"), "invalid").unwrap();
    fs::write(dir.path().join("script.bt"), "BEGIN{exit();}").unwrap();
    symlink(outside.path(), dir.path().join("external")).unwrap();
    symlink(
        outside.path().join("bad.bt"),
        dir.path().join("external.bt"),
    )
    .unwrap();
    symlink("script.bt", dir.path().join("alias.bt")).unwrap();
    symlink("missing", dir.path().join("broken.bt")).unwrap();
    symlink(".", dir.path().join("loop")).unwrap();
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["alias.bt", "."])
        .assert()
        .success()
        .stdout("BEGIN\n{\n    exit();\n}\n");
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--write", "alias.bt", "."])
        .assert()
        .success();
    assert!(fs::symlink_metadata(dir.path().join("alias.bt"))
        .unwrap()
        .file_type()
        .is_symlink());
    assert_eq!(
        fs::read_to_string(dir.path().join("script.bt")).unwrap(),
        "BEGIN\n{\n    exit();\n}\n"
    );
    assert_eq!(
        fs::read_to_string(outside.path().join("bad.bt")).unwrap(),
        "invalid"
    );
}

#[test]
fn cli_directory_discovery_respects_parent_rules_and_git_excludes() {
    let dir = tempdir().unwrap();
    fs::create_dir_all(dir.path().join(".git/info")).unwrap();
    fs::create_dir(dir.path().join("nested")).unwrap();
    fs::write(dir.path().join(".gitignore"), "parent-ignored.bt\n").unwrap();
    fs::write(dir.path().join(".git/info/exclude"), "excluded.bt\n").unwrap();
    for name in ["parent-ignored.bt", "excluded.bt"] {
        fs::write(dir.path().join("nested").join(name), "invalid").unwrap();
    }
    fs::write(dir.path().join("nested/ok.bt"), "BEGIN{exit();}").unwrap();
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .arg("nested")
        .assert()
        .success()
        .stdout("BEGIN\n{\n    exit();\n}\n");
}

#[test]
fn cli_invalid_ignore_rules_fail_before_writing() {
    let dir = tempdir().unwrap();
    fs::write(dir.path().join(".ignore"), "[z-a]\n").unwrap();
    fs::write(dir.path().join("script.bt"), "BEGIN{exit();}").unwrap();
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--write", "."])
        .assert()
        .failure()
        .stdout("");
    assert_eq!(
        fs::read_to_string(dir.path().join("script.bt")).unwrap(),
        "BEGIN{exit();}"
    );
}

#[test]
fn cli_diff_separates_distant_changes_into_context_hunks() {
    let source = format!(
        "BEGIN\n{{\n    @a=1;\n{}    @z=1;\n}}\n",
        "    @m = 1;\n".repeat(20)
    );
    let output = Command::cargo_bin("btfmt")
        .unwrap()
        .args(["--diff", "-"])
        .write_stdin(source)
        .assert()
        .success()
        .get_output()
        .stdout
        .clone();
    let diff = String::from_utf8(output).unwrap();
    assert_eq!(
        diff.lines().filter(|line| line.starts_with("@@ ")).count(),
        2
    );
    assert!(diff.contains("-    @a=1;\n+    @a = 1;"));
    assert!(diff.contains("-    @z=1;\n+    @z = 1;"));
}

#[cfg(unix)]
#[test]
fn cli_diff_quotes_control_characters_in_file_labels() {
    let dir = tempdir().unwrap();
    let name = "line\nbreak.bt";
    fs::write(dir.path().join(name), "BEGIN{exit();}").unwrap();
    Command::cargo_bin("btfmt")
        .unwrap()
        .current_dir(dir.path())
        .args(["--diff", name])
        .assert()
        .success()
        .stdout(predicate::str::starts_with(
            "--- \"line\\nbreak.bt\"\n+++ \"line\\nbreak.bt\"\n",
        ));
}
