package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
)

func FromLens[S, A, HKTES, HKTA any](
	fmap functor.MapType[A, Endomorphism[S], HKTA, HKTES],
) func(Lens[S, A]) Traversal[S, A, HKTES, HKTA] {
	return func(sa Lens[S, A]) Traversal[S, A, HKTES, HKTA] {
		saGet := sa.Get
		saSet := fmap(sa.Set)
		return func(f func(A) HKTA) func(S) HKTES {
			return F.Flow3(
				saGet,
				f,
				saSet,
			)
		}
	}
}
