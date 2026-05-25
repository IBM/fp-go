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

// Traverse converts an Either of some higher kinded type into the higher kinded type of an Either.
// This is a generic traversal operation that works with any applicative functor.
//
// Parameters:
//   - mof: Lifts an Either into the target higher-kinded type
//   - mmap: Maps over the target higher-kinded type
//
// Example (conceptual - requires understanding of higher-kinded types):
//
//	// Traverse an Either[error, Option[int]] to Option[Either[error, int]]
//	result := either.Traverse[int, error, int, option.Option[int], option.Option[either.Either[error, int]]](
//	    option.Of[either.Either[error, int]],
//	    option.Map[int, either.Either[error, int]],
//	)(f)(eitherOfOption)
//
// deprecated: use MakeTraversable instead
func Traverse[A, E, B, HKTB, HKTRB any](
	mof func(Either[E, B]) HKTRB,
	mmap func(Kleisli[E, B, B]) func(HKTB) HKTRB,
) func(func(A) HKTB) func(Either[E, A]) HKTRB {
	return MakeTraversable[A](mof, mmap)
}

// MakeTraversable creates a fully curried traversal function for Either types.
// This function enables traversing an Either by applying a transformation that produces
// a higher-kinded type, then sequencing the result into that higher-kinded type containing an Either.
//
// This is the preferred way to create traversable operations for Either types, replacing
// the deprecated Traverse function. It provides maximum composability by returning a function
// that can be specialized with different transformation functions.
//
// Type Parameters:
//   - A: The input element type contained in the Right case
//   - E: The error type contained in the Left case
//   - B: The output element type after transformation
//   - HKTB: The higher-kinded type containing B (e.g., IO[B], Option[B])
//   - HKTRB: The higher-kinded type containing Either[E, B] (e.g., IO[Either[E, B]])
//
// Parameters:
//   - mof: Function to lift an Either into the target higher-kinded type
//   - mmap: Function to map over the target higher-kinded type
//
// Returns:
//   - A function that takes a transformation function and returns a function that takes
//     an Either and produces the higher-kinded type containing an Either
//
// Behavior:
//   - Left values are lifted directly into the target type as Left (error is preserved)
//   - Right values are transformed by f, then wrapped in Right and lifted into the target type
//
// Example:
//
//	import (
//	    E "github.com/IBM/fp-go/v2/either"
//	    IO "github.com/IBM/fp-go/v2/io"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Create a traversable for Either that works with IO
//	traverseIO := E.MakeTraversable[int, error, string](
//	    IO.Of[E.Either[error, string]],
//	    IO.Map[string, E.Either[error, string]],
//	)
//
//	// Use it to transform Either[error, int] to IO[Either[error, string]]
//	fetchUser := func(id int) IO.IO[string] {
//	    return IO.Of(F.Pipe1(id, strconv.Itoa))
//	}
//
//	result := traverseIO(fetchUser)(E.Right[error](42))  // IO[Right("42")]
//	result2 := traverseIO(fetchUser)(E.Left[int](errors.New("fail")))  // IO[Left(error)]
//
// See Also:
//   - Traverse: Deprecated alias for this function
//   - Sequence: Traversal with identity transformation
func MakeTraversable[A, E, B, HKTB, HKTRB any](
	mof func(Either[E, B]) HKTRB,
	mmap func(Kleisli[E, B, B]) func(HKTB) HKTRB,
) func(func(A) HKTB) func(Either[E, A]) HKTRB {

	left := F.Flow2(Left[B, E], mof)
	right := mmap(Right[E, B])

	return func(f func(A) HKTB) func(Either[E, A]) HKTRB {
		return Fold(left, F.Flow2(f, right))
	}
}

// Sequence converts an Either of some higher kinded type into the higher kinded type of an Either.
// This is the identity version of Traverse - it doesn't transform the values, just swaps the type constructors.
//
// Parameters:
//   - mof: Lifts an Either into the target higher-kinded type
//   - mmap: Maps over the target higher-kinded type
//
// Example (conceptual - requires understanding of higher-kinded types):
//
//	// Sequence an Either[error, Option[int]] to Option[Either[error, int]]
//	result := either.Sequence[error, int, option.Option[int], option.Option[either.Either[error, int]]](
//	    option.Of[either.Either[error, int]],
//	    option.Map[int, either.Either[error, int]],
//	)(eitherOfOption)
func Sequence[E, A, HKTA, HKTRA any](
	mof func(Either[E, A]) HKTRA,
	mmap func(Kleisli[E, A, A]) func(HKTA) HKTRA,
) func(Either[E, HKTA]) HKTRA {
	return Fold(F.Flow2(Left[A, E], mof), mmap(Right[E, A]))
}
