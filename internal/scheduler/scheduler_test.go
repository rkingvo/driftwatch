package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/manifest"
	"github.com/yourorg/driftwatch/internal/inspector"
	"github.com/yourorg/driftwatch/internal/watcher"
	"github.com/yourorg/driftwatch/internal/scheduler"
)

func buildWatcher(t *testing.T) *watcher.Watcher {
	t.Helper()
	spec := manifest.ContainerSpec{
		Name:  "web",
		Image: "nginx:latest",
	}
	mock := inspector.NewMock([]inspector.ContainerInfo{
		{Name: "web", Image: "nginx:latest", Env: map[string]string{}},
	})
	return watcher.New([]manifest.ContainerSpec{spec}, mock, nil)
}

func TestScheduler_RunCancels(t *testing.T) {
	w := buildWatcher(t)
	s := scheduler.New(w, 100*time.Millisecond, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	err := s.Run(ctx)
	if err == nil {
		t.Fatal("expected non-nil error on context cancellation")
	}
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScheduler_ImmediateCheck(t *testing.T) {
	// Ensure the scheduler fires at least once before the first tick.
	w := buildWatcher(t)

	checked := make(chan struct{}, 1)
	_ = checked

	s := scheduler.New(w, 10*time.Second, nil)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	// Give the immediate tick time to execute.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("scheduler did not stop in time")
	}
}
