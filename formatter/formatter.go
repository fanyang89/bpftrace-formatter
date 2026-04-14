package formatter

// Formatter defines the interface for formatting bpftrace scripts
type Formatter interface {
	Format(input string) (string, error)
}
