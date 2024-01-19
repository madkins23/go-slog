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
	bytes atomic.Uint64
	count atomic.Uint64
}

// Write supplies the required io.Writer interface method.
func (cw *CountWriter) Write(p []byte) (n int, err error) {
	cw.bytes.Add(uint64(len(p)))
	cw.count.Add(1)
	return len(p), nil
}

// Bytes returns the number of bytes that the `CountWriter` ignored.
func (cw *CountWriter) Bytes() uint64 {
	return cw.bytes.Load()
}

// Written returns the number of `Write` calls that the `CountWriter` ignored.
func (cw *CountWriter) Written() uint64 {
	return cw.count.Load()
}
