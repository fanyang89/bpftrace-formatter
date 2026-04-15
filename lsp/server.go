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

// Server represents the LSP server and its dependencies.
type Server struct {
	documentStore  DocumentStore
	configResolver ConfigProvider
	logger         *log.Logger
}

// NewServer creates a new Server instance.
func NewServer() *Server {
	s := &Server{
		logger: log.New(os.Stderr, "[LSP] ", log.Ltime|log.Lmicroseconds),
	}
	s.configResolver = NewConfigResolver()
	s.documentStore = NewDocumentStore(s.configResolver)
	return s
}

// Run starts the LSP server over stdio.
func (s *Server) Run() {
	s.logger.Print("server starting")

	handler := s.newHandler()
	srv := server.NewServer(&handler, serverName, false)
	s.logger.Print("handler registered, running stdio")
	srv.RunStdio()
}

// RunServer is a convenience wrapper for backward compatibility.
func RunServer() {
	NewServer().Run()
}
