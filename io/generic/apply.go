package generic

import (
	M "github.com/ibm/fp-go/monoid"
	S "github.com/ibm/fp-go/semigroup"
)

func ApplySemigroup[GA ~func() A, A any](s S.Semigroup[A]) S.Semigroup[GA] {
	return S.ApplySemigroup(MonadMap[GA, func() func(A) A, A, func(A) A], MonadAp[GA, GA, func() func(A) A, A, A], s)
}

func ApplicativeMonoid[GA ~func() A, A any](m M.Monoid[A]) M.Monoid[GA] {
	return M.ApplicativeMonoid(Of[GA, A], MonadMap[GA, func() func(A) A, A, func(A) A], MonadAp[GA, GA, func() func(A) A, A, A], m)
}
