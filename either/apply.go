package either

import (
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

func ApplySemigroup[E, A any](s S.Semigroup[A]) S.Semigroup[Either[E, A]] {
	return S.ApplySemigroup(MonadMap[E, A, func(A) A], MonadAp[A, E, A], s)
}

func ApplicativeMonoid[E, A any](m M.Monoid[A]) M.Monoid[Either[E, A]] {
	return M.ApplicativeMonoid(Of[E, A], MonadMap[E, A, func(A) A], MonadAp[A, E, A], m)
}
