package offset

import (
	"errors"
	"testing"
)

func TestNew_EmptyTracker(t *testing.T) {
	tr := New()
	if len(tr.Services()) != 0 {
		t.Fatal("expected no services")
	}
}

func TestSet_And_Get(t *testing.T) {
	tr := New()
	if err := tr.Set("svc-a", 42); err != nil {
		t.Fatalf("Set: %v", err)
	}
	v, err := tr.Get("svc-a")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	tr := New()
	if err := tr.Set("", 1); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSet_NegativeOffsetReturnsError(t *testing.T) {
	tr := New()
	if err := tr.Set("svc", -1); err == nil {
		t.Fatal("expected error for negative offset")
	}
}

func TestGet_UnknownServiceReturnsError(t *testing.T) {
	tr := New()
	_, err := tr.Get("missing")
	if !errors.Is(err, ErrUnknownService) {
		t.Fatalf("expected ErrUnknownService, got %v", err)
	}
}

func TestAdvance_IncrementsOffset(t *testing.T) {
	tr := New()
	_ = tr.Set("svc", 100)
	if err := tr.Advance("svc", 50); err != nil {
		t.Fatalf("Advance: %v", err)
	}
	v, _ := tr.Get("svc")
	if v != 150 {
		t.Fatalf("expected 150, got %d", v)
	}
}

func TestAdvance_UnknownServiceReturnsError(t *testing.T) {
	tr := New()
	if err := tr.Advance("ghost", 10); !errors.Is(err, ErrUnknownService) {
		t.Fatalf("expected ErrUnknownService, got %v", err)
	}
}

func TestAdvance_NegativeDeltaReturnsError(t *testing.T) {
	tr := New()
	_ = tr.Set("svc", 10)
	if err := tr.Advance("svc", -5); err == nil {
		t.Fatal("expected error for negative delta")
	}
}

func TestServices_ReturnsAll(t *testing.T) {
	tr := New()
	_ = tr.Set("a", 0)
	_ = tr.Set("b", 0)
	svcs := tr.Services()
	if len(svcs) != 2 {
		t.Fatalf("expected 2 services, got %d", len(svcs))
	}
}
