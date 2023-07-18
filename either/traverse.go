package either

import (
	F "github.com/ibm/fp-go/function"
)

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<Either[B]>
HKTA = HKT<A>
HKTB = HKT<B>
*/
func traverse[E, A, B, HKTA, HKTB, HKTRB any](
	_of func(Either[E, B]) HKTRB,
	_map func(HKTB, func(B) Either[E, B]) HKTRB,
) func(Either[E, A], func(A) HKTB) HKTRB {

	left := F.Flow2(Left[B, E], _of)
	right := F.Bind2nd(_map, Right[E, B])

	return func(ta Either[E, A], f func(A) HKTB) HKTRB {
		return MonadFold(ta,
			left,
			F.Flow2(f, right),
		)
	}
}

func Traverse[E, A, B, HKTA, HKTB, HKTRB any](
	_of func(Either[E, B]) HKTRB,
	_map func(HKTB, func(B) Either[E, B]) HKTRB,
) func(func(A) HKTB) func(Either[E, A]) HKTRB {
	delegate := traverse[E, A, B, HKTA](_of, _map)
	return func(f func(A) HKTB) func(Either[E, A]) HKTRB {
		return F.Bind2nd(delegate, f)
	}
}

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRA = HKT<Either[A]>
HKTA = HKT<A>
HKTB = HKT<B>
*/
func Sequence[E, A, HKTA, HKTRA any](
	_of func(Either[E, A]) HKTRA,
	_map func(HKTA, func(A) Either[E, A]) HKTRA,
) func(Either[E, HKTA]) HKTRA {
	return Fold(F.Flow2(Left[A, E], _of), F.Bind2nd(_map, Right[E, A]))
}
