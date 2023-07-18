package record

import (
	G "github.com/IBM/fp-go/internal/record"
)

func TraverseWithIndex[K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(map[K]B) HKTRB,
	fmap func(func(map[K]B) func(B) map[K]B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(K, A) HKTB) func(map[K]A) HKTRB {
	return G.TraverseWithIndex[map[K]A](fof, fmap, fap, f)
}

// HKTA = HKT<A>
// HKTB = HKT<B>
// HKTAB = HKT<func(A)B>
// HKTRB = HKT<map[K]B>
func Traverse[K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(map[K]B) HKTRB,
	fmap func(func(map[K]B) func(B) map[K]B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,
	f func(A) HKTB) func(map[K]A) HKTRB {
	return G.Traverse[map[K]A](fof, fmap, fap, f)
}

// HKTA = HKT[A]
// HKTAA = HKT[func(A)map[K]A]
// HKTRA = HKT[map[K]A]
func Sequence[K comparable, A, HKTA, HKTAA, HKTRA any](
	fof func(map[K]A) HKTRA,
	fmap func(func(map[K]A) func(A) map[K]A) func(HKTRA) HKTAA,
	fap func(HKTA) func(HKTAA) HKTRA,
	ma map[K]HKTA) HKTRA {
	return G.Sequence(fof, fmap, fap, ma)

}
