// Package pattern provides compiled regular-expression pattern matching
// for log entry content within logdrift pipelines.
//
// A Matcher is constructed with a MatchMode and a list of regex strings.
// Two modes are supported:
//
//   - MatchAny: the input matches if at least one pattern matches (logical OR).
//   - MatchAll: the input matches only when every pattern matches (logical AND).
//
// When no patterns are supplied the Matcher is a no-op and Match always
// returns true, making it safe to embed in optional pipeline stages.
//
// Example:
//
//	m, err := pattern.New(pattern.MatchAny, []string{`error`, `panic`})
//	if err != nil { /* handle */ }
//	if m.Match(line) {
//	    // process interesting line
//	}
package pattern
