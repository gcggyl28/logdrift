package label_test

import (
	"testing"

	"github.com/user/logdrift/internal/label"
)

func newTestRegistry(t *testing.T) *label.Registry {
	t.Helper()
	r := label.NewRegistry()
	if err := r.Register("api", "backend"); err != nil {
		t.Fatalf("register api: %v", err)
	}
	if err := r.Register("worker"); err != nil {
		t.Fatalf("register worker: %v", err)
	}
	return r
}

func TestTag_CanonicalName(t *testing.T) {
	tgr := label.NewTagger(newTestRegistry(t))
	e, err := tgr.Tag("api", "hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Service != "api" || e.Line != "hello world" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestTag_AliasResolvesToCanonical(t *testing.T) {
	tgr := label.NewTagger(newTestRegistry(t))
	e, err := tgr.Tag("backend", "log line")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Service != "api" {
		t.Errorf("expected canonical 'api', got %q", e.Service)
	}
}

func TestTag_UnknownServiceReturnsError(t *testing.T) {
	tgr := label.NewTagger(newTestRegistry(t))
	_, err := tgr.Tag("ghost", "line")
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestTagAll_MixedValid(t *testing.T) {
	tgr := label.NewTagger(newTestRegistry(t))
	pairs := [][2]string{
		{"api", "line1"},
		{"unknown", "line2"},
		{"worker", "line3"},
	}
	entries, errs := tgr.TagAll(pairs)
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
}

func TestTagAll_AllValid(t *testing.T) {
	tgr := label.NewTagger(newTestRegistry(t))
	pairs := [][2]string{
		{"api", "a"},
		{"worker", "b"},
	}
	entries, errs := tgr.TagAll(pairs)
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d", len(errs))
	}
}
