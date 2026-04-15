package linecount

import (
	"testing"
	"time"
)

func TestNew_PositiveWindow(t *testing.T) {
	c, err := New(5 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Counter")
	}
}

func TestNew_ZeroWindowReturnsError(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_NegativeWindowReturnsError(t *testing.T) {
	_, err := New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestRate_InitiallyZero(t *testing.T) {
	c, _ := New(5 * time.Second)
	if r := c.Rate("svc"); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestRecord_IncrementsRate(t *testing.T) {
	c, _ := New(5 * time.Second)
	for i := 0; i < 10; i++ {
		c.Record("svc")
	}
	r := c.Rate("svc")
	if r <= 0 {
		t.Fatalf("expected positive rate, got %f", r)
	}
}

func TestRate_IndependentServices(t *testing.T) {
	c, _ := New(5 * time.Second)
	c.Record("alpha")
	c.Record("alpha")
	c.Record("beta")

	ra := c.Rate("alpha")
	rb := c.Rate("beta")
	if ra <= rb {
		t.Fatalf("alpha rate (%f) should exceed beta rate (%f)", ra, rb)
	}
}

func TestServices_ReturnsActiveServices(t *testing.T) {
	c, _ := New(5 * time.Second)
	c.Record("svc-a")
	c.Record("svc-b")

	svcs := c.Services()
	if len(svcs) != 2 {
		t.Fatalf("expected 2 services, got %d", len(svcs))
	}
}

func TestServices_ExcludesExpired(t *testing.T) {
	c, _ := New(50 * time.Millisecond)
	c.Record("old-svc")
	time.Sleep(100 * time.Millisecond)
	svcs := c.Services()
	for _, s := range svcs {
		if s == "old-svc" {
			t.Fatal("expired service should not appear")
		}
	}
}

func TestEviction_RemovesOldEntries(t *testing.T) {
	c, _ := New(50 * time.Millisecond)
	for i := 0; i < 5; i++ {
		c.Record("svc")
	}
	time.Sleep(100 * time.Millisecond)
	r := c.Rate("svc")
	if r != 0 {
		t.Fatalf("expected 0 after eviction, got %f", r)
	}
}
