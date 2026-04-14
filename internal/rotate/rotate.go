// Package rotate provides log rotation detection for file-based tailers.
// It watches for file truncation or inode changes and signals when a
// tailer should reopen its source file.
package rotate

import (
	"errors"
	"os"
	"time"
)

// Watcher monitors a file path for rotation events (truncation or inode
// replacement). When a rotation is detected the Done channel is closed.
type Watcher struct {
	path     string
	interval time.Duration
	initSize int64
	initIno  uint64
	Done     chan struct{}
}

// New creates a Watcher for the given file path, polling at the specified
// interval. interval must be positive. The file must exist at creation time.
func New(path string, interval time.Duration) (*Watcher, error) {
	if interval <= 0 {
		return nil, errors.New("rotate: interval must be positive")
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		path:     path,
		interval: interval,
		initSize: info.Size(),
		initIno:  inode(info),
		Done:     make(chan struct{}),
	}
	return w, nil
}

// Watch begins polling in the background. cancel should be closed to stop the
// watcher when no rotation is needed. Watch returns immediately.
func (w *Watcher) Watch(cancel <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-cancel:
				return
			case <-ticker.C:
				if w.rotated() {
					close(w.Done)
					return
				}
			}
		}
	}()
}

func (w *Watcher) rotated() bool {
	info, err := os.Stat(w.path)
	if err != nil {
		// file disappeared — treat as rotation
		return true
	}
	if inode(info) != w.initIno {
		return true
	}
	if info.Size() < w.initSize {
		return true
	}
	w.initSize = info.Size()
	return false
}
