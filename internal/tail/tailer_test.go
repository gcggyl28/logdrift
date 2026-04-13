package tail_test

import (
	"context"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/tail"
)

// nopCloser wraps a Reader with a no-op Close.
type nopCloser struct{ io.Reader }

func (nopCloser) Close() error { return nil }

func TestTail_EmitsLines(t *testing.T) {
	input := "line one\nline two\nline three\n"
	rc := nopCloser{strings.NewReader(input)}
	tr := tail.NewReaderTailer("svc-a", rc)

	out := make(chan tail.Line, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = tr.Tail(ctx, out)
	}()

	wg.Wait()
	close(out)

	var got []string
	for l := range out {
		got = append(got, l.Text)
		if l.Service != "svc-a" {
			t.Errorf("unexpected service name %q", l.Service)
		}
		if l.Timestamp.IsZero() {
			t.Error("timestamp should not be zero")
		}
	}

	want := []string{"line one", "line two", "line three"}
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("line %d: got %q, want %q", i, got[i], w)
		}
	}
}

func TestTail_CancelStops(t *testing.T) {
	// Use a pipe so the reader blocks indefinitely.
	pr, pw := io.Pipe()
	tr := tail.NewReaderTailer("svc-b", pr)

	out := make(chan tail.Line, 4)
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- tr.Tail(ctx, out)
	}()

	// Write one line then cancel.
	_, _ = pw.Write([]byte("hello\n"))
	time.Sleep(50 * time.Millisecond)
	cancel()
	pw.Close()

	select {
	case err := <-errCh:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Tail did not stop after context cancel")
	}
}
