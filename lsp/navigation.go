package lsp

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/parser"
)

type symbolKind int

const (
	symbolKindVariable symbolKind = iota
	symbolKindMap
)

type symbolOccurrence struct {
	kind          symbolKind
	name          string
	scopeKey      string
	sigil         byte
	rangeValue    protocol.Range
	token         antlr.Token
	isDefinition  bool
	isDeclaration bool
}

func definitionLocationForPosition(doc *Document, pos protocol.Position) (protocol.Location, bool) {
	target, ok := symbolAtPosition(doc, pos)
	if !ok {
		return protocol.Location{}, false
	}

	occurrences := symbolOccurrencesByIdentity(doc, *target)
	if len(occurrences) == 0 {
		return protocol.Location{}, false
	}
	definition, _ := selectDefinitionOccurrence(occurrences)

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

	occurrences := symbolOccurrencesByIdentity(doc, *target)
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

	occurrences := symbolOccurrencesByIdentity(doc, *target)
	if len(occurrences) == 0 {
		return []protocol.Location{}
	}

	declaration, hasDeclaration := selectDeclarationOccurrence(occurrences)
	if !hasDeclaration {
		declaration, hasDeclaration = selectDefinitionOccurrence(occurrences)
	}
	locations := make([]protocol.Location, 0, len(occurrences))
	for _, occurrence := range occurrences {
		if !includeDeclaration && hasDeclaration && sameTokenIndex(occurrence.token, declaration.token) {
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

	occurrences := symbolOccurrencesByIdentity(doc, *target)
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

func symbolOccurrencesByIdentity(doc *Document, target symbolOccurrence) []symbolOccurrence {
	occurrences := collectSymbolOccurrences(doc)
	filtered := make([]symbolOccurrence, 0, len(occurrences))
	for _, occurrence := range occurrences {
		if occurrence.kind != target.kind || occurrence.name != target.name {
			continue
		}
		if (target.kind == symbolKindVariable || target.kind == symbolKindMap) && occurrence.scopeKey != target.scopeKey {
			continue
		}
		filtered = append(filtered, occurrence)
	}
	return filtered
}

func collectSymbolOccurrences(doc *Document) []symbolOccurrence {
	if doc == nil || doc.ParseResult == nil || doc.ParseResult.Tree == nil {
		return nil
	}

	occurrences := make([]symbolOccurrence, 0)
	var walk func(node antlr.Tree, variableScopeKey string, macroMapParamNames map[string]struct{}, macroMapScopeKey string, macroParamDeclarationTokenIndexes map[int]struct{})
	walk = func(node antlr.Tree, variableScopeKey string, macroMapParamNames map[string]struct{}, macroMapScopeKey string, macroParamDeclarationTokenIndexes map[int]struct{}) {
		if node == nil {
			return
		}

		currentScopeKey := variableScopeKey
		currentMacroMapScopeKey := macroMapScopeKey
		currentMacroMapParamNames := macroMapParamNames
		currentMacroParamDeclarationTokenIndexes := macroParamDeclarationTokenIndexes
		switch typed := node.(type) {
		case *parser.ProbeContext:
			currentScopeKey = variableScopeKeyForContext("probe", typed)
			currentMacroMapScopeKey = ""
			currentMacroMapParamNames = nil
			currentMacroParamDeclarationTokenIndexes = nil
		case *parser.Macro_definitionContext:
			currentScopeKey = variableScopeKeyForContext("macro", typed)
			currentMacroMapScopeKey = currentScopeKey
			currentMacroMapParamNames = macroMapParamSet(typed)
			currentMacroParamDeclarationTokenIndexes = macroParamDeclarationTokenIndexSet(typed)
		}

		if terminal, ok := node.(antlr.TerminalNode); ok {
			if occurrence, ok := symbolOccurrenceFromTerminal(doc.Text, terminal, currentScopeKey, currentMacroMapParamNames, currentMacroMapScopeKey, currentMacroParamDeclarationTokenIndexes); ok {
				occurrences = append(occurrences, occurrence)
			}
		}

		for i := 0; i < node.GetChildCount(); i++ {
			walk(node.GetChild(i), currentScopeKey, currentMacroMapParamNames, currentMacroMapScopeKey, currentMacroParamDeclarationTokenIndexes)
		}
	}

	walk(doc.ParseResult.Tree, "", nil, "", nil)
	sortSymbolOccurrences(occurrences)
	return occurrences
}

func macroMapParamSet(ctx *parser.Macro_definitionContext) map[string]struct{} {
	if ctx == nil || ctx.Macro_params() == nil {
		return nil
	}

	result := make(map[string]struct{})
	for _, param := range ctx.Macro_params().AllMacro_param() {
		if param == nil || param.MAP_NAME() == nil {
			continue
		}
		kind, name, _, ok := symbolIdentityFromText(param.MAP_NAME().GetText())
		if !ok || kind != symbolKindMap {
			continue
		}
		result[name] = struct{}{}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func macroParamDeclarationTokenIndexSet(ctx *parser.Macro_definitionContext) map[int]struct{} {
	if ctx == nil || ctx.Macro_params() == nil {
		return nil
	}

	result := make(map[int]struct{})
	for _, param := range ctx.Macro_params().AllMacro_param() {
		if param == nil {
			continue
		}

		if mapParam := param.MAP_NAME(); mapParam != nil && mapParam.GetSymbol() != nil {
			result[mapParam.GetSymbol().GetTokenIndex()] = struct{}{}
		}
		if variableParam := param.VARIABLE(); variableParam != nil && variableParam.GetSymbol() != nil {
			result[variableParam.GetSymbol().GetTokenIndex()] = struct{}{}
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func variableScopeKeyForContext(prefix string, ctx antlr.ParserRuleContext) string {
	if ctx == nil {
		return prefix
	}

	startIndex := -1
	if startToken := ctx.GetStart(); startToken != nil {
		startIndex = startToken.GetTokenIndex()
	}

	stopIndex := -1
	if stopToken := ctx.GetStop(); stopToken != nil {
		stopIndex = stopToken.GetTokenIndex()
	}

	return fmt.Sprintf("%s:%d:%d", prefix, startIndex, stopIndex)
}

func symbolOccurrenceFromTerminal(text string, terminal antlr.TerminalNode, variableScopeKey string, macroMapParamNames map[string]struct{}, macroMapScopeKey string, macroParamDeclarationTokenIndexes map[int]struct{}) (symbolOccurrence, bool) {
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

	scopeKey := ""
	if kind == symbolKindVariable {
		scopeKey = variableScopeKey
	}
	if kind == symbolKindMap {
		if _, ok := macroMapParamNames[name]; ok {
			scopeKey = macroMapScopeKey
		}
	}

	_, isDeclaration := macroParamDeclarationTokenIndexes[token.GetTokenIndex()]

	return symbolOccurrence{
		kind:          kind,
		name:          name,
		scopeKey:      scopeKey,
		sigil:         sigil,
		rangeValue:    rangeValue,
		token:         token,
		isDefinition:  tokenRepresentsAssignmentLHS(text, token, sigil),
		isDeclaration: isDeclaration,
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
	runes := []rune(text)

	start := token.GetStart()
	stop := token.GetStop()
	if start < 0 || stop < start || stop+1 > len(runes) {
		return false
	}

	index := stop + 1
	index = skipASCIIWhitespaceForward(runes, index)

	if sigil == '@' && index < len(runes) && runes[index] == '[' {
		next, ok := consumeBracketExpression(runes, index)
		if !ok {
			return false
		}
		index = skipASCIIWhitespaceForward(runes, next)
	}

	return hasAssignmentOperator(runes, index)
}

func skipASCIIWhitespaceForward(text []rune, index int) int {
	for index < len(text) && isASCIIWhitespaceRune(text[index]) {
		index++
	}
	return index
}

func consumeBracketExpression(text []rune, start int) (int, bool) {
	if start >= len(text) || text[start] != '[' {
		return start, false
	}

	depth := 0
	var quote rune
	escaped := false
	for i := start; i < len(text); i++ {
		if quote != 0 {
			if escaped {
				escaped = false
				continue
			}
			if text[i] == '\\' {
				escaped = true
				continue
			}
			if text[i] == quote {
				quote = 0
			}
			continue
		}

		switch text[i] {
		case '\'', '"':
			quote = text[i]
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

func hasAssignmentOperator(text []rune, index int) bool {
	if index < 0 || index >= len(text) {
		return false
	}

	operators := []string{"<<=", ">>=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "="}
	for _, operator := range operators {
		if !hasRunePrefix(text, index, operator) {
			continue
		}
		if operator == "=" && index+1 < len(text) && text[index+1] == '=' {
			continue
		}
		return true
	}

	return false
}

func isASCIIWhitespaceRune(value rune) bool {
	switch value {
	case ' ', '\t', '\n', '\r', '\f', '\v':
		return true
	default:
		return false
	}
}

func hasRunePrefix(text []rune, index int, prefix string) bool {
	prefixRunes := []rune(prefix)
	if index < 0 || index+len(prefixRunes) > len(text) {
		return false
	}
	for i, char := range prefixRunes {
		if text[index+i] != char {
			return false
		}
	}
	return true
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

	return occurrences[0], false
}

func selectDeclarationOccurrence(occurrences []symbolOccurrence) (symbolOccurrence, bool) {
	if len(occurrences) == 0 {
		return symbolOccurrence{}, false
	}

	for _, occurrence := range occurrences {
		if occurrence.isDeclaration {
			return occurrence, true
		}
	}

	return symbolOccurrence{}, false
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
