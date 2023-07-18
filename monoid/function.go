package monoid

import (
	F "github.com/IBM/fp-go/function"
	S "github.com/IBM/fp-go/semigroup"
)

// FunctionMonoid forms a monoid as long as you can provide a monoid for the codomain.
func FunctionMonoid[A, B any](M Monoid[B]) Monoid[func(A) B] {
	return MakeMonoid(
		S.FunctionSemigroup[A, B](M).Concat,
		F.Constant1[A](M.Empty()),
	)
}
