package iso

import (
	"github.com/IBM/fp-go/v2/optics/iso"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
)

type (
	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Iso represents an isomorphism - a bidirectional transformation between two types.
	Iso[S, A any] = iso.Iso[S, A]

	// Lens is a functional reference to a subpart of a data structure.
	Lens[S, A any] = lens.Lens[S, A]

	// Operator represents a function that transforms one lens into another.
	Operator[S, A, B any] = lens.Operator[S, A, B]
)
