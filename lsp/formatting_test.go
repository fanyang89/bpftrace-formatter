package lsp

import (
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fanyang89/bpftrace-formatter/config"
)

func TestFormatWithTimeout_TimesOut(t *testing.T) {
	resetFormatTasks()
	t.Cleanup(resetFormatTasks)
	oldRunFormat := runFormat
	t.Cleanup(func() { runFormat = oldRunFormat })

	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	runFormat = func(_ *Document, _ *config.Config) formatResult {
		<-release
		return formatResult{text: "ok"}
	}

	doc := &Document{Version: 1, Text: strings.Repeat("kprobe:sys_clone { @x[pid] = count(); }\n", 8)}

	_, err := formatWithTimeout("file:///timeout.bt", doc, config.DefaultConfig(), 10*time.Millisecond)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("unexpected timeout error: %v", err)
	}
}

func TestFormatWithTimeout_RecoversPanic(t *testing.T) {
	resetFormatTasks()
	t.Cleanup(resetFormatTasks)
	oldRunFormat := runFormat
	t.Cleanup(func() { runFormat = oldRunFormat })

	runFormat = func(_ *Document, _ *config.Config) formatResult {
		panic("boom")
	}

	doc := &Document{Version: 1, Text: "kprobe:sys_clone { @x = count(); }\n"}
	_, err := formatWithTimeout("file:///panic.bt", doc, config.DefaultConfig(), time.Second)
	if err == nil {
		t.Fatalf("expected panic error")
	}
	if !strings.Contains(err.Error(), "formatter panic") {
		t.Fatalf("unexpected panic error: %v", err)
	}
}

func TestFormatWithTimeout_ReusesInFlightTaskForSameVersion(t *testing.T) {
	resetFormatTasks()
	t.Cleanup(resetFormatTasks)
	oldRunFormat := runFormat
	t.Cleanup(func() { runFormat = oldRunFormat })

	var calls atomic.Int32
	release := make(chan struct{})
	finished := make(chan struct{}, 1)
	runFormat = func(_ *Document, _ *config.Config) formatResult {
		calls.Add(1)
		<-release
		finished <- struct{}{}
		return formatResult{text: "ok"}
	}

	doc := &Document{Version: 3, Text: "BEGIN { exit(); }\n"}
	_, err := formatWithTimeout("file:///same.bt", doc, config.DefaultConfig(), 10*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("first call expected timeout, got: %v", err)
	}

	_, err = formatWithTimeout("file:///same.bt", doc, config.DefaultConfig(), 10*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("second call expected timeout, got: %v", err)
	}

	if got := calls.Load(); got != 1 {
		t.Fatalf("runFormat call count = %d, want 1", got)
	}

	close(release)
	<-finished
}

func TestFormatWithTimeout_RejectsNewVersionWhenPreviousStillRunning(t *testing.T) {
	resetFormatTasks()
	t.Cleanup(resetFormatTasks)
	oldRunFormat := runFormat
	t.Cleanup(func() { runFormat = oldRunFormat })

	var calls atomic.Int32
	release := make(chan struct{})
	finished := make(chan struct{}, 1)
	runFormat = func(_ *Document, _ *config.Config) formatResult {
		calls.Add(1)
		<-release
		finished <- struct{}{}
		return formatResult{text: "ok"}
	}

	uri := "file:///versioned.bt"
	docV1 := &Document{Version: 1, Text: "BEGIN { exit(); }\n"}
	_, err := formatWithTimeout(uri, docV1, config.DefaultConfig(), 10*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("first call expected timeout, got: %v", err)
	}

	docV2 := &Document{Version: 2, Text: "BEGIN { printf(\"x\"); }\n"}
	_, err = formatWithTimeout(uri, docV2, config.DefaultConfig(), time.Second)
	if err == nil {
		t.Fatalf("expected in-progress error")
	}
	if !errors.Is(err, errFormattingInProgress) {
		t.Fatalf("expected errFormattingInProgress, got: %v", err)
	}

	if got := calls.Load(); got != 1 {
		t.Fatalf("runFormat call count = %d, want 1", got)
	}

	close(release)
	<-finished
}
