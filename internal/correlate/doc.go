// Package correlate groups log entries from multiple services by a shared
// field value (e.g. request_id or trace_id), enabling cross-service
// correlation of log activity in real time.
//
// # Usage
//
//	c, err := correlate.New("request_id", 30*time.Second)
//	src := make(chan correlate.Entry)
//	out := make(chan correlate.Group, 64)
//	r, err := correlate.NewRunner(c, src, out, 10*time.Second)
//	go r.Run(ctx)
//
// Each time an Entry is added whose field matches an existing group, the full
// group (all correlated entries so far) is emitted on the output channel.
// Stale groups are evicted periodically based on the configured TTL.
package correlate
