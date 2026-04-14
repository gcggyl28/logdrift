package aggregate

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNew_ValidWindow(t *testing.T) {
	a, err := New(time.Minute)
	if err != nil || a == nil {
		t.Fatalf("expected valid aggregator, got err=%v", err)
	}
}

func TestNew_ZeroWindow(t *testing.T) {
	a, err := New(0)
	if err != nil || a == nil {
		t.Fatalf("expected valid aggregator with zero window, got err=%v", err)
	}
}

func TestNew_NegativeWindow(t *testing.T) {
	_, err := New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestRecordLine_IncrementsCounter(t *testing.T) {
	a, _ := New(0)
	a.RecordLine("svc-a", epoch)
	a.RecordLine("svc-a", epoch.Add(time.Second))
	snap := a.Snapshot()
	if len(snap) != 1 || snap[0].Lines != 2 {
		t.Fatalf("expected 2 lines, got %+v", snap)
	}
}

func TestRecordDrift_IncrementsCounter(t *testing.T) {
	a, _ := New(0)
	a.RecordDrift("svc-b", epoch)
	a.RecordDrift("svc-b", epoch)
	snap := a.Snapshot()
	if len(snap) != 1 || snap[0].Drifts != 2 {
		t.Fatalf("expected 2 drifts, got %+v", snap)
	}
}

func TestSnapshot_IndependentKeys(t *testing.T) {
	a, _ := New(0)
	a.RecordLine("svc-a", epoch)
	a.RecordLine("svc-b", epoch)
	snap := a.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(snap))
	}
}

func TestSnapshot_IsolatesInternalState(t *testing.T) {
	a, _ := New(0)
	a.RecordLine("svc-a", epoch)
	snap := a.Snapshot()
	snap[0].Lines = 999
	snap2 := a.Snapshot()
	if snap2[0].Lines == 999 {
		t.Fatal("snapshot mutation leaked into aggregator")
	}
}

func TestEviction_RemovesStaleKeys(t *testing.T) {
	a, _ := New(time.Minute)
	a.RecordLine("old", epoch)
	// Trigger eviction with a time far in the future
	a.RecordLine("new", epoch.Add(2*time.Minute))
	snap := a.Snapshot()
	for _, s := range snap {
		if s.Key == "old" {
			t.Fatal("expected old key to be evicted")
		}
	}
}

func TestReset_ClearsData(t *testing.T) {
	a, _ := New(0)
	a.RecordLine("svc-a", epoch)
	a.Reset()
	if len(a.Snapshot()) != 0 {
		t.Fatal("expected empty snapshot after reset")
	}
}

func TestFirstSeen_LastSeen(t *testing.T) {
	a, _ := New(0)
	t1 := epoch
	t2 := epoch.Add(5 * time.Second)
	a.RecordLine("svc", t1)
	a.RecordLine("svc", t2)
	snn	if snap[0].FirstSeen != t1 || snap[0].LastSeen != t2 {
		t.Fatalf("unexpected times: %+v", snap[0])
	}
}
