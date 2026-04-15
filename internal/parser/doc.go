// Package parser provides structured log line parsing for logdrift.
//
// Supported formats:
//
//   - auto    – heuristically selects JSON, logfmt, or common log format
//   - json    – newline-delimited JSON (e.g. zerolog, zap, logrus JSON)
//   - logfmt  – key=value pairs (e.g. logfmt, Go standard library slog text)
//   - common  – Apache / nginx combined/common access log
//
// Usage:
//
//	p, err := parser.New(parser.FormatAuto)
//	if err != nil {
//		log.Fatal(err)
//	}
//	entry, err := p.Parse(rawLine)
//
Each parsed Entry exposes Timestamp, Level, Message, and a Fields map for
arbitrary key-value pairs, as well as the original Raw line for pass-through
use cases such as rendering and diffing.
package parser
