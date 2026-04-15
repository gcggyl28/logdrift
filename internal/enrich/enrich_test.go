package enrich

import (
	"testing"
	"time"
)

func baseEntry() Entry {
	return Entry{
		Service:   "api",
		Line:      "GET /health 200",
		Timestamp: time.Now(),
		Fields:    map[string]string{"level": "info"},
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := New(map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyFields(t *testing.T) {
	_, err := New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty fields map")
	}
}

func TestNew_EmptyKey(t *testing.T) {
	_, err := New(map[string]string{"": "value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNew_EmptyValue(t *testing.T) {
	_, err := New(map[string]string{"env": ""})
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestApply_AddsFields(t *testing.T) {
	e, _ := New(map[string]string{"env": "staging", "region": "us-east-1"})
	out := e.Apply(baseEntry())
	if out.Fields["env"] != "staging" {
		t.Errorf("expected env=staging, got %q", out.Fields["env"])
	}
	if out.Fields["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", out.Fields["region"])
	}
}

func TestApply_EntryFieldsWinOnConflict(t *testing.T) {
	e, _ := New(map[string]string{"level": "debug"})
	out := e.Apply(baseEntry()) // baseEntry has level=info
	if out.Fields["level"] != "info" {
		t.Errorf("expected entry field to win, got %q", out.Fields["level"])
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	e, _ := New(map[string]string{"env": "prod"})
	orig := baseEntry()
	e.Apply(orig)
	if _, ok := orig.Fields["env"]; ok {
		t.Error("original entry should not be mutated")
	}
}

func TestApplyAll_ReturnsCorrectCount(t *testing.T) {
	e, _ := New(map[string]string{"env": "prod"})
	entries := []Entry{baseEntry(), baseEntry(), baseEntry()}
	out := e.ApplyAll(entries)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	for i, en := range out {
		if en.Fields["env"] != "prod" {
			t.Errorf("entry %d missing env field", i)
		}
	}
}
