// Package gapdetect provides time-gap detection across log streams.
//
// A Detector tracks the most-recent timestamp seen for each service.
// Whenever a new entry arrives and the elapsed time since the previous
// entry exceeds the configured threshold, the Detector returns a GapEvent
// describing the silence window.
//
// Typical usage:
//
//	det, err := gapdetect.New(30 * time.Second)
//	if err != nil { ... }
//
//	ev, err := det.Record(entry.Service, entry.Timestamp)
//	if ev != nil {
//	    log.Printf("gap detected: %s", ev.Summary())
//	}
//
// Each service is tracked independently; calling Reset clears the
// baseline for a single service without affecting others.
package gapdetect
