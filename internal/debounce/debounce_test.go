package debounce_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/logdrift/internal/debounce"
)

func TestNew_PositiveWindow(t *testing.T) {
	d, err := debounce.New(10 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil debouncer")
	}
}

func TestNew_ZeroReturnsError(t *testing.T) {
	_, err := debounce.New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_NegativeReturnsError(t *testing.T) {
	_, err := debounce.New(-1 * time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestPush_EmitsAfterQuiet(t *testing.T) {
	d, _ := debounce.New(20 * time.Millisecond)
	d.Push(debounce.Entry{Service: "svc", Line: "hello"})

	select {
	case e := <-d.Out():
		if e.Line != "hello" {
			t.Fatalf("expected 'hello', got %q", e.Line)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for debounced entry")
	}
}

func TestPush_CollapsesBurst(t *testing.T) {
	d, _ := debounce.New(40 * time.Millisecond)
	for i := 0; i < 5; i++ {
		d.Push(debounce.Entry{Service: "svc", Line: "line"})
		time.Sleep(5 * time.Millisecond)
	}

	// Only one entry should arrive.
	count := 0
	timer := time.After(200 * time.Millisecond)
loop:
	for {
		select {
		case <-d.Out():
			count++
		case <-timer:
			break loop
		}
	}
	if count != 1 {
		t.Fatalf("expected 1 debounced entry, got %d", count)
	}
}

func TestDrain_StopsOnCancel(t *testing.T) {
	d, _ := debounce.New(10 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		d.Drain(ctx, func(debounce.Entry) {})
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Drain did not stop after context cancel")
	}
}

func TestPush_IndependentServices(t *testing.T) {
	d, _ := debounce.New(20 * time.Millisecond)
	d.Push(debounce.Entry{Service: "a", Line: "alpha"})
	d.Push(debounce.Entry{Service: "b", Line: "beta"})

	seen := map[string]bool{}
	timer := time.After(300 * time.Millisecond)
	for len(seen) < 2 {
		select {
		case e := <-d.Out():
			seen[e.Service] = true
		case <-timer:
			t.Fatalf("timed out; only saw services: %v", seen)
		}
	}
}
