package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestPositionForOffset_Unicode(t *testing.T) {
	text := "a\U0001F642b\nc"

	cases := []struct {
		name   string
		offset int
		want   protocol.Position
	}{
		{name: "after_emoji", offset: 2, want: protocol.Position{Line: 0, Character: 3}},
		{name: "after_newline", offset: 4, want: protocol.Position{Line: 1, Character: 0}},
		{name: "end", offset: 5, want: protocol.Position{Line: 1, Character: 1}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := PositionForOffset(text, tc.offset)
			if got != tc.want {
				t.Fatalf("PositionForOffset(%d) = %+v, want %+v", tc.offset, got, tc.want)
			}
		})
	}
}

func TestPositionForLineColumn_Unicode(t *testing.T) {
	text := "a\U0001F642b\nc"

	cases := []struct {
		name   string
		line   int
		column int
		want   protocol.Position
	}{
		{name: "line1_after_emoji", line: 1, column: 2, want: protocol.Position{Line: 0, Character: 3}},
		{name: "line1_after_b", line: 1, column: 3, want: protocol.Position{Line: 0, Character: 4}},
		{name: "line2_after_c", line: 2, column: 1, want: protocol.Position{Line: 1, Character: 1}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := PositionForLineColumn(text, tc.line, tc.column)
			if got != tc.want {
				t.Fatalf("PositionForLineColumn(%d, %d) = %+v, want %+v", tc.line, tc.column, got, tc.want)
			}
		})
	}
}
