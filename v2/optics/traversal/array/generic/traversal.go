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
	"github.com/IBM/fp-go/v2/optics/optional"
	OA "github.com/IBM/fp-go/v2/optics/optional/array/generic"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

// FromArray returns a traversal from an array
func FromArray[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof pointed.OfType[GB, HKTRB],
	fmap functor.MapType[GB, func(B) GB, HKTRB, HKTAB],
	fap apply.ApType[HKTB, HKTRB, HKTAB],
) G.Traversal[GA, A, HKTRB, HKTB] {
	return func(f func(A) HKTB) func(GA) HKTRB {
		return AR.Traverse[GA](fof, fmap, fap, f)
	}
}

// At creates a function that focuses a traversal on a specific array index.
//
// This function takes an index and returns a function that transforms a traversal
// focusing on an array into a traversal focusing on the element at that index.
// It works by:
//  1. Creating an Optional that focuses on the array element at the given index
//  2. Converting that Optional into a Traversal
//  3. Composing it with the original traversal
//
// If the index is out of bounds, the traversal will focus on zero elements (no-op),
// following the Optional laws where operations on non-existent values have no effect.
//
// This is particularly useful when you have a nested structure containing arrays
// and want to traverse to a specific element within those arrays.
//
// Type Parameters:
//   - GA: Array type constraint (e.g., []A)
//   - S: The source type of the outer traversal
//   - A: The element type within the array
//   - HKTS: Higher-kinded type for S (functor/applicative context)
//   - HKTGA: Higher-kinded type for GA (functor/applicative context)
//   - HKTA: Higher-kinded type for A (functor/applicative context)
//
// Parameters:
//   - fof: Function to lift GA into the higher-kinded type HKTGA (pure/of operation)
//   - fmap: Function to map over HKTA and produce HKTGA (functor map operation)
//
// Returns:
//   - A function that takes an index and returns a traversal transformer
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/identity"
//	    T "github.com/IBM/fp-go/v2/optics/traversal"
//	    TA "github.com/IBM/fp-go/v2/optics/traversal/array/generic"
//	)
//
//	type Person struct {
//	    Name    string
//	    Hobbies []string
//	}
//
//	// Create a traversal focusing on the hobbies array
//	hobbiesTraversal := T.Id[Person, []string]()
//
//	// Focus on the first hobby (index 0)
//	firstHobby := F.Pipe1(
//	    hobbiesTraversal,
//	    TA.At[[]string, Person, string](
//	        identity.Of[[]string],
//	        identity.Map[string, []string],
//	    )(0),
//	)
//
//	// Modify the first hobby
//	person := Person{Name: "Alice", Hobbies: []string{"reading", "coding"}}
//	updated := firstHobby(func(s string) string {
//	    return s + "!"
//	})(person)
//	// updated.Hobbies: []string{"reading!", "coding"}
//
// See Also:
//   - OA.At: Creates an Optional focusing on an array element
//   - optional.AsTraversal: Converts an Optional to a Traversal
//   - G.Compose: Composes two traversals
func At[GA ~[]A, S, A, HKTS, HKTGA, HKTA any](
	fof pointed.OfType[GA, HKTGA],
	fmap functor.MapType[A, GA, HKTA, HKTGA],
) func(int) func(G.Traversal[S, GA, HKTS, HKTGA]) G.Traversal[S, A, HKTS, HKTA] {
	return F.Flow3(
		OA.At[GA],
		optional.AsTraversal[G.Traversal[GA, A, HKTGA, HKTA]](fof, fmap),
		G.Compose[
			G.Traversal[GA, A, HKTGA, HKTA],
			G.Traversal[S, GA, HKTS, HKTGA],
			G.Traversal[S, A, HKTS, HKTA],
		],
	)
}
