package lsp

import (
	"fmt"
	"net/url"
	"path/filepath"
	"sync"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/config"
)

type Document struct {
	URI         string
	Path        string
	Text        string
	Version     int32
	ParseResult *ParseResult
	Diagnostics []protocol.Diagnostic
	Config      *config.Config
}

type DocumentStore struct {
	mu       sync.RWMutex
	docs     map[string]*Document
	resolver *ConfigResolver
}

func NewDocumentStore(resolver *ConfigResolver) *DocumentStore {
	return &DocumentStore{
		docs:     make(map[string]*Document),
		resolver: resolver,
	}
}

func (s *DocumentStore) Open(uri string, version int32, text string) (*Document, error) {
	return s.upsert(uri, version, text)
}

func (s *DocumentStore) Change(uri string, version int32, text string) (*Document, error) {
	return s.upsert(uri, version, text)
}

func (s *DocumentStore) Close(uri string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.docs, uri)
}

func (s *DocumentStore) Get(uri string) (*Document, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.docs[uri]
	return doc, ok
}

func (s *DocumentStore) upsert(uri string, version int32, text string) (*Document, error) {
	path, err := fileURIToPath(uri)
	if err != nil {
		return nil, err
	}

	parseResult := ParseDocument(text)
	var resolvedConfig *config.Config
	if s.resolver != nil {
		resolvedConfig, err = s.resolver.ResolveForDocument(uri, path)
		if err != nil {
			return nil, err
		}
	} else {
		resolvedConfig = config.DefaultConfig()
	}

	doc := &Document{
		URI:         uri,
		Path:        path,
		Text:        text,
		Version:     version,
		ParseResult: parseResult,
		Diagnostics: parseResult.Diagnostics,
		Config:      resolvedConfig,
	}

	s.mu.Lock()
	s.docs[uri] = doc
	s.mu.Unlock()

	return doc, nil
}

func fileURIToPath(uri string) (string, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "file" {
		return "", fmt.Errorf("unsupported uri scheme: %s", parsed.Scheme)
	}
	if parsed.Path == "" {
		return "", fmt.Errorf("empty uri path")
	}
	return filepath.FromSlash(parsed.Path), nil
}
