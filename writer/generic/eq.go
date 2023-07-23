package generic

import (
	EQ "github.com/IBM/fp-go/eq"
	T "github.com/IBM/fp-go/tuple"
)

// Constructs an equal predicate for a [Writer]
func Eq[GA ~func() T.Tuple2[A, W], W, A any](w EQ.Eq[W], a EQ.Eq[A]) EQ.Eq[GA] {
	return EQ.FromEquals(func(l, r GA) bool {
		ll := l()
		rr := r()

		return a.Equals(ll.F1, rr.F1) && w.Equals(ll.F2, rr.F2)
	})
}

// FromStrictEquals constructs an `Eq` from the canonical comparison function
func FromStrictEquals[GA ~func() T.Tuple2[A, W], W, A comparable]() EQ.Eq[GA] {
	return Eq[GA](EQ.FromStrictEquals[W](), EQ.FromStrictEquals[A]())
}
