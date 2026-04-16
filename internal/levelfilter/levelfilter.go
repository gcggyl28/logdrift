// Package levelfilter provides severity-level-based filtering of log entries.
package levelfilter

import (
	"fmt"
	"strings"
)

// Level represents a numeric log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
}

// Filter drops log entries whose severity is below the configured minimum level.
type Filter struct {
	min      Level
	levelKey string
}

// New creates a Filter that passes only entries at or above minLevel.
// levelKey is the field name in the entry map that holds the severity string.
func New(minLevel, levelKey string) (*Filter, error) {
	if levelKey == "" {
		return nil, fmt.Errorf("levelfilter: levelKey must not be empty")
	}
	min, ok := levelNames[strings.ToLower(minLevel)]
	if !ok {
		return nil, fmt.Errorf("levelfilter: unknown level %q", minLevel)
	}
	return &Filter{min: min, levelKey: levelKey}, nil
}

// Allow returns true when the entry's level is at or above the minimum.
// Entries missing the level field are passed through.
func (f *Filter) Allow(fields map[string]string) bool {
	raw, ok := fields[f.levelKey]
	if !ok {
		return true
	}
	lvl, ok := levelNames[strings.ToLower(raw)]
	if !ok {
		return true
	}
	return lvl >= f.min
}
