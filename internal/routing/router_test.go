package routing

import (
	"testing"
)

func TestNew_EmptyRulesReturnsError(t *testing.T) {
	_, err := New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty rules, got nil")
	}
}

func TestNew_InvalidPatternReturnsError(t *testing.T) {
	_, err := New(map[string]string{"dest": "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestNew_EmptyPatternReturnsError(t *testing.T) {
	_, err := New(map[string]string{"dest": ""})
	if err == nil {
		t.Fatal("expected error for empty pattern, got nil")
	}
}

func TestNew_ValidRules(t *testing.T) {
	r, err := New(map[string]string{"groupA": "^api-", "groupB": "^db-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Routes()) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(r.Routes()))
	}
}

func TestMatch_NoMatch(t *testing.T) {
	r, _ := New(map[string]string{"groupA": "^api-"})
	dests := r.Match("worker-1")
	if len(dests) != 0 {
		t.Fatalf("expected no destinations, got %v", dests)
	}
}

func TestMatch_SingleMatch(t *testing.T) {
	r, _ := New(map[string]string{"groupA": "^api-", "groupB": "^db-"})
	dests := r.Match("api-gateway")
	if len(dests) != 1 || dests[0] != "groupA" {
		t.Fatalf("expected [groupA], got %v", dests)
	}
}

func TestMatch_MultipleMatches(t *testing.T) {
	r, _ := New(map[string]string{
		"all":  ".*",
		"apis": "^api-",
	})
	dests := r.Match("api-users")
	if len(dests) != 2 {
		t.Fatalf("expected 2 destinations, got %v", dests)
	}
}

func TestRoutes_ReturnsCopy(t *testing.T) {
	r, _ := New(map[string]string{"g": "svc"})
	a := r.Routes()
	b := r.Routes()
	if &a[0] == &b[0] {
		t.Fatal("Routes() must return independent copies")
	}
}
