// Package evict provides a time-based eviction policy that removes entries
// from a tracked set once their TTL has elapsed. It is safe for concurrent use.
package evict

import (
	"errors"
	"sync"
	"time"
)

// Entry holds the value and its expiry deadline.
type Entry struct {
	Value     string
	ExpiresAt time.Time
}

// Evictor tracks string keys with individual TTLs and evicts expired ones.
type Evictor struct {
	mu      sync.Mutex
	entries map[string]Entry
	ttl     time.Duration
}

// New creates an Evictor with the given TTL. TTL must be positive.
func New(ttl time.Duration) (*Evictor, error) {
	if ttl <= 0 {
		return nil, errors.New("evict: TTL must be positive")
	}
	return &Evictor{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}, nil
}

// Add inserts or refreshes the key with a new expiry deadline.
func (e *Evictor) Add(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.entries[key] = Entry{
		Value:     key,
		ExpiresAt: time.Now().Add(e.ttl),
	}
}

// Has reports whether the key exists and has not yet expired.
// Expired entries are lazily removed on lookup.
func (e *Evictor) Has(key string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	ent, ok := e.entries[key]
	if !ok {
		return false
	}
	if time.Now().After(ent.ExpiresAt) {
		delete(e.entries, key)
		return false
	}
	return true
}

// Evict removes all entries whose TTL has elapsed. It returns the number
// of entries removed.
func (e *Evictor) Evict() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	removed := 0
	for k, ent := range e.entries {
		if now.After(ent.ExpiresAt) {
			delete(e.entries, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of live (non-expired) entries.
func (e *Evictor) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	count := 0
	for _, ent := range e.entries {
		if !now.After(ent.ExpiresAt) {
			count++
		}
	}
	return count
}
