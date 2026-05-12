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

package itereither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/iterator/iter"
)

// FilterOrElse filters a SeqEither value based on a predicate.
// If the predicate returns true for the Right value, it passes through unchanged.
// If the predicate returns false, it transforms the Right value into a Left using onFalse.
// Left values are passed through unchanged.
//
// This is useful for adding validation or constraints to successful computations,
// converting values that don't meet certain criteria into errors.
//
// Marble diagram:
//
//	Input:  ---R(5)---R(-3)---L(e)---R(10)---|
//	pred(x) = x > 0
//	onFalse(x) = "negative: " + x
//	Output: ---R(5)---L("negative: -3")---L(e)---R(10)---|
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
// Values that fail the predicate are converted to Left.
//
// Parameters:
//   - pred: A predicate function that tests the Right value
//   - onFalse: A function that converts the failing value into an error of type E
//
// Returns:
//   - An Operator that filters SeqEither values based on the predicate
//
// Example:
//
//	// Validate that a number is positive
//	isPositive := N.MoreThan(0)
//	onNegative := S.Format[int]("%d is not positive")
//
//	validatePositive := itereither.FilterOrElse(isPositive, onNegative)
//
//	result1 := validatePositive(itereither.Right[string](42))  // Right(42)
//	result2 := validatePositive(itereither.Right[string](-5))  // Left("-5 is not positive")
//	result3 := validatePositive(itereither.Left[int]("error")) // Left("error")
//
//go:inline
func FilterOrElse[E, A any](pred Predicate[A], onFalse func(A) E) Operator[E, A, A] {
	return ChainEitherK(either.FromPredicate(pred, onFalse))
}

// MonadFilter filters a SeqEither sequence, keeping only Right values that satisfy the predicate.
// Left values are always passed through unchanged, regardless of the predicate.
//
// This function processes each element in the sequence:
//   - If the element is Left, it passes through unchanged
//   - If the element is Right and satisfies the predicate, it passes through
//   - If the element is Right but fails the predicate, it is filtered out
//
// Unlike FilterOrElse, which converts failing Right values to Left, MonadFilter
// simply removes them from the sequence entirely.
//
// Marble diagram:
//
//	Input:  --R(1)--R(2)--L(e)--R(3)--R(4)--R(5)-->
//	pred(x) = x % 2 == 0
//	Output: --------R(2)--L(e)--------R(4)-------->
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
// Odd numbers are filtered out, even numbers and errors pass through.
//
// Type Parameters:
//   - E: The type of the Left value (error type)
//   - A: The type of the Right value (success type)
//
// Parameters:
//   - as: The input SeqEither sequence to filter
//   - pred: A predicate function that tests Right values
//
// Returns:
//   - A SeqEither containing only Left values and Right values that satisfy the predicate
//
// Example:
//
//	seq := iter.From(
//	    E.Right[string](1),
//	    E.Right[string](2),
//	    E.Left[int]("error"),
//	    E.Right[string](3),
//	    E.Right[string](4),
//	)
//	isEven := func(x int) bool { return x%2 == 0 }
//	result := itereither.MonadFilter(seq, isEven)
//	// yields: Right(2), Left("error"), Right(4)
//
// See Also:
//
// Filter is the curried version of MonadFilter.
// FilterOrElse converts failing values to Left instead of removing them.
func MonadFilter[E, A any](as SeqEither[E, A], pred Predicate[A]) SeqEither[E, A] {
	return iter.MonadFilter(as, either.ForAll[E](pred))
}

// Filter returns a function that filters SeqEither elements based on a predicate.
// This is the curried version of MonadFilter, useful for creating reusable filter operations.
//
// The returned function keeps only Right values that satisfy the predicate, while
// passing all Left values through unchanged. Right values that fail the predicate
// are removed from the sequence.
//
// Type Parameters:
//   - E: The type of the Left value (error type)
//   - A: The type of the Right value (success type)
//
// Parameters:
//   - pred: A predicate function that tests Right values
//
// Returns:
//   - An Operator that filters SeqEither sequences based on the predicate
//
// Example:
//
//	// Create a reusable filter for even numbers
//	evens := itereither.Filter[string](func(x int) bool { return x%2 == 0 })
//
//	seq1 := iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3))
//	result1 := evens(seq1)
//	// yields: Right(2)
//
//	seq2 := iter.From(E.Right[string](4), E.Left[int]("error"), E.Right[string](5))
//	result2 := evens(seq2)
//	// yields: Right(4), Left("error")
//
// See Also:
//
// MonadFilter is the non-curried version.
// FilterOrElse converts failing values to Left instead of removing them.
//
//go:inline
func Filter[E, A any](pred func(A) bool) Operator[E, A, A] {
	return iter.Filter(either.ForAll[E](pred))
}
