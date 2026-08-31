package builder

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/common"
	"github.com/IBM/fp-go/v2/option"
)

type (
	Option[A any]   = option.Option[A]
	Lens[S, A any]  = common.Lens[S, A]
	Prism[S, A any] = common.Prism[S, A]

	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Builder[S, A any] struct {
		GetOption func(S) Option[A]

		Set func(A) Endomorphism[S]

		name string
	}

	Kleisli[S, A, B any]  = func(A) Builder[S, B]
	Operator[S, A, B any] = Kleisli[S, Builder[S, A], B]
)
