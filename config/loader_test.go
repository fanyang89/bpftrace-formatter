package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestConfigLoader_findConfigFile_ExplicitPathWins(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "custom.json")
	writeFile(t, cfgPath, "{}")

	cl := NewConfigLoader()
	cl.configFile = cfgPath
	got := cl.findConfigFile()
	if got != cfgPath {
		t.Fatalf("expected %q, got %q", cfgPath, got)
	}
}

func TestConfigLoader_findConfigFile_ExplicitMissingDoesNotFallback(t *testing.T) {
	tmp := t.TempDir()
	// Even if .btfmt.json exists in cwd, specifying a missing config file stops search.
	writeFile(t, filepath.Join(tmp, ".btfmt.json"), "{}")

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	cl := NewConfigLoader()
	cl.configFile = filepath.Join(tmp, "does-not-exist.json")
	got := cl.findConfigFile()
	if got != "" {
		t.Fatalf("expected empty path, got %q", got)
	}
}

func TestConfigLoader_searchUpwards_NearestAncestorWins(t *testing.T) {
	tmp := t.TempDir()
	a := filepath.Join(tmp, "a")
	b := filepath.Join(a, "b")
	c := filepath.Join(b, "c")
	if err := os.MkdirAll(c, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfgA := filepath.Join(a, ".btfmt.json")
	cfgB := filepath.Join(b, ".btfmt.json")
	writeFile(t, cfgA, "{}")
	writeFile(t, cfgB, "{}")

	cl := NewConfigLoader()
	got := cl.searchUpwards(c, ".btfmt.json")
	if got != cfgB {
		t.Fatalf("expected %q, got %q", cfgB, got)
	}
}

func TestConfigLoader_findConfigFile_FallsBackToHome(t *testing.T) {
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	cwd := filepath.Join(tmp, "cwd")
	if err := os.MkdirAll(home, 0755); err != nil {
		t.Fatalf("mkdir home: %v", err)
	}
	if err := os.MkdirAll(cwd, 0755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}
	t.Setenv("HOME", home)

	homeCfg := filepath.Join(home, ".btfmt.json")
	writeFile(t, homeCfg, "{}")

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWd) })
	if err := os.Chdir(cwd); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	cl := NewConfigLoader()
	got := cl.findConfigFile()
	if got != homeCfg {
		t.Fatalf("expected %q, got %q", homeCfg, got)
	}
}

func TestConfigLoader_LoadConfig_InvalidJSONIncludesPath(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "bad.json")
	writeFile(t, cfgPath, "{")

	cl := NewConfigLoader()
	cl.configFile = cfgPath
	_, err := cl.LoadConfig()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "failed to load config from "+cfgPath) {
		t.Fatalf("expected error to include path; got: %v", err)
	}
}
