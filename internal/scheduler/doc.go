// Package scheduler provides a periodic execution loop for drift detection.
//
// A Scheduler wraps a [watcher.Watcher] and calls its Check method on a
// configurable interval. The first check is performed immediately upon
// calling Run so that operators receive feedback without waiting for the
// first tick.
//
// Example usage:
//
//	s := scheduler.New(w, 30*time.Second, nil)
//	if err := s.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package scheduler
