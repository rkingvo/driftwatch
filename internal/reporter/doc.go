// Package reporter provides formatted output for DriftReports.
//
// Two formats are supported:
//
//	"text" — human-readable tabular output printed to stdout.
//	"json" — machine-readable JSON, suitable for piping to other tools.
//
// Usage:
//
//	r := reporter.New(reporter.FormatText, os.Stdout)
//	if err := r.Write(report); err != nil {
//	    log.Fatal(err)
//	}
package reporter
