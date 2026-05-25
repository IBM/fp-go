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

package array

import (
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
	"github.com/IBM/fp-go/v2/internal/traversable"
)

// Traverse maps each element of an array to an effect (HKT), then collects the results
// into an effect of an array. This is like a combination of Map and Sequence.
//
// Unlike Sequence which works with []HKT<A> -> HKT<[]A>, Traverse works with
// []A -> (A -> HKT<B>) -> HKT<[]B>, allowing you to transform elements while sequencing effects.
//
// Type parameters:
//   - HKTB = HKT<B> (e.g., Option[B], Either[E, B])
//   - HKTAB = HKT<func(B)[]B> (intermediate type for applicative)
//   - HKTRB = HKT<[]B> (e.g., Option[[]B], Either[E, []B])
//
// Example:
//
//	import (
//	    "github.com/IBM/fp-go/v2/option"
//	    "strconv"
//	)
//
//	// Parse strings to ints, returning None if any parse fails
//	parseAll := array.Traverse(
//	    option.Of[[]int],
//	    option.Map[[]int, func(int) []int],
//	    option.Ap[[]int, int],
//	    func(s string) option.Option[int] {
//	        if n, err := strconv.Atoi(s); err == nil {
//	            return option.Some(n)
//	        }
//	        return option.None[int]()
//	    },
//	)
//
//	result := parseAll([]string{"1", "2", "3"}) // Some([1, 2, 3])
//	result2 := parseAll([]string{"1", "x", "3"}) // None
//
//go:inline
func Traverse[A, B, HKTB, HKTAB, HKTRB any](
	fof pointed.OfType[[]B, HKTRB],
	fmap functor.MapType[[]B, func(B) []B, HKTRB, HKTAB],
	fap apply.ApType[HKTB, HKTRB, HKTAB],

	f func(A) HKTB) func([]A) HKTRB {
	return array.Traverse[[]A](fof, fmap, fap, f)
}

// MonadTraverse is the monadic version of Traverse that takes the array as a parameter.
// It maps each element of an array to an effect (HKT), then collects the results
// into an effect of an array.
//
// This is useful when you want to apply the traverse operation directly without currying.
//
//go:inline
func MonadTraverse[A, B, HKTB, HKTAB, HKTRB any](
	fof pointed.OfType[[]B, HKTRB],
	fmap func(func([]B) func(B) []B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta []A,
	f func(A) HKTB) HKTRB {

	return array.MonadTraverse(fof, fmap, fap, ta, f)
}

//go:inline
func TraverseWithIndex[A, B, HKTB, HKTAB, HKTRB any](
	fof pointed.OfType[[]B, HKTRB],
	fmap func(func([]B) func(B) []B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(int, A) HKTB) func([]A) HKTRB {
	return array.TraverseWithIndex[[]A](fof, fmap, fap, f)
}

//go:inline
func MonadTraverseWithIndex[A, B, HKTB, HKTAB, HKTRB any](
	fof pointed.OfType[[]B, HKTRB],
	fmap func(func([]B) func(B) []B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta []A,
	f func(int, A) HKTB) HKTRB {

	return array.MonadTraverseWithIndex(fof, fmap, fap, ta, f)
}

func MakeTraverseType[A, B, HKT_F_B, HKT_F_T_B, HKT_F_B_T_B any]() traversable.TraverseType[A, B, []A, []B, HKT_F_B, HKT_F_T_B, HKT_F_B_T_B] {
	return func(
		// ap
		fof_b pointed.OfType[[]B, HKT_F_T_B],
		fmap_b functor.MapType[[]B, func(B) []B, HKT_F_T_B, HKT_F_B_T_B],
		fap_b apply.ApType[HKT_F_B, HKT_F_T_B, HKT_F_B_T_B],

	) func(func(A) HKT_F_B) func([]A) HKT_F_T_B {
		return func(f func(A) HKT_F_B) func([]A) HKT_F_T_B {
			return Traverse(
				fof_b,
				fmap_b,
				fap_b,
				f,
			)
		}
	}
}

// MakeTraversable creates a fully curried traversal function that can be specialized for different effects.
// This is the most abstract form of traverse, returning a function that takes the transformation
// function and then the array, enabling maximum composability and reusability.
//
// This function is useful when you want to create a reusable traversal operation that can be
// applied to different transformation functions and arrays. It's particularly valuable in
// point-free style programming and when building traversal pipelines.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type
//   - HKTB: The effect containing B (e.g., Option[B], Either[E, B])
//   - HKTAB: Intermediate applicative type for function application
//   - HKTRB: The effect containing the result array (e.g., Option[[]B])
//
// Parameters:
//   - fof: Function to lift a value into the effect (Of/Pure from the applicative)
//   - fmap: Function to map over the effect (Map from the functor)
//   - fap: Function to apply an effect of a function to an effect of a value (Ap from the applicative)
//
// Returns:
//   - A function that takes a transformation function and returns a function that takes an array
//     and returns an effect containing the transformed array
//
// Example:
//
//	import (
//	    A "github.com/IBM/fp-go/v2/array"
//	    O "github.com/IBM/fp-go/v2/option"
//	    R "github.com/IBM/fp-go/v2/result"
//	    "strconv"
//	)
//
//	// Create a reusable Option traversal for string arrays
//	traverseOption := A.MakeTraversable[string, int](
//	    O.Of[[]int],
//	    O.Map[[]int, func(int) []int],
//	    O.Ap[int, []int],
//	)
//
//	// Use it with different transformation functions
//	parseInts := traverseOption(R.Eitherize1(strconv.Atoi))
//	result1 := parseInts([]string{"1", "2", "3"})  // Some([1, 2, 3])
//	result2 := parseInts([]string{"1", "x", "3"})  // None
//
// See Also:
//   - Traverse: Non-curried version that takes the transformation function directly
//   - MonadTraverse: Version that takes the array as a direct parameter
func MakeTraversable[A, B, HKTB, HKTAB, HKTRB any](
	fof func([]B) HKTRB,
	fmap func(func([]B) func(B) []B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,
) func(func(A) HKTB) func([]A) HKTRB {
	return array.MakeTraversable[[]A](fof, fmap, fap)
}
