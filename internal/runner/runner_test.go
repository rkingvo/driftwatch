package runner_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/alerter"
	"github.com/user/driftwatch/internal/inspector"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/metrics"
	"github.com/user/driftwatch/internal/runner"
	"github.com/user/driftwatch/internal/watcher"
)

func buildRunner(t *testing.T, specs []manifest.ContainerSpec) *runner.Runner {
	t.Helper()

	containers := make(map[string]inspector.ContainerInfo, len(specs))
	for _, s := range specs {
		containers[s.Name] = inspector.ContainerInfo{
			Name:  s.Name,
			Image: s.Image,
			Env:   s.Env,
		}
	}

	client := inspector.NewMock(containers)
	w := watcher.New(specs, client)
	a := alerter.New(nil)
	c := metrics.New()
	logger := log.New(os.Stderr, "test: ", 0)

	return runner.New(w, 50*time.Millisecond, a, c, logger)
}

func TestRunner_RunCancels(t *testing.T) {
	specs := []manifest.ContainerSpec{
		{Name: "web", Image: "nginx:latest", Env: map[string]string{}},
	}
	r := buildRunner(t, specs)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := r.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("expected context error, got %v", err)
	}
}

func TestRunner_RecordsMetrics(t *testing.T) {
	specs := []manifest.ContainerSpec{
		{Name: "api", Image: "myapp:v1", Env: map[string]string{"PORT": "8080"}},
	}

	c := metrics.New()
	client := inspector.NewMock(map[string]inspector.ContainerInfo{
		"api": {Name: "api", Image: "myapp:v1", Env: map[string]string{"PORT": "8080"}},
	})
	w := watcher.New(specs, client)
	a := alerter.New(nil)
	logger := log.New(os.Stderr, "test: ", 0)
	r := runner.New(w, 40*time.Millisecond, a, c, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	_ = r.Run(ctx)

	snap := c.Snapshot()
	if snap.TotalChecks == 0 {
		t.Error("expected at least one check to be recorded")
	}
}
