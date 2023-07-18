package option

import (
	F "github.com/IBM/fp-go/function"
)

// HKTA = HKT<A>
// HKTOA = HKT<Option<A>>
//
// Sequence converts an option of some higher kinded type into the higher kinded type of an option
func Sequence[A, HKTA, HKTOA any](
	_of func(Option[A]) HKTOA,
	_map func(HKTA, func(A) Option[A]) HKTOA,
) func(Option[HKTA]) HKTOA {
	return Fold(F.Nullary2(None[A], _of), F.Bind2nd(_map, Some[A]))
}
