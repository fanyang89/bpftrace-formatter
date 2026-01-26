package lsp

import (
	"net/url"
	"path/filepath"
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
	if _, err := fileURIToPath("http://example.com"); err == nil {
		t.Fatalf("expected error for non-file scheme")
	}
	if _, err := fileURIToPath("file://"); err == nil {
		t.Fatalf("expected error for empty file path")
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
