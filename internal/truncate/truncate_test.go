package truncate

import (
	"strings"
	"testing"
)

func TestNew_ZeroDisablesTruncation(t *testing.T) {
	tr, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Enabled() {
		t.Fatal("expected truncation to be disabled for maxLen=0")
	}
}

func TestNew_NegativeReturnsError(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative maxLen")
	}
}

func TestNew_SuffixTooLong(t *testing.T) {
	_, err := New(3, WithSuffix("....."))
	if err == nil {
		t.Fatal("expected error when suffix >= maxLen")
	}
}

func TestNew_ValidPositive(t *testing.T) {
	tr, err := New(80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tr.Enabled() {
		t.Fatal("expected truncation to be enabled")
	}
}

func TestApply_Disabled_Passthrough(t *testing.T) {
	tr, _ := New(0)
	long := strings.Repeat("x", 200)
	if got := tr.Apply(long); got != long {
		t.Fatalf("disabled truncator should not modify line")
	}
}

func TestApply_ShortLine_Unchanged(t *testing.T) {
	tr, _ := New(100)
	line := "short line"
	if got := tr.Apply(line); got != line {
		t.Fatalf("short line should not be truncated, got %q", got)
	}
}

func TestApply_ExactLength_Unchanged(t *testing.T) {
	tr, _ := New(10)
	line := "1234567890" // exactly 10 bytes
	if got := tr.Apply(line); got != line {
		t.Fatalf("exact-length line should not be truncated, got %q", got)
	}
}

func TestApply_LongLine_Clipped(t *testing.T) {
	tr, _ := New(10)
	line := "1234567890EXTRA"
	got := tr.Apply(line)
	if len(got) != 10 {
		t.Fatalf("expected length 10, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(got, defaultSuffix) {
		t.Fatalf("expected suffix %q in %q", defaultSuffix, got)
	}
}

func TestApply_CustomSuffix(t *testing.T) {
	tr, _ := New(20, WithSuffix("[cut]"))
	line := strings.Repeat("a", 30)
	got := tr.Apply(line)
	if !strings.HasSuffix(got, "[cut]") {
		t.Fatalf("expected custom suffix in %q", got)
	}
	if len(got) != 20 {
		t.Fatalf("expected length 20, got %d", len(got))
	}
}
