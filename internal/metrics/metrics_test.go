package metrics_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/metrics"
)

func TestCollector_InitialSnapshot_IsZero(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()
	if s.ChecksTotal != 0 || s.DriftTotal != 0 {
		t.Fatalf("expected zero snapshot, got %+v", s)
	}
	if !s.LastCheckAt.IsZero() || !s.LastDriftAt.IsZero() {
		t.Fatal("expected zero timestamps")
	}
}

func TestCollector_RecordCheck(t *testing.T) {
	c := metrics.New()
	before := time.Now()
	c.RecordCheck()
	c.RecordCheck()
	s := c.Snapshot()
	if s.ChecksTotal != 2 {
		t.Fatalf("expected 2 checks, got %d", s.ChecksTotal)
	}
	if s.LastCheckAt.Before(before) {
		t.Fatal("LastCheckAt should be after test start")
	}
}

func TestCollector_RecordDrift(t *testing.T) {
	c := metrics.New()
	c.RecordDrift([]string{"web", "db"})
	s := c.Snapshot()
	if s.DriftTotal != 2 {
		t.Fatalf("expected drift total 2, got %d", s.DriftTotal)
	}
	if len(s.ContainersDrifted) != 2 {
		t.Fatalf("expected 2 drifted containers, got %d", len(s.ContainersDrifted))
	}
}

func TestCollector_RecordDrift_Empty_NoOp(t *testing.T) {
	c := metrics.New()
	c.RecordDrift([]string{})
	s := c.Snapshot()
	if s.DriftTotal != 0 {
		t.Fatal("empty drift slice should not increment counter")
	}
	if !s.LastDriftAt.IsZero() {
		t.Fatal("LastDriftAt should remain zero")
	}
}

func TestCollector_Reset(t *testing.T) {
	c := metrics.New()
	c.RecordCheck()
	c.RecordDrift([]string{"api"})
	c.Reset()
	s := c.Snapshot()
	if s.ChecksTotal != 0 || s.DriftTotal != 0 {
		t.Fatalf("expected zeroed snapshot after Reset, got %+v", s)
	}
	if len(s.ContainersDrifted) != 0 {
		t.Fatal("ContainersDrifted should be empty after Reset")
	}
}

func TestCollector_Snapshot_IsCopy(t *testing.T) {
	c := metrics.New()
	c.RecordDrift([]string{"svc"})
	s := c.Snapshot()
	s.ContainersDrifted = append(s.ContainersDrifted, "mutated")
	s2 := c.Snapshot()
	if len(s2.ContainersDrifted) != 1 {
		t.Fatal("mutating snapshot slice should not affect collector state")
	}
}
