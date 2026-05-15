package differ_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/differ"
	"github.com/yourorg/driftwatch/internal/inspector"
	"github.com/yourorg/driftwatch/internal/manifest"
)

func baseSpec() manifest.ContainerSpec {
	return manifest.ContainerSpec{
		Name:  "web",
		Image: "nginx:1.25",
		Env:   []string{"PORT=8080", "DEBUG=false"},
	}
}

func baseLive() inspector.ContainerInfo {
	return inspector.ContainerInfo{
		Name:  "web",
		Image: "nginx:1.25",
		Env:   []string{"PORT=8080", "DEBUG=false"},
	}
}

func TestDiff_NoDrift(t *testing.T) {
	report := differ.Diff(baseSpec(), baseLive())
	if report.Drifted {
		t.Errorf("expected no drift, got diffs: %v", report.Diffs)
	}
	if len(report.Diffs) != 0 {
		t.Errorf("expected empty diffs, got %d entries", len(report.Diffs))
	}
}

func TestDiff_ImageChanged(t *testing.T) {
	live := baseLive()
	live.Image = "nginx:1.26"

	report := differ.Diff(baseSpec(), live)
	if !report.Drifted {
		t.Fatal("expected drift to be detected")
	}
	if len(report.Diffs) != 1 || report.Diffs[0].Field != "image" {
		t.Errorf("unexpected diffs: %v", report.Diffs)
	}
}

func TestDiff_EnvValueChanged(t *testing.T) {
	live := baseLive()
	live.Env = []string{"PORT=9090", "DEBUG=false"}

	report := differ.Diff(baseSpec(), live)
	if !report.Drifted {
		t.Fatal("expected drift")
	}
	found := false
	for _, d := range report.Diffs {
		if d.Field == "env.PORT" && d.Expected == "8080" && d.Actual == "9090" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected env.PORT diff, got: %v", report.Diffs)
	}
}

func TestDiff_EnvKeyMissing(t *testing.T) {
	live := baseLive()
	live.Env = []string{"PORT=8080"} // DEBUG missing

	report := differ.Diff(baseSpec(), live)
	if !report.Drifted {
		t.Fatal("expected drift")
	}
	found := false
	for _, d := range report.Diffs {
		if d.Field == "env.DEBUG" && d.Actual == "<missing>" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected env.DEBUG missing diff, got: %v", report.Diffs)
	}
}

func TestDiff_ContainerNamePropagated(t *testing.T) {
	report := differ.Diff(baseSpec(), baseLive())
	if report.ContainerName != "web" {
		t.Errorf("expected ContainerName=web, got %q", report.ContainerName)
	}
}
