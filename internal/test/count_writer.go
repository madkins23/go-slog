package test

import (
	"io"
	"sync/atomic"
)

var _ io.Writer = &CountWriter{}

type CountWriter struct {
	count atomic.Uint64
}

func (cw *CountWriter) Write(p []byte) (n int, err error) {
	cw.count.Add(1)
	return len(p), nil
}

func (cw *CountWriter) Written() uint64 {
	return cw.count.Load()
}
