package tail

import (
	"context"
	"sync"
)

// FanIn merges multiple Tailer goroutines into a single Line channel.
// Each tailer is started in its own goroutine; the returned channel is closed
// once all tailers finish (or ctx is cancelled).
func FanIn(ctx context.Context, tailers []*Tailer) <-chan Line {
	out := make(chan Line, len(tailers)*16)
	var wg sync.WaitGroup

	for _, t := range tailers {
		wg.Add(1)
		go func(tr *Tailer) {
			defer wg.Done()
			// Errors are silently swallowed here; callers can inspect
			// individual tailer errors by wrapping Tail themselves.
			_ = tr.Tail(ctx, out)
		}(t)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
