package eq

import (
	F "github.com/IBM/fp-go/function"
)

type Eq[T any] interface {
	Equals(x, y T) bool
}

type eq[T any] struct {
	c func(x, y T) bool
}

func (self eq[T]) Equals(x, y T) bool {
	return self.c(x, y)
}

func strictEq[A comparable](a, b A) bool {
	return a == b
}

// FromStrictEquals constructs an `Eq` from the canonical comparison function
func FromStrictEquals[T comparable]() Eq[T] {
	return FromEquals(strictEq[T])
}

// FromEquals constructs an `Eq` from the comparison function
func FromEquals[T any](c func(x, y T) bool) Eq[T] {
	return eq[T]{c: c}
}

// Empty returns the equals predicate that is always true
func Empty[T any]() Eq[T] {
	return FromEquals(F.Constant2[T, T](true))
}

// Equals returns a predicate to test if one value equals the other under an equals predicate
func Equals[T any](eq Eq[T]) func(T) func(T) bool {
	return func(other T) func(T) bool {
		return F.Bind2nd(eq.Equals, other)
	}
}
