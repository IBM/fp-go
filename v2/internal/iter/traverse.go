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

package iter

import (
	F "github.com/IBM/fp-go/v2/function"
	INTA "github.com/IBM/fp-go/v2/internal/array"
	M "github.com/IBM/fp-go/v2/monoid"
)

// MonadTraverse traverses an iterator sequence, applying an effectful function to each element
// and collecting the results in an applicative context.
//
// This is a fundamental operation in functional programming that allows you to "turn inside out"
// a structure containing effects. It maps each element through a function that produces an effect,
// then sequences all those effects together while preserving the iterator structure.
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(A) bool)
//   - GB: The output iterator type ~func(yield func(B) bool)
//   - A: The input element type
//   - B: The output element type
//   - HKT_B: The higher-kinded type representing an effect containing B
//   - HKT_GB_GB: The higher-kinded type for a function from GB to GB in the effect context
//   - HKT_GB: The higher-kinded type representing an effect containing GB (the result iterator)
//
// Parameters:
//   - fmap_b: Maps a function over HKT_B to produce HKT_GB
//   - fof_gb: Lifts a GB value into the effect context (pure/of operation)
//   - fmap_gb: Maps a function over HKT_GB to produce HKT_GB_GB
//   - fap_gb: Applies an effectful function to an effectful value (ap operation)
//   - ta: The input iterator sequence to traverse
//   - f: The effectful function to apply to each element
//
// Returns:
//   - An effect containing an iterator of transformed elements
//
// Note: We need to pass the applicative operations explicitly because Go doesn't support
// higher-kinded types or template methods on structs/interfaces.
//
// Example (conceptual with Option):
//
//	// Traverse an iterator of strings, parsing each as an integer
//	// If any parse fails, the whole result is None
//	iter := func(yield func(string) bool) {
//	    yield("1")
//	    yield("2")
//	    yield("3")
//	}
//	result := MonadTraverse(..., iter, parseInt) // Some(iterator of [1,2,3]) or None
func MonadTraverse[GA ~func(yield func(A) bool), GB ~func(yield func(B) bool), A, B, HKT_B, HKT_GB_GB, HKT_GB any](
	fmap_b func(HKT_B, func(B) GB) HKT_GB,

	fof_gb OfType[GB, HKT_GB],
	fmap_gb func(HKT_GB, func(GB) func(GB) GB) HKT_GB_GB,
	fap_gb func(HKT_GB_GB, HKT_GB) HKT_GB,

	ta GA,
	f func(A) HKT_B) HKT_GB {

	fof := F.Bind2nd(fmap_b, Of[GB])

	empty := fof_gb(Empty[GB]())

	cb := F.Curry2(Concat[GB])
	concat_gb := F.Bind2nd(fmap_gb, cb)
	concat := func(first HKT_GB, second HKT_GB) HKT_GB {
		return fap_gb(concat_gb(first), second)
	}

	// convert to an array
	hktb := MonadMapToArray[GA, []HKT_B](ta, f)
	return INTA.MonadSequenceSegment(fof, empty, concat, hktb, 0, len(hktb))
}

// Traverse is the curried version of MonadTraverse, returning a function that traverses an iterator.
//
// This version uses type aliases for better readability and is more suitable for partial application
// and function composition. It returns a Kleisli arrow (a function from GA to HKT_GB).
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(A) bool)
//   - GB: The output iterator type ~func(yield func(B) bool)
//   - A: The input element type
//   - B: The output element type
//   - HKT_B: The higher-kinded type representing an effect containing B
//   - HKT_GB_GB: The higher-kinded type for a function from GB to GB in the effect context
//   - HKT_GB: The higher-kinded type representing an effect containing GB
//
// Parameters:
//   - fmap_b: Maps a function over HKT_B to produce HKT_GB
//   - fof_gb: Lifts a GB value into the effect context
//   - fmap_gb: Maps a function over HKT_GB to produce HKT_GB_GB
//   - fap_gb: Applies an effectful function to an effectful value
//   - f: The effectful function to apply to each element (Kleisli arrow)
//
// Returns:
//   - A function that takes an iterator and returns an effect containing an iterator of transformed elements
//
// Example (conceptual):
//
//	parseInts := Traverse[...](fmap, fof, fmap_gb, fap, parseInt)
//	iter := func(yield func(string) bool) { yield("1"); yield("2") }
//	result := parseInts(iter) // Effect containing iterator of integers
func Traverse[GA ~func(yield func(A) bool), GB ~func(yield func(B) bool), A, B, HKT_B, HKT_GB_GB, HKT_GB any](
	fmap_b MapType[B, GB, HKT_B, HKT_GB],

	fof_gb OfType[GB, HKT_GB],
	fmap_gb MapType[GB, Endomorphism[GB], HKT_GB, HKT_GB_GB],
	fap_gb ApType[HKT_GB, HKT_GB, HKT_GB_GB],

	f Kleisli[A, HKT_B]) Kleisli[GA, HKT_GB] {

	fof := fmap_b(Of[GB])
	empty := fof_gb(Empty[GB]())
	cb := F.Curry2(Concat[GB])
	concat_gb := fmap_gb(cb)

	concat := func(first, second HKT_GB) HKT_GB {
		return fap_gb(second)(concat_gb(first))
	}

	return F.Flow2(
		MapToArray[GA, []HKT_B](f),
		INTA.SequenceSegment(fof, empty, concat),
	)
}

// MonadSequence sequences an iterator of effects into an effect containing an iterator.
//
// This is a special case of traverse where the transformation function is the identity.
// It "flips" the nesting of the iterator and effect types, collecting all effects into
// a single effect containing an iterator of values.
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(HKTA) bool)
//   - HKTA: The higher-kinded type representing an effect containing A
//   - HKTRA: The higher-kinded type representing an effect containing an iterator of A
//
// Parameters:
//   - fof: Lifts an HKTA value into the HKTRA context
//   - m: A monoid for combining HKTRA values
//   - ta: The input iterator of effects to sequence
//
// Returns:
//   - An effect containing an iterator of values
//
// Example (conceptual with Option):
//
//	iter := func(yield func(Option[int]) bool) {
//	    yield(Some(1))
//	    yield(Some(2))
//	    yield(Some(3))
//	}
//	result := MonadSequence(..., iter) // Some(iterator of [1,2,3])
//
//	iter2 := func(yield func(Option[int]) bool) {
//	    yield(Some(1))
//	    yield(None)
//	}
//	result2 := MonadSequence(..., iter2) // None
func MonadSequence[GA ~func(yield func(HKTA) bool), HKTA, HKTRA any](
	fof OfType[HKTA, HKTRA],
	m M.Monoid[HKTRA],

	ta GA) HKTRA {

	// convert to an array
	hktb := ToArray[GA, []HKTA](ta)
	return INTA.MonadSequenceSegment(fof, m.Empty(), m.Concat, hktb, 0, len(hktb))
}

// MonadTraverseWithIndex traverses an iterator sequence with index tracking, applying an effectful
// function to each element along with its index.
//
// This is similar to MonadTraverse but the transformation function receives both the element's
// zero-based index and the element itself, useful when the position matters in the transformation.
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(A) bool)
//   - A: The input element type
//   - HKTB: The higher-kinded type representing an effect containing B
//   - HKTRB: The higher-kinded type representing an effect containing an iterator of B
//
// Parameters:
//   - fof: Lifts an HKTB value into the HKTRB context
//   - m: A monoid for combining HKTRB values
//   - ta: The input iterator sequence to traverse
//   - f: The effectful function that takes (index, element) and returns an effect
//
// Returns:
//   - An effect containing an iterator of transformed elements
//
// Example (conceptual):
//
//	iter := func(yield func(string) bool) {
//	    yield("a")
//	    yield("b")
//	}
//	// Add index prefix to each element
//	result := MonadTraverseWithIndex(..., iter, func(i int, s string) Effect[string] {
//	    return Pure(fmt.Sprintf("%d:%s", i, s))
//	}) // Effect containing iterator of ["0:a", "1:b"]
func MonadTraverseWithIndex[GA ~func(yield func(A) bool), A, HKTB, HKTRB any](
	fof func(HKTB) HKTRB,
	m M.Monoid[HKTRB],

	ta GA,
	f func(int, A) HKTB) HKTRB {

	// convert to an array
	hktb := MonadMapToArrayWithIndex[GA, []HKTB](ta, f)
	return INTA.MonadSequenceSegment(fof, m.Empty(), m.Concat, hktb, 0, len(hktb))
}

// Sequence is the curried version of MonadSequence, returning a function that sequences an iterator of effects.
//
// This version is more suitable for partial application and function composition.
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(HKTA) bool)
//   - HKTA: The higher-kinded type representing an effect containing A
//   - HKTRA: The higher-kinded type representing an effect containing an iterator of A
//
// Parameters:
//   - fof: Lifts an HKTA value into the HKTRA context
//   - m: A monoid for combining HKTRA values
//
// Returns:
//   - A function that takes an iterator of effects and returns an effect containing an iterator
//
// Example (conceptual):
//
//	sequenceOptions := Sequence[...](fof, monoid)
//	iter := func(yield func(Option[int]) bool) { yield(Some(1)); yield(Some(2)) }
//	result := sequenceOptions(iter) // Some(iterator of [1,2])
func Sequence[GA ~func(yield func(HKTA) bool), HKTA, HKTRA any](
	fof OfType[HKTA, HKTRA],
	m M.Monoid[HKTRA]) func(GA) HKTRA {

	return func(ma GA) HKTRA {
		return MonadSequence(fof, m, ma)
	}
}

// TraverseWithIndex is the curried version of MonadTraverseWithIndex, returning a function that
// traverses an iterator with index tracking.
//
// This version is more suitable for partial application and function composition.
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(A) bool)
//   - A: The input element type
//   - HKTB: The higher-kinded type representing an effect containing B
//   - HKTRB: The higher-kinded type representing an effect containing an iterator of B
//
// Parameters:
//   - fof: Lifts an HKTB value into the HKTRB context
//   - m: A monoid for combining HKTRB values
//   - f: The effectful function that takes (index, element) and returns an effect
//
// Returns:
//   - A function that takes an iterator and returns an effect containing an iterator of transformed elements
//
// Example (conceptual):
//
//	addIndexPrefix := TraverseWithIndex[...](fof, monoid, func(i int, s string) Effect[string] {
//	    return Pure(fmt.Sprintf("%d:%s", i, s))
//	})
//	iter := func(yield func(string) bool) { yield("a"); yield("b") }
//	result := addIndexPrefix(iter) // Effect containing iterator of ["0:a", "1:b"]
func TraverseWithIndex[GA ~func(yield func(A) bool), A, HKTB, HKTRB any](
	fof func(HKTB) HKTRB,
	m M.Monoid[HKTRB],

	f func(int, A) HKTB) func(GA) HKTRB {

	return func(ma GA) HKTRB {
		return MonadTraverseWithIndex(fof, m, ma, f)
	}
}

// MonadTraverseReduce combines traversal with reduction, applying an effectful transformation
// and accumulating results using a reducer function.
//
// This is a more efficient operation when you want to both transform elements through effects
// and reduce them to a single accumulated value, avoiding intermediate collections.
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(A) bool)
//   - GB: The accumulator type
//   - A: The input element type
//   - B: The transformed element type
//   - HKTB: The higher-kinded type representing an effect containing B
//   - HKTAB: The higher-kinded type for a function from B to GB in the effect context
//   - HKTRB: The higher-kinded type representing an effect containing GB
//
// Parameters:
//   - fof: Lifts a GB value into the effect context
//   - fmap: Maps a function over the effect to produce an effectful function
//   - fap: Applies an effectful function to an effectful value
//   - ta: The input iterator sequence to traverse and reduce
//   - transform: The effectful function to apply to each element
//   - reduce: The reducer function that combines the accumulator with a transformed element
//   - initial: The initial accumulator value
//
// Returns:
//   - An effect containing the final accumulated value
//
// Example (conceptual):
//
//	iter := func(yield func(string) bool) { yield("1"); yield("2"); yield("3") }
//	// Parse strings to ints and sum them
//	result := MonadTraverseReduce(..., iter, parseInt, add, 0)
//	// Returns: Some(6) or None if any parse fails
func MonadTraverseReduce[GA ~func(yield func(A) bool), GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := fmap(F.Curry2(reduce))

	return MonadReduce(ta, func(r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap,
			fap(transform(a)),
		)
	}, fof(initial))
}

// MonadTraverseReduceWithIndex combines indexed traversal with reduction, applying an effectful
// transformation that receives element indices and accumulating results using a reducer function.
//
// This is similar to MonadTraverseReduce but the transformation function also receives the
// zero-based index of each element, useful when position matters in the transformation logic.
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(A) bool)
//   - GB: The accumulator type
//   - A: The input element type
//   - B: The transformed element type
//   - HKTB: The higher-kinded type representing an effect containing B
//   - HKTAB: The higher-kinded type for a function from B to GB in the effect context
//   - HKTRB: The higher-kinded type representing an effect containing GB
//
// Parameters:
//   - fof: Lifts a GB value into the effect context
//   - fmap: Maps a function over the effect to produce an effectful function
//   - fap: Applies an effectful function to an effectful value
//   - ta: The input iterator sequence to traverse and reduce
//   - transform: The effectful function that takes (index, element) and returns an effect
//   - reduce: The reducer function that combines the accumulator with a transformed element
//   - initial: The initial accumulator value
//
// Returns:
//   - An effect containing the final accumulated value
//
// Example (conceptual):
//
//	iter := func(yield func(string) bool) { yield("a"); yield("b"); yield("c") }
//	// Create indexed strings and concatenate
//	result := MonadTraverseReduceWithIndex(..., iter,
//	    func(i int, s string) Effect[string] {
//	        return Pure(fmt.Sprintf("%d:%s", i, s))
//	    },
//	    func(acc, s string) string { return acc + "," + s },
//	    "")
//	// Returns: Effect containing "0:a,1:b,2:c"
func MonadTraverseReduceWithIndex[GA ~func(yield func(A) bool), GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,

	transform func(int, A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := fmap(F.Curry2(reduce))

	return MonadReduceWithIndex(ta, func(idx int, r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap,
			fap(transform(idx, a)),
		)
	}, fof(initial))
}

// TraverseReduce is the curried version of MonadTraverseReduce, returning a function that
// traverses and reduces an iterator.
//
// This version is more suitable for partial application and function composition.
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(A) bool)
//   - GB: The accumulator type
//   - A: The input element type
//   - B: The transformed element type
//   - HKTB: The higher-kinded type representing an effect containing B
//   - HKTAB: The higher-kinded type for a function from B to GB in the effect context
//   - HKTRB: The higher-kinded type representing an effect containing GB
//
// Parameters:
//   - fof: Lifts a GB value into the effect context
//   - fmap: Maps a function over the effect to produce an effectful function
//   - fap: Applies an effectful function to an effectful value
//   - transform: The effectful function to apply to each element
//   - reduce: The reducer function that combines the accumulator with a transformed element
//   - initial: The initial accumulator value
//
// Returns:
//   - A function that takes an iterator and returns an effect containing the accumulated value
//
// Example (conceptual):
//
//	sumParsedInts := TraverseReduce[...](fof, fmap, fap, parseInt, add, 0)
//	iter := func(yield func(string) bool) { yield("1"); yield("2"); yield("3") }
//	result := sumParsedInts(iter) // Some(6) or None if any parse fails
func TraverseReduce[GA ~func(yield func(A) bool), GB, A, B, HKTB, HKTAB, HKTRB any](
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

// TraverseReduceWithIndex is the curried version of MonadTraverseReduceWithIndex, returning a
// function that traverses and reduces an iterator with index tracking.
//
// This version is more suitable for partial application and function composition.
//
// Type Parameters:
//   - GA: The input iterator type ~func(yield func(A) bool)
//   - GB: The accumulator type
//   - A: The input element type
//   - B: The transformed element type
//   - HKTB: The higher-kinded type representing an effect containing B
//   - HKTAB: The higher-kinded type for a function from B to GB in the effect context
//   - HKTRB: The higher-kinded type representing an effect containing GB
//
// Parameters:
//   - fof: Lifts a GB value into the effect context
//   - fmap: Maps a function over the effect to produce an effectful function
//   - fap: Applies an effectful function to an effectful value
//   - transform: The effectful function that takes (index, element) and returns an effect
//   - reduce: The reducer function that combines the accumulator with a transformed element
//   - initial: The initial accumulator value
//
// Returns:
//   - A function that takes an iterator and returns an effect containing the accumulated value
//
// Example (conceptual):
//
//	concatIndexed := TraverseReduceWithIndex[...](fof, fmap, fap,
//	    func(i int, s string) Effect[string] {
//	        return Pure(fmt.Sprintf("%d:%s", i, s))
//	    },
//	    func(acc, s string) string { return acc + "," + s },
//	    "")
//	iter := func(yield func(string) bool) { yield("a"); yield("b") }
//	result := concatIndexed(iter) // Effect containing "0:a,1:b"
func TraverseReduceWithIndex[GA ~func(yield func(A) bool), GB, A, B, HKTB, HKTAB, HKTRB any](
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
