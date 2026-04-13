package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/export"
)

func sampleReport() export.DriftReport {
	return export.DriftReport{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		ServiceA:  "api",
		ServiceB:  "worker",
		Delta:     "-line1\n+line2",
		Severity:  "warn",
	}
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []export.Format{export.FormatJSON, export.FormatText} {
		_, err := export.New(&bytes.Buffer{}, f)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := export.New(&bytes.Buffer{}, export.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := export.New(&buf, export.FormatJSON)

	r := sampleReport()
	if err := ex.Write(r); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	var got export.DriftReport
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got.ServiceA != r.ServiceA {
		t.Errorf("service_a: want %q got %q", r.ServiceA, got.ServiceA)
	}
	if got.Delta != r.Delta {
		t.Errorf("delta: want %q got %q", r.Delta, got.Delta)
	}
}

func TestWrite_Text(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := export.New(&buf, export.FormatText)

	r := sampleReport()
	if err := ex.Write(r); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "api vs worker") {
		t.Errorf("expected service names in output, got: %q", out)
	}
	if !strings.Contains(out, r.Delta) {
		t.Errorf("expected delta in output, got: %q", out)
	}
}
