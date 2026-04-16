// Package offset tracks the read position (byte offset) for each tailed
// service source so that logdrift can resume from where it left off after a
// restart without relying on the full checkpoint subsystem.
package offset

import (
	"errors"
	"fmt"
	"sync"
)

// ErrUnknownService is returned when a service has no offset.
var ErrUnknownService = errors.New("offset: unknown service")

// Tracker holds per-service byte offsets in memory.
type Tracker struct {
	mu      sync.RWMutex
	offsets map[string]int64
}

// New creates an empty Tracker.
func New() *Tracker {
	return &Tracker{offsets: make(map[string]int64)}
}

// Set records the current offset for a service. service must be non-empty.
func (t *Tracker) Set(service string, offset int64) error {
	if service == "" {
		return fmt.Errorf("offset: service name must not be empty")
	}
	if offset < 0 {
		return fmt.Errorf("offset: offset must be >= 0, got %d", offset)
	}
	t.mu.Lock()
	t.offsets[service] = offset
	t.mu.Unlock()
	return nil
}

// Get returns the recorded offset for a service.
// Returns ErrUnknownService if the service has never been set.
func (t *Tracker) Get(service string) (int64, error) {
	t.mu.RLock()
	v, ok := t.offsets[service]
	t.mu.RUnlock()
	if !ok {
		return 0, ErrUnknownService
	}
	return v, nil
}

// Advance increments the stored offset for a service by delta bytes.
// The service must already exist; call Set first to initialise.
func (t *Tracker) Advance(service string, delta int64) error {
	if delta < 0 {
		return fmt.Errorf("offset: delta must be >= 0, got %d", delta)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	v, ok := t.offsets[service]
	if !ok {
		return ErrUnknownService
	}
	t.offsets[service] = v + delta
	return nil
}

// Services returns a snapshot of all tracked service names.
func (t *Tracker) Services() []string {
	t.mu.RLock()
	out := make([]string, 0, len(t.offsets))
	for s := range t.offsets {
		out = append(out, s)
	}
	t.mu.RUnlock()
	return out
}
