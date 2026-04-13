package buffer

import (
	"fmt"
	"testing"
)

func TestNew_ValidCapacity(t *testing.T) {
	b, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Cap() != 10 {
		t.Fatalf("expected cap 10, got %d", b.Cap())
	}
}

func TestNew_InvalidCapacity(t *testing.T) {
	for _, c := range []int{0, -1, -100} {
		_, err := New(c)
		if err == nil {
			t.Fatalf("expected error for capacity %d", c)
		}
	}
}

func TestPush_And_Entries(t *testing.T) {
	b, _ := New(5)
	for i := 0; i < 3; i++ {
		b.Push(Entry{Service: "svc", Line: fmt.Sprintf("line%d", i)})
	}
	got := b.Entries()
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	if got[0].Line != "line0" || got[2].Line != "line2" {
		t.Fatalf("unexpected entries: %v", got)
	}
}

func TestPush_Overflow_EvictsOldest(t *testing.T) {
	b, _ := New(3)
	for i := 0; i < 5; i++ {
		b.Push(Entry{Service: "svc", Line: fmt.Sprintf("line%d", i)})
	}
	got := b.Entries()
	if len(got) != 3 {
		t.Fatalf("expected 3 entries after overflow, got %d", len(got))
	}
	if got[0].Line != "line2" {
		t.Fatalf("expected oldest to be line2, got %s", got[0].Line)
	}
	if got[2].Line != "line4" {
		t.Fatalf("expected newest to be line4, got %s", got[2].Line)
	}
}

func TestLen_Tracking(t *testing.T) {
	b, _ := New(4)
	if b.Len() != 0 {
		t.Fatal("expected initial len 0")
	}
	b.Push(Entry{Service: "a", Line: "x"})
	b.Push(Entry{Service: "b", Line: "y"})
	if b.Len() != 2 {
		t.Fatalf("expected len 2, got %d", b.Len())
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	b, _ := New(4)
	b.Push(Entry{Service: "svc", Line: "hello"})
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", b.Len())
	}
	if len(b.Entries()) != 0 {
		t.Fatal("expected empty entries after reset")
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	b, _ := New(4)
	b.Push(Entry{Service: "svc", Line: "original"})
	snap := b.Entries()
	snap[0].Line = "mutated"
	got := b.Entries()
	if got[0].Line != "original" {
		t.Fatal("Entries should return isolated snapshot")
	}
}
