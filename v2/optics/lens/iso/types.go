package iso

import (
	"github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/option"
)

type (
	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Iso represents an isomorphism - a bidirectional transformation between two types.
	Iso[S, A any] = common.Iso[S, A]

	// Lens is a functional reference to a subpart of a data structure.
	Lens[S, A any] = common.Lens[S, A]

	// Operator represents a function that transforms one lens into another.
	Operator[S, A, B any] = common.LensOperator[S, A, B]
)
