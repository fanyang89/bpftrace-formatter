package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	return config, nil
}

// LoadConfigFrom loads configuration relative to baseDir, using explicitPath if provided.
// If explicitPath is set but missing, defaults are returned without searching.
func LoadConfigFrom(baseDir, explicitPath string, verbose bool) (*Config, error) {
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err == nil {
			if verbose {
				fmt.Printf("Loading configuration from: %s\n", explicitPath)
			}
			config, err := LoadConfig(explicitPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load config from %s: %w", explicitPath, err)
			}
			return config, nil
		}
		if verbose {
			fmt.Printf("Warning: specified config file %s not found\n", explicitPath)
		}
		return DefaultConfig(), nil
	}

	if baseDir != "" {
		if configPath := searchUpwards(baseDir, ".btfmt.json"); configPath != "" {
			if verbose {
				fmt.Printf("Loading configuration from: %s\n", configPath)
			}
			config, err := LoadConfig(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
			}
			return config, nil
		}
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeConfig := filepath.Join(homeDir, ".btfmt.json")
		if _, err := os.Stat(homeConfig); err == nil {
			if verbose {
				fmt.Printf("Loading configuration from: %s\n", homeConfig)
			}
			config, err := LoadConfig(homeConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to load config from %s: %w", homeConfig, err)
			}
			return config, nil
		}
	}

	if verbose {
		fmt.Println("No configuration file found, using defaults")
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
		if cl.verbose {
			fmt.Printf("Warning: specified config file %s not found\n", cl.configFile)
		}
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
