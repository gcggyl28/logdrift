package window

import (
	"testing"
	"time"
)

func TestNew_ValidDuration(t *testing.T) {
	w, err := New(5 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestNew_ZeroDurationReturnsError(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestNew_NegativeDurationReturnsError(t *testing.T) {
	_, err := New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestAdd_And_Total(t *testing.T) {
	w, _ := New(5 * time.Second)
	w.Add(3)
	w.Add(7)
	if got := w.Total(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestTotal_InitiallyZero(t *testing.T) {
	w, _ := New(1 * time.Second)
	if got := w.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	w, _ := New(5 * time.Second)
	w.Add(10)
	w.Reset()
	if got := w.Total(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestBuckets_ReturnsCopy(t *testing.T) {
	w, _ := New(5 * time.Second)
	w.Add(1)
	w.Add(2)
	b := w.Buckets()
	if len(b) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(b))
	}
	// Mutating the copy must not affect internal state.
	b[0].Count = 999
	if w.Total() != 3 {
		t.Fatal("internal state was mutated via Buckets slice")
	}
}

func TestEviction_RemovesExpiredEntries(t *testing.T) {
	w, _ := New(50 * time.Millisecond)
	w.Add(5)
	time.Sleep(80 * time.Millisecond)
	w.Add(2)
	if got := w.Total(); got != 2 {
		t.Fatalf("expected 2 after eviction, got %d", got)
	}
}

func TestBuckets_EmptyAfterEviction(t *testing.T) {
	w, _ := New(30 * time.Millisecond)
	w.Add(4)
	time.Sleep(60 * time.Millisecond)
	if b := w.Buckets(); len(b) != 0 {
		t.Fatalf("expected empty buckets, got %d", len(b))
	}
}
