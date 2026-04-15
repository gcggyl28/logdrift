// Package debounce provides a debouncer that delays processing of log entries
// until a quiet period has elapsed, collapsing rapid bursts into single events.
package debounce

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Entry represents a debounced log entry.
type Entry struct {
	Service string
	Line    string
}

// Debouncer delays forwarding of entries until no new entry arrives within
// the configured quiet window.
type Debouncer struct {
	quiet  time.Duration
	mu     sync.Mutex
	timers map[string]*time.Timer
	out    chan Entry
}

// New creates a Debouncer with the given quiet window.
// quiet must be positive.
func New(quiet time.Duration) (*Debouncer, error) {
	if quiet <= 0 {
		return nil, errors.New("debounce: quiet window must be positive")
	}
	return &Debouncer{
		quiet:  quiet,
		timers: make(map[string]*time.Timer),
		out:    make(chan Entry, 64),
	}, nil
}

// Push schedules e for forwarding after the quiet window.
// If another entry for the same service arrives before the window expires,
// the timer resets and only the latest entry is forwarded.
func (d *Debouncer) Push(e Entry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[e.Service]; ok {
		t.Reset(d.quiet)
		d.timers[e.Service] = time.AfterFunc(d.quiet, func() { d.emit(e) })
		return
	}
	d.timers[e.Service] = time.AfterFunc(d.quiet, func() { d.emit(e) })
}

func (d *Debouncer) emit(e Entry) {
	d.mu.Lock()
	delete(d.timers, e.Service)
	d.mu.Unlock()
	select {
	case d.out <- e:
	default:
	}
}

// Out returns the channel on which debounced entries are delivered.
func (d *Debouncer) Out() <-chan Entry { return d.out }

// Drain blocks until ctx is cancelled, forwarding debounced entries to fn.
func (d *Debouncer) Drain(ctx context.Context, fn func(Entry)) {
	for {
		select {
		case e := <-d.out:
			fn(e)
		case <-ctx.Done():
			return
		}
	}
}
