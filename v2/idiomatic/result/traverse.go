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
func Traverse[A, B, HKTB, HKTRB any](
	mof func(B, error) HKTRB,
	mmap func(Kleisli[B, B]) func(HKTB) HKTRB,
) func(func(A) HKTB) func(A, error) HKTRB {
	return func(f func(A) HKTB) func(A, error) HKTRB {
		right := mmap(Right[B])
		return func(a A, err error) HKTRB {
			if err != nil {
				return mof(Left[B](err))
			}
			return right(f(a))
		}
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
func Sequence[A, HKTA, HKTRA any](
	mof func(A, error) HKTRA,
	mmap func(Kleisli[A, A]) func(HKTA) HKTRA,
) func(hkta HKTA, err error) HKTRA {
	right := mmap(Right[A])
	return func(hkta HKTA, err error) HKTRA {
		if err != nil {
			return mof(Left[A](err))
		}
		return right(hkta)
	}
}
