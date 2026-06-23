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

package either

import (
	F "github.com/IBM/fp-go/v2/function"
)

// AltAllArray combines multiple Either values from an array using the Alt operation.
// It starts with an initial Either and iteratively applies Alt with each Either
// in the array, returning the first Right value encountered or the last Left value
// if all are Left.
//
// The Alt operation returns the first Either if it's Right, otherwise returns the
// alternative. This function chains multiple Alt operations together, effectively
// implementing a "first success" or "fallback chain" pattern for Either values.
//
// Implementation:
//
// This function is semantically equivalent to alt.AltAllArray[Either[E, A]](Alt)(startWith)
// but uses the generic implementation directly. The generic alt.AltAllArray provides
// lazy evaluation through thunks, which enables early break when a Right value is found.
//
// Short-Circuit Behavior:
//
// This function short-circuits on the first Right value but processes all Left values:
//   - If startWith is Right, it returns immediately without examining the array
//   - When iterating, it returns immediately upon finding the first Right value
//   - The array is not fully consumed once a Right value is found
//   - If all elements are Left, the entire array is traversed and the last Left is returned
//
// Relationship to array.Fold and AltMonoid:
//
// AltAllArray is closely related to array.Fold with AltMonoid:
//   - When startWith is Left, AltAllArray(Left[E, A](e))(eithers) is equivalent to
//     array.Fold(AltMonoid[E, A]())(eithers)
//   - When startWith is Right, it's equivalent to prepending startWith to the array
//     before folding: array.Fold(AltMonoid[E, A]())(append([]Either[E, A]{startWith}, eithers...))
//   - AltMonoid is a monoid that uses Alt as its Concat operation and Left as Empty
//   - Both approaches have O(n) time complexity and similar performance
//
// Use AltAllArray when you want to:
//   - Express the "find first Right value" pattern clearly
//   - Specify a custom starting value (not just Left)
//   - Work specifically with Either values in a functional style
//   - Implement fallback chains for error handling
//
// Type Parameters:
//   - E: The type of error/Left value
//   - A: The type of success/Right value
//
// Parameters:
//   - startWith: The initial Either to start the chain with
//
// Returns:
//   - Kleisli[E, []Either[E, A], A]: A function that takes an array of Either values and
//     returns an Either containing the first Right value, or the last Left value
//     if all are Left
//
// See Also:
//   - Alt: The underlying Alt operation
//   - AltAllSeq: Similar function for iterator sequences
//   - AltMonoid: Monoid that uses Alt operation
func AltAllArray[E, A any](startWith Either[E, A]) Kleisli[E, []Either[E, A], A] {
	if IsRight(startWith) {
		return F.Constant1[[]Either[E, A]](startWith)
	}
	return func(as []Either[E, A]) Either[E, A] {
		last := startWith
		for _, e := range as {
			if IsRight(e) {
				return e
			}
			last = e
		}
		return last
	}
}

// AltAllSeq combines multiple Either values from an iterator sequence using the Alt operation.
// It starts with an initial Either and iteratively applies Alt with each Either
// from the sequence, returning the first Right value encountered or the last Left value
// if all are Left.
//
// This function is similar to AltAllArray but works with Go's iterator sequences,
// making it suitable for lazy evaluation and potentially infinite sequences.
//
// Implementation:
//
// This function is semantically equivalent to alt.AltAllSeq[Either[E, A]](Alt)(startWith)
// but uses the generic implementation directly. The generic alt.AltAllSeq provides
// lazy evaluation through thunks, which enables early break when a Right value is found.
//
// Short-Circuit Behavior:
//
// This function short-circuits on the first Right value but processes all Left values:
//   - If startWith is Right, it returns immediately without consuming the sequence
//   - When iterating, it returns immediately upon finding the first Right value
//   - The sequence is not fully consumed once a Right value is found
//   - This makes it safe to use with infinite sequences as long as a Right value exists
//   - If all elements are Left, the entire sequence is consumed and the last Left is returned
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
//   - E: The type of error/Left value
//   - A: The type of success/Right value
//
// Parameters:
//   - startWith: The initial Either to start the chain with
//
// Returns:
//   - Kleisli[E, Seq[Either[E, A]], A]: A function that takes a sequence of Either values and
//     returns an Either containing the first Right value, or the last Left value
//     if all are Left
//
// See Also:
//   - Alt: The underlying Alt operation
//   - AltAllArray: Similar function for arrays
//   - AltMonoid: Monoid that uses Alt operation
func AltAllSeq[E, A any](startWith Either[E, A]) Kleisli[E, Seq[Either[E, A]], A] {
	if IsRight(startWith) {
		return F.Constant1[Seq[Either[E, A]]](startWith)
	}
	return func(as Seq[Either[E, A]]) Either[E, A] {
		last := startWith
		for e := range as {
			if IsRight(e) {
				return e
			}
			last = e
		}
		return last
	}
}
