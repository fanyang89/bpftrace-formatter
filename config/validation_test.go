package config

import (
	"strings"
	"testing"
)

func TestValidate_InvalidBraceStyle(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Blocks.BraceStyle = "invalid"

	errors := cfg.Validate()
	if len(errors) == 0 {
		t.Fatal("expected validation error for invalid brace_style")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err.Error(), "brace_style") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error message to mention brace_style")
	}
}

func TestValidate_IndentSizeTooSmall(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Indent.Size = 0

	errors := cfg.Validate()
	if len(errors) == 0 {
		t.Fatal("expected validation error for indent.size = 0")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err.Error(), "indent.size") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error message to mention indent.size")
	}
}

func TestValidate_IndentSizeTooLarge(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Indent.Size = 17

	errors := cfg.Validate()
	if len(errors) == 0 {
		t.Fatal("expected validation error for indent.size = 17")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err.Error(), "indent.size") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error message to mention indent.size")
	}
}

func TestValidate_MaxLineLengthTooSmall(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LineBreaks.MaxLineLength = 39

	errors := cfg.Validate()
	if len(errors) == 0 {
		t.Fatal("expected validation error for max_line_length = 39")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err.Error(), "max_line_length") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error message to mention max_line_length")
	}
}

func TestValidate_NegativeEmptyLinesBetweenProbes(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LineBreaks.EmptyLinesBetweenProbes = -1

	errors := cfg.Validate()
	if len(errors) == 0 {
		t.Fatal("expected validation error for negative empty_lines_between_probes")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err.Error(), "empty_lines_between_probes") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error message to mention empty_lines_between_probes")
	}
}

func TestValidate_EmptyLinesBetweenProbesTooLarge(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LineBreaks.EmptyLinesBetweenProbes = 6

	errors := cfg.Validate()
	if len(errors) == 0 {
		t.Fatal("expected validation error for empty_lines_between_probes = 6")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err.Error(), "empty_lines_between_probes") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error message to mention empty_lines_between_probes")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := DefaultConfig()
	errors := cfg.Validate()

	if len(errors) != 0 {
		t.Errorf("expected no validation errors for default config, got: %v", errors)
	}
}

func TestValidate_AllBraceStyles(t *testing.T) {
	for _, style := range ValidBraceStyles {
		cfg := DefaultConfig()
		cfg.Blocks.BraceStyle = style

		errors := cfg.Validate()
		if len(errors) != 0 {
			t.Errorf("brace_style %q should be valid, got errors: %v", style, errors)
		}
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Blocks.BraceStyle = "invalid"
	cfg.Indent.Size = 0
	cfg.LineBreaks.MaxLineLength = 10

	errors := cfg.Validate()
	if len(errors) < 3 {
		t.Errorf("expected at least 3 validation errors, got %d: %v", len(errors), errors)
	}
}

func TestValidate_EmptyLinesAfterShebang(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LineBreaks.EmptyLinesAfterShebang = -1

	errors := cfg.Validate()
	if len(errors) == 0 {
		t.Fatal("expected validation error for negative empty_lines_after_shebang")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err.Error(), "empty_lines_after_shebang") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error message to mention empty_lines_after_shebang")
	}
}
