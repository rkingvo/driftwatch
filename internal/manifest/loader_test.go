package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp manifest: %v", err)
	}
	return path
}

func TestLoadFromFile_Valid(t *testing.T) {
	raw := `
version: "1"
containers:
  - name: web
    image: nginx:1.25
    ports:
      - "80:80"
    env:
      APP_ENV: production
    labels:
      team: platform
    restartPolicy: always
`
	path := writeTemp(t, raw)
	m, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(m.Containers))
	}
	c := m.Containers[0]
	if c.Name != "web" {
		t.Errorf("expected name 'web', got %q", c.Name)
	}
	if c.Image != "nginx:1.25" {
		t.Errorf("expected image 'nginx:1.25', got %q", c.Image)
	}
	if c.Env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", c.Env["APP_ENV"])
	}
}

func TestLoadFromFile_MissingFile(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/manifest.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFromFile_EmptyContainers(t *testing.T) {
	raw := `version: "1"\ncontainers: []\n`
	path := writeTemp(t, raw)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected validation error for empty containers")
	}
}

func TestLoadFromFile_DuplicateNames(t *testing.T) {
	raw := `
version: "1"
containers:
  - name: api
    image: myapp:latest
  - name: api
    image: myapp:v2
`
	path := writeTemp(t, raw)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected validation error for duplicate container names")
	}
}

func TestLoadFromFile_MissingImage(t *testing.T) {
	raw := `
version: "1"
containers:
  - name: worker
`
	path := writeTemp(t, raw)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected validation error for missing image")
	}
}
