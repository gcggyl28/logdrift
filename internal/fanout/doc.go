// Package fanout provides a Broadcaster that fans out a single stream of
// log entries to multiple independent subscriber channels.
//
// # Usage
//
//	b, err := fanout.New(src, 32)
//	if err != nil { ... }
//
//	sub1 := b.Subscribe()
//	sub2 := b.Subscribe()
//
//	go b.Run(ctx)
//
// Each subscriber receives every entry. If a subscriber's buffer is full the
// entry is dropped for that subscriber rather than blocking the broadcast loop.
package fanout
