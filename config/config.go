package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all formatting configuration options
type Config struct {
	// Indentation settings
	Indent IndentConfig `json:"indent"`

	// Spacing settings
	Spacing SpacingConfig `json:"spacing"`

	// Line break settings
	LineBreaks LineBreakConfig `json:"line_breaks"`

	// Comment formatting
	Comments CommentConfig `json:"comments"`

	// Probe formatting
	Probes ProbeConfig `json:"probes"`

	// Block formatting
	Blocks BlockConfig `json:"blocks"`
}

// IndentConfig controls indentation behavior
type IndentConfig struct {
	Size      int  `json:"size"`       // Number of spaces/tabs per indent level
	UseSpaces bool `json:"use_spaces"` // Use spaces instead of tabs
}

// SpacingConfig controls spacing around operators and keywords
type SpacingConfig struct {
	AroundOperators   bool `json:"around_operators"`   // Space around =, +, -, etc.
	AroundCommas      bool `json:"around_commas"`      // Space after commas
	AroundParentheses bool `json:"around_parentheses"` // Space inside parentheses
	AroundBrackets    bool `json:"around_brackets"`    // Space inside brackets
	BeforeBlockStart  bool `json:"before_block_start"` // Space before {
	AfterKeywords     bool `json:"after_keywords"`     // Space after if, while, etc.
}

// LineBreakConfig controls line break behavior
type LineBreakConfig struct {
	MaxLineLength           int  `json:"max_line_length"`            // Maximum line length before wrapping
	BreakLongStatements     bool `json:"break_long_statements"`      // Break long statements across lines
	EmptyLinesBetweenProbes int  `json:"empty_lines_between_probes"` // Number of empty lines between probes
	EmptyLinesAfterShebang  int  `json:"empty_lines_after_shebang"`  // Number of empty lines after shebang
}

// CommentConfig controls comment formatting
type CommentConfig struct {
	PreserveInline bool `json:"preserve_inline"` // Keep inline comments on same line
	AlignInline    bool `json:"align_inline"`    // Align inline comments
	IndentLevel    int  `json:"indent_level"`    // Indent level for standalone comments
}

// ProbeConfig controls probe formatting
type ProbeConfig struct {
	AlignPredicates bool `json:"align_predicates"` // Align predicates with probe definitions
	SortProbes      bool `json:"sort_probes"`      // Sort probes alphabetically
	GroupByType     bool `json:"group_by_type"`    // Group probes by type (tracepoint, kprobe, etc.)
}

// BlockConfig controls block formatting
type BlockConfig struct {
	BraceStyle        string `json:"brace_style"`          // "same_line", "next_line", "gnu"
	IndentStatements  bool   `json:"indent_statements"`    // Indent statements inside blocks
	EmptyLineInBlocks bool   `json:"empty_line_in_blocks"` // Add empty lines in large blocks
	AlignAssignments  bool   `json:"align_assignments"`    // Align assignment operators
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Indent: IndentConfig{
			Size:      4,
			UseSpaces: true,
		},
		Spacing: SpacingConfig{
			AroundOperators:   true,
			AroundCommas:      true,
			AroundParentheses: false,
			AroundBrackets:    false,
			BeforeBlockStart:  true,
			AfterKeywords:     true,
		},
		LineBreaks: LineBreakConfig{
			MaxLineLength:           80,
			BreakLongStatements:     true,
			EmptyLinesBetweenProbes: 1,
			EmptyLinesAfterShebang:  1,
		},
		Comments: CommentConfig{
			PreserveInline: true,
			AlignInline:    false,
			IndentLevel:    0,
		},
		Probes: ProbeConfig{
			AlignPredicates: false,
			SortProbes:      false,
			GroupByType:     false,
		},
		Blocks: BlockConfig{
			BraceStyle:        "same_line",
			IndentStatements:  true,
			EmptyLineInBlocks: false,
			AlignAssignments:  false,
		},
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// LoadConfigWithFallback loads config from file, falls back to default if file doesn't exist
func LoadConfigWithFallback(filename string) *Config {
	if config, err := LoadConfig(filename); err == nil {
		return config
	}
	return DefaultConfig()
}

// SaveConfig saves configuration to a JSON file
func (c *Config) SaveConfig(filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
