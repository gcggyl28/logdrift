// Package topology tracks directed relationships between services,
// allowing logdrift to understand upstream/downstream dependencies
// and annotate log entries with propagation context.
package topology

import (
	"errors"
	"fmt"
	"sync"
)

// Edge represents a directed dependency from one service to another.
type Edge struct {
	From string
	To   string
}

// Graph holds the service dependency graph.
type Graph struct {
	mu    sync.RWMutex
	edges map[string][]string // from -> []to
}

// New creates an empty Graph.
func New() *Graph {
	return &Graph{edges: make(map[string][]string)}
}

// AddEdge registers a directed edge from -> to.
// Both service names must be non-empty and distinct.
func (g *Graph) AddEdge(from, to string) error {
	if from == "" {
		return errors.New("topology: from service must not be empty")
	}
	if to == "" {
		return errors.New("topology: to service must not be empty")
	}
	if from == to {
		return fmt.Errorf("topology: self-loop not allowed for service %q", from)
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, existing := range g.edges[from] {
		if existing == to {
			return fmt.Errorf("topology: edge %q->%q already exists", from, to)
		}
	}
	g.edges[from] = append(g.edges[from], to)
	return nil
}

// Downstream returns the direct downstream services of the given service.
func (g *Graph) Downstream(service string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]string, len(g.edges[service]))
	copy(out, g.edges[service])
	return out
}

// Upstream returns all services that have a direct edge to the given service.
func (g *Graph) Upstream(service string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var up []string
	for from, tos := range g.edges {
		for _, to := range tos {
			if to == service {
				up = append(up, from)
				break
			}
		}
	}
	return up
}

// Services returns all service names that appear in at least one edge.
func (g *Graph) Services() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	seen := make(map[string]struct{})
	for from, tos := range g.edges {
		seen[from] = struct{}{}
		for _, to := range tos {
			seen[to] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	return out
}
