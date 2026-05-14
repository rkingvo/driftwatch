package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ContainerSpec represents the desired state of a container as defined
// in a source manifest file.
type ContainerSpec struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Env         map[string]string `yaml:"env"`
	Ports       []string          `yaml:"ports"`
	Labels      map[string]string `yaml:"labels"`
	RestartPolicy string          `yaml:"restartPolicy"`
}

// Manifest holds one or more container specifications loaded from a file.
type Manifest struct {
	Version    string          `yaml:"version"`
	Containers []ContainerSpec `yaml:"containers"`
}

// LoadFromFile reads and parses a YAML manifest file at the given path.
func LoadFromFile(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("manifest: failed to read file %q: %w", path, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("manifest: failed to parse YAML in %q: %w", path, err)
	}

	if err := m.validate(); err != nil {
		return nil, fmt.Errorf("manifest: validation error in %q: %w", path, err)
	}

	return &m, nil
}

// validate performs basic sanity checks on the parsed manifest.
func (m *Manifest) validate() error {
	if len(m.Containers) == 0 {
		return fmt.Errorf("manifest must define at least one container")
	}
	seen := make(map[string]bool)
	for i, c := range m.Containers {
		if c.Name == "" {
			return fmt.Errorf("container at index %d is missing a name", i)
		}
		if c.Image == "" {
			return fmt.Errorf("container %q is missing an image", c.Name)
		}
		if seen[c.Name] {
			return fmt.Errorf("duplicate container name %q", c.Name)
		}
		seen[c.Name] = true
	}
	return nil
}
