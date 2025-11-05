package util

import (
	"container/ring"
)

// Ring is a simple typed ring
type Ring[T any] struct {
	current T
	ring    *ring.Ring
}

// NewRing creates a new ring with the given values
func NewRing[T any](values ...T) *Ring[T] {
	r := ring.New(len(values))
	for _, value := range values {
		r.Value = value
		r = r.Next()
	}
	return &Ring[T]{
		current: r.Value.(T),
		ring:    r,
	}
}

// Value returns the value of the ring
func (r Ring[T]) Value() T {
	return r.current
}

// Forward moves the ring forward and returns the value
func (r *Ring[T]) Forward() T {
	r.ring = r.ring.Next()
	r.current = r.ring.Value.(T)
	return r.current
}

// Backward returns the previous ring
func (r *Ring[T]) Backward() T {
	r.ring = r.ring.Prev()
	r.current = r.ring.Value.(T)
	return r.current
}

// SetValue sets the value of the ring
func (r *Ring[T]) SetValue(value T) {
	r.ring.Value = value
}

func (r *Ring[T]) Move(n int) T {
	r.ring = r.ring.Move(n)
	r.current = r.ring.Value.(T)
	return r.current
}

func (r Ring[T]) Len() int {
	return r.ring.Len()
}

func (r *Ring[T]) Do(f func(T)) {
	fun := func(v any) {
		f(v.(T))
	}
	r.ring.Do(fun)
}
