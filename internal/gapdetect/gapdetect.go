// Package gapdetect identifies time gaps between consecutive log entries
// for a given service, emitting a GapEvent whenever the silence exceeds
// a configurable threshold.
package gapdetect

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// GapEvent is emitted when a gap larger than the threshold is detected.
type GapEvent struct {
	Service  string
	From     time.Time
	To       time.Time
	Duration time.Duration
}

// Detector tracks the last-seen timestamp per service and reports gaps.
type Detector struct {
	mu        sync.Mutex
	threshold time.Duration
	last      map[string]time.Time
}

// New creates a Detector with the given gap threshold.
// threshold must be positive.
func New(threshold time.Duration) (*Detector, error) {
	if threshold <= 0 {
		return nil, errors.New("gapdetect: threshold must be positive")
	}
	return &Detector{
		threshold: threshold,
		last:      make(map[string]time.Time),
	}, nil
}

// Record records a log entry timestamp for service.
// If the gap since the previous entry exceeds the threshold, a non-nil
// *GapEvent is returned; otherwise nil is returned.
func (d *Detector) Record(service string, ts time.Time) (*GapEvent, error) {
	if service == "" {
		return nil, errors.New("gapdetect: service must not be empty")
	}
	if ts.IsZero() {
		return nil, errors.New("gapdetect: timestamp must not be zero")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	prev, seen := d.last[service]
	d.last[service] = ts

	if !seen {
		return nil, nil
	}

	gap := ts.Sub(prev)
	if gap < 0 {
		gap = -gap
	}
	if gap >= d.threshold {
		return &GapEvent{
			Service:  service,
			From:     prev,
			To:       ts,
			Duration: gap,
		}, nil
	}
	return nil, nil
}

// Reset clears the last-seen timestamp for service.
func (d *Detector) Reset(service string) error {
	if service == "" {
		return errors.New("gapdetect: service must not be empty")
	}
	d.mu.Lock()
	delete(d.last, service)
	d.mu.Unlock()
	return nil
}

// Summary returns a human-readable description of a GapEvent.
func (e *GapEvent) Summary() string {
	return fmt.Sprintf("[%s] gap of %s (from %s to %s)",
		e.Service, e.Duration.Round(time.Millisecond),
		e.From.Format(time.RFC3339), e.To.Format(time.RFC3339))
}
