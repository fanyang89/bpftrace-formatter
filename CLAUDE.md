# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based bpftrace script formatter that parses and formats bpftrace scripts with consistent indentation, spacing, and structure. The tool uses a token-based approach to parse bpftrace syntax and reformat it according to configurable rules.

## Development Commands

### Building

```bash
go mod tidy          # Ensure dependencies are up to date
go build -o bpftrace-formatter  # Build the binary
```

### Testing

```bash
go test -v           # Run all tests with verbose output
go test -run TestName  # Run specific test function
```

### Running

```bash
./bpftrace-formatter <file.bt>  # Format a bpftrace script file
```

## Architecture

### Core Components

**Formatter Struct** (`main.go:35-50`):

- Central formatting engine with configurable options
- `indentSize`: Number of spaces per indent level (default: 4)
- `useSpaces`: Use spaces instead of tabs (default: true)
- `probeSpacing`: Empty lines between probes (default: 2)
- `commentIndent`: Indentation for comments (default: 2)

**Token System** (`main.go:11-32`):

- `TokenType`: Enum defining token types (Probe, Predicate, BlockStart, etc.)
- `Token`: Struct containing token type, value, and position information
- Tokenizer breaks input into lexical tokens for processing

### Key Processing Flow

1. **Tokenization** (`main.go:73-149`):
   - Uses regex patterns to identify probe definitions
   - Handles shebang lines and comments
   - Supports predicate filtering syntax (`/condition/`)
   - Special handling for END blocks

2. **Formatting** (`main.go:152-220`):
   - Processes tokens with proper indentation levels
   - Maintains spacing between probes
   - Preserves comments and special syntax
   - Handles block structure and nested indentation

### Supported bpftrace Features

- Probe types: tracepoint, kprobe, uprobe, profile
- Predicate filtering with `/condition/` syntax
- Variable assignments and map operations
- printf statements and function calls
- END blocks for cleanup
- Comments (both `//` and shebang `#!`)

### Configuration Methods

- `SetIndentSize(int)`: Change indentation size
- `SetUseSpaces(bool)`: Toggle between spaces and tabs
- `SetProbeSpacing(int)`: Set spacing between probe definitions

## Testing Strategy

Tests in `main_test.go` cover:

- Basic probe formatting with indentation
- Predicate handling and spacing
- Configuration option validation
- Cross-platform newline normalization

## File Structure

- `main.go`: Core formatting logic and tokenizer
- `main_test.go`: Unit tests
- `example.bt`: Formatted output example
- `test_input.bt`: Test case input file

## Important Patterns

1. **Token-based processing**: All formatting works through the token system, not direct string manipulation
2. **Regex-based parsing**: Uses `regexp.MustCompile(probePattern)` for probe detection
3. **Stateful formatting**: Maintains indentation level and probe spacing state during token processing
4. **Comment preservation**: Special handling to maintain comments and shebang lines in original positions
