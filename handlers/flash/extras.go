package flash

import (
	"log/slog"
	"time"
)

const (
	DefaultTimeFormat = time.RFC3339Nano
)

// Extras defines extra options specific to a flash.Handler.
//
// Using these options it is possible to override some of the log/slog "standard" behavior.
// This supports testing of slog.HandlerOptions.ReplaceAttr functions and may also
// be used to replicate non-standard behavior in other handlers.
type Extras struct {
	// TimeFormat holds the format of the basic time field.
	// If not set defaults to the value of flash.DefaultTimeFormat (= time.RFC3339Nano).
	TimeFormat string

	// LevelNames holds a map from slog.Level to string.
	// If these fields are configured they replace the usual level names.
	// It is possible to configure only some of the level names.
	// Any level name that is not configured will be set to the appropriate
	// slog global constant (e.g. slog.LevelInfo.String()).
	LevelNames map[slog.Level]string

	// LevelKey specifies the JSON field name for the slog.Level for the log records.
	// If this field is not configured the value of slog.LevelKey is used.
	LevelKey string

	// MessageKey specifies the JSON field name for the log message.
	// If this field is not configured the value of slog.MessageKey is used.
	MessageKey string

	// SourceKey specifies the JSON field name for source data.
	// If this field is not configured the value of slog.SourceKey is used.
	SourceKey string

	// TimeKey specifies the JSON field name for the time the record was logged.
	// If this field is not configured the value of slog.TimeKey is used.
	TimeKey string
}

// fixExtras makes certain that an Extras object has been properly created and
// configured with default values.
//
// Use fixExtras(nil) to generate a new, otherwise blank Extras object (mostly useful in testing).
// Using &Extras{} for this purpose will likely end in tears and a pointer exception
// when some piece of code goes looking for one of the field values.
func fixExtras(extras *Extras) *Extras {
	if extras == nil {
		extras = &Extras{}
	}
	if extras.TimeFormat == "" {
		extras.TimeFormat = DefaultTimeFormat
	}
	if extras.LevelNames == nil {
		extras.LevelNames = make(map[slog.Level]string, 4)
	}
	if extras.LevelNames[slog.LevelDebug] == "" {
		extras.LevelNames[slog.LevelDebug] = slog.LevelDebug.String()
	}
	if extras.LevelNames[slog.LevelInfo] == "" {
		extras.LevelNames[slog.LevelInfo] = slog.LevelInfo.String()
	}
	if extras.LevelNames[slog.LevelWarn] == "" {
		extras.LevelNames[slog.LevelWarn] = slog.LevelWarn.String()
	}
	if extras.LevelNames[slog.LevelError] == "" {
		extras.LevelNames[slog.LevelError] = slog.LevelError.String()
	}
	if extras.LevelKey == "" {
		extras.LevelKey = slog.LevelKey
	}
	if extras.MessageKey == "" {
		extras.MessageKey = slog.MessageKey
	}
	if extras.SourceKey == "" {
		extras.SourceKey = slog.SourceKey
	}
	if extras.TimeKey == "" {
		extras.TimeKey = slog.TimeKey
	}
	return extras
}
