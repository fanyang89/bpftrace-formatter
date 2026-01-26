package lsp

import (
	"unicode/utf16"
	"unicode/utf8"

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

// PositionForOffset returns the line and UTF-16 column for a byte offset.
func PositionForOffset(text string, offset int) protocol.Position {
	if offset <= 0 || text == "" {
		return protocol.Position{Line: 0, Character: 0}
	}
	if offset >= len(text) {
		return EndPosition(text)
	}

	var line uint32
	var column uint32
	for i := 0; i < len(text); {
		if i >= offset {
			break
		}
		r, size := utf8.DecodeRuneInString(text[i:])
		if i+size > offset {
			break
		}
		if r == '\n' {
			line++
			column = 0
			i += size
			continue
		}
		column += utf16ColumnWidth(r)
		i += size
	}

	return protocol.Position{Line: line, Character: column}
}

// PositionForLineColumn returns the line and UTF-16 column for a byte column in a line.
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
	lineStart := 0
	for i := 0; i < len(text) && line < targetLine; {
		r, size := utf8.DecodeRuneInString(text[i:])
		if r == '\n' {
			line++
			lineStart = i + size
		}
		i += size
	}

	if line < targetLine {
		return EndPosition(text)
	}

	targetByteIndex := lineStart + column0Based
	if targetByteIndex > len(text) {
		targetByteIndex = len(text)
	}

	var column uint32
	for i := lineStart; i < len(text); {
		if i >= targetByteIndex {
			break
		}
		r, size := utf8.DecodeRuneInString(text[i:])
		if r == '\n' {
			break
		}
		if i+size > targetByteIndex {
			break
		}
		column += utf16ColumnWidth(r)
		i += size
	}

	return protocol.Position{Line: line, Character: column}
}

func utf16ColumnWidth(r rune) uint32 {
	if utf16.RuneLen(r) == 2 {
		return 2
	}
	return 1
}
