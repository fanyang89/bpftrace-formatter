# btfmt LSP

VS Code extension for formatting bpftrace scripts.

## Features

- Syntax highlighting for `.bt` files
- Format on save
- Format document command (`Shift+Alt+F`)

## Installation

Download the `.vsix` file for your platform from [Releases](https://github.com/fanyang89/bpftrace-formatter/releases) and install via "Extensions: Install from VSIX..." command.

The extension includes the btfmt binary - no additional installation required.

## Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `btfmt.serverPath` | `btfmt` | Path to btfmt binary (uses bundled binary by default) |
| `btfmt.configPath` | `""` | Path to `.btfmt.json` configuration file |
