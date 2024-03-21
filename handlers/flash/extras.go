package flash

import (
	"log/slog"
	"time"
)

const (
	DefaultTimeFormat = time.RFC3339Nano
)

// Extras defines extra options specific to a flash.Handler.
type Extras struct {
	// Holds the format of the basic time field.
	// If not set defaults to the value of flash.DefaultTimeFormat (= time.RFC3339Nano).
	TimeFormat string

	LevelNames map[slog.Level]string
	LevelKey   string
	MessageKey string
	SourceKey  string
	TimeKey    string
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
