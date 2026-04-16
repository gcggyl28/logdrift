package levelfilter_test

import (
	"testing"

	"github.com/user/logdrift/internal/levelfilter"
)

func TestNew_ValidLevel(t *testing.T) {
	f, err := levelfilter.New("warn", "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := levelfilter.New("verbose", "level")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestNew_EmptyLevelKey(t *testing.T) {
	_, err := levelfilter.New("info", "")
	if err == nil {
		t.Fatal("expected error for empty levelKey")
	}
}

func TestAllow_AboveMinimum(t *testing.T) {
	f, _ := levelfilter.New("warn", "level")
	if !f.Allow(map[string]string{"level": "error"}) {
		t.Error("expected error to pass warn filter")
	}
}

func TestAllow_AtMinimum(t *testing.T) {
	f, _ := levelfilter.New("warn", "level")
	if !f.Allow(map[string]string{"level": "warn"}) {
		t.Error("expected warn to pass warn filter")
	}
}

func TestAllow_BelowMinimum(t *testing.T) {
	f, _ := levelfilter.New("warn", "level")
	if f.Allow(map[string]string{"level": "info"}) {
		t.Error("expected info to be blocked by warn filter")
	}
}

func TestAllow_MissingField_PassThrough(t *testing.T) {
	f, _ := levelfilter.New("error", "level")
	if !f.Allow(map[string]string{"msg": "hello"}) {
		t.Error("expected missing level field to pass through")
	}
}

func TestAllow_UnknownLevelValue_PassThrough(t *testing.T) {
	f, _ := levelfilter.New("warn", "level")
	if !f.Allow(map[string]string{"level": "trace"}) {
		t.Error("expected unknown level value to pass through")
	}
}

func TestAllow_CaseInsensitive(t *testing.T) {
	f, _ := levelfilter.New("info", "level")
	if !f.Allow(map[string]string{"level": "INFO"}) {
		t.Error("expected case-insensitive match")
	}
}
