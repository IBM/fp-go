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

package option

import (
	F "github.com/IBM/fp-go/v2/function"
)

// AltAllArray combines multiple Options from an array using the Alt operation.
// It starts with an initial Option and iteratively applies Alt with each Option
// in the array, returning the first Some value encountered or the last value if
// all are None.
//
// The Alt operation returns the first Option if it's Some, otherwise returns the
// alternative. This function chains multiple Alt operations together, effectively
// implementing a "first success" or "fallback chain" pattern.
//
// Implementation:
//
// This function is semantically equivalent to alt.AltAllArray[Option[A]](Alt)(startWith)
// but provides an optimized implementation with early break for better performance.
// The direct implementation avoids the overhead of lazy thunk allocations while
// maintaining identical behavior.
//
// Short-Circuit Behavior:
//
// This function short-circuits on the first Some value:
//   - If startWith is Some, it returns immediately without examining the array
//   - Otherwise, it iterates through the array and returns the first Some encountered
//   - The array is not fully consumed once a Some value is found
//   - Only if all elements are None will the entire array be traversed
//
// Relationship to array.Fold and AltMonoid:
//
// AltAllArray is closely related to array.Fold with AltMonoid:
//   - When startWith is None, AltAllArray(None[A]())(options) is equivalent to
//     array.Fold(AltMonoid[A]())(options)
//   - When startWith is Some, it's equivalent to prepending startWith to the array
//     before folding: array.Fold(AltMonoid[A]())(append([]Option[A]{startWith}, options...))
//   - AltMonoid is a monoid that uses Alt as its Concat operation and None as Empty
//   - Both approaches have O(n) time complexity and similar performance
//
// Use AltAllArray when you want to:
//   - Express the "find first Some value" pattern clearly
//   - Specify a custom starting value (not just None)
//   - Work specifically with Options in a functional style
//
// Type Parameters:
//   - A: The type of value contained in the Options
//
// Parameters:
//   - startWith: The initial Option to start the chain with
//
// Returns:
//   - Kleisli[[]Option[A], A]: A function that takes an array of Options and
//     returns an Option containing the first Some value, or the last value if
//     all are None
//
// Example:
//
//	options := []Option[int]{
//	    None[int](),
//	    None[int](),
//	    Some(42),
//	    Some(100),
//	}
//	result := AltAllArray(None[int]())(options) // Some(42)
//
//	allNone := []Option[int]{
//	    None[int](),
//	    None[int](),
//	}
//	result := AltAllArray(Some(10))(allNone) // Some(10)
//
// See Also:
//   - Alt: The underlying Alt operation
//   - AltAllSeq: Similar function for iterator sequences
//   - AltMonoid: Monoid that uses Alt operation
func AltAllArray[A any](startWith Option[A]) Kleisli[[]Option[A], A] {
	// Direct first-Some scan: option's Alt keeps the current value if it is Some,
	// otherwise takes the next one. Folding via the generic lazy Alt would wrap
	// every element in an F.Constant thunk (a closure allocation per element);
	// this loop is allocation-free while preserving identical semantics
	// (startWith priority; first Some in the array; else the last value).
	if IsSome(startWith) {
		return F.Constant1[[]Option[A]](startWith)
	}
	return func(as []Option[A]) Option[A] {
		for _, o := range as {
			if IsSome(o) {
				return o
			}
		}
		return startWith
	}
}

// AltAllSeq combines multiple Options from an iterator sequence using the Alt operation.
// It starts with an initial Option and iteratively applies Alt with each Option
// from the sequence, returning the first Some value encountered or the last value
// if all are None.
//
// This function is similar to AltAllArray but works with Go's iterator sequences,
// making it suitable for lazy evaluation and potentially infinite sequences.
//
// Implementation:
//
// This function is semantically equivalent to alt.AltAllSeq[Option[A]](Alt)(startWith)
// but provides an optimized implementation with early break for better performance.
// The direct implementation avoids the overhead of lazy thunk allocations while
// maintaining identical behavior.
//
// Short-Circuit Behavior:
//
// This function short-circuits on the first Some value:
//   - If startWith is Some, it returns immediately without consuming the sequence
//   - Otherwise, it iterates through the sequence and returns the first Some encountered
//   - The sequence is not fully consumed once a Some value is found
//   - This makes it safe to use with infinite sequences as long as a Some value exists
//   - Only if all elements are None will the entire sequence be consumed
//
// Relationship to Folding:
//
// Like AltAllArray, this function implements a fold operation using the Alt operation.
// The key difference is that it works with iterator sequences instead of arrays,
// enabling:
//   - Lazy evaluation of the sequence
//   - Working with potentially infinite sequences
//   - Memory-efficient processing of large datasets
//   - Composition with other iterator-based operations
//
// The relationship to AltMonoid is the same as AltAllArray, but applied to sequences
// rather than arrays.
//
// Type Parameters:
//   - A: The type of value contained in the Options
//
// Parameters:
//   - startWith: The initial Option to start the chain with
//
// Returns:
//   - Kleisli[Seq[Option[A]], A]: A function that takes a sequence of Options and
//     returns an Option containing the first Some value, or the last value if
//     all are None
//
// Example:
//
//	generator := func(yield func(Option[int]) bool) {
//	    yield(None[int]())
//	    yield(Some(42))
//	    yield(Some(100))
//	}
//	result := AltAllSeq(None[int]())(generator) // Some(42)
//
//	emptyGen := func(yield func(Option[int]) bool) {}
//	result := AltAllSeq(Some(10))(emptyGen) // Some(10)
//
// See Also:
//   - Alt: The underlying Alt operation
//   - AltAllArray: Similar function for arrays
//   - AltMonoid: Monoid that uses Alt operation
func AltAllSeq[A any](startWith Option[A]) Kleisli[Seq[Option[A]], A] {
	if IsSome(startWith) {
		return F.Constant1[Seq[Option[A]]](startWith)
	}
	return func(as Seq[Option[A]]) Option[A] {
		for o := range as {
			if IsSome(o) {
				return o
			}
		}
		return startWith
	}
}
