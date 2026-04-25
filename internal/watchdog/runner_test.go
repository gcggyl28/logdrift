package watchdog

import (
	"context"
	"testing"
	"time"
)

type stubEntry struct{ svc string }

func (s stubEntry) ServiceName() string { return s.svc }

func TestNewRunner_NilWatchdogReturnsError(t *testing.T) {
	src := make(chan LogEntry)
	_, err := NewRunner(nil, src)
	if err == nil {
		t.Fatal("expected error for nil watchdog")
	}
}

func TestNewRunner_NilSrcReturnsError(t *testing.T) {
	w, _ := New(100 * time.Millisecond)
	_, err := NewRunner(w, nil)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestNewRunner_Valid(t *testing.T) {
	w, _ := New(100 * time.Millisecond)
	src := make(chan LogEntry)
	r, err := NewRunner(w, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Runner")
	}
}

func TestRunner_PingsOnEntry(t *testing.T) {
	w, _ := New(200 * time.Millisecond)
	w.Register("svc-x")
	// Backdate so it looks silent.
	w.mu.Lock()
	w.lastSeen["svc-x"] = time.Now().Add(-300 * time.Millisecond)
	w.mu.Unlock()

	src := make(chan LogEntry, 1)
	src <- stubEntry{svc: "svc-x"}

	r, _ := NewRunner(w, src)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	events := r.Run(ctx)

	// Give the runner goroutine time to process the ping.
	time.Sleep(20 * time.Millisecond)

	w.mu.Lock()
	last := w.lastSeen["svc-x"]
	w.mu.Unlock()

	if time.Since(last) > 100*time.Millisecond {
		t.Error("expected lastSeen to be updated by Ping")
	}
	_ = events
}

func TestRunner_CancelStops(t *testing.T) {
	w, _ := New(500 * time.Millisecond)
	src := make(chan LogEntry)
	r, _ := NewRunner(w, src)
	ctx, cancel := context.WithCancel(context.Background())
	events := r.Run(ctx)
	cancel()
	for range events {
	}
}
