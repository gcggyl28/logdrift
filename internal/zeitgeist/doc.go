// Package zeitgeist provides time-bucketing of log entries into fixed-width
// slots. Each slot accumulates per-service line counts, enabling callers to
// observe throughput trends over a rolling window.
//
// Usage:
//
//	bk, err := zeitgeist.New(time.Minute, 60)
//	if err != nil { ... }
//	bk.Record("api", entry.Timestamp)
//	for _, slot := range bk.Snapshot() {
//		fmt.Println(slot.Start, slot.Counts)
//	}
package zeitgeist
