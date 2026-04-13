package buffer

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewFlusher_InvalidInterval(t *testing.T) {
	b, _ := New(4)
	_, err := NewFlusher(b, 0, func(_ []Entry) {})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestFlusher_CallsFnOnTick(t *testing.T) {
	b, _ := New(8)
	b.Push(Entry{Service: "svc", Line: "hello"})
	b.Push(Entry{Service: "svc", Line: "world"})

	var mu sync.Mutex
	var received []Entry

	f, err := NewFlusher(b, 20*time.Millisecond, func(entries []Entry) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, entries...)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	f.Run(ctx)

	mu.Lock()
	defer mu.Unlock()
	if len(received) < 2 {
		t.Fatalf("expected at least 2 flushed entries, got %d", len(received))
	}
}

func TestFlusher_ResetAfterFlush(t *testing.T) {
	b, _ := New(4)
	b.Push(Entry{Service: "svc", Line: "line1"})

	flushCount := 0
	f, _ := NewFlusher(b, 20*time.Millisecond, func(entries []Entry) {
		flushCount++
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	f.Run(ctx)

	if b.Len() != 0 {
		t.Fatalf("buffer should be reset after flush, len=%d", b.Len())
	}
}

func TestFlusher_FinalFlushOnCancel(t *testing.T) {
	b, _ := New(4)
	b.Push(Entry{Service: "svc", Line: "last"})

	var flushed []Entry
	f, _ := NewFlusher(b, 10*time.Second, func(entries []Entry) {
		flushed = append(flushed, entries...)
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()
	f.Run(ctx)

	if len(flushed) != 1 || flushed[0].Line != "last" {
		t.Fatalf("expected final flush with 1 entry, got %v", flushed)
	}
}
