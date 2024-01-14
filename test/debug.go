package test

import (
	"flag"
	"fmt"
)

// Set -debug flag to show extra print statements.
// Command line setting:
//
//	go test ./... -args -debug=<int>
//
// where the <int> is the debug level.
// If the level specified in the Debugf call is less than or equal the flag value
// the debug statement will be printed.
var debug = flag.Uint("debug", 0, "Show debug statements")

// Debugf will only print the specified data if the -debug command flag is set.
// The level field determines whether the statement will be printed.
// The -debug flag must be greater than or equal to the specified level for printing.
func Debugf(level uint, format string, args ...interface{}) {
	if *debug >= level {
		fmt.Printf(format, args...)
	}
}

// DebugLevel returns the level set by the -debug flag or 0 for default.
func DebugLevel() uint {
	return *debug
}
