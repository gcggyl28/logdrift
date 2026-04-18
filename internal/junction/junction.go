// Package junction merges multiple named log entry channels into a single
// stream, tagging each entry with its source service name.
package junction

import (
	"context"
	"sync"

	"github.com/user/logdrift/internal/diff"
)

// Entry is a log line annotated with its originating service.
type Entry struct {
	Service string
	Line    string
}

// Junction fans in named sources into one output channel.
type Junction struct {
	sources map[string]<-chan diff.Line
}

// New creates a Junction from a map of service name → line channel.
// Returns an error when sources is empty.
func New(sources map[string]<-chan diff.Line) (*Junction, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("junction: at least one source is required")
	}
	return &Junction{sources: sources}, nil
}

// Run merges all source channels into the returned channel.
// The output channel is closed when all sources are drained or ctx is cancelled.
func (j *Junction) Run(ctx context.Context) <-chan Entry {
	out := make(chan Entry, 64)
	var wg sync.WaitGroup

	for svc, ch := range j.sources {
		wg.Add(1)
		go func(service string, src <-chan diff.Line) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case line, ok := <-src:
					if !ok {
						return
					}
					select {
					case out <- Entry{Service: service, Line: line.Text}:
					case <-ctx.Done():
						return
					}
				}
			}
		}(svc, ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
