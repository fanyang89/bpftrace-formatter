use btfmt::config::{load_from_base, search_upwards, BraceStyle, Config};
use std::fs;
use tempfile::tempdir;

#[test]
fn default_config_matches_documented_values() {
    let config = Config::default();
    assert_eq!(config.indent.size, 4);
    assert!(config.indent.use_spaces);
    assert!(config.spacing.around_operators);
    assert!(config.spacing.around_commas);
    assert!(!config.spacing.around_parentheses);
    assert_eq!(config.line_breaks.empty_lines_between_probes, 1);
    assert_eq!(config.line_breaks.empty_lines_after_shebang, 1);
    assert_eq!(config.blocks.brace_style, BraceStyle::NextLine);
    assert!(config.blocks.indent_statements);
}

#[test]
fn partial_config_load_merges_with_defaults() {
    let dir = tempdir().unwrap();
    let path = dir.path().join("btfmt.json");
    fs::write(&path, r#"{"indent":{"size":2}}"#).unwrap();

    let config = Config::load(&path).unwrap();
    assert_eq!(config.indent.size, 2);
    assert!(config.indent.use_spaces);
    assert_eq!(config.blocks.brace_style, BraceStyle::NextLine);
}

#[test]
fn save_config_round_trips_and_creates_parent_directory() {
    let dir = tempdir().unwrap();
    let path = dir.path().join("nested/out.json");
    let mut config = Config::default();
    config.indent.size = 8;

    config.save(&path).unwrap();
    let loaded = Config::load(&path).unwrap();
    assert_eq!(loaded.indent.size, 8);
}

#[test]
fn save_rejects_invalid_config_before_creating_a_file() {
    let dir = tempdir().unwrap();
    let path = dir.path().join("invalid.json");
    let mut config = Config::default();
    config.indent.size = 0;

    assert!(config.save(&path).is_err());
    assert!(!path.exists());
}

#[test]
fn validates_invalid_numeric_config_values() {
    let mut config = Config::default();
    config.indent.size = 0;
    assert!(config
        .validate()
        .unwrap_err()
        .to_string()
        .contains("indent.size"));

    let mut config = Config::default();
    config.indent.size = 17;
    assert!(config
        .validate()
        .unwrap_err()
        .to_string()
        .contains("indent.size"));

    let mut config = Config::default();
    config.line_breaks.empty_lines_between_probes = 6;
    assert!(config
        .validate()
        .unwrap_err()
        .to_string()
        .contains("empty_lines_between_probes"));

    let mut config = Config::default();
    config.line_breaks.empty_lines_after_shebang = 6;
    assert!(config
        .validate()
        .unwrap_err()
        .to_string()
        .contains("empty_lines_after_shebang"));
}

#[test]
fn invalid_json_enum_and_unknown_fields_return_errors() {
    let dir = tempdir().unwrap();
    let bad_json = dir.path().join("bad.json");
    fs::write(&bad_json, "{").unwrap();
    assert!(Config::load(&bad_json)
        .unwrap_err()
        .to_string()
        .contains("parsing"));

    let bad_style = dir.path().join("bad-style.json");
    fs::write(&bad_style, r#"{"blocks":{"brace_style":"invalid"}}"#).unwrap();
    assert!(Config::load(&bad_style).is_err());

    let unknown = dir.path().join("unknown.json");
    fs::write(&unknown, r#"{"spacing":{"after_keywords":true}}"#).unwrap();
    assert!(Config::load(&unknown)
        .unwrap_err()
        .to_string()
        .contains("parsing"));
}

#[test]
fn search_upwards_uses_nearest_ancestor() {
    let dir = tempdir().unwrap();
    let a = dir.path().join("a");
    let b = a.join("b");
    let c = b.join("c");
    fs::create_dir_all(&c).unwrap();
    fs::write(a.join(".btfmt.json"), "{}").unwrap();
    fs::write(b.join(".btfmt.json"), "{}").unwrap();

    assert_eq!(
        search_upwards(&c, ".btfmt.json"),
        Some(b.join(".btfmt.json"))
    );
}

#[test]
fn load_from_base_handles_relative_explicit_paths_and_rejects_missing_files() {
    let dir = tempdir().unwrap();
    let config_dir = dir.path().join("config");
    fs::create_dir_all(&config_dir).unwrap();
    fs::write(config_dir.join("btfmt.json"), r#"{"indent":{"size":2}}"#).unwrap();

    let loaded =
        load_from_base(dir.path(), Some(std::path::Path::new("config/btfmt.json"))).unwrap();
    assert_eq!(loaded.indent.size, 2);

    let error = load_from_base(dir.path(), Some(std::path::Path::new("missing.json")))
        .unwrap_err()
        .to_string();
    assert!(error.contains("config file does not exist"));
}
