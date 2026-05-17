package config_test

import (
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.ManifestPath != "manifest.yaml" {
		t.Errorf("expected manifest.yaml, got %q", cfg.ManifestPath)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected text output format, got %q", cfg.OutputFormat)
	}
}

func TestLoadFromEnv_Defaults(t *testing.T) {
	clearEnv(t)

	cfg, err := config.LoadFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected default 30s, got %v", cfg.Interval)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected default text, got %q", cfg.OutputFormat)
	}
}

func TestLoadFromEnv_CustomValues(t *testing.T) {
	clearEnv(t)
	t.Setenv("DRIFTWATCH_MANIFEST", "/etc/driftwatch/manifest.yaml")
	t.Setenv("DRIFTWATCH_INTERVAL", "60")
	t.Setenv("DRIFTWATCH_WEBHOOK_URL", "https://hooks.example.com/drift")
	t.Setenv("DRIFTWATCH_OUTPUT_FORMAT", "json")
	t.Setenv("DOCKER_HOST", "tcp://localhost:2375")

	cfg, err := config.LoadFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ManifestPath != "/etc/driftwatch/manifest.yaml" {
		t.Errorf("unexpected manifest path: %q", cfg.ManifestPath)
	}
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected 60s, got %v", cfg.Interval)
	}
	if cfg.WebhookURL != "https://hooks.example.com/drift" {
		t.Errorf("unexpected webhook url: %q", cfg.WebhookURL)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected json, got %q", cfg.OutputFormat)
	}
	if cfg.DockerHost != "tcp://localhost:2375" {
		t.Errorf("unexpected docker host: %q", cfg.DockerHost)
	}
}

func TestLoadFromEnv_InvalidInterval(t *testing.T) {
	clearEnv(t)
	t.Setenv("DRIFTWATCH_INTERVAL", "not-a-number")

	_, err := config.LoadFromEnv()
	if err == nil {
		t.Fatal("expected error for invalid interval, got nil")
	}
}

func TestLoadFromEnv_ZeroInterval(t *testing.T) {
	clearEnv(t)
	t.Setenv("DRIFTWATCH_INTERVAL", "0")

	_, err := config.LoadFromEnv()
	if err == nil {
		t.Fatal("expected error for zero interval, got nil")
	}
}

func TestLoadFromEnv_InvalidOutputFormat(t *testing.T) {
	clearEnv(t)
	t.Setenv("DRIFTWATCH_OUTPUT_FORMAT", "xml")

	_, err := config.LoadFromEnv()
	if err == nil {
		t.Fatal("expected error for invalid output format, got nil")
	}
}

// clearEnv unsets all driftwatch-related environment variables for test isolation.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"DRIFTWATCH_MANIFEST",
		"DRIFTWATCH_INTERVAL",
		"DRIFTWATCH_WEBHOOK_URL",
		"DRIFTWATCH_OUTPUT_FORMAT",
		"DOCKER_HOST",
	} {
		t.Setenv(key, "")
	}
}
