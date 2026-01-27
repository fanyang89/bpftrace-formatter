package config

import (
	"fmt"
	"strings"
)

// ValidBraceStyles contains the allowed values for brace_style configuration.
var ValidBraceStyles = []string{"same_line", "next_line", "gnu"}

// Validate checks the configuration for invalid values and returns a slice of errors.
// An empty slice indicates the configuration is valid.
func (c *Config) Validate() []error {
	var errors []error

	// Validate BraceStyle
	if !isValidBraceStyle(c.Blocks.BraceStyle) {
		errors = append(errors, fmt.Errorf(
			"invalid brace_style %q: must be one of %s",
			c.Blocks.BraceStyle,
			strings.Join(ValidBraceStyles, ", "),
		))
	}

	// Validate indent size (1-16)
	if c.Indent.Size < 1 || c.Indent.Size > 16 {
		errors = append(errors, fmt.Errorf(
			"indent.size must be between 1 and 16, got %d",
			c.Indent.Size,
		))
	}

	// Validate max line length (minimum 40)
	if c.LineBreaks.MaxLineLength < 40 {
		errors = append(errors, fmt.Errorf(
			"line_breaks.max_line_length must be at least 40, got %d",
			c.LineBreaks.MaxLineLength,
		))
	}

	// Validate empty lines between probes (0-5)
	if c.LineBreaks.EmptyLinesBetweenProbes < 0 || c.LineBreaks.EmptyLinesBetweenProbes > 5 {
		errors = append(errors, fmt.Errorf(
			"line_breaks.empty_lines_between_probes must be 0-5, got %d",
			c.LineBreaks.EmptyLinesBetweenProbes,
		))
	}

	// Validate empty lines after shebang (0-5)
	if c.LineBreaks.EmptyLinesAfterShebang < 0 || c.LineBreaks.EmptyLinesAfterShebang > 5 {
		errors = append(errors, fmt.Errorf(
			"line_breaks.empty_lines_after_shebang must be 0-5, got %d",
			c.LineBreaks.EmptyLinesAfterShebang,
		))
	}

	return errors
}

// isValidBraceStyle checks if the given style is one of the valid brace styles.
func isValidBraceStyle(style string) bool {
	for _, valid := range ValidBraceStyles {
		if style == valid {
			return true
		}
	}
	return false
}
