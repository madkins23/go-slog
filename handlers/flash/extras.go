package flash

import "time"

const (
	DefaultTimeFormat = time.RFC3339Nano
)

// Extras defines extra options specific to a flash.Handler.
type Extras struct {
	TimeFormat string
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
