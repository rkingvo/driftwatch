// Package metrics provides a lightweight, dependency-free telemetry collector
// for driftwatch.
//
// Usage:
//
//	col := metrics.New()
//
//	// After each watcher cycle:
//	col.RecordCheck()
//	col.RecordDrift(driftedNames)
//
//	// Expose over HTTP:
//	http.Handle("/metrics", metrics.Handler(col))
//
// All methods are safe for concurrent use. The Snapshot type is immutable;
// callers may inspect it without holding any lock.
package metrics
