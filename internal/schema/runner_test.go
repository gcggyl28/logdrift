package schema

import (
	"context"
	"testing"
	"time"

	"github.com/user/logdrift/internal/parser"
)

func makeEntryCh(entries []parser.Entry) <-chan parser.Entry {
	ch := make(chan parser.Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func TestNewRunner_NilValidatorReturnsError(t *testing.T) {
	_, err := NewRunner(nil, make(chan parser.Entry))
	if err == nil {
		t.Fatal("expected error for nil validator")
	}
}

func TestNewRunner_NilChannelReturnsError(t *testing.T) {
	v, _ := New(baseFields())
	_, err := NewRunner(v, nil)
	if err == nil {
		t.Fatal("expected error for nil channel")
	}
}

func TestRunner_EmitsResults(t *testing.T) {
	v, _ := New(baseFields())
	entries := []parser.Entry{
		{Fields: map[string]string{"level": "info"}},
		{Fields: map[string]string{}},
	}
	r, _ := NewRunner(v, makeEntryCh(entries))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var results []Result
	for res := range r.Run(ctx) {
		results = append(results, res)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if len(results[0].Violations) != 0 {
		t.Fatalf("first entry should have no violations")
	}
	if len(results[1].Violations) == 0 {
		t.Fatal("second entry should have violations")
	}
}

func TestRunner_CancelStops(t *testing.T) {
	v, _ := New(baseFields())
	blocking := make(chan parser.Entry)
	r, _ := NewRunner(v, blocking)
	ctx, cancel := context.WithCancel(context.Background())
	out := r.Run(ctx)
	cancel()
	select {
	case <-out:
	case <-time.After(time.Second):
		t.Fatal("runner did not stop after context cancel")
	}
}
