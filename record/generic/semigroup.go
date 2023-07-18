package generic

import (
	S "github.com/IBM/fp-go/semigroup"
)

func UnionSemigroup[N ~map[K]V, K comparable, V any](s S.Semigroup[V]) S.Semigroup[N] {
	union := Union[N, K, V](s)
	return S.MakeSemigroup(func(first N, second N) N {
		return union(second)(first)
	})
}
