// Package schema validates parsed log entries against a set of required and
// optional field definitions, reporting missing or mistyped fields.
package schema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/user/logdrift/internal/parser"
)

// FieldType describes the expected type of a log field value.
type FieldType string

const (
	FieldTypeString FieldType = "string"
	FieldTypeNumber FieldType = "number"
	FieldTypeBool   FieldType = "bool"
)

// FieldDef defines a single expected field in a log entry.
type FieldDef struct {
	Name     string
	Type     FieldType
	Required bool
}

// Violation describes a single schema violation found in an entry.
type Violation struct {
	Field   string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("field %q: %s", v.Field, v.Message)
}

// Validator checks log entries against a schema.
type Validator struct {
	fields []FieldDef
}

// New creates a Validator from the provided field definitions.
// Returns an error if no fields are provided or a field name is empty.
func New(fields []FieldDef) (*Validator, error) {
	if len(fields) == 0 {
		return nil, errors.New("schema: at least one field definition required")
	}
	for _, f := range fields {
		if strings.TrimSpace(f.Name) == "" {
			return nil, errors.New("schema: field name must not be empty")
		}
		switch f.Type {
		case FieldTypeString, FieldTypeNumber, FieldTypeBool:
		default:
			return nil, fmt.Errorf("schema: unknown field type %q for field %q", f.Type, f.Name)
		}
	}
	return &Validator{fields: fields}, nil
}

// Validate checks the entry against the schema and returns any violations.
func (v *Validator) Validate(entry parser.Entry) []Violation {
	var violations []Violation
	for _, def := range v.fields {
		val, ok := entry.Fields[def.Name]
		if !ok {
			if def.Required {
				violations = append(violations, Violation{Field: def.Name, Message: "required field missing"})
			}
			continue
		}
		if err := checkType(def.Name, def.Type, val); err != nil {
			violations = append(violations, *err)
		}
	}
	return violations
}

func checkType(name string, ft FieldType, val string) *Violation {
	switch ft {
	case FieldTypeNumber:
		var f float64
		if _, err := fmt.Sscanf(val, "%g", &f); err != nil {
			return &Violation{Field: name, Message: fmt.Sprintf("expected number, got %q", val)}
		}
	case FieldTypeBool:
		low := strings.ToLower(val)
		if low != "true" && low != "false" {
			return &Violation{Field: name, Message: fmt.Sprintf("expected bool, got %q", val)}
		}
	}
	return nil
}
