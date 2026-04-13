package snapshot

import (
	"context"

	"github.com/yourorg/logdrift/internal/tail"
)

// Entry pairs a line of text with the service it came from.
type Entry struct {
	Service string
	Line    string
}

// Collector reads from a fan-in channel of tail.Lines, pushes each line
// into a Snapshot, and forwards entries to an output channel for
// downstream consumers (e.g. the diff pipeline).
type Collector struct {
	snap *Snapshot
	out  chan Entry
}

// NewCollector returns a Collector backed by snap.
// The returned output channel is closed when ctx is done.
func NewCollector(snap *Snapshot) *Collector {
	return &Collector{
		snap: snap,
		out:  make(chan Entry, 64),
	}
}

// Out returns the read-only channel of collected entries.
func (c *Collector) Out() <-chan Entry {
	return c.out
}

// Run consumes lines from src until ctx is cancelled or src is closed.
// Each line is stored in the snapshot and forwarded on Out().
func (c *Collector) Run(ctx context.Context, src <-chan tail.Line) {
	defer close(c.out)
	for {
		select {
		case <-ctx.Done():
			return
		case tl, ok := <-src:
			if !ok {
				return
			}
			c.snap.Push(tl.Service, tl.Text)
			select {
			case c.out <- Entry{Service: tl.Service, Line: tl.Text}:
			case <-ctx.Done():
				return
			}
		}
	}
}
