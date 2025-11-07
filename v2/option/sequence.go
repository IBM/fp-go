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

package option

import (
	F "github.com/IBM/fp-go/v2/function"
)

// Sequence converts an Option of some higher kinded type into the higher kinded type of an Option.
// This is a generic sequencing operation that works with any applicative functor.
//
// Parameters:
//   - mof: wraps an Option in the target higher kinded type
//   - mmap: maps a function over the higher kinded type
//
// Example (conceptual - typically used with other monadic types):
//
//	// Sequencing Option[IO[A]] to IO[Option[A]]
//	result := Sequence[int, IO[int], IO[Option[int]]](
//	    func(opt Option[int]) IO[Option[int]] { return IO.Of(opt) },
//	    func(f func(int) Option[int]) func(IO[int]) IO[Option[int]] { ... },
//	)
func Sequence[A, HKTA, HKTOA any](
	mof func(Option[A]) HKTOA,
	mmap func(Kleisli[A, A]) func(HKTA) HKTOA,
) func(Option[HKTA]) HKTOA {
	return Fold(F.Nullary2(None[A], mof), mmap(Some[A]))
}

// Traverse converts an Option by applying a function that produces a higher kinded type,
// then sequences the result. This combines mapping and sequencing in one operation.
//
// Parameters:
//   - mof: wraps an Option in the target higher kinded type
//   - mmap: maps a function over the higher kinded type
//
// Returns a function that takes a transformation function and an Option, producing
// the higher kinded type containing an Option.
//
// Example (conceptual - typically used with other monadic types):
//
//	// Traversing Option[A] with a function A -> IO[B] to get IO[Option[B]]
//	result := Traverse[int, string, IO[string], IO[Option[string]]](
//	    func(opt Option[string]) IO[Option[string]] { return IO.Of(opt) },
//	    func(f func(string) Option[string]) func(IO[string]) IO[Option[string]] { ... },
//	)
func Traverse[A, B, HKTB, HKTOB any](
	mof func(Option[B]) HKTOB,
	mmap func(Kleisli[B, B]) func(HKTB) HKTOB,
) func(func(A) HKTB) func(Option[A]) HKTOB {
	onNone := F.Nullary2(None[B], mof)
	onSome := mmap(Some[B])
	return func(f func(A) HKTB) func(Option[A]) HKTOB {
		return Fold(onNone, F.Flow2(f, onSome))
	}
}
