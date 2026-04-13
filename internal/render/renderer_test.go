package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_ValidFormats(t *testing.T) {
	formats := []Format{FormatPlain, FormatColored, FormatTimestamp}
	for _, f := range formats {
		r, err := New(f)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", f, err)
		}
		if r == nil {
			t.Errorf("New(%q) returned nil renderer", f)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("neon")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
	if !strings.Contains(err.Error(), "neon") {
		t.Errorf("error message should mention format name, got: %v", err)
	}
}

func TestWriteDrift_Plain(t *testing.T) {
	r, _ := New(FormatPlain)
	var buf bytes.Buffer
	r.SetOutput(&buf)

	r.WriteDrift("svc-a", "svc-b", "+added line\n-removed line")

	out := buf.String()
	if !strings.Contains(out, "--- svc-a") {
		t.Errorf("expected service A header, got: %s", out)
	}
	if !strings.Contains(out, "+++ svc-b") {
		t.Errorf("expected service B header, got: %s", out)
	}
	if !strings.Contains(out, "+added line") {
		t.Errorf("expected delta content, got: %s", out)
	}
}

func TestWriteDrift_Timestamp(t *testing.T) {
	r, _ := New(FormatTimestamp)
	var buf bytes.Buffer
	r.SetOutput(&buf)

	r.WriteDrift("alpha", "beta", "+new")

	out := buf.String()
	// RFC3339 timestamps contain 'T' and 'Z'
	if !strings.Contains(out, "T") || !strings.Contains(out, "Z") {
		t.Errorf("expected RFC3339 timestamp in output, got: %s", out)
	}
}

func TestWriteDrift_Colored_ContainsDelta(t *testing.T) {
	r, _ := New(FormatColored)
	var buf bytes.Buffer
	r.SetOutput(&buf)

	delta := "+line added\n-line removed\n context"
	r.WriteDrift("x", "y", delta)

	out := buf.String()
	// Color codes wrap the text; raw strings may not appear verbatim,
	// but the plain words should still be present.
	if !strings.Contains(out, "line added") {
		t.Errorf("expected 'line added' in colored output, got: %s", out)
	}
	if !strings.Contains(out, "line removed") {
		t.Errorf("expected 'line removed' in colored output, got: %s", out)
	}
}
