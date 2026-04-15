// Package debounce implements entry-level debouncing for log streams.
//
// A Debouncer delays forwarding of log entries until a configurable quiet
// window has elapsed with no new entry for the same service. Rapid bursts
// of lines from a single service are collapsed so that only the most recent
// entry is forwarded, reducing noise in diff and alert pipelines.
//
// Basic usage:
//
//	d, err := debounce.New(50 * time.Millisecond)
//	if err != nil { ... }
//
//	d.Push(debounce.Entry{Service: "api", Line: "request failed"})
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	d.Drain(ctx, func(e debounce.Entry) {
//		fmt.Println(e.Service, e.Line)
//	})
//
// For pipeline integration use NewRunner which wires a source channel
// directly into a Debouncer and calls Drain automatically.
package debounce
