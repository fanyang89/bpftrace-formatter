package config

import (
	"os"
	"testing"
)

func writeFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}
