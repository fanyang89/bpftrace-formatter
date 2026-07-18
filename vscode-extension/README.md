# btfmt LSP

VS Code extension for formatting bpftrace scripts.

## Features

- Syntax highlighting for `.bt` files
- Format on save
- Format document command (`Shift+Alt+F`)
- Context-aware builtin, probe target, `args` field, variable, map, and macro completion
- Document symbols, definitions, references, highlights, and rename

## Installation

Download the `.vsix` file for your platform from [Releases](https://github.com/fanyang89/bpftrace-formatter/releases) and install via "Extensions: Install from VSIX..." command.

The extension includes the btfmt binary - no additional installation required.
Probe target and `args` field completion uses workspace scripts, portable catalogs, and readable kernel metadata without requiring root access.

## Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `btfmt.serverPath` | `btfmt` | Path to btfmt binary (uses bundled binary by default) |
| `btfmt.configPath` | `""` | Path to `.btfmt.json` configuration file |
