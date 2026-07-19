# Configuration

btfmt accepts partial JSON configuration. Missing sections and options use built-in defaults, while unknown fields and invalid values are rejected.

## Search Order

The CLI and LSP intentionally use slightly different fallback rules.

### CLI

1. `--config <path>`
2. The nearest `.btfmt.json` in the current directory or an ancestor
3. `~/.btfmt.json`
4. Built-in defaults

### LSP

1. `btfmt.configPath` from the editor, resolved relative to the matching workspace root when necessary
2. The nearest `.btfmt.json` next to the document or in an ancestor
3. Built-in defaults

Generate the complete default file with:

```bash
btfmt --generate-config
```

Use `--config-output <path>` to choose another destination. Existing files are preserved unless `--force` is also supplied.

## Complete Example

```json
{
  "indent": {
    "size": 4,
    "use_spaces": true
  },
  "spacing": {
    "around_operators": true,
    "around_commas": true,
    "around_parentheses": false,
    "around_brackets": false,
    "before_block_start": true
  },
  "line_breaks": {
    "empty_lines_between_probes": 1,
    "empty_lines_after_shebang": 1
  },
  "comments": {
    "preserve_inline": true
  },
  "blocks": {
    "brace_style": "next_line",
    "indent_statements": true
  }
}
```

## Reference

| Option | Type | Default | Validation |
| --- | --- | --- | --- |
| `indent.size` | integer | `4` | `1` through `16` |
| `indent.use_spaces` | boolean | `true` | Use tabs when `false` |
| `spacing.around_operators` | boolean | `true` | Controls supported binary and assignment operators |
| `spacing.around_commas` | boolean | `true` | Controls spacing after commas |
| `spacing.around_parentheses` | boolean | `false` | Controls spaces inside parentheses |
| `spacing.around_brackets` | boolean | `false` | Controls spaces inside brackets |
| `spacing.before_block_start` | boolean | `true` | Controls spacing before `{` where braces share a line |
| `line_breaks.empty_lines_between_probes` | integer | `1` | `0` through `5` |
| `line_breaks.empty_lines_after_shebang` | integer | `1` | `0` through `5` |
| `comments.preserve_inline` | boolean | `true` | Preserves inline comments where syntax allows |
| `blocks.brace_style` | string | `"next_line"` | `"same_line"`, `"next_line"`, or `"gnu"` |
| `blocks.indent_statements` | boolean | `true` | Indents statements inside blocks |

## VS Code Settings

| Setting | Default | Scope |
| --- | --- | --- |
| `btfmt.serverPath` | `""` | Machine-overridable; empty uses the bundled server, then PATH as a fallback |
| `btfmt.configPath` | `""` | Resource; absolute or relative to the matching workspace root |

Both settings are restricted in VS Code Restricted Mode. Current-document formatting and static language features remain available, while workspace indexing and system metadata access stay disabled.
