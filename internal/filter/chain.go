package filter

// Chain applies multiple Filters in sequence.
// A line passes only if every Filter in the chain matches it.
type Chain struct {
	filters []*Filter
}

// NewChain creates a Chain from the provided filters.
func NewChain(filters ...*Filter) *Chain {
	return &Chain{filters: filters}
}

// Match reports whether line passes all filters in the chain.
func (c *Chain) Match(line string) bool {
	for _, f := range c.filters {
		if !f.Match(line) {
			return false
		}
	}
	return true
}

// Apply reads from in and writes matching lines to the returned channel.
// The returned channel is closed when in is closed.
func (c *Chain) Apply(in <-chan string) <-chan string {
	out := make(chan string, 64)
	go func() {
		defer close(out)
		for line := range in {
			if c.Match(line) {
				out <- line
			}
		}
	}()
	return out
}
