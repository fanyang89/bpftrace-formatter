package config

import (
	"path/filepath"
	"testing"
)

func TestLoadConfig_PartialOverrideMergesWithDefaults(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "cfg.json")
	writeFile(t, cfgPath, `{"indent":{"size":2}}`)

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.Indent.Size != 2 {
		t.Fatalf("expected indent.size=2, got %d", cfg.Indent.Size)
	}
	if cfg.Indent.UseSpaces != true {
		t.Fatalf("expected indent.use_spaces to remain default true")
	}
	if cfg.Blocks.BraceStyle != "next_line" {
		t.Fatalf("expected brace_style to remain default next_line, got %q", cfg.Blocks.BraceStyle)
	}
}

func TestSaveConfig_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "nested", "out.json")

	cfg := DefaultConfig()
	cfg.Indent.Size = 8
	if err := cfg.SaveConfig(cfgPath); err != nil {
		t.Fatalf("SaveConfig returned error: %v", err)
	}

	loaded, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if loaded.Indent.Size != 8 {
		t.Fatalf("expected indent.size=8, got %d", loaded.Indent.Size)
	}
}
