// Package sampler provides probabilistic and deterministic sampling strategies
// for log lines in logdrift.
//
// Two modes are supported:
//
//   - ModeRandom: each log line is forwarded with probability Rate (0.0–1.0).
//     Useful for reducing volume while preserving a representative sample.
//
//   - ModeEveryN: every Nth line is forwarded deterministically.
//     Useful when you need predictable, evenly-spaced samples.
//
// Example:
//
//	s, err := sampler.New(sampler.Config{
//		Mode: sampler.ModeEveryN,
//		N:    10,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if s.Allow() {
//		// forward the line
//	}
package sampler
