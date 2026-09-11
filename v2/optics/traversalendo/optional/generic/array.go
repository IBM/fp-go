// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/apply"
	AR "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
	TG "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

// FromArrayOptional creates a traversal endomorphism from an optional that focuses on an array.
//
// This function is analogous to FromArrayLens but works with an Optional instead of a Lens.
// When the optional value exists (Some), it traverses all elements in the array using the
// provided transformation function. When the optional value does not exist (None), it returns
// the identity endomorphism (no-op), leaving the structure unchanged.
//
// The function is derived by composition:
//  1. FromOptional converts the optional into a traversal that focuses on the whole array
//  2. The array traversable (built from fof, fmap and fap) focuses on each element of that array
//  3. Composing both traversals yields a traversal that focuses on every element of the array
//
// Type Parameters:
//   - GA: Array type constraint (e.g., []A)
//   - S: The source structure type containing the optional array
//   - A: The element type within the array
//   - HKTS: Higher-kinded type for Endomorphism[S]
//   - HKTA: Higher-kinded type for A in the effect context
//   - HKTAA: Higher-kinded type for the array transformation function
//   - HKTRA: Higher-kinded type for the array in the effect context
//
// Parameters:
//   - fof: Function to lift GA into HKTRA (pure/of operation)
//   - fmap: Function to map over the array transformation
//   - fofEndo: Function to lift Endomorphism[S] into HKTS (used for the None branch)
//   - fmapEndo: Function to map the optional setter into the endomorphism context
//   - fap: Applicative apply operation for combining effects
//
// Returns:
//   - A function that takes an Optional and returns a traversal endomorphism
//
// See Also:
//   - FromArrayLens: The lens-based variant for non-optional array fields
//   - FromOptional: For converting a plain optional to a traversal
func FromArrayOptional[GA ~[]A, S, A, HKTS, HKTA, HKTAA, HKTRA any](
	fof pointed.OfType[GA, HKTRA],
	fmap functor.MapType[GA, func(A) GA, HKTRA, HKTAA],
	fofEndo pointed.OfType[Endomorphism[S], HKTS],
	fmapEndo functor.MapType[GA, Endomorphism[S], HKTRA, HKTS],
	fap apply.ApType[HKTA, HKTRA, HKTAA],
) func(Optional[S, GA]) Traversal[S, A, HKTS, HKTA] {
	return F.Flow2(
		FromOptional(fofEndo, fmapEndo),
		TG.Compose[
			Traversal[GA, A, HKTRA, HKTA],
			Traversal[S, GA, HKTS, HKTRA],
			Traversal[S, A, HKTS, HKTA],
		](AR.MakeTraversable[GA](fof, fmap, fap)),
	)
}
