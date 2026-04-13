// Package buffer provides a thread-safe, fixed-capacity ring buffer for
// log entries and a periodic Flusher that drains the buffer on a configurable
// interval.
//
// # RingBuffer
//
// RingBuffer stores [Entry] values (service name + log line) up to a fixed
// capacity. When the buffer is full, the oldest entry is silently overwritten.
// All methods are safe for concurrent use.
//
// # Flusher
//
// Flusher wraps a RingBuffer and calls a user-supplied [FlushFunc] on every
// tick. After each successful flush the buffer is reset. A final flush is
// performed when the context is cancelled so no entries are lost on shutdown.
package buffer
