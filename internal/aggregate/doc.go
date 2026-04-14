// Package aggregate provides per-key statistical aggregation over log streams.
//
// An Aggregator groups log entries by an arbitrary string key (typically a
// service name or label) and tracks line counts, drift counts, and time
// boundaries within an optional sliding window.
//
// Usage:
//
//	agg, err := aggregate.New(5 * time.Minute)
//	if err != nil { ... }
//
//	agg.RecordLine("auth-service", time.Now())
//	agg.RecordDrift("auth-service", time.Now())
//
//	stats := agg.Snapshot()
//	for _, s := range stats {
//		fmt.Printf("%s: lines=%d drifts=%d\n", s.Key, s.Lines, s.Drifts)
//	}
//
// When window is zero, no eviction occurs and all keys are retained until
// Reset is called.
package aggregate
