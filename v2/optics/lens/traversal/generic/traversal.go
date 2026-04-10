package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
	G "github.com/IBM/fp-go/v2/optics/lens/generic"
	TG "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

func Compose[S, A, B, HKTS, HKTA, HKTB any](
	fmap functor.MapType[A, S, HKTA, HKTS],
) func(Traversal[A, B, HKTA, HKTB]) func(Lens[S, A]) Traversal[S, B, HKTS, HKTB] {
	lensTrav := G.AsTraversal[Traversal[S, A, HKTS, HKTA]](fmap)

	return func(ab Traversal[A, B, HKTA, HKTB]) func(Lens[S, A]) Traversal[S, B, HKTS, HKTB] {
		return F.Flow2(
			lensTrav,
			TG.Compose[
				Traversal[A, B, HKTA, HKTB],
				Traversal[S, A, HKTS, HKTA],
				Traversal[S, B, HKTS, HKTB],
			](ab),
		)
	}
}
