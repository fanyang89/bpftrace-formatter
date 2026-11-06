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

## Formatting Suppression

Use standalone line comments to preserve hand-formatted sections in both the CLI
and LSP. No configuration option is required:

```bpftrace
BEGIN {
    // btfmt: off
    @first  = 1;   @second = 2;
    // btfmt: on
    printf("ready\n");
}
```

- The exact, case-sensitive comment text is `// btfmt: off` or `// btfmt: on`.
  Leading and trailing whitespace is allowed; LF and CRLF line
  endings are supported. Extra comment text does not count as a directive.
- The complete `off` line through the complete matching `on` line is preserved
  byte-for-byte, including indentation, blank lines, trailing whitespace and line
  endings. Surrounding code is formatted normally. A suppressed region ending at
  EOF also retains the presence or absence of its final newline.
- Multiple and adjacent regions are supported. Directives do not nest: repeated
  `off` comments while disabled have no effect, and the first standalone `on`
  resumes formatting. An unmatched `off` extends to EOF; an unmatched `on` is an
  ordinary comment and does not disable anything.
- Strings, block comments, inline comments and preprocessor regions do not
  activate directives. Preprocessor regions retain their existing opaque handling;
  place directives outside them to suppress surrounding code. Inline comments
  with exact directive text remain inline even when `comments.preserve_inline` is
  false, so another formatting pass cannot accidentally activate them.
- Suppression affects layout only. The **entire document**, including disabled
  regions, must still parse with the bundled `tree-sitter-bpftrace` grammar, and
  the formatted result is validated too. It cannot hide malformed code or syntax
  unsupported by that grammar. Regions may cross block or expression boundaries;
  the formatter tracks their tokens to resume with the correct block state.

Library integrations can call
`btfmt::suppression::suppression_ranges(source) -> anyhow::Result<Vec<std::ops::Range<usize>>>`.
It validates the source and returns sorted, disjoint, half-open UTF-8 byte ranges
in the original document, including both directive lines. Detect ranges on the
whole document before selecting a formatting range: an earlier `off` can suppress
the selection. Do not independently format fragments that lose this context or
reindent a preserved region; any overlapping edit must retain its bytes exactly.

## VS Code Settings

| Setting | Default | Scope |
| --- | --- | --- |
| `btfmt.serverPath` | `""` | Machine-overridable; empty uses the bundled server, then PATH as a fallback |
| `btfmt.configPath` | `""` | Resource; absolute or relative to the matching workspace root |

Both settings are restricted in VS Code Restricted Mode. Current-document formatting and static language features remain available, while workspace indexing and system metadata access stay disabled.
