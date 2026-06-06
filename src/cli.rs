use crate::config::{load_for_cwd, Config};
use crate::format::format_source;
use crate::lsp;
use anyhow::{Context, Result};
use clap::{Parser, Subcommand};
use std::ffi::OsString;
use std::fs;
use std::io::{self, Write};
use std::path::PathBuf;

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
        process_file(path, &config, write_to_file)?;
        if cli.verbose && write_to_file {
            eprintln!("Formatted: {}", path.display());
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

fn process_file(path: &PathBuf, config: &Config, write_to_file: bool) -> Result<()> {
    let source = fs::read_to_string(path).with_context(|| format!("reading {}", path.display()))?;
    let formatted =
        format_source(&source, config).with_context(|| format!("formatting {}", path.display()))?;
    if write_to_file {
        fs::write(path, formatted).with_context(|| format!("writing {}", path.display()))?;
    } else {
        let mut stdout = io::stdout().lock();
        stdout.write_all(formatted.as_bytes())?;
    }
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
