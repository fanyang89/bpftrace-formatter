package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ConfigLoader handles loading configuration from various sources
type ConfigLoader struct {
	configFile string
	verbose    bool
}

// NewConfigLoader creates a new config loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

// ParseFlags parses command line flags
func (cl *ConfigLoader) ParseFlags() {
	flag.StringVar(&cl.configFile, "config", "", "Path to configuration file (.btfmt.json)")
	flag.StringVar(&cl.configFile, "c", "", "Path to configuration file (short form)")
	flag.BoolVar(&cl.verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&cl.verbose, "v", false, "Enable verbose output (short form)")
	flag.Parse()
}

// LoadConfig loads configuration from file or returns default
func (cl *ConfigLoader) LoadConfig() (*Config, error) {
	configPath := cl.findConfigFile()

	if configPath == "" {
		if cl.verbose {
			fmt.Println("No configuration file found, using defaults")
		}
		return DefaultConfig(), nil
	}

	if cl.verbose {
		fmt.Printf("Loading configuration from: %s\n", configPath)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
	}

	// Validate configuration
	if validationErrors := config.Validate(); len(validationErrors) > 0 {
		return nil, fmt.Errorf("invalid configuration in %s:\n%s", configPath, formatValidationErrors(validationErrors))
	}

	return config, nil
}

// LoadConfigFrom loads configuration relative to baseDir, using explicitPath if provided.
// If explicitPath is set but missing, defaults are returned without searching.
func LoadConfigFrom(baseDir, explicitPath string, verbose bool) (*Config, error) {
	return LoadConfigFromWithLogger(baseDir, explicitPath, verbose, os.Stdout)
}

// formatValidationErrors formats a slice of validation errors into a single error message
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

// LoadConfigFromWithLogger loads configuration with optional verbose logging to logWriter.
func LoadConfigFromWithLogger(baseDir, explicitPath string, verbose bool, logWriter io.Writer) (*Config, error) {
	if logWriter == nil {
		logWriter = io.Discard
	}

	if explicitPath != "" {
		configPath := explicitPath
		if !filepath.IsAbs(explicitPath) && baseDir != "" {
			configPath = filepath.Join(baseDir, explicitPath)
		}
		if _, err := os.Stat(configPath); err == nil {
			if verbose {
				fmt.Fprintf(logWriter, "Loading configuration from: %s\n", configPath)
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
		// Always warn when explicitly specified config file is not found
		fmt.Fprintf(logWriter, "Warning: specified config file %s not found, using defaults\n", configPath)
		return DefaultConfig(), nil
	}

	if baseDir != "" {
		if configPath := searchUpwards(baseDir, ".btfmt.json"); configPath != "" {
			if verbose {
				fmt.Fprintf(logWriter, "Loading configuration from: %s\n", configPath)
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
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeConfig := filepath.Join(homeDir, ".btfmt.json")
		if _, err := os.Stat(homeConfig); err == nil {
			if verbose {
				fmt.Fprintf(logWriter, "Loading configuration from: %s\n", homeConfig)
			}
			config, err := LoadConfig(homeConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to load config from %s: %w", homeConfig, err)
			}
			if validationErrors := config.Validate(); len(validationErrors) > 0 {
				return nil, fmt.Errorf("invalid configuration in %s:\n%s", homeConfig, formatValidationErrors(validationErrors))
			}
			return config, nil
		}
	}

	if verbose {
		fmt.Fprintln(logWriter, "No configuration file found, using defaults")
	}
	return DefaultConfig(), nil
}

// findConfigFile finds the configuration file in order of precedence:
// 1. Command line specified file
// 2. .btfmt.json in current directory
// 3. .btfmt.json in parent directories (up to root)
// 4. ~/.btfmt.json in home directory
func (cl *ConfigLoader) findConfigFile() string {
	// 1. Command line specified file
	if cl.configFile != "" {
		if _, err := os.Stat(cl.configFile); err == nil {
			return cl.configFile
		}
		// Always warn when explicitly specified config file is not found
		fmt.Printf("Warning: specified config file %s not found\n", cl.configFile)
		return ""
	}

	// 2. .btfmt.json in current directory and parent directories
	cwd, err := os.Getwd()
	if err == nil {
		if configPath := cl.searchUpwards(cwd, ".btfmt.json"); configPath != "" {
			return configPath
		}
	}

	// 3. ~/.btfmt.json in home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeConfig := filepath.Join(homeDir, ".btfmt.json")
		if _, err := os.Stat(homeConfig); err == nil {
			return homeConfig
		}
	}

	return ""
}

// searchUpwards searches for a file in the current directory and parent directories
func (cl *ConfigLoader) searchUpwards(startDir, filename string) string {
	return searchUpwards(startDir, filename)
}

func searchUpwards(startDir, filename string) string {
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

// GetVerbose returns whether verbose mode is enabled
func (cl *ConfigLoader) GetVerbose() bool {
	return cl.verbose
}

// GenerateDefaultConfig generates a default configuration file
func (cl *ConfigLoader) GenerateDefaultConfig(outputPath string) error {
	config := DefaultConfig()
	return config.SaveConfig(outputPath)
}
