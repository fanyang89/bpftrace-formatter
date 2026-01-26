# btfmt LSP (VS Code extension)

This extension wires the `btfmt` LSP server into VS Code.

## Requirements

- Build or install `btfmt` so it is available in `PATH`, or set an absolute path in settings.
- The server is started as `btfmt lsp` over stdio.

## Settings

- `btfmt.serverPath`: Path to the `btfmt` executable (default: `btfmt`).
- `btfmt.configPath`: Path to `.btfmt.json` (absolute or workspace-relative).

## Development

```bash
cd vscode-extension
npm install
npm run compile
```

Press `F5` in VS Code to launch the extension host.
