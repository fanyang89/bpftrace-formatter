# btfmt

[English](README.md) | [中文](README_zh-CN.md)

A formatter for bpftrace scripts with VS Code integration.

## Features

- Format bpftrace scripts with consistent indentation, spacing, and structure
- VS Code extension with bundled binary - install and use immediately
- Language Server Protocol (LSP) support with hover, completion, navigation, and rename
- Context-aware completion for builtins, probe targets, `args` fields, visible variables, maps, and macros
- Cross-file definitions, references, and rename for imported maps and macro families
- Configurable formatting rules via JSON configuration file
- Preserves comments and shebangs

## Installation

### VS Code Extension (Recommended)

Install [btfmt - bpftrace Language Support](https://marketplace.visualstudio.com/items?itemName=fanyang89.btfmt) from the VS Code Marketplace or search for `btfmt` in the Extensions view.

Platform-specific VSIX files are also available from [GitHub Releases](https://github.com/fanyang89/bpftrace-formatter/releases):

1. Download the `.vsix` file for your platform (for example, Linux x64)
2. In VS Code, press `Ctrl+Shift+P` and run "Extensions: Install from VSIX..."
3. Select the downloaded file

The extension includes the btfmt binary - no additional formatter installation required.
Probe target and `args` field completion uses workspace scripts, portable catalogs, and readable kernel metadata without requiring root access.

### CLI Binary

Download the pre-built binary from [Releases](https://github.com/fanyang89/bpftrace-formatter/releases):

| Platform      | File                        |
| ------------- | --------------------------- |
| Linux x64     | `btfmt-linux-amd64.tar.gz`  |
| macOS x64     | `btfmt-darwin-amd64.tar.gz` |
| macOS ARM64   | `btfmt-darwin-arm64.tar.gz` |
| Windows x64   | `btfmt-windows-amd64.zip`   |

Linux ARM64 and Windows ARM64 release artifacts are currently deferred.

Extract and add to your PATH:

```bash
tar -xzf btfmt-linux-amd64.tar.gz
sudo mv btfmt /usr/local/bin/
```

### Build from Source

```bash
cargo install --git https://github.com/fanyang89/bpftrace-formatter.git
```

Or clone and build:

```bash
git clone https://github.com/fanyang89/bpftrace-formatter.git
cd bpftrace-formatter
cargo build --release
./target/release/btfmt --help
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

tracepoint:syscalls:sys_enter_openat
/pid == 1234/
{
    @opens[pid] = count();
}
```

### CLI Options

```
btfmt [options] <file.bt|-> [file2.bt ...]

Options:
  -w                     Write result to source file
  -i                     Edit files in place (same as -w)
  -c, -config <file>     Path to configuration file
  -v, -verbose           Enable verbose output
  --check                Exit non-zero if input is not formatted
  -generate-config       Generate default configuration file
  -config-output <file>  Output path for generated configuration
  --force                Overwrite an existing generated configuration
  -version               Show version information
  -help                  Show help message
```

Read from stdin with an explicit `-`:

```bash
cat script.bt | btfmt -
btfmt --check script.bt
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
btfmt -generate-config --force  # overwrite an existing file
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

| Section       | Option                       | Default     | Description                              |
| ------------- | ---------------------------- | ----------- | ---------------------------------------- |
| `indent`      | `size`                       | 4           | Spaces/tabs per indent level             |
| `indent`      | `use_spaces`                 | true        | Use spaces instead of tabs               |
| `spacing`     | `around_operators`           | true        | Space around `=`, `+`, `-`, etc.         |
| `spacing`     | `around_commas`              | true        | Space after commas                       |
| `spacing`     | `around_parentheses`         | false       | Space inside parentheses                 |
| `spacing`     | `around_brackets`            | false       | Space inside brackets                    |
| `spacing`     | `before_block_start`         | true        | Space before `{`                         |
| `line_breaks` | `empty_lines_between_probes` | 1           | Empty lines between probe blocks         |
| `line_breaks` | `empty_lines_after_shebang`  | 1           | Empty lines after shebang                |
| `comments`    | `preserve_inline`            | true        | Preserve inline comments when possible   |
| `blocks`      | `brace_style`                | "next_line" | `"same_line"`, `"next_line"`, or `"gnu"` |
| `blocks`      | `indent_statements`          | true        | Indent statements inside blocks          |

## VS Code Extension

The VS Code extension provides:

- Syntax highlighting for `.bt` files
- Format on save (enable in VS Code settings)
- Format document command (`Shift+Alt+F`)
- Hover documentation for bpftrace builtins
- Context-aware completion for builtins, providers, probe targets, `args` fields, keywords, visible variables, maps, macro parameters, and imported macros
- Document symbols for probes and macros
- Definitions, references, highlights, and rename for lexical variables, maps, and macro families
- Workspace-aware navigation and rename across imported `.bt` files

The language server can also be started directly for other LSP clients:

```bash
btfmt lsp
```

### Extension Settings

| Setting            | Default | Description                                           |
| ------------------ | ------- | ----------------------------------------------------- |
| `btfmt.serverPath` | `btfmt` | Path to btfmt binary (uses bundled binary by default) |
| `btfmt.configPath` | `""`    | Path to `.btfmt.json` configuration file              |

## License

[Unlicense](LICENSE) (Public Domain)
