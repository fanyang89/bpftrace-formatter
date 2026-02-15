package lsp

import (
	"strings"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestSnippetForContext_UnicodeOffsets(t *testing.T) {
	input := "// 注释\nkprobe:sys_clone { @x = count(); }\n"
	result := ParseDocument(input)
	if result == nil || result.Tree == nil || result.Tree.Content() == nil {
		t.Fatalf("expected parse tree")
	}

	probes := result.Tree.Content().AllProbe()
	if len(probes) == 0 {
		t.Fatalf("expected probe")
	}
	probeList := probes[0].Probe_list()
	if probeList == nil {
		t.Fatalf("expected probe list")
	}

	got := snippetForContext(input, probeList)
	if got != "kprobe:sys_clone" {
		t.Fatalf("snippetForContext = %q, want %q", got, "kprobe:sys_clone")
	}
}

func TestHoverForPosition_BuiltinFunctionMarkdown(t *testing.T) {
	input := "BEGIN { printf(\"x\"); }\n"
	result := ParseDocument(input)
	doc := &Document{Text: input, ParseResult: result}

	offset := strings.Index(input, "printf")
	if offset < 0 {
		t.Fatalf("missing printf in input")
	}

	hover := HoverForPosition(doc, PositionForOffset(input, offset+2))
	if hover == nil {
		t.Fatalf("expected hover")
	}
	content, ok := hover.Contents.(protocol.MarkupContent)
	if !ok {
		t.Fatalf("expected markup content, got %T", hover.Contents)
	}
	if content.Kind != protocol.MarkupKindMarkdown {
		t.Fatalf("markup kind = %s, want %s", content.Kind, protocol.MarkupKindMarkdown)
	}
	if !strings.Contains(content.Value, "Builtin Function") {
		t.Fatalf("expected builtin function heading, got %q", content.Value)
	}
	if !strings.Contains(content.Value, "printf(fmt, ...)") {
		t.Fatalf("expected function signature, got %q", content.Value)
	}
}

func TestHoverForPosition_ProbeTypeMarkdown(t *testing.T) {
	input := "kprobe:sys_clone { @x = count(); }\n"
	result := ParseDocument(input)
	doc := &Document{Text: input, ParseResult: result}

	offset := strings.Index(input, "kprobe")
	if offset < 0 {
		t.Fatalf("missing kprobe in input")
	}

	hover := HoverForPosition(doc, PositionForOffset(input, offset+1))
	if hover == nil {
		t.Fatalf("expected hover")
	}
	content, ok := hover.Contents.(protocol.MarkupContent)
	if !ok {
		t.Fatalf("expected markup content, got %T", hover.Contents)
	}
	if content.Kind != protocol.MarkupKindMarkdown {
		t.Fatalf("markup kind = %s, want %s", content.Kind, protocol.MarkupKindMarkdown)
	}
	if !strings.Contains(content.Value, "Probe Type") {
		t.Fatalf("expected probe heading, got %q", content.Value)
	}
	if !strings.Contains(content.Value, "kprobe:function") {
		t.Fatalf("expected probe signature, got %q", content.Value)
	}
}

func TestHoverForPosition_SigilPrefixedVariableFallsBackToSyntaxHover(t *testing.T) {
	input := "BEGIN { $pid = 1; print($pid); }\n"
	result := ParseDocument(input)
	doc := &Document{Text: input, ParseResult: result}

	offset := strings.Index(input, "$pid")
	if offset < 0 {
		t.Fatalf("missing $pid in input")
	}

	hover := HoverForPosition(doc, PositionForOffset(input, offset+1))
	if hover == nil {
		t.Fatalf("expected hover")
	}
	content, ok := hover.Contents.(protocol.MarkupContent)
	if !ok {
		t.Fatalf("expected markup content, got %T", hover.Contents)
	}
	if content.Kind != protocol.MarkupKindPlainText {
		t.Fatalf("markup kind = %s, want %s", content.Kind, protocol.MarkupKindPlainText)
	}
	if strings.Contains(content.Value, "Builtin Constant") {
		t.Fatalf("expected non-semantic hover for sigil-prefixed variable, got %q", content.Value)
	}
}

func TestHoverForPosition_SigilPrefixedMapFallsBackToSyntaxHover(t *testing.T) {
	input := "BEGIN { @count = 1; print(@count); }\n"
	result := ParseDocument(input)
	doc := &Document{Text: input, ParseResult: result}

	offset := strings.Index(input, "@count")
	if offset < 0 {
		t.Fatalf("missing @count in input")
	}

	hover := HoverForPosition(doc, PositionForOffset(input, offset+1))
	if hover == nil {
		t.Fatalf("expected hover")
	}
	content, ok := hover.Contents.(protocol.MarkupContent)
	if !ok {
		t.Fatalf("expected markup content, got %T", hover.Contents)
	}
	if content.Kind != protocol.MarkupKindPlainText {
		t.Fatalf("markup kind = %s, want %s", content.Kind, protocol.MarkupKindPlainText)
	}
	if strings.Contains(content.Value, "Map Function") {
		t.Fatalf("expected non-semantic hover for sigil-prefixed map, got %q", content.Value)
	}
}

func TestHoverForPosition_StringLiteralFallsBackToSyntaxHover(t *testing.T) {
	input := "BEGIN { printf(\"pid\"); }\n"
	result := ParseDocument(input)
	doc := &Document{Text: input, ParseResult: result}

	offset := strings.Index(input, "pid")
	if offset < 0 {
		t.Fatalf("missing pid in input")
	}

	hover := HoverForPosition(doc, PositionForOffset(input, offset+1))
	if hover == nil {
		t.Fatalf("expected hover")
	}
	content, ok := hover.Contents.(protocol.MarkupContent)
	if !ok {
		t.Fatalf("expected markup content, got %T", hover.Contents)
	}
	if content.Kind != protocol.MarkupKindPlainText {
		t.Fatalf("markup kind = %s, want %s", content.Kind, protocol.MarkupKindPlainText)
	}
	if strings.Contains(content.Value, "Builtin Constant") {
		t.Fatalf("expected non-semantic hover inside string literal, got %q", content.Value)
	}
}

func TestHoverForPosition_FieldAccessFallsBackToSyntaxHover(t *testing.T) {
	input := "kprobe:sys_clone { printf(\"%d\", args.pid); }\n"
	result := ParseDocument(input)
	doc := &Document{Text: input, ParseResult: result}

	offset := strings.Index(input, "pid")
	if offset < 0 {
		t.Fatalf("missing pid in input")
	}

	hover := HoverForPosition(doc, PositionForOffset(input, offset+1))
	if hover == nil {
		t.Fatalf("expected hover")
	}
	content, ok := hover.Contents.(protocol.MarkupContent)
	if !ok {
		t.Fatalf("expected markup content, got %T", hover.Contents)
	}
	if content.Kind != protocol.MarkupKindPlainText {
		t.Fatalf("markup kind = %s, want %s", content.Kind, protocol.MarkupKindPlainText)
	}
	if strings.Contains(content.Value, "Builtin Constant") {
		t.Fatalf("expected non-semantic hover on field access, got %q", content.Value)
	}
}

func TestHoverForPosition_MultilineFieldAccessFallsBackToSyntaxHover(t *testing.T) {
	input := "kprobe:sys_clone {\n  printf(\"%d\", args.\n    pid);\n}\n"
	result := ParseDocument(input)
	doc := &Document{Text: input, ParseResult: result}

	offset := strings.Index(input, "pid")
	if offset < 0 {
		t.Fatalf("missing pid in input")
	}

	hover := HoverForPosition(doc, PositionForOffset(input, offset+1))
	if hover == nil {
		t.Fatalf("expected hover")
	}
	content, ok := hover.Contents.(protocol.MarkupContent)
	if !ok {
		t.Fatalf("expected markup content, got %T", hover.Contents)
	}
	if content.Kind != protocol.MarkupKindPlainText {
		t.Fatalf("markup kind = %s, want %s", content.Kind, protocol.MarkupKindPlainText)
	}
	if strings.Contains(content.Value, "Builtin Constant") {
		t.Fatalf("expected non-semantic hover on multiline field access, got %q", content.Value)
	}
}

func TestHoverForPosition_CommentFallsBackToSyntaxHover(t *testing.T) {
	input := "BEGIN { // pid\n  @x = 1; }\n"
	result := ParseDocument(input)
	doc := &Document{Text: input, ParseResult: result}

	offset := strings.Index(input, "pid")
	if offset < 0 {
		t.Fatalf("missing pid in input")
	}

	hover := HoverForPosition(doc, PositionForOffset(input, offset+1))
	if hover == nil {
		t.Fatalf("expected hover")
	}
	content, ok := hover.Contents.(protocol.MarkupContent)
	if !ok {
		t.Fatalf("expected markup content, got %T", hover.Contents)
	}
	if content.Kind != protocol.MarkupKindPlainText {
		t.Fatalf("markup kind = %s, want %s", content.Kind, protocol.MarkupKindPlainText)
	}
	if strings.Contains(content.Value, "Builtin Constant") {
		t.Fatalf("expected non-semantic hover in comment, got %q", content.Value)
	}
}
