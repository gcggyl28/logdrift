// Package normalize provides log line normalization utilities that strip
// or replace volatile fields (timestamps, PIDs, UUIDs) before diffing so
// that structurally identical lines are not reported as drift.
package normalize

import (
	"fmt"
	"regexp"
	"strings"
)

// Mode controls which built-in normalizations are applied.
type Mode string

const (
	ModeNone      Mode = "none"
	ModeTimestamp Mode = "timestamp"
	ModePID       Mode = "pid"
	ModeFull      Mode = "full" // timestamp + pid + uuid
)

var builtinPatterns = map[Mode][]*regexp.Regexp{
	ModeTimestamp: {
		regexp.MustCompile(`\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})?`),
	},
	ModePID: {
		regexp.MustCompile(`\bpid[=:]\s*\d+\b`),
	},
	ModeFull: {
		regexp.MustCompile(`\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})?`),
		regexp.MustCompile(`\bpid[=:]\s*\d+\b`),
		regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`),
	},
}

// Normalizer applies a set of regexp substitutions to log lines.
type Normalizer struct {
	patterns []*regexp.Regexp
	placeholder string
}

// New returns a Normalizer for the given mode plus any extra raw patterns.
// placeholder is the string used to replace matched substrings; if empty
// it defaults to "<?>".
func New(mode Mode, extraPatterns []string, placeholder string) (*Normalizer, error) {
	if placeholder == "" {
		placeholder = "<?>"
	}

	patterns, ok := builtinPatterns[mode]
	if !ok && mode != ModeNone {
		return nil, fmt.Errorf("normalize: unknown mode %q", mode)
	}

	compiled := make([]*regexp.Regexp, len(patterns))
	copy(compiled, patterns)

	for _, raw := range extraPatterns {
		if strings.TrimSpace(raw) == "" {
			return nil, fmt.Errorf("normalize: empty pattern")
		}
		re, err := regexp.Compile(raw)
		if err != nil {
			return nil, fmt.Errorf("normalize: invalid pattern %q: %w", raw, err)
		}
		compiled = append(compiled, re)
	}

	return &Normalizer{patterns: compiled, placeholder: placeholder}, nil
}

// Apply returns the normalized form of line.
func (n *Normalizer) Apply(line string) string {
	for _, re := range n.patterns {
		line = re.ReplaceAllString(line, n.placeholder)
	}
	return line
}
