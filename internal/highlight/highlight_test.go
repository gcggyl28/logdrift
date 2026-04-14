package highlight

import (
	"strings"
	"testing"
)

func TestNew_ValidANSI(t *testing.T) {
	h, err := New([]string{`ERROR`, `WARN`}, StyleANSI, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Highlighter")
	}
}

func TestNew_ValidBraces(t *testing.T) {
	h, err := New([]string{`(?i)error`}, StyleBraces, "")
	if err != nil || h == nil {
		t.Fatalf("unexpected error or nil highlighter: %v", err)
	}
}

func TestNew_InvalidStyle(t *testing.T) {
	_, err := New([]string{`foo`}, Style("neon"), "")
	if err == nil {
		t.Fatal("expected error for unknown style")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]string{`[invalid`}, StyleANSI, "")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New([]string{""}, StyleANSI, "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestApply_NoPatterns_Passthrough(t *testing.T) {
	h, _ := New([]string{}, StyleANSI, "")
	line := "nothing to highlight"
	if got := h.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_BracesWrapsMatch(t *testing.T) {
	h, _ := New([]string{`ERROR`}, StyleBraces, "")
	got := h.Apply("2024/01/01 ERROR something broke")
	if !strings.Contains(got, "[[ERROR]]") {
		t.Errorf("expected [[ERROR]] in output, got %q", got)
	}
}

func TestApply_ANSIWrapsMatch(t *testing.T) {
	h, _ := New([]string{`WARN`}, StyleANSI, "\033[31m")
	got := h.Apply("WARN disk usage high")
	if !strings.Contains(got, "\033[31m") {
		t.Errorf("expected ANSI code in output, got %q", got)
	}
	if !strings.Contains(got, "\033[0m") {
		t.Errorf("expected ANSI reset in output, got %q", got)
	}
}

func TestApply_MultiplePatterns(t *testing.T) {
	h, _ := New([]string{`ERROR`, `WARN`}, StyleBraces, "")
	got := h.Apply("ERROR and WARN in one line")
	if !strings.Contains(got, "[[ERROR]]") || !strings.Contains(got, "[[WARN]]") {
		t.Errorf("expected both patterns highlighted, got %q", got)
	}
}

func TestApply_DefaultANSIColor(t *testing.T) {
	h, _ := New([]string{`INFO`}, StyleANSI, "") // empty color → default yellow
	got := h.Apply("INFO starting up")
	if !strings.Contains(got, "\033[33m") {
		t.Errorf("expected default yellow ANSI code, got %q", got)
	}
}
