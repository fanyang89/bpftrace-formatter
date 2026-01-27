package formatter

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/antlr4-go/antlr/v4"
	"github.com/fanyang89/bpftrace-formatter/config"
	"github.com/fanyang89/bpftrace-formatter/parser"
)

// ASTFormatter formats bpftrace scripts using ANTLR AST
type ASTFormatter struct {
	config         *config.Config
	output         strings.Builder
	indentLevel    int
	lastWasNewline bool
	needIndent     bool
	lineLength     int
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
	return &ASTFormatter{
		config:         cfg,
		indentLevel:    0,
		lastWasNewline: true,
		needIndent:     false,
		lineLength:     0,
	}
}

// Format formats the given bpftrace script using AST
func (f *ASTFormatter) Format(input string) (string, error) {
	// Reset formatter state
	f.output.Reset()
	f.indentLevel = 0
	f.lastWasNewline = true
	f.needIndent = false
	f.lineLength = 0

	// Create ANTLR input stream
	inputStream := antlr.NewInputStream(input)
	lexer := parser.NewbpftraceLexer(inputStream)
	errorListener := newSyntaxErrorListener()
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	bpftraceParser := parser.NewbpftraceParser(tokenStream)
	bpftraceParser.RemoveErrorListeners()
	bpftraceParser.AddErrorListener(errorListener)

	// Parse the program
	tree := bpftraceParser.Program()
	if err := errorListener.Err(); err != nil {
		return "", err
	}

	// Create and use the visitor
	visitor := NewASTVisitor(f)
	visitor.Visit(tree)

	return strings.TrimSpace(f.output.String()), nil
}

// writeString writes a string to the output
func (f *ASTFormatter) writeString(s string) {
	if f.needIndent && !f.isWhitespace(s) {
		f.writeIndent()
	}
	f.output.WriteString(s)

	if s == "" {
		return
	}

	if strings.Contains(s, "\n") {
		lastNewline := strings.LastIndex(s, "\n")
		f.lineLength = len(s) - lastNewline - 1
		f.lastWasNewline = strings.HasSuffix(s, "\n")
		f.needIndent = f.lastWasNewline
		return
	}

	f.lineLength += len(s)
	f.lastWasNewline = false
}

// writeIndent writes the current indentation
func (f *ASTFormatter) writeIndent() {
	f.writeIndentLevel(f.indentLevel)
}

// writeIndentLevel writes indentation for a specific level
func (f *ASTFormatter) writeIndentLevel(level int) {
	if level < 0 {
		level = 0
	}
	var indent string
	if f.config.Indent.UseSpaces {
		indent = strings.Repeat(" ", level*f.config.Indent.Size)
	} else {
		indent = strings.Repeat("\t", level)
	}
	f.output.WriteString(indent)
	f.lineLength += len(indent)
	f.lastWasNewline = false
	f.needIndent = false
}

// writeNewline writes a newline character
func (f *ASTFormatter) writeNewline() {
	f.output.WriteString("\n")
	f.lastWasNewline = true
	f.needIndent = true
	f.lineLength = 0
}

// writeSpace writes a space if spacing is enabled
func (f *ASTFormatter) writeSpace() {
	if f.shouldWrap(1) {
		f.writeNewline()
		f.writeIndent()
		return
	}
	f.output.WriteString(" ")
	f.lineLength++
	f.lastWasNewline = false
}

func (f *ASTFormatter) writeSpaceNoWrap() {
	f.output.WriteString(" ")
	f.lineLength++
	f.lastWasNewline = false
}

// writeSpaceIf writes a space conditionally
func (f *ASTFormatter) writeSpaceIf(condition bool) {
	if condition {
		f.writeSpace()
	}
}

// increaseIndent increases the indentation level
func (f *ASTFormatter) increaseIndent() {
	f.indentLevel++
}

// decreaseIndent decreases the indentation level
func (f *ASTFormatter) decreaseIndent() {
	if f.indentLevel > 0 {
		f.indentLevel--
	}
}

// isWhitespace checks if a string contains only whitespace
func (f *ASTFormatter) isWhitespace(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
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
	if !f.lastWasNewline {
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
	baseIndent := f.indentLevel
	braceIndent := baseIndent
	statementIndent := baseIndent

	switch f.config.Blocks.BraceStyle {
	case "next_line":
		f.writeNewline()
	case "same_line":
		if f.config.Spacing.BeforeBlockStart {
			f.writeSpace()
		}
	case "gnu":
		f.writeNewline()
		braceIndent = baseIndent + 1
	default: // default to same_line
		if f.config.Spacing.BeforeBlockStart {
			f.writeSpace()
		}
	}
	if f.config.Blocks.BraceStyle == "gnu" {
		f.writeIndentLevel(braceIndent)
		f.writeString("{")
	} else {
		f.writeString("{")
	}
	f.writeNewline()

	switch f.config.Blocks.BraceStyle {
	case "gnu":
		if f.config.Blocks.IndentStatements {
			statementIndent = baseIndent + 2
		} else {
			statementIndent = baseIndent + 1
		}
	default:
		if f.config.Blocks.IndentStatements {
			statementIndent = baseIndent + 1
		}
	}
	f.indentLevel = statementIndent
}

// writeBlockEnd writes a block end with appropriate indentation
func (f *ASTFormatter) writeBlockEnd() {
	indentDelta := 0
	switch f.config.Blocks.BraceStyle {
	case "gnu":
		if f.config.Blocks.IndentStatements {
			indentDelta = 2
		} else {
			indentDelta = 1
		}
	default:
		if f.config.Blocks.IndentStatements {
			indentDelta = 1
		}
	}

	parentIndent := f.indentLevel - indentDelta
	if parentIndent < 0 {
		parentIndent = 0
	}
	braceIndent := parentIndent
	if f.config.Blocks.BraceStyle == "gnu" {
		braceIndent = parentIndent + 1
	}
	f.indentLevel = parentIndent
	f.ensureNewline()
	f.writeIndentLevel(braceIndent)
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
	if f.lastWasNewline || f.lineLength == 0 {
		return false
	}
	return f.lineLength+nextLen > f.config.LineBreaks.MaxLineLength
}
