package pipeline

import (
	"errors"
	"fmt"
)

// Pipeline is an ordered sequence of Stages applied to each Entry.
type Pipeline struct {
	stages []Stage
}

// New constructs a Pipeline from the supplied stages.
// At least one stage is required.
func New(stages []Stage) (*Pipeline, error) {
	if len(stages) == 0 {
		return nil, errors.New("pipeline: at least one stage is required")
	}
	copy := make([]Stage, len(stages))
	for i, s := range stages {
		if s.Name == "" {
			return nil, fmt.Errorf("pipeline: stage at index %d has empty name", i)
		}
		copy[i] = s
	}
	return &Pipeline{stages: copy}, nil
}

// Run passes e through every stage in order.
// If any stage returns nil the entry is dropped and (nil, nil) is returned.
// The first non-nil error short-circuits the chain.
func (p *Pipeline) Run(e Entry) (*Entry, error) {
	cur := &e
	for _, s := range p.stages {
		out, err := s.Apply(*cur)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %q: %w", s.Name, err)
		}
		if out == nil {
			return nil, nil
		}
		cur = out
	}
	return cur, nil
}

// Stages returns a copy of the stage slice (for inspection/testing).
func (p *Pipeline) Stages() []Stage {
	out := make([]Stage, len(p.stages))
	copy(out, p.stages)
	return out
}
