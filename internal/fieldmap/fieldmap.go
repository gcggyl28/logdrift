// Package fieldmap provides field renaming and remapping for parsed log entries.
// It allows operators to normalise field names across heterogeneous log sources
// so that downstream diff and alert logic operates on a consistent schema.
package fieldmap

import (
	"errors"
	"fmt"
)

// Rule describes a single field rename: From is the source key, To is the
// destination key. If To is empty the field is dropped from the entry.
type Rule struct {
	From string
	To   string
}

// Mapper applies a set of renaming rules to a map of log fields.
type Mapper struct {
	rules []Rule
	// index for O(1) lookup
	idx map[string]string
}

// New validates rules and returns a ready-to-use Mapper.
// It returns an error if any rule has an empty From field or if two rules
// share the same From key.
func New(rules []Rule) (*Mapper, error) {
	if len(rules) == 0 {
		return nil, errors.New("fieldmap: at least one rule is required")
	}
	idx := make(map[string]string, len(rules))
	for i, r := range rules {
		if r.From == "" {
			return nil, fmt.Errorf("fieldmap: rule[%d]: From must not be empty", i)
		}
		if _, exists := idx[r.From]; exists {
			return nil, fmt.Errorf("fieldmap: duplicate From key %q", r.From)
		}
		idx[r.From] = r.To
	}
	return &Mapper{rules: rules, idx: idx}, nil
}

// Apply returns a new map with the renaming rules applied.
// Fields whose From key is not present in fields are silently ignored.
// Fields with an empty To value are dropped.
// All other fields are passed through unchanged.
func (m *Mapper) Apply(fields map[string]string) map[string]string {
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		if to, matched := m.idx[k]; matched {
			if to != "" {
				out[to] = v
			}
			// empty To → drop the field
		} else {
			out[k] = v
		}
	}
	return out
}

// Rules returns a copy of the rules used to build this Mapper.
func (m *Mapper) Rules() []Rule {
	cp := make([]Rule, len(m.rules))
	copy(cp, m.rules)
	return cp
}
