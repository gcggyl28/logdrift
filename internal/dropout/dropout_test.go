package dropout

import (
	"testing"
)

func TestNew_ZeroDisablesDropping(t *testing.T) {
	d, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 1000; i++ {
		if !d.Allow() {
			t.Fatal("expected all entries allowed when rate=0")
		}
	}
}

func TestNew_RateOne_DropsAll(t *testing.T) {
	d, err := New(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 200; i++ {
		if d.Allow() {
			t.Fatal("expected all entries dropped when rate=1")
		}
	}
}

func TestNew_NegativeRateReturnsError(t *testing.T) {
	_, err := New(-0.1)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNew_RateAboveOneReturnsError(t *testing.T) {
	_, err := New(1.1)
	if err == nil {
		t.Fatal("expected error for rate > 1")
	}
}

func TestStats_TracksCorrectly(t *testing.T) {
	d, _ := New(1.0)
	const n = 50
	for i := 0; i < n; i++ {
		d.Allow()
	}
	total, dropped := d.Stats()
	if total != n {
		t.Fatalf("expected total=%d, got %d", n, total)
	}
	if dropped != n {
		t.Fatalf("expected dropped=%d, got %d", n, dropped)
	}
}

func TestAllow_PartialRate_ApproximateDrop(t *testing.T) {
	d, _ := New(0.5)
	const n = 10_000
	allowed := 0
	for i := 0; i < n; i++ {
		if d.Allow() {
			allowed++
		}
	}
	ratio := float64(allowed) / n
	if ratio < 0.4 || ratio > 0.6 {
		t.Fatalf("expected ~50%% allowed, got %.2f", ratio)
	}
}

func TestRate_ReturnsConfigured(t *testing.T) {
	d, _ := New(0.25)
	if d.Rate() != 0.25 {
		t.Fatalf("expected 0.25, got %v", d.Rate())
	}
}
