package zeitgeist

import (
	"testing"
	"time"
)

func TestNew_ZeroWidthReturnsError(t *testing.T) {
	_, err := New(0, 10)
	if err == nil {
		t.Fatal("expected error for zero width")
	}
}

func TestNew_ZeroMaxSlotsReturnsError(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero maxSlots")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := New(time.Minute, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecord_And_Snapshot(t *testing.T) {
	b, _ := New(time.Minute, 10)
	now := time.Now().Truncate(time.Minute)
	b.Record("svc-a", now)
	b.Record("svc-a", now)
	b.Record("svc-b", now)
	snap := b.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap))
	}
	if snap[0].Counts["svc-a"] != 2 {
		t.Errorf("svc-a: want 2, got %d", snap[0].Counts["svc-a"])
	}
	if snap[0].Counts["svc-b"] != 1 {
		t.Errorf("svc-b: want 1, got %d", snap[0].Counts["svc-b"])
	}
}

func TestSnapshot_Isolation(t *testing.T) {
	b, _ := New(time.Minute, 10)
	now := time.Now().Truncate(time.Minute)
	b.Record("svc", now)
	snap := b.Snapshot()
	snap[0].Counts["svc"] = 999
	snap2 := b.Snapshot()
	if snap2[0].Counts["svc"] != 1 {
		t.Error("snapshot mutation leaked into internal state")
	}
}

func TestEviction_KeepsMaxSlots(t *testing.T) {
	b, _ := New(time.Minute, 3)
	base := time.Now().Truncate(time.Minute)
	for i := 0; i < 5; i++ {
		b.Record("svc", base.Add(time.Duration(i)*time.Minute))
	}
	snap := b.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 buckets after eviction, got %d", len(snap))
	}
	if !snap[0].Start.Equal(base.Add(2 * time.Minute)) {
		t.Errorf("oldest retained bucket mismatch: %v", snap[0].Start)
	}
}
