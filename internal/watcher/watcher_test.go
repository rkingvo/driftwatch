package watcher_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/inspector"
	"github.com/yourorg/driftwatch/internal/manifest"
	"github.com/yourorg/driftwatch/internal/reporter"
	"github.com/yourorg/driftwatch/internal/watcher"
)

func writeTempManifest(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "manifest-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

const sampleManifest = `containers:
  - name: web
    image: nginx:1.25
    env:
      PORT: "8080"
`

func TestWatcher_CheckNoDrift(t *testing.T) {
	path := writeTempManifest(t, sampleManifest)

	mock := inspector.NewMock(map[string]manifest.ContainerInfo{
		"web": {
			Name:  "web",
			Image: "nginx:1.25",
			Env:   map[string]string{"PORT": "8080"},
		},
	})

	var buf bytes.Buffer
	rep := reporter.New(&buf, "text")

	w := watcher.New(path, 10*time.Second, mock, rep)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected reporter output, got none")
	}
}

func TestWatcher_MissingContainer(t *testing.T) {
	path := writeTempManifest(t, sampleManifest)

	// Live state has no containers — watcher should log and not panic.
	mock := inspector.NewMock(map[string]manifest.ContainerInfo{})

	var buf bytes.Buffer
	rep := reporter.New(&buf, "text")

	w := watcher.New(path, 10*time.Second, mock, rep)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}
