# btfmt

A formatter for bpftrace scripts with VS Code integration.

## Features

- Format bpftrace scripts with consistent indentation, spacing, and structure
- VS Code extension with bundled binary - install and use immediately
- Language Server Protocol (LSP) support for editor integration
- Configurable formatting rules via JSON configuration file
- Preserves comments and shebangs

## Installation

### VS Code Extension (Recommended)

Install the `btfmt-lsp` extension from the [Releases](https://github.com/fanyang89/bpftrace-formatter/releases) page:

1. Download the `.vsix` file for your platform (e.g., `btfmt-lsp-0.0.2@linux-x64.vsix`)
2. In VS Code, press `Ctrl+Shift+P` and run "Extensions: Install from VSIX..."
3. Select the downloaded file

The extension includes the btfmt binary - no additional installation required.

### CLI Binary

Download the pre-built binary from [Releases](https://github.com/fanyang89/bpftrace-formatter/releases):

| Platform | File |
|----------|------|
| Linux x64 | `btfmt-linux-amd64.tar.gz` |
| Linux ARM64 | `btfmt-linux-arm64.tar.gz` |
| macOS x64 | `btfmt-darwin-amd64.tar.gz` |
| macOS ARM64 | `btfmt-darwin-arm64.tar.gz` |
| Windows x64 | `btfmt-windows-amd64.zip` |
| Windows ARM64 | `btfmt-windows-arm64.zip` |

Extract and add to your PATH:

```bash
tar -xzf btfmt-linux-amd64.tar.gz
sudo mv btfmt /usr/local/bin/
```

### Build from Source

```bash
go install github.com/fanyang89/bpftrace-formatter/cmd/btfmt@latest
```

Or clone and build:

```bash
git clone https://github.com/fanyang89/bpftrace-formatter.git
cd bpftrace-formatter
go build ./cmd/btfmt
```

## Usage

### Format a file

```bash
btfmt script.bt          # Print formatted output to stdout
btfmt -w script.bt       # Write result back to file
btfmt -w *.bt            # Format multiple files
```

### Example

Before:

```bpftrace
#!/usr/bin/env bpftrace
tracepoint:syscalls:sys_enter_openat{printf("openat: %s\n",str(args.filename));}
tracepoint:syscalls:sys_enter_openat/pid==1234/{@opens[pid]=count();}
```

After:

```bpftrace
#!/usr/bin/env bpftrace

tracepoint:syscalls:sys_enter_openat
{
    printf("openat: %s\n", str(args.filename));
}

tracepoint:syscalls:sys_enter_openat /pid == 1234/
{
    @opens[pid] = count();
}
```

### CLI Options

```
btfmt [options] <file.bt> [file2.bt ...]

Options:
  -w                     Write result to source file
  -i                     Edit files in place (same as -w)
  -c, -config <file>     Path to configuration file
  -v, -verbose           Enable verbose output
  -generate-config       Generate default configuration file
  -version               Show version information
  -help                  Show help message
```

## Configuration

btfmt looks for configuration in this order:

1. File specified with `-config` flag
2. `.btfmt.json` in current directory or parent directories
3. `~/.btfmt.json` in home directory
4. Built-in defaults

Generate a default configuration file:

```bash
btfmt -generate-config
```

Example `.btfmt.json`:

```json
{
  "indent": {
    "size": 4,
    "use_spaces": true
  },
  "spacing": {
    "around_operators": true,
    "around_commas": true
  },
  "line_breaks": {
    "empty_lines_between_probes": 1
  },
  "blocks": {
    "brace_style": "next_line"
  }
}
```

### Configuration Options

| Section | Option | Default | Description |
|---------|--------|---------|-------------|
| `indent` | `size` | 4 | Spaces/tabs per indent level |
| `indent` | `use_spaces` | true | Use spaces instead of tabs |
| `spacing` | `around_operators` | true | Space around `=`, `+`, `-`, etc. |
| `spacing` | `around_commas` | true | Space after commas |
| `spacing` | `before_block_start` | true | Space before `{` |
| `spacing` | `after_keywords` | true | Space after `if`, `while`, etc. |
| `line_breaks` | `empty_lines_between_probes` | 1 | Empty lines between probe blocks |
| `line_breaks` | `empty_lines_after_shebang` | 1 | Empty lines after shebang |
| `blocks` | `brace_style` | "next_line" | `"same_line"`, `"next_line"`, or `"gnu"` |

## VS Code Extension

The VS Code extension provides:

- Syntax highlighting for `.bt` files
- Format on save (enable in VS Code settings)
- Format document command (`Shift+Alt+F`)

### Extension Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `btfmt.serverPath` | `btfmt` | Path to btfmt binary (uses bundled binary by default) |
| `btfmt.configPath` | `""` | Path to `.btfmt.json` configuration file |

## License

[Unlicense](LICENSE) (Public Domain)
