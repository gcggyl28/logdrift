package evict

import (
	"testing"
	"time"
)

func TestNew_PositiveTTL(t *testing.T) {
	e, err := New(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil Evictor")
	}
}

func TestNew_ZeroTTLReturnsError(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestNew_NegativeTTLReturnsError(t *testing.T) {
	_, err := New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative TTL")
	}
}

func TestAdd_And_Has(t *testing.T) {
	e, _ := New(time.Second)
	e.Add("key1")
	if !e.Has("key1") {
		t.Error("expected Has to return true for added key")
	}
}

func TestHas_UnknownKeyReturnsFalse(t *testing.T) {
	e, _ := New(time.Second)
	if e.Has("missing") {
		t.Error("expected Has to return false for unknown key")
	}
}

func TestHas_ExpiredKeyReturnsFalse(t *testing.T) {
	e, _ := New(10 * time.Millisecond)
	e.Add("expiring")
	time.Sleep(30 * time.Millisecond)
	if e.Has("expiring") {
		t.Error("expected Has to return false for expired key")
	}
}

func TestEvict_RemovesExpiredEntries(t *testing.T) {
	e, _ := New(10 * time.Millisecond)
	e.Add("a")
	e.Add("b")
	time.Sleep(30 * time.Millisecond)
	e.Add("c") // still live
	removed := e.Evict()
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if !e.Has("c") {
		t.Error("live entry 'c' should still be present")
	}
}

func TestEvict_NoExpiredEntries(t *testing.T) {
	e, _ := New(time.Second)
	e.Add("x")
	e.Add("y")
	removed := e.Evict()
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
}

func TestLen_CountsLiveEntries(t *testing.T) {
	e, _ := New(50 * time.Millisecond)
	e.Add("a")
	e.Add("b")
	e.Add("c")
	if l := e.Len(); l != 3 {
		t.Errorf("expected Len 3, got %d", l)
	}
	time.Sleep(80 * time.Millisecond)
	if l := e.Len(); l != 0 {
		t.Errorf("expected Len 0 after expiry, got %d", l)
	}
}

func TestAdd_RefreshesExpiry(t *testing.T) {
	e, _ := New(50 * time.Millisecond)
	e.Add("k")
	time.Sleep(30 * time.Millisecond)
	e.Add("k") // refresh
	time.Sleep(30 * time.Millisecond)
	if !e.Has("k") {
		t.Error("expected key to still be live after refresh")
	}
}
