package rotate_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/logdrift/internal/rotate"
)

func TestNew_InvalidInterval(t *testing.T) {
	f := tempFile(t)
	_, err := rotate.New(f, 0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
	_, err = rotate.New(f, -time.Second)
	if err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := rotate.New("/nonexistent/path/file.log", time.Millisecond)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNew_ValidFile(t *testing.T) {
	f := tempFile(t)
	w, err := rotate.New(f, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Done == nil {
		t.Fatal("Done channel should not be nil")
	}
}

func TestWatch_DetectsTruncation(t *testing.T) {
	f := tempFile(t)
	if err := os.WriteFile(f, []byte("initial content\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := rotate.New(f, 5*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cancel := make(chan struct{})
	defer close(cancel)
	w.Watch(cancel)

	// truncate the file to simulate rotation
	if err := os.WriteFile(f, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case <-w.Done:
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for rotation detection")
	}
}

func TestWatch_CancelStops(t *testing.T) {
	f := tempFile(t)
	w, err := rotate.New(f, 5*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cancel := make(chan struct{})
	w.Watch(cancel)
	close(cancel)

	// Done should NOT be closed because we cancelled without rotation
	select {
	case <-w.Done:
		t.Fatal("Done should not fire when cancelled without rotation")
	case <-time.After(80 * time.Millisecond):
		// expected — no rotation occurred
	}
}

// tempFile creates a temporary file and registers cleanup.
func tempFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "rotate-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}
