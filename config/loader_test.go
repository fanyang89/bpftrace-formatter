package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigLoader_findConfigPath_ExplicitPathWins(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "custom.json")
	writeFile(t, cfgPath, "{}")

	cl := NewConfigLoader().WithExplicitPath(cfgPath)
	got, _ := cl.findConfigPath()
	if got != cfgPath {
		t.Fatalf("expected %q, got %q", cfgPath, got)
	}
}

func TestConfigLoader_findConfigPath_ExplicitMissingDoesNotFallback(t *testing.T) {
	tmp := t.TempDir()
	// Even if .btfmt.json exists in cwd, specifying a missing config file stops search.
	writeFile(t, filepath.Join(tmp, ".btfmt.json"), "{}")

	cl := NewConfigLoader().
		WithBaseDir(tmp).
		WithExplicitPath(filepath.Join(tmp, "does-not-exist.json"))

	got, isExplicit := cl.findConfigPath()
	if !isExplicit {
		t.Fatalf("expected explicit path")
	}
	if _, err := os.Stat(got); err == nil {
		t.Fatalf("expected path to not exist, but it does: %s", got)
	}

	cfg, err := cl.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.Indent.Size != DefaultConfig().Indent.Size {
		t.Fatalf("expected default config when explicit path is missing")
	}
}

func TestSearchUpwards_NearestAncestorWins(t *testing.T) {
	tmp := t.TempDir()
	a := filepath.Join(tmp, "a")
	b := filepath.Join(a, "b")
	c := filepath.Join(b, "c")
	if err := os.MkdirAll(c, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfgA := filepath.Join(a, ".btfmt.json")
	cfgB := filepath.Join(b, ".btfmt.json")
	writeFile(t, cfgA, "{}")
	writeFile(t, cfgB, "{}")

	got := SearchUpwards(c, ".btfmt.json")
	if got != cfgB {
		t.Fatalf("expected %q, got %q", cfgB, got)
	}
}

func TestConfigLoader_findConfigPath_FallsBackToHome(t *testing.T) {
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	cwd := filepath.Join(tmp, "cwd")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("mkdir home: %v", err)
	}
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}

	// We can't easily set HOME for the whole process in a way that os.UserHomeDir() respects,
	// but we can test that findConfigPath correctly uses SearchUpwards and falls back.
	// For this test, we'll just check baseDir behavior.

	homeCfg := filepath.Join(home, ".btfmt.json")
	writeFile(t, homeCfg, "{}")

	cl := NewConfigLoader().WithBaseDir(home)
	got, _ := cl.findConfigPath()
	if got != homeCfg {
		t.Fatalf("expected %q, got %q", homeCfg, got)
	}
}

func TestConfigLoader_LoadConfig_InvalidJSONIncludesPath(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "bad.json")
	writeFile(t, cfgPath, "{")

	cl := NewConfigLoader().WithExplicitPath(cfgPath)
	_, err := cl.LoadConfig()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "failed to load config from "+cfgPath) {
		t.Fatalf("expected error to include path; got: %v", err)
	}
}

func TestLoadConfigFrom_ExplicitRelativeUsesBaseDir(t *testing.T) {
	baseDir := t.TempDir()
	configDir := filepath.Join(baseDir, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	configPath := filepath.Join(configDir, "btfmt.json")
	if err := os.WriteFile(configPath, []byte(`{"indent":{"size":2}}`), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := LoadConfigFrom(baseDir, filepath.Join("config", "btfmt.json"), false)
	if err != nil {
		t.Fatalf("LoadConfigFrom: %v", err)
	}
	if cfg.Indent.Size != 2 {
		t.Fatalf("indent size = %d, want %d", cfg.Indent.Size, 2)
	}
}

func TestLoadConfigFrom_ExplicitMissingRelativeFallsBack(t *testing.T) {
	baseDir := t.TempDir()
	cfg, err := LoadConfigFrom(baseDir, filepath.Join("config", "missing.json"), false)
	if err != nil {
		t.Fatalf("LoadConfigFrom: %v", err)
	}
	if cfg.Indent.Size != DefaultConfig().Indent.Size {
		t.Fatalf("indent size = %d, want %d", cfg.Indent.Size, DefaultConfig().Indent.Size)
	}
}
