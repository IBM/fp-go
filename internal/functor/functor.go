package functor

import (
	F "github.com/IBM/fp-go/function"
)

// HKTFGA = HKT[F, HKT[G, A]]
// HKTFGB = HKT[F, HKT[G, B]]
func MonadMap[A, B, HKTGA, HKTGB, HKTFGA, HKTFGB any](fmap func(HKTFGA, func(HKTGA) HKTGB) HKTFGB, gmap func(HKTGA, func(A) B) HKTGB, fa HKTFGA, f func(A) B) HKTFGB {
	return fmap(fa, F.Bind2nd(gmap, f))
}
