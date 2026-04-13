// Package filter provides log line filtering based on patterns.
package filter

import (
	"fmt"
	"regexp"
)

// Mode defines how multiple patterns are combined.
type Mode string

const (
	ModeAny Mode = "any" // line matches if any pattern matches
	ModeAll Mode = "all" // line matches if all patterns match
)

// Filter holds compiled patterns and match mode.
type Filter struct {
	patterns []*regexp.Regexp
	mode     Mode
	invert   bool
}

// New creates a Filter from raw pattern strings.
// mode must be "any" or "all". invert negates the final result.
func New(patterns []string, mode Mode, invert bool) (*Filter, error) {
	if mode != ModeAny && mode != ModeAll {
		return nil, fmt.Errorf("filter: unknown mode %q, want \"any\" or \"all\"", mode)
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Filter{patterns: compiled, mode: mode, invert: invert}, nil
}

// Match reports whether line passes the filter.
func (f *Filter) Match(line string) bool {
	if len(f.patterns) == 0 {
		return true
	}
	var result bool
	switch f.mode {
	case ModeAll:
		result = f.matchAll(line)
	default:
		result = f.matchAny(line)
	}
	if f.invert {
		return !result
	}
	return result
}

func (f *Filter) matchAny(line string) bool {
	for _, re := range f.patterns {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

func (f *Filter) matchAll(line string) bool {
	for _, re := range f.patterns {
		if !re.MatchString(line) {
			return false
		}
	}
	return true
}
