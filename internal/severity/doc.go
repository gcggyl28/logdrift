// Package severity provides log-level classification for parsed log entries.
//
// A Classifier inspects one or more named fields (e.g. "level", "severity",
// "lvl") in a parser.Entry and maps the string value to a canonical Level
// constant (Debug, Info, Warn, Error, Fatal). Field names are tried in the
// order supplied; the first match wins.
//
// Example:
//
//	c, err := severity.New("level", "severity")
//	if err != nil { ... }
//	lvl := c.Classify(entry)  // severity.LevelError, etc.
package severity
