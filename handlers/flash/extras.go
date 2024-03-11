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
}

// fixExtras makes certain that an Extras object has been properly created and
// configured with default values.
//
// Use fixExtras(nil) to generate a new, otherwise blank Extras object (mostly useful in testing).
// Using &Extras{} for this purpose will likely end in tears and a pointer exception
// when some piece of code goes looking for extras.TimeFormat.
func fixExtras(extras *Extras) *Extras {
	if extras == nil {
		extras = &Extras{}
	}
	if extras.TimeFormat == "" {
		extras.TimeFormat = DefaultTimeFormat
	}
	return extras
}
