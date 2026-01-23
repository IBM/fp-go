package builder

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

func MakeBuilder[S, A any](get func(S) Option[A], set func(A) Endomorphism[S], name string) Builder[S, A] {
	return Builder[S, A]{
		GetOption: get,
		Set:       set,
		name:      name,
	}
}

func ComposeLensPrism[S, A, B any](r Prism[A, B]) func(Lens[S, A]) Builder[S, B] {
	return func(l Lens[S, A]) Builder[S, B] {
		return MakeBuilder(
			F.Flow2(
				l.Get,
				r.GetOption,
			),
			F.Flow2(
				r.ReverseGet,
				l.Set,
			),
			fmt.Sprintf("Compose[%s -> %s]", l, r),
		)
	}
}
