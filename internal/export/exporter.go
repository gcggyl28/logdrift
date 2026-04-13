// Package export provides functionality for writing log drift reports
// to external sinks such as files or stdout in structured formats.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Format represents the output format for exported drift reports.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// DriftReport holds a single drift event for export.
type DriftReport struct {
	Timestamp time.Time `json:"timestamp"`
	ServiceA  string    `json:"service_a"`
	ServiceB  string    `json:"service_b"`
	Delta     string    `json:"delta"`
	Severity  string    `json:"severity"`
}

// Exporter writes DriftReports to an io.Writer in a configured format.
type Exporter struct {
	format Format
	w      io.Writer
}

// New creates a new Exporter. Returns an error if the format is unsupported.
func New(w io.Writer, format Format) (*Exporter, error) {
	switch format {
	case FormatJSON, FormatText:
		// valid
	default:
		return nil, fmt.Errorf("export: unsupported format %q", format)
	}
	return &Exporter{format: format, w: w}, nil
}

// Write serialises and writes a DriftReport to the underlying writer.
func (e *Exporter) Write(r DriftReport) error {
	switch e.format {
	case FormatJSON:
		return e.writeJSON(r)
	case FormatText:
		return e.writeText(r)
	}
	return nil
}

func (e *Exporter) writeJSON(r DriftReport) error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("export: json marshal: %w", err)
	}
	_, err = fmt.Fprintf(e.w, "%s\n", data)
	return err
}

func (e *Exporter) writeText(r DriftReport) error {
	_, err := fmt.Fprintf(e.w, "[%s] %s vs %s (%s)\n%s\n",
		r.Timestamp.Format(time.RFC3339),
		r.ServiceA, r.ServiceB,
		r.Severity,
		r.Delta,
	)
	return err
}
