package lsp

import (
	"log"
	"os"

	"github.com/tliron/glsp/server"
)

const (
	serverName    = "btfmt"
	serverVersion = "dev"
)

// RunServer starts the LSP server over stdio.
func RunServer() {
	// Log to stderr so VSCode can capture it
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Print("[LSP] server starting")

	handler := newHandler()
	srv := server.NewServer(&handler, serverName, false)
	log.Print("[LSP] handler registered, running stdio")
	srv.RunStdio()
}
