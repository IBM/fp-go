package array

import (
	F "github.com/ibm/fp-go/function"
)

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<GB>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func MonadTraverse[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	_of func(GB) HKTRB,
	_map func(HKTRB, func(GB) func(B) GB) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	ta GA,
	f func(A) HKTB) HKTRB {
	return MonadTraverseReduce(_of, _map, _ap, ta, f, Append[GB, B], Empty[GB]())
}

func Traverse[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	_of func(GB) HKTRB,
	_map func(HKTRB, func(GB) func(B) GB) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	f func(A) HKTB) func(GA) HKTRB {

	return func(ma GA) HKTRB {
		return MonadTraverse(_of, _map, _ap, ma, f)
	}
}

func MonadTraverseReduce[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	_of func(GB) HKTRB,
	_map func(HKTRB, func(GB) func(B) GB) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	ta GA,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := F.Bind2nd(_map, F.Curry2(reduce))

	return Reduce(ta, func(r HKTRB, a A) HKTRB {
		return _ap(
			mmap(r),
			transform(a),
		)
	}, _of(initial))
}

func TraverseReduce[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	_of func(GB) HKTRB,
	_map func(HKTRB, func(GB) func(B) GB) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) func(GA) HKTRB {
	return func(ta GA) HKTRB {
		return MonadTraverseReduce(_of, _map, _ap, ta, transform, reduce, initial)
	}
}
