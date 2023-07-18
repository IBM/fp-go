package generic

import (
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

func UnionMonoid[N ~map[K]V, K comparable, V any](s S.Semigroup[V]) M.Monoid[N] {
	return M.MakeMonoid(
		UnionSemigroup[N](s).Concat,
		Empty[N](),
	)
}
