package generic

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/common"
	TG "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

type (
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Lens[S, A any] = common.Lens[S, A]

	Optional[S, A any] = common.Optional[S, A]

	// HKTES = HKT[Endomorphism[S]]
	Traversal[S, A, HKTES, HKTA any] = TG.Traversal[S, A, HKTES, HKTA]

	Monoid[T any] = monoid.Monoid[T]
)
