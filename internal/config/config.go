// Package config loads and validates driftwatch daemon configuration
// from environment variables and/or a YAML/JSON config file.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration for the driftwatch daemon.
type Config struct {
	// ManifestPath is the path to the container manifest file.
	ManifestPath string

	// Interval is how often the watcher checks for drift.
	Interval time.Duration

	// WebhookURL is an optional URL to POST drift alerts to.
	WebhookURL string

	// OutputFormat controls reporter output: "text" or "json".
	OutputFormat string

	// DockerHost overrides the Docker socket/host (optional).
	DockerHost string
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		ManifestPath: "manifest.yaml",
		Interval:     30 * time.Second,
		OutputFormat: "text",
	}
}

// LoadFromEnv reads configuration from environment variables,
// falling back to defaults for any unset values.
func LoadFromEnv() (Config, error) {
	cfg := DefaultConfig()

	if v := os.Getenv("DRIFTWATCH_MANIFEST"); v != "" {
		cfg.ManifestPath = v
	}

	if v := os.Getenv("DRIFTWATCH_INTERVAL"); v != "" {
		secs, err := strconv.Atoi(v)
		if err != nil {
			return cfg, fmt.Errorf("config: DRIFTWATCH_INTERVAL must be an integer (seconds): %w", err)
		}
		if secs <= 0 {
			return cfg, fmt.Errorf("config: DRIFTWATCH_INTERVAL must be positive, got %d", secs)
		}
		cfg.Interval = time.Duration(secs) * time.Second
	}

	if v := os.Getenv("DRIFTWATCH_WEBHOOK_URL"); v != "" {
		cfg.WebhookURL = v
	}

	if v := os.Getenv("DRIFTWATCH_OUTPUT_FORMAT"); v != "" {
		if v != "text" && v != "json" {
			return cfg, fmt.Errorf("config: DRIFTWATCH_OUTPUT_FORMAT must be \"text\" or \"json\", got %q", v)
		}
		cfg.OutputFormat = v
	}

	if v := os.Getenv("DOCKER_HOST"); v != "" {
		cfg.DockerHost = v
	}

	return cfg, nil
}
