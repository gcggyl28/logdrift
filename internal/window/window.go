// Package window provides a sliding time-window counter for tracking
// event rates over a configurable rolling duration.
package window

import (
	"errors"
	"sync"
	"time"
)

// Entry holds a single timestamped event count.
type Entry struct {
	At    time.Time
	Count int64
}

// Window is a thread-safe sliding time-window counter.
type Window struct {
	mu       sync.Mutex
	duration time.Duration
	buckets  []Entry
}

// New creates a Window with the given duration.
// Returns an error if duration is non-positive.
func New(d time.Duration) (*Window, error) {
	if d <= 0 {
		return nil, errors.New("window: duration must be positive")
	}
	return &Window{duration: d}, nil
}

// Add records n events at the current time.
func (w *Window) Add(n int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(time.Now())
	w.buckets = append(w.buckets, Entry{At: time.Now(), Count: n})
}

// Total returns the sum of all event counts within the current window.
func (w *Window) Total() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(time.Now())
	var sum int64
	for _, b := range w.buckets {
		sum += b.Count
	}
	return sum
}

// Reset discards all buckets.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buckets = w.buckets[:0]
}

// Buckets returns a snapshot of the current in-window entries.
func (w *Window) Buckets() []Entry {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(time.Now())
	out := make([]Entry, len(w.buckets))
	copy(out, w.buckets)
	return out
}

// evict removes entries older than the window duration.
// Caller must hold w.mu.
func (w *Window) evict(now time.Time) {
	cutoff := now.Add(-w.duration)
	i := 0
	for i < len(w.buckets) && w.buckets[i].At.Before(cutoff) {
		i++
	}
	w.buckets = w.buckets[i:]
}
