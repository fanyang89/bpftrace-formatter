package formatter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanyang89/bpftrace-formatter/config"
)

func TestASTFormatter_FormatsBpftraceToolsTree(t *testing.T) {
	root := filepath.Join("..", "bpftrace", "tools")
	if !dirExists(root) {
		t.Skip("bpftrace/tools not found")
	}

	files, err := collectBtFiles(root)
	if err != nil {
		t.Fatalf("collect %q: %v", root, err)
	}
	if len(files) == 0 {
		t.Fatalf("no .bt files found under %q", root)
	}

	f := NewASTFormatter(config.DefaultConfig())
	for _, path := range files {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		out, err := f.Format(string(b))
		if err != nil {
			t.Fatalf("format %s: %v", path, err)
		}
		if out == "" {
			t.Fatalf("format %s: empty output", path)
		}
	}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func collectBtFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".bt") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
