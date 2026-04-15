// Package coalesce merges log entries from multiple services that arrive
// within a short time window into a single grouped event. This is useful
// for correlating related log lines that are emitted nearly simultaneously
// across services (e.g. a shared request ID or timestamp bucket).
package coalesce

import (
	"errors"
	"sync"
	"time"
)

// Entry represents a single log line attributed to a named service.
type Entry struct {
	Service string
	Line    string
	At      time.Time
}

// Group is a set of entries coalesced within one time window.
type Group []Entry

// Coalescer buffers incoming entries and flushes them as a Group once the
// window elapses with no new arrivals, or when Flush is called explicitly.
type Coalescer struct {
	mu      sync.Mutex
	window  time.Duration
	buf     []Entry
	timer   *time.Timer
	out     chan Group
	closed  bool
}

// New creates a Coalescer with the given idle window duration.
// window must be positive. The returned channel receives flushed groups.
func New(window time.Duration) (*Coalescer, <-chan Group, error) {
	if window <= 0 {
		return nil, nil, errors.New("coalesce: window must be positive")
	}
	c := &Coalescer{
		window: window,
		out:    make(chan Group, 64),
	}
	return c, c.out, nil
}

// Push adds an entry to the current window buffer, resetting the idle timer.
func (c *Coalescer) Push(e Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	if e.At.IsZero() {
		e.At = time.Now()
	}
	c.buf = append(c.buf, e)
	if c.timer != nil {
		c.timer.Reset(c.window)
	} else {
		c.timer = time.AfterFunc(c.window, c.flush)
	}
}

// Flush immediately drains the current buffer as a Group.
func (c *Coalescer) Flush() {
	c.flush()
}

func (c *Coalescer) flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.buf) == 0 {
		return
	}
	g := make(Group, len(c.buf))
	copy(g, c.buf)
	c.buf = c.buf[:0]
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.out <- g
}

// Close flushes any remaining entries and closes the output channel.
func (c *Coalescer) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	if len(c.buf) > 0 {
		g := make(Group, len(c.buf))
		copy(g, c.buf)
		c.buf = nil
		c.out <- g
	}
	close(c.out)
}
