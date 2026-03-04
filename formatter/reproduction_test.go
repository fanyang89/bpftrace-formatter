package formatter

import (
	"testing"

	"github.com/fanyang89/bpftrace-formatter/config"
)

func TestReproduction_IfElseNextLine(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Blocks.BraceStyle = "next_line"

	input := "BEGIN{if(1){exit();}else{exit();}}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	// Desired: else on new line, no trailing semicolon for the if block
	want := "BEGIN\n" +
		"{\n" +
		"    if (1)\n" +
		"    {\n" +
		"        exit();\n" +
		"    }\n" +
		"    else\n" +
		"    {\n" +
		"        exit();\n" +
		"    }\n" +
		"}"

	if got != want {
		t.Errorf("IF/ELSE failed\ngot:\n%q\nwant:\n%q", got, want)
	}
}

func TestReproduction_RedundantSemicolon(t *testing.T) {
	cfg := config.DefaultConfig()
	input := "BEGIN{if(1){exit();} exit();}"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    if (1)\n" +
		"    {\n" +
		"        exit();\n" +
		"    }\n" +
		"    exit();\n" +
		"}"
	if got != want {
		t.Errorf("Semicolon failed\ngot:\n%q\nwant:\n%q", got, want)
	}
}

func TestReproduction_CommentSpacing(t *testing.T) {
	cfg := config.DefaultConfig()
	input := "BEGIN{exit();}//comment\n//nospace"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "BEGIN\n" +
		"{\n" +
		"    exit();\n" +
		"}\n" +
		"// comment\n" +
		"// nospace"

	// Note: Current implementation might not even put // comment on a new line if it's inline but the want here is after improvement
	if got != want {
		t.Errorf("Comment spacing failed\ngot:\n%q\nwant:\n%q", got, want)
	}
}

func TestReproduction_InsidePredicates(t *testing.T) {
	cfg := config.DefaultConfig()
	// Default is likely true (with spaces), we'll test both when implemented
	// For now, let's assume we want to support false (no spaces)
	cfg.Spacing.InsidePredicates = false

	input := "kprobe:vfs_read /pid==1234/ { exit(); }"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "kprobe:vfs_read\n/pid == 1234/\n{\n    exit();\n}"
	if got != want {
		t.Errorf("InsidePredicates failed\ngot:\n%q\nwant:\n%q", got, want)
	}
}

func TestReproduction_NewlineBetweenSpecifiers(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Probes.NewlineBetweenSpecifiers = false

	input := "kprobe:f1,kprobe:f2 { exit(); }"
	got, err := NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "kprobe:f1, kprobe:f2\n{\n    exit();\n}"
	if got != want {
		t.Errorf("NewlineBetweenSpecifiers failed\ngot:\n%q\nwant:\n%q", got, want)
	}
}
