// Package dropout provides a probabilistic line-dropping mechanism that
// sheds load by discarding a configurable fraction of log entries when
// the inbound rate exceeds a threshold.
package dropout

import (
	"errors"
	"math/rand"
	"sync/atomic"
)

// Dropper drops entries with probability Rate when enabled.
type Dropper struct {
	rate    float64 // 0.0 – 1.0; 0 disables
	dropped atomic.Uint64
	total   atomic.Uint64
}

// New returns a Dropper with the given drop rate.
// rate must be in [0.0, 1.0]; 0 disables dropping.
func New(rate float64) (*Dropper, error) {
	if rate < 0 || rate > 1 {
		return nil, errors.New("dropout: rate must be between 0.0 and 1.0")
	}
	return &Dropper{rate: rate}, nil
}

// Allow returns true if the entry should be forwarded.
func (d *Dropper) Allow() bool {
	d.total.Add(1)
	if d.rate == 0 {
		return true
	}
	if rand.Float64() < d.rate { //nolint:gosec
		d.dropped.Add(1)
		return false
	}
	return true
}

// Stats returns total seen and total dropped counts.
func (d *Dropper) Stats() (total, dropped uint64) {
	return d.total.Load(), d.dropped.Load()
}

// Rate returns the configured drop rate.
func (d *Dropper) Rate() float64 { return d.rate }
