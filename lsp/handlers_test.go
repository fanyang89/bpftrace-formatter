package lsp

import (
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/formatter"
)

func setupTestState(t *testing.T) string {
	t.Helper()

	shutdownRequested.Store(false)
	configSupported.Store(false)

	initSettingsMu.Lock()
	initSettings = nil
	initSettingsMu.Unlock()

	configResolver = NewConfigResolver()
	// Force default config without searching workspace/home.
	configResolver.SetSettings(map[string]any{
		"btfmt": map[string]any{
			"configPath": filepath.Join(t.TempDir(), "does-not-exist.json"),
		},
	})

	documentStore = NewDocumentStore(configResolver)

	path := filepath.Join(t.TempDir(), "test.bt")
	u := url.URL{Scheme: "file", Path: filepath.ToSlash(path)}
	return u.String()
}

func TestInitialize_ReturnsServerCapabilities(t *testing.T) {
	setupTestState(t)

	rootURI := protocol.DocumentUri("file:///tmp")
	params := &protocol.InitializeParams{
		Capabilities: protocol.ClientCapabilities{},
		RootURI:      &rootURI,
		InitializationOptions: map[string]any{
			"btfmt": map[string]any{"indent": map[string]any{"size": 2}},
		},
	}

	resultAny, err := initialize(nil, params)
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}
	result, ok := resultAny.(protocol.InitializeResult)
	if !ok {
		t.Fatalf("initialize: unexpected result type %T", resultAny)
	}

	if value, ok := result.Capabilities.DocumentFormattingProvider.(bool); !ok || !value {
		t.Fatalf("initialize: expected formatting provider")
	}
	if value, ok := result.Capabilities.HoverProvider.(bool); !ok || !value {
		t.Fatalf("initialize: expected hover provider")
	}
	if value, ok := result.Capabilities.DocumentSymbolProvider.(bool); !ok || !value {
		t.Fatalf("initialize: expected document symbol provider")
	}
	if result.ServerInfo == nil || result.ServerInfo.Name == "" {
		t.Fatalf("initialize: expected server info")
	}
}

func TestDidOpen_PopulatesDiagnosticsForInvalidDocument(t *testing.T) {
	uri := setupTestState(t)

	invalid := "kprobe:sys_clone { @x = count( }\n"
	openParams := &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        protocol.DocumentUri(uri),
			LanguageID: "bpftrace",
			Version:    protocol.Integer(1),
			Text:       invalid,
		},
	}

	if err := didOpen(nil, openParams); err != nil {
		t.Fatalf("didOpen: %v", err)
	}

	doc, ok := documentStore.Get(uri)
	if !ok || doc == nil {
		t.Fatalf("expected document in store")
	}
	if len(doc.Diagnostics) == 0 {
		t.Fatalf("expected diagnostics for invalid document")
	}
}

func TestDidChangeConfiguration_RefreshesDocumentConfig(t *testing.T) {
	uri := setupTestState(t)

	input := "kprobe:sys_clone { @x = count(); }\n"
	if err := didOpen(nil, &protocol.DidOpenTextDocumentParams{TextDocument: protocol.TextDocumentItem{URI: protocol.DocumentUri(uri), LanguageID: "bpftrace", Version: protocol.Integer(1), Text: input}}); err != nil {
		t.Fatalf("didOpen: %v", err)
	}

	doc, ok := documentStore.Get(uri)
	if !ok || doc == nil {
		t.Fatalf("expected document in store")
	}
	if doc.Config == nil || doc.Config.Indent.Size != 4 {
		t.Fatalf("expected default indent size 4")
	}

	settings := map[string]any{
		"btfmt": map[string]any{
			"configPath": filepath.Join(t.TempDir(), "missing.json"),
			"indent":     map[string]any{"size": 2},
		},
	}
	if err := didChangeConfiguration(nil, &protocol.DidChangeConfigurationParams{Settings: settings}); err != nil {
		t.Fatalf("didChangeConfiguration: %v", err)
	}

	doc, ok = documentStore.Get(uri)
	if !ok || doc == nil {
		t.Fatalf("expected document in store")
	}
	if doc.Config == nil || doc.Config.Indent.Size != 2 {
		t.Fatalf("expected updated indent size 2")
	}
}

func TestDidChange_UpdatesDocumentTextAndDiagnostics(t *testing.T) {
	valid := "kprobe:sys_clone { @x = count(); }\n"

	cases := []struct {
		name   string
		change any
	}{
		{
			name:   "whole",
			change: protocol.TextDocumentContentChangeEventWhole{Text: "kprobe:sys_clone { @x = count( }\n"},
		},
		{
			name:   "whole_ptr",
			change: &protocol.TextDocumentContentChangeEventWhole{Text: "kprobe:sys_clone { @x = count( }\n"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			uri := setupTestState(t)
			if err := didOpen(nil, &protocol.DidOpenTextDocumentParams{TextDocument: protocol.TextDocumentItem{URI: protocol.DocumentUri(uri), LanguageID: "bpftrace", Version: protocol.Integer(1), Text: valid}}); err != nil {
				t.Fatalf("didOpen: %v", err)
			}

			params := &protocol.DidChangeTextDocumentParams{
				TextDocument: protocol.VersionedTextDocumentIdentifier{
					TextDocumentIdentifier: protocol.TextDocumentIdentifier{URI: protocol.DocumentUri(uri)},
					Version:                protocol.Integer(2),
				},
				ContentChanges: []any{
					tc.change,
				},
			}
			if err := didChange(nil, params); err != nil {
				t.Fatalf("didChange: %v", err)
			}

			doc, ok := documentStore.Get(uri)
			if !ok || doc == nil {
				t.Fatalf("expected document in store")
			}
			if doc.Text == valid {
				t.Fatalf("expected document text to change")
			}
			if len(doc.Diagnostics) == 0 {
				t.Fatalf("expected diagnostics after invalid change")
			}
		})
	}
}

func TestDidFormat_MatchesFormatterOutput(t *testing.T) {
	uri := setupTestState(t)

	input := "kprobe:sys_clone{@x[pid]=count();}\n"
	if err := didOpen(nil, &protocol.DidOpenTextDocumentParams{TextDocument: protocol.TextDocumentItem{URI: protocol.DocumentUri(uri), LanguageID: "bpftrace", Version: protocol.Integer(1), Text: input}}); err != nil {
		t.Fatalf("didOpen: %v", err)
	}

	edits, err := didFormat(nil, &protocol.DocumentFormattingParams{TextDocument: protocol.TextDocumentIdentifier{URI: protocol.DocumentUri(uri)}})
	if err != nil {
		t.Fatalf("didFormat: %v", err)
	}
	if len(edits) != 1 {
		t.Fatalf("didFormat: expected 1 edit, got %d", len(edits))
	}

	doc, ok := documentStore.Get(uri)
	if !ok || doc == nil {
		t.Fatalf("expected document in store")
	}
	want, err := formatter.NewASTFormatter(doc.Config).Format(input)
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if !strings.HasSuffix(want, "\n") {
		want += "\n"
	}
	if edits[0].NewText != want {
		t.Fatalf("didFormat: formatted text mismatch")
	}
}

func TestDidHover_ReturnsContents(t *testing.T) {
	uri := setupTestState(t)

	input := "kprobe:sys_clone { @x = count(); }\n"
	if err := didOpen(nil, &protocol.DidOpenTextDocumentParams{TextDocument: protocol.TextDocumentItem{URI: protocol.DocumentUri(uri), LanguageID: "bpftrace", Version: protocol.Integer(1), Text: input}}); err != nil {
		t.Fatalf("didOpen: %v", err)
	}

	offset := strings.Index(input, "sys_clone")
	if offset < 0 {
		t.Fatalf("expected probe marker in input")
	}

	hover, err := didHover(nil, &protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: protocol.DocumentUri(uri)},
			Position:     PositionForOffset(input, offset+2),
		},
	})
	if err != nil {
		t.Fatalf("didHover: %v", err)
	}
	if hover == nil {
		t.Fatalf("expected hover")
	}
	content, ok := hover.Contents.(protocol.MarkupContent)
	if !ok || strings.TrimSpace(content.Value) == "" {
		t.Fatalf("expected hover contents")
	}
}

func TestDidDocumentSymbol_ReturnsSymbols(t *testing.T) {
	uri := setupTestState(t)

	input := "kprobe:sys_clone { @x = count(); }\n"
	if err := didOpen(nil, &protocol.DidOpenTextDocumentParams{TextDocument: protocol.TextDocumentItem{URI: protocol.DocumentUri(uri), LanguageID: "bpftrace", Version: protocol.Integer(1), Text: input}}); err != nil {
		t.Fatalf("didOpen: %v", err)
	}

	resultAny, err := didDocumentSymbol(nil, &protocol.DocumentSymbolParams{TextDocument: protocol.TextDocumentIdentifier{URI: protocol.DocumentUri(uri)}})
	if err != nil {
		t.Fatalf("didDocumentSymbol: %v", err)
	}
	symbols, ok := resultAny.([]protocol.DocumentSymbol)
	if !ok {
		t.Fatalf("didDocumentSymbol: unexpected result type %T", resultAny)
	}
	if len(symbols) == 0 {
		t.Fatalf("expected at least one document symbol")
	}
	if symbols[0].Name == "" {
		t.Fatalf("expected non-empty symbol name")
	}
}
