package jitter

import (
	"testing"
	"time"
)

func TestNew_ZeroDisablesJitter(t *testing.T) {
	j, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j.Enabled() {
		t.Fatal("expected jitter to be disabled for zero max")
	}
}

func TestNew_NegativeReturnsError(t *testing.T) {
	_, err := New(-1 * time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative max")
	}
}

func TestNew_PositiveIsEnabled(t *testing.T) {
	j, err := New(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !j.Enabled() {
		t.Fatal("expected jitter to be enabled for positive max")
	}
	if j.Max() != 100*time.Millisecond {
		t.Fatalf("expected max 100ms, got %s", j.Max())
	}
}

func TestDelay_DisabledAlwaysZero(t *testing.T) {
	j, _ := New(0)
	for i := 0; i < 20; i++ {
		if d := j.Delay(); d != 0 {
			t.Fatalf("expected 0, got %s", d)
		}
	}
}

func TestDelay_WithinBounds(t *testing.T) {
	max := 50 * time.Millisecond
	j, err := New(max)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 200; i++ {
		d := j.Delay()
		if d < 0 || d > max {
			t.Fatalf("delay %s out of bounds [0, %s]", d, max)
		}
	}
}

func TestDelay_ProducesVariance(t *testing.T) {
	j, _ := New(1 * time.Second)
	seen := make(map[time.Duration]struct{})
	for i := 0; i < 50; i++ {
		seen[j.Delay()] = struct{}{}
	}
	if len(seen) < 2 {
		t.Fatal("expected variance in delays, got none")
	}
}
