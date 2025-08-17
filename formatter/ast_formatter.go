package formatter

import (
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
}

// NewASTFormatter creates a new AST-based formatter
func NewASTFormatter(cfg *config.Config) *ASTFormatter {
	return &ASTFormatter{
		config:         cfg,
		indentLevel:    0,
		lastWasNewline: true,
		needIndent:     false,
	}
}

// Format formats the given bpftrace script using AST
func (f *ASTFormatter) Format(input string) (string, error) {
	// Reset formatter state
	f.output.Reset()
	f.indentLevel = 0
	f.lastWasNewline = true
	f.needIndent = false

	// Create ANTLR input stream
	inputStream := antlr.NewInputStream(input)
	lexer := parser.NewbpftraceLexer(inputStream)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	bpftraceParser := parser.NewbpftraceParser(tokenStream)

	// Parse the program
	tree := bpftraceParser.Program()

	// Create and use the visitor
	visitor := NewASTVisitor(f)
	visitor.Visit(tree)

	return strings.TrimSpace(f.output.String()), nil
}

// writeString writes a string to the output
func (f *ASTFormatter) writeString(s string) {
	if f.needIndent && !f.isWhitespace(s) {
		f.writeIndent()
		f.needIndent = false
	}
	f.output.WriteString(s)
	f.lastWasNewline = strings.HasSuffix(s, "\n")
	if f.lastWasNewline {
		f.needIndent = true
	}
}

// writeIndent writes the current indentation
func (f *ASTFormatter) writeIndent() {
	if f.config.Indent.UseSpaces {
		f.output.WriteString(strings.Repeat(" ", f.indentLevel*f.config.Indent.Size))
	} else {
		f.output.WriteString(strings.Repeat("\t", f.indentLevel))
	}
}

// writeNewline writes a newline character
func (f *ASTFormatter) writeNewline() {
	f.output.WriteString("\n")
	f.lastWasNewline = true
	f.needIndent = true
}

// writeSpace writes a space if spacing is enabled
func (f *ASTFormatter) writeSpace() {
	f.output.WriteString(" ")
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
	if f.config.Spacing.BeforeBlockStart {
		f.writeSpace()
	}
	f.writeString("{")
	f.writeNewline()
	f.increaseIndent()
}

// writeBlockEnd writes a block end with appropriate indentation
func (f *ASTFormatter) writeBlockEnd() {
	f.decreaseIndent()
	f.ensureNewline()
	f.writeString("}")
}
