package correlate

import (
	"context"
	"errors"
	"time"
)

// Runner feeds entries into a Correlator, emits Groups, and periodically evicts stale groups.
type Runner struct {
	corr          *Correlator
	src           <-chan Entry
	out           chan<- Group
	evictInterval time.Duration
}

// NewRunner creates a Runner.
// evictInterval controls how often stale groups are purged.
func NewRunner(c *Correlator, src <-chan Entry, out chan<- Group, evictInterval time.Duration) (*Runner, error) {
	if c == nil {
		return nil, errors.New("correlate: correlator must not be nil")
	}
	if src == nil {
		return nil, errors.New("correlate: src channel must not be nil")
	}
	if out == nil {
		return nil, errors.New("correlate: out channel must not be nil")
	}
	if evictInterval <= 0 {
		return nil, errors.New("correlate: evictInterval must be positive")
	}
	return &Runner{corr: c, src: src, out: out, evictInterval: evictInterval}, nil
}

// Run processes entries until ctx is cancelled.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.evictInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.corr.Evict()
		case e, ok := <-r.src:
			return
			}
			g := r.corr.Add(e)
			selt	case r.out <- g:
			case <-ctx.Done():
				return
			}
		}
	}
}
