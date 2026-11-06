# Semantic highlighting

`btfmt lsp` supports `textDocument/semanticTokens/full` for open documents.
The server advertises full-document tokens (no range or delta requests). Tokens
are ordered, non-overlapping, single-line segments with UTF-16 positions and
lengths. Multiline comments are split at line endings; CR and LF are excluded.

## Token legend

The advertised `tokenTypes` array has the following stable order:

| Index | Type | Classification |
| --- | --- | --- |
| 0 | `variable` | Lexically scoped `$name` variables and known builtin values |
| 1 | `map` | `@name` and anonymous `@` maps, excluding indexes |
| 2 | `parameter` | Macro parameters (bare, `$`, or `@`) and their resolved uses; `$1` script arguments |
| 3 | `macro` | Macro declarations and references resolved in the current document |
| 4 | `function` | Calls to builtin functions in the bundled catalog |
| 5 | `property` | Field/member identifiers |
| 6 | `type` | Primitive types and identifiers in type syntax |
| 7 | `event` | Probe providers, including their grammar-supported aliases |
| 8 | `comment` | Line/block comments and hashbangs |
| 9 | `string` | String literals |
| 10 | `number` | Integer literals and probe numeric arguments |
| 11 | `keyword` | Grammar-recognized control/declaration keywords and booleans |

The `tokenModifiers` array is `declaration`, `defaultLibrary`, in that order:

- Bit 0 (`declaration`): explicit declarations recorded by symbol analysis,
  including `let`, macro names/parameters, and loop bindings. An initial
  assignment alone does not receive this modifier.
- Bit 1 (`defaultLibrary`): known builtin function calls/values, script
  arguments, and probe providers.

`map` is a custom token type with `variable` as its VS Code supertype and a
TextMate scope mapping. All other token types and modifiers are standard LSP
names. Macro parameters take precedence over builtin spellings; property,
probe-target, and type positions are not mistaken for builtin values or calls.

## Editing and fallback

Each request uses the current in-memory document snapshot, parser tree, and
symbol analysis. Incomplete syntax uses that snapshot's recovery tree/index;
only recognized constructs receive tokens. Unknown calls and unresolved
cross-file macros retain lexical highlighting. No previous version's tokens
are cached; empty or closed documents return an empty token array.

Highlighting works in Restricted Mode without reading workspace files,
resolving imports, inspecting the system, or executing bpftrace. It does not
provide type inference or runtime validation. Classification is limited by
the bundled grammar and builtin catalog (`tree-sitter-bpftrace 0.3.2`).

The extension's TextMate grammar remains the fallback when semantic
highlighting is disabled, unsupported by the theme/client, or unavailable
while the server starts. VS Code controls it through
`editor.semanticHighlighting.enabled`; token colors can be customized with
`editor.semanticTokenColorCustomizations`.
