package schema

import (
	"testing"

	"github.com/user/logdrift/internal/parser"
)

func baseFields() []FieldDef {
	return []FieldDef{
		{Name: "level", Type: FieldTypeString, Required: true},
		{Name: "status", Type: FieldTypeNumber, Required: false},
		{Name: "ok", Type: FieldTypeBool, Required: false},
	}
}

func TestNew_EmptyFieldsReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_EmptyFieldNameReturnsError(t *testing.T) {
	_, err := New([]FieldDef{{Name: "", Type: FieldTypeString}})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_UnknownTypeReturnsError(t *testing.T) {
	_, err := New([]FieldDef{{Name: "x", Type: "uuid"}})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := New(baseFields())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_NoViolations(t *testing.T) {
	v, _ := New(baseFields())
	entry := parser.Entry{Fields: map[string]string{"level": "info", "status": "200", "ok": "true"}}
	if vv := v.Validate(entry); len(vv) != 0 {
		t.Fatalf("expected no violations, got %v", vv)
	}
}

func TestValidate_MissingRequiredField(t *testing.T) {
	v, _ := New(baseFields())
	entry := parser.Entry{Fields: map[string]string{"status": "200"}}
	vv := v.Validate(entry)
	if len(vv) != 1 || vv[0].Field != "level" {
		t.Fatalf("expected missing 'level' violation, got %v", vv)
	}
}

func TestValidate_WrongNumberType(t *testing.T) {
	v, _ := New(baseFields())
	entry := parser.Entry{Fields: map[string]string{"level": "info", "status": "not-a-number"}}
	vv := v.Validate(entry)
	if len(vv) != 1 || vv[0].Field != "status" {
		t.Fatalf("expected type violation for 'status', got %v", vv)
	}
}

func TestValidate_WrongBoolType(t *testing.T) {
	v, _ := New(baseFields())
	entry := parser.Entry{Fields: map[string]string{"level": "info", "ok": "yes"}}
	vv := v.Validate(entry)
	if len(vv) != 1 || vv[0].Field != "ok" {
		t.Fatalf("expected type violation for 'ok', got %v", vv)
	}
}

func TestValidate_OptionalMissingIsOK(t *testing.T) {
	v, _ := New(baseFields())
	entry := parser.Entry{Fields: map[string]string{"level": "warn"}}
	if vv := v.Validate(entry); len(vv) != 0 {
		t.Fatalf("expected no violations for missing optional fields, got %v", vv)
	}
}
