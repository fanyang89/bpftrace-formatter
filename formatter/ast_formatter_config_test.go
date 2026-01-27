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
		"    printf(\"x\", 1 + 2\n" +
		"    + 3 + 4);\n" +
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
