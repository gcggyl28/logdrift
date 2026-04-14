package redact

import (
	"strings"
	"testing"
)

func TestNew_NoPatterns(t *testing.T) {
	r, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Enabled() {
		t.Fatal("expected Enabled() == false with no patterns")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]string{"["})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New([]string{""})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_ValidPatterns(t *testing.T) {
	r, err := New([]string{`password=\S+`, `token=[A-Za-z0-9]+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Enabled() {
		t.Fatal("expected Enabled() == true")
	}
}

func TestApply_NoPatterns_Passthrough(t *testing.T) {
	r, _ := New(nil)
	line := "user logged in with password=secret"
	if got := r.Apply(line); got != line {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestApply_MasksMatch(t *testing.T) {
	r, err := New([]string{`password=\S+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Apply("user login password=hunter2 ok")
	if strings.Contains(got, "hunter2") {
		t.Fatalf("sensitive value not redacted: %q", got)
	}
	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("expected mask in output: %q", got)
	}
}

func TestApply_CustomMask(t *testing.T) {
	r, err := New([]string{`token=\S+`}, WithMask("***"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Apply("auth token=abc123 granted")
	if !strings.Contains(got, "***") {
		t.Fatalf("expected custom mask in output: %q", got)
	}
}

func TestApply_MultiplePatterns(t *testing.T) {
	r, err := New([]string{`password=\S+`, `ssn=\d+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Apply("password=secret ssn=123456789")
	if strings.Contains(got, "secret") || strings.Contains(got, "123456789") {
		t.Fatalf("not all sensitive values redacted: %q", got)
	}
}
