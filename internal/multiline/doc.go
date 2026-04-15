// Package multiline combines consecutive log lines that belong to a single
// logical event — such as Java stack traces or Python tracebacks — into one
// string before the event is forwarded to the rest of the logdrift pipeline.
//
// # Modes
//
// "prefix" — a new event begins whenever a line matches the configured
// regular expression. All subsequent non-matching lines are appended to that
// event until the next match arrives.
//
// "continuation" — a line that matches the pattern is treated as a
// continuation of the current event. The event is flushed when a
// non-matching line is seen.
//
// # Timeout flushing
//
// TimeoutFlusher runs a background goroutine that periodically checks whether
// the Combiner has been idle for longer than Timeout and, if so, forces a
// flush. This prevents incomplete events from being held in memory
// indefinitely when a service stops logging.
package multiline
