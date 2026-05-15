package differ

import (
	"fmt"

	"github.com/yourorg/driftwatch/internal/inspector"
	"github.com/yourorg/driftwatch/internal/manifest"
)

// Diff compares a ContainerSpec from the manifest against live ContainerInfo
// from the inspector and returns a DriftReport describing any differences.
func Diff(spec manifest.ContainerSpec, live inspector.ContainerInfo) manifest.DriftReport {
	report := manifest.DriftReport{
		ContainerName: spec.Name,
		Drifted:       false,
		Diffs:         []manifest.DriftEntry{},
	}

	if spec.Image != live.Image {
		report.Drifted = true
		report.Diffs = append(report.Diffs, manifest.DriftEntry{
			Field:    "image",
			Expected: spec.Image,
			Actual:   live.Image,
		})
	}

	expectedEnv := toEnvMap(spec.Env)
	actualEnv := toEnvMap(live.Env)

	for k, ev := range expectedEnv {
		av, ok := actualEnv[k]
		if !ok {
			report.Drifted = true
			report.Diffs = append(report.Diffs, manifest.DriftEntry{
				Field:    fmt.Sprintf("env.%s", k),
				Expected: ev,
				Actual:   "<missing>",
			})
		} else if av != ev {
			report.Drifted = true
			report.Diffs = append(report.Diffs, manifest.DriftEntry{
				Field:    fmt.Sprintf("env.%s", k),
				Expected: ev,
				Actual:   av,
			})
		}
	}

	return report
}

// toEnvMap converts a slice of "KEY=VALUE" strings into a map.
func toEnvMap(pairs []string) map[string]string {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		for i := 0; i < len(p); i++ {
			if p[i] == '=' {
				m[p[:i]] = p[i+1:]
				break
			}
		}
	}
	return m
}
