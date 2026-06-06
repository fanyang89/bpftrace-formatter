use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};
use std::env;
use std::fs;
use std::path::{Path, PathBuf};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(default)]
pub struct Config {
    pub indent: IndentConfig,
    pub spacing: SpacingConfig,
    pub line_breaks: LineBreakConfig,
    pub comments: CommentConfig,
    pub probes: ProbeConfig,
    pub blocks: BlockConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(default)]
pub struct IndentConfig {
    pub size: usize,
    pub use_spaces: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(default)]
pub struct SpacingConfig {
    pub around_operators: bool,
    pub around_commas: bool,
    pub around_parentheses: bool,
    pub around_brackets: bool,
    pub before_block_start: bool,
    pub after_keywords: bool,
    pub inside_predicates: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(default)]
pub struct LineBreakConfig {
    pub max_line_length: usize,
    pub break_long_statements: bool,
    pub empty_lines_between_probes: usize,
    pub empty_lines_after_shebang: usize,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(default)]
pub struct CommentConfig {
    pub preserve_inline: bool,
    pub indent_level: usize,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(default)]
pub struct ProbeConfig {
    pub align_predicates: bool,
    pub newline_between_specifiers: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(default)]
pub struct BlockConfig {
    pub brace_style: BraceStyle,
    pub indent_statements: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Default)]
#[serde(rename_all = "snake_case")]
pub enum BraceStyle {
    SameLine,
    #[default]
    NextLine,
    Gnu,
}

impl Default for IndentConfig {
    fn default() -> Self {
        Config::default().indent
    }
}

impl Default for SpacingConfig {
    fn default() -> Self {
        Config::default().spacing
    }
}

impl Default for LineBreakConfig {
    fn default() -> Self {
        Config::default().line_breaks
    }
}

impl Default for CommentConfig {
    fn default() -> Self {
        Config::default().comments
    }
}

impl Default for ProbeConfig {
    fn default() -> Self {
        Config::default().probes
    }
}

impl Default for BlockConfig {
    fn default() -> Self {
        Config::default().blocks
    }
}

impl Default for Config {
    fn default() -> Self {
        Self {
            indent: IndentConfig {
                size: 4,
                use_spaces: true,
            },
            spacing: SpacingConfig {
                around_operators: true,
                around_commas: true,
                around_parentheses: false,
                around_brackets: false,
                before_block_start: true,
                after_keywords: true,
                inside_predicates: true,
            },
            line_breaks: LineBreakConfig {
                max_line_length: 80,
                break_long_statements: true,
                empty_lines_between_probes: 1,
                empty_lines_after_shebang: 1,
            },
            comments: CommentConfig {
                preserve_inline: true,
                indent_level: 0,
            },
            probes: ProbeConfig {
                align_predicates: false,
                newline_between_specifiers: false,
            },
            blocks: BlockConfig {
                brace_style: BraceStyle::NextLine,
                indent_statements: true,
            },
        }
    }
}

impl Config {
    pub fn validate(&self) -> Result<()> {
        if !(1..=16).contains(&self.indent.size) {
            anyhow::bail!(
                "indent.size must be between 1 and 16, got {}",
                self.indent.size
            );
        }
        if self.line_breaks.empty_lines_between_probes > 5 {
            anyhow::bail!(
                "line_breaks.empty_lines_between_probes must be 0-5, got {}",
                self.line_breaks.empty_lines_between_probes
            );
        }
        if self.line_breaks.empty_lines_after_shebang > 5 {
            anyhow::bail!(
                "line_breaks.empty_lines_after_shebang must be 0-5, got {}",
                self.line_breaks.empty_lines_after_shebang
            );
        }
        Ok(())
    }

    pub fn load(path: &Path) -> Result<Self> {
        let data =
            fs::read_to_string(path).with_context(|| format!("reading {}", path.display()))?;
        let config: Self =
            serde_json::from_str(&data).with_context(|| format!("parsing {}", path.display()))?;
        config.validate()?;
        Ok(config)
    }

    pub fn save(&self, path: &Path) -> Result<()> {
        if let Some(parent) = path.parent().filter(|p| !p.as_os_str().is_empty()) {
            fs::create_dir_all(parent).with_context(|| format!("creating {}", parent.display()))?;
        }
        let data = serde_json::to_string_pretty(self)?;
        fs::write(path, data + "\n").with_context(|| format!("writing {}", path.display()))
    }
}

pub fn load_for_cwd(explicit_path: Option<&Path>) -> Result<Config> {
    if let Some(path) = explicit_path {
        if path.exists() {
            return Config::load(path);
        }
        return Ok(Config::default());
    }

    let cwd = env::current_dir().unwrap_or_else(|_| PathBuf::from("."));
    if let Some(path) = search_upwards(&cwd, ".btfmt.json") {
        return Config::load(&path);
    }

    if let Some(home) = home_dir() {
        let path = home.join(".btfmt.json");
        if path.exists() {
            return Config::load(&path);
        }
    }

    Ok(Config::default())
}

pub fn load_from_base(base_dir: &Path, explicit_path: Option<&Path>) -> Result<Config> {
    if let Some(path) = explicit_path {
        let path = if path.is_absolute() {
            path.to_path_buf()
        } else {
            base_dir.join(path)
        };
        if path.exists() {
            return Config::load(&path);
        }
        return Ok(Config::default());
    }

    if let Some(path) = search_upwards(base_dir, ".btfmt.json") {
        return Config::load(&path);
    }
    Ok(Config::default())
}

pub fn search_upwards(start: &Path, filename: &str) -> Option<PathBuf> {
    let mut current = start.to_path_buf();
    loop {
        let candidate = current.join(filename);
        if candidate.exists() {
            return Some(candidate);
        }
        if !current.pop() {
            return None;
        }
    }
}

fn home_dir() -> Option<PathBuf> {
    env::var_os("HOME").map(PathBuf::from)
}
