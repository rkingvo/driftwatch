// Package reporter formats and outputs DriftReports to various sinks.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/driftwatch/internal/manifest"
)

// Format controls how the report is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes a DriftReport to an io.Writer.
type Reporter struct {
	format Format
	out    io.Writer
}

// New creates a Reporter with the given format and output writer.
func New(format Format, out io.Writer) *Reporter {
	return &Reporter{format: format, out: out}
}

// Write renders the report according to the configured format.
func (r *Reporter) Write(report manifest.DriftReport) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(report)
	default:
		return r.writeText(report)
	}
}

func (r *Reporter) writeJSON(report manifest.DriftReport) error {
	enc := json.NewEncoder(r.out)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func (r *Reporter) writeText(report manifest.DriftReport) error {
	if !report.HasDrift {
		fmt.Fprintln(r.out, "✔  No drift detected.")
		return nil
	}

	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CONTAINER\tFIELD\tEXPECTED\tACTUAL")
	for _, d := range report.Drifts {
		for _, f := range d.Fields {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.Container, f.Field, f.Expected, f.Actual)
		}
	}
	return w.Flush()
}
