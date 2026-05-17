// Package alerter provides notification hooks for drift events.
package alerter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourorg/driftwatch/internal/manifest"
)

// Level represents the severity of a drift alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Alert holds the data for a single drift notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Report    manifest.DriftReport
}

// Notifier is the interface implemented by alert backends.
type Notifier interface {
	Notify(a Alert) error
}

// Alerter dispatches drift reports to one or more Notifiers.
type Alerter struct {
	notifiers []Notifier
}

// New creates an Alerter with the provided Notifiers.
func New(notifiers ...Notifier) *Alerter {
	return &Alerter{notifiers: notifiers}
}

// Send evaluates the report and dispatches an alert if drift is detected.
func (a *Alerter) Send(report manifest.DriftReport) error {
	if !report.HasDrift() {
		return nil
	}
	alert := Alert{
		Timestamp: time.Now().UTC(),
		Level:     LevelWarn,
		Report:    report,
	}
	var errs []error
	for _, n := range a.notifiers {
		if err := n.Notify(alert); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("alerter: %d notifier(s) failed: %v", len(errs), errs[0])
	}
	return nil
}

// LogNotifier writes alerts to an io.Writer (defaults to os.Stderr).
type LogNotifier struct {
	out io.Writer
}

// NewLogNotifier creates a LogNotifier writing to w; if w is nil os.Stderr is used.
func NewLogNotifier(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stderr
	}
	return &LogNotifier{out: w}
}

// Notify writes a human-readable drift alert line.
func (l *LogNotifier) Notify(a Alert) error {
	_, err := fmt.Fprintf(l.out, "[%s] %s drift detected in %d container(s)\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		len(a.Report.Containers),
	)
	return err
}
