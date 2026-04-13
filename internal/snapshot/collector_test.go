package snapshot_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/snapshot"
	"github.com/yourorg/logdrift/internal/tail"
)

func makeSrc(lines ...tail.Line) chan tail.Line {
	ch := make(chan tail.Line, len(lines))
	for _, l := range lines {
		ch <- l
	}
	close(ch)
	return ch
}

func TestCollector_ForwardsEntries(t *testing.T) {
	snap := snapshot.New(10)
	col := snapshot.NewCollector(snap)

	src := makeSrc(
		tail.Line{Service: "api", Text: "hello"},
		tail.Line{Service: "api", Text: "world"},
	)

	ctx := context.Background()
	go col.Run(ctx, src)

	var got []snapshot.Entry
	for e := range col.Out() {
		got = append(got, e)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Line != "hello" || got[1].Line != "world" {
		t.Fatalf("unexpected entries: %v", got)
	}
}

func TestCollector_PushesToSnapshot(t *testing.T) {
	snap := snapshot.New(10)
	col := snapshot.NewCollector(snap)

	src := makeSrc(tail.Line{Service: "db", Text: "connected"})
	ctx := context.Background()
	go col.Run(ctx, src)

	for range col.Out() {
	}

	lines := snap.Lines("db")
	if len(lines) != 1 || lines[0] != "connected" {
		t.Fatalf("unexpected snapshot lines: %v", lines)
	}
}

func TestCollector_CancelStops(t *testing.T) {
	snap := snapshot.New(10)
	col := snapshot.NewCollector(snap)

	// unbuffered, never written — blocks until cancel
	src := make(chan tail.Line)
	ctx, cancel := context.WithCancel(context.Background())

	go col.Run(ctx, src)
	cancel()

	select {
	case _, ok := <-col.Out():
		if ok {
			t.Fatal("expected Out() to be closed after cancel")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for collector to stop")
	}
}
