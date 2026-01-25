# BPFTrace Formatter

A Go-based formatter for bpftrace scripts built on ANTLR4.

## Features

- Grammar-based parsing and formatting for bpftrace scripts
- Configurable indentation, spacing, line breaks, comments, probes, and blocks
- Preserves shebangs and comments
- Supports formatting multiple files and in-place edits

## Build

```bash
go mod tidy
task build
# or
go build ./cmd/btfmt
```

## Usage

### Basic

```bash
./btfmt <file.bt>
```

### In-place / write

```bash
./btfmt -i <file.bt>
./btfmt -w <file1.bt> <file2.bt>
```

### Configuration

```bash
./btfmt -generate-config
./btfmt -config <path/to/.btfmt.json> <file.bt>
```

Options:

- `-c`, `-config <file>`: Path to configuration file
- `-i`: Edit files in place
- `-w`: Write result to source file instead of stdout
- `-v`, `-verbose`: Enable verbose output
- `-generate-config`: Generate default configuration file
- `-config-output <file>`: Output path for generated config (default: .btfmt.json)
- `-help`: Show help message

### Example

Input file (`testdata/test_input.bt`):

```bpftrace
#!/usr/bin/env bpftrace
tracepoint:syscalls:sys_enter_openat{printf("openat: %s\n",str(args.filename));}
tracepoint:syscalls:sys_enter_openat2{printf("openat2: %s\n",str(args->filename));}
tracepoint:syscalls:sys_enter_openat/pid==1234/{@opens[pid]=count();}
```

Output:

```bpftrace
#!/usr/bin/env bpftrace
tracepoint:syscalls:sys_enter_openat {
    printf("openat: %s\n",str(args.filename));
}

tracepoint:syscalls:sys_enter_openat2 {
    printf("openat2: %s\n",str(args->filename));
}

tracepoint:syscalls:sys_enter_openat/pid==1234/ {
    @opens[pid] = count();
}
```

## Configuration

The formatter loads configuration in this order:

1. File specified with `-config`
2. `.btfmt.json` in the current directory or parent directories
3. `~/.btfmt.json` in the home directory
4. Built-in defaults

Configuration is JSON with these top-level sections: `indent`, `spacing`,
`line_breaks`, `comments`, `probes`, `blocks`.

Example:

```json
{
  "indent": { "size": 4, "use_spaces": true },
  "spacing": { "around_operators": true, "around_commas": true },
  "line_breaks": { "empty_lines_between_probes": 1, "max_line_length": 80 },
  "comments": { "preserve_inline": true },
  "probes": { "sort_probes": false },
  "blocks": { "brace_style": "next_line" }
}
```

## Grammar Compilation

```bash
task compile-grammar
```

## Testing

```bash
task test
# or
go test ./...
```

## Project Structure

- `cmd/btfmt/`: CLI entrypoint
- `formatter/`: AST formatter and visitor
- `config/`: Configuration types and loader
- `parser/`: Generated ANTLR parser/lexer
- `bpftrace.g4`: Grammar definition
- `testdata/`: Test scripts and fixtures

## Supported Syntax Features

- Probe definitions and predicates
- Statements (if/while/for, return, print/printf, clear/delete)
- Expressions (logical, arithmetic, relational, unary)
- Built-in functions and maps
- Comments and shebang handling

## Development

This project is developed in Go. Contributions and issue reports are welcome.
