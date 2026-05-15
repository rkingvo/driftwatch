// Package differ compares a desired container state described by a
// manifest.ContainerSpec against the live state returned by the inspector,
// producing a manifest.DriftReport that lists every field that has diverged.
//
// Usage:
//
//	report := differ.Diff(spec, liveInfo)
//	if report.Drifted {
//		// handle drift
//	}
package differ
