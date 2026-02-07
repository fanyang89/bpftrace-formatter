package lsp

import (
	"fmt"
	"log"
	"strings"
	"time"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/config"
	"github.com/fanyang89/bpftrace-formatter/formatter"
)

const formatTimeout = 30 * time.Second

type formatResult struct {
	text string
	err  error
}

func formatWithTimeout(doc *Document, cfg *config.Config, timeout time.Duration) (string, error) {
	ch := make(chan formatResult, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- formatResult{err: fmt.Errorf("formatter panic: %v", r)}
			}
		}()
		start := time.Now()
		f := formatter.NewASTFormatter(cfg)
		// Reuse the parse tree from didOpen/didChange when available,
		// avoiding a redundant ANTLR parse.
		if doc.ParseResult != nil && doc.ParseResult.Tree != nil && len(doc.ParseResult.Diagnostics) == 0 {
			text := f.FormatTree(doc.ParseResult.Tree)
			log.Printf("[format] FormatTree took %s", time.Since(start))
			ch <- formatResult{text: text}
		} else {
			formatted, err := f.Format(doc.Text)
			log.Printf("[format] Format took %s", time.Since(start))
			ch <- formatResult{text: formatted, err: err}
		}
	}()
	select {
	case res := <-ch:
		return res.text, res.err
	case <-time.After(timeout):
		return "", fmt.Errorf("formatting timed out after %s", timeout)
	}
}

func formatDocument(uri string) ([]protocol.TextEdit, error) {
	doc, ok := documentStore.Get(uri)
	if !ok || doc == nil {
		return []protocol.TextEdit{}, nil
	}

	cfg := doc.Config
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	formatted, err := formatWithTimeout(doc, cfg, formatTimeout)
	if err != nil {
		return nil, err
	}

	// Ensure trailing newline for consistency with CLI
	if !strings.HasSuffix(formatted, "\n") {
		formatted += "\n"
	}

	edits := []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   EndPosition(doc.Text),
			},
			NewText: formatted,
		},
	}

	return edits, nil
}
