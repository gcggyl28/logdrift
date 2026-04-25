package stutter

import (
	"testing"
	"time"
)

func TestNew_ZeroWindowReturnsError(t *testing.T) {
	_, err := New(0, 3)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_NegativeWindowReturnsError(t *testing.T) {
	_, err := New(-time.Second, 3)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNew_ThresholdOneLessThanMinReturnsError(t *testing.T) {
	_, err := New(time.Second, 1)
	if err == nil {
		t.Fatal("expected error for threshold < 2")
	}
}

func TestNew_Valid(t *testing.T) {
	d, err := New(time.Second, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestRecord_BelowThresholdNoEvent(t *testing.T) {
	d, _ := New(time.Minute, 3)
	now := time.Now()
	if ev := d.Record("svc", "msg", now); ev != nil {
		t.Fatalf("expected nil, got event on first occurrence")
	}
	if ev := d.Record("svc", "msg", now.Add(time.Second)); ev != nil {
		t.Fatalf("expected nil below threshold")
	}
}

func TestRecord_AtThresholdEmitsEvent(t *testing.T) {
	d, _ := New(time.Minute, 3)
	now := time.Now()
	d.Record("svc", "err", now)
	d.Record("svc", "err", now.Add(time.Second))
	ev := d.Record("svc", "err", now.Add(2*time.Second))
	if ev == nil {
		t.Fatal("expected event at threshold")
	}
	if ev.Count != 3 {
		t.Fatalf("expected count 3, got %d", ev.Count)
	}
	if ev.Service != "svc" {
		t.Fatalf("expected service svc, got %s", ev.Service)
	}
}

func TestRecord_ContinuesEmittingAfterThreshold(t *testing.T) {
	d, _ := New(time.Minute, 2)
	now := time.Now()
	d.Record("svc", "loop", now)
	ev1 := d.Record("svc", "loop", now.Add(time.Second))
	ev2 := d.Record("svc", "loop", now.Add(2*time.Second))
	if ev1 == nil || ev2 == nil {
		t.Fatal("expected events after threshold reached")
	}
	if ev2.Count != 3 {
		t.Fatalf("expected count 3, got %d", ev2.Count)
	}
}

func TestRecord_WindowExpiredResetsState(t *testing.T) {
	d, _ := New(500*time.Millisecond, 2)
	now := time.Now()
	d.Record("svc", "msg", now)
	// second occurrence is outside the window
	ev := d.Record("svc", "msg", now.Add(time.Second))
	if ev != nil {
		t.Fatal("expected nil after window expiry reset")
	}
}

func TestRecord_DifferentLinesIndependent(t *testing.T) {
	d, _ := New(time.Minute, 2)
	now := time.Now()
	d.Record("svc", "alpha", now)
	ev := d.Record("svc", "beta", now.Add(time.Millisecond))
	if ev != nil {
		t.Fatal("different lines should not trigger stutter")
	}
}

func TestReset_ClearsServiceState(t *testing.T) {
	d, _ := New(time.Minute, 2)
	now := time.Now()
	d.Record("svc", "err", now)
	d.Reset("svc")
	// after reset the count starts fresh — threshold should not be reached
	ev := d.Record("svc", "err", now.Add(time.Second))
	if ev != nil {
		t.Fatal("expected nil after reset")
	}
}
