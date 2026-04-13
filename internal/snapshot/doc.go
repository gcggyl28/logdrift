// Package snapshot provides a thread-safe, fixed-size rolling buffer of
// log lines keyed by service name.
//
// # Overview
//
// A [Snapshot] stores the most recent N lines for each service seen by
// logdrift. Lines are added via [Snapshot.Push] and retrieved via
// [Snapshot.Lines]. When the buffer for a service is full the oldest line
// is evicted automatically.
//
// # Collector
//
// [Collector] sits on top of a Snapshot and consumes a fan-in stream of
// [tail.Line] values (produced by [tail.FanIn]). For every line it:
//  1. Pushes the text into the Snapshot for the originating service.
//  2. Forwards an [Entry] on its output channel so downstream stages
//     (e.g. the diff pipeline) can react in real time.
//
// Both types are safe for concurrent use.
package snapshot
