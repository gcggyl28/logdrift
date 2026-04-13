// Package filter implements log-line filtering for logdrift.
//
// A Filter matches individual log lines against one or more regular
// expressions. Filters can operate in "any" mode (OR semantics) or
// "all" mode (AND semantics), and can be inverted to act as exclusion
// filters.
//
// Multiple filters can be composed into a Chain, which passes a line
// only when every filter in the chain matches. Chain.Apply wraps a
// string channel and returns a new channel that emits only the lines
// that survive the chain, making it easy to plug filtering into the
// existing fan-in / pipeline architecture.
//
// Typical usage:
//
//	f, err := filter.New([]string{"error", "warn"}, filter.ModeAny, false)
//	chain := filter.NewChain(f)
//	filtered := chain.Apply(rawLines)
package filter
