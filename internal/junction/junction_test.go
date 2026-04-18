package junction_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/user/logdrift/internal/diff"
	"github.com/user/logdrift/internal/junction"
)

func makeLineCh(lines []string) <-chan diff.Line {
	ch := make(chan diff.Line, len(lines))
	for _, l := range lines {
		ch <- diff.Line{Text: l}
	}
	close(ch)
	return ch
}

func TestNew_EmptySourcesReturnsError(t *testing.T) {
	_, err := junction.New(map[string]<-chan diff.Line{})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestNew_ValidSources(t *testing.T) {
	src := map[string]<-chan diff.Line{"svc": makeLineCh(nil)}
	_, err := junction.New(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_EmitsAllEntries(t *testing.T) {
	src := map[string]<-chan diff.Line{
		"alpha": makeLineCh([]string{"a1", "a2"}),
		"beta":  makeLineCh([]string{"b1"}),
	}
	j, err := junction.New(src)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out := j.Run(ctx)
	var got []junction.Entry
	for e := range out {
		got = append(got, e)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestRun_TagsServiceCorrectly(t *testing.T) {
	src := map[string]<-chan diff.Line{
		"myservice": makeLineCh([]string{"hello"}),
	}
	j, _ := junction.New(src)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	out := j.Run(ctx)
	e := <-out
	if e.Service != "myservice" {
		t.Errorf("expected service 'myservice', got %q", e.Service)
	}
	if e.Line != "hello" {
		t.Errorf("expected line 'hello', got %q", e.Line)
	}
}

func TestRun_CancelStopsEarly(t *testing.T) {
	blocking := make(chan diff.Line) // never sends
	src := map[string]<-chan diff.Line{"svc": blocking}
	j, _ := junction.New(src)

	ctx, cancel := context.WithCancel(context.Background())
	out := j.Run(ctx)
	cancel()

	select {
	case _, ok := <-out:
		if ok {
			// drain is fine
		}
	case <-time.After(time.Second):
		t.Fatal("output channel not closed after cancel")
	}
	_ = sort.Search // suppress import
}
