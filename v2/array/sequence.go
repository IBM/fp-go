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

package array

import (
	"github.com/IBM/fp-go/v2/internal/array"
	M "github.com/IBM/fp-go/v2/monoid"
	O "github.com/IBM/fp-go/v2/option"
)

func MonadSequence[HKTA, HKTRA any](
	fof func(HKTA) HKTRA,
	m M.Monoid[HKTRA],
	ma []HKTA) HKTRA {
	return array.MonadSequence(fof, m.Empty(), m.Concat, ma)
}

// Sequence takes an array where elements are HKT<A> (higher kinded type) and,
// using an applicative of that HKT, returns an HKT of []A.
//
// For example, it can turn:
//   - []Either[error, string] into Either[error, []string]
//   - []Option[int] into Option[[]int]
//
// Sequence requires an Applicative of the HKT you are targeting. To turn an
// []Either[E, A] into an Either[E, []A], it needs an Applicative for Either.
// To turn an []Option[A] into an Option[[]A], it needs an Applicative for Option.
//
// Note: We need to pass the members of the applicative explicitly because Go does not
// support higher kinded types or template methods on structs or interfaces.
//
// Type parameters:
//   - HKTA = HKT<A> (e.g., Option[A], Either[E, A])
//   - HKTRA = HKT<[]A> (e.g., Option[[]A], Either[E, []A])
//   - HKTFRA = HKT<func(A)[]A> (e.g., Option[func(A)[]A])
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/option"
//
//	opts := []option.Option[int]{
//	    option.Some(1),
//	    option.Some(2),
//	    option.Some(3),
//	}
//
//	seq := array.Sequence(
//	    option.Of[[]int],
//	    option.MonadMap[[]int, func(int) []int],
//	    option.MonadAp[[]int, int],
//	)
//	result := seq(opts) // Some([1, 2, 3])
func Sequence[HKTA, HKTRA any](
	fof func(HKTA) HKTRA,
	m M.Monoid[HKTRA],
) func([]HKTA) HKTRA {
	return array.Sequence[[]HKTA](fof, m.Empty(), m.Concat)
}

// ArrayOption returns a function to convert a sequence of options into an option of a sequence.
// If all options are Some, returns Some containing an array of all values.
// If any option is None, returns None.
//
// Example:
//
//	opts := []option.Option[int]{
//	    option.Some(1),
//	    option.Some(2),
//	    option.Some(3),
//	}
//	result := array.ArrayOption[int]()(opts) // Some([1, 2, 3])
//
//	opts2 := []option.Option[int]{
//	    option.Some(1),
//	    option.None[int](),
//	    option.Some(3),
//	}
//	result2 := array.ArrayOption[int]()(opts2) // None
func ArrayOption[A any](ma []Option[A]) Option[[]A] {
	return MonadSequence(
		O.Map(Of[A]),
		O.ApplicativeMonoid(Monoid[A]()),
		ma,
	)
}
