package topology

import (
	"strings"
	"testing"
)

func TestNewAnnotator_NilGraphReturnsError(t *testing.T) {
	if _, err := NewAnnotator(nil); err == nil {
		t.Fatal("expected error for nil graph")
	}
}

func TestNewAnnotator_Valid(t *testing.T) {
	g := New()
	if _, err := NewAnnotator(g); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnnotate_NilEntryNoOp(t *testing.T) {
	g := New()
	a, _ := NewAnnotator(g)
	a.Annotate(nil) // must not panic
}

func TestAnnotate_AddsUpstreamAndDownstream(t *testing.T) {
	g := New()
	_ = g.AddEdge("api", "db")
	_ = g.AddEdge("worker", "db")
	a, _ := NewAnnotator(g)

	e := &Entry{Service: "db", Fields: map[string]string{}}
	a.Annotate(e)

	if _, ok := e.Fields[FieldUpstream]; !ok {
		t.Fatal("expected upstream field")
	}
	if _, ok := e.Fields[FieldDownstream]; ok {
		t.Fatal("db has no downstream, field should be absent")
	}
}

func TestAnnotate_UpstreamContainsBothServices(t *testing.T) {
	g := New()
	_ = g.AddEdge("api", "db")
	_ = g.AddEdge("worker", "db")
	a, _ := NewAnnotator(g)

	e := &Entry{Service: "db"}
	a.Annotate(e)

	upVal := e.Fields[FieldUpstream]
	if !strings.Contains(upVal, "api") || !strings.Contains(upVal, "worker") {
		t.Fatalf("expected both api and worker in upstream, got %q", upVal)
	}
}

func TestAnnotate_NilFieldsInitialised(t *testing.T) {
	g := New()
	_ = g.AddEdge("api", "db")
	a, _ := NewAnnotator(g)

	e := &Entry{Service: "api"} // Fields is nil
	a.Annotate(e)

	if e.Fields == nil {
		t.Fatal("expected Fields to be initialised")
	}
	if _, ok := e.Fields[FieldDownstream]; !ok {
		t.Fatal("expected downstream field for api")
	}
}

func TestAnnotate_UnknownService_NoFields(t *testing.T) {
	g := New()
	a, _ := NewAnnotator(g)

	e := &Entry{Service: "ghost", Fields: map[string]string{}}
	a.Annotate(e)

	if len(e.Fields) != 0 {
		t.Fatalf("expected no fields for unknown service, got %v", e.Fields)
	}
}
