package multiline

import (
	"strings"
	"testing"
	"time"
)

func TestNew_ValidPrefix(t *testing.T) {
	_, err := New(Config{Mode: ModePrefix, Pattern: `^\d{4}-`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidMode(t *testing.T) {
	_, err := New(Config{Mode: "bad", Pattern: `^\d`})
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(Config{Mode: ModePrefix, Pattern: ""})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Config{Mode: ModePrefix, Pattern: `[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_NegativeMaxLines(t *testing.T) {
	_, err := New(Config{Mode: ModePrefix, Pattern: `^\d`, MaxLines: -1})
	if err == nil {
		t.Fatal("expected error for negative max_lines")
	}
}

func TestPush_PrefixMode_BuffersUntilNextStart(t *testing.T) {
	c, _ := New(Config{Mode: ModePrefix, Pattern: `^START`})

	if ev, ok := c.Push("START first"); ok {
		t.Fatalf("unexpected emit on first line: %q", ev)
	}
	c.Push("  continuation 1")
	c.Push("  continuation 2")

	ev, ok := c.Push("START second")
	if !ok {
		t.Fatal("expected event when new start line arrives")
	}
	if !strings.Contains(ev, "START first") || !strings.Contains(ev, "continuation 2") {
		t.Errorf("event missing lines: %q", ev)
	}
}

func TestPush_ContinuationMode_EmitsOnNonMatch(t *testing.T) {
	c, _ := New(Config{Mode: ModeContinuation, Pattern: `^\s+`})

	c.Push("  line 1")
	c.Push("  line 2")

	ev, ok := c.Push("new event")
	if !ok {
		t.Fatal("expected flush when non-continuation arrives")
	}
	if !strings.Contains(ev, "line 1") {
		t.Errorf("event missing buffered lines: %q", ev)
	}
}

func TestPush_MaxLines_ForcesFlush(t *testing.T) {
	c, _ := New(Config{Mode: ModeContinuation, Pattern: `^\s+`, MaxLines: 3})

	c.Push("  a")
	c.Push("  b")
	ev, ok := c.Push("  c")
	if !ok {
		t.Fatal("expected flush at max_lines")
	}
	lines := strings.Split(ev, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestFlush_EmptyBufferReturnsFalse(t *testing.T) {
	c, _ := New(Config{Mode: ModePrefix, Pattern: `^START`})
	_, ok := c.Flush()
	if ok {
		t.Fatal("flush of empty buffer should return false")
	}
}

func TestFlush_UpdatesLastFlush(t *testing.T) {
	c, _ := New(Config{Mode: ModePrefix, Pattern: `^START`, Timeout: 100 * time.Millisecond})
	c.Push("START foo")
	before := c.LastFlush
	time.Sleep(5 * time.Millisecond)
	c.Flush()
	if !c.LastFlush.After(before) {
		t.Error("LastFlush should be updated after flush")
	}
}
