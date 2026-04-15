package multiline

import (
	"context"
	"time"
)

// TimeoutFlusher periodically flushes a Combiner when no new lines have
// arrived within the configured Timeout, preventing incomplete events from
// being held indefinitely.
type TimeoutFlusher struct {
	combiner *Combiner
	out      chan<- string
}

// NewTimeoutFlusher returns a TimeoutFlusher that writes flushed events to out.
// It returns an error when the combiner's Timeout is zero (disabled).
func NewTimeoutFlusher(c *Combiner, out chan<- string) (*TimeoutFlusher, error) {
	if c.cfg.Timeout <= 0 {
		return nil, errorf("multiline: timeout must be > 0 for TimeoutFlusher")
	}
	return &TimeoutFlusher{combiner: c, out: out}, nil
}

// Run starts the background flush loop and blocks until ctx is cancelled.
func (tf *TimeoutFlusher) Run(ctx context.Context) {
	ticker := time.NewTicker(tf.combiner.cfg.Timeout / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			if ev, ok := tf.combiner.Flush(); ok {
				select {
				case tf.out <- ev:
				default:
				}
			}
			return
		case <-ticker.C:
			if time.Since(tf.combiner.LastFlush) >= tf.combiner.cfg.Timeout {
				if ev, ok := tf.combiner.Flush(); ok {
					select {
					case tf.out <- ev:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}
}

// errorf is a small helper to avoid importing fmt in this file.
func errorf(msg string) error {
	return &multilineError{msg: msg}
}

type multilineError struct{ msg string }

func (e *multilineError) Error() string { return e.msg }
