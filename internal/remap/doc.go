// Package remap provides field-value remapping for log entries.
//
// A Remapper is configured with one or more Rules, each targeting a specific
// field and supplying a lookup table of old→new values. Fields or values that
// do not match any rule are passed through unchanged.
//
// Example:
//
//	r, err := remap.New([]remap.Rule{
//		{
//			Field:   "level",
//			Mapping: map[string]string{"warn": "warning", "err": "error"},
//		},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	out := r.Apply(entry)
package remap
