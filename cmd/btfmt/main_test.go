package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fanyang89/bpftrace-formatter/config"
	"github.com/fanyang89/bpftrace-formatter/formatter"
)

func TestProcessFile_WriteToFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "in.bt")

	input := "BEGIN{printf(\"x\",1);}"
	if err := os.WriteFile(path, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	cfg := config.DefaultConfig()
	if err := processFile(path, cfg, true, false); err != nil {
		t.Fatalf("processFile returned error: %v", err)
	}

	gotBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	want, err := formatter.NewASTFormatter(cfg).Format(input)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	if string(gotBytes) != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", string(gotBytes), want)
	}
}
