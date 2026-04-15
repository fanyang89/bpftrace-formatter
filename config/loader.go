package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ConfigLoader handles loading configuration from various sources.
type ConfigLoader struct {
	baseDir      string
	explicitPath string
	verbose      bool
	logWriter    io.Writer
}

// NewConfigLoader creates a new config loader with default settings.
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		logWriter: io.Discard,
	}
}

// WithBaseDir sets the base directory for relative configuration discovery.
func (cl *ConfigLoader) WithBaseDir(baseDir string) *ConfigLoader {
	cl.baseDir = baseDir
	return cl
}

// WithExplicitPath sets an explicit path to a configuration file.
func (cl *ConfigLoader) WithExplicitPath(path string) *ConfigLoader {
	cl.explicitPath = path
	return cl
}

// WithVerbose enables or disables verbose logging.
func (cl *ConfigLoader) WithVerbose(verbose bool) *ConfigLoader {
	cl.verbose = verbose
	return cl
}

// WithLogger sets the writer for verbose logging.
func (cl *ConfigLoader) WithLogger(w io.Writer) *ConfigLoader {
	cl.logWriter = w
	return cl
}

// LoadConfig loads configuration from file or returns default.
// It searches in the following order:
// 1. Explicitly provided path (if any)
// 2. .btfmt.json in baseDir or its parents (if baseDir is set)
// 3. ~/.btfmt.json in user's home directory
// 4. Built-in defaults
func (cl *ConfigLoader) LoadConfig() (*Config, error) {
	configPath, isExplicit := cl.findConfigPath()

	if configPath == "" {
		if cl.verbose {
			fmt.Fprintln(cl.logWriter, "No configuration file found, using defaults")
		}
		return DefaultConfig(), nil
	}

	if _, err := os.Stat(configPath); err != nil {
		if isExplicit {
			fmt.Fprintf(cl.logWriter, "Warning: specified config file %s not found, using defaults\n", configPath)
		} else if cl.verbose {
			fmt.Fprintf(cl.logWriter, "No configuration file found at %s, using defaults\n", configPath)
		}
		return DefaultConfig(), nil
	}

	if cl.verbose {
		fmt.Fprintf(cl.logWriter, "Loading configuration from: %s\n", configPath)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
	}

	if validationErrors := config.Validate(); len(validationErrors) > 0 {
		return nil, fmt.Errorf("invalid configuration in %s:\n%s", configPath, formatValidationErrors(validationErrors))
	}

	return config, nil
}

// findConfigPath determines the path to the configuration file based on priority.
// Returns the path and a boolean indicating if it was explicitly requested.
func (cl *ConfigLoader) findConfigPath() (string, bool) {
	if cl.explicitPath != "" {
		path := cl.explicitPath
		if !filepath.IsAbs(path) && cl.baseDir != "" {
			path = filepath.Join(cl.baseDir, path)
		}
		return path, true
	}

	if cl.baseDir != "" {
		if path := SearchUpwards(cl.baseDir, ".btfmt.json"); path != "" {
			return path, false
		}
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeConfig := filepath.Join(homeDir, ".btfmt.json")
		if _, err := os.Stat(homeConfig); err == nil {
			return homeConfig, false
		}
	}

	return "", false
}

// SearchUpwards searches for a file in the startDir and its parent directories.
func SearchUpwards(startDir, filename string) string {
	currentDir := startDir

	for {
		configPath := filepath.Join(currentDir, filename)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return ""
}

// LoadConfigFrom is a convenience wrapper for backward compatibility.
func LoadConfigFrom(baseDir, explicitPath string, verbose bool) (*Config, error) {
	return NewConfigLoader().
		WithBaseDir(baseDir).
		WithExplicitPath(explicitPath).
		WithVerbose(verbose).
		WithLogger(os.Stdout).
		LoadConfig()
}

// LoadConfigFromWithLogger is a convenience wrapper for backward compatibility.
func LoadConfigFromWithLogger(baseDir, explicitPath string, verbose bool, logWriter io.Writer) (*Config, error) {
	return NewConfigLoader().
		WithBaseDir(baseDir).
		WithExplicitPath(explicitPath).
		WithVerbose(verbose).
		WithLogger(logWriter).
		LoadConfig()
}

// GenerateDefaultConfig generates a default configuration file.
func (cl *ConfigLoader) GenerateDefaultConfig(outputPath string) error {
	return DefaultConfig().SaveConfig(outputPath)
}

// formatValidationErrors formats a slice of validation errors into a single error message.
func formatValidationErrors(validationErrors []error) string {
	if len(validationErrors) == 0 {
		return ""
	}
	if len(validationErrors) == 1 {
		return validationErrors[0].Error()
	}

	var msg strings.Builder
	msg.WriteString("multiple validation errors:\n")
	for _, err := range validationErrors {
		msg.WriteString(fmt.Sprintf("  - %s\n", err.Error()))
	}
	return msg.String()
}
