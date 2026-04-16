// Package timefield extracts and normalises a timestamp from a parsed log entry.
package timefield

import (
	"fmt"
	"time"

	"github.com/angryboat/logdrift/internal/parser"
)

// Extractor reads a named field from a log entry and parses it as a time value.
type Extractor struct {
	field   string
	layouts []string
	fallback bool
}

// Option configures an Extractor.
type Option func(*Extractor)

// WithLayouts sets the list of time layouts tried in order.
func WithLayouts(layouts ...string) Option {
	return func(e *Extractor) { e.layouts = layouts }
}

// WithFallback causes Extract to return time.Now() when the field is absent or
// unparseable instead of returning an error.
func WithFallback() Option {
	return func(e *Extractor) { e.fallback = true }
}

var defaultLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	time.DateTime,
}

// New creates an Extractor that reads timestamps from field.
func New(field string, opts ...Option) (*Extractor, error) {
	if field == "" {
		return nil, fmt.Errorf("timefield: field name must not be empty")
	}
	e := &Extractor{field: field, layouts: defaultLayouts}
	for _, o := range opts {
		o(e)
	}
	if len(e.layouts) == 0 {
		return nil, fmt.Errorf("timefield: at least one layout is required")
	}
	return e, nil
}

// Extract returns the parsed timestamp from entry, or an error.
func (e *Extractor) Extract(entry parser.Entry) (time.Time, error) {
	raw, ok := entry.Fields[e.field]
	if !ok {
		if e.fallback {
			return time.Now(), nil
		}
		return time.Time{}, fmt.Errorf("timefield: field %q not present in entry", e.field)
	}
	for _, layout := range e.layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	if e.fallback {
		return time.Now(), nil
	}
	return time.Time{}, fmt.Errorf("timefield: could not parse %q with any known layout", raw)
}
