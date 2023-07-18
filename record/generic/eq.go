package generic

import (
	E "github.com/ibm/fp-go/eq"
)

func equals[M ~map[K]V, K comparable, V any](left, right M, eq func(V, V) bool) bool {
	if len(left) != len(right) {
		return false
	}
	for k, v1 := range left {
		if v2, ok := right[k]; !ok || !eq(v1, v2) {
			return false
		}
	}
	return true
}

func Eq[M ~map[K]V, K comparable, V any](e E.Eq[V]) E.Eq[M] {
	eq := e.Equals
	return E.FromEquals(func(left, right M) bool {
		return equals(left, right, eq)
	})
}
