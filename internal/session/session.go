// Package session orchestrates tailing, diffing, and rendering for a
// logdrift run. It wires together the tail, diff, and render layers.
package session

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/user/logdrift/internal/config"
	"github.com/user/logdrift/internal/diff"
	"github.com/user/logdrift/internal/render"
	"github.com/user/logdrift/internal/tail"
)

// Session holds the runtime state for a single logdrift invocation.
type Session struct {
	cfg      *config.Config
	pipeline *diff.Pipeline
	renderer render.Renderer
	out      io.Writer
}

// New constructs a Session from the provided config, writing output to out.
func New(cfg *config.Config, out io.Writer) (*Session, error) {
	differ, err := diff.New(cfg.DiffMode)
	if err != nil {
		return nil, fmt.Errorf("session: %w", err)
	}

	r, err := render.New(cfg.OutputFormat)
	if err != nil {
		return nil, fmt.Errorf("session: %w", err)
	}

	return &Session{
		cfg:      cfg,
		pipeline: diff.NewPipeline(differ),
		renderer: r,
		out:      out,
	}, nil
}

// Run starts tailing all configured services, diffs incoming lines, and
// renders any drift until ctx is cancelled.
func (s *Session) Run(ctx context.Context) error {
	var tailers []<-chan tail.Line

	for _, svc := range s.cfg.Services {
		t, err := tail.NewFileTailer(ctx, svc.Name, svc.Path)
		if err != nil {
			return fmt.Errorf("session: service %q: %w", svc.Name, err)
		}
		tailers = append(tailers, t)
	}

	merged := tail.FanIn(ctx, tailers...)

	window := time.Duration(s.cfg.WindowMS) * time.Millisecond

	for result := range s.pipeline.Run(ctx, merged, window) {
		if err := s.renderer.WriteDrift(s.out, result); err != nil {
			return fmt.Errorf("session: render: %w", err)
		}
	}

	return nil
}
