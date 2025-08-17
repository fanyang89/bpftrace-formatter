# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based bpftrace script formatter that uses ANTLR4 to parse and format bpftrace scripts with consistent indentation, spacing, and structure. The tool provides a robust, grammar-based approach to parsing bpftrace syntax and reformatting it according to configurable rules.

## Development Commands

### Building

```bash
go mod tidy          # Ensure dependencies are up to date
task build           # Build the binary using Taskfile
go build ./cmd/btfmt # Direct build command
```

### Grammar Compilation

```bash
task compile-grammar # Compile bpftrace.g4 grammar file using ANTLR4
uvx --from antlr4-tools antlr4 -Dlanguage=Go -o parser bpftrace.g4  # Direct command
```

### Testing

```bash
go test -v           # Run all tests with verbose output
go test -run TestName  # Run specific test function
```

### Running

```bash
./btfmt <file.bt>    # Format a bpftrace script file
./btfmt -i <file.bt> # Format file in place
./btfmt -w <file.bt> # Write result to source file
```

## Architecture

### Core Components

**ANTLR Grammar** (`bpftrace.g4`):

- Comprehensive grammar definition for bpftrace syntax
- Supports probe types, predicates, expressions, statements, and built-in functions
- Defines lexical structure including operators, keywords, and literals

**AST Formatter** (`formatter/ast_formatter.go`):

- Central formatting engine with configurable options
- Handles indentation, spacing, and structural formatting
- Provides methods for writing operators, keywords, brackets, and blocks with appropriate spacing

**AST Visitor** (`formatter/ast_visitor.go`):

- Implements visitor pattern for traversing the ANTLR parse tree
- Handles specific formatting for each AST node type (probes, statements, expressions, etc.)
- Maintains state for proper spacing between probes and structural elements

**Configuration System** (`config/`):

- `config.go`: Defines configuration structure with default values
- `loader.go`: Handles configuration loading from files with fallback hierarchy

### Key Processing Flow

1. **Grammar Compilation**: ANTLR4 compiles `bpftrace.g4` into Go parser code
2. **Tokenization**: Input is tokenized using the generated lexer
3. **Parsing**: ANTLR parser builds AST from token stream
4. **AST Traversal**: Visitor pattern traverses the AST and applies formatting rules
5. **Output Generation**: Formatter writes formatted output with proper indentation and spacing

### Configuration Hierarchy

The formatter loads configuration in this order:

1. File specified with `-config` flag
2. `.btfmt.json` in current directory or parent directories (searching upwards)
3. `~/.btfmt.json` in home directory
4. Built-in defaults if no config file is found

### Supported bpftrace Features

- **Probe Types**: tracepoint, kprobe, kretprobe, uprobe, uretprobe, usdt, profile, interval, software, hardware, watchpoint, asyncwatchpoint, BEGIN, END
- **Predicates**: `/condition/` syntax with proper formatting
- **Statements**: assignments, function calls, control flow (if/while/for), return, clear, delete, exit, print, printf
- **Expressions**: Logical, arithmetic, relational, shift, and unary operations
- **Built-in Functions**: printf, sprintf, system, exit, count, sum, avg, min, max, stats, hist, lhist, delete, clear, print, cat, join, time, strftime, str, strerror, kaddr, uaddr, ntop, pton, reg, kstack, ustack, ksym, usym, cgroupid, macaddr, nsecs, elapsed, cpu, pid, tid, uid, gid, comm, curtask, rand, ctx, args, retval, probe, username
- **Data Types**: Variables, maps, strings, numbers (decimal, hex, octal, binary)
- **Comments**: Both `//` and `#` style comments with preservation
- **Shebang**: `#!` line handling at the beginning of files

### Configuration Options

**Indentation**:

- `size`: Number of spaces/tabs per indent level (default: 4)
- `use_spaces`: Use spaces instead of tabs (default: true)

**Spacing**:

- `around_operators`: Space around =, +, -, etc. (default: true)
- `around_commas`: Space after commas (default: true)
- `around_parentheses`: Space inside parentheses (default: false)
- `around_brackets`: Space inside brackets (default: false)
- `before_block_start`: Space before { (default: true)
- `after_keywords`: Space after if, while, etc. (default: true)

**Line Breaks**:

- `max_line_length`: Maximum line length before wrapping (default: 80)
- `break_long_statements`: Break long statements across lines (default: true)
- `empty_lines_between_probes`: Number of empty lines between probes (default: 1)
- `empty_lines_after_shebang`: Number of empty lines after shebang (default: 1)

**Comments**:

- `preserve_inline`: Keep inline comments on same line (default: true)
- `align_inline`: Align inline comments (default: false)
- `indent_level`: Indent level for standalone comments (default: 0)

**Probes**:

- `align_predicates`: Align predicates with probe definitions (default: false)
- `sort_probes`: Sort probes alphabetically (default: false)
- `group_by_type`: Group probes by type (default: false)

**Blocks**:

- `brace_style`: "same_line", "next_line", "gnu" (default: "next_line")
- `indent_statements`: Indent statements inside blocks (default: true)
- `empty_line_in_blocks`: Add empty lines in large blocks (default: false)
- `align_assignments`: Align assignment operators (default: false)

## File Structure

- `bpftrace.g4`: ANTLR4 grammar definition for bpftrace syntax
- `cmd/btfmt/main.go`: Main entry point with CLI argument handling
- `config/`: Configuration loading and management
- `formatter/`: AST-based formatting logic
- `parser/`: Generated ANTLR4 parser and lexer code
- `testdata/`: Test input files for validation
- `Taskfile.yml`: Task automation for build and grammar compilation

## Important Patterns

1. **Grammar-First Approach**: All parsing is based on the formal ANTLR grammar, not regex patterns
2. **Visitor Pattern**: The AST visitor handles each node type with specific formatting logic
3. **Configuration Hierarchy**: Config files are searched in a specific order with defaults
4. **Stateful Formatting**: The formatter maintains indentation level and spacing state during traversal
5. **Comprehensive Coverage**: The grammar supports nearly all bpftrace language features

## Testing Strategy

The formatter is tested with various bpftrace script examples in the `testdata/` directory:

- `test_input.bt`: Basic probe formatting test
- `test_script.bt`: Complex script with multiple probe types
- `test_operators.bt`: Operator and expression formatting tests
