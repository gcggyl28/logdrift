// Package alert implements threshold-based alerting for logdrift.
//
// An Alerter is created with a Config that specifies warn and error
// thresholds (as drift-line counts) and a rolling time window.  The
// caller is responsible for counting drifted lines within that window
// (typically via the snapshot package) and passing the count to
// Alerter.Evaluate.  When a threshold is breached, Evaluate returns a
// populated *Alert describing the breach level, affected service, and
// a human-readable message suitable for display in the terminal UI or
// writing to a log file.
//
// Typical usage:
//
//	cfg := alert.Config{
//		WarnThreshold:  10,
//		ErrorThreshold: 25,
//		Window:         time.Minute,
//	}
//	a, err := alert.New(cfg)
//	if err != nil { /* handle */ }
//	if al := a.Evaluate(service, driftCount); al != nil {
//		fmt.Println(al.Message)
//	}
package alert
