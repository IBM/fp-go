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
//
// deprecated: use MakeTraversable instead
func Traverse[A, B, HKTB, HKTOB any](
	mof func(Option[B]) HKTOB,
	mmap func(Kleisli[B, B]) func(HKTB) HKTOB,
) func(func(A) HKTB) func(Option[A]) HKTOB {
	return MakeTraversable[A](mof, mmap)
}

// MakeTraversable creates a fully curried traversal function for Option types.
// This function enables traversing an Option by applying a transformation that produces
// a higher-kinded type, then sequencing the result into that higher-kinded type containing an Option.
//
// This is the preferred way to create traversable operations for Option types, replacing
// the deprecated Traverse function. It provides maximum composability by returning a function
// that can be specialized with different transformation functions.
//
// Type Parameters:
//   - A: The input element type contained in the Option
//   - B: The output element type after transformation
//   - HKTB: The higher-kinded type containing B (e.g., IO[B], Either[E, B])
//   - HKTOB: The higher-kinded type containing Option[B] (e.g., IO[Option[B]])
//
// Parameters:
//   - mof: Function to lift an Option into the target higher-kinded type
//   - mmap: Function to map over the target higher-kinded type
//
// Returns:
//   - A function that takes a transformation function and returns a function that takes
//     an Option and produces the higher-kinded type containing an Option
//
// Behavior:
//   - None values are lifted directly into the target type as None
//   - Some values are transformed by f, then wrapped in Some and lifted into the target type
//
// Example:
//
//	import (
//	    O "github.com/IBM/fp-go/v2/option"
//	    IO "github.com/IBM/fp-go/v2/io"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Create a traversable for Option that works with IO
//	traverseIO := O.MakeTraversable[int, string](
//	    IO.Of[O.Option[string]],
//	    IO.Map[string, O.Option[string]],
//	)
//
//	// Use it to transform Option[int] to IO[Option[string]]
//	fetchUser := func(id int) IO.IO[string] {
//	    return IO.Of(F.Pipe1(id, strconv.Itoa))
//	}
//
//	result := traverseIO(fetchUser)(O.Some(42))  // IO[Some("42")]
//	result2 := traverseIO(fetchUser)(O.None[int]())  // IO[None]
//
// See Also:
//   - Traverse: Deprecated alias for this function
//   - Sequence: Traversal with identity transformation
func MakeTraversable[A, B, HKTB, HKTOB any](
	mof func(Option[B]) HKTOB,
	mmap func(Kleisli[B, B]) func(HKTB) HKTOB,
) func(func(A) HKTB) func(Option[A]) HKTOB {
	onNone := F.Nullary2(None[B], mof)
	onSome := mmap(Some[B])
	return func(f func(A) HKTB) func(Option[A]) HKTOB {
		return Fold(onNone, F.Flow2(f, onSome))
	}
}
