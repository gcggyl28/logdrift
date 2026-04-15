package checkpoint_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/logdrift/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_EmptyPathReturnsError(t *testing.T) {
	_, err := checkpoint.New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_CreatesFileOnFirstSet(t *testing.T) {
	p := tempPath(t)
	s, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set("svc-a", 42); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestGet_UnknownServiceReturnsErrNoCheckpoint(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	_, err := s.Get("missing")
	if !errors.Is(err, checkpoint.ErrNoCheckpoint) {
		t.Fatalf("expected ErrNoCheckpoint, got %v", err)
	}
}

func TestSet_And_Get(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	if err := s.Set("svc-b", 1024); err != nil {
		t.Fatalf("Set: %v", err)
	}
	e, err := s.Get("svc-b")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e.Offset != 1024 {
		t.Errorf("expected offset 1024, got %d", e.Offset)
	}
	if e.Service != "svc-b" {
		t.Errorf("expected service svc-b, got %s", e.Service)
	}
	if e.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestSet_Persists_AcrossReload(t *testing.T) {
	p := tempPath(t)
	s1, _ := checkpoint.New(p)
	_ = s1.Set("svc-c", 9999)

	s2, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, err := s2.Get("svc-c")
	if err != nil {
		t.Fatalf("Get after reload: %v", err)
	}
	if e.Offset != 9999 {
		t.Errorf("expected 9999, got %d", e.Offset)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	p := tempPath(t)
	s, _ := checkpoint.New(p)
	_ = s.Set("svc-d", 7)
	if err := s.Delete("svc-d"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Get("svc-d")
	if !errors.Is(err, checkpoint.ErrNoCheckpoint) {
		t.Errorf("expected ErrNoCheckpoint after delete, got %v", err)
	}
}

func TestDelete_Persists_AcrossReload(t *testing.T) {
	p := tempPath(t)
	s1, _ := checkpoint.New(p)
	_ = s1.Set("svc-e", 3)
	_ = s1.Delete("svc-e")

	s2, _ := checkpoint.New(p)
	_, err := s2.Get("svc-e")
	if !errors.Is(err, checkpoint.ErrNoCheckpoint) {
		t.Errorf("expected ErrNoCheckpoint after reload, got %v", err)
	}
}
