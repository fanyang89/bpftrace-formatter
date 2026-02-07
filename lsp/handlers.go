package lsp

import (
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

	if configSupported.Load() {
		section := "btfmt"
		params := protocol.ConfigurationParams{
			Items: []protocol.ConfigurationItem{{Section: &section}},
		}
		var result []any
		context.Call(protocol.ServerWorkspaceConfiguration, params, &result)
		if settings := settingsFromConfigurationResult(result); settings != nil {
			configResolver.SetSettings(settings)
			return nil
		}
	}

	if settings := getInitSettings(); settings != nil {
		configResolver.SetSettings(normalizeSettingsMap(settings))
	}

	return nil
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
	text := ""
	change := params.ContentChanges[0]
	switch value := change.(type) {
	case protocol.TextDocumentContentChangeEvent:
		text = value.Text
	case *protocol.TextDocumentContentChangeEvent:
		text = value.Text
	case protocol.TextDocumentContentChangeEventWhole:
		text = value.Text
	case *protocol.TextDocumentContentChangeEventWhole:
		text = value.Text
	default:
		return nil
	}
	doc, err := documentStore.Change(params.TextDocument.URI, params.TextDocument.Version, text)
	if err != nil {
		return err
	}

	version := protocol.UInteger(params.TextDocument.Version)
	publishDiagnostics(context, params.TextDocument.URI, &version, doc.Diagnostics)
	return nil
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
	syncKind := protocol.TextDocumentSyncKindFull
	openClose := true

	return protocol.ServerCapabilities{
		TextDocumentSync: &protocol.TextDocumentSyncOptions{
			OpenClose: &openClose,
			Change:    &syncKind,
		},
		DocumentFormattingProvider: true,
		HoverProvider:              true,
		DocumentSymbolProvider:     true,
	}
}
