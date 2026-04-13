package snapshot_test

import (
	"sort"
	"testing"

	"github.com/yourorg/logdrift/internal/snapshot"
)

func TestNew_DefaultWindow(t *testing.T) {
	s := snapshot.New(0)
	if s == nil {
		t.Fatal("expected non-nil Snapshot")
	}
}

func TestPush_And_Lines(t *testing.T) {
	s := snapshot.New(3)
	s.Push("svc", "line1")
	s.Push("svc", "line2")
	s.Push("svc", "line3")

	lines := s.Lines("svc")
	if len(lines) != 3 {
		t.Fatalf("want 3 lines, got %d", len(lines))
	}
}

func TestPush_Eviction(t *testing.T) {
	s := snapshot.New(2)
	s.Push("svc", "old")
	s.Push("svc", "keep1")
	s.Push("svc", "keep2")

	lines := s.Lines("svc")
	if lines[0] != "keep1" || lines[1] != "keep2" {
		t.Fatalf("unexpected lines after eviction: %v", lines)
	}
}

func TestLines_UnknownService(t *testing.T) {
	s := snapshot.New(10)
	lines := s.Lines("missing")
	if lines == nil {
		t.Fatal("expected non-nil slice for unknown service")
	}
	if len(lines) != 0 {
		t.Fatalf("expected empty slice, got %v", lines)
	}
}

func TestServices(t *testing.T) {
	s := snapshot.New(10)
	s.Push("alpha", "a")
	s.Push("beta", "b")

	svcs := s.Services()
	sort.Strings(svcs)
	if len(svcs) != 2 || svcs[0] != "alpha" || svcs[1] != "beta" {
		t.Fatalf("unexpected services: %v", svcs)
	}
}

func TestReset(t *testing.T) {
	s := snapshot.New(10)
	s.Push("svc", "line")
	s.Reset("svc")

	lines := s.Lines("svc")
	if len(lines) != 0 {
		t.Fatalf("expected empty after reset, got %v", lines)
	}
}

func TestLines_ReturnsCopy(t *testing.T) {
	s := snapshot.New(10)
	s.Push("svc", "original")

	lines := s.Lines("svc")
	lines[0] = "mutated"

	fresh := s.Lines("svc")
	if fresh[0] != "original" {
		t.Fatal("Lines should return a copy, not a reference")
	}
}
