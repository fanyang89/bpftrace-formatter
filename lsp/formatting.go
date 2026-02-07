package lsp

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/config"
	"github.com/fanyang89/bpftrace-formatter/formatter"
)

const formatTimeout = 30 * time.Second

type formatResult struct {
	text string
	err  error
}

type formatTask struct {
	version int32
	done    chan struct{}
	result  formatResult
}

var (
	errFormattingInProgress = errors.New("formatting already in progress")
	formatTasksMu           sync.Mutex
	formatTasks             = map[string]*formatTask{}
)

var runFormat = func(doc *Document, cfg *config.Config) formatResult {
	start := time.Now()
	f := formatter.NewASTFormatter(cfg)
	// Reuse the parse tree from didOpen/didChange when available,
	// avoiding a redundant ANTLR parse.
	if doc.ParseResult != nil && doc.ParseResult.Tree != nil && len(doc.ParseResult.Diagnostics) == 0 {
		text := f.FormatTree(doc.ParseResult.Tree)
		log.Printf("[format] FormatTree took %s", time.Since(start))
		return formatResult{text: text}
	}
	formatted, err := f.Format(doc.Text)
	log.Printf("[format] Format took %s", time.Since(start))
	return formatResult{text: formatted, err: err}
}

func getOrStartFormatTask(uri string, doc *Document, cfg *config.Config) (*formatTask, bool) {
	formatTasksMu.Lock()
	if task, exists := formatTasks[uri]; exists {
		formatTasksMu.Unlock()
		return task, true
	}

	task := &formatTask{
		version: doc.Version,
		done:    make(chan struct{}),
	}
	formatTasks[uri] = task
	formatTasksMu.Unlock()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				task.result = formatResult{err: fmt.Errorf("formatter panic: %v", r)}
			}
			close(task.done)
			formatTasksMu.Lock()
			if current, exists := formatTasks[uri]; exists && current == task {
				delete(formatTasks, uri)
			}
			formatTasksMu.Unlock()
		}()
		task.result = runFormat(doc, cfg)
	}()

	return task, false
}

func resetFormatTasks() {
	formatTasksMu.Lock()
	defer formatTasksMu.Unlock()
	formatTasks = map[string]*formatTask{}
}

func removeFormatTaskIfCurrent(uri string, task *formatTask) {
	formatTasksMu.Lock()
	defer formatTasksMu.Unlock()
	if current, exists := formatTasks[uri]; exists && current == task {
		delete(formatTasks, uri)
	}
}

func formatWithTimeout(uri string, doc *Document, cfg *config.Config, timeout time.Duration) (string, error) {
	task, reused := getOrStartFormatTask(uri, doc, cfg)
	if reused && task.version != doc.Version {
		return "", fmt.Errorf("%w for %s (running version=%d, requested=%d)", errFormattingInProgress, uri, task.version, doc.Version)
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-task.done:
		return task.result.text, task.result.err
	case <-timer.C:
		// If this caller timed out, drop the task entry so follow-up requests
		// can retry instead of being blocked by a potentially stuck worker.
		removeFormatTaskIfCurrent(uri, task)
		return "", fmt.Errorf("formatting timed out after %s", timeout)
	}
}

func formatDocument(uri string) ([]protocol.TextEdit, error) {
	doc, ok := documentStore.Get(uri)
	if !ok || doc == nil {
		return []protocol.TextEdit{}, nil
	}

	cfg := doc.Config
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	formatted, err := formatWithTimeout(uri, doc, cfg, formatTimeout)
	if err != nil {
		return nil, err
	}

	// Ensure trailing newline for consistency with CLI
	if !strings.HasSuffix(formatted, "\n") {
		formatted += "\n"
	}

	edits := []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   EndPosition(doc.Text),
			},
			NewText: formatted,
		},
	}

	return edits, nil
}
