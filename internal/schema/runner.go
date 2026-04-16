package schema

import (
	"context"
	"fmt"

	"github.com/user/logdrift/internal/parser"
)

// Result pairs an entry with any schema violations found.
type Result struct {
	Entry      parser.Entry
	Violations []Violation
}

// Runner streams entries from src through the Validator and emits Results.
type Runner struct {
	v   *Validator
	src <-chan parser.Entry
}

// NewRunner creates a Runner that validates entries from src.
func NewRunner(v *Validator, src <-chan parser.Entry) (*Runner, error) {
	if v == nil {
		return nil, fmt.Errorf("schema: validator must not be nil")
	}
	if src == nil {
		return nil, fmt.Errorf("schema: source channel must not be nil")
	}
	return &Runner{v: v, src: src}, nil
}

// Run reads entries from src, validates each, and sends Results to the returned
// channel. The channel is closed when ctx is cancelled or src is exhausted.
func (r *Runner) Run(ctx context.Context) <-chan Result {
	out := make(chan Result, 64)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-r.src:
				if !ok {
					return
				}
				violations := r.v.Validate(entry)
				select {
				case out <- Result{Entry: entry, Violations: violations}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
