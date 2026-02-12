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
	"github.com/IBM/fp-go/v2/option"
)

// MonadChainOptionK chains a function that returns an Option into a sequence,
// filtering out None values and unwrapping Some values.
//
// This is useful for operations that may or may not produce a value for each element
// in the sequence. Only the successful (Some) results are included in the output sequence,
// while None values are filtered out.
//
// This is the monadic form that takes the sequence as the first parameter.
//
// RxJS Equivalent: [concatMap] combined with [filter] - https://rxjs.dev/api/operators/concatMap
//
// Type parameters:
//   - A: The element type of the input sequence
//   - B: The element type of the output sequence (wrapped in Option by the function)
//
// Parameters:
//   - as: The input sequence to transform
//   - f: A function that takes an element and returns an Option[B]
//
// Returns:
//
//	A new sequence containing only the unwrapped Some values
//
// Example:
//
//	import (
//	    "strconv"
//	    F "github.com/IBM/fp-go/v2/function"
//	    O "github.com/IBM/fp-go/v2/option"
//	    I "github.com/IBM/fp-go/v2/iterator/iter"
//	)
//
//	// Parse strings to integers, filtering out invalid ones
//	parseNum := func(s string) O.Option[int] {
//	    if n, err := strconv.Atoi(s); err == nil {
//	        return O.Some(n)
//	    }
//	    return O.None[int]()
//	}
//
//	seq := I.From("1", "invalid", "2", "3", "bad")
//	result := I.MonadChainOptionK(seq, parseNum)
//	// yields: 1, 2, 3 (invalid strings are filtered out)
func MonadChainOptionK[A, B any](as Seq[A], f option.Kleisli[A, B]) Seq[B] {
	return MonadFilterMap(as, f)
}

// ChainOptionK returns an operator that chains a function returning an Option into a sequence,
// filtering out None values and unwrapping Some values.
//
// This is the curried version of [MonadChainOptionK], useful for function composition
// and creating reusable transformations.
//
// RxJS Equivalent: [concatMap] combined with [filter] - https://rxjs.dev/api/operators/concatMap
//
// Type parameters:
//   - A: The element type of the input sequence
//   - B: The element type of the output sequence (wrapped in Option by the function)
//
// Parameters:
//   - f: A function that takes an element and returns an Option[B]
//
// Returns:
//
//	An Operator that transforms Seq[A] to Seq[B], filtering out None values
//
// Example:
//
//	import (
//	    "strconv"
//	    F "github.com/IBM/fp-go/v2/function"
//	    O "github.com/IBM/fp-go/v2/option"
//	    I "github.com/IBM/fp-go/v2/iterator/iter"
//	)
//
//	// Create a reusable parser operator
//	parsePositive := I.ChainOptionK(func(x int) O.Option[int] {
//	    if x > 0 {
//	        return O.Some(x)
//	    }
//	    return O.None[int]()
//	})
//
//	result := F.Pipe1(
//	    I.From(-1, 2, -3, 4, 5),
//	    parsePositive,
//	)
//	// yields: 2, 4, 5 (negative numbers are filtered out)
//
//go:inline
func ChainOptionK[A, B any](f option.Kleisli[A, B]) Operator[A, B] {
	return FilterMap(f)
}

// FlatMapOptionK is an alias for [ChainOptionK].
//
// This provides a more familiar name for developers coming from other functional
// programming languages or libraries where "flatMap" is the standard terminology
// for the monadic bind operation.
//
// Type parameters:
//   - A: The element type of the input sequence
//   - B: The element type of the output sequence (wrapped in Option by the function)
//
// Parameters:
//   - f: A function that takes an element and returns an Option[B]
//
// Returns:
//
//	An Operator that transforms Seq[A] to Seq[B], filtering out None values
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    O "github.com/IBM/fp-go/v2/option"
//	    I "github.com/IBM/fp-go/v2/iterator/iter"
//	)
//
//	// Validate and transform data
//	validateAge := I.FlatMapOptionK(func(age int) O.Option[string] {
//	    if age >= 18 && age <= 120 {
//	        return O.Some(fmt.Sprintf("Valid age: %d", age))
//	    }
//	    return O.None[string]()
//	})
//
//	result := F.Pipe1(
//	    I.From(15, 25, 150, 30),
//	    validateAge,
//	)
//	// yields: "Valid age: 25", "Valid age: 30"
//
//go:inline
func FlatMapOptionK[A, B any](f option.Kleisli[A, B]) Operator[A, B] {
	return ChainOptionK(f)
}
