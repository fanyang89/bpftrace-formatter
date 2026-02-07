package formatter

import "testing"

func TestParseBpftrace_ValidInput(t *testing.T) {
	input := "kprobe:sys_clone { @x = count(); }\n"

	tree, err := ParseBpftrace(input)
	if err != nil {
		t.Fatalf("ParseBpftrace returned error: %v", err)
	}
	if tree == nil {
		t.Fatalf("expected parse tree")
	}
}

func TestParseBpftrace_InvalidInput(t *testing.T) {
	input := "kprobe:sys_clone { @x = count( }\n"

	tree, err := ParseBpftrace(input)
	if err == nil {
		t.Fatalf("expected parse error")
	}
	if tree != nil {
		t.Fatalf("expected nil parse tree on parse error")
	}
}
