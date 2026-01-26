package lsp

import (
	"strings"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestDocumentSymbols_NilDoc(t *testing.T) {
	symbols := DocumentSymbols(nil)
	if len(symbols) != 0 {
		t.Fatalf("expected no symbols, got %d", len(symbols))
	}
}

func TestDocumentSymbols_IncludesConfigProbeMacro(t *testing.T) {
	input := strings.Join([]string{
		"config = { foo = bar }",
		"macro demo(x) { @x = count(); }",
		"kprobe:sys_clone { @x = count(); }",
		"",
	}, "\n")

	doc := &Document{Text: input, ParseResult: ParseDocument(input)}
	symbols := DocumentSymbols(doc)
	if len(symbols) != 3 {
		t.Fatalf("expected 3 symbols, got %d", len(symbols))
	}

	if symbols[0].Name != "config" {
		t.Fatalf("config symbol name = %q, want %q", symbols[0].Name, "config")
	}
	if symbols[0].Kind != protocol.SymbolKindModule {
		t.Fatalf("config symbol kind = %d, want %d", symbols[0].Kind, protocol.SymbolKindModule)
	}

	if symbols[1].Name != "kprobe:sys_clone" {
		t.Fatalf("probe symbol name = %q, want %q", symbols[1].Name, "kprobe:sys_clone")
	}
	if symbols[1].Kind != protocol.SymbolKindEvent {
		t.Fatalf("probe symbol kind = %d, want %d", symbols[1].Kind, protocol.SymbolKindEvent)
	}

	if symbols[2].Name != "macro demo(...)" {
		t.Fatalf("macro symbol name = %q, want %q", symbols[2].Name, "macro demo(...)")
	}
	if symbols[2].Kind != protocol.SymbolKindFunction {
		t.Fatalf("macro symbol kind = %d, want %d", symbols[2].Kind, protocol.SymbolKindFunction)
	}
}
