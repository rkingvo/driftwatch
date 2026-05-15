package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/manifest"
	"github.com/driftwatch/internal/reporter"
)

func noDriftReport() manifest.DriftReport {
	return manifest.DriftReport{HasDrift: false, Drifts: nil}
}

func driftReport() manifest.DriftReport {
	return manifest.DriftReport{
		HasDrift: true,
		Drifts: []manifest.ContainerDrift{
			{
				Container: "web",
				Fields: []manifest.FieldDrift{
					{Field: "image", Expected: "nginx:1.24", Actual: "nginx:1.25"},
					{Field: "env.PORT", Expected: "8080", Actual: "9090"},
				},
			},
		},
	}
}

func TestReporter_TextNoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)
	if err := r.Write(noDriftReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected 'No drift' in output, got: %s", buf.String())
	}
}

func TestReporter_TextWithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)
	if err := r.Write(driftReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"web", "image", "nginx:1.24", "nginx:1.25", "env.PORT"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in text output\n%s", want, out)
		}
	}
}

func TestReporter_JSONWithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatJSON, &buf)
	if err := r.Write(driftReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got manifest.DriftReport
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if !got.HasDrift {
		t.Error("expected HasDrift=true in JSON output")
	}
	if len(got.Drifts) != 1 || got.Drifts[0].Container != "web" {
		t.Errorf("unexpected drifts in JSON: %+v", got.Drifts)
	}
}

func TestReporter_JSONNoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatJSON, &buf)
	if err := r.Write(noDriftReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got manifest.DriftReport
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if got.HasDrift {
		t.Error("expected HasDrift=false")
	}
}
