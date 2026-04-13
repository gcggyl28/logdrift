package render

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Format controls how drift output is rendered to the terminal.
type Format string

const (
	FormatPlain    Format = "plain"
	FormatColored  Format = "colored"
	FormatTimestamp Format = "timestamp"
)

// Renderer writes formatted diff output to a writer.
type Renderer struct {
	format Format
	out    io.Writer

	addColor *color.Color
	delColor *color.Color
	headColor *color.Color
}

// New returns a Renderer for the given format. Writes to os.Stdout by default.
func New(format Format) (*Renderer, error) {
	switch format {
	case FormatPlain, FormatColored, FormatTimestamp:
	default:
		return nil, fmt.Errorf("render: unknown format %q", format)
	}
	return &Renderer{
		format:    format,
		out:       os.Stdout,
		addColor:  color.New(color.FgGreen),
		delColor:  color.New(color.FgRed),
		headColor: color.New(color.FgCyan, color.Bold),
	}, nil
}

// SetOutput redirects render output (useful for testing).
func (r *Renderer) SetOutput(w io.Writer) {
	r.out = w
}

// WriteDrift formats and writes a drift block to the output.
func (r *Renderer) WriteDrift(serviceA, serviceB, delta string) {
	prefix := ""
	if r.format == FormatTimestamp {
		prefix = fmt.Sprintf("[%s] ", time.Now().UTC().Format(time.RFC3339))
	}

	header := fmt.Sprintf("%s--- %s  +++ %s", prefix, serviceA, serviceB)

	if r.format == FormatColored || r.format == FormatTimestamp {
		r.headColor.Fprintln(r.out, header)
		for _, line := range strings.Split(delta, "\n") {
			switch {
			case strings.HasPrefix(line, "+"):
				r.addColor.Fprintln(r.out, line)
			case strings.HasPrefix(line, "-"):
				r.delColor.Fprintln(r.out, line)
			default:
				fmt.Fprintln(r.out, line)
			}
		}
	} else {
		fmt.Fprintln(r.out, header)
		fmt.Fprintln(r.out, delta)
	}
}
