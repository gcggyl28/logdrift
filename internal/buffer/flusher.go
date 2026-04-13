package buffer

import (
	"context"
	"time"
)

// FlushFunc is called by Flusher with a snapshot of buffered entries.
type FlushFunc func(entries []Entry)

// Flusher periodically drains a RingBuffer and calls a FlushFunc.
type Flusher struct {
	buf      *RingBuffer
	interval time.Duration
	fn       FlushFunc
}

// NewFlusher creates a Flusher that calls fn every interval with buffered entries.
// interval must be > 0.
func NewFlusher(buf *RingBuffer, interval time.Duration, fn FlushFunc) (*Flusher, error) {
	if interval <= 0 {
		return nil, ErrInvalidCapacity // reuse sentinel; callers check > 0
	}
	if fn == nil {
		panic("buffer: FlushFunc must not be nil")
	}
	return &Flusher{buf: buf, interval: interval, fn: fn}, nil
}

// Run starts the flush loop. It blocks until ctx is cancelled.
func (f *Flusher) Run(ctx context.Context) {
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			// final flush before exit
			if entries := f.buf.Entries(); len(entries) > 0 {
				f.fn(entries)
			}
			return
		case <-ticker.C:
			if entries := f.buf.Entries(); len(entries) > 0 {
				f.fn(entries)
				f.buf.Reset()
			}
		}
	}
}
