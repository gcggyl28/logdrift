package transform

import (
	"testing"
)

func TestNew_EmptyRulesReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil rules")
	}
	_, err = New([]Rule{})
	if err == nil {
		t.Fatal("expected error for empty rules slice")
	}
}

func TestNew_TrimPrefixMissingArgReturnsError(t *testing.T) {
	_, err := New([]Rule{{Op: OpTrimPrefix, Arg: ""}})
	if err == nil {
		t.Fatal("expected error when Arg is empty for TrimPrefix")
	}
}

func TestNew_TrimSuffixMissingArgReturnsError(t *testing.T) {
	_, err := New([]Rule{{Op: OpTrimSuffix, Arg: ""}})
	if err == nil {
		t.Fatal("expected error when Arg is empty for TrimSuffix")
	}
}

func TestNew_ValidRules(t *testing.T) {
	_, err := New([]Rule{{Op: OpTrimSpace}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_Uppercase(t *testing.T) {
	tr, _ := New([]Rule{{Op: OpUppercase}})
	if got := tr.Apply("hello"); got != "HELLO" {
		t.Fatalf("got %q, want %q", got, "HELLO")
	}
}

func TestApply_Lowercase(t *testing.T) {
	tr, _ := New([]Rule{{Op: OpLowercase}})
	if got := tr.Apply("WORLD"); got != "world" {
		t.Fatalf("got %q, want %q", got, "world")
	}
}

func TestApply_TrimSpace(t *testing.T) {
	tr, _ := New([]Rule{{Op: OpTrimSpace}})
	if got := tr.Apply("  padded  "); got != "padded" {
		t.Fatalf("got %q, want %q", got, "padded")
	}
}

func TestApply_TrimPrefix(t *testing.T) {
	tr, _ := New([]Rule{{Op: OpTrimPrefix, Arg: "INFO "}})
	if got := tr.Apply("INFO message"); got != "message" {
		t.Fatalf("got %q, want %q", got, "message")
	}
}

func TestApply_TrimSuffix(t *testing.T) {
	tr, _ := New([]Rule{{Op: OpTrimSuffix, Arg: "\n"}})
	if got := tr.Apply("line\n"); got != "line" {
		t.Fatalf("got %q, want %q", got, "line")
	}
}

func TestApply_ChainedRules(t *testing.T) {
	tr, _ := New([]Rule{
		{Op: OpTrimSpace},
		{Op: OpLowercase},
		{Op: OpTrimPrefix, Arg: "debug "},
	})
	got := tr.Apply("  DEBUG message  ")
	want := "message"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApply_NoMatchingPrefix_Unchanged(t *testing.T) {
	tr, _ := New([]Rule{{Op: OpTrimPrefix, Arg: "ERR "}})
	if got := tr.Apply("INFO msg"); got != "INFO msg" {
		t.Fatalf("got %q, want %q", got, "INFO msg")
	}
}
