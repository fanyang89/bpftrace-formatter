package lsp

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/fanyang89/bpftrace-formatter/config"
)

func TestSendFormatResult_ContextCanceledDoesNotBlock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ch := make(chan formatResult)
	done := make(chan struct{})

	go func() {
		sendFormatResult(ctx, ch, formatResult{text: "x"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("sendFormatResult blocked after context cancellation")
	}
}

func TestFormatWithTimeout_TimesOut(t *testing.T) {
	doc := &Document{Text: strings.Repeat("kprobe:sys_clone { @x[pid] = count(); }\n", 2000)}

	_, err := formatWithTimeout(doc, config.DefaultConfig(), -time.Nanosecond)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("unexpected timeout error: %v", err)
	}
}

func TestFormatWithTimeout_RecoversPanic(t *testing.T) {
	doc := &Document{Text: "kprobe:sys_clone { @x = count(); }\n"}

	_, err := formatWithTimeout(doc, nil, time.Second)
	if err == nil {
		t.Fatalf("expected panic error")
	}
	if !strings.Contains(err.Error(), "formatter panic") {
		t.Fatalf("unexpected panic error: %v", err)
	}
}
