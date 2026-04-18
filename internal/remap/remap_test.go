package remap

import (
	"testing"
)

func TestNew_EmptyRulesReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyFieldReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", Mapping: map[string]string{"a": "b"}}})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptyMappingReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "level", Mapping: map[string]string{}}})
	if err == nil {
		t.Fatal("expected error for empty mapping")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := New([]Rule{{Field: "level", Mapping: map[string]string{"warn": "warning"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_RemapsMatchingValue(t *testing.T) {
	r, _ := New([]Rule{{Field: "level", Mapping: map[string]string{"warn": "warning", "err": "error"}}})
	out := r.Apply(map[string]string{"level": "warn", "msg": "hello"})
	if out["level"] != "warning" {
		t.Errorf("expected 'warning', got %q", out["level"])
	}
	if out["msg"] != "hello" {
		t.Errorf("expected msg unchanged, got %q", out["msg"])
	}
}

func TestApply_UnknownValueUnchanged(t *testing.T) {
	r, _ := New([]Rule{{Field: "level", Mapping: map[string]string{"warn": "warning"}}})
	out := r.Apply(map[string]string{"level": "info"})
	if out["level"] != "info" {
		t.Errorf("expected 'info' unchanged, got %q", out["level"])
	}
}

func TestApply_MissingFieldUnchanged(t *testing.T) {
	r, _ := New([]Rule{{Field: "level", Mapping: map[string]string{"warn": "warning"}}})
	out := r.Apply(map[string]string{"msg": "hello"})
	if _, ok := out["level"]; ok {
		t.Error("expected level field to be absent")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	r, _ := New([]Rule{
		{Field: "level", Mapping: map[string]string{"err": "error"}},
		{Field: "env", Mapping: map[string]string{"prod": "production"}},
	})
	out := r.Apply(map[string]string{"level": "err", "env": "prod"})
	if out["level"] != "error" {
		t.Errorf("expected 'error', got %q", out["level"])
	}
	if out["env"] != "production" {
		t.Errorf("expected 'production', got %q", out["env"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	r, _ := New([]Rule{{Field: "level", Mapping: map[string]string{"warn": "warning"}}})
	input := map[string]string{"level": "warn"}
	r.Apply(input)
	if input["level"] != "warn" {
		t.Error("Apply must not mutate the input map")
	}
}
