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

package iterresult

import (
	"github.com/IBM/fp-go/v2/iterator/itereither"
)

// FilterOrElse filters a SeqResult value based on a predicate.
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
//   - An Operator that filters SeqResult values based on the predicate
//
// Example:
//
//	// Validate that a number is positive
//	isPositive := func(x int) bool { return x > 0 }
//	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }
//
//	validatePositive := FilterOrElse(isPositive, onNegative)
//
//	result1 := validatePositive(Right(42))  // Right(42)
//	result2 := validatePositive(Right(-5))  // Left(error: "-5 is not positive")
//	result3 := validatePositive(Left[int](errors.New("error"))) // Left(error: "error")
//
//go:inline
func FilterOrElse[A any](pred Predicate[A], onFalse func(A) error) Operator[A, A] {
	return itereither.FilterOrElse(pred, onFalse)
}

// MonadFilter filters a SeqResult sequence, keeping only Ok values that satisfy the predicate.
// Error values are always passed through unchanged, regardless of the predicate.
//
// This function processes each element in the sequence:
//   - If the element is an error, it passes through unchanged
//   - If the element is Ok and satisfies the predicate, it passes through
//   - If the element is Ok but fails the predicate, it is filtered out
//
// Unlike FilterOrElse, which converts failing Ok values to errors, MonadFilter
// simply removes them from the sequence entirely.
//
// Marble diagram:
//
//	Input:  --Ok(1)--Ok(2)--Err(e)--Ok(3)--Ok(4)--Ok(5)-->
//	pred(x) = x % 2 == 0
//	Output: ---------Ok(2)--Err(e)---------Ok(4)--------->
//
// Where Ok(x) represents a successful result and Err(e) represents an error.
// Odd numbers are filtered out, even numbers and errors pass through.
//
// Type Parameters:
//   - A: The type of the Ok value (success type)
//
// Parameters:
//   - as: The input SeqResult sequence to filter
//   - pred: A predicate function that tests Ok values
//
// Returns:
//   - A SeqResult containing only error values and Ok values that satisfy the predicate
//
// Example:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/result"
//	    "github.com/IBM/fp-go/v2/iterator/iter"
//	)
//
//	seq := iter.From(
//	    R.Of(1),
//	    R.Of(2),
//	    R.Error[int](errors.New("error")),
//	    R.Of(3),
//	    R.Of(4),
//	)
//	isEven := func(x int) bool { return x%2 == 0 }
//	result := iterresult.MonadFilter(seq, isEven)
//	// yields: Ok(2), Err(error), Ok(4)
//
// See Also:
//
// Filter is the curried version of MonadFilter.
// FilterOrElse converts failing values to errors instead of removing them.
func MonadFilter[A any](as SeqResult[A], pred Predicate[A]) SeqResult[A] {
	return itereither.MonadFilter(as, pred)
}

// Filter returns a function that filters SeqResult elements based on a predicate.
// This is the curried version of MonadFilter, useful for creating reusable filter operations.
//
// The returned function keeps only Ok values that satisfy the predicate, while
// passing all error values through unchanged. Ok values that fail the predicate
// are removed from the sequence.
//
// Type Parameters:
//   - A: The type of the Ok value (success type)
//
// Parameters:
//   - pred: A predicate function that tests Ok values
//
// Returns:
//   - An Operator that filters SeqResult sequences based on the predicate
//
// Example:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/result"
//	    "github.com/IBM/fp-go/v2/iterator/iter"
//	)
//
//	// Create a reusable filter for even numbers
//	evens := iterresult.Filter(func(x int) bool { return x%2 == 0 })
//
//	seq1 := iter.From(R.Of(1), R.Of(2), R.Of(3))
//	result1 := evens(seq1)
//	// yields: Ok(2)
//
//	seq2 := iter.From(R.Of(4), R.Error[int](errors.New("error")), R.Of(5))
//	result2 := evens(seq2)
//	// yields: Ok(4), Err(error)
//
// See Also:
//
// MonadFilter is the non-curried version.
// FilterOrElse converts failing values to errors instead of removing them.
//
//go:inline
func Filter[A any](pred func(A) bool) Operator[A, A] {
	return itereither.Filter[error](pred)
}
