package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fanyang89/bpftrace-formatter/config"
	"github.com/fanyang89/bpftrace-formatter/formatter"
)

func main() {
	// Parse command line flags
	var (
		generateConfig = flag.Bool("generate-config", false, "Generate default configuration file")
		configOutput   = flag.String("config-output", ".btfmt.json", "Output path for generated config")
		inPlace        = flag.Bool("i", false, "Edit files in place")
		write          = flag.Bool("w", false, "Write result to source file instead of stdout")
		help           = flag.Bool("help", false, "Show help message")
	)

	// Create config loader
	configLoader := config.NewConfigLoader()
	configLoader.ParseFlags()

	// Handle help
	if *help {
		printUsage()
		return
	}

	// Handle config generation
	if *generateConfig {
		err := configLoader.GenerateDefaultConfig(*configOutput)
		if err != nil {
			fmt.Printf("Error generating config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated default configuration at: %s\n", *configOutput)
		return
	}

	// Get remaining arguments (input files)
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error: No input files specified")
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := configLoader.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Process each input file
	for _, filename := range args {
		err := processFile(filename, cfg, *inPlace || *write, configLoader.GetVerbose())
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", filename, err)
			os.Exit(1)
		}
	}
}

// processFile processes a single bpftrace file
func processFile(filename string, cfg *config.Config, writeToFile bool, verbose bool) error {
	// Read input file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if verbose {
		fmt.Printf("Processing: %s\n", filename)
	}

	// Create formatter and format the content
	astFormatter := formatter.NewASTFormatter(cfg)
	formatted, err := astFormatter.Format(string(content))
	if err != nil {
		return fmt.Errorf("formatting: %w", err)
	}

	// Output result
	if writeToFile {
		// Write back to the original file
		err = os.WriteFile(filename, []byte(formatted), 0644)
		if err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		if verbose {
			fmt.Printf("Formatted: %s\n", filename)
		}
	} else {
		// Print to stdout
		fmt.Print(formatted)
		if !strings.HasSuffix(formatted, "\n") {
			fmt.Println()
		}
	}

	return nil
}

// printUsage prints usage information
func printUsage() {
	fmt.Printf(`bpftrace-formatter - A formatter for bpftrace scripts

Usage:
  %s [options] <file.bt> [file2.bt ...]

Options:
  -c, -config <file>     Path to configuration file
  -i                     Edit files in place
  -w                     Write result to source file instead of stdout
  -v, -verbose           Enable verbose output
  -generate-config       Generate default configuration file
  -config-output <file>  Output path for generated config (default: .btfmt.json)
  -help                  Show this help message

Examples:
  # Format a file and print to stdout
  %s script.bt
  
  # Format a file in place
  %s -i script.bt
  
  # Format multiple files
  %s -w file1.bt file2.bt file3.bt
  
  # Generate default configuration
  %s -generate-config
  
  # Use custom configuration
  %s -config my-config.json script.bt

Configuration:
  The formatter looks for configuration files in this order:
  1. File specified with -config flag
  2. .btfmt.json in current directory or parent directories
  3. ~/.btfmt.json in home directory
  4. Built-in defaults if no config file is found

`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
}
