// Package tagfilter provides filtering of log entries by their service label tags.
package tagfilter

import (
	"errors"
	"fmt"
)

// Entry is the minimal interface tagfilter operates on.
type Entry struct {
	Service string
	Tags    map[string]string
	Line    string
}

// Mode controls how multiple required tags are evaluated.
type Mode string

const (
	ModeAny Mode = "any"
	ModeAll Mode = "all"
)

// Filter allows or denies entries based on tag key/value pairs.
type Filter struct {
	mode Mode
	tags map[string]string
}

// New creates a Filter that matches entries whose tags satisfy the given
// key/value pairs under the specified mode ("any" or "all").
func New(mode Mode, tags map[string]string) (*Filter, error) {
	if mode != ModeAny && mode != ModeAll {
		return nil, fmt.Errorf("tagfilter: unknown mode %q", mode)
	}
	if len(tags) == 0 {
		return nil, errors.New("tagfilter: at least one tag required")
	}
	for k, v := range tags {
		if k == "" {
			return nil, errors.New("tagfilter: tag key must not be empty")
		}
		if v == "" {
			return nil, fmt.Errorf("tagfilter: value for key %q must not be empty", k)
		}
	}
	copy := make(map[string]string, len(tags))
	for k, v := range tags {
		copy[k] = v
	}
	return &Filter{mode: mode, tags: copy}, nil
}

// Allow returns true when the entry's tags satisfy the filter criteria.
func (f *Filter) Allow(e Entry) bool {
	for k, want := range f.tags {
		got, ok := e.Tags[k]
		matched := ok && got == want
		if f.mode == ModeAny && matched {
			return true
		}
		if f.mode == ModeAll && !matched {
			return false
		}
	}
	return f.mode == ModeAll
}
