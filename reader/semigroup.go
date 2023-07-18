package reader

import (
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

func ApplySemigroup[R, A any](
	_map func(func(R) A, func(A) func(A) A) func(R, func(A) A),
	_ap func(func(R, func(A) A), func(R) A) func(R) A,

	s S.Semigroup[A],
) S.Semigroup[func(R) A] {
	return S.ApplySemigroup(_map, _ap, s)
}

func ApplicativeMonoid[R, A any](
	_of func(A) func(R) A,
	_map func(func(R) A, func(A) func(A) A) func(R, func(A) A),
	_ap func(func(R, func(A) A), func(R) A) func(R) A,

	m M.Monoid[A],
) M.Monoid[func(R) A] {
	return M.ApplicativeMonoid(_of, _map, _ap, m)
}
