package filter_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/filter"
)

func TestChain_EmptyPassesAll(t *testing.T) {
	c := filter.NewChain()
	if !c.Match("anything") {
		t.Error("empty chain should pass all lines")
	}
}

func TestChain_MultipleFilters(t *testing.T) {
	f1, _ := filter.New([]string{"error"}, filter.ModeAny, false)
	f2, _ := filter.New([]string{"disk"}, filter.ModeAny, false)
	c := filter.NewChain(f1, f2)

	if !c.Match("error: disk full") {
		t.Error("expected match when both filters satisfied")
	}
	if c.Match("error: network") {
		t.Error("expected no match when second filter not satisfied")
	}
	if c.Match("disk: ok") {
		t.Error("expected no match when first filter not satisfied")
	}
}

func TestChain_Apply(t *testing.T) {
	f, _ := filter.New([]string{"keep"}, filter.ModeAny, false)
	c := filter.NewChain(f)

	in := make(chan string, 4)
	in <- "keep this"
	in <- "drop this"
	in <- "keep me too"
	in <- "ignore"
	close(in)

	out := c.Apply(in)
	var results []string
	for line := range out {
		results = append(results, line)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(results), results)
	}
	if results[0] != "keep this" || results[1] != "keep me too" {
		t.Errorf("unexpected results: %v", results)
	}
}
