package optional

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
)

func AsTraversal[R ~func(func(A) HKTA) func(S) HKTS, S, A, HKTS, HKTA any](
	fof pointed.OfType[S, HKTS],
	fmap functor.MapType[A, S, HKTA, HKTS],
) func(Optional[S, A]) R {
	return func(sa Optional[S, A]) R {
		return func(f func(A) HKTA) func(S) HKTS {
			return func(s S) HKTS {
				return F.Pipe2(
					s,
					sa.GetOption,
					O.Fold(
						lazy.Of(fof(s)),
						F.Flow2(
							f,
							fmap(func(a A) S {
								return sa.Set(a)(s)
							}),
						),
					),
				)
			}
		}
	}
}
