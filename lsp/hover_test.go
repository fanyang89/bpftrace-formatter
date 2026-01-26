package lsp

import "testing"

func TestSnippetForContext_UnicodeOffsets(t *testing.T) {
	input := "// 注释\nkprobe:sys_clone { @x = count(); }\n"
	result := ParseDocument(input)
	if result == nil || result.Tree == nil || result.Tree.Content() == nil {
		t.Fatalf("expected parse tree")
	}

	probes := result.Tree.Content().AllProbe()
	if len(probes) == 0 {
		t.Fatalf("expected probe")
	}
	probeList := probes[0].Probe_list()
	if probeList == nil {
		t.Fatalf("expected probe list")
	}

	got := snippetForContext(input, probeList)
	if got != "kprobe:sys_clone" {
		t.Fatalf("snippetForContext = %q, want %q", got, "kprobe:sys_clone")
	}
}
