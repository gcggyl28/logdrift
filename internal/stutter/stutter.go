// Package stutter detects repeated identical log lines emitted in rapid
// succession from the same service — a common symptom of tight error loops.
package stutter

import (
	"errors"
	"sync"
	"time"
)

// Event is emitted when a stutter condition is detected.
type Event struct {
	Service string
	Line    string
	Count   int
	First   time.Time
	Last    time.Time
}

type entry struct {
	line  string
	count int
	first time.Time
	last  time.Time
}

// Detector tracks repeated lines per service within a sliding window.
type Detector struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	state     map[string]*entry
}

// New creates a Detector that fires when the same line from a service
// repeats at least threshold times within window.
func New(window time.Duration, threshold int) (*Detector, error) {
	if window <= 0 {
		return nil, errors.New("stutter: window must be positive")
	}
	if threshold < 2 {
		return nil, errors.New("stutter: threshold must be at least 2")
	}
	return &Detector{
		window:    window,
		threshold: threshold,
		state:     make(map[string]*entry),
	}, nil
}

// Record records a line for the given service. It returns a non-nil Event
// exactly when the stutter threshold is first reached or on every subsequent
// occurrence within the same window.
func (d *Detector) Record(service, line string, at time.Time) *Event {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := service + "\x00" + line
	e, ok := d.state[key]

	if ok && at.Sub(e.last) > d.window {
		// window expired — reset
		ok = false
		delete(d.state, key)
	}

	if !ok {
		d.state[key] = &entry{line: line, count: 1, first: at, last: at}
		return nil
	}

	e.count++
	e.last = at

	if e.count >= d.threshold {
		return &Event{
			Service: service,
			Line:    line,
			Count:   e.count,
			First:   e.first,
			Last:    e.last,
		}
	}
	return nil
}

// Reset clears all tracked state for a service.
func (d *Detector) Reset(service string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for k := range d.state {
		if len(k) > len(service) && k[:len(service)] == service {
			delete(d.state, k)
		}
	}
}
