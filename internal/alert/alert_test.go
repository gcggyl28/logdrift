package alert_test

import (
	"testing"
	"time"

	"github.com/user/logdrift/internal/alert"
)

func validConfig() alert.Config {
	return alert.Config{
		WarnThreshold:  5,
		ErrorThreshold: 10,
		Window:         30 * time.Second,
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := alert.New(validConfig())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNew_InvalidWarnThreshold(t *testing.T) {
	cfg := validConfig()
	cfg.WarnThreshold = 0
	_, err := alert.New(cfg)
	if err == nil {
		t.Fatal("expected error for zero WarnThreshold")
	}
}

func TestNew_ErrorThresholdLessThanWarn(t *testing.T) {
	cfg := validConfig()
	cfg.ErrorThreshold = 3
	_, err := alert.New(cfg)
	if err == nil {
		t.Fatal("expected error when ErrorThreshold < WarnThreshold")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	cfg := validConfig()
	cfg.Window = 0
	_, err := alert.New(cfg)
	if err == nil {
		t.Fatal("expected error for zero Window")
	}
}

func TestEvaluate_NoAlert(t *testing.T) {
	a, _ := alert.New(validConfig())
	if got := a.Evaluate("svc", 3); got != nil {
		t.Fatalf("expected nil alert, got %+v", got)
	}
}

func TestEvaluate_WarnLevel(t *testing.T) {
	a, _ := alert.New(validConfig())
	got := a.Evaluate("svc", 5)
	if got == nil {
		t.Fatal("expected alert, got nil")
	}
	if got.Level != alert.LevelWarn {
		t.Fatalf("expected warn level, got %s", got.Level)
	}
	if got.Service != "svc" {
		t.Fatalf("unexpected service %q", got.Service)
	}
}

func TestEvaluate_ErrorLevel(t *testing.T) {
	a, _ := alert.New(validConfig())
	got := a.Evaluate("api", 10)
	if got == nil {
		t.Fatal("expected alert, got nil")
	}
	if got.Level != alert.LevelError {
		t.Fatalf("expected error level, got %s", got.Level)
	}
}

func TestEvaluate_MessageNotEmpty(t *testing.T) {
	a, _ := alert.New(validConfig())
	got := a.Evaluate("worker", 7)
	if got == nil {
		t.Fatal("expected alert")
	}
	if got.Message == "" {
		t.Fatal("expected non-empty message")
	}
}
