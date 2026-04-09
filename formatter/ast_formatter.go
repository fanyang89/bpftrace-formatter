package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/antlr4-go/antlr/v4"
	"github.com/fanyang89/bpftrace-formatter/config"
	"github.com/fanyang89/bpftrace-formatter/parser"
)

// ASTFormatter formats bpftrace scripts using ANTLR AST
//
// Thread Safety: This type is NOT thread-safe and should not be used concurrently
// from multiple goroutines. Create a new instance for each concurrent formatting operation.
type ASTFormatter struct {
	config *config.Config
	output bytes.Buffer
	visitor *ASTVisitor

	// Formatting state - reset before each FormatTree call
	state formatterState

	// Cached indentation strings to avoid repeated loops and allocations
	indentCache []string
}

// formatterState encapsulates the mutable state during formatting
type formatterState struct {
	indentLevel    int
	lastWasNewline bool
	needIndent     bool
	lineLength     int
	pendingSpace   bool
}

type syntaxErrorListener struct {
	*antlr.DefaultErrorListener
	errors []error
}

func newSyntaxErrorListener() *syntaxErrorListener {
	return &syntaxErrorListener{DefaultErrorListener: antlr.NewDefaultErrorListener()}
}

func (l *syntaxErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	l.errors = append(l.errors, fmt.Errorf("line %d:%d: %s", line, column, msg))
}

func (l *syntaxErrorListener) Err() error {
	count := len(l.errors)
	if count == 0 {
		return nil
	}
	if count == 1 {
		return fmt.Errorf("parse failed: %w", l.errors[0])
	}
	return fmt.Errorf("parse failed with %d syntax error(s): %v", count, l.errors[0])
}

// NewASTFormatter creates a new AST-based formatter
func NewASTFormatter(cfg *config.Config) *ASTFormatter {
	f := &ASTFormatter{
		config:      cfg,
		indentCache: make([]string, 0, 10),
	}
	f.visitor = NewASTVisitor(f)
	f.prepareIndentCache()
	return f
}

const maxIndentCache = 32

// prepareIndentCache pre-calculates indentation strings for common levels
func (f *ASTFormatter) prepareIndentCache() {
	if f.indentCache != nil {
		return
	}

	f.indentCache = make([]string, maxIndentCache)
	indentChar := "\t"
	if f.config.Indent.UseSpaces {
		indentChar = strings.Repeat(" ", f.config.Indent.Size)
	}

	for i := 0; i < maxIndentCache; i++ {
		f.indentCache[i] = strings.Repeat(indentChar, i)
	}
}

// Format formats the given bpftrace script using AST
func (f *ASTFormatter) Format(input string) (string, error) {
	tree, err := ParseBpftrace(input)
	if err != nil {
		return "", err
	}
	return f.FormatTree(tree), nil
}

// FormatTree formats a pre-parsed AST. This avoids re-parsing when the tree
// is already available (e.g. from the LSP document store).
func (f *ASTFormatter) FormatTree(tree antlr.Tree) string {
	// Initialize visitor with this formatter if not already done
	if f.visitor.formatter == nil {
		f.visitor.formatter = f
	}

	// Reset all formatting state
	f.resetState()
	f.output.Reset()

	// Visit the AST tree
	f.visitor.Visit(tree)

	return string(bytes.TrimRightFunc(f.output.Bytes(), unicode.IsSpace))
}

// resetState resets all mutable state to initial values
func (f *ASTFormatter) resetState() {
	f.state = formatterState{
		indentLevel:    0,
		lastWasNewline: true,
		needIndent:     false,
		lineLength:     0,
		pendingSpace:   false,
	}
}

// ParseBpftrace parses a bpftrace script and returns the AST.
// It uses a two-stage strategy: SLL prediction mode first (fast), falling back
// to full LL mode only when the input is ambiguous.
func ParseBpftrace(input string) (parser.IProgramContext, error) {
	inputStream := antlr.NewInputStream(input)
	lexer := parser.NewbpftraceLexer(inputStream)
	errorListener := newSyntaxErrorListener()
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)

	tree, sllErr := parseProgramWithMode(tokenStream, antlr.PredictionModeSLL)
	if lexerErr := errorListener.Err(); lexerErr == nil && sllErr == nil {
		return tree, nil
	}

	tokenStream.Seek(0)
	tree, llErr := parseProgramWithMode(tokenStream, antlr.PredictionModeLL)

	if lexerErr := errorListener.Err(); lexerErr != nil {
		return nil, lexerErr
	}
	if llErr != nil {
		return nil, llErr
	}
	return tree, nil
}

func parseProgramWithMode(tokenStream *antlr.CommonTokenStream, mode int) (parser.IProgramContext, error) {
	errorListener := newSyntaxErrorListener()

	p := parser.NewbpftraceParser(tokenStream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errorListener)
	p.GetInterpreter().SetPredictionMode(mode)

	tree := p.Program()
	if err := errorListener.Err(); err != nil {
		return nil, err
	}
	return tree, nil
}

// writeString writes a string to the output
func (f *ASTFormatter) writeString(s string) {
	if s == "" {
		return
	}

	isWS := f.isWhitespace(s)
	f.handlePendingSpace(s, isWS)
	f.writeWithIndent(s, isWS)
	f.updateLineTracking(s)
}

// handlePendingSpace handles any pending space before writing content
func (f *ASTFormatter) handlePendingSpace(s string, isWS bool) {
	if !f.state.pendingSpace || isWS {
		return
	}

	if f.state.needIndent {
		f.state.pendingSpace = false
		return
	}

	tokenLen := len(s)
	if idx := strings.IndexByte(s, '\n'); idx >= 0 {
		tokenLen = idx
	}

	if f.shouldWrap(1 + tokenLen) {
		f.state.pendingSpace = false
		f.writeNewline()
	} else {
		f.state.pendingSpace = false
		f.output.WriteByte(' ')
		f.state.lineLength++
		f.state.lastWasNewline = false
		f.state.needIndent = false
	}
}

// writeWithIndent writes the string with proper indentation if needed
func (f *ASTFormatter) writeWithIndent(s string, isWS bool) {
	if f.state.needIndent && !isWS {
		f.writeIndent()
	}
	f.output.WriteString(s)
}

// updateLineTracking updates line length and newline state after writing
func (f *ASTFormatter) updateLineTracking(s string) {
	idx := strings.LastIndexByte(s, '\n')
	if idx >= 0 {
		f.state.lineLength = len(s) - idx - 1
		f.state.lastWasNewline = idx == len(s)-1
		f.state.needIndent = f.state.lastWasNewline
	} else {
		f.state.lineLength += len(s)
		f.state.lastWasNewline = false
	}
}

// writeIndent writes the current indentation
func (f *ASTFormatter) writeIndent() {
	f.writeIndentLevel(f.state.indentLevel)
}

// writeIndentLevel writes indentation for a specific level
func (f *ASTFormatter) writeIndentLevel(level int) {
	f.state.pendingSpace = false
	if level < 0 {
		level = 0
	}

	// Use cached indentation if available
	if level < len(f.indentCache) {
		indentStr := f.indentCache[level]
		f.output.WriteString(indentStr)
		f.state.lineLength += len(indentStr)
		f.state.lastWasNewline = false
		f.state.needIndent = false
		return
	}

	// Fallback for extremely deep nesting
	count := level
	indentChar := byte('\t')
	if f.config.Indent.UseSpaces {
		indentChar = ' '
		count = level * f.config.Indent.Size
	}

	for i := 0; i < count; i++ {
		f.output.WriteByte(indentChar)
	}

	f.state.lineLength += count
	f.state.lastWasNewline = false
	f.state.needIndent = false
}

// writeNewline writes a newline character
func (f *ASTFormatter) writeNewline() {
	f.output.WriteByte('\n')
	f.state.lastWasNewline = true
	f.state.needIndent = true
	f.state.lineLength = 0
	f.state.pendingSpace = false
}

// writeSpace writes a space if spacing is enabled
func (f *ASTFormatter) writeSpace() {
	f.state.pendingSpace = true
}

func (f *ASTFormatter) writeSpaceNoWrap() {
	f.state.pendingSpace = false
	f.output.WriteByte(' ')
	f.state.lineLength++
	f.state.lastWasNewline = false
}

// writeSpaceIf writes a space conditionally
func (f *ASTFormatter) writeSpaceIf(condition bool) {
	if condition {
		f.writeSpace()
	}
}

// increaseIndent increases the indentation level
func (f *ASTFormatter) increaseIndent() {
	f.state.indentLevel++
}

// decreaseIndent decreases the indentation level
func (f *ASTFormatter) decreaseIndent() {
	if f.state.indentLevel > 0 {
		f.state.indentLevel--
	}
}

// isWhitespace checks if a string contains only whitespace
func (f *ASTFormatter) isWhitespace(s string) bool {
	if len(s) == 0 {
		return true
	}

	// Fast path for ASCII
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b > 127 {
			// Fallback to unicode.IsSpace for non-ASCII
			for _, r := range s[i:] {
				if !unicode.IsSpace(r) {
					return false
				}
			}
			return true
		}
		// ASCII whitespace characters: ' ', '\t', '\n', '\v', '\f', '\r'
		if !(b == ' ' || (b >= 0x09 && b <= 0x0d)) {
			return false
		}
	}
	return true
}

// writeEmptyLines writes the specified number of empty lines
func (f *ASTFormatter) writeEmptyLines(count int) {
	for i := 0; i < count; i++ {
		f.writeNewline()
	}
}

// ensureNewline ensures the output ends with a newline
func (f *ASTFormatter) ensureNewline() {
	if !f.state.lastWasNewline {
		f.writeNewline()
	}
}

// writeOperator writes an operator with appropriate spacing
func (f *ASTFormatter) writeOperator(op string) {
	if f.config.Spacing.AroundOperators {
		f.writeSpace()
	}
	f.writeString(op)
	if f.config.Spacing.AroundOperators {
		f.writeSpace()
	}
}

// writeKeyword writes a keyword with appropriate spacing
func (f *ASTFormatter) writeKeyword(keyword string) {
	f.writeString(keyword)
	if f.config.Spacing.AfterKeywords {
		f.writeSpace()
	}
}

// writeComma writes a comma with appropriate spacing
func (f *ASTFormatter) writeComma() {
	f.writeString(",")
	if f.config.Spacing.AroundCommas {
		f.writeSpace()
	}
}

// writeOpenParen writes an opening parenthesis with appropriate spacing
func (f *ASTFormatter) writeOpenParen() {
	f.writeString("(")
	if f.config.Spacing.AroundParentheses {
		f.writeSpace()
	}
}

// writeCloseParen writes a closing parenthesis with appropriate spacing
func (f *ASTFormatter) writeCloseParen() {
	if f.config.Spacing.AroundParentheses {
		f.writeSpace()
	}
	f.writeString(")")
}

// writeOpenBracket writes an opening bracket with appropriate spacing
func (f *ASTFormatter) writeOpenBracket() {
	f.writeString("[")
	if f.config.Spacing.AroundBrackets {
		f.writeSpace()
	}
}

// writeCloseBracket writes a closing bracket with appropriate spacing
func (f *ASTFormatter) writeCloseBracket() {
	if f.config.Spacing.AroundBrackets {
		f.writeSpace()
	}
	f.writeString("]")
}

// writeBlockStart writes a block start with appropriate spacing and indentation
func (f *ASTFormatter) writeBlockStart() {
	baseIndent := f.state.indentLevel
	braceIndent := f.calculateBraceIndent(baseIndent)

	f.writeBracePrefix()
	f.writeOpeningBrace(braceIndent)

	statementIndent := f.calculateStatementIndent(baseIndent)
	f.state.indentLevel = statementIndent
}

// calculateBraceIndent calculates the indent level for the opening brace
func (f *ASTFormatter) calculateBraceIndent(baseIndent int) int {
	if f.config.Blocks.BraceStyle == "gnu" {
		return baseIndent + 1
	}
	return baseIndent
}

// calculateStatementIndent calculates the indent level for statements inside the block
func (f *ASTFormatter) calculateStatementIndent(baseIndent int) int {
	switch f.config.Blocks.BraceStyle {
	case "gnu":
		if f.config.Blocks.IndentStatements {
			return baseIndent + 2
		}
		return baseIndent + 1
	default:
		if f.config.Blocks.IndentStatements {
			return baseIndent + 1
		}
		return baseIndent
	}
}

// writeBracePrefix writes the appropriate content before the opening brace
func (f *ASTFormatter) writeBracePrefix() {
	switch f.config.Blocks.BraceStyle {
	case "next_line", "gnu":
		f.writeNewline()
	case "same_line":
		if f.config.Spacing.BeforeBlockStart {
			f.writeSpace()
		}
	default:
		if f.config.Spacing.BeforeBlockStart {
			f.writeSpace()
		}
	}
}

// writeOpeningBrace writes the opening brace with proper indentation
func (f *ASTFormatter) writeOpeningBrace(indentLevel int) {
	if f.config.Blocks.BraceStyle == "gnu" {
		f.writeIndentLevel(indentLevel)
	}
	f.writeString("{")
	f.writeNewline()
}

// writeBlockEnd writes a block end with appropriate indentation
func (f *ASTFormatter) writeBlockEnd() {
	indentDelta := f.calculateIndentDelta()
	parentIndent := f.state.indentLevel - indentDelta
	if parentIndent < 0 {
		parentIndent = 0
	}

	braceIndent := parentIndent
	if f.config.Blocks.BraceStyle == "gnu" {
		braceIndent = parentIndent + 1
	}

	f.state.indentLevel = parentIndent
	f.writeClosingBrace(braceIndent)
}

// calculateIndentDelta calculates how much to decrease indent when exiting a block
func (f *ASTFormatter) calculateIndentDelta() int {
	switch f.config.Blocks.BraceStyle {
	case "gnu":
		if f.config.Blocks.IndentStatements {
			return 2
		}
		return 1
	default:
		if f.config.Blocks.IndentStatements {
			return 1
		}
		return 0
	}
}

// writeClosingBrace writes the closing brace with proper indentation
func (f *ASTFormatter) writeClosingBrace(indentLevel int) {
	f.ensureNewline()
	f.writeIndentLevel(indentLevel)
	f.writeString("}")
}

// writeSemicolon writes a semicolon
func (f *ASTFormatter) writeSemicolon() {
	f.writeString(";")
}

func (f *ASTFormatter) shouldWrap(nextLen int) bool {
	if !f.config.LineBreaks.BreakLongStatements {
		return false
	}
	if f.config.LineBreaks.MaxLineLength <= 0 {
		return false
	}
	if f.state.lastWasNewline || f.state.lineLength == 0 {
		return false
	}
	return f.state.lineLength+nextLen > f.config.LineBreaks.MaxLineLength
}
