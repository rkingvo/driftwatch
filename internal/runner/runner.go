// Package runner wires together the scheduler, alerter, and metrics
// collector into a single long-running component that can be started and
// stopped via context cancellation.
package runner

import (
	"context"
	"log"
	"time"

	"github.com/user/driftwatch/internal/alerter"
	"github.com/user/driftwatch/internal/metrics"
	"github.com/user/driftwatch/internal/scheduler"
	"github.com/user/driftwatch/internal/watcher"
)

// Runner orchestrates periodic drift checks and forwards results to the
// alerter and metrics collector.
type Runner struct {
	sched     *scheduler.Scheduler
	alerter   *alerter.Alerter
	collector *metrics.Collector
	logger    *log.Logger
}

// New creates a Runner from the provided dependencies.
func New(
	w *watcher.Watcher,
	interval time.Duration,
	a *alerter.Alerter,
	c *metrics.Collector,
	logger *log.Logger,
) *Runner {
	s := scheduler.New(w, interval)
	return &Runner{
		sched:     s,
		alerter:   a,
		collector: c,
		logger:    logger,
	}
}

// Run starts the drift-check loop and blocks until ctx is cancelled.
// Each report produced by the scheduler is forwarded to the alerter and
// recorded in the metrics collector.
func (r *Runner) Run(ctx context.Context) error {
	reports, err := r.sched.Start(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case report, ok := <-reports:
			if !ok {
				return nil
			}
			r.collector.RecordCheck(report)
			if err := r.alerter.Notify(ctx, report); err != nil {
				r.logger.Printf("runner: alert notification failed: %v", err)
			}
		}
	}
}
