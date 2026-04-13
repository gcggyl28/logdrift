package session_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/logdrift/internal/config"
	"github.com/user/logdrift/internal/session"
)

func baseConfig() *config.Config {
	return &config.Config{
		DiffMode:     "unified",
		OutputFormat: "plain",
		WindowMS:     50,
	}
}

func TestNew_ValidConfig(t *testing.T) {
	cfg := baseConfig()
	var buf bytes.Buffer
	s, err := session.New(cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil session")
	}
}

func TestNew_InvalidDiffMode(t *testing.T) {
	cfg := baseConfig()
	cfg.DiffMode = "bogus"
	_, err := session.New(cfg, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for invalid diff mode")
	}
}

func TestNew_InvalidOutputFormat(t *testing.T) {
	cfg := baseConfig()
	cfg.OutputFormat = "xml"
	_, err := session.New(cfg, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for invalid output format")
	}
}

func TestRun_NoServices_ReturnsQuickly(t *testing.T) {
	cfg := baseConfig()
	cfg.Services = nil

	var buf bytes.Buffer
	s, err := session.New(cfg, &buf)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if err := s.Run(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_BadServicePath_ReturnsError(t *testing.T) {
	cfg := baseConfig()
	cfg.Services = []config.Service{
		{Name: "svc", Path: "/nonexistent/path/to/log.log"},
	}

	var buf bytes.Buffer
	s, err := session.New(cfg, &buf)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	ctx := context.Background()
	err = s.Run(ctx)
	if err == nil || !strings.Contains(err.Error(), "svc") {
		t.Fatalf("expected service error, got: %v", err)
	}
}
