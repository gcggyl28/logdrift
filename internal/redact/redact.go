// Package redact provides log line redaction to mask sensitive values
// such as tokens, passwords, and PII before they are displayed or exported.
package redact

import (
	"fmt"
	"regexp"
)

// Redactor masks substrings in log lines that match configured patterns.
type Redactor struct {
	patterns []*regexp.Regexp
	mask     string
}

// Option configures a Redactor.
type Option func(*Redactor)

// WithMask overrides the default replacement string.
func WithMask(mask string) Option {
	return func(r *Redactor) {
		r.mask = mask
	}
}

// New creates a Redactor from a slice of regex pattern strings.
// Returns an error if any pattern fails to compile.
func New(patterns []string, opts ...Option) (*Redactor, error) {
	r := &Redactor{mask: "[REDACTED]"}
	for _, o := range opts {
		o(r)
	}
	for _, p := range patterns {
		if p == "" {
			return nil, fmt.Errorf("redact: empty pattern is not allowed")
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("redact: invalid pattern %q: %w", p, err)
		}
		r.patterns = append(r.patterns, re)
	}
	return r, nil
}

// Apply returns a copy of line with all matching substrings replaced by the mask.
func (r *Redactor) Apply(line string) string {
	for _, re := range r.patterns {
		line = re.ReplaceAllString(line, r.mask)
	}
	return line
}

// Enabled reports whether the redactor has any active patterns.
func (r *Redactor) Enabled() bool {
	return len(r.patterns) > 0
}
