package timefield_test

import (
	"testing"
	"time"

	"github.com/angryboat/logdrift/internal/parser"
	"github.com/angryboat/logdrift/internal/timefield"
)

func TestNewStamper_EmptyFieldReturnsError(t *testing.T) {
	_, err := timefield.NewStamper("", time.RFC3339)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewStamper_EmptyLayoutReturnsError(t *testing.T) {
	_, err := timefield.NewStamper("ts", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewStamper_Valid(t *testing.T) {
	_, err := timefield.NewStamper("ts", time.RFC3339)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStamper_Apply_AddsFieldWhenAbsent(t *testing.T) {
	s, _ := timefield.NewStamper("ts", time.RFC3339)
	out := s.Apply(parser.Entry{Fields: nil})
	if out.Fields["ts"] == "" {
		t.Fatal("expected ts field to be populated")
	}
}

func TestStamper_Apply_PreservesExistingTimestamp(t *testing.T) {
	s, _ := timefield.NewStamper("ts", time.RFC3339)
	original := "2024-01-02T03:04:05Z"
	out := s.Apply(parser.Entry{Fields: map[string]string{"ts": original}})
	if out.Fields["ts"] != original {
		t.Fatalf("expected %q got %q", original, out.Fields["ts"])
	}
}

func TestStamper_Apply_NilFields_Initialised(t *testing.T) {
	s, _ := timefield.NewStamper("ts", time.RFC3339)
	out := s.Apply(parser.Entry{})
	if out.Fields == nil {
		t.Fatal("Fields map should be initialised")
	}
}
