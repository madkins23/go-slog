package flash

import "sync"

// -----------------------------------------------------------------------------

type genPool[T any] struct {
	pool sync.Pool
}

func newGenPool[T any]() genPool[T] {
	return genPool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return new(T)
			},
		},
	}
}

// borrow returns an object like get as well as
// a function that will return that object via put.
// The latter is convenient for defer statements.
func (p *genPool[T]) borrow() (*T, func()) {
	x := p.pool.Get().(*T)
	return x, func() {
		p.put(x)
	}
}

// get is a generic wrapper around sync.Pool's Get method.
func (p *genPool[T]) get() *T {
	return p.pool.Get().(*T)
}

// put is a generic wrapper around sync.Pool's Put method.
func (p *genPool[T]) put(x *T) {
	p.pool.Put(x)
}

// -----------------------------------------------------------------------------

// arrayPool is a generic wrapper around a sync.Pool.
//
// Inspired by https://github.com/mkmik/syncpool
type arrayPool[T any] struct {
	pool sync.Pool
}

// New creates a new arrayPool with the provided new Function.
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

// borrow returns an object like get as well as
// a function that will return that object via put.
// The latter is convenient for defer statements.
func (p *arrayPool[T]) borrow() ([]T, func()) {
	x := p.pool.Get().([]T)
	return x, func() {
		p.put(x)
	}
}

// get is a generic wrapper around sync.Pool's Get method.
func (p *arrayPool[T]) get() []T {
	return p.pool.Get().([]T)
}

// put is a generic wrapper around sync.Pool's Put method.
func (p *arrayPool[T]) put(x []T) {
	// The x[:0] is supposed to reset len(x) to zero but leave cap(x) and
	// the underlying array space intact for reuse.
	p.pool.Put(x[:0])
}
