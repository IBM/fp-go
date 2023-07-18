package record

import (
	E "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/record/generic"
)

func Eq[K comparable, V any](e E.Eq[V]) E.Eq[map[K]V] {
	return G.Eq[map[K]V, K, V](e)
}
