package record

import (
	G "github.com/ibm/fp-go/record/generic"
	S "github.com/ibm/fp-go/semigroup"
)

func UnionSemigroup[K comparable, V any](s S.Semigroup[V]) S.Semigroup[map[K]V] {
	return G.UnionSemigroup[map[K]V](s)
}
