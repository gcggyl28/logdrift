// Package correlate links log entries across services by a shared field value.
package correlate

import (
	"errors"
	"sync"
	"time"
)

// Entry is a minimal log entry used for correlation.
type Entry struct {
	Service string
	Line    string
	Fields  map[string]string
	At      time.Time
}

// Group holds entries that share the same correlation key value.
type Group struct {
	Key     string
	Value   string
	Entries []Entry
}

// Correlator indexes entries by a field value and emits groups.
type Correlator struct {
	field  string
	ttl    time.Duration
	mu     sync.Mutex
	bucket map[string]*group
}

type group struct {
	entries []Entry
	updated time.Time
}

// New creates a Correlator that groups entries sharing the same value for field.
// ttl controls how long a group is retained without new entries.
func New(field string, ttl time.Duration) (*Correlator, error) {
	if field == "" {
		return nil, errors.New("correlate: field must not be empty")
	}
	if ttl <= 0 {
		return nil, errors.New("correlate: ttl must be positive")
	}
	return &Correlator{field: field, ttl: ttl, bucket: make(map[string]*group)}, nil
}

// Add indexes an entry. Returns the current group for its correlation value.
func (c *Correlator) Add(e Entry) Group {
	v := e.Fields[c.field]
	c.mu.Lock()
	defer c.mu.Unlock()
	g, ok := c.bucket[v]
	if !ok {
		g = &group{}
		c.bucket[v] = g
	}
	g.entries = append(g.entries, e)
	g.updated = time.Now()
	out := make([]Entry, len(g.entries))
	copy(out, g.entries)
	return Group{Key: c.field, Value: v, Entries: out}
}

// Evict removes groups that have not been updated within ttl.
func (c *Correlator) Evict() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, g := range c.bucket {
		if now.Sub(g.updated) > c.ttl {
			delete(c.bucket, k)
		}
	}
}

// Len returns the number of active correlation groups.
func (c *Correlator) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.bucket)
}
