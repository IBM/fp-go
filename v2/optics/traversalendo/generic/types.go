package generic

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/optional"
	TG "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

type (
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Lens[S, A any] = lens.Lens[S, A]

	Optional[S, A any] = optional.Optional[S, A]

	// HKTS = HKT[Endomorphism[S]]
	Traversal[S, A, HKTS, HKTA any] = TG.Traversal[S, A, HKTS, HKTA]

	Monoid[T any] = monoid.Monoid[T]
)
