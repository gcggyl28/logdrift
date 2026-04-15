// Package enrich attaches structured metadata fields to log entries
// before they are passed downstream for diffing or rendering.
package enrich

import (
	"errors"
	"fmt"
	"time"
)

// Entry is a log line augmented with metadata.
type Entry struct {
	Service   string
	Line      string
	Timestamp time.Time
	Fields    map[string]string
}

// Enricher attaches a fixed set of key/value fields to every entry.
type Enricher struct {
	fields map[string]string
}

// New creates an Enricher that will stamp each entry with the provided
// static fields. At least one field is required.
func New(fields map[string]string) (*Enricher, error) {
	if len(fields) == 0 {
		return nil, errors.New("enrich: at least one field is required")
	}
	for k, v := range fields {
		if k == "" {
			return nil, errors.New("enrich: field key must not be empty")
		}
		if v == "" {
			return nil, fmt.Errorf("enrich: value for key %q must not be empty", k)
		}
	}
	copy := make(map[string]string, len(fields))
	for k, v := range fields {
		copy[k] = v
	}
	return &Enricher{fields: copy}, nil
}

// Apply returns a new Entry with the enricher's static fields merged into
// the entry's existing Fields map. Existing keys are not overwritten.
func (e *Enricher) Apply(entry Entry) Entry {
	out := entry
	merged := make(map[string]string, len(entry.Fields)+len(e.fields))
	for k, v := range e.fields {
		merged[k] = v
	}
	for k, v := range entry.Fields {
		merged[k] = v // entry fields win on conflict
	}
	out.Fields = merged
	return out
}

// ApplyAll enriches a slice of entries in place, returning the results.
func (e *Enricher) ApplyAll(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, en := range entries {
		out[i] = e.Apply(en)
	}
	return out
}
