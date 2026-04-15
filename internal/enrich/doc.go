// Package enrich provides a lightweight enrichment stage that stamps
// structured metadata fields onto log entries as they flow through the
// logdrift pipeline.
//
// # Overview
//
// An [Enricher] holds a fixed map of key/value pairs that are merged into
// every [Entry] it processes. Fields already present on the entry take
// precedence, so enrichment acts as a "default values" layer rather than
// an override.
//
// # Usage
//
//	e, err := enrich.New(map[string]string{
//		"env":    "production",
//		"region": "us-east-1",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	enriched := e.Apply(entry)
package enrich
