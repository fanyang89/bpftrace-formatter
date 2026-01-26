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
