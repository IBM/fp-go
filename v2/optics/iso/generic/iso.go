package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
	I "github.com/IBM/fp-go/v2/optics/iso"
)

// AsTraversal converts a iso to a traversal
func AsTraversal[R ~func(func(A) HKTA) func(S) HKTS, S, A, HKTS, HKTA any](
	fmap functor.MapType[A, S, HKTA, HKTS],
) func(I.Iso[S, A]) R {
	return func(sa I.Iso[S, A]) R {
		saSet := fmap(sa.ReverseGet)
		return func(f func(A) HKTA) func(S) HKTS {
			return F.Flow3(
				sa.Get,
				f,
				saSet,
			)
		}
	}
}
