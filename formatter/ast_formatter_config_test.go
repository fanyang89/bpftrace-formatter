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
