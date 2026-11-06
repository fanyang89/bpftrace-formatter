# Selection formatting

The LSP advertises `documentRangeFormattingProvider` and handles
`textDocument/rangeFormatting` (VS Code's **Format Selection**). It uses the same
formatter and document configuration lookup as full-document formatting. As with
full-document formatting, editor `tabSize` and `insertSpaces` options do not
override `.btfmt.json` or `btfmt.configPath`.

## Effective range

A nonempty selection expands to complete syntax boundaries:

- A partial expression expands to its containing statement, including its
  semicolon. Multiple statements expand to the contiguous siblings they touch.
- Selecting a condition or a block brace expands to the owning control-flow
  statement or probe. An `if` includes its `else` branches.
- Selecting probe names or predicates expands to complete probes. A selection
  spanning probes includes those probes and the text between them.
- Leading horizontal whitespace, trailing horizontal whitespace and at most one
  following line ending are included. Trailing comments on the same source line
  belong to the selected statement. Comments between selected statements remain
  included; preceding standalone comments are included only when selected.

The returned edit's `range` is this effective range, which may exceed the original
selection. All bytes outside it remain untouched, including existing probe
separators and line endings. Formatting is taken from the corresponding syntax in
the full-document result, so nested indentation, brace styles, tabs, predicates
and comment configuration match full formatting. A selected statement inside a
one-line block may introduce a newline before and after itself without changing
the surrounding braces. Selected line endings use the formatter's LF convention;
unselected CRLF line endings are preserved.

## Safe no-ops and errors

Empty or whitespace-only selections, comment-only selections, documents with
syntax errors, and selections touching preprocessor directives, shebangs or raw
C structures return no edits. Ordinary statements following a preprocessor
directive can still be formatted. Formatting also returns no edits if syntax
cannot be safely matched to the full-format result or the replacement would
produce invalid syntax. An already formatted selection returns no edits.

Standalone `// btfmt: off` and `// btfmt: on` comments are protected. Selections
whose effective range touches either directive or the disabled region return no
edits; an unmatched `off` protects the rest of the document. Enabled neighboring
statements remain formatable. Directive detection uses parsed line comments,
not text embedded in strings or block comments.

Positions use UTF-16 code units. Reversed ranges, out-of-document positions and
positions splitting a surrogate pair produce an `InvalidParams` error. Config
loading and shared formatter failures also return a descriptive `InvalidParams`
error. The implementation formats the full document internally, then extracts
one bounded edit; it does not currently optimize work to the selected fragment.
