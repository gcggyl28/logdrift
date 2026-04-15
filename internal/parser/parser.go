// Package parser provides structured log line parsing for common log formats.
package parser

import (
	"fmt"
	"regexp"
	"time"
)

// Format identifies a log line format.
type Format string

const (
	FormatAuto    Format = "auto"
	FormatJSON    Format = "json"
	FormatLogfmt  Format = "logfmt"
	FormatCommon  Format = "common" // Apache/nginx common log
)

// Entry holds the parsed fields extracted from a log line.
type Entry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]string
	Raw       string
}

// Parser parses raw log lines into structured Entry values.
type Parser struct {
	format  Format
	parseFn func(string) (Entry, error)
}

var commonLogRe = regexp.MustCompile(
	`^(\S+) \S+ \S+ \[([^\]]+)\] "([^"]+)" (\d+) (\d+|-)`)

// New constructs a Parser for the given format.
func New(format Format) (*Parser, error) {
	p := &Parser{format: format}
	switch format {
	case FormatJSON:
		p.parseFn = parseJSON
	case FormatLogfmt:
		p.parseFn = parseLogfmt
	case FormatCommon:
		p.parseFn = parseCommon
	case FormatAuto:
		p.parseFn = parseAuto
	default:
		return nil, fmt.Errorf("parser: unknown format %q", format)
	}
	return p, nil
}

// Parse parses a single raw log line.
func (p *Parser) Parse(line string) (Entry, error) {
	return p.parseFn(line)
}
