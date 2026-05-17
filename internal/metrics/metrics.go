// Package metrics exposes Prometheus-style counters and gauges for
// driftwatch runtime telemetry.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of collected metrics.
type Snapshot struct {
	ChecksTotal    int64
	DriftTotal     int64
	LastCheckAt    time.Time
	LastDriftAt    time.Time
	ContainersDrifted []string
}

// Collector accumulates runtime statistics about drift checks.
type Collector struct {
	mu                sync.RWMutex
	checksTotal       int64
	driftTotal        int64
	lastCheckAt       time.Time
	lastDriftAt       time.Time
	containersDrifted []string
}

// New returns a zeroed Collector.
func New() *Collector {
	return &Collector{}
}

// RecordCheck increments the total check counter and updates the timestamp.
func (c *Collector) RecordCheck() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checksTotal++
	c.lastCheckAt = time.Now()
}

// RecordDrift increments the drift counter and records which containers drifted.
func (c *Collector) RecordDrift(containers []string) {
	if len(containers) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.driftTotal += int64(len(containers))
	c.lastDriftAt = time.Now()
	c.containersDrifted = append(c.containersDrifted, containers...)
}

// Snapshot returns an immutable copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	names := make([]string, len(c.containersDrifted))
	copy(names, c.containersDrifted)
	return Snapshot{
		ChecksTotal:       c.checksTotal,
		DriftTotal:        c.driftTotal,
		LastCheckAt:       c.lastCheckAt,
		LastDriftAt:       c.lastDriftAt,
		ContainersDrifted: names,
	}
}

// Reset clears all accumulated metrics.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checksTotal = 0
	c.driftTotal = 0
	c.lastCheckAt = time.Time{}
	c.lastDriftAt = time.Time{}
	c.containersDrifted = nil
}
