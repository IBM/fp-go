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

package traversal

import (
	C "github.com/IBM/fp-go/v2/constant"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

// Id is the identity constructor of a traversal
func Id[S, A any]() Traversal[S, S, A, A] {
	return F.Identity[func(S) A]
}

// Modify applies a transformation function to a traversal
func Modify[S, A any](f Endomorphism[A]) func(Traversal[S, A, S, A]) Endomorphism[S] {
	return identity.Flap[Endomorphism[S]](f)
}

// Set sets a constant value for all values of the traversal
func Set[S, A any](a A) func(Traversal[S, A, S, A]) Endomorphism[S] {
	return Modify[S](F.Constant1[A](a))
}

// FoldMap maps each target to a `Monoid` and combines the result
func FoldMap[S, M, A any](f func(A) M) func(sa Traversal[S, A, C.Const[M, S], C.Const[M, A]]) func(S) M {
	return G.FoldMap[S](f)
}

// Fold maps each target to a `Monoid` and combines the result
func Fold[S, A any](sa Traversal[S, A, C.Const[A, S], C.Const[A, A]]) func(S) A {
	return G.Fold(sa)
}

// GetAll gets all the targets of a traversal
func GetAll[A, S any](s S) func(sa Traversal[S, A, C.Const[[]A, S], C.Const[[]A, A]]) []A {
	return G.GetAll[[]A](s)
}

// Compose composes two traversables
func Compose[
	S, HKTS, A, B, HKTA, HKTB any](ab Traversal[A, B, HKTA, HKTB]) func(Traversal[S, A, HKTS, HKTA]) Traversal[S, B, HKTS, HKTB] {
	return G.Compose[
		Traversal[A, B, HKTA, HKTB],
		Traversal[S, A, HKTS, HKTA],
		Traversal[S, B, HKTS, HKTB]](ab)
}

// Filter creates a function that filters the targets of a traversal based on a predicate.
//
// This function allows you to refine a traversal to only focus on values that satisfy
// a given predicate. It works by converting the predicate into a prism, then converting
// that prism into a traversal, and finally composing it with the original traversal.
//
// The filtering is selective: when modifying values through the filtered traversal,
// only values that satisfy the predicate will be transformed. Values that don't
// satisfy the predicate remain unchanged.
//
// Type Parameters:
//   - S: The source type
//   - A: The focus type (the values being filtered)
//   - HKTS: Higher-kinded type for S (functor/applicative context)
//   - HKTA: Higher-kinded type for A (functor/applicative context)
//
// Parameters:
//   - fof: Function to lift A into the higher-kinded type HKTA (pure/of operation)
//   - fmap: Function to map over HKTA (functor map operation)
//
// Returns:
//   - A function that takes a predicate and returns an endomorphism on traversals
//
// Example:
//
//	import (
//	    AR "github.com/IBM/fp-go/v2/array"
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/identity"
//	    N "github.com/IBM/fp-go/v2/number"
//	    AI "github.com/IBM/fp-go/v2/optics/traversal/array/identity"
//	)
//
//	// Create a traversal for array elements
//	arrayTraversal := AI.FromArray[int]()
//	baseTraversal := F.Pipe1(
//	    Id[[]int, []int](),
//	    Compose[[]int, []int, []int, int](arrayTraversal),
//	)
//
//	// Filter to only positive numbers
//	isPositive := N.MoreThan(0)
//	filteredTraversal := F.Pipe1(
//	    baseTraversal,
//	    Filter[[]int, int](identity.Of[int], identity.Map[int, int])(isPositive),
//	)
//
//	// Double only positive numbers
//	numbers := []int{-2, -1, 0, 1, 2, 3}
//	result := filteredTraversal(func(n int) int { return n * 2 })(numbers)
//	// result: [-2, -1, 0, 2, 4, 6]
//
// See Also:
//   - prism.FromPredicate: Creates a prism from a predicate
//   - prism.AsTraversal: Converts a prism to a traversal
//   - Compose: Composes two traversals
func Filter[S, HKTS, A, HKTA any](
	fof pointed.OfType[A, HKTA],
	fmap functor.MapType[A, A, HKTA, HKTA],
) func(Predicate[A]) Endomorphism[Traversal[S, A, HKTS, HKTA]] {
	return G.Filter[S, HKTS](fof, fmap)
}
