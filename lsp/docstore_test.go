package lsp

import (
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanyang89/bpftrace-formatter/config"
)

func TestFileURIToPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "space dir")
	uri := url.URL{Scheme: "file", Path: filepath.ToSlash(path)}

	got, err := fileURIToPath(uri.String())
	if err != nil {
		t.Fatalf("fileURIToPath: %v", err)
	}
	if got != path {
		t.Fatalf("fileURIToPath = %q, want %q", got, path)
	}
}

func TestFileURIToPath_Errors(t *testing.T) {
	if _, err := fileURIToPath("file://"); err == nil {
		t.Fatalf("expected error for empty file path")
	}
}

func TestFileURIToPath_EmptySchemeErrorMessage(t *testing.T) {
	_, err := fileURIToPath("no-scheme-path")
	if err == nil {
		t.Fatalf("expected error for URI without scheme")
	}
	if !strings.Contains(err.Error(), "missing uri scheme") {
		t.Fatalf("expected 'missing uri scheme' in error, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "no-scheme-path") {
		t.Fatalf("expected URI in error message, got %q", err.Error())
	}
}

func TestFileURIToPath_NonFileScheme(t *testing.T) {
	got, err := fileURIToPath("untitled:Untitled-1")
	if err != nil {
		t.Fatalf("fileURIToPath: %v", err)
	}
	if got != "" {
		t.Fatalf("fileURIToPath = %q, want empty string", got)
	}
}

func TestFileURIToPath_VscodeRemote(t *testing.T) {
	uri := "vscode-remote://ssh-remote+host/home/user/workspace/test.bt"

	got, err := fileURIToPath(uri)
	if err != nil {
		t.Fatalf("fileURIToPath: %v", err)
	}
	want := filepath.FromSlash("/home/user/workspace/test.bt")
	if got != want {
		t.Fatalf("fileURIToPath = %q, want %q", got, want)
	}
}

func TestFileURIToPath_WindowsForms(t *testing.T) {
	cases := []struct {
		name string
		uri  string
		want string
	}{
		{
			name: "drive_letter",
			uri:  "file:///C:/Windows/System32",
			want: filepath.FromSlash("C:/Windows/System32"),
		},
		{
			name: "unc_host",
			uri:  "file://server/share/path/file.bt",
			want: filepath.FromSlash("//server/share/path/file.bt"),
		},
		{
			name: "drive_host",
			uri:  "file://C:/Windows/System32",
			want: filepath.FromSlash("C:/Windows/System32"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := fileURIToPath(tc.uri)
			if err != nil {
				t.Fatalf("fileURIToPath(%q): %v", tc.uri, err)
			}
			if got != tc.want {
				t.Fatalf("fileURIToPath(%q) = %q, want %q", tc.uri, got, tc.want)
			}
		})
	}
}

func TestDocumentStoreOpen_DefaultConfig(t *testing.T) {
	store := NewDocumentStore(nil)
	uri := url.URL{Scheme: "file", Path: filepath.ToSlash(filepath.Join(t.TempDir(), "test.bt"))}

	doc, err := store.Open(uri.String(), 1, "kprobe:sys_clone { @x = count(); }\n")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if doc.Config == nil {
		t.Fatalf("expected config")
	}
	if doc.Config.Indent.Size != config.DefaultConfig().Indent.Size {
		t.Fatalf("indent size = %d, want %d", doc.Config.Indent.Size, config.DefaultConfig().Indent.Size)
	}
}

func TestDocumentStoreOpen_NonFileURI(t *testing.T) {
	store := NewDocumentStore(nil)
	uri := "untitled:Untitled-1"

	doc, err := store.Open(uri, 1, "kprobe:sys_clone { @x = count(); }\n")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if doc.Path != "" {
		t.Fatalf("path = %q, want empty string", doc.Path)
	}
}

func TestDocumentStoreAllDocs(t *testing.T) {
	store := NewDocumentStore(nil)
	tmpDir := t.TempDir()

	uri1 := (&url.URL{Scheme: "file", Path: filepath.ToSlash(filepath.Join(tmpDir, "a.bt"))}).String()
	uri2 := (&url.URL{Scheme: "file", Path: filepath.ToSlash(filepath.Join(tmpDir, "b.bt"))}).String()

	if _, err := store.Open(uri1, 1, "kprobe:sys_clone { @x = count(); }\n"); err != nil {
		t.Fatalf("Open: %v", err)
	}
	if _, err := store.Open(uri2, 3, "kprobe:sys_clone { @x = count( }\n"); err != nil {
		t.Fatalf("Open: %v", err)
	}

	snapshots := store.AllDocs()
	if len(snapshots) != 2 {
		t.Fatalf("AllDocs: expected 2 snapshots, got %d", len(snapshots))
	}

	byURI := make(map[string]DocSnapshot)
	for _, s := range snapshots {
		byURI[s.URI] = s
	}

	s1, ok := byURI[uri1]
	if !ok {
		t.Fatalf("AllDocs: missing %s", uri1)
	}
	if s1.Version != 1 {
		t.Fatalf("AllDocs: version = %d, want 1", s1.Version)
	}

	s2, ok := byURI[uri2]
	if !ok {
		t.Fatalf("AllDocs: missing %s", uri2)
	}
	if s2.Version != 3 {
		t.Fatalf("AllDocs: version = %d, want 3", s2.Version)
	}
	if len(s2.Diagnostics) == 0 {
		t.Fatalf("AllDocs: expected diagnostics for invalid document")
	}
}

func TestDocumentStoreAllDocs_Empty(t *testing.T) {
	store := NewDocumentStore(nil)
	snapshots := store.AllDocs()
	if len(snapshots) != 0 {
		t.Fatalf("AllDocs: expected 0 snapshots, got %d", len(snapshots))
	}
}

func TestDocumentStoreRefreshConfigs(t *testing.T) {
	resolver := NewConfigResolver()
	missing := filepath.Join(t.TempDir(), "missing.json")
	resolver.SetSettings(map[string]any{
		"btfmt": map[string]any{
			"configPath": missing,
			"indent":     map[string]any{"size": 2},
		},
	})
	store := NewDocumentStore(resolver)
	uri := url.URL{Scheme: "file", Path: filepath.ToSlash(filepath.Join(t.TempDir(), "test.bt"))}

	doc, err := store.Open(uri.String(), 1, "kprobe:sys_clone { @x = count(); }\n")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if doc.Config == nil || doc.Config.Indent.Size != 2 {
		t.Fatalf("expected indent size 2, got %v", doc.Config)
	}

	resolver.SetSettings(map[string]any{
		"btfmt": map[string]any{
			"configPath": missing,
			"indent":     map[string]any{"size": 6},
		},
	})
	if err := store.RefreshConfigs(); err != nil {
		t.Fatalf("RefreshConfigs: %v", err)
	}
	doc, ok := store.Get(uri.String())
	if !ok || doc == nil {
		t.Fatalf("expected document in store")
	}
	if doc.Config == nil || doc.Config.Indent.Size != 6 {
		t.Fatalf("expected indent size 6, got %v", doc.Config)
	}
}
