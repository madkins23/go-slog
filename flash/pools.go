package flash

import "sync"

// -----------------------------------------------------------------------------

// A arrayPool is a generic wrapper around a sync.Pool.
type arrayPool[T any] struct {
	pool sync.Pool
}

// New creates a new arrayPool with the provided new function.
//
// The equivalent sync.Pool construct is "sync.Pool{New: fn}"
func newArrayPool[T any](size uint) arrayPool[T] {
	return arrayPool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]T, 0, size)
			},
		},
	}
}

// Get is a generic wrapper around sync.Pool's Get method.
func (p *arrayPool[T]) get() []T {
	return p.pool.Get().([]T)
}

// Put is a generic wrapper around sync.Pool's Put method.
func (p *arrayPool[T]) put(x []T) {
	p.pool.Put(x[:0])
}
