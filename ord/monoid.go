package ord

import (
	F "github.com/IBM/fp-go/function"
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

// Semigroup implements a two level ordering
func Semigroup[A any]() S.Semigroup[Ord[A]] {
	return S.MakeSemigroup(func(first, second Ord[A]) Ord[A] {
		return FromCompare(func(a, b A) int {
			ox := first.Compare(a, b)
			if ox != 0 {
				return ox
			}
			return second.Compare(a, b)
		})
	})
}

// Monoid implements a two level ordering such that
// - its `Concat(ord1, ord2)` operation will order first by `ord1`, and then by `ord2`
// - its `Empty` value is an `Ord` that always considers compared elements equal
func Monoid[A any]() M.Monoid[Ord[A]] {
	return M.MakeMonoid(Semigroup[A]().Concat, FromCompare(F.Constant2[A, A](0)))
}

// MaxSemigroup returns a semigroup where `concat` will return the maximum, based on the provided order.
func MaxSemigroup[A any](O Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(Max(O))
}

// MaxSemigroup returns a semigroup where `concat` will return the minimum, based on the provided order.
func MinSemigroup[A any](O Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(Min(O))
}
