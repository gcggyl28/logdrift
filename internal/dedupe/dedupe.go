// Package dedupe provides log-line deduplication by suppressing consecutive
// identical messages within a configurable window.
package dedupe

import (
	"errors"
	"sync"
	"time"
)

// Entry holds the last-seen state for a single service.
type entry struct {
	line  string
	count int
	seenAt time.Time
}

// Deduper suppresses repeated identical log lines within a time window.
type Deduper struct {
	window time.Duration
	mu     sync.Mutex
	state  map[string]*entry
}

// New creates a Deduper with the given suppression window.
// A zero window disables deduplication (every line passes through).
// A negative window returns an error.
func New(window time.Duration) (*Deduper, error) {
	if window < 0 {
		return nil, errors.New("dedupe: window must be >= 0")
	}
	return &Deduper{
		window: window,
		state:  make(map[string]*entry),
	}, nil
}

// Allow returns true if the line should be forwarded for the given service.
// When deduplication is disabled (window == 0) it always returns true.
func (d *Deduper) Allow(service, line string) bool {
	if d.window == 0 {
		return true
	}

	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()

	e, ok := d.state[service]
	if !ok || now.Sub(e.seenAt) > d.window || e.line != line {
		d.state[service] = &entry{line: line, count: 1, seenAt: now}
		return true
	}
	e.count++
	return false
}

// Count returns how many times the last line for a service was suppressed
// (i.e. seen but not forwarded). Returns 0 for unknown services.
func (d *Deduper) Count(service string) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	if e, ok := d.state[service]; ok {
		return e.count
	}
	return 0
}

// Reset clears all tracked state.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = make(map[string]*entry)
}
