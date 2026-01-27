package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	var stdout, stderr bytes.Buffer
	if err := processFile(path, cfg, true, false, &stdout, &stderr); err != nil {
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
	if !strings.HasSuffix(want, "\n") {
		want += "\n"
	}

	if string(gotBytes) != want {
		t.Fatalf("unexpected output\n--- got ---\n%s\n--- want ---\n%s\n", string(gotBytes), want)
	}
}

func TestProcessFile_PreservesPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file mode preservation is not reliable on Windows")
	}

	tmp := t.TempDir()
	path := filepath.Join(tmp, "perm.bt")

	input := "BEGIN{printf(\"x\",1);}"
	if err := os.WriteFile(path, []byte(input), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if err := os.Chmod(path, 0o751); err != nil {
		t.Fatalf("chmod input: %v", err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat input: %v", err)
	}

	cfg := config.DefaultConfig()
	var stdout, stderr bytes.Buffer
	if err := processFile(path, cfg, true, false, &stdout, &stderr); err != nil {
		t.Fatalf("processFile returned error: %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat output: %v", err)
	}
	if after.Mode().Perm() != before.Mode().Perm() {
		t.Fatalf("mode = %v, want %v", after.Mode().Perm(), before.Mode().Perm())
	}
}

func TestProcessFile_ReadError(t *testing.T) {
	cfg := config.DefaultConfig()
	var stdout, stderr bytes.Buffer
	err := processFile("/nonexistent/file.bt", cfg, false, false, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "reading file") {
		t.Errorf("expected 'reading file' error, got: %v", err)
	}
}

func TestProcessFile_FormatError(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "invalid.bt")

	// Invalid bpftrace syntax
	input := "{ { { invalid syntax"
	if err := os.WriteFile(path, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	cfg := config.DefaultConfig()
	var stdout, stderr bytes.Buffer
	err := processFile(path, cfg, false, false, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for invalid syntax")
	}
	if !strings.Contains(err.Error(), "formatting") {
		t.Errorf("expected 'formatting' error, got: %v", err)
	}
}

func TestProcessFile_WriteToStdout(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.bt")

	input := "BEGIN{printf(\"test\");}"
	if err := os.WriteFile(path, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	cfg := config.DefaultConfig()
	var stdout, stderr bytes.Buffer
	if err := processFile(path, cfg, false, false, &stdout, &stderr); err != nil {
		t.Fatalf("processFile returned error: %v", err)
	}

	if stdout.Len() == 0 {
		t.Error("expected output to stdout")
	}
	if !strings.Contains(stdout.String(), "BEGIN") {
		t.Error("stdout should contain formatted output")
	}
}

func TestProcessFile_VerboseMode(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.bt")

	input := "BEGIN{printf(\"test\");}"
	if err := os.WriteFile(path, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	cfg := config.DefaultConfig()
	var stdout, stderr bytes.Buffer
	if err := processFile(path, cfg, false, true, &stdout, &stderr); err != nil {
		t.Fatalf("processFile returned error: %v", err)
	}

	if !strings.Contains(stderr.String(), "Processing:") {
		t.Error("verbose mode should print processing message to stderr")
	}
}

func TestRun_NoInputFiles(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"btfmt"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error when no input files specified")
	}
	if !strings.Contains(err.Error(), "no input files") {
		t.Errorf("expected 'no input files' error, got: %v", err)
	}
}

func TestRun_GenerateConfig(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "test-config.json")

	var stdout, stderr bytes.Buffer
	err := run([]string{"btfmt", "-generate-config", "-config-output", configPath}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	if !strings.Contains(stdout.String(), "Generated default configuration") {
		t.Error("expected success message in stdout")
	}
}

func TestRun_HelpFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"btfmt", "-help"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if !strings.Contains(stdout.String(), "Usage:") {
		t.Error("help output should contain Usage:")
	}
	if !strings.Contains(stdout.String(), "-config") {
		t.Error("help output should mention -config flag")
	}
}

func TestRun_MultipleFiles(t *testing.T) {
	tmp := t.TempDir()

	// Create multiple test files
	files := []string{
		filepath.Join(tmp, "file1.bt"),
		filepath.Join(tmp, "file2.bt"),
		filepath.Join(tmp, "file3.bt"),
	}

	for i, path := range files {
		input := "BEGIN{printf(\"file%d\");}"
		input = strings.Replace(input, "%d", string(rune('1'+i)), 1)
		if err := os.WriteFile(path, []byte(input), 0644); err != nil {
			t.Fatalf("write input %d: %v", i, err)
		}
	}

	var stdout, stderr bytes.Buffer
	args := append([]string{"btfmt", "-w"}, files...)
	err := run(args, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	// Verify all files were processed
	for _, path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		if !strings.Contains(string(content), "BEGIN") {
			t.Errorf("file %s was not formatted correctly", path)
		}
	}
}

func TestRun_ConfigLoading(t *testing.T) {
	tmp := t.TempDir()

	// Create a config file with non-default settings
	configPath := filepath.Join(tmp, ".btfmt.json")
	configContent := `{
		"indent": {
			"size": 2,
			"use_spaces": true
		},
		"blocks": {
			"brace_style": "same_line"
		}
	}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Create a test file
	testPath := filepath.Join(tmp, "test.bt")
	input := "BEGIN{printf(\"test\");}"
	if err := os.WriteFile(testPath, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := run([]string{"btfmt", "-config", configPath, testPath}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	// With same_line brace style, the output should have { on the same line
	output := stdout.String()
	if !strings.Contains(output, "BEGIN {") {
		t.Errorf("expected same_line brace style, got: %s", output)
	}
}

func TestRun_InvalidConfig(t *testing.T) {
	tmp := t.TempDir()

	// Create a config file with invalid brace_style
	configPath := filepath.Join(tmp, ".btfmt.json")
	configContent := `{
		"blocks": {
			"brace_style": "invalid_style"
		}
	}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Create a test file
	testPath := filepath.Join(tmp, "test.bt")
	input := "BEGIN{printf(\"test\");}"
	if err := os.WriteFile(testPath, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := run([]string{"btfmt", "-config", configPath, testPath}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
	if !strings.Contains(err.Error(), "brace_style") {
		t.Errorf("expected error mentioning brace_style, got: %v", err)
	}
}

func TestRun_InPlaceFlag(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.bt")

	input := "BEGIN{printf(\"test\",1);}"
	if err := os.WriteFile(path, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := run([]string{"btfmt", "-i", path}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	// File should be modified in place
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	// Should be formatted (different from input)
	if string(content) == input {
		t.Error("file should have been formatted in place")
	}
}

func TestRun_WriteFlag(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.bt")

	input := "BEGIN{printf(\"test\",1);}"
	if err := os.WriteFile(path, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := run([]string{"btfmt", "-w", path}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	// File should be modified
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	// Should be formatted (different from input)
	if string(content) == input {
		t.Error("file should have been written with formatted content")
	}
}

func TestRun_VerboseFlag(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.bt")

	input := "BEGIN{printf(\"test\");}"
	if err := os.WriteFile(path, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := run([]string{"btfmt", "-v", path}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if !strings.Contains(stderr.String(), "Processing:") {
		t.Error("verbose mode should output processing info to stderr")
	}
}

func TestParseFlags_AllOptions(t *testing.T) {
	args := []string{
		"btfmt",
		"-generate-config",
		"-config-output", "custom.json",
		"-config", "my-config.json",
		"-i",
		"-w",
		"-verbose",
		"-help",
		"file1.bt",
		"file2.bt",
	}

	opts, err := parseFlags(args)
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}

	if !opts.generateConfig {
		t.Error("generateConfig should be true")
	}
	if opts.configOutput != "custom.json" {
		t.Errorf("configOutput = %q, want %q", opts.configOutput, "custom.json")
	}
	if opts.configFile != "my-config.json" {
		t.Errorf("configFile = %q, want %q", opts.configFile, "my-config.json")
	}
	if !opts.inPlace {
		t.Error("inPlace should be true")
	}
	if !opts.write {
		t.Error("write should be true")
	}
	if !opts.verbose {
		t.Error("verbose should be true")
	}
	if !opts.help {
		t.Error("help should be true")
	}
	if len(opts.files) != 2 {
		t.Errorf("files count = %d, want 2", len(opts.files))
	}
}

func TestParseFlags_ShortForms(t *testing.T) {
	args := []string{"btfmt", "-c", "config.json", "-v", "file.bt"}

	opts, err := parseFlags(args)
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}

	if opts.configFile != "config.json" {
		t.Errorf("configFile = %q, want %q", opts.configFile, "config.json")
	}
	if !opts.verbose {
		t.Error("verbose should be true (short form -v)")
	}
}

func TestParseFlags_Defaults(t *testing.T) {
	args := []string{"btfmt", "file.bt"}

	opts, err := parseFlags(args)
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}

	if opts.generateConfig {
		t.Error("generateConfig should default to false")
	}
	if opts.configOutput != ".btfmt.json" {
		t.Errorf("configOutput = %q, want %q", opts.configOutput, ".btfmt.json")
	}
	if opts.configFile != "" {
		t.Errorf("configFile = %q, want empty", opts.configFile)
	}
	if opts.inPlace {
		t.Error("inPlace should default to false")
	}
	if opts.write {
		t.Error("write should default to false")
	}
	if opts.verbose {
		t.Error("verbose should default to false")
	}
	if opts.help {
		t.Error("help should default to false")
	}
}

func TestRun_LSPHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"btfmt", "lsp", "-help"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if !strings.Contains(stderr.String(), "Usage:") {
		t.Error("LSP help should contain Usage:")
	}
}
