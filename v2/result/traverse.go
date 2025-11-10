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

package result

import (
	"github.com/IBM/fp-go/v2/either"
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
//	// Traverse an Result[Option[int]] to Option[Result[int]]
//	result := either.Traverse[int, error, int, option.Option[int], option.Option[either.Result[int]]](
//	    option.Of[either.Result[int]],
//	    option.Map[int, either.Result[int]],
//	)(f)(eitherOfOption)
//
//go:inline
func Traverse[A, B, HKTB, HKTRB any](
	mof func(Result[B]) HKTRB,
	mmap func(Kleisli[B, B]) func(HKTB) HKTRB,
) func(func(A) HKTB) func(Result[A]) HKTRB {
	return either.Traverse[A](mof, mmap)
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
//	// Sequence an Result[Option[int]] to Option[Result[int]]
//	result := either.Sequence[error, int, option.Option[int], option.Option[either.Result[int]]](
//	    option.Of[either.Result[int]],
//	    option.Map[int, either.Result[int]],
//	)(eitherOfOption)
//
//go:inline
func Sequence[A, HKTA, HKTRA any](
	mof func(Result[A]) HKTRA,
	mmap func(Kleisli[A, A]) func(HKTA) HKTRA,
) func(Result[HKTA]) HKTRA {
	return either.Sequence(mof, mmap)
}
