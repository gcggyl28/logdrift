package pattern_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/pattern"
)

func TestNew_ValidAny(t *testing.T) {
	m, err := pattern.New(pattern.MatchAny, []string{`error`, `warn`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Mode() != pattern.MatchAny {
		t.Errorf("mode = %q, want %q", m.Mode(), pattern.MatchAny)
	}
	if m.Len() != 2 {
		t.Errorf("len = %d, want 2", m.Len())
	}
}

func TestNew_ValidAll(t *testing.T) {
	m, err := pattern.New(pattern.MatchAll, []string{`\d+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Mode() != pattern.MatchAll {
		t.Errorf("mode = %q, want %q", m.Mode(), pattern.MatchAll)
	}
}

func TestNew_InvalidMode(t *testing.T) {
	_, err := pattern.New("none", nil)
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := pattern.New(pattern.MatchAny, []string{`[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := pattern.New(pattern.MatchAny, []string{""})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestMatch_NoPatterns_AlwaysTrue(t *testing.T) {
	for _, mode := range []pattern.MatchMode{pattern.MatchAny, pattern.MatchAll} {
		m, _ := pattern.New(mode, nil)
		if !m.Match("anything") {
			t.Errorf("mode %q with no patterns should match everything", mode)
		}
	}
}

func TestMatch_AnyMode(t *testing.T) {
	m, _ := pattern.New(pattern.MatchAny, []string{`error`, `panic`})

	if !m.Match("error: disk full") {
		t.Error("expected match on 'error'")
	}
	if !m.Match("panic: nil pointer") {
		t.Error("expected match on 'panic'")
	}
	if m.Match("info: all good") {
		t.Error("expected no match")
	}
}

func TestMatch_AllMode(t *testing.T) {
	m, _ := pattern.New(pattern.MatchAll, []string{`error`, `disk`})

	if !m.Match("error: disk full") {
		t.Error("expected match when both patterns present")
	}
	if m.Match("error: network timeout") {
		t.Error("expected no match when only one pattern present")
	}
	if m.Match("disk full warning") {
		t.Error("expected no match when only second pattern present")
	}
}

func TestMatch_CaseInsensitive(t *testing.T) {
	m, _ := pattern.New(pattern.MatchAny, []string{`(?i)ERROR`})
	if !m.Match("error: something failed") {
		t.Error("expected case-insensitive match")
	}
}
