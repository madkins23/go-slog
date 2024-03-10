package flash

import "time"

const (
	DefaultTimeFormat = time.RFC3339Nano
)

// Extras defines extra options specific to a flash.Handler.
type Extras struct {
	// Holds the format of the basic time field.
	// If not set defaults to the value of flash.DefaultTimeFormat (= time.RFC3339Nano).
	TimeFormat string

	// Causes embedded UTF8 to be escaped in all quoted strings.
	EscapeUTF8 bool
}

func fixExtras(extras *Extras) *Extras {
	if extras == nil {
		extras = &Extras{}
	}
	if extras.TimeFormat == "" {
		extras.TimeFormat = DefaultTimeFormat
	}
	return extras
}
