// Package redact provides pattern-based redaction of sensitive values in log lines.
//
// A Redactor is configured with one or more regular expressions; any substring
// matching a pattern is replaced with a configurable mask string (default
// "[REDACTED]") before the line is passed downstream for display or export.
//
// Usage:
//
//	r, err := redact.New([]string{`password=\S+`, `token=[A-Za-z0-9]+`})
//	if err != nil {
//		log.Fatal(err)
//	}
//	safe := r.Apply(rawLine)
//
// The mask string can be overridden via WithMask:
//
//	r, _ := redact.New(patterns, redact.WithMask("<hidden>"))
package redact
