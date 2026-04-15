// Package linecount provides a sliding-window line-rate counter for
// logdrift services.
//
// Each service's line events are timestamped and stored in a ring of
// timestamps. Entries older than the configured window are evicted lazily
// on every Record or Rate call, keeping memory proportional to the actual
// event rate rather than to the total number of lines ever seen.
//
// Typical usage:
//
//	c, err := linecount.New(10 * time.Second)
//	if err != nil { ... }
//
//	// in your log-ingestion loop:
//	c.Record(entry.Service)
//
//	// periodically report:
//	fmt.Printf("auth: %.2f lines/sec\n", c.Rate("auth"))
package linecount
