// Package alert provides threshold-based alerting for log drift events.
// When the number of differing lines within a rolling window exceeds a
// configured threshold, an Alert is emitted on the output channel.
package alert

import (
	"fmt"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Alert carries information about a drift threshold breach.
type Alert struct {
	Service   string
	Level     Level
	DriftCount int
	Threshold  int
	At         time.Time
	Message    string
}

// Config holds the alerting configuration.
type Config struct {
	// WarnThreshold triggers a warn-level alert when drift count reaches this value.
	WarnThreshold int
	// ErrorThreshold triggers an error-level alert.
	ErrorThreshold int
	// Window is the rolling time window over which drifts are counted.
	Window time.Duration
}

// Alerter evaluates drift counts against configured thresholds.
type Alerter struct {
	cfg Config
}

// New returns an Alerter configured with cfg.
// Returns an error if thresholds are non-positive or ErrorThreshold < WarnThreshold.
func New(cfg Config) (*Alerter, error) {
	if cfg.WarnThreshold <= 0 {
		return nil, fmt.Errorf("alert: WarnThreshold must be > 0, got %d", cfg.WarnThreshold)
	}
	if cfg.ErrorThreshold <= 0 {
		return nil, fmt.Errorf("alert: ErrorThreshold must be > 0, got %d", cfg.ErrorThreshold)
	}
	if cfg.ErrorThreshold < cfg.WarnThreshold {
		return nil, fmt.Errorf("alert: ErrorThreshold (%d) must be >= WarnThreshold (%d)",
			cfg.ErrorThreshold, cfg.WarnThreshold)
	}
	if cfg.Window <= 0 {
		return nil, fmt.Errorf("alert: Window must be > 0")
	}
	return &Alerter{cfg: cfg}, nil
}

// Evaluate checks driftCount for service against configured thresholds.
// Returns a non-nil *Alert when a threshold is breached, otherwise nil.
func (a *Alerter) Evaluate(service string, driftCount int) *Alert {
	var lvl Level
	switch {
	case driftCount >= a.cfg.ErrorThreshold:
		lvl = LevelError
	case driftCount >= a.cfg.WarnThreshold:
		lvl = LevelWarn
	default:
		return nil
	}
	return &Alert{
		Service:    service,
		Level:      lvl,
		DriftCount: driftCount,
		Threshold:  a.cfg.WarnThreshold,
		At:         time.Now(),
		Message: fmt.Sprintf("[%s] service %q exceeded drift threshold: %d drifted lines in window %s",
			lvl, service, driftCount, a.cfg.Window),
	}
}
