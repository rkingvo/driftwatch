// Package server wires together the HTTP endpoints exposed by driftwatch,
// including metrics, health checks, and a status summary endpoint.
package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/example/driftwatch/internal/healthz"
	"github.com/example/driftwatch/internal/metrics"
)

// Server is a self-contained HTTP server that aggregates all driftwatch
// diagnostic endpoints under a single listener.
type Server struct {
	httpServer *http.Server
}

// Config holds the options used to construct a Server.
type Config struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Addr:            ":8080",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 15 * time.Second,
	}
}

// New constructs a Server using the provided Config, metrics Collector, and
// healthz Handler pair.
func New(cfg Config, col *metrics.Collector, hz *healthz.Health) *Server {
	mux := http.NewServeMux()

	mux.Handle("/metrics", metrics.Handler(col))
	mux.Handle("/healthz/live", hz.LiveHandler())
	mux.Handle("/healthz/ready", hz.ReadyHandler())

	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Addr,
			Handler:      mux,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

// Start begins listening on the configured address. It blocks until the
// context is cancelled, then performs a graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutCtx)
	}
}

// Addr returns the address the underlying http.Server is configured to use.
func (s *Server) Addr() string {
	return s.httpServer.Addr
}

// Listener returns a net.Listener bound to an OS-assigned port. Useful in
// tests where a fixed port is undesirable.
func Listener() (net.Listener, error) {
	return net.Listen("tcp", "127.0.0.1:0")
}
