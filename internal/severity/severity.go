// Package severity classifies log entries by severity level.
package severity

import (
	"fmt"
	"strings"

	"github.com/yourorg/logdrift/internal/parser"
)

// Level represents a log severity level.
type Level int

const (
	LevelUnknown Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// Classifier maps log entry fields to a severity Level.
type Classifier struct {
	fields []string // candidate field names, e.g. ["level","severity","lvl"]
}

// New creates a Classifier that inspects the given field names in order.
// At least one field name must be provided.
func New(fields ...string) (*Classifier, error) {
	for _, f := range fields {
		if strings.TrimSpace(f) == "" {
			return nil, fmt.Errorf("severity: field name must not be empty")
		}
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("severity: at least one field name required")
	}
	return &Classifier{fields: fields}, nil
}

// Classify returns the Level for the given parsed entry.
func (c *Classifier) Classify(e parser.Entry) Level {
	for _, f := range c.fields {
		v, ok := e.Fields[f]
		if !ok {
			continue
		}
		return parse.Sprintf("%v", v))
	}
	return LevelUnfunc parse(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug", "dbg", "trace":
		return LevelDebug
	case "info", "information":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error", "err":
		return LevelError
	case "fatal", "critical", "crit", "panic":
		return LevelFatal
	default:
		return LevelUnknown
	}
}
