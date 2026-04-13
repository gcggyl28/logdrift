// Package throttle provides a per-service line-rate throttler that
// suppresses log entries when a service emits lines faster than a
// configured burst window, preventing diff noise from log storms.
package throttle

import (
	"errors"
	"sync"
	"time"
)

// Throttle tracks per-service emission rates and decides whether a
// given line should be forwarded downstream.
type Throttle struct {
	maxPerWindow int
	window       time.Duration
	mu           sync.Mutex
	buckets      map[string]*bucket
}

type bucket struct {
	count     int
	windowEnd time.Time
}

// New creates a Throttle that allows at most maxPerWindow lines per
// service within the given window duration.
// A maxPerWindow of 0 disables throttling (all lines pass).
func New(maxPerWindow int, window time.Duration) (*Throttle, error) {
	if maxPerWindow < 0 {
		return nil, errors.New("throttle: maxPerWindow must be >= 0")
	}
	if window <= 0 && maxPerWindow > 0 {
		return nil, errors.New("throttle: window must be positive when maxPerWindow > 0")
	}
	return &Throttle{
		maxPerWindow: maxPerWindow,
		window:       window,
		buckets:      make(map[string]*bucket),
	}, nil
}

// Allow reports whether the next line from service should be forwarded.
// It is safe for concurrent use.
func (t *Throttle) Allow(service string) bool {
	if t.maxPerWindow == 0 {
		return true
	}
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	b, ok := t.buckets[service]
	if !ok || now.After(b.windowEnd) {
		t.buckets[service] = &bucket{count: 1, windowEnd: now.Add(t.window)}
		return true
	}
	if b.count >= t.maxPerWindow {
		return false
	}
	b.count++
	return true
}

// Reset clears the rate counters for all services.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buckets = make(map[string]*bucket)
}
