# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`btfmt` is a Rust-based bpftrace formatter with CLI and LSP support. Parsing is provided by the `tree-sitter-bpftrace` crate, while formatting, configuration, and language-server behavior live in this repository.

## Development Commands

### Building

```bash
task build
cargo build --bin btfmt
cargo build --release --bin btfmt
```

### Testing

```bash
task test
cargo test --all-targets --all-features
cargo test --test formatter
cargo test lsp_smoke_script_passes
```

### Formatting And Linting

```bash
task fmt
cargo fmt --all -- --check
cargo clippy --all-targets --all-features -- -D warnings
```

### Running

```bash
./btfmt <file.bt>
./btfmt -i <file.bt>
./btfmt -w <file.bt>
./btfmt lsp
```

## Architecture

- `src/main.rs`: process entrypoint
- `src/cli.rs`: CLI argument handling and file processing
- `src/config.rs`: JSON configuration types, defaults, validation, and loading
- `src/parse.rs`: tree-sitter parser setup and diagnostics
- `src/format.rs`: formatter implementation
- `src/lsp.rs`: tower-lsp server implementation
- `src/text.rs`: byte offset and LSP range helpers
- `tests/`: integration tests for CLI, config, formatter, LSP, and text helpers
- `tests/testdata/`: formatter fixtures
- `tests/testdata/golden/`: expected Rust formatter output

## Parser Policy

- Use `tree-sitter-bpftrace` from crates.io.
- Do not commit generated parser sources.
- If grammar coverage needs changes, prefer updating or patching the upstream tree-sitter grammar.

## Configuration Hierarchy

CLI configuration is loaded in this order:

1. File specified with `-config` / `-c`
2. `.btfmt.json` in the current directory or parent directories
3. `~/.btfmt.json`
4. Built-in defaults

LSP formatting uses an explicit `btfmt.configPath` when provided; otherwise it searches from the document directory.

## Testing Strategy

- Unit tests cover small helpers and formatter basics.
- Integration tests cover config, CLI, formatter golden output, text ranges, and LSP smoke behavior.
- `tests/lsp_smoke.rs` runs `scripts/lsp_smoke.py` against a Cargo-built binary.
- `tests/formatter.rs` verifies exact golden output for `tests/testdata/*.bt` and smoke-formats `bpftrace/tools`.

## Release Notes

Native GitHub Actions runners build binaries and platform-specific VSIX packages. GoReleaser generates checksums and changelog content, then publishes the assembled artifacts for tags matching the versions in `Cargo.toml` and `vscode-extension/package.json`.

Release builds currently publish:

- Linux x64
- macOS x64
- macOS ARM64
- Windows x64

Linux ARM64 and Windows ARM64 artifacts are intentionally deferred.
