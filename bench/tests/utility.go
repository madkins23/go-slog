package tests

import (
	"runtime/debug"
	"testing"
)

// -----------------------------------------------------------------------------
// Utility methods.

func recoverAndFailOnPanic(b *testing.B) {
	r := recover()
	failOnPanic(b, r)
}

func failOnPanic(b *testing.B, r interface{}) {
	if r != nil {
		b.Errorf("test panicked: %v\n%s", r, debug.Stack())
		b.FailNow()
	}
}

const (
	message = "This is a message"
)
