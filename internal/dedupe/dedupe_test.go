package dedupe

import (
	"testing"
	"time"
)

func TestNew_ZeroDisablesDedup(t *testing.T) {
	d, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Deduper")
	}
}

func TestNew_NegativeReturnsError(t *testing.T) {
	_, err := New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNew_PositiveWindow(t *testing.T) {
	d, err := New(5 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Deduper")
	}
}

func TestAllow_DisabledAlwaysTrue(t *testing.T) {
	d, _ := New(0)
	for i := 0; i < 5; i++ {
		if !d.Allow("svc", "same line") {
			t.Fatalf("disabled deduper should always allow (iteration %d)", i)
		}
	}
}

func TestAllow_FirstOccurrenceAllowed(t *testing.T) {
	d, _ := New(time.Minute)
	if !d.Allow("svc", "hello") {
		t.Fatal("first occurrence should be allowed")
	}
}

func TestAllow_DuplicateSuppressed(t *testing.T) {
	d, _ := New(time.Minute)
	d.Allow("svc", "hello")
	if d.Allow("svc", "hello") {
		t.Fatal("duplicate within window should be suppressed")
	}
}

func TestAllow_DifferentLineAllowed(t *testing.T) {
	d, _ := New(time.Minute)
	d.Allow("svc", "hello")
	if !d.Allow("svc", "world") {
		t.Fatal("different line should be allowed")
	}
}

func TestAllow_IndependentServices(t *testing.T) {
	d, _ := New(time.Minute)
	d.Allow("svc-a", "line")
	if !d.Allow("svc-b", "line") {
		t.Fatal("same line for different service should be allowed")
	}
}

func TestCount_TracksSuppressions(t *testing.T) {
	d, _ := New(time.Minute)
	d.Allow("svc", "msg")
	d.Allow("svc", "msg")
	d.Allow("svc", "msg")
	if got := d.Count("svc"); got != 2 {
		t.Fatalf("expected 2 suppressions, got %d", got)
	}
}

func TestCount_UnknownServiceReturnsZero(t *testing.T) {
	d, _ := New(time.Minute)
	if got := d.Count("unknown"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsState(t *testing.T) {
	d, _ := New(time.Minute)
	d.Allow("svc", "msg")
	d.Allow("svc", "msg")
	d.Reset()
	if !d.Allow("svc", "msg") {
		t.Fatal("after reset, line should be allowed again")
	}
	if got := d.Count("svc"); got != 1 {
		t.Fatalf("expected count 1 after reset, got %d", got)
	}
}
