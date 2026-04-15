package coalesce

import (
	"testing"
	"time"
)

func TestNew_ZeroWindowReturnsError(t *testing.T) {
	_, _, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_NegativeWindowReturnsError(t *testing.T) {
	_, _, err := New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNew_PositiveWindowSucceeds(t *testing.T) {
	c, ch, err := New(50 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil || ch == nil {
		t.Fatal("expected non-nil coalescer and channel")
	}
	c.Close()
}

func TestFlush_EmitsGroup(t *testing.T) {
	c, ch, _ := New(500 * time.Millisecond)
	c.Push(Entry{Service: "svc-a", Line: "hello"})
	c.Push(Entry{Service: "svc-b", Line: "world"})
	c.Flush()

	select {
	case g := <-ch:
		if len(g) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(g))
		}
		if g[0].Service != "svc-a" || g[1].Service != "svc-b" {
			t.Errorf("unexpected services: %v", g)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for group")
	}
	c.Close()
}

func TestIdleWindowFlushes(t *testing.T) {
	c, ch, _ := New(40 * time.Millisecond)
	c.Push(Entry{Service: "svc-a", Line: "line1"})

	select {
	case g := <-ch:
		if len(g) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(g))
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for idle flush")
	}
	c.Close()
}

func TestClose_FlushesPendingEntries(t *testing.T) {
	c, ch, _ := New(500 * time.Millisecond)
	c.Push(Entry{Service: "svc-x", Line: "pending"})
	c.Close()

	var groups []Group
	for g := range ch {
		groups = append(groups, g)
	}
	if len(groups) != 1 || len(groups[0]) != 1 {
		t.Fatalf("expected 1 group with 1 entry, got %v", groups)
	}
}

func TestPush_SetsTimestamp(t *testing.T) {
	c, ch, _ := New(40 * time.Millisecond)
	before := time.Now()
	c.Push(Entry{Service: "svc", Line: "ts-check"})

	select {
	case g := <-ch:
		if g[0].At.Before(before) {
			t.Errorf("entry timestamp should be >= push time")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out")
	}
	c.Close()
}

func TestClose_IdempotentDoesNotPanic(t *testing.T) {
	c, _, _ := New(50 * time.Millisecond)
	c.Close()
	c.Close() // second close must not panic
}
