package array

import "github.com/ibm/fp-go/internal/array"

func Traverse[A, B, HKTB, HKTAB, HKTRB any](
	_of func([]B) HKTRB,
	_map func(HKTRB, func([]B) func(B) []B) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	f func(A) HKTB) func([]A) HKTRB {
	return array.Traverse[[]A](_of, _map, _ap, f)
}

func MonadTraverse[A, B, HKTB, HKTAB, HKTRB any](
	_of func([]B) HKTRB,
	_map func(HKTRB, func([]B) func(B) []B) HKTAB,
	_ap func(HKTAB, HKTB) HKTRB,

	ta []A,
	f func(A) HKTB) HKTRB {

	return array.MonadTraverse(_of, _map, _ap, ta, f)
}
