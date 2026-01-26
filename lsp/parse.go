package lsp

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/parser"
)

type ParseResult struct {
	Tree        parser.IProgramContext
	Tokens      *antlr.CommonTokenStream
	Diagnostics []protocol.Diagnostic
}

type parseError struct {
	start   int
	stop    int
	line    int
	column  int
	message string
}

type syntaxErrorListener struct {
	*antlr.DefaultErrorListener
	errors []parseError
}

func newSyntaxErrorListener() *syntaxErrorListener {
	return &syntaxErrorListener{DefaultErrorListener: antlr.NewDefaultErrorListener()}
}

func (l *syntaxErrorListener) SyntaxError(_ antlr.Recognizer, offendingSymbol any, line, column int, msg string, _ antlr.RecognitionException) {
	start := -1
	stop := -1
	if token, ok := offendingSymbol.(antlr.Token); ok && token != nil {
		start = token.GetStart()
		stop = token.GetStop()
	}
	l.errors = append(l.errors, parseError{start: start, stop: stop, line: line, column: column, message: msg})
}

func ParseDocument(input string) *ParseResult {
	listener := newSyntaxErrorListener()
	inputStream := antlr.NewInputStream(input)
	lexer := parser.NewbpftraceLexer(inputStream)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(listener)

	stream := antlr.NewCommonTokenStream(lexer, 0)
	bpftraceParser := parser.NewbpftraceParser(stream)
	bpftraceParser.RemoveErrorListeners()
	bpftraceParser.AddErrorListener(listener)

	tree := bpftraceParser.Program()

	return &ParseResult{
		Tree:        tree,
		Tokens:      stream,
		Diagnostics: diagnosticsFromErrors(input, listener.errors),
	}
}

func diagnosticsFromErrors(input string, errors []parseError) []protocol.Diagnostic {
	if len(errors) == 0 {
		return nil
	}

	diagnostics := make([]protocol.Diagnostic, 0, len(errors))
	for _, parseErr := range errors {
		var startPosition protocol.Position
		var endPosition protocol.Position
		if parseErr.start >= 0 && parseErr.stop >= 0 {
			startPosition = PositionForOffset(input, parseErr.start)
			endPosition = PositionForOffset(input, parseErr.stop+1)
		} else {
			startPosition = PositionForLineColumn(input, parseErr.line, parseErr.column)
			endPosition = PositionForLineColumn(input, parseErr.line, parseErr.column+1)
		}

		severity := protocol.DiagnosticSeverityError
		source := "btfmt"

		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: startPosition,
				End:   endPosition,
			},
			Severity: &severity,
			Source:   &source,
			Message:  fmt.Sprintf("line %d:%d: %s", parseErr.line, parseErr.column, parseErr.message),
		})
	}

	return diagnostics
}
