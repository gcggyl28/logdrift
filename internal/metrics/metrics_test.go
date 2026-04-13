package metrics

import (
	"testing"
	"time"
)

func TestNew_InitialState(t *testing.T) {
	tr := New()
	snap := tr.Snapshot()

	if len(snap.LinesReceived) != 0 {
		t.Fatalf("expected empty lines map, got %v", snap.LinesReceived)
	}
	if snap.DriftEvents != 0 {
		t.Fatalf("expected 0 drift events, got %d", snap.DriftEvents)
	}
	if snap.StartedAt.IsZero() {
		t.Fatal("StartedAt should not be zero")
	}
}

func TestRecordLine_IncrementsCounter(t *testing.T) {
	tr := New()
	tr.RecordLine("svc-a")
	tr.RecordLine("svc-a")
	tr.RecordLine("svc-b")

	snap := tr.Snapshot()
	if snap.LinesReceived["svc-a"] != 2 {
		t.Errorf("svc-a: want 2, got %d", snap.LinesReceived["svc-a"])
	}
	if snap.LinesReceived["svc-b"] != 1 {
		t.Errorf("svc-b: want 1, got %d", snap.LinesReceived["svc-b"])
	}
}

func TestRecordDrift_IncrementsCounter(t *testing.T) {
	tr := New()
	tr.RecordDrift()
	tr.RecordDrift()

	snap := tr.Snapshot()
	if snap.DriftEvents != 2 {
		t.Errorf("want 2 drift events, got %d", snap.DriftEvents)
	}
}

func TestSnapshot_IsolatesInternalState(t *testing.T) {
	tr := New()
	tr.RecordLine("svc-a")

	snap := tr.Snapshot()
	snap.LinesReceived["svc-a"] = 999 // mutate copy

	snap2 := tr.Snapshot()
	if snap2.LinesReceived["svc-a"] != 1 {
		t.Error("snapshot mutation leaked into tracker state")
	}
}

func TestSnapshot_CapturedAt_IsRecent(t *testing.T) {
	tr := New()
	before := time.Now()
	snap := tr.Snapshot()
	after := time.Now()

	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v not in [%v, %v]", snap.CapturedAt, before, after)
	}
}

func TestReset_ZeroesCounters(t *testing.T) {
	tr := New()
	tr.RecordLine("svc-a")
	tr.RecordDrift()

	tr.Reset()
	snap := tr.Snapshot()

	if len(snap.LinesReceived) != 0 {
		t.Errorf("expected empty lines after reset, got %v", snap.LinesReceived)
	}
	if snap.DriftEvents != 0 {
		t.Errorf("expected 0 drift events after reset, got %d", snap.DriftEvents)
	}
}
