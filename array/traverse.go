package array

import "github.com/IBM/fp-go/internal/array"

func Traverse[A, B, HKTB, HKTAB, HKTRB any](
	fof func([]B) HKTRB,
	fmap func(func([]B) func(B) []B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(A) HKTB) func([]A) HKTRB {
	return array.Traverse[[]A](fof, fmap, fap, f)
}

func MonadTraverse[A, B, HKTB, HKTAB, HKTRB any](
	fof func([]B) HKTRB,
	fmap func(func([]B) func(B) []B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta []A,
	f func(A) HKTB) HKTRB {

	return array.MonadTraverse(fof, fmap, fap, ta, f)
}
