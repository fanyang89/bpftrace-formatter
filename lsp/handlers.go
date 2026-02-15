package lsp

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	shutdownRequested atomic.Bool
	configResolver    = NewConfigResolver()
	documentStore     = NewDocumentStore(configResolver)
	configSupported   atomic.Bool
	initSettingsMu    sync.Mutex
	initSettings      map[string]any
)

func newHandler() protocol.Handler {
	return protocol.Handler{
		Initialize:                      initialize,
		Initialized:                     initialized,
		Shutdown:                        shutdown,
		Exit:                            exit,
		WorkspaceDidChangeConfiguration: didChangeConfiguration,
		TextDocumentDidOpen:             didOpen,
		TextDocumentDidChange:           didChange,
		TextDocumentDidClose:            didClose,
		TextDocumentFormatting:          didFormat,
		TextDocumentHover:               didHover,
		TextDocumentDocumentSymbol:      didDocumentSymbol,
		TextDocumentCompletion:          didCompletion,
	}
}

func initialize(_ *glsp.Context, params *protocol.InitializeParams) (any, error) {
	if params != nil {
		if roots := workspaceRootsFromParams(params); len(roots) > 0 {
			configResolver.SetWorkspaceRoots(roots)
		}
		supported := params.Capabilities.Workspace != nil && params.Capabilities.Workspace.Configuration != nil && *params.Capabilities.Workspace.Configuration
		configSupported.Store(supported)

		if settings := settingsFromParams(params); settings != nil {
			initSettingsMu.Lock()
			initSettings = settings
			initSettingsMu.Unlock()
			if !supported {
				configResolver.SetSettings(normalizeSettingsMap(settings))
			}
		}
	}

	capabilities := serverCapabilities()
	version := serverVersion

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    serverName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, _ *protocol.InitializedParams) error {
	if context == nil {
		return nil
	}

	applyInitSettingsFallback()

	if configSupported.Load() {
		// Fetch configuration asynchronously to avoid blocking the message loop.
		// A synchronous call here would deadlock because the response arrives
		// on the same message loop that is waiting for this handler to return.
		go func() {
			section := "btfmt"
			params := protocol.ConfigurationParams{
				Items: []protocol.ConfigurationItem{{Section: &section}},
			}
			var result []any
			context.Call(protocol.ServerWorkspaceConfiguration, params, &result)
			applyWorkspaceConfigurationResult(result)
		}()
		return nil
	}

	return nil
}

func applyInitSettingsFallback() {
	if settings := getInitSettings(); settings != nil {
		configResolver.SetSettings(normalizeSettingsMap(settings))
	}
}

func applyWorkspaceConfigurationResult(result []any) {
	if settings := settingsFromConfigurationResult(result); settings != nil {
		configResolver.SetSettings(settings)
		return
	}
	applyInitSettingsFallback()
}

func didOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	if params == nil {
		return nil
	}

	doc, err := documentStore.Open(params.TextDocument.URI, params.TextDocument.Version, params.TextDocument.Text)
	if err != nil {
		return err
	}

	version := protocol.UInteger(params.TextDocument.Version)
	publishDiagnostics(context, params.TextDocument.URI, &version, doc.Diagnostics)
	return nil
}

func didChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	if params == nil || len(params.ContentChanges) == 0 {
		return nil
	}

	uri := string(params.TextDocument.URI)
	text := ""
	hasBase := false
	if existingDoc, ok := documentStore.Get(uri); ok && existingDoc != nil {
		text = existingDoc.Text
		hasBase = true
	}

	updatedText, err := applyContentChangesToText(text, hasBase, params.ContentChanges)
	if err != nil {
		return err
	}

	doc, err := documentStore.Change(params.TextDocument.URI, params.TextDocument.Version, updatedText)
	if err != nil {
		return err
	}

	version := protocol.UInteger(params.TextDocument.Version)
	publishDiagnostics(context, params.TextDocument.URI, &version, doc.Diagnostics)
	return nil
}

func applyContentChangesToText(text string, hasBase bool, changes []any) (string, error) {
	currentText := text
	baseAvailable := hasBase

	for _, change := range changes {
		switch value := change.(type) {
		case protocol.TextDocumentContentChangeEvent:
			var err error
			currentText, baseAvailable, err = applyTextDocumentContentChange(currentText, baseAvailable, value.Range, value.Text)
			if err != nil {
				return "", err
			}
		case *protocol.TextDocumentContentChangeEvent:
			if value == nil {
				continue
			}
			var err error
			currentText, baseAvailable, err = applyTextDocumentContentChange(currentText, baseAvailable, value.Range, value.Text)
			if err != nil {
				return "", err
			}
		case protocol.TextDocumentContentChangeEventWhole:
			currentText = value.Text
			baseAvailable = true
		case *protocol.TextDocumentContentChangeEventWhole:
			if value == nil {
				continue
			}
			currentText = value.Text
			baseAvailable = true
		default:
			return "", fmt.Errorf("unsupported text document change type %T", change)
		}
	}

	return currentText, nil
}

func applyTextDocumentContentChange(currentText string, hasBase bool, changeRange *protocol.Range, newText string) (string, bool, error) {
	if changeRange == nil {
		return newText, true, nil
	}
	if !hasBase {
		return "", false, fmt.Errorf("incremental change received before document content is available")
	}

	updatedText, err := applyIncrementalChange(currentText, *changeRange, newText)
	if err != nil {
		return "", hasBase, err
	}

	return updatedText, true, nil
}

func applyIncrementalChange(text string, changeRange protocol.Range, newText string) (string, error) {
	start := offsetForPosition(text, changeRange.Start)
	end := offsetForPosition(text, changeRange.End)
	if start > end {
		return "", fmt.Errorf("invalid incremental change range: start=%d end=%d", start, end)
	}
	if start < 0 || end < 0 || start > len(text) || end > len(text) {
		return "", fmt.Errorf("incremental change range out of bounds: start=%d end=%d len=%d", start, end, len(text))
	}
	return text[:start] + newText + text[end:], nil
}

func didClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	if params == nil {
		return nil
	}

	documentStore.Close(params.TextDocument.URI)
	publishDiagnostics(context, params.TextDocument.URI, nil, []protocol.Diagnostic{})
	return nil
}

func publishDiagnostics(context *glsp.Context, uri string, version *protocol.UInteger, diagnostics []protocol.Diagnostic) {
	if context == nil {
		return
	}

	params := protocol.PublishDiagnosticsParams{
		URI:         uri,
		Version:     version,
		Diagnostics: diagnostics,
	}
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, params)
}

func didChangeConfiguration(context *glsp.Context, params *protocol.DidChangeConfigurationParams) error {
	if params == nil {
		return nil
	}

	settings, ok := params.Settings.(map[string]any)
	if !ok {
		return nil
	}

	configResolver.SetSettings(normalizeSettingsMap(settings))
	if err := documentStore.RefreshConfigs(); err != nil {
		return err
	}

	for _, snap := range documentStore.AllDocs() {
		version := protocol.UInteger(snap.Version)
		publishDiagnostics(context, snap.URI, &version, snap.Diagnostics)
	}
	return nil
}

func didFormat(_ *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	if params == nil {
		return []protocol.TextEdit{}, nil
	}
	return formatDocument(params.TextDocument.URI)
}

func didHover(_ *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	if params == nil {
		return nil, nil
	}

	doc, ok := documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return nil, nil
	}

	return HoverForPosition(doc, params.Position), nil
}

func didDocumentSymbol(_ *glsp.Context, params *protocol.DocumentSymbolParams) (any, error) {
	if params == nil {
		return []protocol.DocumentSymbol{}, nil
	}

	doc, ok := documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return []protocol.DocumentSymbol{}, nil
	}

	return DocumentSymbols(doc), nil
}

func didCompletion(_ *glsp.Context, params *protocol.CompletionParams) (any, error) {
	if params == nil {
		return []protocol.CompletionItem{}, nil
	}

	doc, ok := documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return defaultCompletions(), nil
	}

	return CompletionForPosition(doc, params.Position), nil
}

func getInitSettings() map[string]any {
	initSettingsMu.Lock()
	defer initSettingsMu.Unlock()
	return initSettings
}

func shutdown(_ *glsp.Context) error {
	shutdownRequested.Store(true)
	return nil
}

func exit(_ *glsp.Context) error {
	if shutdownRequested.Load() {
		os.Exit(0)
	}
	os.Exit(1)
	return nil
}

func serverCapabilities() protocol.ServerCapabilities {
	syncKind := protocol.TextDocumentSyncKindIncremental
	openClose := true

	// Completion trigger characters
	triggerChars := []string{"@", "$", ":", "."}

	return protocol.ServerCapabilities{
		TextDocumentSync: &protocol.TextDocumentSyncOptions{
			OpenClose: &openClose,
			Change:    &syncKind,
		},
		DocumentFormattingProvider: true,
		HoverProvider:              true,
		DocumentSymbolProvider:     true,
		CompletionProvider: &protocol.CompletionOptions{
			TriggerCharacters: triggerChars,
		},
	}
}
