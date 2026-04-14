// Package aggregate groups log entries by a key (e.g. service, label) and
// computes basic statistics over a sliding time window.
package aggregate

import (
	"errors"
	"sync"
	"time"
)

// Stats holds aggregated counters for a single key.
type Stats struct {
	Key      string
	Lines    int64
	Drifts   int64
	FirstSeen time.Time
	LastSeen  time.Time
}

// Aggregator accumulates per-key statistics.
type Aggregator struct {
	mu     sync.Mutex
	window time.Duration
	data   map[string]*Stats
}

// New creates an Aggregator with the given sliding window duration.
// A zero window disables eviction.
func New(window time.Duration) (*Aggregator, error) {
	if window < 0 {
		return nil, errors.New("aggregate: window must be non-negative")
	}
	return &Aggregator{
		window: window,
		data:   make(map[string]*Stats),
	}, nil
}

// RecordLine increments the line counter for key at time t.
func (a *Aggregator) RecordLine(key string, t time.Time) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.evict(t)
	s := a.getOrCreate(key, t)
	s.Lines++
	s.LastSeen = t
}

// RecordDrift increments the drift counter for key at time t.
func (a *Aggregator) RecordDrift(key string, t time.Time) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.evict(t)
	s := a.getOrCreate(key, t)
	s.Drifts++
	s.LastSeen = t
}

// Snapshot returns a copy of all current stats.
func (a *Aggregator) Snapshot() []Stats {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]Stats, 0, len(a.data))
	for _, s := range a.data {
		copy := *s
		out = append(out, copy)
	}
	return out
}

// Reset clears all accumulated data.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.data = make(map[string]*Stats)
}

func (a *Aggregator) getOrCreate(key string, t time.Time) *Stats {
	if s, ok := a.data[key]; ok {
		return s
	}
	s := &Stats{Key: key, FirstSeen: t, LastSeen: t}
	a.data[key] = s
	return s
}

func (a *Aggregator) evict(now time.Time) {
	if a.window == 0 {
		return
	}
	cutoff := now.Add(-a.window)
	for k, s := range a.data {
		if s.LastSeen.Before(cutoff) {
			delete(a.data, k)
		}
	}
}
