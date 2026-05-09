// internal/sandbox/server.go
package sandbox

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/testmind-hq/caseforge/internal/spec"
)

// SandboxServer is a local HTTP server that mocks an OpenAPI-described API.
type SandboxServer struct {
	srv        *http.Server
	addr       string // resolved after Start()
	logCleanup func() // closes log file(s) opened by buildLogger
}

// NewSandboxServer builds (but does not start) a SandboxServer for the given spec.
func NewSandboxServer(ps *spec.ParsedSpec, opts Options) *SandboxServer {
	logger, cleanup := buildLogger(opts)

	gen := newResponseGenerator(opts.formatOrDefault())
	store := newMemStateStore()

	mux := http.NewServeMux()
	for _, op := range ps.Operations {
		pattern := fmt.Sprintf("%s %s", strings.ToUpper(op.Method), op.Path)
		mux.HandleFunc(pattern, makeHandler(op, store, gen, logger))
	}

	srv := &http.Server{Handler: mux}

	return &SandboxServer{
		srv:        srv,
		logCleanup: cleanup,
	}
}

// Start listens on the given host:port and begins serving in a goroutine.
// When port is 0, the OS assigns a random available port.
// Addr() returns the actual address after Start returns.
func (s *SandboxServer) Start(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("sandbox listen: %w", err)
	}
	s.addr = ln.Addr().String()
	go s.srv.Serve(ln) //nolint:errcheck
	return nil
}

// Addr returns "host:port" of the listening server. Call after Start.
func (s *SandboxServer) Addr() string { return s.addr }

// Shutdown gracefully stops the server within the given context deadline.
func (s *SandboxServer) Shutdown(ctx context.Context) error {
	defer s.logCleanup()
	return s.srv.Shutdown(ctx)
}
