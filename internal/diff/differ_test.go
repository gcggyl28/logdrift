package diff

import (
	"strings"
	"testing"
)

func TestNew_ValidModes(t *testing.T) {
	modes := []Mode{ModeUnified, ModeSideBySide, ModeTimestamp}
	for _, m := range modes {
		d, err := New(m)
		if err != nil {
			t.Errorf("expected no error for mode %q, got %v", m, err)
		}
		if d == nil {
			t.Errorf("expected non-nil Differ for mode %q", m)
		}
	}
}

func TestNew_InvalidMode(t *testing.T) {
	_, err := New("bogus")
	if err == nil {
		t.Fatal("expected error for invalid mode, got nil")
	}
}

func TestCompare_NoDrift(t *testing.T) {
	d, _ := New(ModeUnified)
	a := Entry{Service: "svc-a", Line: "hello world"}
	b := Entry{Service: "svc-b", Line: "hello world"}
	res := d.Compare(a, b)
	if res.Drifted {
		t.Error("expected no drift for identical lines")
	}
	if res.Delta != "" {
		t.Errorf("expected empty delta, got %q", res.Delta)
	}
}

func TestCompare_UnifiedDrift(t *testing.T) {
	d, _ := New(ModeUnified)
	a := Entry{Service: "svc-a", Line: "line A"}
	b := Entry{Service: "svc-b", Line: "line B"}
	res := d.Compare(a, b)
	if !res.Drifted {
		t.Fatal("expected drift for different lines")
	}
	if !strings.Contains(res.Delta, "--- svc-a") {
		t.Errorf("expected unified header in delta, got: %s", res.Delta)
	}
	if !strings.Contains(res.Delta, "+++ svc-b") {
		t.Errorf("expected unified header in delta, got: %s", res.Delta)
	}
}

func TestCompare_SideBySide(t *testing.T) {
	d, _ := New(ModeSideBySide)
	a := Entry{Service: "alpha", Line: "foo"}
	b := Entry{Service: "beta", Line: "bar"}
	res := d.Compare(a, b)
	if !res.Drifted {
		t.Fatal("expected drift")
	}
	if !strings.Contains(res.Delta, "|") {
		t.Errorf("expected side-by-side separator in delta, got: %s", res.Delta)
	}
}

func TestCompare_Timestamp(t *testing.T) {
	d, _ := New(ModeTimestamp)
	a := Entry{Service: "x", Line: "2024-01-01T00:00:00 msg"}
	b := Entry{Service: "y", Line: "2024-01-01T00:00:05 msg"}
	res := d.Compare(a, b)
	if !res.Drifted {
		t.Fatal("expected drift")
	}
	if !strings.Contains(res.Delta, "timestamp drift") {
		t.Errorf("expected timestamp message, got: %s", res.Delta)
	}
}
