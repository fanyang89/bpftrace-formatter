package lsp

import (
	"net/url"
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

func TestWorkspaceRootFromParams_Priority(t *testing.T) {
	workspace := t.TempDir()
	root := t.TempDir()
	rootPath := t.TempDir()

	rootURI := fileURIForPath(root)
	params := &protocol.InitializeParams{
		WorkspaceFolders: []protocol.WorkspaceFolder{{URI: fileURIForPath(workspace)}},
		RootURI:          &rootURI,
		RootPath:         &rootPath,
	}

	got := workspaceRootFromParams(params)
	if got != workspace {
		t.Fatalf("workspaceRootFromParams = %q, want %q", got, workspace)
	}
}

func TestWorkspaceRootFromParams_Fallbacks(t *testing.T) {
	root := t.TempDir()
	rootURI := fileURIForPath(root)
	params := &protocol.InitializeParams{RootURI: &rootURI}

	got := workspaceRootFromParams(params)
	if got != root {
		t.Fatalf("workspaceRootFromParams rootURI = %q, want %q", got, root)
	}

	rootPath := t.TempDir()
	params = &protocol.InitializeParams{RootPath: &rootPath}
	got = workspaceRootFromParams(params)
	if got != rootPath {
		t.Fatalf("workspaceRootFromParams rootPath = %q, want %q", got, rootPath)
	}
}
