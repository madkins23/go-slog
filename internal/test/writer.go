package test

import (
	"io"
	"sync/atomic"
)

var _ io.Writer = &CountWriter{}

// CountWriter is an io.Writer that throws away all input but
// counts `Write` calls and the number of bytes that would have been written.
// This is used during benchmarking.
type CountWriter struct {
	count atomic.Uint64
}

// Write supplies the required io.Writer interface method.
func (cw *CountWriter) Write(p []byte) (n int, err error) {
	if len(p) > 0 && p[len(p)-1] == '\n' {
		cw.count.Add(1)
	}
	return len(p), nil
}

// Written returns the number of `Write` calls that the `CountWriter` ignored.
func (cw *CountWriter) Written() uint64 {
	return cw.count.Load()
}
