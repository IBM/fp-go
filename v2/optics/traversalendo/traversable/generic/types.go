package generic

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/internal/traversable"
	"github.com/IBM/fp-go/v2/optics/common"
	TG "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

type (
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Lens[S, A any] = common.Lens[S, A]

	Optional[S, A any] = common.Optional[S, A]

	// HKTS = HKT[Endomorphism[S]]
	Traversal[S, A, HKTS, HKTA any] = TG.Traversal[S, A, HKTS, HKTA]

	Traversable[A, HKT_F_B, HKT_T_A, HKT_F_T_B any] = traversable.Traversable[A, HKT_F_B, HKT_T_A, HKT_F_T_B]
)
