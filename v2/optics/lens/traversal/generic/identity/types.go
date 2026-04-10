package identity

import (
	"github.com/IBM/fp-go/v2/optics/lens"
	T "github.com/IBM/fp-go/v2/optics/traversal"
)

type (

	// Lens is a functional reference to a subpart of a data structure.
	Lens[S, A any] = lens.Lens[S, A]

	Traversal[S, A, HKTS, HKTA any] = T.Traversal[S, A, HKTS, HKTA]
)
