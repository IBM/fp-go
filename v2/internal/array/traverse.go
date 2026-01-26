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
	F "github.com/IBM/fp-go/v2/function"
)

// MonadSequenceSegment sequences a segment of an array of effects using a divide-and-conquer approach.
// It recursively splits the array segment in half, sequences each half, and concatenates the results.
//
// This function is optimized for performance by using a divide-and-conquer strategy that reduces
// the depth of nested function calls compared to a linear fold approach.
//
// Type parameters:
//   - HKTB: The higher-kinded type containing values (e.g., Option[B], Either[E, B])
//   - HKTRB: The higher-kinded type containing an array of values (e.g., Option[[]B], Either[E, []B])
//
// Parameters:
//   - fof: Function to lift a single HKTB into HKTRB
//   - empty: The empty/identity value for HKTRB
//   - concat: Function to concatenate two HKTRB values
//   - fbs: The array of effects to sequence
//   - start: The starting index of the segment (inclusive)
//   - end: The ending index of the segment (exclusive)
//
// Returns:
//   - HKTRB: The sequenced result for the segment
//
// The function handles three cases:
//   - Empty segment (end - start == 0): returns empty
//   - Single element (end - start == 1): returns fof(fbs[start])
//   - Multiple elements: recursively divides and conquers
func MonadSequenceSegment[HKTB, HKTRB any](
	fof func(HKTB) HKTRB,
	empty HKTRB,
	concat func(HKTRB, HKTRB) HKTRB,
	fbs []HKTB,
	start, end int,
) HKTRB {

	switch end - start {
	case 0:
		return empty
	case 1:
		return fof(fbs[start])
	default:
		mid := (start + end) / 2
		return concat(
			MonadSequenceSegment(fof, empty, concat, fbs, start, mid),
			MonadSequenceSegment(fof, empty, concat, fbs, mid, end),
		)
	}
}

// SequenceSegment creates a function that sequences a segment of an array of effects.
// Unlike MonadSequenceSegment, this returns a curried function that can be reused.
//
// This function builds a computation tree at construction time, which can be more efficient
// when the same sequencing pattern needs to be applied multiple times to arrays of the same length.
//
// Type parameters:
//   - HKTB: The higher-kinded type containing values
//   - HKTRB: The higher-kinded type containing an array of values
//
// Parameters:
//   - fof: Function to lift a single HKTB into HKTRB
//   - empty: The empty/identity value for HKTRB
//   - concat: Function to concatenate two HKTRB values
//
// Returns:
//   - A function that takes an array of HKTB and returns HKTRB
func SequenceSegment[HKTB, HKTRB any](
	fof func(HKTB) HKTRB,
	empty HKTRB,
	concat func(HKTRB, HKTRB) HKTRB,
) func([]HKTB) HKTRB {

	concat_f := func(left, right func([]HKTB) HKTRB) func([]HKTB) HKTRB {
		return func(fbs []HKTB) HKTRB {
			return concat(left(fbs), right(fbs))
		}
	}
	empty_f := F.Constant1[[]HKTB](empty)
	at := func(idx int) func([]HKTB) HKTRB {
		return func(fbs []HKTB) HKTRB {
			return fof(fbs[idx])
		}
	}

	var divide func(start, end int) func([]HKTB) HKTRB
	divide = func(start, end int) func([]HKTB) HKTRB {
		switch end - start {
		case 0:
			return empty_f
		case 1:
			return at(start)
		default:
			mid := (start + end) / 2
			left := divide(start, mid)
			right := divide(mid, end)

			return concat_f(left, right)
		}
	}

	// TODO this could be cached by length
	get_divide := func(len int) func([]HKTB) HKTRB {
		return divide(0, len)
	}

	return func(fbs []HKTB) HKTRB {
		return get_divide(len(fbs))(fbs)
	}
}

// MonadTraverse maps each element of an array to an effect, then sequences the results.
// This is the monadic version that takes the array as a direct parameter.
//
// Traverse combines mapping and sequencing in one operation. It's useful when you want to
// transform each element of an array into an effect (like Option, Either, IO, etc.) and
// then collect all those effects into a single effect containing an array.
//
// We need to pass the members of the applicative explicitly, because golang does neither
// support higher kinded types nor template methods on structs or interfaces.
//
// Type parameters:
//   - GA: The input array type (e.g., []A)
//   - GB: The output array type (e.g., []B)
//   - A: The input element type
//   - B: The output element type
//   - HKTB: HKT<B> - The effect containing B (e.g., Option[B])
//   - HKTAB: HKT<func(B)GB> - Intermediate applicative type
//   - HKTRB: HKT<GB> - The effect containing the result array (e.g., Option[[]B])
//
// Parameters:
//   - fof: Function to lift a value into the effect (Of/Pure)
//   - fmap: Function to map over the effect (Map)
//   - fap: Function to apply an effect of a function to an effect of a value (Ap)
//   - ta: The input array to traverse
//   - f: The function to apply to each element, producing an effect
//
// Returns:
//   - HKTRB: An effect containing the array of transformed values
//
// Example:
//
//	If any element produces None, the entire result is None.
//	If all elements produce Some, the result is Some containing all values.
func MonadTraverse[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,
	f func(A) HKTB) HKTRB {
	return MonadTraverseReduce(fof, fmap, fap, ta, f, Append[GB, B], Empty[GB]())
}

// MonadTraverseWithIndex is like MonadTraverse but the transformation function also receives the index.
// This is useful when the transformation depends on the element's position in the array.
//
// Type parameters: Same as MonadTraverse
//
// Parameters:
//   - fof: Function to lift a value into the effect (Of/Pure)
//   - fmap: Function to map over the effect (Map)
//   - fap: Function to apply an effect of a function to an effect of a value (Ap)
//   - ta: The input array to traverse
//   - f: The function to apply to each element with its index, producing an effect
//
// Returns:
//   - HKTRB: An effect containing the array of transformed values
func MonadTraverseWithIndex[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,
	f func(int, A) HKTB) HKTRB {
	return MonadTraverseReduceWithIndex(fof, fmap, fap, ta, f, Append[GB, B], Empty[GB]())
}

// Traverse creates a curried function that maps each element to an effect and sequences the results.
// This is the curried version of MonadTraverse, useful for partial application and composition.
//
// Type parameters: Same as MonadTraverse
//
// Parameters:
//   - fof: Function to lift a value into the effect (Of/Pure)
//   - fmap: Function to map over the effect (Map)
//   - fap: Function to apply an effect of a function to an effect of a value (Ap)
//   - f: The function to apply to each element, producing an effect
//
// Returns:
//   - A function that takes an array and returns an effect containing the transformed array
func Traverse[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(A) HKTB) func(GA) HKTRB {

	return func(ma GA) HKTRB {
		return MonadTraverse(fof, fmap, fap, ma, f)
	}
}

// TraverseWithIndex creates a curried function like Traverse but with index-aware transformation.
// This is the curried version of MonadTraverseWithIndex.
//
// Type parameters: Same as MonadTraverse
//
// Parameters:
//   - fof: Function to lift a value into the effect (Of/Pure)
//   - fmap: Function to map over the effect (Map)
//   - fap: Function to apply an effect of a function to an effect of a value (Ap)
//   - f: The function to apply to each element with its index, producing an effect
//
// Returns:
//   - A function that takes an array and returns an effect containing the transformed array
func TraverseWithIndex[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(int, A) HKTB) func(GA) HKTRB {

	return func(ma GA) HKTRB {
		return MonadTraverseWithIndex(fof, fmap, fap, ma, f)
	}
}

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<GB>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func MonadSequence[GA ~[]HKTA, HKTA, HKTRA any](
	fof func(HKTA) HKTRA,
	empty HKTRA,
	concat func(HKTRA, HKTRA) HKTRA,

	ta GA) HKTRA {
	return MonadSequenceSegment(fof, empty, concat, ta, 0, len(ta))
}

func Sequence[GA ~[]HKTA, HKTA, HKTRA any](
	fof func(HKTA) HKTRA,
	empty HKTRA,
	concat func(HKTRA, HKTRA) HKTRA,
) func(GA) HKTRA {

	return func(ma GA) HKTRA {
		return MonadSequence(fof, empty, concat, ma)
	}
}

func MonadTraverseReduce[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := fmap(F.Curry2(reduce))

	return Reduce(ta, func(r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap,
			fap(transform(a)),
		)
	}, fof(initial))
}

func MonadTraverseReduceWithIndex[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,

	transform func(int, A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := fmap(F.Curry2(reduce))

	return ReduceWithIndex(ta, func(idx int, r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap,
			fap(transform(idx, a)),
		)
	}, fof(initial))
}

func TraverseReduce[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) func(GA) HKTRB {
	return func(ta GA) HKTRB {
		return MonadTraverseReduce(fof, fmap, fap, ta, transform, reduce, initial)
	}
}

// TraverseReduceWithIndex creates a curried function for index-aware custom reduction during traversal.
// This is the curried version of MonadTraverseReduceWithIndex.
//
// Type parameters: Same as MonadTraverseReduce
//
// Parameters: Same as TraverseReduce, except:
//   - transform: Function that takes index and element, producing an effect
//
// Returns:
//   - A function that takes an array and returns an effect containing the accumulated value
func TraverseReduceWithIndex[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	transform func(int, A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) func(GA) HKTRB {
	return func(ta GA) HKTRB {
		return MonadTraverseReduceWithIndex(fof, fmap, fap, ta, transform, reduce, initial)
	}
}
