// Package remap rewrites field values using a static lookup table.
package remap

import (
	"errors"
	"fmt"
)

// Rule maps values of a single field through a lookup table.
type Rule struct {
	Field   string
	Mapping map[string]string
}

// Remapper applies value remapping rules to log entry fields.
type Remapper struct {
	rules []Rule
}

// New creates a Remapper from the given rules.
// Each rule must have a non-empty Field and at least one mapping entry.
func New(rules []Rule) (*Remapper, error) {
	if len(rules) == 0 {
		return nil, errors.New("remap: at least one rule is required")
	}
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("remap: rule %d has empty field name", i)
		}
		if len(r.Mapping) == 0 {
			return nil, fmt.Errorf("remap: rule %d (%q) has empty mapping", i, r.Field)
		}
	}
	return &Remapper{rules: rules}, nil
}

// Apply rewrites fields in entry according to the configured rules.
// Fields not present in an entry or values not found in a mapping are left unchanged.
func (r *Remapper) Apply(entry map[string]string) map[string]string {
	out := make(map[string]string, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, rule := range r.rules {
		if cur, ok := out[rule.Field]; ok {
			if mapped, found := rule.Mapping[cur]; found {
				out[rule.Field] = mapped
			}
		}
	}
	return out
}
