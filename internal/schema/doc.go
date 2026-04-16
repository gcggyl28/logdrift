// Package schema provides field-level validation for parsed log entries.
//
// A Validator is built from a slice of FieldDef values, each specifying a
// field name, its expected type (string, number, or bool), and whether the
// field is required. Calling Validate on a parser.Entry returns a slice of
// Violation values describing any missing required fields or type mismatches.
//
// For streaming use, Runner wraps a Validator and reads from a channel of
// parser.Entry values, emitting Result values that pair each entry with its
// violations. The runner stops when the source channel is closed or the
// context is cancelled.
//
// Example:
//
//	v, err := schema.New([]schema.FieldDef{
//	    {Name: "level", Type: schema.FieldTypeString, Required: true},
//	    {Name: "status", Type: schema.FieldTypeNumber},
//	})
//	violations := v.Validate(entry)
package schema
