package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fanyang89/bpftrace-formatter/config"
	"github.com/fanyang89/bpftrace-formatter/formatter"
	"github.com/fanyang89/bpftrace-formatter/lsp"
)

// runOptions holds the parsed command line options
type runOptions struct {
	generateConfig bool
	configOutput   string
	configFile     string
	inPlace        bool
	write          bool
	verbose        bool
	help           bool
	files          []string
}

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run executes the main program logic with the given arguments and output writers
func run(args []string, stdout, stderr io.Writer) error {
	// Handle LSP subcommand
	if len(args) > 1 && args[1] == "lsp" {
		return runLSP(args[2:], stderr)
	}

	opts, err := parseFlags(args)
	if err != nil {
		return err
	}

	// Handle help
	if opts.help {
		printUsageTo(stdout)
		return nil
	}

	// Handle config generation
	if opts.generateConfig {
		loader := config.NewConfigLoader()
		if err := loader.GenerateDefaultConfig(opts.configOutput); err != nil {
			return fmt.Errorf("error generating config: %w", err)
		}
		fmt.Fprintf(stdout, "Generated default configuration at: %s\n", opts.configOutput)
		return nil
	}

	// Check for input files
	if len(opts.files) == 0 {
		return fmt.Errorf("no input files specified")
	}

	// Load configuration
	cfg, err := loadConfig(opts.configFile, opts.verbose, stderr)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Process each input file
	for _, filename := range opts.files {
		err := processFile(filename, cfg, opts.inPlace || opts.write, opts.verbose, stdout, stderr)
		if err != nil {
			return fmt.Errorf("error processing %s: %w", filename, err)
		}
	}

	return nil
}

// runLSP handles the LSP subcommand
func runLSP(args []string, stderr io.Writer) error {
	lspFlags := flag.NewFlagSet("lsp", flag.ContinueOnError)
	lspFlags.SetOutput(stderr)
	help := lspFlags.Bool("help", false, "Show help message")
	lspFlags.Usage = func() {
		fmt.Fprintf(stderr, "Usage:\n  btfmt lsp [-help]\n")
	}
	if err := lspFlags.Parse(args); err != nil {
		return err
	}
	if *help {
		lspFlags.Usage()
		return nil
	}
	lsp.RunServer()
	return nil
}

// parseFlags parses command line arguments and returns options
func parseFlags(args []string) (*runOptions, error) {
	opts := &runOptions{}

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard) // Suppress default error output

	fs.BoolVar(&opts.generateConfig, "generate-config", false, "Generate default configuration file")
	fs.StringVar(&opts.configOutput, "config-output", ".btfmt.json", "Output path for generated config")
	fs.StringVar(&opts.configFile, "config", "", "Path to configuration file (.btfmt.json)")
	fs.StringVar(&opts.configFile, "c", "", "Path to configuration file (short form)")
	fs.BoolVar(&opts.inPlace, "i", false, "Edit files in place")
	fs.BoolVar(&opts.write, "w", false, "Write result to source file instead of stdout")
	fs.BoolVar(&opts.verbose, "verbose", false, "Enable verbose output")
	fs.BoolVar(&opts.verbose, "v", false, "Enable verbose output (short form)")
	fs.BoolVar(&opts.help, "help", false, "Show help message")

	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}

	opts.files = fs.Args()
	return opts, nil
}

// loadConfig loads configuration from file or defaults
func loadConfig(configFile string, verbose bool, stderr io.Writer) (*config.Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	cfg, err := config.LoadConfigFrom(cwd, configFile, false)
	if err != nil {
		return nil, err
	}

	if configFile != "" {
		configPath := configFile
		if !filepath.IsAbs(configFile) && cwd != "" {
			configPath = filepath.Join(cwd, configFile)
		}
		if _, err := os.Stat(configPath); err == nil {
			if verbose {
				fmt.Fprintf(stderr, "Using configuration file: %s\n", configPath)
			}
		} else {
			fmt.Fprintf(stderr, "Warning: specified config file %s not found\n", configPath)
		}
	}

	return cfg, nil
}

// processFile processes a single bpftrace file
func processFile(filename string, cfg *config.Config, writeToFile bool, verbose bool, stdout, stderr io.Writer) error {
	// Read input file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if verbose {
		fmt.Fprintf(stderr, "Processing: %s\n", filename)
	}

	// Create formatter and format the content
	astFormatter := formatter.NewASTFormatter(cfg)
	formatted, err := astFormatter.Format(string(content))
	if err != nil {
		return fmt.Errorf("formatting: %w", err)
	}

	if !strings.HasSuffix(formatted, "\n") {
		formatted += "\n"
	}

	// Output result
	if writeToFile {
		// Write back to the original file without resetting permissions.
		if err := writeFilePreserveMode(filename, []byte(formatted)); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		if verbose {
			fmt.Fprintf(stderr, "Formatted: %s\n", filename)
		}
	} else {
		// Print to stdout
		fmt.Fprint(stdout, formatted)
	}

	return nil
}

func writeFilePreserveMode(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}

	n, err := file.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	closeErr := file.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return nil
}

// printUsageTo prints usage information to the given writer
func printUsageTo(w io.Writer) {
	fmt.Fprintf(w, `bpftrace-formatter - A formatter for bpftrace scripts

Usage:
  btfmt [options] <file.bt> [file2.bt ...]

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
  btfmt script.bt

  # Format a file in place
  btfmt -i script.bt

  # Format multiple files
  btfmt -w file1.bt file2.bt file3.bt

  # Generate default configuration
  btfmt -generate-config

  # Use custom configuration
  btfmt -config my-config.json script.bt

Configuration:
  The formatter looks for configuration files in this order:
  1. File specified with -config flag
  2. .btfmt.json in current directory or parent directories
  3. ~/.btfmt.json in home directory
  4. Built-in defaults if no config file is found

`)
}
