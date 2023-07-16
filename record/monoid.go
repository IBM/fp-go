package record

import (
	M "github.com/ibm/fp-go/monoid"
	G "github.com/ibm/fp-go/record/generic"
	S "github.com/ibm/fp-go/semigroup"
)

func UnionMonoid[K comparable, V any](s S.Semigroup[V]) M.Monoid[map[K]V] {
	return G.UnionMonoid[map[K]V](s)
}
