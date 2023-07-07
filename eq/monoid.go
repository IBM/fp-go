package eq

import (
	M "github.com/ibm/fp-go/monoid"
	S "github.com/ibm/fp-go/semigroup"
)

func Semigroup[A any]() S.Semigroup[Eq[A]] {
	return S.MakeSemigroup(func(x, y Eq[A]) Eq[A] {
		return FromEquals(func(a, b A) bool {
			return x.Equals(a, b) && y.Equals(a, b)
		})
	})
}

func Monoid[A any]() M.Monoid[Eq[A]] {
	return M.MakeMonoid(Semigroup[A]().Concat, Empty[A]())
}
