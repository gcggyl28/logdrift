package checkpoint

import (
	"context"
	"time"
)

// Tracker wraps a Store and periodically flushes updated offsets
// for a single service as lines are consumed.
type Tracker struct {
	store   *Store
	service string
	offset  int64
	interval time.Duration
}

// NewTracker creates a Tracker that records progress for service into store.
// interval controls how often the position is persisted; use 0 to persist
// on every Advance call.
func NewTracker(store *Store, service string, interval time.Duration) (*Tracker, error) {
	if store == nil {
		return nil, ErrNoCheckpoint
	}
	if service == "" {
		return nil, errEmptyService
	}
	t := &Tracker{store: store, service: service, interval: interval}
	if e, err := store.Get(service); err == nil {
		t.offset = e.Offset
	}
	return t, nil
}

// Offset returns the current tracked byte offset.
func (t *Tracker) Offset() int64 { return t.offset }

// Advance records that n additional bytes have been consumed and persists
// if no interval is set (synchronous mode).
func (t *Tracker) Advance(n int64) error {
	t.offset += n
	if t.interval == 0 {
		return t.store.Set(t.service, t.offset)
	}
	return nil
}

// Run starts the periodic flush loop. It blocks until ctx is cancelled.
func (t *Tracker) Run(ctx context.Context) error {
	if t.interval <= 0 {
		<-ctx.Done()
		return ctx.Err()
	}
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := t.store.Set(t.service, t.offset); err != nil {
				return err
			}
		case <-ctx.Done():
			// Final flush on shutdown.
			_ = t.store.Set(t.service, t.offset)
			return ctx.Err()
		}
	}
}

// sentinel so we do not import fmt.
type cpErr string

func (e cpErr) Error() string { return string(e) }

const errEmptyService cpErr = "checkpoint: service name must not be empty"
