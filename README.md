# BPFTrace Formatter

A Go tool for formatting bpftrace scripts using an ANTLR-generated parser and an AST visitor.

## Features

- Grammar-based parsing and formatting for bpftrace scripts
- Consistent indentation, spacing, and line breaks
- Preserves comments and shebangs
- Configurable formatting options via JSON
- CLI supports formatting multiple files and in-place edits

## Build

### From source

```bash
go mod tidy
go build ./cmd/btfmt
```

### Taskfile

```bash
task build
```

## Usage

### Basic usage

```bash
./btfmt <file.bt>
```

### Format in place

```bash
./btfmt -i script.bt
```

### Format multiple files

```bash
./btfmt -w file1.bt file2.bt
```

### Generate a default configuration

```bash
./btfmt -generate-config
```

### Use a custom configuration

```bash
./btfmt -config .btfmt.json script.bt
```

### Options

- `-c`, `-config <file>`: Path to configuration file
- `-i`: Edit files in place
- `-w`: Write result to source file instead of stdout
- `-v`, `-verbose`: Enable verbose output
- `-generate-config`: Generate default configuration file
- `-config-output <file>`: Output path for generated config (default: .btfmt.json)
- `-help`: Show help message

### Language Server (LSP)

Run the formatter as an LSP server:

```bash
btfmt lsp
```

The server communicates over stdio (LSP JSON-RPC). Your client/editor must handle JSON-RPC framing (for example, `Content-Length` headers).

## Configuration

The formatter loads configuration in the following order:

1. File specified with `-config`
2. `.btfmt.json` in the current directory or parent directories
3. `~/.btfmt.json`
4. Built-in defaults

If `-config` is set but the file does not exist, the CLI warns and uses built-in defaults without searching for `.btfmt.json`.

### LSP configuration resolution

When running `btfmt lsp`, configuration is resolved per document using:

1. `btfmt.configPath` from client settings, if set and the file exists
2. `.btfmt.json` in the workspace root (or the document directory if no workspace root is provided), searching parent directories up to the filesystem root
3. `~/.btfmt.json`
4. Built-in defaults

If `btfmt.configPath` is set but the file does not exist, the server uses built-in defaults and does not fall back to searching for `.btfmt.json`.

Client-provided settings (everything under `btfmt` except `configPath`) are merged on top of the resolved file config; client settings take precedence.

A minimal example configuration:

```json
{
  "indent": {
    "size": 4,
    "use_spaces": true
  },
  "line_breaks": {
    "empty_lines_between_probes": 1
  }
}
```

Configuration is JSON with these top-level sections: `indent`, `spacing`,
`line_breaks`, `comments`, `probes`, `blocks`.

Run `./btfmt -generate-config` to generate a full default configuration file.

## Example

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

## Grammar Compilation

```bash
task compile-grammar
# or
uvx --from antlr4-tools antlr4 -Dlanguage=Go -o parser bpftrace.g4
```

## Testing

```bash
task test
# or
go test ./...
```

## Project Structure

- `cmd/btfmt/`: CLI entry point and tests
- `formatter/`: AST formatter and visitor
- `config/`: configuration types and loader
- `parser/`: generated ANTLR parser/lexer (do not edit by hand)
- `bpftrace.g4`: grammar definition
- `testdata/`: input fixtures and golden output

## Supported Syntax Features

- Probe definitions and predicates
- Statements (if/while/for, return, print/printf, clear/delete)
- Expressions (logical, arithmetic, relational, unary)
- Built-in functions and maps
- Comments and shebang handling

## Limitations

- Some edge-case bpftrace syntax may need further polishing
- Complex multi-line blocks may still require manual review

## Development

This project is developed in Go. Contributions and issue reports are welcome.
