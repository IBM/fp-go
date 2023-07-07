package ord

import (
	O "github.com/ibm/fp-go/ord"
	S "github.com/ibm/fp-go/semigroup"
)

// Max gets a semigroup where `concat` will return the maximum, based on the provided order.
func Max[A any](o O.Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(O.Max(o))
}

// Min gets a semigroup where `concat` will return the minimum, based on the provided order.
func Min[A any](o O.Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(O.Min(o))
}
