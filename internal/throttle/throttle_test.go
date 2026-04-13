package throttle

import (
	"testing"
	"time"
)

func TestNew_ZeroDisablesThrottling(t *testing.T) {
	th, err := New(0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 1000; i++ {
		if !th.Allow("svc") {
			t.Fatal("expected all lines to pass when maxPerWindow=0")
		}
	}
}

func TestNew_NegativeReturnsError(t *testing.T) {
	_, err := New(-1, time.Second)
	if err == nil {
		t.Fatal("expected error for negative maxPerWindow")
	}
}

func TestNew_ZeroWindowWithPositiveMax(t *testing.T) {
	_, err := New(5, 0)
	if err == nil {
		t.Fatal("expected error for zero window with positive maxPerWindow")
	}
}

func TestAllow_BurstUpToMax(t *testing.T) {
	th, err := New(3, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if !th.Allow("svc") {
			t.Fatalf("expected Allow=true on call %d", i+1)
		}
	}
	if th.Allow("svc") {
		t.Fatal("expected Allow=false after burst exhausted")
	}
}

func TestAllow_IndependentServices(t *testing.T) {
	th, _ := New(2, time.Minute)
	th.Allow("a")
	th.Allow("a")
	if !th.Allow("b") {
		t.Fatal("service b should not be affected by service a's quota")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	th, err := New(1, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !th.Allow("svc") {
		t.Fatal("first call should be allowed")
	}
	if th.Allow("svc") {
		t.Fatal("second call within window should be denied")
	}
	time.Sleep(30 * time.Millisecond)
	if !th.Allow("svc") {
		t.Fatal("call after window expiry should be allowed")
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	th, _ := New(1, time.Minute)
	th.Allow("svc")
	if th.Allow("svc") {
		t.Fatal("should be denied before reset")
	}
	th.Reset()
	if !th.Allow("svc") {
		t.Fatal("should be allowed after reset")
	}
}
