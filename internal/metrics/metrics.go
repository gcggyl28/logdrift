// Package metrics tracks per-service log line counters and drift event counts
// for the lifetime of a logdrift session.
package metrics

import (
	"sync"
	"time"
)

// Snapshot is an immutable view of the current metrics state.
type Snapshot struct {
	LinesReceived map[string]int64
	DriftEvents   int64
	StartedAt     time.Time
	CapturedAt    time.Time
}

// Tracker accumulates metrics in a thread-safe manner.
type Tracker struct {
	mu            sync.RWMutex
	linesReceived map[string]int64
	driftEvents   int64
	startedAt     time.Time
}

// New creates a new Tracker with the clock set to now.
func New() *Tracker {
	return &Tracker{
		linesReceived: make(map[string]int64),
		startedAt:     time.Now(),
	}
}

// RecordLine increments the line counter for the given service.
func (t *Tracker) RecordLine(service string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.linesReceived[service]++
}

// RecordDrift increments the global drift event counter.
func (t *Tracker) RecordDrift() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.driftEvents++
}

// Snapshot returns an immutable copy of the current state.
func (t *Tracker) Snapshot() Snapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()

	lines := make(map[string]int64, len(t.linesReceived))
	for k, v := range t.linesReceived {
		lines[k] = v
	}
	return Snapshot{
		LinesReceived: lines,
		DriftEvents:   t.driftEvents,
		StartedAt:     t.startedAt,
		CapturedAt:    time.Now(),
	}
}

// Reset zeroes all counters without changing the start time.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.linesReceived = make(map[string]int64)
	t.driftEvents = 0
}
