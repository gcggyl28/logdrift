// Package jitter adds configurable random jitter to inter-line delays,
// useful when replaying or simulating log streams to avoid thundering-herd
// effects in downstream consumers.
package jitter

import (
	"fmt"
	"math/rand"
	"time"
)

// Jitter computes a randomised delay within [0, max].
type Jitter struct {
	max     time.Duration
	rng     *rand.Rand
	enabled bool
}

// New creates a Jitter with the given maximum spread.
// A zero max disables jitter (Delay always returns 0).
// A negative max returns an error.
func New(max time.Duration) (*Jitter, error) {
	if max < 0 {
		return nil, fmt.Errorf("jitter: max must be >= 0, got %s", max)
	}
	return &Jitter{
		max:     max,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		enabled: max > 0,
	}, nil
}

// Enabled reports whether jitter will produce non-zero delays.
func (j *Jitter) Enabled() bool { return j.enabled }

// Max returns the configured maximum spread.
func (j *Jitter) Max() time.Duration { return j.max }

// Delay returns a random duration in [0, max].
// If jitter is disabled it always returns 0.
func (j *Jitter) Delay() time.Duration {
	if !j.enabled {
		return 0
	}
	return time.Duration(j.rng.Int63n(int64(j.max) + 1))
}

// Sleep blocks for a random duration in [0, max].
func (j *Jitter) Sleep() {
	if d := j.Delay(); d > 0 {
		time.Sleep(d)
	}
}
