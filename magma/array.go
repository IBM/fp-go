package magma

import (
	F "github.com/IBM/fp-go/function"
	AR "github.com/IBM/fp-go/internal/array"
)

func GenericMonadConcatAll[GA ~[]A, A any](m Magma[A]) func(GA, A) A {
	return func(as GA, first A) A {
		return AR.Reduce(as, m.Concat, first)
	}
}

// GenericConcatAll concats all items using the semigroup and a starting value
func GenericConcatAll[GA ~[]A, A any](m Magma[A]) func(A) func(GA) A {
	ca := GenericMonadConcatAll[GA](m)
	return func(a A) func(GA) A {
		return F.Bind2nd(ca, a)
	}
}

func MonadConcatAll[A any](m Magma[A]) func([]A, A) A {
	return GenericMonadConcatAll[[]A](m)
}

// ConcatAll concats all items using the semigroup and a starting value
func ConcatAll[A any](m Magma[A]) func(A) func([]A) A {
	return GenericConcatAll[[]A](m)
}
