package prefix_test

import (
	"testing"

	"github.com/logdrift/logdrift/internal/prefix"
	"github.com/logdrift/logdrift/internal/snapshot"
)

func entry(svc, line string) snapshot.Entry {
	return snapshot.Entry{Service: svc, Line: line}
}

func TestNew_NoOptionsReturnsError(t *testing.T) {
	_, err := prefix.New()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNew_GlobalOnly(t *testing.T) {
	_, err := prefix.New(prefix.WithGlobal("[INFO] "))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_ServiceOnly(t *testing.T) {
	_, err := prefix.New(prefix.WithService("api", "[api] "))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_GlobalPrefix(t *testing.T) {
	pr, _ := prefix.New(prefix.WithGlobal(">> "))
	out, err := pr.Apply(entry("svc", "hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Line != ">> hello" {
		t.Fatalf("got %q", out.Line)
	}
}

func TestApply_ServiceOverridesGlobal(t *testing.T) {
	pr, _ := prefix.New(
		prefix.WithGlobal("[global] "),
		prefix.WithService("api", "[api] "),
	)
	out, err := pr.Apply(entry("api", "request"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Line != "[api] request" {
		t.Fatalf("got %q", out.Line)
	}
}

func TestApply_UnknownServiceUsesGlobal(t *testing.T) {
	pr, _ := prefix.New(
		prefix.WithGlobal("[global] "),
		prefix.WithService("api", "[api] "),
	)
	out, err := pr.Apply(entry("worker", "task done"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Line != "[global] task done" {
		t.Fatalf("got %q", out.Line)
	}
}

func TestApply_NoGlobalNoServiceMatch_NoChange(t *testing.T) {
	pr, _ := prefix.New(prefix.WithService("api", "[api] "))
	out, err := pr.Apply(entry("worker", "task done"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Line != "task done" {
		t.Fatalf("got %q", out.Line)
	}
}
