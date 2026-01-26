package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestParseDocument_ValidHasNoDiagnostics(t *testing.T) {
	input := "kprobe:sys_clone { @x = count(); }\n"
	result := ParseDocument(input)
	if result == nil {
		t.Fatalf("expected parse result")
	}
	if len(result.Diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %d", len(result.Diagnostics))
	}
}

func TestParseDocument_InvalidHasDiagnostics(t *testing.T) {
	input := "kprobe:sys_clone { @x = count( }\n"
	result := ParseDocument(input)
	if result == nil {
		t.Fatalf("expected parse result")
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected diagnostics")
	}
}

func TestDiagnosticsFromErrors_FallbackPositions(t *testing.T) {
	input := "kprobe:sys_clone {\n  @x = count();\n}\n"
	errors := []parseError{{start: -1, stop: -1, line: 2, column: 3, message: "boom"}}
	diagnostics := diagnosticsFromErrors(input, errors)
	if len(diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diagnostics))
	}
	diagnostic := diagnostics[0]
	wantStart := PositionForLineColumn(input, 2, 3)
	if diagnostic.Range.Start != wantStart {
		t.Fatalf("start = %+v, want %+v", diagnostic.Range.Start, wantStart)
	}
	wantMessage := "line 2:3: boom"
	if diagnostic.Message != wantMessage {
		t.Fatalf("message = %q, want %q", diagnostic.Message, wantMessage)
	}
	if diagnostic.Severity == nil || *diagnostic.Severity != protocol.DiagnosticSeverityError {
		t.Fatalf("expected error severity")
	}
	if diagnostic.Source == nil || *diagnostic.Source != "btfmt" {
		t.Fatalf("expected source btfmt")
	}
}
