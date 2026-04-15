// Package routing provides log-line routing based on service name patterns,
// directing entries to one or more named output channels.
package routing

import (
	"fmt"
	"regexp"
	"sync"
)

// Route maps a compiled pattern to a destination label.
type Route struct {
	pattern *regexp.Regexp
	Dest    string
}

// Router holds a set of routes and dispatches log lines to matching destinations.
type Router struct {
	mu     sync.RWMutex
	routes []Route
}

// New creates a Router from a map of destination→pattern strings.
// An error is returned if any pattern fails to compile.
func New(rules map[string]string) (*Router, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("routing: at least one rule is required")
	}
	r := &Router{}
	for dest, pat := range rules {
		if pat == "" {
			return nil, fmt.Errorf("routing: empty pattern for destination %q", dest)
		}
		compiled, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("routing: invalid pattern %q for destination %q: %w", pat, dest, err)
		}
		r.routes = append(r.routes, Route{pattern: compiled, Dest: dest})
	}
	return r, nil
}

// Match returns all destination labels whose pattern matches the given service name.
// If no route matches, an empty slice is returned.
func (r *Router) Match(service string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var dests []string
	for _, route := range r.routes {
		if route.pattern.MatchString(service) {
			dests = append(dests, route.Dest)
		}
	}
	return dests
}

// Routes returns a snapshot of all registered routes.
func (r *Router) Routes() []Route {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Route, len(r.routes))
	copy(out, r.routes)
	return out
}
