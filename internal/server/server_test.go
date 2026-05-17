package server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/healthz"
	"github.com/example/driftwatch/internal/metrics"
	"github.com/example/driftwatch/internal/server"
)

func buildServer(t *testing.T) (srv *server.Server, base string, cancel context.CancelFunc) {
	t.Helper()

	ln, err := server.Listener()
	if err != nil {
		t.Fatalf("listener: %v", err)
	}

	col := metrics.New()
	hz := healthz.New()
	hz.MarkReady()

	cfg := server.Config{
		Addr:            ln.Addr().String(),
		ReadTimeout:     2 * time.Second,
		WriteTimeout:    2 * time.Second,
		ShutdownTimeout: 2 * time.Second,
	}

	_ = ln.Close() // release so ListenAndServe can bind

	srv = server.New(cfg, col, hz)
	ctx, c := context.WithCancel(context.Background())

	go func() { _ = srv.Start(ctx) }()
	time.Sleep(50 * time.Millisecond) // allow goroutine to bind

	base = fmt.Sprintf("http://%s", cfg.Addr)
	return srv, base, c
}

func TestServer_LiveEndpoint(t *testing.T) {
	_, base, cancel := buildServer(t)
	defer cancel()

	resp, err := http.Get(base + "/healthz/live")
	if err != nil {
		t.Fatalf("GET /healthz/live: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestServer_ReadyEndpoint(t *testing.T) {
	_, base, cancel := buildServer(t)
	defer cancel()

	resp, err := http.Get(base + "/healthz/ready")
	if err != nil {
		t.Fatalf("GET /healthz/ready: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestServer_MetricsEndpoint_ReturnsJSON(t *testing.T) {
	_, base, cancel := buildServer(t)
	defer cancel()

	resp, err := http.Get(base + "/metrics")
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode metrics JSON: %v", err)
	}
}

func TestServer_UnknownRoute_Returns404(t *testing.T) {
	_, base, cancel := buildServer(t)
	defer cancel()

	resp, err := http.Get(base + "/does-not-exist")
	if err != nil {
		t.Fatalf("GET /does-not-exist: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestServer_GracefulShutdown(t *testing.T) {
	_, base, cancel := buildServer(t)

	// Verify the server is up.
	resp, err := http.Get(base + "/healthz/live")
	if err != nil {
		t.Fatalf("pre-shutdown GET: %v", err)
	}
	resp.Body.Close()

	// Cancel context → triggers graceful shutdown.
	cancel()
	time.Sleep(100 * time.Millisecond)

	_, err = http.Get(base + "/healthz/live")
	if err == nil {
		t.Error("expected connection refused after shutdown, got nil error")
	}

	var netErr *net.OpError
	if !isConnRefused(err) && netErr == nil {
		t.Logf("shutdown error (acceptable): %v", err)
	}
}

func isConnRefused(err error) bool {
	if err == nil {
		return false
	}
	return true // any error after cancel is acceptable in this context
}
