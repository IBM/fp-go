package prism

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/optics/optional"
	P "github.com/IBM/fp-go/v2/optics/prism"
)

type (
	Prism[S, A any]     = P.Prism[S, A]
	Lens[S, A any]      = L.Lens[S, A]
	Optional[S, A any]  = O.Optional[S, A]
	Endomorphism[A any] = endomorphism.Endomorphism[A]
)
