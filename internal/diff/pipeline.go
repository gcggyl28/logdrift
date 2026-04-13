// Package diff provides line-level diffing primitives and a streaming
// pipeline that groups lines from multiple services within a sliding time
// window before comparing them.
package diff

import (
	"context"
	"time"

	"github.com/user/logdrift/internal/tail"
)

// Result holds the output of a single diff comparison.
type Result struct {
	ServiceA string
	ServiceB string
	LineA    string
	LineB    string
	Delta    string
	HasDrift bool
}

// Pipeline buffers lines from a merged channel within a time window and emits
// diff Results for every pair of lines that arrive in the same window.
type Pipeline struct {
	differ Differ
}

// NewPipeline creates a Pipeline backed by the given Differ.
func NewPipeline(d Differ) *Pipeline {
	return &Pipeline{differ: d}
}

// Run consumes lines from src, groups them by window duration, and emits
// diff Results on the returned channel. The channel is closed when ctx is
// done or src is exhausted.
func (p *Pipeline) Run(ctx context.Context, src <-chan tail.Line, window time.Duration) <-chan Result {
	out := make(chan Result, 16)

	go func() {
		defer close(out)

		buf := make([]tail.Line, 0, 8)
		ticker := time.NewTicker(window)
		defer ticker.Stop()

		flush := func() {
			for i := 0; i+1 < len(buf); i += 2 {
				a, b := buf[i], buf[i+1]
				delta, hasDrift := p.differ.Compare(a.Text, b.Text)
				select {
				case out <- Result{
					ServiceA:  a.Service,
					ServiceB:  b.Service,
					LineA:     a.Text,
					LineB:     b.Text,
					Delta:     delta,
					HasDrift: hasDrift,
				}:
				case <-ctx.Done():
					return
				}
			}
			buf = buf[:0]
		}

		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-src:
				if !ok {
					flush()
					return
				}
				buf = append(buf, line)
			case <-ticker.C:
				flush()
			}
		}
	}()

	return out
}
