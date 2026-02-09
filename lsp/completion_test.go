package lsp

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestProbeTypeCompletions(t *testing.T) {
	items := probeTypeCompletions()
	if len(items) == 0 {
		t.Error("expected probe type completions")
	}

	// Check for some expected probe types
	found := make(map[string]bool)
	for _, item := range items {
		found[item.Label] = true
	}

	expected := []string{"BEGIN", "END", "kprobe", "kretprobe", "tracepoint", "uprobe"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected probe type %q not found", e)
		}
	}
}

func TestFunctionCompletions(t *testing.T) {
	items := functionCompletions("")
	if len(items) == 0 {
		t.Error("expected function completions")
	}

	// Check for some expected functions
	found := make(map[string]bool)
	for _, item := range items {
		found[item.Label] = true
	}

	expected := []string{"printf", "print", "str", "ksym", "usym"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected function %q not found", e)
		}
	}
}

func TestFunctionCompletionsWithPrefix(t *testing.T) {
	items := functionCompletions("pr")
	if len(items) == 0 {
		t.Error("expected filtered function completions")
	}

	for _, item := range items {
		if item.Label[:2] != "pr" {
			t.Errorf("expected item starting with 'pr', got %q", item.Label)
		}
	}
}

func TestMapFunctionCompletions(t *testing.T) {
	items := mapFunctionCompletions("")
	if len(items) == 0 {
		t.Error("expected map function completions")
	}

	// Check for some expected map functions
	found := make(map[string]bool)
	for _, item := range items {
		found[item.Label] = true
	}

	expected := []string{"count", "sum", "avg", "min", "max", "hist"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected map function %q not found", e)
		}
	}
}

func TestStatementCompletions(t *testing.T) {
	items := statementCompletions("")
	if len(items) == 0 {
		t.Error("expected statement completions")
	}

	// Should include keywords
	found := make(map[string]bool)
	for _, item := range items {
		found[item.Label] = true
	}

	expected := []string{"if", "else", "while", "for", "return"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected keyword %q not found", e)
		}
	}
}

func TestDefaultCompletions(t *testing.T) {
	items := defaultCompletions()
	if len(items) == 0 {
		t.Error("expected default completions")
	}
}

func TestDetermineCompletionContext(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		line     uint32
		char     uint32
		wantKind completionContextKind
	}{
		{
			name:     "empty file",
			text:     "",
			line:     0,
			char:     0,
			wantKind: contextProbeStart,
		},
		{
			name:     "at start of file",
			text:     "k",
			line:     0,
			char:     1,
			wantKind: contextProbeStart,
		},
		{
			name:     "map name after @",
			text:     "kprobe:foo { @",
			line:     0,
			char:     14,
			wantKind: contextMapName,
		},
		{
			name:     "variable after $",
			text:     "kprobe:foo { $",
			line:     0,
			char:     14,
			wantKind: contextVariable,
		},
		{
			name:     "inside block",
			text:     "kprobe:foo { ",
			line:     0,
			char:     13,
			wantKind: contextStatement,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseResult := ParseDocument(tt.text)
			doc := &Document{
				Text:        tt.text,
				ParseResult: parseResult,
			}
			pos := protocol.Position{Line: tt.line, Character: tt.char}
			ctx := determineCompletionContext(doc, pos)
			if ctx.kind != tt.wantKind {
				t.Errorf("determineCompletionContext() = %v, want %v", ctx.kind, tt.wantKind)
			}
		})
	}
}

func TestCollectMapNames(t *testing.T) {
	text := `kprobe:foo {
		@counter = count();
		@data[pid] = 1;
		@stats = sum(@data[pid]);
	}`
	parseResult := ParseDocument(text)
	doc := &Document{
		Text:        text,
		ParseResult: parseResult,
	}

	maps := collectMapNames(doc)
	if len(maps) == 0 {
		t.Error("expected to find map names")
	}

	found := make(map[string]bool)
	for _, name := range maps {
		found[name] = true
	}

	// Note: the exact names depend on how the lexer tokenizes
	// At minimum we should find some maps
	if len(found) == 0 {
		t.Error("expected to find at least one map name")
	}
}

func TestCollectVariableNames(t *testing.T) {
	text := `kprobe:foo {
		$x = 1;
		$y = $x + 2;
		printf("%d\n", $y);
	}`
	parseResult := ParseDocument(text)
	doc := &Document{
		Text:        text,
		ParseResult: parseResult,
	}

	vars := collectVariableNames(doc)
	if len(vars) == 0 {
		t.Error("expected to find variable names")
	}

	found := make(map[string]bool)
	for _, name := range vars {
		found[name] = true
	}

	// Note: the exact names depend on how the lexer tokenizes
	// At minimum we should find some variables
	if len(found) == 0 {
		t.Error("expected to find at least one variable name")
	}
}

func TestIsValidIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"foo", true},
		{"_foo", true},
		{"foo123", true},
		{"Foo", true},
		{"123foo", false},
		{"", false},
		{"foo-bar", false},
		{"foo bar", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isValidIdentifier(tt.input)
			if got != tt.want {
				t.Errorf("isValidIdentifier(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCompletionForPosition(t *testing.T) {
	text := `kprobe:vfs_read {
		@count = count();
		$x = 1;
		printf("%d\n", $x);
	}`
	parseResult := ParseDocument(text)
	doc := &Document{
		Text:        text,
		ParseResult: parseResult,
	}

	// Test completion at various positions
	pos := protocol.Position{Line: 1, Character: 2}
	items := CompletionForPosition(doc, pos)
	if len(items) == 0 {
		t.Error("expected completion items")
	}
}

func TestCompletionForNilDocument(t *testing.T) {
	items := CompletionForPosition(nil, protocol.Position{})
	if len(items) == 0 {
		t.Error("expected default completions for nil document")
	}
}
