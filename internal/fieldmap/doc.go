// Package fieldmap provides field renaming and dropping for structured log
// entries produced by the parser package.
//
// # Overview
//
// Different services often emit the same semantic data under different key
// names — one service writes "msg", another writes "message", a third writes
// "log". fieldmap lets you define a set of [Rule] values that are applied to
// the parsed field map before the entry is forwarded to the diff pipeline.
//
// # Usage
//
//	rules := []fieldmap.Rule{
//		{From: "msg",    To: "message"},  // rename
//		{From: "secret", To: ""},         // drop
//	}
//	m, err := fieldmap.New(rules)
//	if err != nil { ... }
//	normalised := m.Apply(parsedFields)
//
// Fields that do not match any rule are passed through unchanged.
package fieldmap
