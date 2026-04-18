package tagfilter

import (
	"testing"
)

func TestNew_UnknownModeReturnsError(t *testing.T) {
	_, err := New("bad", map[string]string{"env": "prod"})
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNew_EmptyTagsReturnsError(t *testing.T) {
	_, err := New(ModeAny, map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty tags")
	}
}

func TestNew_EmptyKeyReturnsError(t *testing.T) {
	_, err := New(ModeAny, map[string]string{"": "prod"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNew_EmptyValueReturnsError(t *testing.T) {
	_, err := New(ModeAny, map[string]string{"env": ""})
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestNew_Valid(t *testing.T) {
	f, err := New(ModeAll, map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestAllow_AnyMode_MatchesOneOfMany(t *testing.T) {
	f, _ := New(ModeAny, map[string]string{"env": "prod", "region": "us-east"})
	e := Entry{Tags: map[string]string{"env": "prod", "region": "eu-west"}}
	if !f.Allow(e) {
		t.Error("expected allow when at least one tag matches in any mode")
	}
}

func TestAllow_AnyMode_NoMatch(t *testing.T) {
	f, _ := New(ModeAny, map[string]string{"env": "prod"})
	e := Entry{Tags: map[string]string{"env": "staging"}}
	if f.Allow(e) {
		t.Error("expected deny when no tags match in any mode")
	}
}

func TestAllow_AllMode_RequiresAllTags(t *testing.T) {
	f, _ := New(ModeAll, map[string]string{"env": "prod", "region": "us-east"})
	e := Entry{Tags: map[string]string{"env": "prod", "region": "us-east"}}
	if !f.Allow(e) {
		t.Error("expected allow when all tags match")
	}
}

func TestAllow_AllMode_PartialMatchDenied(t *testing.T) {
	f, _ := New(ModeAll, map[string]string{"env": "prod", "region": "us-east"})
	e := Entry{Tags: map[string]string{"env": "prod"}}
	if f.Allow(e) {
		t.Error("expected deny when only partial tags match in all mode")
	}
}

func TestAllow_NilTags_Denied(t *testing.T) {
	f, _ := New(ModeAny, map[string]string{"env": "prod"})
	e := Entry{Tags: nil}
	if f.Allow(e) {
		t.Error("expected deny for entry with nil tags")
	}
}
