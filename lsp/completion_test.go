package lsp

import (
	"reflect"
	"strings"
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

	expected := []string{"BEGIN", "END", "kprobe", "kretprobe", "tracepoint", "uprobe", "asyncwatchpoint"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected probe type %q not found", e)
		}
	}

	notExpected := []string{"kfunc", "kretfunc", "iter", "rawtracepoint"}
	for _, e := range notExpected {
		if found[e] {
			t.Errorf("unexpected unsupported probe type %q found", e)
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

func TestFunctionCompletions_UseSnippetInsertText(t *testing.T) {
	previous := snippetSupported.Load()
	snippetSupported.Store(true)
	t.Cleanup(func() { snippetSupported.Store(previous) })

	items := functionCompletions("printf")
	if len(items) == 0 {
		t.Fatalf("expected printf completion")
	}

	var found bool
	for _, item := range items {
		if item.Label != "printf" {
			continue
		}
		found = true
		if item.InsertText == nil || !strings.HasPrefix(*item.InsertText, "printf(") {
			t.Fatalf("expected printf snippet insert text, got %+v", item.InsertText)
		}
		if item.InsertTextFormat == nil || *item.InsertTextFormat != protocol.InsertTextFormatSnippet {
			t.Fatalf("expected snippet insert format")
		}
		if item.SortText == nil || *item.SortText == "" {
			t.Fatalf("expected sortText for printf item")
		}
		if item.FilterText == nil || *item.FilterText != "printf" {
			t.Fatalf("expected filterText=printf")
		}
	}

	if !found {
		t.Fatalf("missing printf completion item")
	}
}

func TestFunctionCompletions_FallsBackToPlainTextWithoutSnippetSupport(t *testing.T) {
	previous := snippetSupported.Load()
	snippetSupported.Store(false)
	t.Cleanup(func() { snippetSupported.Store(previous) })

	items := functionCompletions("printf")
	if len(items) == 0 {
		t.Fatalf("expected printf completion")
	}

	var found bool
	for _, item := range items {
		if item.Label != "printf" {
			continue
		}
		found = true
		if item.InsertText == nil || *item.InsertText != "printf()" {
			t.Fatalf("expected plain insert text printf(), got %+v", item.InsertText)
		}
		if item.InsertTextFormat != nil {
			t.Fatalf("expected nil insert text format without snippet support")
		}
	}

	if !found {
		t.Fatalf("missing printf completion item")
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

	expected := []string{"if", "else", "while", "for", "return", "pid"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected completion %q not found", e)
		}
	}
}

func TestVariableCompletionsDoesNotIncludeBuiltinConstants(t *testing.T) {
	text := `kprobe:foo {
		$var = 1;
	}`
	parseResult := ParseDocument(text)
	doc := &Document{Text: text, ParseResult: parseResult}

	items := variableCompletions(doc, "")
	for _, item := range items {
		if item.Label == "pid" || item.Label == "tid" {
			t.Fatalf("unexpected builtin constant %q in variable completions", item.Label)
		}
	}
}

func TestDefaultCompletions(t *testing.T) {
	items := defaultCompletions()
	if len(items) == 0 {
		t.Error("expected default completions")
	}
}

func TestGetMapAssignmentPrefix(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		wantPrefix string
		wantOk     bool
	}{
		{"right after =", "@x = ", "", true},
		{"typing function", "@x = cou", "cou", true},
		{"with spaces", "@x =   sum", "sum", true},
		{"anonymous map right after =", "@ = ", "", true},
		{"anonymous indexed map typing", "@[pid] = su", "su", true},
		{"named indexed map with spaces", "@x [pid] = ", "", true},
		{"anonymous indexed map with spaces", "@ [pid] = su", "su", true},
		{"nested map index", "@dst[@src] = ", "", true},
		{"nested map index with condition", "if (@a) @b[@c] = ", "", true},
		{"commented map assignment", "// @x = ", "", false},
		{"unterminated string map assignment", "printf(\"@x = ", "", false},
		{"no @", "x = count", "", false},
		{"greater-equal comparison", "if (@x >= ", "", false},
		{"less-equal comparison", "if (@x <= ", "", false},
		{"double-equal comparison", "if (@x == ", "", false},
		{"not-equal comparison", "if (@x != ", "", false},
		{"map read in condition then scalar assignment", "if (@x) $y = ", "", false},
		{"map read in condition then map assignment", "if (@x) @y = ", "", true},
		{"rhs non-map identifier", "@x = pid", "", false},
		{"has operator", "@x = 1 + ", "", false},
		{"has paren", "@x = foo(", "", false},
		{"complete call", "@x = count()", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, ok := getMapAssignmentPrefix(tt.line)
			if ok != tt.wantOk {
				t.Errorf("getMapAssignmentPrefix(%q) ok = %v, want %v", tt.line, ok, tt.wantOk)
			}
			if ok && prefix != tt.wantPrefix {
				t.Errorf("getMapAssignmentPrefix(%q) prefix = %q, want %q", tt.line, prefix, tt.wantPrefix)
			}
		})
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
			name:     "probe type with colon",
			text:     "kprobe:",
			line:     0,
			char:     7,
			wantKind: contextProbeStart,
		},
		{
			name:     "probe target typing",
			text:     "kprobe:vfs_re",
			line:     0,
			char:     13,
			wantKind: contextProbeStart,
		},
		{
			name:     "uprobe path target typing",
			text:     "uprobe:/bin/bash:readl",
			line:     0,
			char:     22,
			wantKind: contextProbeStart,
		},
		{
			name:     "usdt path target typing",
			text:     "usdt:/usr/lib/libc.so.6:pro",
			line:     0,
			char:     27,
			wantKind: contextProbeStart,
		},
		{
			name:     "inline closed block then probe should be probe context",
			text:     "BEGIN { printf(\"x\"); } k",
			line:     0,
			char:     24,
			wantKind: contextProbeStart,
		},
		{
			name:     "line starts with closing brace then probe should be probe context",
			text:     "BEGIN {\n} k",
			line:     1,
			char:     3,
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
			name:     "map name typing",
			text:     "kprobe:foo { @cou",
			line:     0,
			char:     17,
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
			name:     "variable typing",
			text:     "kprobe:foo { $va",
			line:     0,
			char:     16,
			wantKind: contextVariable,
		},
		{
			name:     "inside block",
			text:     "kprobe:foo { ",
			line:     0,
			char:     13,
			wantKind: contextStatement,
		},
		{
			name:     "map function after =",
			text:     "kprobe:foo { @x = ",
			line:     0,
			char:     18,
			wantKind: contextMapFunction,
		},
		{
			name:     "map function after anonymous map =",
			text:     "kprobe:foo { @ = ",
			line:     0,
			char:     17,
			wantKind: contextMapFunction,
		},
		{
			name:     "map function typing",
			text:     "kprobe:foo { @x = cou",
			line:     0,
			char:     21,
			wantKind: contextMapFunction,
		},
		{
			name:     "map function after spaced map index",
			text:     "kprobe:foo { @x [pid] = ",
			line:     0,
			char:     24,
			wantKind: contextMapFunction,
		},
		{
			name:     "map function typing after spaced anonymous index",
			text:     "kprobe:foo { @ [pid] = su",
			line:     0,
			char:     25,
			wantKind: contextMapFunction,
		},
		{
			name:     "map function after nested map index",
			text:     "kprobe:foo { @dst[@src] = ",
			line:     0,
			char:     26,
			wantKind: contextMapFunction,
		},
		{
			name:     "map function after conditional nested map assignment",
			text:     "kprobe:foo { if (@a) @b[@c] = ",
			line:     0,
			char:     30,
			wantKind: contextMapFunction,
		},
		{
			name:     "non-map identifier after map assignment should be statement",
			text:     "kprobe:foo { @x = pid",
			line:     0,
			char:     21,
			wantKind: contextStatement,
		},
		{
			name:     "expression after map assignment - should be statement",
			text:     "kprobe:foo { @x = 1 + ",
			line:     0,
			char:     22,
			wantKind: contextStatement,
		},
		{
			name:     "inside function call after map assignment - should be statement",
			text:     "kprobe:foo { @x = foo(",
			line:     0,
			char:     22,
			wantKind: contextStatement,
		},
		{
			name:     "after semicolon in map line - should be statement",
			text:     "kprobe:foo { @x = count(); ",
			line:     0,
			char:     27,
			wantKind: contextStatement,
		},
		{
			name:     "comparison >= should be statement",
			text:     "kprobe:foo { if (@x >= ",
			line:     0,
			char:     23,
			wantKind: contextStatement,
		},
		{
			name:     "comparison <= should be statement",
			text:     "kprobe:foo { if (@x <= ",
			line:     0,
			char:     23,
			wantKind: contextStatement,
		},
		{
			name:     "comparison == should be statement",
			text:     "kprobe:foo { if (@x == ",
			line:     0,
			char:     23,
			wantKind: contextStatement,
		},
		{
			name:     "comparison != should be statement",
			text:     "kprobe:foo { if (@x != ",
			line:     0,
			char:     23,
			wantKind: contextStatement,
		},
		{
			name:     "map read then scalar assignment should be statement",
			text:     "kprobe:foo { if (@x) $y = ",
			line:     0,
			char:     27,
			wantKind: contextStatement,
		},
		{
			name:     "top-level after string brace should be probe context",
			text:     "BEGIN { printf(\"{\"); }\nk",
			line:     1,
			char:     1,
			wantKind: contextProbeStart,
		},
		{
			name:     "top-level after single-quoted brace should be probe context",
			text:     "BEGIN { printf('{'); }\nk",
			line:     1,
			char:     1,
			wantKind: contextProbeStart,
		},
		{
			name:     "probe predicate should be statement context",
			text:     "kprobe:vfs_read /pid == ",
			line:     0,
			char:     24,
			wantKind: contextStatement,
		},
		{
			name:     "probe predicate without leading space should be statement context",
			text:     "kprobe:vfs_read/pid==",
			line:     0,
			char:     21,
			wantKind: contextStatement,
		},
		{
			name:     "uprobe predicate without leading space should be statement context",
			text:     "uprobe:/bin/bash:readl/pid==",
			line:     0,
			char:     28,
			wantKind: contextStatement,
		},
		{
			name:     "map marker in string should be statement",
			text:     "kprobe:foo { printf(\"@",
			line:     0,
			char:     22,
			wantKind: contextStatement,
		},
		{
			name:     "map marker in comment should be statement",
			text:     "kprobe:foo { // @",
			line:     0,
			char:     17,
			wantKind: contextStatement,
		},
		{
			name:     "variable marker in comment should be statement",
			text:     "kprobe:foo { // $",
			line:     0,
			char:     17,
			wantKind: contextStatement,
		},
		{
			name:     "map assignment in comment should be statement",
			text:     "kprobe:foo { // @x = ",
			line:     0,
			char:     21,
			wantKind: contextStatement,
		},
		{
			name:     "map assignment in unterminated string should be statement",
			text:     "kprobe:foo { printf(\"@x = ",
			line:     0,
			char:     26,
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

func TestExtractLastWord(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{name: "identifier", line: "foo bar", want: "bar"},
		{name: "operator token", line: "@x +", want: ""},
		{name: "brace token", line: "{", want: ""},
		{name: "equals token", line: "@x =", want: ""},
		{name: "trailing punctuation", line: "if (@x)", want: "x"},
		{name: "function arg without space", line: "printf(pi", want: "pi"},
		{name: "paren expression without space", line: "if (pi", want: "pi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractLastWord(tt.line)
			if got != tt.want {
				t.Errorf("extractLastWord(%q) = %q, want %q", tt.line, got, tt.want)
			}
		})
	}
}

func TestMarkerPrefixInCode(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		marker   byte
		want     string
		wantOkay bool
	}{
		{name: "map marker", line: "@co", marker: '@', want: "co", wantOkay: true},
		{name: "variable marker", line: "$va", marker: '$', want: "va", wantOkay: true},
		{name: "map marker in string", line: "printf(\"@\")", marker: '@', want: "", wantOkay: false},
		{name: "map marker in comment", line: "// @foo", marker: '@', want: "", wantOkay: false},
		{name: "non-identifier suffix", line: "@x+", marker: '@', want: "", wantOkay: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := markerPrefixInCode(tt.line, tt.marker)
			if ok != tt.wantOkay {
				t.Fatalf("markerPrefixInCode(%q, %q) ok = %v, want %v", tt.line, tt.marker, ok, tt.wantOkay)
			}
			if ok && got != tt.want {
				t.Fatalf("markerPrefixInCode(%q, %q) = %q, want %q", tt.line, tt.marker, got, tt.want)
			}
		})
	}
}

func TestIsInStringOrComment(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{name: "normal code", line: "@x = ", want: false},
		{name: "in comment", line: "// @x = ", want: true},
		{name: "in unterminated string", line: "printf(\"@x = ", want: true},
		{name: "closed string then code", line: "printf(\"@\"); @x = ", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isInStringOrComment(tt.line)
			if got != tt.want {
				t.Errorf("isInStringOrComment(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestHasPredicateStart(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{name: "kprobe predicate spaced", line: "kprobe:vfs_read /pid==1/", want: true},
		{name: "kprobe predicate compact", line: "kprobe:vfs_read/pid==1/", want: true},
		{name: "uprobe path only", line: "uprobe:/bin/bash:readl", want: false},
		{name: "uprobe predicate compact", line: "uprobe:/bin/bash:readl/pid==1/", want: true},
		{name: "uprobe path with double slash", line: "uprobe:/tmp//bin/bash:readl/pid==1/", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasPredicateStart(tt.line)
			if got != tt.want {
				t.Errorf("hasPredicateStart(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestTopLevelLineTail(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{name: "simple line", line: "kprobe:", want: "kprobe:"},
		{name: "after closed inline block", line: "BEGIN { printf(\"x\"); } k", want: "k"},
		{name: "ignore comment after close", line: "BEGIN { } // comment", want: ""},
		{name: "closing brace from previous line", line: "} kprobe:", want: "kprobe:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topLevelLineTail(tt.line)
			if got != tt.want {
				t.Errorf("topLevelLineTail(%q) = %q, want %q", tt.line, got, tt.want)
			}
		})
	}
}

func TestIsInsideBlockIgnoresBracesInStringsAndComments(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "string brace ignored",
			text: "kprobe:foo { printf(\"{\"); }\n",
			want: false,
		},
		{
			name: "comment brace ignored",
			text: "kprobe:foo { // }\n",
			want: true,
		},
		{
			name: "escaped quote in string",
			text: "kprobe:foo { printf(\"\\\"{\\\"\"); }\n",
			want: false,
		},
		{
			name: "single-quoted brace ignored",
			text: "kprobe:foo { printf('{'); }\n",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isInsideBlock(tt.text)
			if got != tt.want {
				t.Errorf("isInsideBlock(%q) = %v, want %v", tt.text, got, tt.want)
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

func TestCollectMapNames_Sorted(t *testing.T) {
	text := `kprobe:foo {
		@z = count();
		@a = count();
		@m = count();
	}`
	parseResult := ParseDocument(text)
	doc := &Document{Text: text, ParseResult: parseResult}

	maps := collectMapNames(doc)
	if !reflect.DeepEqual(maps, []string{"a", "m", "z"}) {
		t.Fatalf("collectMapNames sorted = %v, want [a m z]", maps)
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

func TestCollectVariableNames_Sorted(t *testing.T) {
	text := `kprobe:foo {
		$z = 1;
		$a = 2;
		$m = 3;
	}`
	parseResult := ParseDocument(text)
	doc := &Document{Text: text, ParseResult: parseResult}

	vars := collectVariableNames(doc)
	if !reflect.DeepEqual(vars, []string{"a", "m", "z"}) {
		t.Fatalf("collectVariableNames sorted = %v, want [a m z]", vars)
	}
}

func TestMapCompletions_SortedOrder(t *testing.T) {
	text := `kprobe:foo {
		@z = count();
		@a = count();
	}`
	parseResult := ParseDocument(text)
	doc := &Document{Text: text, ParseResult: parseResult}

	items := mapCompletions(doc, "")
	if len(items) < 2 {
		t.Fatalf("expected at least 2 map completions, got %d", len(items))
	}
	if items[0].Label != "@a" || items[1].Label != "@z" {
		t.Fatalf("map completion order = [%s %s], want [@a @z]", items[0].Label, items[1].Label)
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
