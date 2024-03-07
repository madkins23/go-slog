package flash

import "sync"

// -----------------------------------------------------------------------------

type bufferPool sync.Pool

func newBufferPool(size uint) *bufferPool {
	return &bufferPool{
		New: func() any {
			return make([]byte, 0, size)
		},
	}
}

// GetBuffer acquires a new or recycled buffer via sync.Pool.
func (p *bufferPool) GetBuffer() []byte {
	return (*sync.Pool)(p).Get().([]byte)
}

// PutBuffer returns a used buffer for reuse via sync.Pool.
func (p *bufferPool) PutBuffer(b []byte) {
	// Using b[:0] should set len=0 but leave capacity and underlying array intact.
	// This should be what make([], 0, size) would generate, but reuses the buffer.
	// The result should be a "clean up" and refresh of the buffer for reuse.
	(*sync.Pool)(p).Put(b[:0])
}

// -----------------------------------------------------------------------------
