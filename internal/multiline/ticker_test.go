package multiline

import (
	"context"
	"testing"
	"time"
)

func TestNewTimeoutFlusher_ZeroTimeoutReturnsError(t *testing.T) {
	c, _ := New(Config{Mode: ModePrefix, Pattern: `^START`})
	_, err := NewTimeoutFlusher(c, make(chan string, 1))
	if err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestNewTimeoutFlusher_ValidTimeout(t *testing.T) {
	c, _ := New(Config{Mode: ModePrefix, Pattern: `^START`, Timeout: 50 * time.Millisecond})
	_, err := NewTimeoutFlusher(c, make(chan string, 1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTimeoutFlusher_FlushesOnTimeout(t *testing.T) {
	c, _ := New(Config{Mode: ModePrefix, Pattern: `^START`, Timeout: 40 * time.Millisecond})
	out := make(chan string, 4)
	tf, _ := NewTimeoutFlusher(c, out)

	c.Push("START foo")
	c.Push("  continuation")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go tf.Run(ctx)

	select {
	case ev := <-out:
		if ev == "" {
			t.Fatal("received empty event")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for flushed event")
	}
}

func TestTimeoutFlusher_FinalFlushOnCancel(t *testing.T) {
	c, _ := New(Config{Mode: ModePrefix, Pattern: `^START`, Timeout: 10 * time.Second})
	out := make(chan string, 4)
	tf, _ := NewTimeoutFlusher(c, out)

	c.Push("START pending")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		tf.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not return after cancel")
	}

	select {
	case ev := <-out:
		if ev == "" {
			t.Fatal("expected non-empty final flush")
		}
	default:
		t.Fatal("expected final flush event in channel")
	}
}

func TestTimeoutFlusher_NoFlush_WhenBufferEmpty(t *testing.T) {
	c, _ := New(Config{Mode: ModePrefix, Pattern: `^START`, Timeout: 30 * time.Millisecond})
	out := make(chan string, 4)
	tf, _ := NewTimeoutFlusher(c, out)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	tf.Run(ctx)

	if len(out) != 0 {
		t.Errorf("expected no events for empty combiner, got %d", len(out))
	}
}
