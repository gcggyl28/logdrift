// Package highlight provides keyword-based line highlighting for log output.
// It wraps matched substrings with ANSI escape codes or custom markers
// depending on the configured style.
package highlight

import (
	"fmt"
	"regexp"
	"strings"
)

// Style controls how matched text is rendered.
type Style string

const (
	StyleANSI  Style = "ansi"  // wrap with ANSI colour codes
	StyleBraces Style = "braces" // wrap with [[ … ]]
)

// Highlighter applies keyword highlighting to log lines.
type Highlighter struct {
	patterns []*regexp.Regexp
	style    Style
	color    string // ANSI colour code, e.g. "\033[33m"
}

// New creates a Highlighter for the given keyword patterns and style.
// patterns must be valid regular expressions; style must be "ansi" or "braces".
func New(patterns []string, style Style, color string) (*Highlighter, error) {
	switch style {
	case StyleANSI, StyleBraces:
	default:
		return nil, fmt.Errorf("highlight: unknown style %q", style)
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if strings.TrimSpace(p) == "" {
			return nil, fmt.Errorf("highlight: pattern must not be empty")
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("highlight: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	return &Highlighter{patterns: compiled, style: style, color: color}, nil
}

// Apply returns the line with all pattern matches highlighted.
// If no patterns are configured the original line is returned unchanged.
func (h *Highlighter) Apply(line string) string {
	if len(h.patterns) == 0 {
		return line
	}
	for _, re := range h.patterns {
		line = re.ReplaceAllStringFunc(line, func(match string) string {
			return h.wrap(match)
		})
	}
	return line
}

func (h *Highlighter) wrap(s string) string {
	switch h.style {
	case StyleBraces:
		return "[[" + s + "]]"
	default: // StyleANSI
		code := h.color
		if code == "" {
			code = "\033[33m" // default yellow
		}
		return code + s + "\033[0m"
	}
}
