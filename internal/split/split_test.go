package split

import (
	"testing"
)

func TestNew_EmptyDelimiterReturnsError(t *testing.T) {
	_, err := New("", []string{"a"}, false)
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestNew_NoFieldsReturnsError(t *testing.T) {
	_, err := New("|", []string{}, false)
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_EmptyFieldNameReturnsError(t *testing.T) {
	_, err := New("|", []string{"a", ""}, false)
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_Valid(t *testing.T) {
	s, err := New("|", []string{"level", "msg"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil splitter")
	}
}

func TestApply_BasicSplit(t *testing.T) {
	s, _ := New("|", []string{"level", "msg"}, false)
	out := s.Apply("INFO|hello world")
	if out["level"] != "INFO" {
		t.Errorf("expected INFO got %q", out["level"])
	}
	if out["msg"] != "hello world" {
		t.Errorf("expected 'hello world' got %q", out["msg"])
	}
}

func TestApply_TrimWhitespace(t *testing.T) {
	s, _ := New("|", []string{"level", "msg"}, true)
	out := s.Apply(" INFO | hello ")
	if out["level"] != "INFO" {
		t.Errorf("expected INFO got %q", out["level"])
	}
	if out["msg"] != "hello" {
		t.Errorf("expected 'hello' got %q", out["msg"])
	}
}

func TestApply_MissingSegmentsProduceEmpty(t *testing.T) {
	s, _ := New("|", []string{"a", "b", "c"}, false)
	out := s.Apply("only")
	if out["a"] != "only" {
		t.Errorf("unexpected a: %q", out["a"])
	}
	if out["b"] != "" || out["c"] != "" {
		t.Error("expected empty strings for missing segments")
	}
}

func TestApply_ExtraSegmentsCollapsedIntoLast(t *testing.T) {
	s, _ := New("|", []string{"level", "rest"}, false)
	out := s.Apply("INFO|part1|part2|part3")
	if out["rest"] != "part1|part2|part3" {
		t.Errorf("unexpected rest: %q", out["rest"])
	}
}
