package formatter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanyang89/bpftrace-formatter/config"
)

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	// Normalize line endings.
	return strings.ReplaceAll(string(b), "\r\n", "\n")
}

func trimTrailingNewlines(s string) string {
	return strings.TrimRight(s, "\n")
}

func TestASTFormatter_GoldenFiles_DefaultConfig(t *testing.T) {
	fixturesDir := filepath.Join("..", "testdata")
	goldenDir := filepath.Join("..", "testdata", "golden")

	tests := []struct {
		name       string
		inputPath  string
		goldenPath string
	}{
		{
			name:       "test_input",
			inputPath:  filepath.Join(fixturesDir, "test_input.bt"),
			goldenPath: filepath.Join(goldenDir, "test_input.bt"),
		},
		{
			name:       "test_script",
			inputPath:  filepath.Join(fixturesDir, "test_script.bt"),
			goldenPath: filepath.Join(goldenDir, "test_script.bt"),
		},
		{
			name:       "test_operators",
			inputPath:  filepath.Join(fixturesDir, "test_operators.bt"),
			goldenPath: filepath.Join(goldenDir, "test_operators.bt"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := readFile(t, tc.inputPath)
			want := trimTrailingNewlines(readFile(t, tc.goldenPath))

			f := NewASTFormatter(config.DefaultConfig())
			got, err := f.Format(input)
			if err != nil {
				t.Fatalf("Format returned error: %v", err)
			}

			if got != want {
				t.Fatalf("golden mismatch\n--- got ---\n%s\n--- want ---\n%s\n", got, want)
			}
		})
	}
}
