package timefield_test

import (
	"testing"
	"time"

	"github.com/angryboat/logdrift/internal/parser"
	"github.com/angryboat/logdrift/internal/timefield"
)

func entry(fields map[string]string) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestNew_EmptyFieldReturnsError(t *testing.T) {
	_, err := timefield.New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptyLayoutsReturnsError(t *testing.T) {
	_, err := timefield.New("ts", timefield.WithLayouts())
	if err == nil {
		t.Fatal("expected error for empty layouts")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := timefield.New("ts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExtract_RFC3339(t *testing.T) {
	ext, _ := timefield.New("ts")
	now := time.Now().UTC().Truncate(time.Second)
	e := entry(map[string]string{"ts": now.Format(time.RFC3339)})
	got, err := ext.Extract(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(now) {
		t.Fatalf("expected %v got %v", now, got)
	}
}

func TestExtract_MissingField_NoFallback(t *testing.T) {
	ext, _ := timefield.New("ts")
	_, err := ext.Extract(entry(nil))
	if err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestExtract_MissingField_WithFallback(t *testing.T) {
	ext, _ := timefield.New("ts", timefield.WithFallback())
	before := time.Now()
	got, err := ext.Extract(entry(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Before(before) {
		t.Fatal("fallback time should be >= before")
	}
}

func TestExtract_UnparsableValue_NoFallback(t *testing.T) {
	ext, _ := timefield.New("ts")
	_, err := ext.Extract(entry(map[string]string{"ts": "not-a-time"}))
	if err == nil {
		t.Fatal("expected error for unparseable value")
	}
}

func TestExtract_CustomLayout(t *testing.T) {
	const layout = "02/Jan/2006"
	ext, _ := timefield.New("ts", timefield.WithLayouts(layout))
	e := entry(map[string]string{"ts": "15/Jun/2024"})
	_, err := ext.Extract(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
