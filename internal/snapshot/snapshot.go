// Package snapshot maintains a rolling window of recent log lines
// per service, used to produce meaningful diffs across streams.
package snapshot

import "sync"

// Snapshot holds the most recent N lines for a named service.
type Snapshot struct {
	mu      sync.RWMutex
	window  int
	buckets map[string][]string
}

// New creates a Snapshot that retains up to window lines per service.
// If window is <= 0 it defaults to 100.
func New(window int) *Snapshot {
	if window <= 0 {
		window = 100
	}
	return &Snapshot{
		window:  window,
		buckets: make(map[string][]string),
	}
}

// Push appends line to the buffer for service, evicting the oldest
// entry when the window is full.
func (s *Snapshot) Push(service, line string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	buf := s.buckets[service]
	if len(buf) >= s.window {
		buf = buf[1:]
	}
	s.buckets[service] = append(buf, line)
}

// Lines returns a copy of the current buffer for service.
func (s *Snapshot) Lines(service string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	buf := s.buckets[service]
	out := make([]string, len(buf))
	copy(out, buf)
	return out
}

// Services returns the names of all services that have at least one line.
func (s *Snapshot) Services() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.buckets))
	for k := range s.buckets {
		names = append(names, k)
	}
	return names
}

// Reset clears all buffered lines for service.
func (s *Snapshot) Reset(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.buckets, service)
}
