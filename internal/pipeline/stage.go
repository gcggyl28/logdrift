// Package pipeline wires together the per-line processing stages
// (filter, redact, normalize, transform, enrich) into a single ordered
// chain that can be applied to every log entry before it reaches the
// diff or render layer.
package pipeline

import (
	"errors"
	"fmt"
)

// Entry is the unit of work flowing through the pipeline.
type Entry struct {
	Service string
	Line    string
	Fields  map[string]string
}

// StageFn is a function that mutates or drops an Entry.
// Returning (nil, nil) signals the entry should be dropped.
type StageFn func(e Entry) (*Entry, error)

// Stage wraps a named StageFn.
type Stage struct {
	Name string
	fn   StageFn
}

// NewStage creates a Stage with the given name and function.
func NewStage(name string, fn StageFn) (Stage, error) {
	if name == "" {
		return Stage{}, errors.New("pipeline: stage name must not be empty")
	}
	if fn == nil {
		return Stage{}, fmt.Errorf("pipeline: stage %q has nil function", name)
	}
	return Stage{Name: name, fn: fn}, nil
}

// Apply runs the stage function against e.
func (s Stage) Apply(e Entry) (*Entry, error) {
	return s.fn(e)
}
