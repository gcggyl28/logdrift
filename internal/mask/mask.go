// Package mask provides field-level masking for log entries,
// replacing sensitive field values with a configurable placeholder.
package mask

import (
	"errors"
	"fmt"
)

const defaultPlaceholder = "***"

// Masker replaces the values of named fields with a placeholder string.
type Masker struct {
	fields      map[string]struct{}
	placeholder string
}

// Option configures a Masker.
type Option func(*Masker)

// WithPlaceholder overrides the default placeholder string.
func WithPlaceholder(p string) Option {
	return func(m *Masker) { m.placeholder = p }
}

// New creates a Masker that masks the given field names.
// At least one field name must be supplied.
func New(fields []string, opts ...Option) (*Masker, error) {
	if len(fields) == 0 {
		return nil, errors.New("mask: at least one field name required")
	}
	m := &Masker{
		fields:      make(map[string]struct{}, len(fields)),
		placeholder: defaultPlaceholder,
	}
	for _, opt := range opts {
		opt(m)
	}
	if m.placeholder == "" {
		return nil, errors.New("mask: placeholder must not be empty")
	}
	for _, f := range fields {
		if f == "" {
			return nil, fmt.Errorf("mask: field name must not be empty")
		}
		m.fields[f] = struct{}{}
	}
	return m, nil
}

// Apply returns a copy of the entry map with sensitive fields replaced.
// Fields not present in the entry are ignored.
func (m *Masker) Apply(entry map[string]string) map[string]string {
	out := make(map[string]string, len(entry))
	for k, v := range entry {
		if _, ok := m.fields[k]; ok {
			out[k] = m.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// Fields returns the set of masked field names.
func (m *Masker) Fields() []string {
	out := make([]string, 0, len(m.fields))
	for f := range m.fields {
		out = append(out, f)
	}
	return out
}
