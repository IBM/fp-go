package effect

import (
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/monoid"
)

func ApplicativeMonoid[C, A any](m monoid.Monoid[A]) Monoid[Effect[C, A]] {
	return readerreaderioresult.ApplicativeMonoid[C](m)
}

func AlternativeMonoid[C, A any](m monoid.Monoid[A]) Monoid[Effect[C, A]] {
	return readerreaderioresult.AlternativeMonoid[C](m)
}
