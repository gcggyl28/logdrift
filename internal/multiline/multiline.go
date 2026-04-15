// Package multiline provides a combiner that merges multi-line log events
// (e.g. Java stack traces, Python tracebacks) into a single logical entry
// before the line is forwarded downstream.
package multiline

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// Mode controls how continuation lines are detected.
type Mode string

const (
	// ModePrefix treats any line matching the pattern as the START of a new event.
	ModePrefix Mode = "prefix"
	// ModeContinuation treats any line matching the pattern as a CONTINUATION of the current event.
	ModeContinuation Mode = "continuation"
)

// Config holds options for the Combiner.
type Config struct {
	Mode    Mode
	Pattern string
	MaxLines int
	Timeout  time.Duration
}

// Combiner accumulates lines and emits complete multi-line events.
type Combiner struct {
	cfg     Config
	re      *regexp.Regexp
	buf     []string
	LastFlush time.Time
}

// New validates cfg and returns a ready Combiner.
func New(cfg Config) (*Combiner, error) {
	if cfg.Mode != ModePrefix && cfg.Mode != ModeContinuation {
		return nil, errors.New("multiline: mode must be \"prefix\" or \"continuation\"")
	}
	if strings.TrimSpace(cfg.Pattern) == "" {
		return nil, errors.New("multiline: pattern must not be empty")
	}
	re, err := regexp.Compile(cfg.Pattern)
	if err != nil {
		return nil, errors.New("multiline: invalid pattern: " + err.Error())
	}
	if cfg.MaxLines < 0 {
		return nil, errors.New("multiline: max_lines must be >= 0")
	}
	if cfg.Timeout < 0 {
		return nil, errors.New("multiline: timeout must be >= 0")
	}
	return &Combiner{cfg: cfg, re: re, LastFlush: time.Now()}, nil
}

// Push feeds a new raw line into the combiner.
// It returns a complete event string and true when one is ready, or "", false
// when the line was buffered and more input is needed.
func (c *Combiner) Push(line string) (string, bool) {
	switch c.cfg.Mode {
	case ModePrefix:
		if c.re.MatchString(line) {
			return c.swapBuf(line)
		}
		c.buf = append(c.buf, line)
	case ModeContinuation:
		if c.re.MatchString(line) {
			c.buf = append(c.buf, line)
		} else {
			return c.swapBuf(line)
		}
	}
	if c.cfg.MaxLines > 0 && len(c.buf) >= c.cfg.MaxLines {
		return c.Flush()
	}
	return "", false
}

// Flush drains whatever is buffered and returns it as a complete event.
func (c *Combiner) Flush() (string, bool) {
	if len(c.buf) == 0 {
		return "", false
	}
	event := strings.Join(c.buf, "\n")
	c.buf = c.buf[:0]
	c.LastFlush = time.Now()
	return event, true
}

// swapBuf emits the current buffer as an event and starts a fresh one with line.
func (c *Combiner) swapBuf(line string) (string, bool) {
	prev := c.buf
	c.buf = []string{line}
	if len(prev) == 0 {
		return "", false
	}
	c.LastFlush = time.Now()
	return strings.Join(prev, "\n"), true
}
