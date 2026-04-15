package debounce

import (
	"context"
)

// Source is a channel of raw entries to be debounced.
type Source <-chan Entry

// Runner reads from src, debounces each entry, and forwards results to fn.
type Runner struct {
	d   *Debouncer
	src Source
}

// NewRunner creates a Runner that feeds src into d.
func NewRunner(d *Debouncer, src Source) *Runner {
	return &Runner{d: d, src: src}
}

// Run starts the read loop and the drain loop.
// It blocks until ctx is cancelled.
func (r *Runner) Run(ctx context.Context, fn func(Entry)) {
	go func() {
		for {
			select {
			case e, ok := <-r.src:
				if !ok {
					return
				}
				r.d.Push(e)
			case <-ctx.Done():
				return
			}
		}
	}()
	r.d.Drain(ctx, fn)
}
