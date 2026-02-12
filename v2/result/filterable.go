// Copyright (c) 2025 IBM Corp.
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

// Package result provides filterable operations for Result types.
//
// This package implements the Fantasy Land Filterable specification:
// https://github.com/fantasyland/fantasy-land#filterable
//
// Since Result[A] is an alias for Either[error, A], these functions are
// thin wrappers around the corresponding either package functions, specialized
// for the common case where the error type is Go's built-in error interface.
package result

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/option"
)

// Partition separates a [Result] value into a [Pair] based on a predicate function.
// It returns a function that takes a Result and produces a Pair of Result values,
// where the first element contains values that fail the predicate and the second
// contains values that pass the predicate.
//
// This function implements the Filterable specification's partition operation:
// https://github.com/fantasyland/fantasy-land#filterable
//
// The behavior is as follows:
//   - If the input is an error (Left), both elements of the resulting Pair will be the same error
//   - If the input is Ok (Right) and the predicate returns true, the result is (Err(empty), Ok(value))
//   - If the input is Ok (Right) and the predicate returns false, the result is (Ok(value), Err(empty))
//
// Parameters:
//   - p: A predicate function that tests values of type A
//   - empty: The default error to use when creating error Results for partitioning
//
// Returns:
//
//	A function that takes a Result[A] and returns a Pair where:
//	  - First element: Result values that fail the predicate (or original error)
//	  - Second element: Result values that pass the predicate (or original error)
//
// Example:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/result"
//	    N "github.com/IBM/fp-go/v2/number"
//	    P "github.com/IBM/fp-go/v2/pair"
//	    "errors"
//	)
//
//	// Partition positive and non-positive numbers
//	isPositive := N.MoreThan(0)
//	partition := R.Partition(isPositive, errors.New("not positive"))
//
//	// Ok value that passes predicate
//	result1 := partition(R.Of(5))
//	// result1 = Pair(Err("not positive"), Ok(5))
//
//	// Ok value that fails predicate
//	result2 := partition(R.Of(-3))
//	// result2 = Pair(Ok(-3), Err("not positive"))
//
//	// Error passes through unchanged in both positions
//	result3 := partition(R.Error[int](errors.New("original error")))
//	// result3 = Pair(Err("original error"), Err("original error"))
//
//go:inline
func Partition[A any](p Predicate[A], empty error) func(Result[A]) Pair[Result[A], Result[A]] {
	return either.Partition(p, empty)
}

// Filter creates a filtering operation for [Result] values based on a predicate function.
// It returns a function that takes a Result and produces a Result, where Ok values
// that fail the predicate are converted to error Results with the provided error.
//
// This function implements the Filterable specification's filter operation:
// https://github.com/fantasyland/fantasy-land#filterable
//
// The behavior is as follows:
//   - If the input is an error, it passes through unchanged
//   - If the input is Ok and the predicate returns true, the Ok value passes through unchanged
//   - If the input is Ok and the predicate returns false, it's converted to Err(empty)
//
// Parameters:
//   - p: A predicate function that tests values of type A
//   - empty: The default error to use when filtering out Ok values that fail the predicate
//
// Returns:
//
//	An Operator function that takes a Result[A] and returns a Result[A] where:
//	  - Error values pass through unchanged
//	  - Ok values that pass the predicate remain as Ok
//	  - Ok values that fail the predicate become Err(empty)
//
// Example:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/result"
//	    N "github.com/IBM/fp-go/v2/number"
//	    "errors"
//	)
//
//	// Filter to keep only positive numbers
//	isPositive := N.MoreThan(0)
//	filterPositive := R.Filter(isPositive, errors.New("not positive"))
//
//	// Ok value that passes predicate - remains Ok
//	result1 := filterPositive(R.Of(5))
//	// result1 = Ok(5)
//
//	// Ok value that fails predicate - becomes Err
//	result2 := filterPositive(R.Of(-3))
//	// result2 = Err("not positive")
//
//	// Error passes through unchanged
//	result3 := filterPositive(R.Error[int](errors.New("original error")))
//	// result3 = Err("original error")
//
//	// Chaining filters
//	isEven := func(n int) bool { return n%2 == 0 }
//	filterEven := R.Filter(isEven, errors.New("not even"))
//
//	result4 := filterEven(filterPositive(R.Of(4)))
//	// result4 = Ok(4) - passes both filters
//
//	result5 := filterEven(filterPositive(R.Of(3)))
//	// result5 = Err("not even") - passes first, fails second
//
//go:inline
func Filter[A any](p Predicate[A], empty error) Operator[A, A] {
	return either.Filter(p, empty)
}

// FilterMap combines filtering and mapping operations for [Result] values using an [Option]-returning function.
// It returns a function that takes a Result[A] and produces a Result[B], where Ok values
// are transformed by applying the function f. If f returns Some(B), the result is Ok(B). If f returns
// None, the result is Err(empty).
//
// This function implements the Filterable specification's filterMap operation:
// https://github.com/fantasyland/fantasy-land#filterable
//
// The behavior is as follows:
//   - If the input is an error, it passes through with its error value preserved
//   - If the input is Ok and f returns Some(B), the result is Ok(B)
//   - If the input is Ok and f returns None, the result is Err(empty)
//
// Parameters:
//   - f: An Option Kleisli function that transforms values of type A to Option[B]
//   - empty: The default error to use when f returns None
//
// Returns:
//
//	An Operator function that takes a Result[A] and returns a Result[B] where:
//	  - Error values pass through with error preserved
//	  - Ok values are transformed by f: Some(B) becomes Ok(B), None becomes Err(empty)
//
// Example:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/result"
//	    O "github.com/IBM/fp-go/v2/option"
//	    "errors"
//	    "strconv"
//	)
//
//	// Parse string to int, filtering out invalid values
//	parseInt := func(s string) O.Option[int] {
//	    if n, err := strconv.Atoi(s); err == nil {
//	        return O.Some(n)
//	    }
//	    return O.None[int]()
//	}
//	filterMapInt := R.FilterMap(parseInt, errors.New("invalid number"))
//
//	// Valid number string - transforms to Ok(int)
//	result1 := filterMapInt(R.Of("42"))
//	// result1 = Ok(42)
//
//	// Invalid number string - becomes Err
//	result2 := filterMapInt(R.Of("abc"))
//	// result2 = Err("invalid number")
//
//	// Error passes through with error preserved
//	result3 := filterMapInt(R.Error[string](errors.New("original error")))
//	// result3 = Err("original error")
//
//go:inline
func FilterMap[A, B any](f option.Kleisli[A, B], empty error) Operator[A, B] {
	return either.FilterMap(f, empty)
}

// PartitionMap separates and transforms a [Result] value into a [Pair] of Result values using a mapping function.
// It returns a function that takes a Result[A] and produces a Pair of Result values, where the mapping
// function f transforms the Ok value into Either[B, C]. The result is partitioned based on whether f
// produces a Left or Right value.
//
// This function implements the Filterable specification's partitionMap operation:
// https://github.com/fantasyland/fantasy-land#filterable
//
// The behavior is as follows:
//   - If the input is an error, both elements of the resulting Pair will be errors with the original error
//   - If the input is Ok and f returns Left(B), the result is (Ok(B), Err(empty))
//   - If the input is Ok and f returns Right(C), the result is (Err(empty), Ok(C))
//
// Parameters:
//   - f: A Kleisli function that transforms values of type A to Either[B, C]
//   - empty: The default error to use when creating error Results for partitioning
//
// Returns:
//
//	A function that takes a Result[A] and returns a Pair[Result[B], Result[C]] where:
//	  - If input is error: (Err(original_error), Err(original_error))
//	  - If f returns Left(B): (Ok(B), Err(empty))
//	  - If f returns Right(C): (Err(empty), Ok(C))
//
// Example:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/result"
//	    E "github.com/IBM/fp-go/v2/either"
//	    P "github.com/IBM/fp-go/v2/pair"
//	    "errors"
//	    "strconv"
//	)
//
//	// Classify and transform numbers: negative -> error message, positive -> squared value
//	classifyNumber := func(n int) E.Either[string, int] {
//	    if n < 0 {
//	        return E.Left[int]("negative: " + strconv.Itoa(n))
//	    }
//	    return E.Right[string](n * n)
//	}
//	partitionMap := R.PartitionMap(classifyNumber, errors.New("not classified"))
//
//	// Positive number - goes to right side as squared value
//	result1 := partitionMap(R.Of(5))
//	// result1 = Pair(Err("not classified"), Ok(25))
//
//	// Negative number - goes to left side with error message
//	result2 := partitionMap(R.Of(-3))
//	// result2 = Pair(Ok("negative: -3"), Err("not classified"))
//
//	// Original error - appears in both positions
//	result3 := partitionMap(R.Error[int](errors.New("original error")))
//	// result3 = Pair(Err("original error"), Err("original error"))
//
//go:inline
func PartitionMap[A, B, C any](f either.Kleisli[B, A, C], empty error) func(Result[A]) Pair[Result[B], Result[C]] {
	return either.PartitionMap(f, empty)
}
