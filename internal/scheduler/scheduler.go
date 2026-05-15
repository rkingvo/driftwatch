// Package scheduler provides periodic drift-check execution.
package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/driftwatch/internal/watcher"
)

// Scheduler runs a drift check on a fixed interval.
type Scheduler struct {
	w        *watcher.Watcher
	interval time.Duration
	logger   *log.Logger
}

// New creates a Scheduler that will invoke w.Check every interval.
func New(w *watcher.Watcher, interval time.Duration, logger *log.Logger) *Scheduler {
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		w:        w,
		interval: interval,
		logger:   logger,
	}
}

// Run blocks and executes drift checks until ctx is cancelled.
// It performs an immediate check before entering the tick loop.
func (s *Scheduler) Run(ctx context.Context) error {
	s.tick(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.tick(ctx)
		case <-ctx.Done():
			s.logger.Println("scheduler: stopping")
			return ctx.Err()
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	reports, err := s.w.Check(ctx)
	if err != nil {
		s.logger.Printf("scheduler: check error: %v", err)
		return
	}
	for _, r := range reports {
		if r.HasDrift {
			s.logger.Printf("scheduler: drift detected in container %q", r.ContainerName)
		}
	}
}
