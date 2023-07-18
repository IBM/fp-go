package option

import (
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

func ApplySemigroup[A any](s S.Semigroup[A]) S.Semigroup[Option[A]] {
	return S.ApplySemigroup(MonadMap[A, func(A) A], MonadAp[A, A], s)
}

func ApplicativeMonoid[A any](m M.Monoid[A]) M.Monoid[Option[A]] {
	return M.ApplicativeMonoid(Of[A], MonadMap[A, func(A) A], MonadAp[A, A], m)
}
