package lsp

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// HandlerContext wraps the server and its state for request handling.
type HandlerContext struct {
	server            *Server
	shutdownRequested atomic.Bool
	configSupported   atomic.Bool
	snippetSupported  atomic.Bool
	initSettingsMu    sync.Mutex
	initSettings      map[string]any
}

func (s *Server) newHandler() protocol.Handler {
	ctx := &HandlerContext{
		server: s,
	}

	return protocol.Handler{
		Initialize:                      ctx.initialize,
		Initialized:                     ctx.initialized,
		Shutdown:                        ctx.shutdown,
		Exit:                            ctx.exit,
		WorkspaceDidChangeConfiguration: ctx.didChangeConfiguration,
		TextDocumentDidOpen:             ctx.didOpen,
		TextDocumentDidChange:           ctx.didChange,
		TextDocumentDidClose:            ctx.didClose,
		TextDocumentFormatting:          ctx.didFormat,
		TextDocumentHover:               ctx.didHover,
		TextDocumentDocumentHighlight:   ctx.didDocumentHighlight,
		TextDocumentDocumentSymbol:      ctx.didDocumentSymbol,
		TextDocumentCompletion:          ctx.didCompletion,
		TextDocumentDefinition:          ctx.didDefinition,
		TextDocumentReferences:          ctx.didReferences,
		TextDocumentRename:              ctx.didRename,
		TextDocumentPrepareRename:       ctx.didPrepareRename,
	}
}

func (ctx *HandlerContext) initialize(_ *glsp.Context, params *protocol.InitializeParams) (any, error) {
	if params != nil {
		if roots := workspaceRootsFromParams(params); len(roots) > 0 {
			ctx.server.configResolver.SetWorkspaceRoots(roots)
		}
		ctx.snippetSupported.Store(clientSupportsCompletionSnippet(params))
		supported := params.Capabilities.Workspace != nil && params.Capabilities.Workspace.Configuration != nil && *params.Capabilities.Workspace.Configuration
		ctx.configSupported.Store(supported)

		if settings := settingsFromParams(params); settings != nil {
			ctx.initSettingsMu.Lock()
			ctx.initSettings = settings
			ctx.initSettingsMu.Unlock()
			if !supported {
				ctx.server.configResolver.SetSettings(normalizeSettingsMap(settings))
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

func clientSupportsCompletionSnippet(params *protocol.InitializeParams) bool {
	if params == nil || params.Capabilities.TextDocument == nil {
		return false
	}
	completion := params.Capabilities.TextDocument.Completion
	if completion == nil || completion.CompletionItem == nil || completion.CompletionItem.SnippetSupport == nil {
		return false
	}
	return *completion.CompletionItem.SnippetSupport
}

func (ctx *HandlerContext) initialized(context *glsp.Context, _ *protocol.InitializedParams) error {
	if context == nil {
		return nil
	}

	ctx.applyInitSettingsFallback()

	if ctx.configSupported.Load() {
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
			ctx.applyWorkspaceConfigurationResult(result)
		}()
		return nil
	}

	return nil
}

func (ctx *HandlerContext) applyInitSettingsFallback() {
	if settings := ctx.getInitSettings(); settings != nil {
		ctx.server.configResolver.SetSettings(normalizeSettingsMap(settings))
	}
}

func (ctx *HandlerContext) applyWorkspaceConfigurationResult(result []any) {
	if settings := settingsFromConfigurationResult(result); settings != nil {
		ctx.server.configResolver.SetSettings(settings)
		return
	}
	ctx.applyInitSettingsFallback()
}

func (ctx *HandlerContext) didOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	if params == nil {
		return nil
	}

	doc, err := ctx.server.documentStore.Open(params.TextDocument.URI, params.TextDocument.Version, params.TextDocument.Text)
	if err != nil {
		return err
	}

	version := protocol.UInteger(params.TextDocument.Version)
	ctx.publishDiagnostics(context, params.TextDocument.URI, &version, doc.Diagnostics)
	return nil
}

func (ctx *HandlerContext) didChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	if params == nil || len(params.ContentChanges) == 0 {
		return nil
	}

	uri := string(params.TextDocument.URI)
	text := ""
	hasBase := false
	if existingDoc, ok := ctx.server.documentStore.Get(uri); ok && existingDoc != nil {
		text = existingDoc.Text
		hasBase = true
	}

	updatedText, err := applyContentChangesToText(text, hasBase, params.ContentChanges)
	if err != nil {
		return err
	}

	doc, err := ctx.server.documentStore.Change(params.TextDocument.URI, params.TextDocument.Version, updatedText)
	if err != nil {
		return err
	}

	version := protocol.UInteger(params.TextDocument.Version)
	ctx.publishDiagnostics(context, params.TextDocument.URI, &version, doc.Diagnostics)
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

func (ctx *HandlerContext) didClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	if params == nil {
		return nil
	}

	ctx.server.documentStore.Close(params.TextDocument.URI)
	ctx.publishDiagnostics(context, params.TextDocument.URI, nil, []protocol.Diagnostic{})
	return nil
}

func (ctx *HandlerContext) publishDiagnostics(context *glsp.Context, uri string, version *protocol.UInteger, diagnostics []protocol.Diagnostic) {
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

func (ctx *HandlerContext) didChangeConfiguration(context *glsp.Context, params *protocol.DidChangeConfigurationParams) error {
	if params == nil {
		return nil
	}

	settings, ok := params.Settings.(map[string]any)
	if !ok {
		return nil
	}

	ctx.server.configResolver.SetSettings(normalizeSettingsMap(settings))
	if err := ctx.server.documentStore.RefreshConfigs(); err != nil {
		return err
	}

	for _, snap := range ctx.server.documentStore.AllDocs() {
		version := protocol.UInteger(snap.Version)
		ctx.publishDiagnostics(context, snap.URI, &version, snap.Diagnostics)
	}
	return nil
}

func (ctx *HandlerContext) didFormat(_ *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	if params == nil {
		return []protocol.TextEdit{}, nil
	}
	return ctx.formatDocument(params.TextDocument.URI)
}

func (ctx *HandlerContext) didHover(_ *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	if params == nil {
		return nil, nil
	}

	doc, ok := ctx.server.documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return nil, nil
	}

	return HoverForPosition(doc, params.Position), nil
}

func (ctx *HandlerContext) didDocumentHighlight(_ *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	if params == nil {
		return []protocol.DocumentHighlight{}, nil
	}

	doc, ok := ctx.server.documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return []protocol.DocumentHighlight{}, nil
	}

	return documentHighlightsForPosition(doc, params.Position), nil
}

func (ctx *HandlerContext) didDocumentSymbol(_ *glsp.Context, params *protocol.DocumentSymbolParams) (any, error) {
	if params == nil {
		return []protocol.DocumentSymbol{}, nil
	}

	doc, ok := ctx.server.documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return []protocol.DocumentSymbol{}, nil
	}

	return DocumentSymbols(doc), nil
}

func (ctx *HandlerContext) didCompletion(_ *glsp.Context, params *protocol.CompletionParams) (any, error) {
	if params == nil {
		return []protocol.CompletionItem{}, nil
	}

	doc, ok := ctx.server.documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return ctx.defaultCompletions(), nil
	}

	return ctx.CompletionForPosition(doc, params.Position), nil
}

func (ctx *HandlerContext) didDefinition(_ *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	if params == nil {
		return []protocol.Location{}, nil
	}

	doc, ok := ctx.server.documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return []protocol.Location{}, nil
	}

	location, ok := definitionLocationForPosition(doc, params.Position)
	if !ok {
		return []protocol.Location{}, nil
	}

	return []protocol.Location{location}, nil
}

func (ctx *HandlerContext) didReferences(_ *glsp.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	if params == nil {
		return []protocol.Location{}, nil
	}

	doc, ok := ctx.server.documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return []protocol.Location{}, nil
	}

	return referencesForPosition(doc, params.Position, params.Context.IncludeDeclaration), nil
}

func (ctx *HandlerContext) didRename(_ *glsp.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	if params == nil {
		return nil, nil
	}

	doc, ok := ctx.server.documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return nil, nil
	}

	return renameWorkspaceEditForPosition(doc, params.Position, params.NewName)
}

func (ctx *HandlerContext) didPrepareRename(_ *glsp.Context, params *protocol.PrepareRenameParams) (any, error) {
	if params == nil {
		return nil, nil
	}

	doc, ok := ctx.server.documentStore.Get(params.TextDocument.URI)
	if !ok || doc == nil {
		return nil, nil
	}

	rangeValue, ok := prepareRenameRangeForPosition(doc, params.Position)
	if !ok {
		return nil, nil
	}

	return *rangeValue, nil
}

func (ctx *HandlerContext) getInitSettings() map[string]any {
	ctx.initSettingsMu.Lock()
	defer ctx.initSettingsMu.Unlock()
	return ctx.initSettings
}

func (ctx *HandlerContext) shutdown(_ *glsp.Context) error {
	ctx.shutdownRequested.Store(true)
	return nil
}

func (ctx *HandlerContext) exit(_ *glsp.Context) error {
	if ctx.shutdownRequested.Load() {
		os.Exit(0)
	}
	os.Exit(1)
	return nil
}

func serverCapabilities() protocol.ServerCapabilities {
	syncKind := protocol.TextDocumentSyncKindIncremental
	openClose := true
	prepareRename := true

	// Completion trigger characters
	triggerChars := []string{"@", "$", ":", "."}

	return protocol.ServerCapabilities{
		TextDocumentSync: &protocol.TextDocumentSyncOptions{
			OpenClose: &openClose,
			Change:    &syncKind,
		},
		DocumentFormattingProvider: true,
		HoverProvider:              true,
		DocumentHighlightProvider:  true,
		DocumentSymbolProvider:     true,
		DefinitionProvider:         true,
		ReferencesProvider:         true,
		RenameProvider:             &protocol.RenameOptions{PrepareProvider: &prepareRename},
		CompletionProvider: &protocol.CompletionOptions{
			TriggerCharacters: triggerChars,
		},
	}
}
