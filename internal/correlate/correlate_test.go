package correlate_test

import (
	"testing"
	"time"

	"github.com/user/logdrift/internal/correlate"
)

func entry(svc, reqID string) correlate.Entry {
	return correlate.Entry{
		Service: svc,
		Line:    "msg",
		Fields:  map[string]string{"request_id": reqID},
		At:      time.Now(),
	}
}

func TestNew_EmptyFieldReturnsError(t *testing.T) {
	_, err := correlate.New("", time.Second)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNew_ZeroTTLReturnsError(t *testing.T) {
	_, err := correlate.New("request_id", 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNew_Valid(t *testing.T) {
	c, err := correlate.New("request_id", time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil correlator")
	}
}

func TestAdd_GroupsByFieldValue(t *testing.T) {
	c, _ := correlate.New("request_id", time.Second)
	c.Add(entry("api", "abc"))
	g := c.Add(entry("worker", "abc"))
	if g.Value != "abc" {
		t.Fatalf("expected value abc, got %s", g.Value)
	}
	if len(g.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(g.Entries))
	}
}

func TestAdd_SeparateGroups(t *testing.T) {
	c, _ := correlate.New("request_id", time.Second)
	c.Add(entry("api", "abc"))
	g := c.Add(entry("api", "xyz"))
	if g.Value != "xyz" || len(g.Entries) != 1 {
		t.Fatal("expected isolated group for xyz")
	}
	if c.Len() != 2 {
		t.Fatalf("expected 2 groups, got %d", c.Len())
	}
}

func TestEvict_RemovesStaleGroups(t *testing.T) {
	c, _ := correlate.New("request_id", 10*time.Millisecond)
	c.Add(entry("api", "stale"))
	time.Sleep(30 * time.Millisecond)
	c.Evict()
	if c.Len() != 0 {
		t.Fatalf("expected 0 groups after eviction, got %d", c.Len())
	}
}

func TestEvict_KeepsFreshGroups(t *testing.T) {
	c, _ := correlate.New("request_id", time.Second)
	c.Add(entry("api", "fresh"))
	c.Evict()
	if c.Len() != 1 {
		t.Fatalf("expected 1 group, got %d", c.Len())
	}
}
