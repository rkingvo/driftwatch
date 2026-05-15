// Package watcher provides periodic drift detection by comparing
// running containers against their source manifests on a schedule.
package watcher

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/driftwatch/internal/differ"
	"github.com/yourorg/driftwatch/internal/inspector"
	"github.com/yourorg/driftwatch/internal/manifest"
	"github.com/yourorg/driftwatch/internal/reporter"
)

// Watcher polls for config drift at a fixed interval.
type Watcher struct {
	manifestPath string
	interval     time.Duration
	inspector    inspector.Client
	reporter     *reporter.Reporter
}

// New creates a Watcher with the given manifest path, poll interval,
// Docker client, and reporter.
func New(manifestPath string, interval time.Duration, client inspector.Client, rep *reporter.Reporter) *Watcher {
	return &Watcher{
		manifestPath: manifestPath,
		interval:     interval,
		inspector:    client,
		reporter:     rep,
	}
}

// Run starts the watch loop. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run once immediately before waiting for the first tick.
	if err := w.check(ctx); err != nil {
		log.Printf("watcher: check error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.check(ctx); err != nil {
				log.Printf("watcher: check error: %v", err)
			}
		}
	}
}

// check loads the manifest, inspects live containers, diffs each one,
// and writes a report via the reporter.
func (w *Watcher) check(ctx context.Context) error {
	spec, err := manifest.LoadFromFile(w.manifestPath)
	if err != nil {
		return err
	}

	live, err := w.inspector.InspectAll(ctx)
	if err != nil {
		return err
	}

	var reports []manifest.DriftReport
	for _, cs := range spec.Containers {
		info, ok := live[cs.Name]
		if !ok {
			log.Printf("watcher: container %q not found in live state", cs.Name)
			continue
		}
		reports = append(reports, differ.Diff(cs, info))
	}

	return w.reporter.Write(reports)
}
