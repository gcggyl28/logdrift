// Package zeitgeist provides time-based bucketing of log entries into
// fixed-duration slots for trend analysis across services.
package zeitgeist

import (
	"errors"
	"sync"
	"time"
)

// Bucket holds aggregated counts for a single time slot.
type Bucket struct {
	Start  time.Time
	Counts map[string]int // keyed by service
}

// Bucketer partitions log lines into fixed-width time buckets.
type Bucketer struct {
	mu       sync.Mutex
	width    time.Duration
	buckets  []*Bucket
	maxSlots int
}

// New creates a Bucketer with the given slot width and maximum retained slots.
func New(width time.Duration, maxSlots int) (*Bucketer, error) {
	if width <= 0 {
		return nil, errors.New("zeitgeist: slot width must be positive")
	}
	if maxSlots <= 0 {
		return nil, errors.New("zeitgeist: maxSlots must be positive")
	}
	return &Bucketer{width: width, maxSlots: maxSlots}, nil
}

// Record increments the counter for service in the bucket that contains t.
func (b *Bucketer) Record(service string, t time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	slot := t.Truncate(b.width)
	for _, bk := range b.buckets {
		if bk.Start.Equal(slot) {
			bk.Counts[service]++
			return
		}
	}
	bk := &Bucket{Start: slot, Counts: map[string]int{service: 1}}
	b.buckets = append(b.buckets, bk)
	if len(b.buckets) > b.maxSlots {
		b.buckets = b.buckets[len(b.buckets)-b.maxSlots:]
	}
}

// Snapshot returns a copy of all retained buckets ordered oldest-first.
func (b *Bucketer) Snapshot() []Bucket {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]Bucket, len(b.buckets))
	for i, bk := range b.buckets {
		copy := Bucket{Start: bk.Start, Counts: make(map[string]int, len(bk.Counts))}
		for k, v := range bk.Counts {
			copy.Counts[k] = v
		}
		out[i] = copy
	}
	return out
}
