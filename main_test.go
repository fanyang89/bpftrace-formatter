package main

import (
	"strings"
	"testing"
)

func TestFormatterBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "simple probe",
			input: `tracepoint:syscalls:sys_enter_openat { printf("openat: %s\n", str(args.filename)); }`,
			expected: `tracepoint:syscalls:sys_enter_openat {
    printf("openat: %s\n", str(args.filename));
}`,
		},
		{
			name:  "probe with predicate",
			input: `tracepoint:syscalls:sys_enter_openat /pid == 1234/ { printf("openat: %s\n", str(args.filename)); }`,
			expected: `tracepoint:syscalls:sys_enter_openat /pid == 1234/ {
    printf("openat: %s\n", str(args.filename));
}`,
		},
	}

	formatter := NewFormatter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.Format(tt.input)
			// Normalize newlines for comparison
			result = normalizeNewlines(result)
			expected := normalizeNewlines(tt.expected)
			if result != expected {
				t.Errorf("Format() mismatch:\nGot:\n%s\n\nWant:\n%s", result, expected)
			}
		})
	}
}

func TestFormatterConfiguration(t *testing.T) {
	formatter := NewFormatter()

	// Test configuration methods
	formatter.SetIndentSize(2)
	formatter.SetUseSpaces(false)
	formatter.SetProbeSpacing(1)

	if formatter.indentSize != 2 {
		t.Errorf("Expected indent size 2, got %d", formatter.indentSize)
	}

	if formatter.useSpaces != false {
		t.Errorf("Expected useSpaces false, got %v", formatter.useSpaces)
	}

	if formatter.probeSpacing != 1 {
		t.Errorf("Expected probe spacing 1, got %d", formatter.probeSpacing)
	}
}

func normalizeNewlines(s string) string {
	// Replace Windows-style newlines with Unix-style
	s = strings.ReplaceAll(s, "\r\n", "\n")
	// Remove trailing newlines
	s = strings.TrimRight(s, "\n")
	return s
}
