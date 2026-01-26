package lsp

import (
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/config"
	"github.com/fanyang89/bpftrace-formatter/formatter"
)

func formatDocument(uri string) ([]protocol.TextEdit, error) {
	doc, ok := documentStore.Get(uri)
	if !ok || doc == nil {
		return []protocol.TextEdit{}, nil
	}

	cfg := doc.Config
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	formatted, err := formatter.NewASTFormatter(cfg).Format(doc.Text)
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
