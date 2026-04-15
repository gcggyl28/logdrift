// Package transform applies a sequence of text transformations to log line
// content before the line is forwarded to downstream consumers.
package transform

import (
	"errors"
	"strings"
)

// Op is a single transformation operation.
type Op int

const (
	OpUppercase Op = iota
	OpLowercase
	OpTrimSpace
	OpTrimPrefix
	OpTrimSuffix
)

// Rule describes one transformation step.
type Rule struct {
	Op     Op
	Arg    string // used by TrimPrefix / TrimSuffix
}

// Transformer applies an ordered list of Rules to a string.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer for the given rules.
// It returns an error when the rules slice is empty or a TrimPrefix/TrimSuffix
// rule is missing its argument.
func New(rules []Rule) (*Transformer, error) {
	if len(rules) == 0 {
		return nil, errors.New("transform: at least one rule is required")
	}
	for i, r := range rules {
		if (r.Op == OpTrimPrefix || r.Op == OpTrimSuffix) && r.Arg == "" {
			return nil, errors.New("transform: rule %d requires a non-empty)
		}
		_ = i
	}
	return &Transformer{rules: rules}, nil
}

// Apply runs every rule in order and returns the transformed string.
func (t *Transformer) Apply(s string) string {
	for _, r := range t.rules {
		switch r.Op {
		case OpUppercase:
			s = strings.ToUpper(s)
		case OpLowercase:
			s = strings.ToLower(s)
		case OpTrimSpace:
			s = strings.TrimSpace(s)
		case OpTrimPrefix:
			s = strings.TrimPrefix(s, r.Arg)
		case OpTrimSuffix:
			s = strings.TrimSuffix(s, r.Arg)
		}
	}
	return s
}
