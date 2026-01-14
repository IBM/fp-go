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

import F "github.com/IBM/fp-go/v2/function"

// MonadExtend applies a function to an Either value, where the function receives the entire Either as input.
// This is the Extend (or Comonad) operation that allows computations to depend on the context.
//
// If the Either is Left, it returns Left unchanged without applying the function.
// If the Either is Right, it applies the function to the entire Either and wraps the result in a Right.
//
// This operation is useful when you need to perform computations that depend on whether
// a value is present (Right) or absent (Left), not just on the value itself.
//
// Type Parameters:
//   - E: The error type (Left channel)
//   - A: The input value type (Right channel)
//   - B: The output value type
//
// Parameters:
//   - fa: The Either value to extend
//   - f: Function that takes the entire Either[E, A] and produces a value of type B
//
// Returns:
//   - Either[E, B]: Left if input was Left, otherwise Right containing the result of f(fa)
//
// Example:
//
//	// Count how many times we've seen a Right value
//	counter := func(e either.Either[error, int]) int {
//	    return either.Fold(
//	        func(err error) int { return 0 },
//	        func(n int) int { return 1 },
//	    )(e)
//	}
//	result := either.MonadExtend(either.Right[error](42), counter) // Right(1)
//	result := either.MonadExtend(either.Left[int](errors.New("err")), counter) // Left(error)
//
//go:inline
func MonadExtend[E, A, B any](fa Either[E, A], f func(Either[E, A]) B) Either[E, B] {
	if fa.isLeft {
		return Left[B](fa.l)
	}
	return Of[E](f(fa))
}

// Extend is the curried version of [MonadExtend].
// It returns a function that applies the given function to an Either value.
//
// This is useful for creating reusable transformations that depend on the Either context.
//
// Type Parameters:
//   - E: The error type (Left channel)
//   - A: The input value type (Right channel)
//   - B: The output value type
//
// Parameters:
//   - f: Function that takes the entire Either[E, A] and produces a value of type B
//
// Returns:
//   - Operator[E, A, B]: A function that transforms Either[E, A] to Either[E, B]
//
// Example:
//
//	// Create a reusable extender that extracts metadata
//	getMetadata := either.Extend(func(e either.Either[error, string]) string {
//	    return either.Fold(
//	        func(err error) string { return "error: " + err.Error() },
//	        func(s string) string { return "value: " + s },
//	    )(e)
//	})
//	result := getMetadata(either.Right[error]("hello")) // Right("value: hello")
//
//go:inline
func Extend[E, A, B any](f func(Either[E, A]) B) Operator[E, A, B] {
	return F.Bind2nd(MonadExtend[E, A, B], f)
}
