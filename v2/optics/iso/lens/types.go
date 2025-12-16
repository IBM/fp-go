// Package lens provides utilities for converting isomorphisms to lenses.
//
// This package bridges the gap between isomorphisms (bidirectional transformations)
// and lenses (focused accessors), allowing isomorphisms to be used wherever lenses
// are expected.
package lens

import (
	"github.com/IBM/fp-go/v2/optics/iso"
	L "github.com/IBM/fp-go/v2/optics/lens"
)

type (
	// Lens is a type alias for the standard lens from the optics/lens package.
	// A lens provides a composable way to focus on a field within a structure,
	// with operations to get and set values immutably.
	//
	// Type Parameters:
	//   - S: The source/structure type (the whole)
	//   - A: The focus/field type (the part)
	//
	// See github.com/IBM/fp-go/v2/optics/lens for full documentation.
	Lens[S, A any] = L.Lens[S, A]

	// Iso is a type alias for an isomorphism from the optics/iso package.
	// An isomorphism represents a bidirectional transformation between two types
	// without loss of information. It consists of two functions (Get and ReverseGet)
	// that are inverses of each other.
	//
	// Type Parameters:
	//   - S: The source type
	//   - A: The target type
	//
	// Isomorphisms can be converted to lenses using IsoAsLens or IsoAsLensRef,
	// which allows them to be used in lens compositions and operations.
	//
	// See github.com/IBM/fp-go/v2/optics/iso for full documentation.
	Iso[S, A any] = iso.Iso[S, A]
)
