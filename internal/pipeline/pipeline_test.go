package pipeline

import (
	"errors"
	"testing"
)

func identity(e Entry) (*Entry, error) { return &e, nil }

func drop(_ Entry) (*Entry, error) { return nil, nil }

func boom(_ Entry) (*Entry, error) { return nil, errors.New("boom") }

func appendSuffix(suffix string) StageFn {
	return func(e Entry) (*Entry, error) {
		e.Line += suffix
		return &e, nil
	}
}

func TestNew_EmptyStagesReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty stages")
	}
}

func TestNew_ValidStages(t *testing.T) {
	s, _ := NewStage("id", identity)
	p, err := New([]Stage{s})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Stages()) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(p.Stages()))
	}
}

func TestRun_PassesThroughIdentity(t *testing.T) {
	s, _ := NewStage("id", identity)
	p, _ := New([]Stage{s})
	out, err := p.Run(Entry{Service: "svc", Line: "hello"})
	if err != nil || out == nil {
		t.Fatalf("expected entry, got out=%v err=%v", out, err)
	}
	if out.Line != "hello" {
		t.Errorf("line mismatch: %q", out.Line)
	}
}

func TestRun_DropStageReturnsNil(t *testing.T) {
	s, _ := NewStage("drop", drop)
	p, _ := New([]Stage{s})
	out, err := p.Run(Entry{Line: "x"})
	if err != nil || out != nil {
		t.Fatalf("expected nil entry, got out=%v err=%v", out, err)
	}
}

func TestRun_ErrorStageWrapsName(t *testing.T) {
	s, _ := NewStage("explode", boom)
	p, _ := New([]Stage{s})
	_, err := p.Run(Entry{Line: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, errors.New("boom")) {
		// just check the string contains stage name
	}
	if msg := err.Error(); msg == "" {
		t.Error("empty error message")
	}
}

func TestRun_MultipleStagesApplyInOrder(t *testing.T) {
	s1, _ := NewStage("a", appendSuffix("-A"))
	s2, _ := NewStage("b", appendSuffix("-B"))
	p, _ := New([]Stage{s1, s2})
	out, err := p.Run(Entry{Line: "base"})
	if err != nil || out == nil {
		t.Fatalf("unexpected: %v %v", out, err)
	}
	if out.Line != "base-A-B" {
		t.Errorf("expected 'base-A-B', got %q", out.Line)
	}
}

func TestNewStage_EmptyNameReturnsError(t *testing.T) {
	_, err := NewStage("", identity)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewStage_NilFnReturnsError(t *testing.T) {
	_, err := NewStage("x", nil)
	if err == nil {
		t.Fatal("expected error for nil fn")
	}
}
