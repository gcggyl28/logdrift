package filter_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/filter"
)

func TestNew_InvalidMode(t *testing.T) {
	_, err := filter.New([]string{"foo"}, "bad", false)
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := filter.New([]string{"[invalid"}, filter.ModeAny, false)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestMatch_NoPatterns(t *testing.T) {
	f, err := filter.New(nil, filter.ModeAny, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match("anything") {
		t.Error("empty filter should match everything")
	}
}

func TestMatch_AnyMode(t *testing.T) {
	f, _ := filter.New([]string{"error", "warn"}, filter.ModeAny, false)
	if !f.Match("ERROR: disk full") {
		t.Error("expected match for 'error' pattern")
	}
	if f.Match("INFO: all good") {
		t.Error("expected no match for info line")
	}
}

func TestMatch_AllMode(t *testing.T) {
	f, _ := filter.New([]string{"error", "disk"}, filter.ModeAll, false)
	if !f.Match("error: disk full") {
		t.Error("expected match when all patterns present")
	}
	if f.Match("error: network timeout") {
		t.Error("expected no match when only one pattern present")
	}
}

func TestMatch_Invert(t *testing.T) {
	f, _ := filter.New([]string{"debug"}, filter.ModeAny, true)
	if f.Match("debug: verbose output") {
		t.Error("inverted filter should reject debug lines")
	}
	if !f.Match("info: service started") {
		t.Error("inverted filter should pass non-debug lines")
	}
}

func TestMatch_CaseInsensitivePattern(t *testing.T) {
	f, _ := filter.New([]string{"(?i)error"}, filter.ModeAny, false)
	if !f.Match("ERROR: something bad") {
		t.Error("expected case-insensitive match")
	}
}
