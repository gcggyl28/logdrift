// Package truncate provides line-length truncation for log entries.
// Lines exceeding the configured maximum byte length are clipped and
// optionally suffixed with an ellipsis marker so consumers can detect
// that content was removed.
package truncate

import (
	"errors"
	"fmt"
)

const defaultSuffix = "..."

// Truncator clips log lines that exceed a maximum byte length.
type Truncator struct {
	maxLen int
	suffix string
	enabled bool
}

// Option configures a Truncator.
type Option func(*Truncator)

// WithSuffix overrides the default ellipsis suffix appended to truncated lines.
func WithSuffix(s string) Option {
	return func(t *Truncator) { t.suffix = s }
}

// New creates a Truncator that clips lines longer than maxLen bytes.
// A maxLen of 0 disables truncation. A negative maxLen returns an error.
func New(maxLen int, opts ...Option) (*Truncator, error) {
	if maxLen < 0 {
		return nil, errors.New("truncate: maxLen must be >= 0")
	}
	t := &Truncator{
		maxLen:  maxLen,
		suffix:  defaultSuffix,
		enabled: maxLen > 0,
	}
	for _, o := range opts {
		o(t)
	}
	if t.enabled && len(t.suffix) >= maxLen {
		return nil, fmt.Errorf("truncate: suffix length (%d) must be less than maxLen (%d)", len(t.suffix), maxLen)
	}
	return t, nil
}

// Apply returns the line unchanged when truncation is disabled or the line fits
// within maxLen. Otherwise it clips the line and appends the suffix.
func (t *Truncator) Apply(line string) string {
	if !t.enabled || len(line) <= t.maxLen {
		return line
	}
	keep := t.maxLen - len(t.suffix)
	return line[:keep] + t.suffix
}

// Enabled reports whether the Truncator will actually clip lines.
func (t *Truncator) Enabled() bool { return t.enabled }
