package manifest

import "time"

// ContainerSpec defines the expected configuration for a single container.
type ContainerSpec struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Env         map[string]string `yaml:"env,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	RestartPolicy string          `yaml:"restart_policy,omitempty"`
}

// Manifest represents the top-level structure of a driftwatch manifest file.
type Manifest struct {
	Version    string          `yaml:"version"`
	Containers []ContainerSpec `yaml:"containers"`
}

// DriftField describes a single field that has drifted from its expected value.
type DriftField struct {
	Field    string `json:"field"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}

// DriftReport holds the result of comparing a running container against its spec.
type DriftReport struct {
	ContainerName string       `json:"container_name"`
	Drifted       bool         `json:"drifted"`
	Fields        []DriftField `json:"fields,omitempty"`
	CheckedAt     time.Time    `json:"checked_at"`
}
