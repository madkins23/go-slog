package test

import (
	"io"
	"sync/atomic"
)

var _ io.Writer = &CountWriter{}

type CountWriter struct {
	bytes atomic.Uint64
	count atomic.Uint64
}

func (cw *CountWriter) Write(p []byte) (n int, err error) {
	cw.bytes.Add(uint64(len(p)))
	cw.count.Add(1)
	return len(p), nil
}

func (cw *CountWriter) Bytes() uint64 {
	return cw.bytes.Load()
}

func (cw *CountWriter) Written() uint64 {
	return cw.count.Load()
}
