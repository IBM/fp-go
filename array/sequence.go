package array

import (
	F "github.com/ibm/fp-go/function"
	O "github.com/ibm/fp-go/option"
)

// We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

// HKTA = HKT<A>
// HKTRA = HKT<[]A>
// HKTFRA = HKT<func(A)[]A>

// Sequence takes an `Array` where elements are `HKT<A>` (higher kinded type) and,
// using an applicative of that `HKT`, returns an `HKT` of `[]A`.
// e.g. it can turn an `[]Either[error, string]` into an `Either[error, []string]`.
//
// Sequence requires an `Applicative` of the `HKT` you are targeting, e.g. to turn an
// `[]Either[E, A]` into an `Either[E, []A]`, it needs an
// Applicative` for `Either`, to to turn an `[]Option[A]` into an `Option[ []A]`,
// it needs an `Applicative` for `Option`.
func Sequence[A, HKTA, HKTRA, HKTFRA any](
	_of func([]A) HKTRA,
	_map func(HKTRA, func([]A) func(A) []A) HKTFRA,
	_ap func(HKTFRA, HKTA) HKTRA,
) func([]HKTA) HKTRA {
	ca := F.Curry2(Append[A])
	return Reduce(func(fas HKTRA, fa HKTA) HKTRA {
		return _ap(
			_map(fas, ca),
			fa,
		)
	}, _of(Empty[A]()))
}

// ArrayOption returns a function to convert sequence of options into an option of a sequence
func ArrayOption[A any]() func([]O.Option[A]) O.Option[[]A] {
	return Sequence(
		O.Of[[]A],
		O.MonadMap[[]A, func(A) []A],
		O.MonadAp[[]A, A],
	)
}
