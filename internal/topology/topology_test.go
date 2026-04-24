package topology

import (
	"testing"
)

func TestAddEdge_Valid(t *testing.T) {
	g := New()
	if err := g.AddEdge("api", "db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddEdge_EmptyFrom(t *testing.T) {
	g := New()
	if err := g.AddEdge("", "db"); err == nil {
		t.Fatal("expected error for empty from")
	}
}

func TestAddEdge_EmptyTo(t *testing.T) {
	g := New()
	if err := g.AddEdge("api", ""); err == nil {
		t.Fatal("expected error for empty to")
	}
}

func TestAddEdge_SelfLoop(t *testing.T) {
	g := New()
	if err := g.AddEdge("api", "api"); err == nil {
		t.Fatal("expected error for self-loop")
	}
}

func TestAddEdge_Duplicate(t *testing.T) {
	g := New()
	_ = g.AddEdge("api", "db")
	if err := g.AddEdge("api", "db"); err == nil {
		t.Fatal("expected error for duplicate edge")
	}
}

func TestDownstream_Known(t *testing.T) {
	g := New()
	_ = g.AddEdge("api", "db")
	_ = g.AddEdge("api", "cache")
	down := g.Downstream("api")
	if len(down) != 2 {
		t.Fatalf("expected 2 downstream, got %d", len(down))
	}
}

func TestDownstream_Unknown(t *testing.T) {
	g := New()
	if got := g.Downstream("unknown"); len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestUpstream_Known(t *testing.T) {
	g := New()
	_ = g.AddEdge("api", "db")
	_ = g.AddEdge("worker", "db")
	up := g.Upstream("db")
	if len(up) != 2 {
		t.Fatalf("expected 2 upstream, got %d", len(up))
	}
}

func TestUpstream_Unknown(t *testing.T) {
	g := New()
	if got := g.Upstream("ghost"); len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestServices_ReturnsAll(t *testing.T) {
	g := New()
	_ = g.AddEdge("api", "db")
	_ = g.AddEdge("api", "cache")
	svc := g.Services()
	if len(svc) != 3 {
		t.Fatalf("expected 3 services, got %d: %v", len(svc), svc)
	}
}
