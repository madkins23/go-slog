package infra

import (
	"runtime"
	"strings"
)

// CurrentFunctionName checks up the call stack for the name of the current test function.
// Only the last part of the function name (after the last period) is returned.
// The function name is found by checking for a prefix of "Test".
// If no test function is found "Unknown" is returned.
func CurrentFunctionName() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	more := true
	for more {
		var frame runtime.Frame
		frame, more = frames.Next()
		parts := strings.Split(frame.Function, ".")
		functionName := parts[len(parts)-1]
		if strings.HasPrefix(functionName, "Test") {
			return functionName
		}
	}
	return "Unknown"
}
