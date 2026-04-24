package topology

import (
	"errors"
	"strings"
)

const (
	// FieldUpstream is the entry field key written by Annotator.
	FieldUpstream = "upstream"
	// FieldDownstream is the entry field key written by Annotator.
	FieldDownstream = "downstream"
)

// Entry is a minimal log entry interface expected by Annotator.
type Entry struct {
	Service string
	Fields  map[string]string
}

// Annotator enriches log entries with topology metadata.
type Annotator struct {
	graph *Graph
}

// NewAnnotator creates an Annotator backed by the given Graph.
func NewAnnotator(g *Graph) (*Annotator, error) {
	if g == nil {
		return nil, errors.New("topology: graph must not be nil")
	}
	return &Annotator{graph: g}, nil
}

// Annotate adds upstream and downstream fields to the entry in-place.
// It is safe to call concurrently.
func (a *Annotator) Annotate(e *Entry) {
	if e == nil {
		return
	}
	if e.Fields == nil {
		e.Fields = make(map[string]string)
	}
	up := a.graph.Upstream(e.Service)
	down := a.graph.Downstream(e.Service)
	if len(up) > 0 {
		e.Fields[FieldUpstream] = strings.Join(up, ",")
	}
	if len(down) > 0 {
		e.Fields[FieldDownstream] = strings.Join(down, ",")
	}
}
