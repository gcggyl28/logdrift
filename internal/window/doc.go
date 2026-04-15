// Package window provides a sliding time-window counter used to measure
// event rates over a rolling duration.
//
// # Overview
//
// A Window accumulates timestamped event counts and automatically evicts
// entries that fall outside the configured duration on every read or write.
// All operations are safe for concurrent use.
//
// # Usage
//
//	w, err := window.New(30 * time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	w.Add(1)          // record one event now
//	total := w.Total() // sum of events in the last 30 s
package window
