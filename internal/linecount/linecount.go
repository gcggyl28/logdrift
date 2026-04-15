// Package linecount tracks per-service line rates over a sliding time window.
package linecount

import (
	"errors"
	"sync"
	"time"
)

// Counter records line events and exposes a rate (lines/sec) per service.
type Counter struct {
	mu       sync.Mutex
	window   time.Duration
	buckets  map[string][]time.Time
}

// New creates a Counter with the given sliding window duration.
// Returns an error if window is non-positive.
func New(window time.Duration) (*Counter, error) {
	if window <= 0 {
		return nil, errors.New("linecount: window must be positive")
	}
	return &Counter{
		window:  window,
		buckets: make(map[string][]time.Time),
	}, nil
}

// Record registers a line event for the named service at the current time.
func (c *Counter) Record(service string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.buckets[service] = append(c.buckets[service], now)
	c.evict(service, now)
}

// Rate returns the number of lines recorded for service within the window,
// divided by the window duration in seconds (lines/sec).
func (c *Counter) Rate(service string) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.evict(service, now)
	count := len(c.buckets[service])
	if count == 0 {
		return 0
	}
	return float64(count) / c.window.Seconds()
}

// Services returns the names of all services that have recorded at least one
// line within the current window.
func (c *Counter) Services() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	out := make([]string, 0, len(c.buckets))
	for svc := range c.buckets {
		c.evict(svc, now)
		if len(c.buckets[svc]) > 0 {
			out = append(out, svc)
		}
	}
	return out
}

// evict removes timestamps older than the window. Must be called with mu held.
func (c *Counter) evict(service string, now time.Time) {
	cutoff := now.Add(-c.window)
	ts := c.buckets[service]
	i := 0
	for i < len(ts) && ts[i].Before(cutoff) {
		i++
	}
	c.buckets[service] = ts[i:]
}
