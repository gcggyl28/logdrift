// Package throttle implements per-service line-rate throttling for
// logdrift.
//
// During log storms a single noisy service can flood the diff pipeline
// with thousands of near-identical lines, producing meaningless drift
// reports and overwhelming downstream consumers.
//
// A [Throttle] tracks how many lines each service has emitted within a
// sliding time window and drops lines once the configured limit is
// reached. When the window expires the counter resets automatically,
// so bursts are absorbed without permanently silencing a service.
//
// Usage:
//
//	th, err := throttle.New(100, time.Second)
//	if err != nil { ... }
//
//	if th.Allow(entry.Service) {
//		downstream <- entry
//	}
//
// Setting maxPerWindow to 0 disables throttling entirely, which is the
// default when no throttle configuration is provided.
package throttle
