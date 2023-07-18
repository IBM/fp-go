package record

import (
	M "github.com/IBM/fp-go/monoid"
	G "github.com/IBM/fp-go/record/generic"
	S "github.com/IBM/fp-go/semigroup"
)

func UnionMonoid[K comparable, V any](s S.Semigroup[V]) M.Monoid[map[K]V] {
	return G.UnionMonoid[map[K]V](s)
}
