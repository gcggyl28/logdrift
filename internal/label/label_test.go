package label_test

import (
	"testing"

	"github.com/user/logdrift/internal/label"
)

func TestRegister_Valid(t *testing.T) {
	r := label.NewRegistry()
	if err := r.Register("api", "api-service", "backend"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegister_EmptyName(t *testing.T) {
	r := label.NewRegistry()
	if err := r.Register(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_DuplicateLabel(t *testing.T) {
	r := label.NewRegistry()
	_ = r.Register("api")
	if err := r.Register("api"); err == nil {
		t.Fatal("expected ErrDuplicateLabel")
	}
}

func TestRegister_DuplicateAlias(t *testing.T) {
	r := label.NewRegistry()
	_ = r.Register("api", "shared-alias")
	if err := r.Register("worker", "shared-alias"); err == nil {
		t.Fatal("expected ErrDuplicateLabel for duplicate alias")
	}
}

func TestResolve_ByCanonical(t *testing.T) {
	r := label.NewRegistry()
	_ = r.Register("api")
	got, err := r.Resolve("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "api" {
		t.Errorf("expected 'api', got %q", got)
	}
}

func TestResolve_ByAlias(t *testing.T) {
	r := label.NewRegistry()
	_ = r.Register("api", "backend")
	got, err := r.Resolve("backend")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "api" {
		t.Errorf("expected canonical 'api', got %q", got)
	}
}

func TestResolve_NotFound(t *testing.T) {
	r := label.NewRegistry()
	_, err := r.Resolve("unknown")
	if err == nil {
		t.Fatal("expected ErrNotFound")
	}
}

func TestList_ReturnsAllCanonical(t *testing.T) {
	r := label.NewRegistry()
	_ = r.Register("api", "backend")
	_ = r.Register("worker")
	list := r.List()
	if len(list) != 2 {
		t.Errorf("expected 2 labels, got %d", len(list))
	}
}

func TestList_AliasesNotIncluded(t *testing.T) {
	r := label.NewRegistry()
	_ = r.Register("api", "alias1", "alias2")
	for _, l := range r.List() {
		if l == "alias1" || l == "alias2" {
			t.Errorf("alias should not appear in List: %q", l)
		}
	}
}
