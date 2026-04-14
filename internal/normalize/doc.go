// Package normalize provides log line normalization before diffing.
//
// Volatile fields such as timestamps, process IDs, and UUIDs change on every
// run and would cause every line to appear as drift even when the underlying
// log structure is identical. Normalizer replaces these fields with a stable
// placeholder so that the diff engine can focus on meaningful differences.
//
// Usage:
//
//	n, err := normalize.New(normalize.ModeFull, nil, "")
//	if err != nil { ... }
//	normalizedLine := n.Apply(rawLine)
//
// Modes:
//
//	- none      – no substitutions; lines are passed through unchanged.
//	- timestamp – replace ISO-8601 / RFC-3339 timestamps.
//	- pid       – replace pid=N / pid:N occurrences.
//	- full      – apply all built-in substitutions (timestamp + pid + uuid).
//
// Additional patterns can be supplied as raw regular expressions via the
// extraPatterns argument of New.
package normalize
