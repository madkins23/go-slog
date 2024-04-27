package flash

import "sync"

// This file contains generic memory pool definitions wrapped around sync.Pool.

// TODO: Doesn't seem a need to clean these up and move them into go-slog/internal
//       until there is another use for them outside of handlers/flash.

// -----------------------------------------------------------------------------

// genPool is a generic wrapper around a sync.Pool for a pool of objects.
//
// Inspired by https://github.com/mkmik/syncpool
type genPool[T any] struct {
	pool sync.Pool
}

// New creates a new genPool with the provided New function which returns a new object.
//
// The equivalent sync.Pool construct is "sync.Pool{New: fn}"
func newGenPool[T any]() genPool[T] {
	return genPool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return new(T)
			},
		},
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

// arrayPool is a generic wrapper around a sync.Pool for a pool of arrays of objects.
//
// Inspired by https://github.com/mkmik/syncpool
type arrayPool[T any] struct {
	pool sync.Pool
}

// New creates a new arrayPool with the provided New function which returns an empty array of objects.
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

// get is a generic wrapper around sync.Pool's Get method.
func (p *arrayPool[T]) get() []T {
	return p.pool.Get().([]T)
}

// put is a generic wrapper around sync.Pool's Put method.
func (p *arrayPool[T]) put(x []T) {
	p.pool.Put(x)
}
