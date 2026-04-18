// Package junction provides a Junction type that merges multiple named
// log-line channels into a single stream of tagged Entry values.
//
// Each source channel is associated with a service name. Lines emitted
// from a source are forwarded to the output channel with the originating
// service name attached, enabling downstream consumers to identify where
// each log line came from without maintaining per-service goroutines
// themselves.
//
// Usage:
//
//	sources := map[string]<-chan diff.Line{
//		"api":    apiCh,
//		"worker": workerCh,
//	}
//	j, err := junction.New(sources)
//	if err != nil { ... }
//	for entry := range j.Run(ctx) {
//		fmt.Printf("[%s] %s\n", entry.Service, entry.Line)
//	}
package junction
