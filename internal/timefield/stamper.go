package timefield

import (
	"fmt"
	"time"

	"github.com/angryboat/logdrift/internal/parser"
)

// Stamper enriches entries that are missing a timestamp field by writing the
// current wall-clock time into them.
type Stamper struct {
	extractor *Extractor
	field     string
	layout    string
}

// NewStamper creates a Stamper that uses field for storage and layout for
// formatting. If the entry already contains a parseable timestamp the entry is
// returned unchanged.
func NewStamper(field, layout string, opts ...Option) (*Stamper, error) {
	if field == "" {
		return nil, fmt.Errorf("timefield: stamper field must not be empty")
	}
	if layout == "" {
		return nil, fmt.Errorf("timefield: stamper layout must not be empty")
	}
	ext, err := New(field, opts...)
	if err != nil {
		return nil, err
	}
	return &Stamper{extractor: ext, field: field, layout: layout}, nil
}

// Apply stamps entry with the current time when the timestamp field is absent.
func (s *Stamper) Apply(entry parser.Entry) parser.Entry {
	if _, err := s.extractor.Extract(entry); err == nil {
		return entry
	}
	if entry.Fields == nil {
		entry.Fields = make(map[string]string)
	}
	entry.Fields[s.field] = time.Now().UTC().Format(s.layout)
	return entry
}
