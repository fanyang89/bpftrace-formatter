package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fanyang89/bpftrace-formatter/config"
	"github.com/fanyang89/bpftrace-formatter/formatter"
	"github.com/fanyang89/bpftrace-formatter/lsp"
)

// version is set at build time via ldflags
var version = "dev"

// runOptions holds the parsed command line options
type runOptions struct {
	generateConfig bool
	configOutput   string
	configFile     string
	inPlace        bool
	write          bool
	verbose        bool
	help           bool
	showVersion    bool
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
		printUsageTo(stderr)
		return err
	}

	// Handle help
	if opts.help {
		printUsageTo(stdout)
		return nil
	}

	// Handle version
	if opts.showVersion {
		fmt.Fprintf(stdout, "btfmt version %s\n", version)
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
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}
	cfg, err := config.LoadConfigFromWithLogger(cwd, opts.configFile, opts.verbose, stderr)
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
	fs.BoolVar(&opts.showVersion, "version", false, "Show version information")

	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}

	opts.files = fs.Args()
	return opts, nil
}

// readInputFile reads the content of a file
func readInputFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}
	return string(content), nil
}

// formatContent formats bpftrace script content using the provided configuration
func formatContent(content string, cfg *config.Config) (string, error) {
	var f formatter.Formatter = formatter.NewASTFormatter(cfg)
	formatted, err := f.Format(content)
	if err != nil {
		return "", fmt.Errorf("formatting: %w", err)
	}

	if !strings.HasSuffix(formatted, "\n") {
		formatted += "\n"
	}
	return formatted, nil
}

// processFile processes a single bpftrace file
func processFile(filename string, cfg *config.Config, writeToFile bool, verbose bool, stdout, stderr io.Writer) error {
	// Read input file
	content, err := readInputFile(filename)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Fprintf(stderr, "Processing: %s\n", filename)
	}

	// Format the content
	formatted, err := formatContent(content, cfg)
	if err != nil {
		return err
	}

	// Output result
	if writeToFile {
		// Write back to the original file without resetting permissions.
		if err := writeFilePreserveMode(filename, []byte(formatted), stderr); err != nil {
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

// createTempFile creates a temporary file, writes data to it, and syncs it
func createTempFile(dir, base string, data []byte, logWriter io.Writer) (string, *os.File, error) {
	tempFile, err := os.CreateTemp(dir, base+".tmp-*")
	if err != nil {
		return "", nil, err
	}
	tempName := tempFile.Name()

	n, err := tempFile.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err != nil {
		if closeErr := tempFile.Close(); closeErr != nil && logWriter != nil {
			fmt.Fprintf(logWriter, "Warning: failed to close temp file: %v\n", closeErr)
		}
		if removeErr := os.Remove(tempName); removeErr != nil && logWriter != nil {
			fmt.Fprintf(logWriter, "Warning: failed to remove temp file %s: %v\n", tempName, removeErr)
		}
		return "", nil, err
	}

	if err := tempFile.Sync(); err != nil {
		if closeErr := tempFile.Close(); closeErr != nil && logWriter != nil {
			fmt.Fprintf(logWriter, "Warning: failed to close temp file: %v\n", closeErr)
		}
		if removeErr := os.Remove(tempName); removeErr != nil && logWriter != nil {
			fmt.Fprintf(logWriter, "Warning: failed to remove temp file %s: %v\n", tempName, removeErr)
		}
		return "", nil, err
	}

	return tempName, tempFile, nil
}

// applyFilePermissions copies file permissions from info to the temp file
func applyFilePermissions(tempFile *os.File, info os.FileInfo) error {
	mode := info.Mode().Perm() | (info.Mode() & (os.ModeSetuid | os.ModeSetgid | os.ModeSticky))
	return tempFile.Chmod(mode)
}

// renameTempFileOverwrite renames tempFile over target file, handling Windows quirks
func renameTempFileOverwrite(tempName, filename string, logWriter io.Writer) error {
	err := os.Rename(tempName, filename)
	if err == nil {
		return nil
	}

	// On Windows, rename fails if target exists. Remove it first.
	if runtime.GOOS == "windows" {
		var renameErr error
		if removeErr := os.Remove(filename); removeErr == nil {
			renameErr = os.Rename(tempName, filename)
			if renameErr == nil {
				return nil
			}
			err = renameErr
		}
	}

	if removeErr := os.Remove(tempName); removeErr != nil && logWriter != nil {
		fmt.Fprintf(logWriter, "Warning: failed to remove temp file %s: %v\n", tempName, removeErr)
	}
	return err
}

func writeFilePreserveMode(filename string, data []byte, logWriter io.Writer) error {
	info, err := os.Lstat(filename)
	if err != nil {
		return err
	}

	if shouldWriteInPlace(info) {
		return writeFileTruncate(filename, data)
	}

	if err := ensureWritable(filename); err != nil {
		return err
	}

	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	tempName, tempFile, err := createTempFile(dir, base, data, logWriter)
	if err != nil {
		return writeFileTruncate(filename, data)
	}

	if err := applyFilePermissions(tempFile, info); err != nil {
		if closeErr := tempFile.Close(); closeErr != nil && logWriter != nil {
			fmt.Fprintf(logWriter, "Warning: failed to close temp file: %v\n", closeErr)
		}
		if removeErr := os.Remove(tempName); removeErr != nil && logWriter != nil {
			fmt.Fprintf(logWriter, "Warning: failed to remove temp file %s: %v\n", tempName, removeErr)
		}
		return err
	}

	if err := tempFile.Close(); err != nil {
		if removeErr := os.Remove(tempName); removeErr != nil && logWriter != nil {
			fmt.Fprintf(logWriter, "Warning: failed to remove temp file %s: %v\n", tempName, removeErr)
		}
		return err
	}

	return renameTempFileOverwrite(tempName, filename, logWriter)
}

func ensureWritable(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	return file.Close()
}

func shouldWriteInPlace(info os.FileInfo) bool {
	if info == nil {
		return false
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return true
	}
	if mismatch, ok := fileOwnerMismatch(info); ok && mismatch {
		return true
	}
	nlink, ok := fileLinkCount(info)
	if !ok {
		return true
	}
	return nlink > 1
}

func writeFileTruncate(filename string, data []byte) error {
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
  -version               Show version information
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
