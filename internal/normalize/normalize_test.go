package normalize

import (
	"testing"
)

func TestNew_ValidModes(t *testing.T) {
	for _, mode := range []Mode{ModeNone, ModeTimestamp, ModePID, ModeFull} {
		_, err := New(mode, nil, "")
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", mode, err)
		}
	}
}

func TestNew_InvalidMode(t *testing.T) {
	_, err := New("bogus", nil, "")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNew_EmptyExtraPattern(t *testing.T) {
	_, err := New(ModeNone, []string{"   "}, "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidExtraPattern(t *testing.T) {
	_, err := New(ModeNone, []string{"[invalid"}, "")
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestApply_None_Passthrough(t *testing.T) {
	n, _ := New(ModeNone, nil, "")
	input := "2024-01-02T15:04:05Z pid=42 hello"
	if got := n.Apply(input); got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_Timestamp(t *testing.T) {
	n, _ := New(ModeTimestamp, nil, "<T>")
	got := n.Apply("2024-01-02T15:04:05Z service started")
	expected := "<T> service started"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestApply_PID(t *testing.T) {
	n, _ := New(ModePID, nil, "<P>")
	got := n.Apply("process pid=1234 exited")
	expected := "process <P> exited"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestApply_Full_StripsBoth(t *testing.T) {
	n, _ := New(ModeFull, nil, "_")
	got := n.Apply("2024-06-01T12:00:00Z pid=99 req=550e8400-e29b-41d4-a716-446655440000 ok")
	for _, bad := range []string{"2024", "pid=99", "550e8400"} {
		if contains(got, bad) {
			t.Errorf("output still contains %q: %q", bad, got)
		}
	}
}

func TestApply_ExtraPattern(t *testing.T) {
	n, _ := New(ModeNone, []string{`\d+ms`}, "<DUR>")
	got := n.Apply("request took 350ms total")
	expected := "request took <DUR> total"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestApply_DefaultPlaceholder(t *testing.T) {
	n, _ := New(ModeTimestamp, nil, "")
	got := n.Apply("2024-01-01T00:00:00Z boot")
	if !contains(got, "<?>") {
		t.Errorf("expected default placeholder in output, got %q", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
