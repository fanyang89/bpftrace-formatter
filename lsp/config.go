package lsp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/config"
)

type ConfigResolver struct {
	mu             sync.Mutex
	workspaceRoots []string
	settings       map[string]any
	cache          map[string]*config.Config
}

func NewConfigResolver() *ConfigResolver {
	return &ConfigResolver{cache: make(map[string]*config.Config)}
}

func (r *ConfigResolver) SetWorkspaceRoot(root string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if root == "" {
		r.workspaceRoots = nil
	} else {
		r.workspaceRoots = []string{filepath.Clean(root)}
	}
	r.cache = make(map[string]*config.Config)
}

func (r *ConfigResolver) SetWorkspaceRoots(roots []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workspaceRoots = normalizeWorkspaceRoots(roots)
	r.cache = make(map[string]*config.Config)
}

func (r *ConfigResolver) SetSettings(settings map[string]any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.settings = settings
	r.cache = make(map[string]*config.Config)
}

func (r *ConfigResolver) ResolveForDocument(uri string, docPath string) (*config.Config, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cached, ok := r.cache[uri]; ok {
		return cached, nil
	}

	baseDir := selectWorkspaceRoot(docPath, r.workspaceRoots)
	if baseDir == "" {
		if docPath == "" {
			if len(r.workspaceRoots) > 0 {
				baseDir = r.workspaceRoots[0]
			}
		} else {
			baseDir = filepath.Dir(docPath)
		}
	}

	settingsOverride, configPath := extractBtfmtSettings(r.settings)

	fileConfig, err := config.LoadConfigFrom(baseDir, configPath, false)
	if err != nil {
		return nil, err
	}

	mergedConfig := fileConfig
	if len(settingsOverride) > 0 {
		merged, warnings, err := config.MergeConfig(fileConfig, settingsOverride)
		if err != nil {
			return nil, err
		}
		for _, warning := range warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		}
		mergedConfig = merged
	}

	r.cache[uri] = mergedConfig
	return mergedConfig, nil
}

func extractBtfmtSettings(settings map[string]any) (map[string]any, string) {
	if settings == nil {
		return nil, ""
	}

	raw, ok := settings["btfmt"]
	if !ok {
		return nil, ""
	}

	btfmtSettings, ok := raw.(map[string]any)
	if !ok {
		return nil, ""
	}

	configPath := ""
	if value, ok := btfmtSettings["configPath"]; ok {
		if asString, ok := value.(string); ok {
			configPath = asString
		}
	}

	override := map[string]any{}
	for key, value := range btfmtSettings {
		if key == "configPath" {
			continue
		}
		override[key] = value
	}

	if len(override) == 0 {
		return nil, configPath
	}

	return override, configPath
}

func workspaceRootFromParams(params *protocol.InitializeParams) string {
	roots := workspaceRootsFromParams(params)
	if len(roots) > 0 {
		return roots[0]
	}
	return ""
}

func workspaceRootsFromParams(params *protocol.InitializeParams) []string {
	if params == nil {
		return nil
	}

	if len(params.WorkspaceFolders) > 0 {
		roots := make([]string, 0, len(params.WorkspaceFolders))
		for _, folder := range params.WorkspaceFolders {
			if path, err := fileURIToPath(folder.URI); err == nil && path != "" {
				roots = append(roots, path)
			}
		}
		if len(roots) > 0 {
			return roots
		}
	}

	if params.RootURI != nil {
		if path, err := fileURIToPath(string(*params.RootURI)); err == nil {
			return []string{path}
		}
	}

	if params.RootPath != nil {
		return []string{*params.RootPath}
	}

	return nil
}

func settingsFromParams(params *protocol.InitializeParams) map[string]any {
	if params == nil || params.InitializationOptions == nil {
		return nil
	}

	if asMap, ok := params.InitializationOptions.(map[string]any); ok {
		return asMap
	}

	return nil
}

func settingsFromConfigurationResult(result []any) map[string]any {
	if len(result) == 0 {
		return nil
	}

	settings, ok := result[0].(map[string]any)
	if !ok {
		return nil
	}

	return normalizeSettingsMap(settings)
}

func normalizeSettingsMap(settings map[string]any) map[string]any {
	if settings == nil {
		return nil
	}
	if _, ok := settings["btfmt"]; ok {
		return settings
	}

	return map[string]any{"btfmt": settings}
}

func normalizeWorkspaceRoots(roots []string) []string {
	if len(roots) == 0 {
		return nil
	}
	normalized := make([]string, 0, len(roots))
	for _, root := range roots {
		if root == "" {
			continue
		}
		normalized = append(normalized, filepath.Clean(root))
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

func selectWorkspaceRoot(docPath string, roots []string) string {
	if docPath == "" || len(roots) == 0 {
		return ""
	}
	docPath = filepath.Clean(docPath)
	best := ""
	bestLen := -1
	for _, root := range roots {
		if root == "" {
			continue
		}
		cleanRoot := filepath.Clean(root)
		if !isSubpath(cleanRoot, docPath) {
			continue
		}
		if len(cleanRoot) > bestLen {
			best = cleanRoot
			bestLen = len(cleanRoot)
		}
	}
	return best
}

func isSubpath(root string, path string) bool {
	if root == "" || path == "" {
		return false
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	prefix := ".." + string(filepath.Separator)
	if rel == ".." || strings.HasPrefix(rel, prefix) {
		return false
	}
	return true
}
