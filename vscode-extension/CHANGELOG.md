# Changelog

## 0.4.0

### Added

- Format selected statements and probes with syntax-aware range expansion and preserved surrounding content.
- Preserve hand-formatted regions with standalone `// btfmt: off` and `// btfmt: on` comments in the CLI and editor.
- Add semantic highlighting for variables, maps, macro parameters, and builtins, with TextMate fallback.
- Search workspace maps, macros, and probe blocks with support for unsaved edits, ignore rules, and Restricted Mode.
- Show builtin signature help and active-parameter hints, including while typing incomplete calls.
- Format directories recursively from the CLI and preview unified diffs with `--diff` or `--check --diff`.

### Fixed

- Resolve Marketplace README images and documentation links from the extension subdirectory.
- Refresh workspace symbols when `.ignore` or `.gitignore` files change.
- Match canonical workspace roots to LSP file paths on Windows so unsaved symbols remain searchable.
- Preserve workspace-root path aliases when indexing files and resolving imports on macOS and Windows.
- Unify encoded file URI keys so open buffers replace disk symbols on Windows.

## 0.3.3

- Add an embedded tracepoint catalog for completion without kernel metadata access.

## 0.3.2

- Keep workspace probe targets ahead of large kernel metadata result sets.
- Make completion smoke tests deterministic on hosts with readable tracefs metadata.

## 0.3.1

- Add rootless probe target and `args` field completion.
- Add native language status, structured logs, and serialized server restarts.
- Limit workspace and system access in Restricted Mode.
- Refresh Marketplace identity, documentation, and packaging.

## 0.3.0

- Add contextual completion and builtin hover documentation.
- Add cross-file navigation and rename for imported maps and macro families.
- Add workspace indexing and immutable document caches.
