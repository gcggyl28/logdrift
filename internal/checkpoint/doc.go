// Package checkpoint provides durable read-position tracking for logdrift
// service log streams.
//
// A [Store] serialises per-service byte offsets to a JSON file so that
// logdrift can resume tailing from the correct position after a process
// restart, avoiding both duplicate processing and missed lines.
//
// A [Tracker] wraps a Store for a single service and accumulates byte
// progress via [Tracker.Advance]. In synchronous mode (interval == 0)
// each Advance call writes through to the Store immediately. In periodic
// mode a background goroutine started with [Tracker.Run] flushes on a
// configurable tick, with a guaranteed final flush on context cancellation.
//
// Typical usage:
//
//	store, err := checkpoint.New("/var/lib/logdrift/checkpoint.json")
//	tracker, err := checkpoint.NewTracker(store, "api", 5*time.Second)
//	go tracker.Run(ctx)
//	// For each line read:
//	tracker.Advance(int64(len(line)))
package checkpoint
