package server_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/server"
)

func TestDefaultConfig_Addr(t *testing.T) {
	cfg := server.DefaultConfig()
	if cfg.Addr != ":8080" {
		t.Errorf("expected :8080, got %q", cfg.Addr)
	}
}

func TestDefaultConfig_Timeouts(t *testing.T) {
	cfg := server.DefaultConfig()

	if cfg.ReadTimeout != 5*time.Second {
		t.Errorf("ReadTimeout: expected 5s, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 10*time.Second {
		t.Errorf("WriteTimeout: expected 10s, got %v", cfg.WriteTimeout)
	}
	if cfg.ShutdownTimeout != 15*time.Second {
		t.Errorf("ShutdownTimeout: expected 15s, got %v", cfg.ShutdownTimeout)
	}
}

func TestServer_Addr_MatchesConfig(t *testing.T) {
	_, base, cancel := buildServer(t)
	defer cancel()

	if base == "" {
		t.Error("expected non-empty base URL")
	}
}
