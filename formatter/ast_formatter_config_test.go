package formatter

import (
	"testing"

	"github.com/fanyang89/bpftrace-formatter/config"
)

func TestASTFormatter_BraceStyleSameLine(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Blocks.BraceStyle = "same_line"

	input := "kprobe:sys_execve{exit();}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "kprobe:sys_execve {\n" +
		"    exit();\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_ParenAndCommaSpacing(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Spacing.AroundParentheses = true
	cfg.Spacing.AroundCommas = false

	input := "BEGIN{printf(\"a\",1,2);}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    printf( \"a\",1,2 );\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_BracketAndOperatorSpacing(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Spacing.AroundBrackets = true
	cfg.Spacing.AroundOperators = false

	// Use the grammar's "@" map form so '[' and ']' are real tokens
	// (MAP_NAME can lex a single token that includes brackets).
	input := "BEGIN{@[pid]=count();}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    @[ pid ]=count();\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_IndentWithTabs(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Indent.UseSpaces = false

	input := "BEGIN{printf(\"x\");}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"\tprintf(\"x\");\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_EmptyBlock(t *testing.T) {
	got, err := NewASTFormatter(config.DefaultConfig()).Format("BEGIN{}")
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_BraceStyleGNU(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Blocks.BraceStyle = "gnu"

	input := "BEGIN{exit();}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"    {\n" +
		"        exit();\n" +
		"    }"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_IndentStatementsDisabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Blocks.IndentStatements = false

	input := "BEGIN{printf(\"x\");}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"printf(\"x\");\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_AlignPredicates(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Probes.AlignPredicates = true

	input := "tracepoint:syscalls:sys_enter_openat/pid==1234/{exit();}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "tracepoint:syscalls:sys_enter_openat / pid == 1234 /\n" +
		"{\n" +
		"    exit();\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_AlignPredicatesNoWrap(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Probes.AlignPredicates = true
	cfg.LineBreaks.MaxLineLength = 10
	cfg.LineBreaks.BreakLongStatements = true

	input := "tracepoint:syscalls:sys_enter_openat/pid/{exit();}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "tracepoint:syscalls:sys_enter_openat / pid /\n" +
		"{\n" +
		"    exit();\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_PreserveInlineComments(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Comments.PreserveInline = true

	input := "BEGIN{printf(\"x\"); // hello\n}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    printf(\"x\"); // hello\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_PreserveInlineCommentNoWrap(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Comments.PreserveInline = true
	cfg.LineBreaks.MaxLineLength = 10
	cfg.LineBreaks.BreakLongStatements = true

	input := "BEGIN{printf(\"long\"); // note\n}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    printf(\"long\"); // note\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_CommentIndentLevel(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Comments.IndentLevel = 1

	input := "BEGIN{exit();}\n// hello"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    exit();\n" +
		"}\n" +
		"    // hello"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_LeadingCommentIndentPreserved(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Comments.IndentLevel = 1

	input := "// hello\nBEGIN{exit();}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "    // hello\n" +
		"BEGIN\n" +
		"{\n" +
		"    exit();\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_CommentBetweenProbes(t *testing.T) {
	cfg := config.DefaultConfig()

	input := "BEGIN{exit();}\n// next probe\nEND{exit();}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    exit();\n" +
		"}\n" +
		"\n" +
		"// next probe\n" +
		"END\n" +
		"{\n" +
		"    exit();\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_BreakLongStatements(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.LineBreaks.MaxLineLength = 20
	cfg.LineBreaks.BreakLongStatements = true

	input := "BEGIN{printf(\"x\",1+2+3+4);}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    printf(\"x\", 1 +\n" +
		"    2 + 3 + 4);\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_BreakLongStatements_TokenAwareWrap(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.LineBreaks.MaxLineLength = 18
	cfg.LineBreaks.BreakLongStatements = true

	input := "BEGIN{printf(\"x\",1+2345);}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    printf(\"x\", 1\n" +
		"    + 2345);\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_IfElseBraceStyleNextLine(t *testing.T) {
	cfg := config.DefaultConfig()

	input := "BEGIN{if(1){exit();}else{exit();}}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    if (1)\n" +
		"    {\n" +
		"        exit();\n" +
		"    } else\n" +
		"    {\n" +
		"        exit();\n" +
		"    };\n" +
		"}"

	if got != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
	}
}

func TestASTFormatter_IndentSize(t *testing.T) {
	tests := []struct {
		name       string
		indentSize int
		want       string
	}{
		{
			name:       "indent size 2",
			indentSize: 2,
			want: "BEGIN\n" +
				"{\n" +
				"  exit();\n" +
				"}",
		},
		{
			name:       "indent size 8",
			indentSize: 8,
			want: "BEGIN\n" +
				"{\n" +
				"        exit();\n" +
				"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.Indent.Size = tt.indentSize
			cfg.Indent.UseSpaces = true

			input := "BEGIN{exit();}"
			got, err := NewASTFormatter(cfg).Format(input)
			if err != nil {
				t.Fatalf("Format returned error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, tt.want)
			}
		})
	}
}

func TestASTFormatter_BeforeBlockStartSpacing(t *testing.T) {
	tests := []struct {
		name             string
		beforeBlockStart bool
		want             string
	}{
		{
			name:             "space before block start enabled",
			beforeBlockStart: true,
			want: "BEGIN {\n" +
				"    exit();\n" +
				"}",
		},
		{
			name:             "space before block start disabled",
			beforeBlockStart: false,
			want: "BEGIN{\n" +
				"    exit();\n" +
				"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.Blocks.BraceStyle = "same_line"
			cfg.Spacing.BeforeBlockStart = tt.beforeBlockStart

			input := "BEGIN{exit();}"
			got, err := NewASTFormatter(cfg).Format(input)
			if err != nil {
				t.Fatalf("Format returned error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, tt.want)
			}
		})
	}
}

func TestASTFormatter_AfterKeywordsSpacing(t *testing.T) {
	tests := []struct {
		name          string
		afterKeywords bool
		want          string
	}{
		{
			name:          "space after keywords enabled",
			afterKeywords: true,
			want: "BEGIN\n" +
				"{\n" +
				"    if (1)\n" +
				"    {\n" +
				"        exit();\n" +
				"    };\n" +
				"}",
		},
		{
			name:          "space after keywords disabled",
			afterKeywords: false,
			want: "BEGIN\n" +
				"{\n" +
				"    if(1)\n" +
				"    {\n" +
				"        exit();\n" +
				"    };\n" +
				"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.Spacing.AfterKeywords = tt.afterKeywords

			input := "BEGIN{if(1){exit();}}"
			got, err := NewASTFormatter(cfg).Format(input)
			if err != nil {
				t.Fatalf("Format returned error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, tt.want)
			}
		})
	}
}

func TestASTFormatter_EmptyLinesBetweenProbes(t *testing.T) {
	tests := []struct {
		name       string
		emptyLines int
		want       string
	}{
		{
			name:       "zero empty lines between probes",
			emptyLines: 0,
			want: "BEGIN\n" +
				"{\n" +
				"    exit();\n" +
				"}\n" +
				"END\n" +
				"{\n" +
				"    exit();\n" +
				"}",
		},
		{
			name:       "one empty line between probes",
			emptyLines: 1,
			want: "BEGIN\n" +
				"{\n" +
				"    exit();\n" +
				"}\n" +
				"\n" +
				"END\n" +
				"{\n" +
				"    exit();\n" +
				"}",
		},
		{
			name:       "two empty lines between probes",
			emptyLines: 2,
			want: "BEGIN\n" +
				"{\n" +
				"    exit();\n" +
				"}\n" +
				"\n" +
				"\n" +
				"END\n" +
				"{\n" +
				"    exit();\n" +
				"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.LineBreaks.EmptyLinesBetweenProbes = tt.emptyLines

			input := "BEGIN{exit();}END{exit();}"
			got, err := NewASTFormatter(cfg).Format(input)
			if err != nil {
				t.Fatalf("Format returned error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, tt.want)
			}
		})
	}
}

func TestASTFormatter_EmptyLinesAfterShebang(t *testing.T) {
	tests := []struct {
		name       string
		emptyLines int
		want       string
	}{
		{
			name:       "zero empty lines after shebang",
			emptyLines: 0,
			want: "#!/usr/bin/env bpftrace\n" +
				"BEGIN\n" +
				"{\n" +
				"    exit();\n" +
				"}",
		},
		{
			name:       "one empty line after shebang",
			emptyLines: 1,
			want: "#!/usr/bin/env bpftrace\n" +
				"\n" +
				"BEGIN\n" +
				"{\n" +
				"    exit();\n" +
				"}",
		},
		{
			name:       "two empty lines after shebang",
			emptyLines: 2,
			want: "#!/usr/bin/env bpftrace\n" +
				"\n" +
				"\n" +
				"BEGIN\n" +
				"{\n" +
				"    exit();\n" +
				"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.LineBreaks.EmptyLinesAfterShebang = tt.emptyLines

			input := "#!/usr/bin/env bpftrace\nBEGIN{exit();}"
			got, err := NewASTFormatter(cfg).Format(input)
			if err != nil {
				t.Fatalf("Format returned error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", got, tt.want)
			}
		})
	}
}
