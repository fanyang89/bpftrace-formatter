package lsp

import (
	"unicode/utf16"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

// EndPosition returns the last line and UTF-16 column for the text.
func EndPosition(text string) protocol.Position {
	if text == "" {
		return protocol.Position{Line: 0, Character: 0}
	}

	var line uint32
	var column uint32
	for _, r := range text {
		if r == '\n' {
			line++
			column = 0
			continue
		}
		column += utf16ColumnWidth(r)
	}

	return protocol.Position{Line: line, Character: column}
}

// PositionForOffset returns the line and UTF-16 column for a rune offset.
func PositionForOffset(text string, offset int) protocol.Position {
	if offset <= 0 || text == "" {
		return protocol.Position{Line: 0, Character: 0}
	}

	var line uint32
	var column uint32
	var runeIndex int
	for _, r := range text {
		if runeIndex >= offset {
			break
		}
		if r == '\n' {
			line++
			column = 0
		} else {
			column += utf16ColumnWidth(r)
		}
		runeIndex++
	}

	if runeIndex < offset {
		return EndPosition(text)
	}

	return protocol.Position{Line: line, Character: column}
}

// PositionForLineColumn returns the line and UTF-16 column for a rune column in a line.
func PositionForLineColumn(text string, line1Based int, column0Based int) protocol.Position {
	if text == "" {
		return protocol.Position{Line: 0, Character: 0}
	}
	if line1Based <= 1 {
		line1Based = 1
	}
	if column0Based < 0 {
		column0Based = 0
	}

	targetLine := uint32(line1Based - 1)
	var line uint32
	var column uint32
	var columnRunes int
	for _, r := range text {
		if line < targetLine {
			if r == '\n' {
				line++
				column = 0
				columnRunes = 0
			}
			continue
		}
		if line > targetLine {
			break
		}
		if columnRunes >= column0Based {
			break
		}
		if r == '\n' {
			break
		}
		column += utf16ColumnWidth(r)
		columnRunes++
	}

	if line < targetLine {
		return EndPosition(text)
	}

	return protocol.Position{Line: line, Character: column}
}

func utf16ColumnWidth(r rune) uint32 {
	if utf16.RuneLen(r) == 2 {
		return 2
	}
	return 1
}
