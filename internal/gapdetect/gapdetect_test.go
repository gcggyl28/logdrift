package gapdetect

import (
	"testing"
	"time"
)

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestNew_PositiveThreshold(t *testing.T) {
	d, err := New(5 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestNew_ZeroThresholdReturnsError(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_NegativeThresholdReturnsError(t *testing.T) {
	_, err := New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative threshold")
	}
}

func TestRecord_FirstEntryNoGap(t *testing.T) {
	d, _ := New(5 * time.Second)
	ev, err := d.Record("svc", base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev != nil {
		t.Fatalf("expected nil event for first entry, got %v", ev)
	}
}

func TestRecord_SmallGapNoEvent(t *testing.T) {
	d, _ := New(5 * time.Second)
	d.Record("svc", base)
	ev, err := d.Record("svc", base.Add(2*time.Second))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev != nil {
		t.Fatalf("expected nil event for small gap, got %v", ev)
	}
}

func TestRecord_LargeGapEmitsEvent(t *testing.T) {
	d, _ := New(5 * time.Second)
	d.Record("svc", base)
	ev, err := d.Record("svc", base.Add(10*time.Second))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev == nil {
		t.Fatal("expected gap event")
	}
	if ev.Service != "svc" {
		t.Errorf("service: got %q, want %q", ev.Service, "svc")
	}
	if ev.Duration != 10*time.Second {
		t.Errorf("duration: got %v, want %v", ev.Duration, 10*time.Second)
	}
}

func TestRecord_EmptyServiceReturnsError(t *testing.T) {
	d, _ := New(time.Second)
	_, err := d.Record("", base)
	if err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestRecord_ZeroTimestampReturnsError(t *testing.T) {
	d, _ := New(time.Second)
	_, err := d.Record("svc", time.Time{})
	if err == nil {
		t.Fatal("expected error for zero timestamp")
	}
}

func TestReset_ClearsState(t *testing.T) {
	d, _ := New(5 * time.Second)
	d.Record("svc", base)
	d.Reset("svc")
	ev, err := d.Record("svc", base.Add(10*time.Second))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev != nil {
		t.Fatalf("expected nil after reset, got %v", ev)
	}
}

func TestReset_EmptyServiceReturnsError(t *testing.T) {
	d, _ := New(time.Second)
	if err := d.Reset(""); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestGapEvent_Summary(t *testing.T) {
	ev := &GapEvent{
		Service:  "api",
		From:     base,
		To:       base.Add(30 * time.Second),
		Duration: 30 * time.Second,
	}
	s := ev.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}

func TestRecord_IndependentServices(t *testing.T) {
	d, _ := New(5 * time.Second)
	d.Record("a", base)
	d.Record("b", base)

	// large gap for a
	ev, _ := d.Record("a", base.Add(20*time.Second))
	if ev == nil {
		t.Fatal("expected gap event for service a")
	}
	// small gap for b
	ev2, _ := d.Record("b", base.Add(1*time.Second))
	if ev2 != nil {
		t.Fatalf("expected no gap event for service b, got %v", ev2)
	}
}
