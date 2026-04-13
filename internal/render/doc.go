// Package render provides terminal rendering for logdrift drift output.
//
// Three formats are supported:
//
//   - plain     – plain text, no ANSI codes, suitable for piping / log files.
//   - colored   – ANSI-colored diff lines (green additions, red deletions).
//   - timestamp – like colored but each drift block is prefixed with an
//                 RFC3339 UTC timestamp so you can correlate events over time.
//
// Basic usage:
//
//	r, err := render.New(render.FormatColored)
//	if err != nil {
//		log.Fatal(err)
//	}
//	r.WriteDrift("service-a", "service-b", diffDelta)
package render
