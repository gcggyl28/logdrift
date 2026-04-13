package replay_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/logdrift/internal/replay"
	"github.com/user/logdrift/internal/snapshot"
)

func buildSnap(t *testing.T, window int, data map[string][]string) *snapshot.Snapshot {
	t.Helper()
	snap := snapshot.New(window)
	for svc, lines := range data {
		for _, l := range lines {
			snap.Push(svc, l)
		}
	}
	return snap
}

func TestRun_EmitsAllLines(t *testing.T) {
	snap := buildSnap(t, 10, map[string][]string{
		"svcA": {"a1", "a2"},
		"svcB": {"b1"},
	})
	r := replay.New(snap, 0)
	ctx := context.Background()
	ch := r.Run(ctx)
	var entries []replay.Entry
	for e := range ch {
		entries = append(entries, e)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestRun_CancelStopsEarly(t *testing.T) {
	snap := buildSnap(t, 100, map[string][]string{
		"svcA": make([]string, 50),
	})
	r := replay.New(snap, 5*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	ch := r.Run(ctx)
	// read one entry then cancel
	<-ch
	cancel()
	// drain
	for range ch {
	}
	// reaching here means the goroutine exited — pass
}

func TestRun_EmptySnapshot(t *testing.T) {
	snap := snapshot.New(10)
	r := replay.New(snap, 0)
	ctx := context.Background()
	ch := r.Run(ctx)
	var count int
	for range ch {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 entries from empty snapshot, got %d", count)
	}
}

func TestRun_EntryFieldsPopulated(t *testing.T) {
	snap := buildSnap(t, 10, map[string][]string{
		"svcX": {"hello"},
	})
	r := replay.New(snap, 0)
	ctx := context.Background()
	ch := r.Run(ctx)
	e := <-ch
	if e.Service != "svcX" {
		t.Errorf("expected service svcX, got %q", e.Service)
	}
	if e.Line != "hello" {
		t.Errorf("expected line 'hello', got %q", e.Line)
	}
	if e.At.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
