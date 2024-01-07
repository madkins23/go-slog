package test

import (
	"flag"
	"fmt"
)

// Set -debug flag to show extra print statements.
// Command line setting:
//
//	go test ./... -args -debug
var debug = flag.Uint("debug", 0, "Show debug statements")

// Debugf will only print the specified data if the -debug command flag is set.
// The format strings will be wrapped with '>>> ' before and '\n' after.
func Debugf(level uint, format string, args ...interface{}) {
	if *debug >= level {
		fmt.Printf(">>> "+format+"\n", args...)
	}
}
