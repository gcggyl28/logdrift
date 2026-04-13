// Package export provides structured export of drift reports produced by
// logdrift's diff pipeline.
//
// Supported output formats:
//
//	"json" – one JSON object per line (NDJSON), suitable for log aggregation
//	         pipelines and machine consumption.
//
//	"text" – human-readable single-line summary followed by the raw delta,
//	         suitable for writing to plain log files or stdout redirection.
//
// Typical usage:
//
//	ex, err := export.New(os.Stdout, export.FormatJSON)
//	if err != nil { ... }
//	ex.Write(export.DriftReport{ ... })
package export
