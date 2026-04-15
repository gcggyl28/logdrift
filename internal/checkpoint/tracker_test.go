package checkpoint_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/logdrift/internal/checkpoint"
)

func TestNewTracker_NilStoreReturnsError(t *testing.T) {
	_, err := checkpoint.NewTracker(nil, "svc", 0)
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestNewTracker_EmptyServiceReturnsError(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	_, err := checkpoint.NewTracker(s, "", 0)
	if err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestNewTracker_RestoresOffset(t *testing.T) {
	p := tempPath(t)
	s, _ := checkpoint.New(p)
	_ = s.Set("svc-x", 512)

	tr, err := checkpoint.NewTracker(s, "svc-x", 0)
	if err != nil {
		t.Fatalf("NewTracker: %v", err)
	}
	if tr.Offset() != 512 {
		t.Errorf("expected restored offset 512, got %d", tr.Offset())
	}
}

func TestAdvance_SynchronousMode(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	tr, _ := checkpoint.NewTracker(s, "svc-y", 0)

	if err := tr.Advance(100); err != nil {
		t.Fatalf("Advance: %v", err)
	}
	if tr.Offset() != 100 {
		t.Errorf("expected 100, got %d", tr.Offset())
	}
	e, _ := s.Get("svc-y")
	if e.Offset != 100 {
		t.Errorf("expected persisted offset 100, got %d", e.Offset)
	}
}

func TestAdvance_AsyncMode_DoesNotPersistImmediately(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	tr, _ := checkpoint.NewTracker(s, "svc-z", time.Hour)

	_ = tr.Advance(200)
	if tr.Offset() != 200 {
		t.Errorf("expected in-memory offset 200, got %d", tr.Offset())
	}
	_, err := s.Get("svc-z")
	if err == nil {
		t.Error("expected no persisted entry yet in async mode")
	}
}

func TestRun_FinalFlushOnCancel(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	tr, _ := checkpoint.NewTracker(s, "svc-flush", time.Hour)
	_ = tr.Advance(777)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- tr.Run(ctx) }()

	cancel()
	<-done

	e, err := s.Get("svc-flush")
	if err != nil {
		t.Fatalf("Get after cancel: %v", err)
	}
	if e.Offset != 777 {
		t.Errorf("expected final flush offset 777, got %d", e.Offset)
	}
}
