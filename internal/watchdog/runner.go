package watchdog

import (
	"context"
	"errors"
)

// LogEntry is the minimal interface the Runner needs from upstream.
type LogEntry interface {
	ServiceName() string
}

// Runner wires a source channel of LogEntry values into a Watchdog,
// calling Ping on each entry and forwarding silence Events downstream.
type Runner struct {
	dog *Watchdog
	src <-chan LogEntry
}

// NewRunner creates a Runner.
// dog and src must both be non-nil.
func NewRunner(dog *Watchdog, src <-chan LogEntry) (*Runner, error) {
	if dog == nil {
		return nil, errors.New("watchdog: runner requires a non-nil Watchdog")
	}
	if src == nil {
		return nil, errors.New("watchdog: runner requires a non-nil source channel")
	}
	return &Runner{dog: dog, src: src}, nil
}

// Run drains src, pings the watchdog for each entry, and returns the silence
// event channel produced by the watchdog.
// The goroutine exits when ctx is cancelled or src is closed.
func (r *Runner) Run(ctx context.Context) <-chan Event {
	events := r.dog.Run(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-r.src:
				if !ok {
					return
				}
				r.dog.Ping(entry.ServiceName())
			}
		}
	}()
	return events
}
