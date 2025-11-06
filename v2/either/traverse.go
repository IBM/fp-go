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
func Traverse[A, E, B, HKTB, HKTRB any](
	mof func(Either[E, B]) HKTRB,
	mmap func(func(B) Either[E, B]) func(HKTB) HKTRB,
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
	mmap func(func(A) Either[E, A]) func(HKTA) HKTRA,
) func(Either[E, HKTA]) HKTRA {
	return Fold(F.Flow2(Left[A, E], mof), mmap(Right[E, A]))
}
