package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "logdrift.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	raw := `
services:
  - name: api
    command: "tail -f /var/log/api.log"
    color: cyan
  - name: worker
    command: "tail -f /var/log/worker.log"
diff_mode: inline
buffer_size: 100
`
	cfg, err := Load(writeTemp(t, raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(cfg.Services))
	}
	if cfg.DiffMode != "inline" {
		t.Errorf("expected diff_mode inline, got %q", cfg.DiffMode)
	}
	if cfg.BufferSize != 100 {
		t.Errorf("expected buffer_size 100, got %d", cfg.BufferSize)
	}
}

func TestLoad_Defaults(t *testing.T) {
	raw := `
services:
  - name: svc
    command: echo hello
`
	cfg, err := Load(writeTemp(t, raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DiffMode != "unified" {
		t.Errorf("expected default diff_mode unified, got %q", cfg.DiffMode)
	}
	if cfg.BufferSize != 200 {
		t.Errorf("expected default buffer_size 200, got %d", cfg.BufferSize)
	}
}

func TestLoad_NoServices(t *testing.T) {
	raw := `diff_mode: unified\nbuffer_size: 50\n`
	_, err := Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for missing services")
	}
}

func TestLoad_InvalidDiffMode(t *testing.T) {
	raw := `
services:
  - name: svc
    command: echo hi
diff_mode: fancy
`
	_, err := Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for invalid diff_mode")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/logdrift.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
