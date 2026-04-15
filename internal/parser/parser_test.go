package parser

import (
	"testing"
	"time"
)

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []Format{FormatJSON, FormatLogfmt, FormatCommon, FormatAuto} {
		_, err := New(f)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestParse_JSON(t *testing.T) {
	p, _ := New(FormatJSON)
	line := `{"level":"info","msg":"hello world","time":"2024-01-02T15:04:05Z","svc":"api"}`
	e, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Level != "INFO" {
		t.Errorf("level: got %q want INFO", e.Level)
	}
	if e.Message != "hello world" {
		t.Errorf("message: got %q", e.Message)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if e.Fields["svc"] != "api" {
		t.Errorf("fields[svc]: got %q", e.Fields["svc"])
	}
}

func TestParse_Logfmt(t *testing.T) {
	p, _ := New(FormatLogfmt)
	line := `level=warn msg="disk full" host=srv1 ts=2024-01-02T15:04:05Z`
	e, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Level != "WARN" {
		t.Errorf("level: got %q want WARN", e.Level)
	}
	if e.Message != "disk full" {
		t.Errorf("message: got %q", e.Message)
	}
	if e.Fields["host"] != "srv1" {
		t.Errorf("fields[host]: got %q", e.Fields["host"])
	}
}

func TestParse_Common(t *testing.T) {
	p, _ := New(FormatCommon)
	line := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	e, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Fields["status"] != "200" {
		t.Errorf("status: got %q", e.Fields["status"])
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestParse_Auto_SelectsJSON(t *testing.T) {
	p, _ := New(FormatAuto)
	e, err := p.Parse(`{"msg":"ok"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Message != "ok" {
		t.Errorf("message: got %q", e.Message)
	}
}

func TestParse_InvalidJSON_ReturnsError(t *testing.T) {
	p, _ := New(FormatJSON)
	_, err := p.Parse("{not valid json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestEntry_RawPreserved(t *testing.T) {
	p, _ := New(FormatLogfmt)
	raw := "level=debug msg=test"
	e, _ := p.Parse(raw)
	if e.Raw != raw {
		t.Errorf("raw: got %q want %q", e.Raw, raw)
	}
}

func TestEntry_ZeroTimestampWhenMissing(t *testing.T) {
	p, _ := New(FormatLogfmt)
	e, _ := p.Parse("level=info msg=hello")
	if !e.Timestamp.Equal(time.Time{}) {
		t.Errorf("expected zero timestamp, got %v", e.Timestamp)
	}
}
