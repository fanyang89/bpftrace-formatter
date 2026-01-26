package lsp

import "github.com/tliron/glsp/server"

const (
	serverName    = "btfmt"
	serverVersion = "dev"
)

// RunServer starts the LSP server over stdio.
func RunServer() {
	handler := newHandler()
	server := server.NewServer(&handler, serverName, false)
	server.RunStdio()
}
