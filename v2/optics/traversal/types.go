package traversal

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Traversal[S, A, HKTS, HKTA any] = G.Traversal[S, A, HKTS, HKTA]

	Predicate[A any] = predicate.Predicate[A]
)
