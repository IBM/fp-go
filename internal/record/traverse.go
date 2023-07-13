package record

import (
	F "github.com/ibm/fp-go/function"
)

// createEmpty creates a new empty, read-write map
// this is different to Empty which creates a new read-only empty map
func createEmpty[N ~map[K]A, K comparable, A any]() N {
	return make(N)
}

// inserts the key/value pair into a read-write map for performance
// order of parameters is adjusted to be curryable
func addKey[N ~map[K]A, K comparable, A any](key K, m N, value A) N {
	m[key] = value
	return m
}

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<MB>
HKTA = HKT<A>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func traverseWithIndex[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	_of func(MB) HKTRB,
	_map func(HKTRB, func(MB) func(B) MB) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	ta MA, f func(K, A) HKTB) HKTRB {
	// this function inserts a value into a map with a given key
	cb := F.Curry3(addKey[MB, K, B])

	return ReduceWithIndex(ta, func(k K, r HKTRB, a A) HKTRB {
		return _ap(
			_map(r, cb(k)),
			f(k, a),
		)
	}, _of(createEmpty[MB]()))
}

func MonadTraverse[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	_of func(MB) HKTRB,
	_map func(HKTRB, func(MB) func(B) MB) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	r MA, f func(A) HKTB) HKTRB {
	return traverseWithIndex(_of, _map, _ap, r, F.Ignore1of2[K](f))
}

func TraverseWithIndex[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	_of func(MB) HKTRB,
	_map func(HKTRB, func(MB) func(B) MB) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	f func(K, A) HKTB) func(MA) HKTRB {

	return func(ma MA) HKTRB {
		return traverseWithIndex(_of, _map, _ap, ma, f)
	}
}

// HKTA = HKT<A>
// HKTB = HKT<B>
// HKTAB = HKT<func(A)B>
// HKTRB = HKT<MB>
func Traverse[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	_of func(MB) HKTRB,
	_map func(HKTRB, func(MB) func(B) MB) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	f func(A) HKTB) func(MA) HKTRB {

	return func(ma MA) HKTRB {
		return MonadTraverse(_of, _map, _ap, ma, f)
	}
}

// HKTA = HKT[A]
// HKTAA = HKT[func(A)MA]
// HKTRA = HKT[MA]
func Sequence[MA ~map[K]A, MKTA ~map[K]HKTA, K comparable, A, HKTA, HKTAA, HKTRA any](
	_of func(MA) HKTRA,
	_map func(HKTRA, func(MA) func(A) MA) HKTAA,
	_ap func(HKTAA, HKTA) HKTRA,

	ma MKTA) HKTRA {
	return MonadTraverse(_of, _map, _ap, ma, F.Identity[HKTA])
}
