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
	"github.com/IBM/fp-go/v2/internal/apply"
	AR "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
	TG "github.com/IBM/fp-go/v2/optics/traversalendo/traversable/generic"
)

// FromArrayLens creates a traversal endomorphism from a lens that focuses on an array.
//
// This function converts a lens that accesses an array field into a traversal endomorphism
// that can traverse all elements within that array. The resulting traversal works with
// Endomorphism[S] as its effect type, making it composable with other traversal endomorphisms
// using monoid operations.
//
// The function is FromTraversableLens specialized to the array traversable built from
// fof, fmap and fap:
//  1. The lens getter extracts the array from the source structure
//  2. The array traversable applies the transformation function to every element, left to right
//  3. The lens setter, lifted with fmapEndo, turns the new array into an Endomorphism[S]
//
// This is particularly useful when building complex traversals that need to access and
// modify array fields within a larger structure. The endomorphism-based approach allows
// you to combine multiple such traversals using Concat and Empty operations.
//
// Type Parameters:
//   - GA: Array type (e.g., []A)
//   - S: The source structure type containing the array
//   - A: The element type within the array
//   - HKTS: Higher-kinded type for Endomorphism[S]
//   - HKTA: Higher-kinded type for A in the effect context
//   - HKTAA: Higher-kinded type for the partially applied array builder func(A) GA
//   - HKTRA: Higher-kinded type for the array GA in the effect context
//
// Parameters:
//   - fof: Function to lift GA into HKTRA (pure/of operation)
//   - fmap: Function to map the curried array builder over HKTRA
//   - fmapEndo: Function to map the lens setter into the endomorphism context
//   - fap: Applicative apply operation for combining effects
//
// Returns:
//   - A function that takes a lens and returns a traversal endomorphism
//
// Example:
//
//	import (
//	    A "github.com/IBM/fp-go/v2/array"
//	    thunk "github.com/IBM/fp-go/v2/context/readerioresult"
//	    "github.com/IBM/fp-go/v2/endomorphism"
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/monoid"
//	    "github.com/IBM/fp-go/v2/optics/lens"
//	    TLA "github.com/IBM/fp-go/v2/optics/traversalendo/array/generic"
//	    TE "github.com/IBM/fp-go/v2/optics/traversalendo/generic"
//	    TLL "github.com/IBM/fp-go/v2/optics/traversalendo/lens/generic"
//	)
//
//	type Person struct {
//	    Name    string
//	    Hobbies []string
//	}
//
//	// Create a lens for the Hobbies field
//	hobbiesLens := lens.MakeLens(
//	    func(p Person) []string { return p.Hobbies },
//	    func(p Person, h []string) Person {
//	        p.Hobbies = h
//	        return p
//	    },
//	)
//
//	// Convert to a traversal endomorphism over all hobbies
//	hobbiesTrav := TLA.FromArrayLens[[]string, Person, string](
//	    thunk.Of[[]string],
//	    thunk.Map[[]string, func(string) []string],
//	    thunk.Map[[]string, endomorphism.Endomorphism[Person]],
//	    thunk.Ap[[]string, string],
//	)(hobbiesLens)
//
//	// Use in a monoid fold to combine with other traversal endomorphisms
//	m := TE.MakeMonoid[string, thunk.ReaderIOResult[string], Person](
//	    thunk.Of,
//	    thunk.Map,
//	    thunk.Ap,
//	)
//	combined := F.Pipe1(
//	    A.From(hobbiesTrav, TLL.FromLens[Person, string](thunk.Map)(nameLens)),
//	    monoid.Fold(m),
//	)
//
// See Also:
//   - traversable/generic.FromTraversableLens: The generalization to any traversable structure
//   - optional/generic.FromArrayOptional: The variant for arrays focused by an optional
//   - generic.ToTraversal: Convert a traversal endomorphism to a regular traversal
//   - generic.MakeMonoid: Create a monoid for combining traversal endomorphisms
func FromArrayLens[GA ~[]A, S, A, HKTS, HKTA, HKTAA, HKTRA any](
	fof pointed.OfType[GA, HKTRA],
	fmap functor.MapType[GA, func(A) GA, HKTRA, HKTAA],
	fmapEndo functor.MapType[GA, Endomorphism[S], HKTRA, HKTS],
	fap apply.ApType[HKTA, HKTRA, HKTAA],
) func(Lens[S, GA]) Traversal[S, A, HKTS, HKTA] {
	return TG.FromTraversableLens[A, HKTA](fmapEndo)(AR.MakeTraversable[GA](fof, fmap, fap))
}
