package mask_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/mask"
)

func TestNew_NoFieldsReturnsError(t *testing.T) {
	_, err := mask.New(nil)
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_EmptyFieldNameReturnsError(t *testing.T) {
	_, err := mask.New([]string{"valid", ""})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_EmptyPlaceholderReturnsError(t *testing.T) {
	_, err := mask.New([]string{"token"}, mask.WithPlaceholder(""))
	if err == nil {
		t.Fatal("expected error for empty placeholder")
	}
}

func TestNew_Valid(t *testing.T) {
	m, err := mask.New([]string{"password", "token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil masker")
	}
}

func TestApply_MasksSensitiveFields(t *testing.T) {
	m, _ := mask.New([]string{"password"}, mask.WithPlaceholder("[REDACTED]"))
	entry := map[string]string{"user": "alice", "password": "s3cr3t"}
	out := m.Apply(entry)
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out["password"])
	}
	if out["user"] != "alice" {
		t.Errorf("expected alice, got %q", out["user"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	m, _ := mask.New([]string{"secret"})
	entry := map[string]string{"secret": "abc123", "msg": "hello"}
	_ = m.Apply(entry)
	if entry["secret"] != "abc123" {
		t.Error("original entry was mutated")
	}
}

func TestApply_MissingFieldIgnored(t *testing.T) {
	m, _ := mask.New([]string{"token"})
	entry := map[string]string{"msg": "hello"}
	out := m.Apply(entry)
	if _, ok := out["token"]; ok {
		t.Error("unexpected token key in output")
	}
	if out["msg"] != "hello" {
		t.Errorf("expected hello, got %q", out["msg"])
	}
}

func TestFields_ReturnsRegisteredFields(t *testing.T) {
	m, _ := mask.New([]string{"api_key", "password"})
	fields := m.Fields()
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}
}
