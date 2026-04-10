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
	I "github.com/IBM/fp-go/v2/identity"
	AR "github.com/IBM/fp-go/v2/optics/traversal/array/generic"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

// FromArray creates a traversal for array elements using the Identity functor.
//
// This is a specialized version of the generic FromArray that uses the Identity
// functor, which provides the simplest possible computational context (no context).
// This makes it ideal for straightforward array transformations where you want to
// modify elements directly without additional effects.
//
// The Identity functor means that operations are applied directly to values without
// wrapping them in any additional structure. This results in clean, efficient
// traversals that simply map functions over array elements.
//
// Type Parameters:
//   - GA: Array type constraint (e.g., []A)
//   - A: The element type within the array
//
// Returns:
//   - A Traversal that can transform all elements in an array
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    T "github.com/IBM/fp-go/v2/optics/traversal"
//	    TI "github.com/IBM/fp-go/v2/optics/traversal/array/generic/identity"
//	)
//
//	// Create a traversal for integer arrays
//	arrayTraversal := TI.FromArray[[]int, int]()
//
//	// Compose with identity traversal
//	traversal := F.Pipe1(
//	    T.Id[[]int, []int](),
//	    T.Compose[[]int, []int, []int, int](arrayTraversal),
//	)
//
//	// Double all numbers in the array
//	numbers := []int{1, 2, 3, 4, 5}
//	doubled := traversal(func(n int) int { return n * 2 })(numbers)
//	// doubled: []int{2, 4, 6, 8, 10}
//
// See Also:
//   - AR.FromArray: Generic version with configurable functor
//   - I.Of: Identity functor's pure/of operation
//   - I.Map: Identity functor's map operation
//   - I.Ap: Identity functor's applicative operation
func FromArray[GA ~[]A, A any]() G.Traversal[GA, A, GA, A] {
	return AR.FromArray[GA](
		I.Of[GA],
		I.Map[GA, func(A) GA],
		I.Ap[GA, A],
	)
}

// At creates a function that focuses a traversal on a specific array index using the Identity functor.
//
// This is a specialized version of the generic At that uses the Identity functor,
// providing the simplest computational context for array element access. It transforms
// a traversal focusing on an array into a traversal focusing on the element at the
// specified index.
//
// The Identity functor means operations are applied directly without additional wrapping,
// making this ideal for straightforward element modifications. If the index is out of
// bounds, the traversal focuses on zero elements (no-op).
//
// Type Parameters:
//   - GA: Array type constraint (e.g., []A)
//   - S: The source type of the outer traversal
//   - A: The element type within the array
//
// Parameters:
//   - idx: The zero-based index to focus on
//
// Returns:
//   - A function that transforms a traversal on arrays into a traversal on a specific element
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    T "github.com/IBM/fp-go/v2/optics/traversal"
//	    TI "github.com/IBM/fp-go/v2/optics/traversal/array/generic/identity"
//	)
//
//	type Person struct {
//	    Name    string
//	    Hobbies []string
//	}
//
//	// Create a traversal focusing on hobbies
//	hobbiesTraversal := T.Id[Person, []string]()
//
//	// Focus on the second hobby (index 1)
//	secondHobby := F.Pipe1(
//	    hobbiesTraversal,
//	    TI.At[[]string, Person, string](1),
//	)
//
//	// Modify the second hobby
//	person := Person{Name: "Alice", Hobbies: []string{"reading", "coding", "gaming"}}
//	updated := secondHobby(func(s string) string {
//	    return s + "!"
//	})(person)
//	// updated.Hobbies: []string{"reading", "coding!", "gaming"}
//
//	// Out of bounds index is a no-op
//	outOfBounds := F.Pipe1(
//	    hobbiesTraversal,
//	    TI.At[[]string, Person, string](10),
//	)
//	unchanged := outOfBounds(func(s string) string {
//	    return s + "!"
//	})(person)
//	// unchanged.Hobbies: []string{"reading", "coding", "gaming"} (no change)
//
// See Also:
//   - AR.At: Generic version with configurable functor
//   - I.Of: Identity functor's pure/of operation
//   - I.Map: Identity functor's map operation
func At[GA ~[]A, S, A any](idx int) func(G.Traversal[S, GA, S, GA]) G.Traversal[S, A, S, A] {
	return AR.At[GA, S, A, S](
		I.Of[GA],
		I.Map[A, GA],
	)(idx)
}
