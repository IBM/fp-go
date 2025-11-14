package stateless

import (
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
)

// ToSeq converts the stateless [Iterator] to an idiomatic go iterator
func ToSeq[T any](it Iterator[T]) Seq[T] {
	current := Current[T]
	return func(yield Predicate[T]) {
		next, ok := O.Unwrap(it())
		for ok && yield(current(next)) {
			next, ok = O.Unwrap(Next(next)())
		}
	}
}

// ToSeq2 converts the stateless [Iterator] to an idiomatic go iterator
func ToSeq2[K, V any](it Iterator[Pair[K, V]]) Seq2[K, V] {
	current := Current[Pair[K, V]]
	return func(yield func(K, V) bool) {
		yp := P.Paired(yield)
		next, ok := O.Unwrap(it())
		for ok && yp(current(next)) {
			next, ok = O.Unwrap(Next(next)())
		}
	}
}
