package lsp

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type symbolKind int

const (
	symbolKindVariable symbolKind = iota
	symbolKindMap
)

type symbolOccurrence struct {
	kind         symbolKind
	name         string
	sigil        byte
	rangeValue   protocol.Range
	token        antlr.Token
	isDefinition bool
}

func definitionLocationForPosition(doc *Document, pos protocol.Position) (protocol.Location, bool) {
	target, ok := symbolAtPosition(doc, pos)
	if !ok {
		return protocol.Location{}, false
	}

	occurrences := symbolOccurrencesByIdentity(doc, target.kind, target.name)
	definition, ok := selectDefinitionOccurrence(occurrences)
	if !ok {
		return protocol.Location{}, false
	}

	return protocol.Location{
		URI:   protocol.DocumentUri(doc.URI),
		Range: definition.rangeValue,
	}, true
}

func documentHighlightsForPosition(doc *Document, pos protocol.Position) []protocol.DocumentHighlight {
	target, ok := symbolAtPosition(doc, pos)
	if !ok {
		return []protocol.DocumentHighlight{}
	}

	occurrences := symbolOccurrencesByIdentity(doc, target.kind, target.name)
	highlights := make([]protocol.DocumentHighlight, 0, len(occurrences))
	for _, occurrence := range occurrences {
		kind := protocol.DocumentHighlightKindRead
		if occurrence.isDefinition {
			kind = protocol.DocumentHighlightKindWrite
		}
		highlights = append(highlights, protocol.DocumentHighlight{
			Range: occurrence.rangeValue,
			Kind:  &kind,
		})
	}

	return highlights
}

func prepareRenameRangeForPosition(doc *Document, pos protocol.Position) (*protocol.Range, bool) {
	target, ok := symbolAtPosition(doc, pos)
	if !ok {
		return nil, false
	}

	rangeValue := target.rangeValue
	return &rangeValue, true
}

func referencesForPosition(doc *Document, pos protocol.Position, includeDeclaration bool) []protocol.Location {
	target, ok := symbolAtPosition(doc, pos)
	if !ok {
		return []protocol.Location{}
	}

	occurrences := symbolOccurrencesByIdentity(doc, target.kind, target.name)
	if len(occurrences) == 0 {
		return []protocol.Location{}
	}

	definition, hasDefinition := selectDefinitionOccurrence(occurrences)
	locations := make([]protocol.Location, 0, len(occurrences))
	for _, occurrence := range occurrences {
		if !includeDeclaration && hasDefinition && sameTokenIndex(occurrence.token, definition.token) {
			continue
		}
		locations = append(locations, protocol.Location{
			URI:   protocol.DocumentUri(doc.URI),
			Range: occurrence.rangeValue,
		})
	}

	return locations
}

func renameWorkspaceEditForPosition(doc *Document, pos protocol.Position, newName string) (*protocol.WorkspaceEdit, error) {
	target, ok := symbolAtPosition(doc, pos)
	if !ok {
		return nil, nil
	}

	normalizedName, err := normalizeRenameIdentifier(newName)
	if err != nil {
		return nil, err
	}

	occurrences := symbolOccurrencesByIdentity(doc, target.kind, target.name)
	if len(occurrences) == 0 {
		return nil, nil
	}

	edits := make([]protocol.TextEdit, 0, len(occurrences))
	for _, occurrence := range occurrences {
		replacement := string(occurrence.sigil) + normalizedName
		edits = append(edits, protocol.TextEdit{
			Range:   occurrence.rangeValue,
			NewText: replacement,
		})
	}

	if len(edits) == 0 {
		return nil, nil
	}

	changes := map[protocol.DocumentUri][]protocol.TextEdit{
		protocol.DocumentUri(doc.URI): edits,
	}

	return &protocol.WorkspaceEdit{Changes: changes}, nil
}

func normalizeRenameIdentifier(newName string) (string, error) {
	name := strings.TrimSpace(newName)
	if name == "" {
		return "", fmt.Errorf("rename target must not be empty")
	}

	if name[0] == '$' || name[0] == '@' {
		name = name[1:]
	}

	if !isValidIdentifier(name) {
		return "", fmt.Errorf("invalid rename target %q", newName)
	}

	return name, nil
}

func symbolAtPosition(doc *Document, pos protocol.Position) (*symbolOccurrence, bool) {
	occurrences := collectSymbolOccurrences(doc)
	if len(occurrences) == 0 {
		return nil, false
	}

	best := -1
	bestSpan := 0
	for i, occurrence := range occurrences {
		if !positionInRange(pos, occurrence.rangeValue) {
			continue
		}
		span := tokenSpan(occurrence.token)
		if best < 0 || span < bestSpan {
			best = i
			bestSpan = span
		}
	}

	if best < 0 {
		return nil, false
	}

	selected := occurrences[best]
	return &selected, true
}

func symbolOccurrencesByIdentity(doc *Document, kind symbolKind, name string) []symbolOccurrence {
	occurrences := collectSymbolOccurrences(doc)
	filtered := make([]symbolOccurrence, 0, len(occurrences))
	for _, occurrence := range occurrences {
		if occurrence.kind == kind && occurrence.name == name {
			filtered = append(filtered, occurrence)
		}
	}
	return filtered
}

func collectSymbolOccurrences(doc *Document) []symbolOccurrence {
	if doc == nil || doc.ParseResult == nil || doc.ParseResult.Tree == nil {
		return nil
	}

	occurrences := make([]symbolOccurrence, 0)
	var walk func(node antlr.Tree)
	walk = func(node antlr.Tree) {
		if node == nil {
			return
		}

		if terminal, ok := node.(antlr.TerminalNode); ok {
			if occurrence, ok := symbolOccurrenceFromTerminal(doc.Text, terminal); ok {
				occurrences = append(occurrences, occurrence)
			}
		}

		for i := 0; i < node.GetChildCount(); i++ {
			walk(node.GetChild(i))
		}
	}

	walk(doc.ParseResult.Tree)
	sortSymbolOccurrences(occurrences)
	return occurrences
}

func symbolOccurrenceFromTerminal(text string, terminal antlr.TerminalNode) (symbolOccurrence, bool) {
	if terminal == nil {
		return symbolOccurrence{}, false
	}

	kind, name, sigil, ok := symbolIdentityFromText(terminal.GetText())
	if !ok {
		return symbolOccurrence{}, false
	}

	rangeValue, ok := rangeForTerminal(text, terminal)
	if !ok {
		return symbolOccurrence{}, false
	}

	token := terminal.GetSymbol()
	if token == nil {
		return symbolOccurrence{}, false
	}

	return symbolOccurrence{
		kind:         kind,
		name:         name,
		sigil:        sigil,
		rangeValue:   rangeValue,
		token:        token,
		isDefinition: tokenRepresentsAssignmentLHS(text, token, sigil),
	}, true
}

func symbolIdentityFromText(tokenText string) (symbolKind, string, byte, bool) {
	if len(tokenText) < 2 {
		return 0, "", 0, false
	}

	sigil := tokenText[0]
	if sigil != '$' && sigil != '@' {
		return 0, "", 0, false
	}

	name := tokenText[1:]
	if !isValidIdentifier(name) {
		return 0, "", 0, false
	}

	if sigil == '$' {
		return symbolKindVariable, name, sigil, true
	}
	return symbolKindMap, name, sigil, true
}

func tokenRepresentsAssignmentLHS(text string, token antlr.Token, sigil byte) bool {
	if token == nil {
		return false
	}

	start := token.GetStart()
	stop := token.GetStop()
	if start < 0 || stop < start || stop+1 > len(text) {
		return false
	}

	index := stop + 1
	index = skipASCIIWhitespaceForward(text, index)

	if sigil == '@' && index < len(text) && text[index] == '[' {
		next, ok := consumeBracketExpression(text, index)
		if !ok {
			return false
		}
		index = skipASCIIWhitespaceForward(text, next)
	}

	return hasAssignmentOperator(text, index)
}

func skipASCIIWhitespaceForward(text string, index int) int {
	for index < len(text) && isASCIIWhitespace(text[index]) {
		index++
	}
	return index
}

func consumeBracketExpression(text string, start int) (int, bool) {
	if start >= len(text) || text[start] != '[' {
		return start, false
	}

	depth := 0
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return i + 1, true
			}
		}
	}

	return start, false
}

func hasAssignmentOperator(text string, index int) bool {
	if index < 0 || index >= len(text) {
		return false
	}

	operators := []string{"<<=", ">>=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "="}
	for _, operator := range operators {
		if !strings.HasPrefix(text[index:], operator) {
			continue
		}
		if operator == "=" && index+1 < len(text) && text[index+1] == '=' {
			continue
		}
		return true
	}

	return false
}

func selectDefinitionOccurrence(occurrences []symbolOccurrence) (symbolOccurrence, bool) {
	if len(occurrences) == 0 {
		return symbolOccurrence{}, false
	}

	for _, occurrence := range occurrences {
		if occurrence.isDefinition {
			return occurrence, true
		}
	}

	return occurrences[0], true
}

func sortSymbolOccurrences(occurrences []symbolOccurrence) {
	sort.SliceStable(occurrences, func(i, j int) bool {
		left := occurrences[i].token
		right := occurrences[j].token
		if left == nil || right == nil {
			return i < j
		}
		if left.GetStart() == right.GetStart() {
			return left.GetStop() < right.GetStop()
		}
		return left.GetStart() < right.GetStart()
	})
}

func sameTokenIndex(left antlr.Token, right antlr.Token) bool {
	if left == nil || right == nil {
		return false
	}
	return left.GetTokenIndex() == right.GetTokenIndex()
}

func tokenSpan(token antlr.Token) int {
	if token == nil {
		return 0
	}
	start := token.GetStart()
	stop := token.GetStop()
	if start < 0 || stop < start {
		return 0
	}
	return stop - start
}
