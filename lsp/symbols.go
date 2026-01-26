package lsp

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/parser"
)

// DocumentSymbols returns document symbols for LSP.
func DocumentSymbols(doc *Document) []protocol.DocumentSymbol {
	return documentSymbols(doc)
}

func documentSymbols(doc *Document) []protocol.DocumentSymbol {
	if doc == nil || doc.ParseResult == nil || doc.ParseResult.Tree == nil {
		return []protocol.DocumentSymbol{}
	}

	program := doc.ParseResult.Tree
	text := doc.Text
	items := make([]protocol.DocumentSymbol, 0)

	if config := program.Config_preamble(); config != nil {
		if section := config.Config_section(); section != nil {
			rangeFull := rangeForContext(text, section)
			selectionRange := rangeFull
			if token := section.CONFIG(); token != nil {
				if tokenRange, ok := rangeForTerminal(text, token); ok {
					selectionRange = tokenRange
				}
			}
			items = append(items, protocol.DocumentSymbol{
				Name:           "config",
				Kind:           protocol.SymbolKindModule,
				Range:          rangeFull,
				SelectionRange: selectionRange,
			})
		}
	}

	content := program.Content()
	if content == nil {
		return items
	}

	for _, probe := range content.AllProbe() {
		name := "probe"
		if probeList := probe.Probe_list(); probeList != nil {
			name = singleLineText(probeList.GetText())
		}
		rangeFull := rangeForContext(text, probe)
		selectionRange := rangeFull
		if probeList := probe.Probe_list(); probeList != nil {
			selectionRange = rangeForContext(text, probeList)
		}
		items = append(items, protocol.DocumentSymbol{
			Name:           name,
			Kind:           protocol.SymbolKindEvent,
			Range:          rangeFull,
			SelectionRange: selectionRange,
		})
	}

	for _, macro := range content.AllMacro_definition() {
		rangeFull := rangeForContext(text, macro)
		selectionRange := rangeFull
		if identifier := macro.IDENTIFIER(); identifier != nil {
			if tokenRange, ok := rangeForTerminal(text, identifier); ok {
				selectionRange = tokenRange
			}
		}

		items = append(items, protocol.DocumentSymbol{
			Name:           macroName(macro),
			Kind:           protocol.SymbolKindFunction,
			Range:          rangeFull,
			SelectionRange: selectionRange,
		})
	}

	return items
}

func rangeForContext(text string, ctx antlr.ParserRuleContext) protocol.Range {
	if ctx == nil {
		return protocol.Range{}
	}
	startToken := ctx.GetStart()
	stopToken := ctx.GetStop()
	return rangeForTokens(text, startToken, stopToken)
}

func rangeForTerminal(text string, node antlr.TerminalNode) (protocol.Range, bool) {
	if node == nil {
		return protocol.Range{}, false
	}
	return rangeForToken(text, node.GetSymbol())
}

func rangeForToken(text string, token antlr.Token) (protocol.Range, bool) {
	if token == nil {
		return protocol.Range{}, false
	}
	start := token.GetStart()
	stop := token.GetStop()
	if start < 0 || stop < 0 {
		return protocol.Range{}, false
	}
	return rangeForOffsets(text, start, stop+1), true
}

func rangeForTokens(text string, startToken antlr.Token, stopToken antlr.Token) protocol.Range {
	start := 0
	if startToken != nil && startToken.GetStart() >= 0 {
		start = startToken.GetStart()
	}
	end := start
	if stopToken != nil && stopToken.GetStop() >= 0 {
		end = stopToken.GetStop() + 1
	}
	if end < start {
		end = start
	}
	return rangeForOffsets(text, start, end)
}

func rangeForOffsets(text string, start int, end int) protocol.Range {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	return protocol.Range{
		Start: PositionForOffset(text, start),
		End:   PositionForOffset(text, end),
	}
}

func singleLineText(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

func macroName(macro parser.IMacro_definitionContext) string {
	name := "macro"
	identifier := ""
	if node := macro.IDENTIFIER(); node != nil {
		identifier = strings.TrimSpace(node.GetText())
	}
	if identifier != "" {
		name = "macro " + identifier
	}
	if macro.Macro_params() != nil {
		if identifier == "" {
			name = "macro (...)"
		} else {
			name += "(...)"
		}
	}
	return name
}
