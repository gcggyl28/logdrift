// Package transform provides a lightweight, composable string transformation
// pipeline for log line content.
//
// A Transformer is constructed from an ordered slice of Rules. Each Rule
// specifies one Op (Uppercase, Lowercase, TrimSpace, TrimPrefix, or
// TrimSuffix) and an optional argument string required by the Trim* ops.
//
// Example usage:
//
//	tr, err := transform.New([]transform.Rule{
//		{Op: transform.OpTrimSpace},
//		{Op: transform.OpLowercase},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	out := tr.Apply("  WARN something happened  ")
//	// out == "warn something happened"
package transform
