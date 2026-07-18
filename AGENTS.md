# AGENTS

This file is the operational guide for agentic coding tools in this repository.
It summarizes the commands and style conventions observed in the codebase.

## Repository Overview

- Rust crate: `btfmt` (`Cargo.toml`)
- Purpose: format bpftrace scripts using `tree-sitter-bpftrace`, with CLI and LSP support
- Runtime integration: VS Code extension launches `btfmt lsp`

## Communication Rule

- Use Chinese when chatting with users; use English when writing code or documentation.

## Commands (Build / Test / Format)

### Taskfile

- Build: `task build`
- Test: `task test`
- Test (tools fixtures): `task test-tools`
- Format Rust code: `task fmt`
- LSP smoke test: `task lsp-smoke`
- Validate release config: `goreleaser check .goreleaser.yaml`

### Direct Cargo Commands

- Build: `cargo build --bin btfmt`
- Build release: `cargo build --release --bin btfmt`
- Test: `cargo test --all-targets --all-features`
- Format: `cargo fmt --all`
- Format check: `cargo fmt --all -- --check`
- Lint: `cargo clippy --all-targets --all-features -- -D warnings`

### Run The Formatter

- Format to stdout: `./btfmt <file.bt>`
- Write in place: `./btfmt -i <file.bt>`
- Write to file: `./btfmt -w <file.bt>`
- Start LSP server: `./btfmt lsp`

## Running A Single Rust Test

- Run one integration test target: `cargo test --test formatter`
- Run one test by name: `cargo test formats_bpftrace_tools_tree`
- Run with output: `cargo test -- --nocapture`

## Project Layout

- `src/cli.rs`: CLI entry point logic
- `src/config.rs`: configuration types and loader
- `src/format.rs`: formatter implementation
- `src/lsp.rs`: LSP server implementation
- `src/parse.rs`: tree-sitter parsing and diagnostics
- `src/text.rs`: text offset/range helpers
- `tests/`: Rust integration tests
- `tests/testdata/`: input fixtures
- `tests/testdata/golden/`: expected Rust formatter output
- `bpftrace/tools/`: upstream bpftrace tool scripts for acceptance tests
- `vscode-extension/`: VS Code extension client

## Parser Policy

- The parser is provided by the `tree-sitter-bpftrace` crate from crates.io.
- Do not add generated parser sources to this repository.
- If grammar behavior needs changes, prefer upstreaming or pinning/updating `tree-sitter-bpftrace`.

## Configuration Files

- Default config values live in `src/config.rs`.
- Example config is `.btfmt.json` at repository root.
- Config search order:
  1. `-config` flag or explicit LSP config path
  2. `.btfmt.json` in the document/current directory or parents
  3. `~/.btfmt.json` for CLI
  4. built-in defaults

## Testing Conventions

- Unit tests may live beside Rust modules.
- Integration tests live in `tests/*.rs`.
- Golden tests read fixtures from `tests/testdata/` and compare to `tests/testdata/golden/`.
- Acceptance tests parse/format files under `bpftrace/tools`.
- LSP behavior is covered by `tests/lsp_smoke.rs`, which runs `scripts/lsp_smoke.py` against the Cargo-built binary.

## Code Style Guidelines

### Formatting

- Use `cargo fmt --all`.
- Keep code idiomatic and simple; prefer small functions when behavior is independently testable.

### Imports

- Use standard Rust import grouping as produced by rustfmt.

### Naming

- Public Rust types/functions: `PascalCase` for types, `snake_case` for functions.
- Test names: descriptive `snake_case`.

### Errors And Control Flow

- Use `anyhow::Result` at CLI/application boundaries.
- Use early returns for error handling and simple branch exits.
- CLI should print user-facing errors and exit non-zero through `src/main.rs`.

### Comments

- Add comments only for non-obvious behavior or integration constraints.

## Linting / Static Analysis

- Primary checks: `cargo fmt --all -- --check`, `cargo clippy --all-targets --all-features -- -D warnings`, and `cargo test --all-targets --all-features`.

## Editing Guidance For Agents

- Keep formatter behavior changes small and update golden fixtures/tests when output changes.
- Maintain config defaults and JSON field names when adding config options.
- Keep CLI and LSP formatting behavior consistent where feasible.
- Avoid committing build outputs such as `target/` or `btfmt`.

## Release Automation

- `.github/workflows/release.yml` builds binaries and platform-specific VSIX packages on native runners.
- `.goreleaser.yaml` generates checksums and changelog content, then publishes the assembled artifacts.
- Release tags must match the versions in `Cargo.toml` and `vscode-extension/package.json`.
- Manual release workflow runs generate an unpublished snapshot artifact bundle.

## Quick References

- Main entrypoint: `src/main.rs`
- CLI: `src/cli.rs`
- Formatter: `src/format.rs`
- LSP: `src/lsp.rs`
- Config: `src/config.rs`
