## 2026-03-02 - [Optimization of ASTFormatter Core Paths]
**Learning:** In a visitor-based formatter, core string handling and indentation logic are frequent operations. Using `strings.Repeat` for indentation and multiple `strings` package scans on every token introduces unnecessary heap allocations and CPU cycles. ASCII fast paths for common whitespace checks can significantly bypass expensive Unicode processing.

**Action:**
1. Replace `strings.Repeat` with a manual loop calling `output.WriteByte` for small, frequent allocations like indentation.
2. Cache intermediate results (like `isWhitespace`) within a single token's processing to avoid redundant re-computation.
3. Use specialized byte-level scans (e.g., `strings.LastIndexByte`) instead of more general string-level functions when looking for newlines in tokens.
4. Provide ASCII fast paths for character-level classification.
