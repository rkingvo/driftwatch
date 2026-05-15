package manifest

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDriftReport_JSONRoundtrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	report := DriftReport{
		ContainerName: "web",
		Drifted:       true,
		CheckedAt:     now,
		Fields: []DriftField{
			{Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"},
		},
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded DriftReport
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.ContainerName != report.ContainerName {
		t.Errorf("expected container name %q, got %q", report.ContainerName, decoded.ContainerName)
	}
	if !decoded.Drifted {
		t.Error("expected drifted to be true")
	}
	if len(decoded.Fields) != 1 {
		t.Fatalf("expected 1 drift field, got %d", len(decoded.Fields))
	}
	if decoded.Fields[0].Field != "image" {
		t.Errorf("expected field name 'image', got %q", decoded.Fields[0].Field)
	}
}

func TestDriftReport_NoDrift(t *testing.T) {
	report := DriftReport{
		ContainerName: "worker",
		Drifted:       false,
		CheckedAt:     time.Now(),
	}

	if len(report.Fields) != 0 {
		t.Errorf("expected no drift fields, got %d", len(report.Fields))
	}
}

func TestContainerSpec_Defaults(t *testing.T) {
	spec := ContainerSpec{
		Name:  "api",
		Image: "myapp:latest",
	}

	if spec.RestartPolicy != "" {
		t.Errorf("expected empty restart policy, got %q", spec.RestartPolicy)
	}
	if spec.Env != nil {
		t.Error("expected nil env map")
	}
	if len(spec.Ports) != 0 {
		t.Errorf("expected no ports, got %d", len(spec.Ports))
	}
}
