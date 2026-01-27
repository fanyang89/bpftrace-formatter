package lsp

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
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

func (s *DocumentStore) RefreshConfigs() error {
	if s.resolver == nil {
		return nil
	}

	type docRef struct {
		uri  string
		path string
	}

	s.mu.RLock()
	docs := make([]docRef, 0, len(s.docs))
	for uri, doc := range s.docs {
		if doc == nil {
			continue
		}
		docs = append(docs, docRef{uri: uri, path: doc.Path})
	}
	s.mu.RUnlock()

	for _, doc := range docs {
		resolvedConfig, err := s.resolver.ResolveForDocument(doc.uri, doc.path)
		if err != nil {
			return err
		}
		s.mu.Lock()
		if current, ok := s.docs[doc.uri]; ok && current != nil {
			current.Config = resolvedConfig
		}
		s.mu.Unlock()
	}

	return nil
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
	if parsed.Scheme == "" {
		return "", fmt.Errorf("unsupported uri scheme: %s", parsed.Scheme)
	}
	if parsed.Scheme != "file" {
		return "", nil
	}
	if parsed.Path == "" && parsed.Host == "" {
		return "", fmt.Errorf("empty uri path")
	}
	unescapedPath, err := url.PathUnescape(parsed.Path)
	if err != nil {
		return "", err
	}

	if parsed.Host != "" {
		if isDriveHost(parsed.Host) {
			if strings.HasPrefix(unescapedPath, "/") {
				unescapedPath = parsed.Host + unescapedPath
			} else if unescapedPath == "" {
				unescapedPath = parsed.Host
			} else {
				unescapedPath = parsed.Host + "/" + unescapedPath
			}
		} else {
			unescapedPath = "//" + parsed.Host + unescapedPath
		}
	}

	unescapedPath = trimLeadingDriveSlash(unescapedPath)

	return filepath.FromSlash(unescapedPath), nil
}

func isDriveHost(host string) bool {
	if len(host) != 2 {
		return false
	}
	if host[1] != ':' {
		return false
	}
	c := host[0]
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func trimLeadingDriveSlash(path string) string {
	if len(path) < 3 {
		return path
	}
	if path[0] != '/' {
		return path
	}
	c := path[1]
	if ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) && path[2] == ':' {
		return path[1:]
	}
	return path
}
