package lsp

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func hoverForPosition(doc *Document, pos protocol.Position) *protocol.Hover {
	if doc == nil || doc.ParseResult == nil || doc.ParseResult.Tree == nil {
		return nil
	}

	var bestCtx antlr.ParserRuleContext
	var bestRange protocol.Range
	bestSpan := 0
	matched := false

	var walk func(node antlr.Tree)
	walk = func(node antlr.Tree) {
		if node == nil {
			return
		}

		if ctx, ok := node.(antlr.ParserRuleContext); ok {
			startToken := ctx.GetStart()
			stopToken := ctx.GetStop()
			if startToken != nil && stopToken != nil {
				start := startToken.GetStart()
				stop := stopToken.GetStop()
				if start >= 0 && stop >= start {
					currentRange := protocol.Range{
						Start: PositionForOffset(doc.Text, start),
						End:   PositionForOffset(doc.Text, stop+1),
					}
					if positionInRange(pos, currentRange) {
						span := stop - start
						if !matched || span < bestSpan {
							bestCtx = ctx
							bestRange = currentRange
							bestSpan = span
							matched = true
						}
					}
				}
			}
		}

		for i := 0; i < node.GetChildCount(); i++ {
			child := node.GetChild(i)
			if child != nil {
				walk(child)
			}
		}
	}

	walk(doc.ParseResult.Tree)
	if !matched {
		return nil
	}

	if identifier, identifierRange, sigilPrefixed, ok := identifierAtPosition(doc.Text, pos); ok && !sigilPrefixed {
		if markdown, ok := semanticHoverMarkdown(identifier); ok {
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: markdown,
				},
				Range: &identifierRange,
			}
		}
	}

	rule := ruleNameForContext(bestCtx)
	snippet := snippetForContext(doc.Text, bestCtx)
	value := fmt.Sprintf("%s: %s", rule, snippet)

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindPlainText,
			Value: value,
		},
		Range: &bestRange,
	}
}

// HoverForPosition returns hover info for a document position.
func HoverForPosition(doc *Document, pos protocol.Position) *protocol.Hover {
	return hoverForPosition(doc, pos)
}

func semanticHoverMarkdown(identifier string) (string, bool) {
	if probe, ok := lookupProbeType(identifier); ok {
		return fmt.Sprintf("**Probe Type** `%s`\n\n```bpftrace\n%s\n```\n\n%s", probe.name, probe.detail, probe.doc), true
	}

	if function, ok := lookupMapFunction(identifier); ok {
		return fmt.Sprintf("**Map Function** `%s`\n\n```bpftrace\n%s\n```\n\n%s", function.name, function.detail, function.doc), true
	}

	if function, ok := lookupBuiltinFunction(identifier); ok {
		if strings.Contains(function.detail, "(") {
			return fmt.Sprintf("**Builtin Function** `%s`\n\n```bpftrace\n%s\n```\n\n%s", function.name, function.detail, function.doc), true
		}
		return fmt.Sprintf("**Builtin Constant** `%s`\n\n%s", function.name, function.doc), true
	}

	if keywordDoc, ok := keywordHoverDoc(identifier); ok {
		return keywordDoc, true
	}

	return "", false
}

func lookupProbeType(identifier string) (struct {
	name    string
	detail  string
	doc     string
	example string
}, bool) {
	for _, probe := range probeTypes {
		if probe.name == identifier {
			return probe, true
		}
	}
	return struct {
		name    string
		detail  string
		doc     string
		example string
	}{}, false
}

func lookupMapFunction(identifier string) (struct {
	name   string
	detail string
	doc    string
}, bool) {
	for _, function := range mapFunctions {
		if function.name == identifier {
			return function, true
		}
	}
	return struct {
		name   string
		detail string
		doc    string
	}{}, false
}

func lookupBuiltinFunction(identifier string) (struct {
	name   string
	detail string
	doc    string
}, bool) {
	for _, function := range builtinFunctions {
		if function.name == identifier {
			return function, true
		}
	}
	return struct {
		name   string
		detail string
		doc    string
	}{}, false
}

func keywordHoverDoc(identifier string) (string, bool) {
	switch identifier {
	case "if":
		return "**Keyword** `if`\n\nConditional branch.", true
	case "else":
		return "**Keyword** `else`\n\nAlternative branch for `if`.", true
	case "while":
		return "**Keyword** `while`\n\nLoop while condition is true.", true
	case "for":
		return "**Keyword** `for`\n\nLoop over ranges or map keys.", true
	case "return":
		return "**Keyword** `return`\n\nReturn from current probe or macro.", true
	case "sizeof":
		return "**Keyword** `sizeof`\n\nReturn the size of a type or expression.", true
	default:
		return "", false
	}
}

func identifierAtPosition(text string, pos protocol.Position) (string, protocol.Range, bool, bool) {
	offset := offsetForPosition(text, pos)
	start, end, ok := identifierByteRangeAtOffset(text, offset)
	if !ok {
		return "", protocol.Range{}, false, false
	}

	startRune := runeOffsetForByteOffset(text, start)
	endRune := runeOffsetForByteOffset(text, end)
	rangeValue := protocol.Range{
		Start: PositionForOffset(text, startRune),
		End:   PositionForOffset(text, endRune),
	}

	sigilPrefixed := start > 0 && isSigilPrefix(text[start-1])
	return text[start:end], rangeValue, sigilPrefixed, true
}

func isSigilPrefix(value byte) bool {
	return value == '$' || value == '@'
}

func identifierByteRangeAtOffset(text string, offset int) (int, int, bool) {
	if text == "" {
		return 0, 0, false
	}
	if offset < 0 {
		offset = 0
	}
	if offset > len(text) {
		offset = len(text)
	}

	if offset == len(text) || !isIdentifierByte(text[offset]) {
		if offset == 0 || !isIdentifierByte(text[offset-1]) {
			return 0, 0, false
		}
		offset--
	}

	start := offset
	for start > 0 && isIdentifierByte(text[start-1]) {
		start--
	}

	end := offset + 1
	for end < len(text) && isIdentifierByte(text[end]) {
		end++
	}

	if start >= end {
		return 0, 0, false
	}
	return start, end, true
}

func isIdentifierByte(value byte) bool {
	return (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z') || (value >= '0' && value <= '9') || value == '_'
}

func runeOffsetForByteOffset(text string, byteOffset int) int {
	if byteOffset <= 0 {
		return 0
	}
	if byteOffset >= len(text) {
		return len([]rune(text))
	}
	return len([]rune(text[:byteOffset]))
}

func positionInRange(pos protocol.Position, rangeValue protocol.Range) bool {
	if !positionLessOrEqual(rangeValue.Start, pos) {
		return false
	}
	return positionLess(pos, rangeValue.End)
}

func positionLess(a protocol.Position, b protocol.Position) bool {
	if a.Line < b.Line {
		return true
	}
	if a.Line > b.Line {
		return false
	}
	return a.Character < b.Character
}

func positionLessOrEqual(a protocol.Position, b protocol.Position) bool {
	if a.Line < b.Line {
		return true
	}
	if a.Line > b.Line {
		return false
	}
	return a.Character <= b.Character
}

func ruleNameForContext(ctx antlr.ParserRuleContext) string {
	if ctx == nil {
		return "node"
	}
	if provider, ok := ctx.(interface{ GetParser() antlr.Parser }); ok {
		parser := provider.GetParser()
		if parser != nil {
			ruleNames := parser.GetRuleNames()
			index := ctx.GetRuleIndex()
			if index >= 0 && index < len(ruleNames) {
				name := strings.TrimSpace(ruleNames[index])
				if name != "" {
					return name
				}
			}
		}
	}
	return "node"
}

func snippetForContext(text string, ctx antlr.ParserRuleContext) string {
	if ctx == nil {
		return ""
	}
	startToken := ctx.GetStart()
	stopToken := ctx.GetStop()
	if startToken != nil && stopToken != nil {
		start := startToken.GetStart()
		stop := stopToken.GetStop()
		if start >= 0 && stop >= start {
			if snippet, ok := sliceByRuneOffsets(text, start, stop+1); ok {
				return truncateRunes(strings.TrimSpace(snippet), 120)
			}
		}
	}

	snippet := strings.TrimSpace(ctx.GetText())
	return truncateRunes(snippet, 120)
}

func truncateRunes(value string, limit int) string {
	if limit <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}

func sliceByRuneOffsets(text string, start int, end int) (string, bool) {
	if start < 0 || end < start {
		return "", false
	}
	runes := []rune(text)
	if start > len(runes) {
		return "", false
	}
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end]), true
}
