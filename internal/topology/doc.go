// Package topology provides a directed service dependency graph and an
// Annotator that enriches log entries with upstream/downstream context.
//
// # Graph
//
// Graph tracks directed edges between named services. Edges are added with
// AddEdge and queried with Upstream and Downstream. All methods are safe
// for concurrent use.
//
// # Annotator
//
// Annotator wraps a Graph and writes topology metadata into log entry
// fields. For each entry it resolves the upstream and downstream service
// sets and stores them as comma-separated strings under the keys
// "upstream" and "downstream".
//
// Example:
//
//	g := topology.New()
//	_ = g.AddEdge("api", "db")
//	a, _ := topology.NewAnnotator(g)
//	a.Annotate(entry)
package topology
