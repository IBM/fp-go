package generic

import (
	M "github.com/IBM/fp-go/v2/monoid"
)

//go:inline
func ApplicativeMonoid[GA ~func(R) A, R, A any](m M.Monoid[A]) M.Monoid[GA] {
	return M.ApplicativeMonoid(
		Of[GA, R, A],
		MonadMap[GA, func(R) func(A) A],
		MonadAp[GA, GA, func(R) func(A) A],
		m,
	)
}
