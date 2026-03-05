## 2026-03-02 - [Optimization of ASTFormatter Core Paths]
**Learning:** In a visitor-based formatter, core string handling and indentation logic are frequent operations. Using `strings.Repeat` for indentation and multiple `strings` package scans on every token introduces unnecessary heap allocations and CPU cycles. ASCII fast paths for common whitespace checks can significantly bypass expensive Unicode processing.

**Action:**
1. Replace `strings.Repeat` with a manual loop calling `output.WriteByte` for small, frequent allocations like indentation.
2. Cache intermediate results (like `isWhitespace`) within a single token's processing to avoid redundant re-computation.
3. Use specialized byte-level scans (e.g., `strings.LastIndexByte`) instead of more general string-level functions when looking for newlines in tokens.
4. Provide ASCII fast paths for character-level classification.

## 2026-03-03 - [Reducing GC Pressure and Buffer Overhead]
**Learning:** Reusing complex objects like `ASTVisitor` and using `bytes.Buffer` instead of `strings.Builder` can significantly reduce GC pressure in performance-critical paths. `strings.Builder.Reset()` clears the underlying capacity, whereas `bytes.Buffer.Reset()` preserves it, making it more efficient for repeated operations like formatting in an LSP server. Additionally, byte-level operations (like `bytes.TrimRight`) are often faster than their string equivalents as they avoid Unicode overhead and extra allocations.

**Action:**
1. Reuse `ASTVisitor` instances in `ASTFormatter` and manually reset state.
2. Prefer `bytes.Buffer` over `strings.Builder` when capacity reuse is desired across operations.
3. Use `switch` statements for frequently called constant set checks (like operators) to avoid slice allocations.
