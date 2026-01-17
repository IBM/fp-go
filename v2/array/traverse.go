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
