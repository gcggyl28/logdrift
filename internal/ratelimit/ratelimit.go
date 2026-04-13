// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently log lines are forwarded downstream in a logdrift session.
package ratelimit

import (
	"errors"
	"time"
)

// Config holds the parameters for the rate limiter.
type Config struct {
	// LinesPerSecond is the maximum number of lines allowed per second.
	// A value of 0 disables rate limiting.
	LinesPerSecond int
}

// Limiter controls throughput of log lines using a token-bucket algorithm.
type Limiter struct {
	tokens   float64
	max      float64
	rate     float64 // tokens per nanosecond
	lastTick time.Time
	disabled bool
}

// New creates a Limiter from cfg. Returns an error if LinesPerSecond is negative.
func New(cfg Config) (*Limiter, error) {
	if cfg.LinesPerSecond < 0 {
		return nil, errors.New("ratelimit: LinesPerSecond must be >= 0")
	}
	if cfg.LinesPerSecond == 0 {
		return &Limiter{disabled: true}, nil
	}
	max := float64(cfg.LinesPerSecond)
	return &Limiter{
		tokens:   max,
		max:      max,
		rate:     max / float64(time.Second),
		lastTick: time.Now(),
	}, nil
}

// Allow reports whether a single line may pass through at the current time.
// It refills tokens based on elapsed time since the last call.
func (l *Limiter) Allow() bool {
	if l.disabled {
		return true
	}
	now := time.Now()
	elapsed := float64(now.Sub(l.lastTick))
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}
	if l.tokens >= 1.0 {
		l.tokens -= 1.0
		return true
	}
	return false
}

// Disabled reports whether rate limiting is turned off.
func (l *Limiter) Disabled() bool {
	return l.disabled
}
