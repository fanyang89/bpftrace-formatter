use crate::config::{load_for_cwd, Config};
use crate::file_io;
use crate::format::format_source;
use crate::lsp;
use anyhow::{Context, Result};
use clap::{Parser, Subcommand};
use std::ffi::OsString;
use std::fs;
use std::io::{self, Read, Write};
use std::path::{Path, PathBuf};

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

    #[arg(
        long = "check",
        conflicts_with_all = ["write", "in_place", "generate_config"]
    )]
    check: bool,

    #[arg(long = "generate-config")]
    generate_config: bool,

    #[arg(long = "force", requires = "generate_config")]
    force: bool,

    #[arg(long = "config-output", default_value = ".btfmt.json")]
    config_output: PathBuf,

    #[arg(value_name = "FILE", allow_hyphen_values = true)]
    files: Vec<PathBuf>,
}

#[derive(Debug, Subcommand)]
enum Command {
    Lsp,
}

#[derive(Debug)]
struct FormattedInput {
    path: PathBuf,
    label: String,
    source: String,
    formatted: String,
}

impl FormattedInput {
    fn changed(&self) -> bool {
        self.formatted != self.source
    }
}

pub async fn run() -> Result<()> {
    let cli = Cli::parse_from(normalized_args());
    if matches!(cli.command, Some(Command::Lsp)) {
        return lsp::run_server().await;
    }

    if cli.generate_config {
        let config = Config::default();
        if cli.force {
            config.save(&cli.config_output)?;
        } else {
            if cli.config_output.exists() {
                anyhow::bail!(
                    "configuration already exists at {}; use --force to overwrite it",
                    cli.config_output.display()
                );
            }
            config.save_new(&cli.config_output)?;
        }
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
    let stdin_count = cli.files.iter().filter(|path| is_stdin(path)).count();
    if stdin_count > 1 {
        anyhow::bail!("stdin may only be specified once");
    }
    if write_to_file && stdin_count != 0 {
        anyhow::bail!("cannot write stdin in place");
    }

    let mut inputs = Vec::with_capacity(cli.files.len());
    for path in &cli.files {
        let label = input_label(path);
        if cli.verbose {
            eprintln!("Processing: {label}");
        }
        let source = read_input(path)?;
        let formatted =
            format_source(&source, &config).with_context(|| format!("formatting {label}"))?;
        inputs.push(FormattedInput {
            path: path.clone(),
            label,
            source,
            formatted,
        });
    }

    let mut check_failures = Vec::new();
    for input in &inputs {
        let changed = input.changed();

        if cli.check {
            if changed {
                check_failures.push(input.label.clone());
            }
            if cli.verbose {
                let status = if changed {
                    "Would reformat"
                } else {
                    "Unchanged"
                };
                eprintln!("{status}: {}", input.label);
            }
        } else if write_to_file {
            if changed {
                file_io::write_atomic(&input.path, input.formatted.as_bytes())?;
            }
        } else {
            let mut stdout = io::stdout().lock();
            stdout.write_all(input.formatted.as_bytes())?;
        }

        if cli.verbose && write_to_file {
            let status = if changed { "Formatted" } else { "Unchanged" };
            eprintln!("{status}: {}", input.label);
        }
    }

    if !check_failures.is_empty() {
        anyhow::bail!("format check failed for: {}", check_failures.join(", "));
    }
    Ok(())
}

fn normalized_args() -> Vec<OsString> {
    std::env::args_os()
        .map(|arg| match arg.to_string_lossy().as_ref() {
            "-config" => OsString::from("--config"),
            "-generate-config" => OsString::from("--generate-config"),
            "-config-output" => OsString::from("--config-output"),
            "-force" => OsString::from("--force"),
            "-verbose" => OsString::from("--verbose"),
            "-version" => OsString::from("--version"),
            "-help" => OsString::from("--help"),
            _ => arg,
        })
        .collect()
}

fn is_stdin(path: &Path) -> bool {
    path == Path::new("-")
}

fn input_label(path: &Path) -> String {
    if is_stdin(path) {
        "<stdin>".to_string()
    } else {
        path.display().to_string()
    }
}

fn read_input(path: &Path) -> Result<String> {
    if is_stdin(path) {
        let mut source = String::new();
        io::stdin()
            .lock()
            .read_to_string(&mut source)
            .context("reading stdin")?;
        Ok(source)
    } else {
        fs::read_to_string(path).with_context(|| format!("reading {}", path.display()))
    }
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
