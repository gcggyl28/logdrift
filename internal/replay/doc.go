// Package replay provides a Replayer that re-emits log lines stored in a
// [snapshot.Snapshot] over a channel, enabling offline or delayed drift
// analysis without a live log source.
//
// # Usage
//
//	snap := snapshot.New(200)
//	// ... populate snap via a Collector ...
//
//	r := replay.New(snap, 0) // 0 delay = as fast as possible
//	for entry := range r.Run(ctx) {
//		fmt.Printf("[%s] %s\n", entry.Service, entry.Line)
//	}
//
// Lines are emitted in round-robin order across services so that the diff
// pipeline receives interleaved input similar to a live stream.
package replay
