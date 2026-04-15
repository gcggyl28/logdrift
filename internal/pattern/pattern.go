// Package pattern provides glob and regex pattern matching utilities
// used to route, filter, or annotate log entries by their content.
package pattern

import (
	"fmt"
	"regexp"
)

// MatchMode controls how multiple patterns are evaluated.
type MatchMode string

const (
	MatchAny MatchMode = "any" // entry matches if at least one pattern matches
	MatchAll MatchMode = "all" // entry matches only if every pattern matches
)

// Matcher holds compiled patterns and applies them to strings.
type Matcher struct {
	mode     MatchMode
	patterns []*regexp.Regexp
}

// New compiles each pattern string and returns a Matcher.
// mode must be "any" or "all". An empty patterns slice is valid and
// matches every input regardless of mode.
func New(mode MatchMode, patterns []string) (*Matcher, error) {
	switch mode {
	case MatchAny, MatchAll:
	default:
		return nil, fmt.Errorf("pattern: unknown mode %q, want \"any\" or \"all\"", mode)
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if p == "" {
			return nil, fmt.Errorf("pattern: empty pattern string")
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("pattern: compile %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	return &Matcher{mode: mode, patterns: compiled}, nil
}

// Match reports whether s satisfies the configured mode and patterns.
// If no patterns are configured, Match always returns true.
func (m *Matcher) Match(s string) bool {
	if len(m.patterns) == 0 {
		return true
	}

	switch m.mode {
	case MatchAll:
		for _, re := range m.patterns {
			if !re.MatchString(s) {
				return false
			}
		}
		return true
	default: // MatchAny
		for _, re := range m.patterns {
			if re.MatchString(s) {
				return true
			}
		}
		return false
	}
}

// Mode returns the configured MatchMode.
func (m *Matcher) Mode() MatchMode { return m.mode }

// Len returns the number of compiled patterns.
func (m *Matcher) Len() int { return len(m.patterns) }
