package watchdog

import (
	"context"
	"testing"
	"time"
)

func TestNew_PositiveThreshold(t *testing.T) {
	w, err := New(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil Watchdog")
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

func TestPing_ResetsTimer(t *testing.T) {
	w, _ := New(50 * time.Millisecond)
	w.Register("svc-a")
	w.Ping("svc-a")
	// Ping just after registration; silence clock should reset — no event yet.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	events := w.Run(ctx)
	for e := range events {
		if e.Service == "svc-a" {
			t.Errorf("unexpected silence event for svc-a: %v", e)
		}
	}
}

func TestRun_EmitsSilenceEvent(t *testing.T) {
	w, _ := New(40 * time.Millisecond)
	w.Register("svc-b")
	// Backdate last-seen so the service appears silent immediately.
	w.mu.Lock()
	w.lastSeen["svc-b"] = time.Now().Add(-100 * time.Millisecond)
	w.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	events := w.Run(ctx)
	select {
	case e := <-events:
		if e.Service != "svc-b" {
			t.Errorf("expected svc-b, got %s", e.Service)
		}
		if e.SilentFor < 40*time.Millisecond {
			t.Errorf("expected SilentFor >= 40ms, got %v", e.SilentFor)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for silence event")
	}
}

func TestRun_CancelStops(t *testing.T) {
	w, _ := New(500 * time.Millisecond)
	w.Register("svc-c")
	ctx, cancel := context.WithCancel(context.Background())
	events := w.Run(ctx)
	cancel()
	// Channel must close after cancel.
	for range events {
	}
}

func TestRegister_IdempotentTimestamp(t *testing.T) {
	w, _ := New(100 * time.Millisecond)
	w.Register("svc-d")
	w.mu.Lock()
	first := w.lastSeen["svc-d"]
	w.mu.Unlock()

	time.Sleep(5 * time.Millisecond)
	w.Register("svc-d") // second call must not overwrite
	w.mu.Lock()
	second := w.lastSeen["svc-d"]
	w.mu.Unlock()

	if !first.Equal(second) {
		t.Error("Register overwrote existing timestamp")
	}
}
