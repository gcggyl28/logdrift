package fieldmap

import (
	"testing"
)

func TestNew_EmptyRulesReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyFromReturnsError(t *testing.T) {
	_, err := New([]Rule{{From: "", To: "dest"}})
	if err == nil {
		t.Fatal("expected error for empty From")
	}
}

func TestNew_DuplicateFromReturnsError(t *testing.T) {
	_, err := New([]Rule{
		{From: "msg", To: "message"},
		{From: "msg", To: "text"},
	})
	if err == nil {
		t.Fatal("expected error for duplicate From key")
	}
}

func TestNew_ValidRules(t *testing.T) {
	m, err := New([]Rule{{From: "msg", To: "message"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Mapper")
	}
}

func TestApply_RenamesField(t *testing.T) {
	m, _ := New([]Rule{{From: "msg", To: "message"}})
	out := m.Apply(map[string]string{"msg": "hello", "level": "info"})
	if out["message"] != "hello" {
		t.Errorf("expected message=hello, got %q", out["message"])
	}
	if _, ok := out["msg"]; ok {
		t.Error("original key msg should have been removed")
	}
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %q", out["level"])
	}
}

func TestApply_DropsFieldWhenToIsEmpty(t *testing.T) {
	m, _ := New([]Rule{{From: "secret", To: ""}})
	out := m.Apply(map[string]string{"secret": "s3cr3t", "msg": "ok"})
	if _, ok := out["secret"]; ok {
		t.Error("secret should have been dropped")
	}
	if out["msg"] != "ok" {
		t.Errorf("expected msg=ok, got %q", out["msg"])
	}
}

func TestApply_PassthroughUnmatchedFields(t *testing.T) {
	m, _ := New([]Rule{{From: "ts", To: "timestamp"}})
	out := m.Apply(map[string]string{"level": "warn", "body": "x"})
	if out["level"] != "warn" || out["body"] != "x" {
		t.Error("unmatched fields should pass through unchanged")
	}
}

func TestApply_MissingFromKeyIgnored(t *testing.T) {
	m, _ := New([]Rule{{From: "nonexistent", To: "other"}})
	out := m.Apply(map[string]string{"msg": "hi"})
	if out["msg"] != "hi" {
		t.Errorf("expected msg=hi, got %q", out["msg"])
	}
	if _, ok := out["other"]; ok {
		t.Error("other should not appear when source key is absent")
	}
}

func TestRules_ReturnsCopy(t *testing.T) {
	orig := []Rule{{From: "a", To: "b"}}
	m, _ := New(orig)
	got := m.Rules()
	got[0].From = "mutated"
	if m.Rules()[0].From != "a" {
		t.Error("Rules() should return an isolated copy")
	}
}
