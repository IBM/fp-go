package io

import (
	G "github.com/ibm/fp-go/io/generic"
	M "github.com/ibm/fp-go/monoid"
	S "github.com/ibm/fp-go/semigroup"
)

func ApplySemigroup[A any](s S.Semigroup[A]) S.Semigroup[IO[A]] {
	return G.ApplySemigroup[IO[A]](s)
}

func ApplicativeMonoid[A any](m M.Monoid[A]) M.Monoid[IO[A]] {
	return G.ApplicativeMonoid[IO[A]](m)
}
