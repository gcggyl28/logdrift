package debounce_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/user/logdrift/internal/debounce"
)

func TestRunner_ForwardsDebounced(t *testing.T) {
	src := make(chan debounce.Entry, 4)
	d, _ := debounce.New(20 * time.Millisecond)
	r := debounce.NewRunner(d, src)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var mu sync.Mutex
	var received []debounce.Entry

	go r.Run(ctx, func(e debounce.Entry) {
		mu.Lock()
		received = append(received, e)
		mu.Unlock()
	})

	src <- debounce.Entry{Service: "x", Line: "msg"}

	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) == 0 {
		t.Fatal("expected at least one forwarded entry")
	}
	if received[0].Line != "msg" {
		t.Fatalf("expected 'msg', got %q", received[0].Line)
	}
}

func TestRunner_StopsOnContextCancel(t *testing.T) {
	src := make(chan debounce.Entry)
	d, _ := debounce.New(10 * time.Millisecond)
	r := debounce.NewRunner(d, src)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		defer close(done)
		r.Run(ctx, func(debounce.Entry) {})
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Runner did not stop after context cancel")
	}
}
