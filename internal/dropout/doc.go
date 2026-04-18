// Package dropout implements probabilistic load-shedding for log entry streams.
//
// A Dropper is constructed with a drop rate in [0.0, 1.0].  Each call to
// Allow() returns false with that probability, allowing callers to discard
// the entry and reduce downstream pressure.
//
// A rate of 0 disables dropping entirely (all entries pass through).
// A rate of 1 drops every entry.
//
// Example:
//
//	dr, err := dropout.New(0.2) // drop ~20 % of entries
//	if err != nil {
//		log.Fatal(err)
//	}
//	if dr.Allow() {
//		// forward entry
//	}
package dropout
