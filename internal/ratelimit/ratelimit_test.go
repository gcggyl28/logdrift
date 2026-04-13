package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/ratelimit"
)

func TestNew_ZeroDisablesLimiting(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{LinesPerSecond: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !l.Disabled() {
		t.Fatal("expected limiter to be disabled when LinesPerSecond=0")
	}
}

func TestNew_NegativeReturnsError(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{LinesPerSecond: -1})
	if err == nil {
		t.Fatal("expected error for negative LinesPerSecond")
	}
}

func TestNew_PositiveIsEnabled(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{LinesPerSecond: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Disabled() {
		t.Fatal("expected limiter to be enabled")
	}
}

func TestAllow_DisabledAlwaysTrue(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{LinesPerSecond: 0})
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatal("disabled limiter should always allow")
		}
	}
}

func TestAllow_BurstUpToCapacity(t *testing.T) {
	const rate = 5
	l, _ := ratelimit.New(ratelimit.Config{LinesPerSecond: rate})

	allowed := 0
	for i := 0; i < rate*2; i++ {
		if l.Allow() {
			allowed++
		}
	}
	// Initial tokens == rate, so exactly `rate` calls should succeed immediately.
	if allowed != rate {
		t.Fatalf("expected %d allowed, got %d", rate, allowed)
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{LinesPerSecond: 100})

	// Drain the bucket.
	for i := 0; i < 100; i++ {
		l.Allow()
	}
	// Should be empty now.
	if l.Allow() {
		t.Fatal("bucket should be empty after draining")
	}

	// Wait long enough to refill at least one token.
	time.Sleep(20 * time.Millisecond)

	if !l.Allow() {
		t.Fatal("expected at least one token after sleep")
	}
}
