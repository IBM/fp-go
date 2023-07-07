package semigroup

import (
	F "github.com/ibm/fp-go/function"
	M "github.com/ibm/fp-go/magma"
)

type Semigroup[A any] interface {
	M.Magma[A]
}

type semigroup[A any] struct {
	c func(A, A) A
}

func (self semigroup[A]) Concat(x A, y A) A {
	return self.c(x, y)
}

func MakeSemigroup[A any](c func(A, A) A) Semigroup[A] {
	return semigroup[A]{c: c}
}

// Reverse returns The dual of a `Semigroup`, obtained by swapping the arguments of `concat`.
func Reverse[A any](m Semigroup[A]) Semigroup[A] {
	return MakeSemigroup(M.Reverse[A](m).Concat)
}

// FunctionSemigroup forms a semigroup as long as you can provide a semigroup for the codomain.
func FunctionSemigroup[A, B any](S Semigroup[B]) Semigroup[func(A) B] {
	return MakeSemigroup(func(f func(A) B, g func(A) B) func(A) B {
		return func(a A) B {
			return S.Concat(f(a), g(a))
		}
	})
}

// First always returns the first argument.
func First[A any]() Semigroup[A] {
	return MakeSemigroup(F.First[A, A])
}

// Last always returns the last argument.
func Last[A any]() Semigroup[A] {
	return MakeSemigroup(F.Second[A, A])
}
