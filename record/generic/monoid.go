package generic

import (
	M "github.com/ibm/fp-go/monoid"
	S "github.com/ibm/fp-go/semigroup"
)

func UnionMonoid[N ~map[K]V, K comparable, V any](s S.Semigroup[V]) M.Monoid[N] {
	return M.MakeMonoid(
		UnionSemigroup[N](s).Concat,
		Empty[N](),
	)
}
