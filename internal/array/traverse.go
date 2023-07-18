package array

import (
	F "github.com/IBM/fp-go/function"
)

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<GB>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func MonadTraverse[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,
	f func(A) HKTB) HKTRB {
	return MonadTraverseReduce(fof, fmap, fap, ta, f, Append[GB, B], Empty[GB]())
}

func Traverse[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(A) HKTB) func(GA) HKTRB {

	return func(ma GA) HKTRB {
		return MonadTraverse(fof, fmap, fap, ma, f)
	}
}

func MonadTraverseReduce[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := fmap(F.Curry2(reduce))

	return Reduce(ta, func(r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap,
			fap(transform(a)),
		)
	}, fof(initial))
}

func TraverseReduce[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) func(GA) HKTRB {
	return func(ta GA) HKTRB {
		return MonadTraverseReduce(fof, fmap, fap, ta, transform, reduce, initial)
	}
}
