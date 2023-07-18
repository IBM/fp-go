package array

import (
	E "github.com/IBM/fp-go/eq"
)

func equals[T any](left []T, right []T, eq func(T, T) bool) bool {
	if len(left) != len(right) {
		return false
	}
	for i, v1 := range left {
		v2 := right[i]
		if !eq(v1, v2) {
			return false
		}
	}
	return true
}

func Eq[T any](e E.Eq[T]) E.Eq[[]T] {
	eq := e.Equals
	return E.FromEquals(func(left, right []T) bool {
		return equals(left, right, eq)
	})
}
