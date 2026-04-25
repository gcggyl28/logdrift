// Package watchdog monitors service log sources for silence — periods where
// no log lines are emitted — and fires an alert when the silence exceeds a
// configurable threshold.
package watchdog

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Event is emitted when a service has been silent for longer than the
// configured threshold.
type Event struct {
	Service  string
	SilentFor time.Duration
	At        time.Time
}

// Watchdog tracks last-seen timestamps per service and fires events on silence.
type Watchdog struct {
	threshold time.Duration
	mu        sync.Mutex
	lastSeen  map[string]time.Time
}

// New creates a Watchdog with the given silence threshold.
// threshold must be positive.
func New(threshold time.Duration) (*Watchdog, error) {
	if threshold <= 0 {
		return nil, errors.New("watchdog: threshold must be positive")
	}
	return &Watchdog{
		threshold: threshold,
		lastSeen:  make(map[string]time.Time),
	}, nil
}

// Ping records that a line was received from the given service.
func (w *Watchdog) Ping(service string) {
	w.mu.Lock()
	w.lastSeen[service] = time.Now()
	w.mu.Unlock()
}

// Register adds a service to be watched without recording a line.
// The silence clock starts from the moment of registration.
func (w *Watchdog) Register(service string) {
	w.mu.Lock()
	if _, ok := w.lastSeen[service]; !ok {
		w.lastSeen[service] = time.Now()
	}
	w.mu.Unlock()
}

// Run starts the watchdog loop, sending Events to the returned channel
// whenever a registered service has been silent for longer than the threshold.
// The loop ticks at threshold/2 for responsiveness.
func (w *Watchdog) Run(ctx context.Context) <-chan Event {
	out := make(chan Event, 16)
	tick := w.threshold / 2
	if tick < time.Millisecond {
		tick = time.Millisecond
	}
	go func() {
		defer close(out)
		t := time.NewTicker(tick)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-t.C:
				w.mu.Lock()
				for svc, last := range w.lastSeen {
					if silent := now.Sub(last); silent >= w.threshold {
						select {
						case out <- Event{Service: svc, SilentFor: silent, At: now}:
						default:
						}
					}
				}
				w.mu.Unlock()
			}
		}
	}()
	return out
}
