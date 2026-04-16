package fanout_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/logdrift/internal/fanout"
	"github.com/user/logdrift/internal/tail"
)

func makeSrc(entries []tail.Entry) chan tail.Entry {
	ch := make(chan tail.Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func TestNew_InvalidBufSz(t *testing.T) {
	ch := make(chan tail.Entry)
	_, err := fanout.New(ch, 0)
	if err == nil {
		t.Fatal("expected error for bufSz=0")
	}
}

func TestNew_Valid(t *testing.T) {
	ch := make(chan tail.Entry)
	_, err := fanout.New(ch, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBroadcast_AllSubscribersReceive(t *testing.T) {
	entries := []tail.Entry{
		{Service: "svc", Line: "line1"},
		{Service: "svc", Line: "line2"},
	}
	src := makeSrc(entries)
	b, _ := fanout.New(src, 8)

	s1 := b.Subscribe()
	s2 := b.Subscribe()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go b.Run(ctx)

	for i, want := range []string{"line1", "line2"} {
		for _, sub := range []<-chan tail.Entry{s1, s2} {
			select {
			case got := <-sub:
				if got.Line != want {
					t.Errorf("sub entry %d: got %q want %q", i, got.Line, want)
				}
			case <-ctx.Done():
				t.Fatalf("timeout waiting for entry %d", i)
			}
		}
	}
}

func TestBroadcast_CancelStops(t *testing.T) {
	src := make(chan tail.Entry)
	b, _ := fanout.New(src, 4)
	_ = b.Subscribe()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		b.Run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after cancel")
	}
}

func TestBroadcast_ClosedSubsOnDone(t *testing.T) {
	src := makeSrc(nil)
	b, _ := fanout.New(src, 4)
	sub := b.Subscribe()

	ctx := context.Background()
	b.Run(ctx)

	_, open := <-sub
	if open {
		t.Fatal("subscriber channel should be closed after src closes")
	}
}
