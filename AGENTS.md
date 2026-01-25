# AGENTS

This file is the operational guide for agentic coding tools in this repository.
It summarizes the commands and style conventions observed in the codebase.

## Repository Overview

- Go module: `github.com/fanyang89/bpftrace-formatter` (`go.mod`)
- Go version: `1.24.6` (`go.mod`)
- Purpose: format bpftrace scripts using ANTLR-generated parser + Go visitor

## Communication Rule

- Use Chinese when chatting with users; use English when writing code or documentation.

## Commands (Build / Test / Format)

### Taskfile

- Build: `task build`
- Test: `task test`
- Test (tools fixtures): `task test-tools`
- Format Go code: `task fmt`
- Compile grammar: `task compile-grammar`

### Direct Go commands

- Dependencies: `go mod tidy`
- Build: `go build ./cmd/btfmt`
- Build (README example): `go build -o bpftrace-formatter`
- Test (all): `go test ./...`
- Test (verbose): `go test -v`
- Format: `go fmt ./...`
- Grammar compile (direct): `uvx --from antlr4-tools antlr4 -Dlanguage=Go -o parser bpftrace.g4`

### Run the formatter

- Format to stdout: `./btfmt <file.bt>`
- Write in place: `./btfmt -i <file.bt>`
- Write to file: `./btfmt -w <file.bt>`

## Running a Single Test (Go)

Go does not expose a custom test runner here; use `go test` with `-run`.

- Run a single test across all packages:
  `go test ./... -run TestName`
- Run a single test in one package:
  `go test ./formatter -run TestASTFormatter_GoldenFiles_DefaultConfig`
- Run tests with fixtures from bpftrace tools:
  `BPFTRACE_TOOLS_DIR=bpftrace/tools go test ./... -run TestName`

## Project Layout

- `cmd/btfmt/`: CLI entry point and tests
- `formatter/`: AST formatter + visitor + golden tests
- `config/`: configuration types and loader
- `parser/`: generated ANTLR parser/lexer (do not edit by hand)
- `bpftrace.g4`: grammar source (edit here, then regenerate parser)
- `testdata/`: input fixtures
- `testdata/golden/`: expected formatted output
- `bpftrace/tools/`: upstream bpftrace tool scripts for acceptance tests

## Generated Code Policy

- Files in `parser/` are generated from `bpftrace.g4`.
- Do not hand-edit `parser/*.go`; regenerate via `task compile-grammar`.
- Grammar changes should update golden tests and fixtures as needed.

## Configuration Files

- Default config values in `config/config.go`.
- Example config in `.btfmt.json` (root).
- Config search order (see `cmd/btfmt/main.go`):
  1. `-config` flag
  2. `.btfmt.json` in cwd or parents
  3. `~/.btfmt.json`
  4. built-in defaults

## Testing Conventions

- Tests live alongside packages (`*_test.go`).
- Golden tests read fixtures from `testdata/` and compare to `testdata/golden/`.
- Tests use `t.Helper()`, `t.TempDir()`, `t.Cleanup()`, and `t.Fatalf()`.
- Acceptance tests may read `bpftrace/tools` via `BPFTRACE_TOOLS_DIR`.

## Code Style Guidelines (Observed)

### Formatting

- Use `gofmt`/`go fmt` formatting.
- Tabs for indentation (standard Go formatting).
- Long expected strings in tests are built via `"..." +` concatenation.

### Imports

- Group standard library imports first, then a blank line, then module imports.
- Prefer explicit imports over dot/blank unless required (none used here).

### Naming

- Exported types/functions: `PascalCase` with doc comments.
- Unexported helpers: `camelCase` with inline comments for sections.
- Test names: `TestXxx` matching Go conventions.

### Errors and Control Flow

- Early returns on error (`if err != nil { return ... }`).
- Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- CLI layer prints user-facing errors and exits non-zero; helpers return errors.

### Comments

- Use doc comments for exported APIs.
- Use inline comments to clarify formatting steps in visitors/formatters.

## Linting / Static Analysis

- No lint configuration detected (`.golangci.yml` / `.editorconfig` not present).
- Use `go fmt` and `go test` as primary checks.

## External Rules Files

- `.cursor/rules/`: not present.
- `.cursorrules`: not present.
- `.github/copilot-instructions.md`: not present.

## Editing Guidance for Agents

- Avoid modifying `parser/` (generated).
- Keep formatter changes small and update golden fixtures/tests if output changes.
- Maintain config defaults and JSON tags when adding config options.
- Keep file I/O and CLI behavior consistent with `cmd/btfmt/main.go`.

## Quick References

- Main entrypoint: `cmd/btfmt/main.go`
- Formatting core: `formatter/ast_formatter.go`
- Visitor logic: `formatter/ast_visitor.go`
- Config loader: `config/loader.go`
- Default config: `config/config.go`
