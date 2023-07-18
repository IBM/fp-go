package monoid

import (
	S "github.com/IBM/fp-go/semigroup"
)

type Monoid[A any] interface {
	S.Semigroup[A]
	Empty() A
}

type monoid[A any] struct {
	c func(A, A) A
	e A
}

func (self monoid[A]) Concat(x A, y A) A {
	return self.c(x, y)
}

func (self monoid[A]) Empty() A {
	return self.e
}

// MakeMonoid creates a monoid given a concat function and an empty element
func MakeMonoid[A any](c func(A, A) A, e A) Monoid[A] {
	return monoid[A]{c: c, e: e}
}

// Reverse returns the dual of a `Monoid`, obtained by swapping the arguments of `Concat`.
func Reverse[A any](m Monoid[A]) Monoid[A] {
	return MakeMonoid(S.Reverse[A](m).Concat, m.Empty())
}
