package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TokenType represents the type of token
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenProbe
	TokenPredicate
	TokenBlockStart
	TokenBlockEnd
	TokenStatement
	TokenComment
	TokenNewline
	TokenWhitespace
)

// Token represents a lexical token
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// Formatter handles bpftrace script formatting
type Formatter struct {
	indentSize    int
	useSpaces     bool
	probeSpacing  int
	commentIndent int
}

// NewFormatter creates a new formatter with default settings
func NewFormatter() *Formatter {
	return &Formatter{
		indentSize:    4,
		useSpaces:     true,
		probeSpacing:  2,
		commentIndent: 2,
	}
}

// Format formats a bpftrace script
func (f *Formatter) Format(input string) string {
	// Pre-process: handle shebang and comments
	lines := strings.Split(input, "\n")
	var processedLines []string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			processedLines = append(processedLines, "")
			continue
		}

		// Add blank line after shebang
		if strings.HasPrefix(line, "#!") && i < len(lines)-1 {
			processedLines = append(processedLines, line)
			processedLines = append(processedLines, "")
			continue
		}

		processedLines = append(processedLines, line)
	}

	processedInput := strings.Join(processedLines, "\n")
	tokens := f.tokenize(processedInput)
	return f.formatTokens(tokens)
}

// tokenize breaks the input into tokens
func (f *Formatter) tokenize(input string) []Token {
	var tokens []Token
	scanner := bufio.NewScanner(strings.NewReader(input))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimRight(scanner.Text(), " \t")

		if line == "" {
			tokens = append(tokens, Token{Type: TokenNewline, Value: "\n", Line: lineNum})
			continue
		}

		// Handle comments
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			tokens = append(tokens, Token{Type: TokenComment, Value: line, Line: lineNum})
			tokens = append(tokens, Token{Type: TokenNewline, Value: "\n", Line: lineNum})
			continue
		}

		// Handle shebang
		if strings.HasPrefix(line, "#!") {
			tokens = append(tokens, Token{Type: TokenComment, Value: line, Line: lineNum})
			tokens = append(tokens, Token{Type: TokenNewline, Value: "\n", Line: lineNum})
			continue
		}

		// Use regex to find probe patterns
		probePattern := `^(\s*)([a-zA-Z_][a-zA-Z0-9_]*:[^{\s]*)\s*(/[^/]+/)?\s*\{(.*)\}\s*$`
		re := regexp.MustCompile(probePattern)

		if matches := re.FindStringSubmatch(line); matches != nil {
			probeDef := strings.TrimSpace(matches[2])
			predicate := matches[3]
			content := strings.TrimSpace(matches[4])

			// Add probe token
			tokens = append(tokens, Token{Type: TokenProbe, Value: probeDef, Line: lineNum})

			// Add predicate if present
			if predicate != "" {
				tokens = append(tokens, Token{Type: TokenPredicate, Value: predicate, Line: lineNum})
			}

			// Add block and content
			tokens = append(tokens, Token{Type: TokenBlockStart, Value: "{", Line: lineNum})
			if content != "" {
				// Split content by semicolons and create separate statements
				statements := strings.Split(content, ";")
				for i, stmt := range statements {
					stmt = strings.TrimSpace(stmt)
					if stmt != "" {
						tokens = append(tokens, Token{Type: TokenStatement, Value: stmt + ";", Line: lineNum})
						// Add newline after each statement except the last one
						if i < len(statements)-1 {
							tokens = append(tokens, Token{Type: TokenNewline, Value: "\n", Line: lineNum})
						}
					}
				}
			}
			tokens = append(tokens, Token{Type: TokenBlockEnd, Value: "}", Line: lineNum})
		} else {
			// Handle standalone statements (like END blocks)
			if strings.Contains(line, "{") || strings.Contains(line, "}") {
				// Simple block handling for non-probe statements
				if strings.HasPrefix(line, "END") {
					// Special handling for END blocks
					tokens = append(tokens, Token{Type: TokenProbe, Value: "END", Line: lineNum})
					content := strings.Trim(strings.TrimPrefix(line, "END"), " {}")
					if content != "" {
						tokens = append(tokens, Token{Type: TokenBlockStart, Value: "{", Line: lineNum})
						// Split content by semicolons and create separate statements
						statements := strings.Split(content, ";")
						for i, stmt := range statements {
							stmt = strings.TrimSpace(stmt)
							if stmt != "" {
								tokens = append(tokens, Token{Type: TokenStatement, Value: stmt + ";", Line: lineNum})
								// Add newline after each statement except the last one
								if i < len(statements)-1 {
									tokens = append(tokens, Token{Type: TokenNewline, Value: "\n", Line: lineNum})
								}
							}
						}
						tokens = append(tokens, Token{Type: TokenBlockEnd, Value: "}", Line: lineNum})
					}
				} else {
					tokens = append(tokens, Token{Type: TokenStatement, Value: line, Line: lineNum})
				}
			} else {
				tokens = append(tokens, Token{Type: TokenStatement, Value: line, Line: lineNum})
			}
		}

		tokens = append(tokens, Token{Type: TokenNewline, Value: "\n", Line: lineNum})
	}

	return tokens
}

// formatTokens formats the tokens into a properly indented script
func (f *Formatter) formatTokens(tokens []Token) string {
	var result strings.Builder
	indentLevel := 0
	lastTokenWasProbe := false
	needIndent := false

	for i, token := range tokens {
		switch token.Type {
		case TokenProbe:
			// Add spacing between probes
			if lastTokenWasProbe && i > 0 {
				for j := 0; j < f.probeSpacing; j++ {
					result.WriteString("\n")
				}
			}
			result.WriteString(token.Value)
			lastTokenWasProbe = true
			needIndent = false

		case TokenPredicate:
			result.WriteString(" ")
			result.WriteString(strings.TrimSpace(token.Value))
			needIndent = false

		case TokenBlockStart:
			result.WriteString(" {\n")
			indentLevel++
			needIndent = true

		case TokenBlockEnd:
			indentLevel--
			if indentLevel >= 0 {
				result.WriteString("\n")
				result.WriteString(f.indent(indentLevel))
				result.WriteString("}")
			}
			needIndent = false

		case TokenStatement:
			if strings.TrimSpace(token.Value) != "" {
				if needIndent {
					result.WriteString(f.indent(indentLevel))
					needIndent = false
				}
				result.WriteString(strings.TrimSpace(token.Value))
			}

		case TokenComment:
			if needIndent {
				result.WriteString(f.indent(indentLevel))
				needIndent = false
			}
			result.WriteString(token.Value)

		case TokenNewline:
			if i < len(tokens)-1 && tokens[i+1].Type != TokenBlockEnd {
				result.WriteString("\n")
				if tokens[i+1].Type == TokenStatement || tokens[i+1].Type == TokenComment {
					needIndent = true
				}
			}

		case TokenEOF:
			// Do nothing
		}
	}

	return strings.TrimSpace(result.String())
}

// indent returns the appropriate indentation string
func (f *Formatter) indent(level int) string {
	if f.useSpaces {
		return strings.Repeat(" ", level*f.indentSize)
	}
	return strings.Repeat("\t", level)
}

// SetIndentSize sets the indentation size
func (f *Formatter) SetIndentSize(size int) {
	f.indentSize = size
}

// SetUseSpaces sets whether to use spaces instead of tabs
func (f *Formatter) SetUseSpaces(useSpaces bool) {
	f.useSpaces = useSpaces
}

// SetProbeSpacing sets the number of blank lines between probes
func (f *Formatter) SetProbeSpacing(spacing int) {
	f.probeSpacing = spacing
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bpftrace-formatter <file.bt>")
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	formatter := NewFormatter()
	formatted := formatter.Format(string(content))

	fmt.Println(formatted)
}
