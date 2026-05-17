package alerter_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/alerter"
	"github.com/yourorg/driftwatch/internal/manifest"
)

// stubNotifier records calls and optionally returns an error.
type stubNotifier struct {
	called int
	err    error
}

func (s *stubNotifier) Notify(alerter.Alert) error {
	s.called++
	return s.err
}

func noDrift() manifest.DriftReport  { return manifest.DriftReport{} }
func hasDrift() manifest.DriftReport {
	return manifest.DriftReport{
		Containers: []manifest.ContainerDrift{
			{Name: "api", Diffs: []manifest.FieldDiff{{Field: "image", Want: "v1", Got: "v2"}}},
		},
	}
}

func TestAlerter_NoDrift_NoNotify(t *testing.T) {
	stub := &stubNotifier{}
	a := alerter.New(stub)
	if err := a.Send(noDrift()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.called != 0 {
		t.Errorf("expected 0 calls, got %d", stub.called)
	}
}

func TestAlerter_Drift_Notified(t *testing.T) {
	stub := &stubNotifier{}
	a := alerter.New(stub)
	if err := a.Send(hasDrift()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.called != 1 {
		t.Errorf("expected 1 call, got %d", stub.called)
	}
}

func TestAlerter_NotifierError_Propagated(t *testing.T) {
	stub := &stubNotifier{err: errors.New("boom")}
	a := alerter.New(stub)
	err := a.Send(hasDrift())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("error should mention underlying cause: %v", err)
	}
}

func TestLogNotifier_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	ln := alerter.NewLogNotifier(&buf)
	alrt := alerter.Alert{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Level:     alerter.LevelWarn,
		Report:    hasDrift(),
	}
	if err := ln.Notify(alrt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "1 container") {
		t.Errorf("expected container count in output, got: %s", out)
	}
}
