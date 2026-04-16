// Package fanout broadcasts log entries from a single source channel
// to multiple subscriber channels.
package fanout

import (
	"context"
	"sync"

	"github.com/user/logdrift/internal/tail"
)

// Broadcaster reads from a source and writes each entry to all registered
// subscribers. Slow subscribers are dropped rather than blocking the broadcast.
type Broadcaster struct {
	src   <-chan tail.Entry
	mu    sync.RWMutex
	subs  []chan tail.Entry
	bufSz int
}

// New creates a Broadcaster that reads from src.
// bufSz is the per-subscriber channel buffer size; must be >= 1.
func New(src <-chan tail.Entry, bufSz int) (*Broadcaster, error) {
	if bufSz < 1 {
		return nil, fmt.Errorf("fanout: bufSz must be >= 1, got %d", bufSz)
	}
	return &Broadcaster{src: src, bufSz: bufSz}, nil
}

// Subscribe returns a new channel that will receive all future entries.
func (b *Broadcaster) Subscribe() <-chan tail.Entry {
	ch := make(chan tail.Entry, b.bufSz)
	b.mu.Lock()
	b.subs = append(b.subs, ch)
	b.mu.Unlock()
	return ch
}

// Run starts broadcasting until ctx is cancelled or src is closed.
func (b *Broadcaster) Run(ctx context.Context) {
	defer b.closeAll()
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-b.src:
			if !ok {
				return
			}
			b.mu.RLock()
			for _, ch := range b.subs {
				select {
				case ch <- entry:
				default:
				}
			}
			b.mu.RUnlock()
		}
	}
}

func (b *Broadcaster) closeAll() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subs {
		close(ch)
	}
	b.subs = nil
}
