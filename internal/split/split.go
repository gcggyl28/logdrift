// Package split provides a line splitter that breaks raw log lines into
// named fields using a configurable delimiter.
package split

import (
	"errors"
	"fmt"
	"strings"
)

// Splitter splits a raw line into a map of named fields.
type Splitter struct {
	delimiter string
	fields    []string
	trim      bool
}

// New returns a Splitter that splits on delimiter and maps positional
// segments to the provided field names. If trim is true, whitespace is
// stripped from each segment. At least one field name is required.
func New(delimiter string, fields []string, trim bool) (*Splitter, error) {
	if delimiter == "" {
		return nil, errors.New("split: delimiter must not be empty")
	}
	if len(fields) == 0 {
		return nil, errors.New("split: at least one field name is required")
	}
	for i, f := range fields {
		if f == "" {
			return nil, fmt.Errorf("split: field name at index %d must not be empty", i)
		}
	}
	return &Splitter{delimiter: delimiter, fields: fields, trim: trim}, nil
}

// Apply splits line and returns a map of field name to segment value.
// Extra segments beyond the number of field names are collected under
// the last field name, joined by the delimiter. Missing segments produce
// empty string values.
func (s *Splitter) Apply(line string) map[string]string {
	parts := strings.SplitN(line, s.delimiter, len(s.fields))
	out := make(map[string]string, len(s.fields))
	for i, name := range s.fields {
		var val string
		if i < len(parts) {
			val = parts[i]
		}
		if s.trim {
			val = strings.TrimSpace(val)
		}
		out[name] = val
	}
	return out
}
