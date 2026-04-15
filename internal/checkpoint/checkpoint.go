// Package checkpoint persists and restores log read positions
// so that logdrift can resume tailing after a restart without
// re-processing already-seen lines.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// ErrNoCheckpoint is returned when no checkpoint exists for a service.
var ErrNoCheckpoint = errors.New("checkpoint: no entry found")

// Entry holds the persisted read position for a single service.
type Entry struct {
	Service   string    `json:"service"`
	Offset    int64     `json:"offset"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store persists checkpoint entries to a JSON file on disk.
type Store struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Entry
}

// New opens (or creates) a checkpoint store at the given file path.
func New(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("checkpoint: path must not be empty")
	}
	s := &Store{path: path, entries: make(map[string]Entry)}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Get returns the persisted Entry for a service, or ErrNoCheckpoint.
func (s *Store) Get(service string) (Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[service]
	if !ok {
		return Entry{}, ErrNoCheckpoint
	}
	return e, nil
}

// Set updates the offset for a service and flushes to disk.
func (s *Store) Set(service string, offset int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[service] = Entry{
		Service:   service,
		Offset:    offset,
		UpdatedAt: time.Now().UTC(),
	}
	return s.flush()
}

// Delete removes a service entry and flushes to disk.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, service)
	return s.flush()
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	for _, e := range entries {
		s.entries[e.Service] = e
	}
	return nil
}

func (s *Store) flush() error {
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
