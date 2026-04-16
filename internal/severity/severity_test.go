package severity_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/parser"
	"github.com/yourorg/logdrift/internal/severity"
)

func entry(fields map[string]any) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestNew_NoFieldsReturnsError(t *testing.T) {
	_, err := severity.New()
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_EmptyFieldNameReturnsError(t *testing.T) {
	_, err := severity.New("level", "")
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := severity.New("level", "severity")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClassify_KnownLevels(t *testing.T) {
	c, _ := severity.New("level")
	cases := []struct {
		val  string
		want severity.Level
	}{
		{"debug", severity.LevelDebug},
		{"trace", severity.LevelDebug},
		{"info", severity.LevelInfo},
		{"information", severity.LevelInfo},
		{"warn", severity.LevelWarn},
		{"warning", severity.LevelWarn},
		{"error", severity.LevelError},
		{"err", severity.LevelError},
		{"fatal", severity.LevelFatal},
		{"critical", severity.LevelFatal},
		{"panic", severity.LevelFatal},
		{"bogus", severity.LevelUnknown},
	}
	for _, tc := range cases {
		got := c.Classify(entry(map[string]any{"level": tc.val}))
		if got != tc.want {
			t.Errorf("val=%q: got %v, want %v", tc.val, got, tc.want)
		}
	}
}

func TestClassify_FallsBackToSecondField(t *testing.T) {
	c, _ := severity.New("level", "severity")
	e := entry(map[string]any{"severity": "warn"})
	if got := c.Classify(e); got != severity.LevelWarn {
		t.Fatalf("expected warn, got %v", got)
	}
}

func TestClassify_MissingFieldReturnsUnknown(t *testing.T) {
	c, _ := severity.New("level")
	if got := c.Classify(entry(map[string]any{})); got != severity.LevelUnknown {
		t.Fatalf("expected unknown, got %v", got)
	}
}

func TestLevel_String(t *testing.T) {
	if severity.LevelError.String() != "error" {
		t.Fatalf("unexpected string for LevelError")
	}
	if severity.LevelUnknown.String() != "unknown" {
		t.Fatalf("unexpected string for LevelUnknown")
	}
}
