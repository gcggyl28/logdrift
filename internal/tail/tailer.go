package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Line represents a single log line emitted from a named service.
type Line struct {
	Service   string
	Text      string
	Timestamp time.Time
}

// Tailer tails a log source and emits lines on a channel.
type Tailer struct {
	Service string
	Reader  io.ReadCloser
}

// NewFileTailer opens the given file path and returns a Tailer for it.
func NewFileTailer(service, path string) (*Tailer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// Seek to end so we only tail new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return nil, err
	}
	return &Tailer{Service: service, Reader: f}, nil
}

// NewReaderTailer wraps an arbitrary io.ReadCloser (useful for testing / stdin).
func NewReaderTailer(service string, r io.ReadCloser) *Tailer {
	return &Tailer{Service: service, Reader: r}
}

// Tail reads lines from the underlying reader and sends them to out until ctx
// is cancelled or the reader returns a permanent error.
func (t *Tailer) Tail(ctx context.Context, out chan<- Line) error {
	defer t.Reader.Close()
	scanner := bufio.NewScanner(t.Reader)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if scanner.Scan() {
			out <- Line{
				Service:   t.Service,
				Text:      scanner.Text(),
				Timestamp: time.Now(),
			}
			continue
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		// EOF — poll briefly before retrying (file may grow).
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
		scanner = bufio.NewScanner(t.Reader)
	}
}
