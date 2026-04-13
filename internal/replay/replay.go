// Package replay provides utilities for replaying historical log lines
// from a snapshot into a diff pipeline, enabling post-hoc drift analysis.
package replay

import (
	"context"
	"time"

	"github.com/user/logdrift/internal/snapshot"
)

// Entry represents a single replayed log line attributed to a service.
type Entry struct {
	Service string
	Line    string
	At      time.Time
}

// Replayer reads lines from a Snapshot and emits them in order.
type Replayer struct {
	snap    *snapshot.Snapshot
	delay   time.Duration
	services []string
}

// New creates a Replayer for the given snapshot.
// delay controls the pause between emitted entries (0 for no delay).
func New(snap *snapshot.Snapshot, delay time.Duration) *Replayer {
	return &Replayer{
		snap:  snap,
		delay: delay,
		services: snap.Services(),
	}
}

// Run emits all buffered log entries across services on the returned channel.
// Entries are interleaved in round-robin order across services.
// The channel is closed when all entries have been emitted or ctx is cancelled.
func (r *Replayer) Run(ctx context.Context) <-chan Entry {
	out := make(chan Entry, 64)
	go func() {
		defer close(out)
		lines := make(map[string][]string, len(r.services))
		maxLen := 0
		for _, svc := range r.services {
			l := r.snap.Lines(svc)
			lines[svc] = l
			if len(l) > maxLen {
				maxLen = len(l)
			}
		}
		for i := 0; i < maxLen; i++ {
			for _, svc := range r.services {
				if i >= len(lines[svc]) {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case out <- Entry{Service: svc, Line: lines[svc][i], At: time.Now()}:
				}
				if r.delay > 0 {
					select {
					case <-ctx.Done():
						return
					case <-time.After(r.delay):
					}
				}
			}
		}
	}()
	return out
}
