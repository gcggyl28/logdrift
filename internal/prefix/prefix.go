// Package prefix prepends a static or per-service string to log line text.
package prefix

import (
	"errors"
	"fmt"

	"github.com/logdrift/logdrift/internal/snapshot"
)

// Prefixer prepends a configured string to entry lines.
type Prefixer struct {
	global  string
	service map[string]string
}

// Option configures a Prefixer.
type Option func(*Prefixer)

// WithGlobal sets a prefix applied to every entry.
func WithGlobal(p string) Option {
	return func(pr *Prefixer) { pr.global = p }
}

// WithService sets a prefix for a specific service, overriding the global one.
func WithService(service, p string) Option {
	return func(pr *Prefixer) { pr.service[service] = p }
}

// New creates a Prefixer. At least one option must produce a non-empty prefix.
func New(opts ...Option) (*Prefixer, error) {
	pr := &Prefixer{service: make(map[string]string)}
	for _, o := range opts {
		o(pr)
	}
	if pr.global == "" && len(pr.service) == 0 {
		return nil, errors.New("prefix: at least one prefix must be configured")
	}
	return pr, nil
}

// Apply prepends the appropriate prefix to entry.Line and returns the modified entry.
func (pr *Prefixer) Apply(e snapshot.Entry) (snapshot.Entry, error) {
	p := pr.global
	if svc, ok := pr.service[e.Service]; ok {
		p = svc
	}
	if p != "" {
		e.Line = fmt.Sprintf("%s%s", p, e.Line)
	}
	return e, nil
}
