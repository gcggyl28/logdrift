// Package dedupe implements log-line deduplication for logdrift.
//
// A Deduper suppresses consecutive identical log lines emitted by the same
// service within a configurable time window. This prevents noisy, repetitive
// log entries from flooding the diff pipeline and obscuring meaningful drift.
//
// # Usage
//
//	d, err := dedupe.New(5 * time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if d.Allow("api", line) {
//		// forward line to the diff pipeline
//	}
//
// Setting the window to zero disables deduplication entirely; every line
// passes through. A negative window is rejected with an error.
package dedupe
