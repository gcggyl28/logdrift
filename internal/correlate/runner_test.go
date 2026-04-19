package correlate_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/logdrift/internal/correlate"
)

func TestNewRunner_NilCorrelatorReturnsError(t *testing.T) {
	src := make(chan correlate.Entry)
	out := make(chan correlate.Group)
	_, err := correlate.NewRunner(nil, src, out, time.Second)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewRunner_NilSrcReturnsError(t *testing.T) {
	c, _ := correlate.New("request_id", time.Second)
	out := make(chan correlate.Group)
	_, err := correlate.NewRunner(c, nil, out, time.Second)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewRunner_ZeroEvictIntervalReturnsError(t *testing.T) {
	c, _ := correlate.New("request_id", time.Second)
	src := make(chan correlate.Entry)
	out := make(chan correlate.Group)
	_, err := correlate.NewRunner(c, src, out, 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRunner_EmitsGroup(t *testing.T) {
	c, _ := correlate.New("request_id", time.Second)
	src := make(chan correlate.Entry, 2)
	out := make(chan correlate.Group, 2)
	r, err := correlate.NewRunner(c, src, out, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	src <- entry("api", "abc")
	src <- entry("worker", "abc")

	var last correlate.Group
	for i := 0; i < 2; i++ {
		select {
		case g := <-out:
			last = g
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for group")
		}
	}
	if len(last.Entries) != 2 {
		t.Fatalf("expected 2 entries in last group, got %d", len(last.Entries))
	}
}

func TestRunner_CancelStops(t *testing.T) {
	c, _ := correlate.New("request_id", time.Second)
	src := make(chan correlate.Entry)
	out := make(chan correlate.Group)
	r, _ := correlate.NewRunner(c, src, out, 100*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { r.Run(ctx); close(done) }()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runner did not stop after cancel")
	}
}
