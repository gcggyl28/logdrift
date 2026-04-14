// Package label provides service label management for log entries,
// allowing consistent tagging and lookup of log sources by name or alias.
package label

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// ErrDuplicateLabel is returned when a label or alias is already registered.
var ErrDuplicateLabel = errors.New("label: duplicate label or alias")

// ErrNotFound is returned when a label cannot be resolved.
var ErrNotFound = errors.New("label: not found")

// Registry maps service names and their optional aliases to a canonical label.
type Registry struct {
	mu      sync.RWMutex
	labels  map[string]string // canonical -> canonical
	aliases map[string]string // alias -> canonical
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		labels:  make(map[string]string),
		aliases: make(map[string]string),
	}
}

// Register adds a canonical label with optional aliases.
// Returns ErrDuplicateLabel if the label or any alias is already registered.
func (r *Registry) Register(name string, aliases ...string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("label: name must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.labels[name]; ok {
		return ErrDuplicateLabel
	}
	for _, a := range aliases {
		a = strings.TrimSpace(a)
		if _, ok := r.aliases[a]; ok {
			return ErrDuplicateLabel
		}
	}
	r.labels[name] = name
	for _, a := range aliases {
		r.aliases[strings.TrimSpace(a)] = name
	}
	return nil
}

// Resolve returns the canonical label for a given name or alias.
func (r *Registry) Resolve(nameOrAlias string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.labels[nameOrAlias]; ok {
		return nameOrAlias, nil
	}
	if canon, ok := r.aliases[nameOrAlias]; ok {
		return canon, nil
	}
	return "", ErrNotFound
}

// List returns all registered canonical labels in insertion order (sorted).
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.labels))
	for k := range r.labels {
		out = append(out, k)
	}
	return out
}
