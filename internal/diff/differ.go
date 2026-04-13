package diff

import (
	"fmt"
	"strings"
)

// Mode defines how log lines are compared across services.
type Mode string

const (
	ModeUnified Mode = "unified"
	ModeSideBySide Mode = "side-by-side"
	ModeTimestamp Mode = "timestamp"
)

// Entry represents a single log line from a named service.
type Entry struct {
	Service string
	Line    string
}

// Result holds the output of a diff operation.
type Result struct {
	Left    Entry
	Right   Entry
	Drifted bool
	Delta   string
}

// Differ compares log entries across services.
type Differ struct {
	mode Mode
}

// New creates a new Differ with the given mode.
func New(mode Mode) (*Differ, error) {
	switch mode {
	case ModeUnified, ModeSideBySide, ModeTimestamp:
		return &Differ{mode: mode}, nil
	default:
		return nil, fmt.Errorf("unknown diff mode: %q", mode)
	}
}

// Compare takes two log entries and returns a Result indicating drift.
func (d *Differ) Compare(a, b Entry) Result {
	if a.Line == b.Line {
		return Result{Left: a, Right: b, Drifted: false}
	}

	var delta string
	switch d.mode {
	case ModeSideBySide:
		delta = fmt.Sprintf("[%s] %s  |  [%s] %s", a.Service, a.Line, b.Service, b.Line)
	case ModeTimestamp:
		delta = fmt.Sprintf("timestamp drift between [%s] and [%s]", a.Service, b.Service)
	default: // unified
		delta = unifiedDelta(a, b)
	}

	return Result{Left: a, Right: b, Drifted: true, Delta: delta}
}

func unifiedDelta(a, b Entry) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("--- %s\n", a.Service))
	sb.WriteString(fmt.Sprintf("-  %s\n", a.Line))
	sb.WriteString(fmt.Sprintf("+++ %s\n", b.Service))
	sb.WriteString(fmt.Sprintf("+  %s", b.Line))
	return sb.String()
}
