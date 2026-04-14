package sampler

import (
	"testing"
)

func TestNew_ValidRandom(t *testing.T) {
	s, err := New(Config{Mode: ModeRandom, Rate: 0.5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestNew_ValidEveryN(t *testing.T) {
	s, err := New(Config{Mode: ModeEveryN, N: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(Config{Mode: ModeRandom, Rate: 1.5})
	if err == nil {
		t.Fatal("expected error for rate > 1.0")
	}
}

func TestNew_InvalidN(t *testing.T) {
	_, err := New(Config{Mode: ModeEveryN, N: 0})
	if err == nil {
		t.Fatal("expected error for N < 1")
	}
}

func TestNew_UnknownMode(t *testing.T) {
	_, err := New(Config{Mode: "bogus"})
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestAllow_EveryN_ExactInterval(t *testing.T) {
	s, _ := New(Config{Mode: ModeEveryN, N: 3})
	results := make([]bool, 9)
	for i := range results {
		results[i] = s.Allow()
	}
	// Expect true at indices 2, 5, 8 (every 3rd call)
	for i, got := range results {
		want := (i+1)%3 == 0
		if got != want {
			t.Errorf("Allow() call %d: got %v, want %v", i+1, got, want)
		}
	}
}

func TestAllow_Random_AlwaysAllow(t *testing.T) {
	s, _ := New(Config{Mode: ModeRandom, Rate: 1.0})
	for i := 0; i < 100; i++ {
		if !s.Allow() {
			t.Fatalf("rate=1.0 should always allow, failed on call %d", i)
		}
	}
}

func TestAllow_Random_NeverAllow(t *testing.T) {
	s, _ := New(Config{Mode: ModeRandom, Rate: 0.0})
	for i := 0; i < 100; i++ {
		if s.Allow() {
			t.Fatalf("rate=0.0 should never allow, passed on call %d", i)
		}
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	s, _ := New(Config{Mode: ModeEveryN, N: 2})
	s.Allow() // counter = 1
	s.Reset() // counter = 0
	// After reset, first Allow() should be false (counter goes to 1, not 2)
	if s.Allow() {
		t.Error("expected false immediately after reset with N=2")
	}
}
