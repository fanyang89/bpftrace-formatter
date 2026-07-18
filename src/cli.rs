use crate::config::{load_for_cwd, Config};
use crate::format::format_source;
use crate::lsp;
use anyhow::{Context, Result};
use clap::{Parser, Subcommand};
use std::ffi::OsString;
use std::fs;
use std::io::{self, Write};
use std::path::{Path, PathBuf};
use tempfile::NamedTempFile;

#[derive(Debug, Parser)]
#[command(name = "btfmt", version, about = "Format bpftrace scripts")]
struct Cli {
    #[command(subcommand)]
    command: Option<Command>,

    #[arg(short = 'c', long = "config")]
    config: Option<PathBuf>,

    #[arg(short = 'w', long = "write")]
    write: bool,

    #[arg(short = 'i', long = "in-place")]
    in_place: bool,

    #[arg(short = 'v', long = "verbose")]
    verbose: bool,

    #[arg(long = "generate-config")]
    generate_config: bool,

    #[arg(long = "config-output", default_value = ".btfmt.json")]
    config_output: PathBuf,

    #[arg(value_name = "FILE")]
    files: Vec<PathBuf>,
}

#[derive(Debug, Subcommand)]
enum Command {
    Lsp,
}

pub async fn run() -> Result<()> {
    let cli = Cli::parse_from(normalized_args());
    if matches!(cli.command, Some(Command::Lsp)) {
        return lsp::run_server().await;
    }

    if cli.generate_config {
        Config::default().save(&cli.config_output)?;
        println!(
            "Generated default configuration at: {}",
            cli.config_output.display()
        );
        return Ok(());
    }

    if cli.files.is_empty() {
        anyhow::bail!("no input files specified");
    }

    let config = load_for_cwd(cli.config.as_deref())?;
    let write_to_file = cli.write || cli.in_place;
    for path in &cli.files {
        if cli.verbose {
            eprintln!("Processing: {}", path.display());
        }
        let changed = process_file(path, &config, write_to_file)?;
        if cli.verbose && write_to_file {
            let status = if changed { "Formatted" } else { "Unchanged" };
            eprintln!("{status}: {}", path.display());
        }
    }
    Ok(())
}

fn normalized_args() -> Vec<OsString> {
    std::env::args_os()
        .map(|arg| match arg.to_string_lossy().as_ref() {
            "-config" => OsString::from("--config"),
            "-generate-config" => OsString::from("--generate-config"),
            "-config-output" => OsString::from("--config-output"),
            "-verbose" => OsString::from("--verbose"),
            "-version" => OsString::from("--version"),
            "-help" => OsString::from("--help"),
            _ => arg,
        })
        .collect()
}

fn process_file(path: &Path, config: &Config, write_to_file: bool) -> Result<bool> {
    let source = fs::read_to_string(path).with_context(|| format!("reading {}", path.display()))?;
    let formatted =
        format_source(&source, config).with_context(|| format!("formatting {}", path.display()))?;
    if write_to_file {
        if formatted == source {
            return Ok(false);
        }
        write_atomic(path, formatted.as_bytes())?;
    } else {
        let mut stdout = io::stdout().lock();
        stdout.write_all(formatted.as_bytes())?;
    }
    Ok(true)
}

fn write_atomic(path: &Path, contents: &[u8]) -> Result<()> {
    let target = if fs::symlink_metadata(path)?.file_type().is_symlink() {
        fs::canonicalize(path).with_context(|| format!("resolving {}", path.display()))?
    } else {
        path.to_path_buf()
    };
    let metadata = fs::metadata(&target)?;
    let parent = target
        .parent()
        .filter(|parent| !parent.as_os_str().is_empty())
        .unwrap_or_else(|| Path::new("."));
    let mut temp = NamedTempFile::new_in(parent)
        .with_context(|| format!("creating temporary file in {}", parent.display()))?;
    temp.as_file()
        .set_permissions(metadata.permissions())
        .with_context(|| format!("preserving permissions for {}", target.display()))?;
    temp.write_all(contents)
        .with_context(|| format!("writing temporary file for {}", target.display()))?;
    temp.as_file()
        .sync_all()
        .with_context(|| format!("syncing temporary file for {}", target.display()))?;
    temp.persist(&target)
        .map_err(|err| err.error)
        .with_context(|| format!("replacing {}", target.display()))?;
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn normalizes_legacy_long_flags() {
        let args = vec![OsString::from("btfmt"), OsString::from("-generate-config")];
        let parsed = Cli::parse_from(args.into_iter().map(
            |arg| match arg.to_string_lossy().as_ref() {
                "-generate-config" => OsString::from("--generate-config"),
                value => OsString::from(value),
            },
        ));
        assert!(parsed.generate_config);
    }
}
