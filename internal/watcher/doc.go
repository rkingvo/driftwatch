// Package watcher implements a periodic drift-detection loop.
//
// A Watcher loads the manifest file on every tick, inspects all running
// containers via the Docker client, computes a diff for each container
// listed in the manifest, and forwards the resulting DriftReports to a
// Reporter for output.
//
// Typical usage:
//
//	client, _ := inspector.New()
//	rep := reporter.New(os.Stdout, "text")
//	w := watcher.New("manifest.yaml", 30*time.Second, client, rep)
//	if err := w.Run(ctx); err != nil && err != context.Canceled {
//		log.Fatal(err)
//	}
package watcher
