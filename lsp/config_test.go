package lsp

import (
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func fileURIForPath(path string) protocol.DocumentUri {
	uri := url.URL{Scheme: "file", Path: filepath.ToSlash(path)}
	return protocol.DocumentUri(uri.String())
}

func TestExtractBtfmtSettings(t *testing.T) {
	cases := []struct {
		name         string
		settings     map[string]any
		wantOverride map[string]any
		wantPath     string
	}{
		{
			name:         "nil_settings",
			settings:     nil,
			wantOverride: nil,
			wantPath:     "",
		},
		{
			name:         "missing_btfmt",
			settings:     map[string]any{"other": true},
			wantOverride: nil,
			wantPath:     "",
		},
		{
			name:         "btfmt_not_map",
			settings:     map[string]any{"btfmt": "bad"},
			wantOverride: nil,
			wantPath:     "",
		},
		{
			name:         "config_path_only",
			settings:     map[string]any{"btfmt": map[string]any{"configPath": "/tmp/missing.json"}},
			wantOverride: nil,
			wantPath:     "/tmp/missing.json",
		},
		{
			name: "override_and_path",
			settings: map[string]any{
				"btfmt": map[string]any{
					"configPath": "/tmp/missing.json",
					"indent":     map[string]any{"size": 2},
				},
			},
			wantOverride: map[string]any{"indent": map[string]any{"size": 2}},
			wantPath:     "/tmp/missing.json",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotOverride, gotPath := extractBtfmtSettings(tc.settings)
			if !reflect.DeepEqual(gotOverride, tc.wantOverride) {
				t.Fatalf("extractBtfmtSettings override = %#v, want %#v", gotOverride, tc.wantOverride)
			}
			if gotPath != tc.wantPath {
				t.Fatalf("extractBtfmtSettings path = %q, want %q", gotPath, tc.wantPath)
			}
		})
	}
}

func TestSettingsFromConfigurationResult(t *testing.T) {
	got := settingsFromConfigurationResult(nil)
	if got != nil {
		t.Fatalf("settingsFromConfigurationResult(nil) = %#v, want nil", got)
	}

	input := []any{map[string]any{"indent": map[string]any{"size": 2}}}
	got = settingsFromConfigurationResult(input)
	want := map[string]any{"btfmt": map[string]any{"indent": map[string]any{"size": 2}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("settingsFromConfigurationResult = %#v, want %#v", got, want)
	}
}

func TestWorkspaceRootsFromParams_Priority(t *testing.T) {
	workspace := t.TempDir()
	second := t.TempDir()
	root := t.TempDir()
	rootPath := t.TempDir()

	rootURI := fileURIForPath(root)
	params := &protocol.InitializeParams{
		WorkspaceFolders: []protocol.WorkspaceFolder{
			{URI: fileURIForPath(workspace)},
			{URI: fileURIForPath(second)},
		},
		RootURI:  &rootURI,
		RootPath: &rootPath,
	}

	got := workspaceRootsFromParams(params)
	want := []string{workspace, second}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("workspaceRootsFromParams = %#v, want %#v", got, want)
	}
}

func TestWorkspaceRootsFromParams_Fallbacks(t *testing.T) {
	root := t.TempDir()
	rootURI := fileURIForPath(root)
	params := &protocol.InitializeParams{RootURI: &rootURI}

	got := workspaceRootsFromParams(params)
	want := []string{root}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("workspaceRootsFromParams rootURI = %#v, want %#v", got, want)
	}

	rootPath := t.TempDir()
	params = &protocol.InitializeParams{RootPath: &rootPath}
	got = workspaceRootsFromParams(params)
	want = []string{rootPath}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("workspaceRootsFromParams rootPath = %#v, want %#v", got, want)
	}
}

func TestWorkspaceRootsFromParams_RemoteRootURI(t *testing.T) {
	rootPath := t.TempDir()
	rootURI := protocol.DocumentUri("vscode-remote://ssh-remote+host/home/user")
	params := &protocol.InitializeParams{RootURI: &rootURI, RootPath: &rootPath}

	got := workspaceRootsFromParams(params)
	want := []string{filepath.FromSlash("/home/user")}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("workspaceRootsFromParams remote rootURI = %#v, want %#v", got, want)
	}
}

func TestConfigResolver_UsesMatchingWorkspaceRoot(t *testing.T) {
	rootA := t.TempDir()
	rootB := t.TempDir()

	if err := os.WriteFile(filepath.Join(rootA, ".btfmt.json"), []byte(`{"indent":{"size":2}}`), 0o644); err != nil {
		t.Fatalf("WriteFile rootA: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rootB, ".btfmt.json"), []byte(`{"indent":{"size":6}}`), 0o644); err != nil {
		t.Fatalf("WriteFile rootB: %v", err)
	}

	docPath := filepath.Join(rootB, "src", "probe.bt")
	if err := os.MkdirAll(filepath.Dir(docPath), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	resolver := NewConfigResolver()
	resolver.SetWorkspaceRoots([]string{rootA, rootB})

	uri := url.URL{Scheme: "file", Path: filepath.ToSlash(docPath)}
	cfg, err := resolver.ResolveForDocument(uri.String(), docPath)
	if err != nil {
		t.Fatalf("ResolveForDocument: %v", err)
	}
	if cfg == nil || cfg.Indent.Size != 6 {
		t.Fatalf("indent size = %d, want %d", cfg.Indent.Size, 6)
	}
}

func TestConfigResolver_UsesWorkspaceRootWhenDocPathEmpty(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".btfmt.json"), []byte(`{"indent":{"size":2}}`), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	resolver := NewConfigResolver()
	resolver.SetWorkspaceRoots([]string{root})

	cfg, err := resolver.ResolveForDocument("untitled:Untitled-1", "")
	if err != nil {
		t.Fatalf("ResolveForDocument: %v", err)
	}
	if cfg == nil || cfg.Indent.Size != 2 {
		t.Fatalf("indent size = %d, want %d", cfg.Indent.Size, 2)
	}
}

func TestConfigResolver_ReflectsUpdatedConfigFile(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, ".btfmt.json")
	if err := os.WriteFile(configPath, []byte(`{"indent":{"size":2}}`), 0o644); err != nil {
		t.Fatalf("WriteFile initial: %v", err)
	}

	docPath := filepath.Join(root, "probe.bt")
	resolver := NewConfigResolver()
	resolver.SetWorkspaceRoots([]string{root})

	uri := fileURIForPath(docPath)
	cfg, err := resolver.ResolveForDocument(string(uri), docPath)
	if err != nil {
		t.Fatalf("ResolveForDocument initial: %v", err)
	}
	if cfg == nil || cfg.Indent.Size != 2 {
		t.Fatalf("initial indent size = %d, want %d", cfg.Indent.Size, 2)
	}

	if err := os.WriteFile(configPath, []byte(`{"indent":{"size":6}}`), 0o644); err != nil {
		t.Fatalf("WriteFile updated: %v", err)
	}

	cfg, err = resolver.ResolveForDocument(string(uri), docPath)
	if err != nil {
		t.Fatalf("ResolveForDocument updated: %v", err)
	}
	if cfg == nil || cfg.Indent.Size != 6 {
		t.Fatalf("updated indent size = %d, want %d", cfg.Indent.Size, 6)
	}
}
