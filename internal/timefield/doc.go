// Package timefield provides utilities for extracting, parsing, and stamping
// timestamps within structured log entries.
//
// # Extractor
//
// Extractor reads a named field from a parser.Entry and attempts to parse it
// as a time.Time value using a configurable list of layouts. When WithFallback
// is set, missing or unparseable fields yield time.Now() instead of an error.
//
// # Stamper
//
// Stamper wraps an Extractor and writes the current UTC time into the
// configured field whenever the entry does not already carry a valid timestamp.
// It is useful as a pre-processing step before entries reach the diff pipeline.
package timefield
