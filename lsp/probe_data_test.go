package lsp

import (
	"strings"
	"testing"
)

func TestExtractProbeTarget(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		wantProbeType string
		wantPrefix    string
	}{
		{
			name:          "kprobe with prefix",
			line:          "kprobe:vfs_r",
			wantProbeType: "kprobe",
			wantPrefix:    "vfs_r",
		},
		{
			name:          "kprobe empty prefix",
			line:          "kprobe:",
			wantProbeType: "kprobe",
			wantPrefix:    "",
		},
		{
			name:          "kretprobe with prefix",
			line:          "kretprobe:do_e",
			wantProbeType: "kretprobe",
			wantPrefix:    "do_e",
		},
		{
			name:          "tracepoint syscall name",
			line:          "tracepoint:syscalls:sys_enter_rea",
			wantProbeType: "tracepoint",
			wantPrefix:    "sys_enter_rea",
		},
		{
			name:          "tracepoint category",
			line:          "tracepoint:sche",
			wantProbeType: "tracepoint",
			wantPrefix:    "sche",
		},
		{
			name:          "software event",
			line:          "software:page",
			wantProbeType: "software",
			wantPrefix:    "page",
		},
		{
			name:          "hardware event",
			line:          "hardware:cache",
			wantProbeType: "hardware",
			wantPrefix:    "cache",
		},
		{
			name:          "no probe type",
			line:          "vfs_read",
			wantProbeType: "",
			wantPrefix:    "",
		},
		{
			name:          "probe with space should not match",
			line:          "kprobe:vfs_read {",
			wantProbeType: "",
			wantPrefix:    "",
		},
		{
			name:          "probe in string should not match",
			line:          `"kprobe:vfs_r"`,
			wantProbeType: "",
			wantPrefix:    "",
		},
		{
			name:          "uprobe with prefix",
			line:          "uprobe:/bin/bash:read",
			wantProbeType: "uprobe",
			wantPrefix:    "read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			probeType, prefix := extractProbeTarget(tt.line)
			if probeType != tt.wantProbeType {
				t.Errorf("extractProbeTarget(%q) probeType = %q, want %q", tt.line, probeType, tt.wantProbeType)
			}
			if prefix != tt.wantPrefix {
				t.Errorf("extractProbeTarget(%q) prefix = %q, want %q", tt.line, prefix, tt.wantPrefix)
			}
		})
	}
}

func TestProbeTargetCompletions_Kprobe(t *testing.T) {
	items := probeTargetCompletions("kprobe", "vfs_r")
	if len(items) == 0 {
		t.Error("expected kprobe completions for vfs_r")
	}

	found := make(map[string]bool)
	for _, item := range items {
		found[item.Label] = true
		if !strings.HasPrefix(item.Label, "vfs_r") {
			t.Errorf("unexpected completion %q", item.Label)
		}
	}

	expected := []string{"vfs_read", "vfs_readv", "vfs_readlink", "vfs_rename"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected completion %q not found", e)
		}
	}
}

func TestProbeTargetCompletions_Tracepoint(t *testing.T) {
	items := probeTargetCompletions("tracepoint", "sys_enter_rea")
	if len(items) == 0 {
		t.Error("expected tracepoint completions for sys_enter_rea")
	}

	found := make(map[string]bool)
	for _, item := range items {
		found[item.Label] = true
		if !strings.HasPrefix(item.Label, "sys_enter_rea") {
			t.Errorf("unexpected completion %q", item.Label)
		}
	}

	expected := []string{"sys_enter_read", "sys_enter_readv"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected completion %q not found", e)
		}
	}
}

func TestProbeTargetCompletions_Software(t *testing.T) {
	items := probeTargetCompletions("software", "page")
	if len(items) == 0 {
		t.Error("expected software event completions for page")
	}

	found := make(map[string]bool)
	for _, item := range items {
		found[item.Label] = true
	}

	expected := []string{"page-faults", "page-faults-user"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected completion %q not found", e)
		}
	}
}

func TestProbeTargetCompletions_Hardware(t *testing.T) {
	items := probeTargetCompletions("hardware", "cache")
	if len(items) == 0 {
		t.Error("expected hardware event completions for cache")
	}

	found := make(map[string]bool)
	for _, item := range items {
		found[item.Label] = true
	}

	if !found["cache-references"] || !found["cache-misses"] {
		t.Error("expected cache-references and cache-misses completions")
	}
}

func TestProbeTargetCompletions_EmptyPrefix(t *testing.T) {
	items := probeTargetCompletions("kprobe", "")
	if len(items) == 0 {
		t.Error("expected all kprobe completions for empty prefix")
	}

	if len(items) < 100 {
		t.Errorf("expected at least 100 kprobe completions, got %d", len(items))
	}
}

func TestProbeTargetCompletions_InvalidType(t *testing.T) {
	items := probeTargetCompletions("invalid", "test")
	if items != nil {
		t.Error("expected nil for invalid probe type")
	}
}

func TestGetProbeDefinitions(t *testing.T) {
	tests := []struct {
		probeType string
		wantCount int
	}{
		{"kprobe", 100},
		{"kretprobe", 100},
		{"tracepoint", 200},
		{"software", 10},
		{"hardware", 30},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.probeType, func(t *testing.T) {
			probes := getProbeDefinitions(tt.probeType)
			if tt.wantCount == 0 {
				if probes != nil {
					t.Errorf("expected nil for probe type %q", tt.probeType)
				}
				return
			}

			if len(probes) < tt.wantCount {
				t.Errorf("getProbeDefinitions(%q) returned %d probes, want at least %d",
					tt.probeType, len(probes), tt.wantCount)
			}

			for _, probe := range probes {
				if probe.Name == "" {
					t.Error("probe name should not be empty")
				}
				if probe.Description == "" {
					t.Errorf("probe %q should have description", probe.Name)
				}
			}
		})
	}
}

func TestProbeDataIntegrity(t *testing.T) {
	seen := make(map[string]string)

	checkProbes := func(probes []ProbeDefinition, source string) {
		for _, probe := range probes {
			if existing, ok := seen[probe.Name]; ok {
				if existing != source {
					t.Errorf("duplicate probe name %q in %s (already in %s)",
						probe.Name, source, existing)
				}
			}
			seen[probe.Name] = source
		}
	}

	checkProbes(commonKprobes, "kprobes")
	checkProbes(commonTracepoints, "tracepoints")
	checkProbes(softwareEvents, "software")
	checkProbes(hardwareEvents, "hardware")
}

func TestProbeCompletionHasDocumentation(t *testing.T) {
	items := probeTargetCompletions("kprobe", "vfs_read")
	if len(items) == 0 {
		t.Fatal("expected vfs_read completion")
	}

	item := items[0]
	if item.Documentation == "" {
		t.Error("completion item should have documentation")
	}
	if item.Detail == nil || *item.Detail == "" {
		t.Error("completion item should have detail")
	}
	if item.SortText == nil || *item.SortText == "" {
		t.Error("completion item should have sort text")
	}
}

func TestProbeCompletionSorted(t *testing.T) {
	items := probeTargetCompletions("kprobe", "vfs_")
	if len(items) < 2 {
		t.Skip("need at least 2 items to test sorting")
	}

	for i := 1; i < len(items); i++ {
		prev := items[i-1].Label
		curr := items[i].Label
		if prev > curr {
			t.Errorf("completions not sorted: %q > %q", prev, curr)
		}
	}
}
