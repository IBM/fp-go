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

package state

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/iso"
	"github.com/IBM/fp-go/v2/pair"
)

// IMap is a profunctor-like operation for State that transforms both the state and value using an isomorphism.
// It applies an isomorphism f to the state (contravariantly via Get and covariantly via ReverseGet)
// and a function g to the value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Convert the input state from S2 to S1 using the isomorphism's Get function
//   - Run the State computation with S1
//   - Convert the output state back from S1 to S2 using the isomorphism's ReverseGet function
//   - Transform the result value from A to B using g
//
// Type Parameters:
//   - A: The original value type produced by the State
//   - S2: The new state type
//   - S1: The original state type expected by the State
//   - B: The new output value type
//
// Parameters:
//   - f: An isomorphism between S2 and S1
//   - g: Function to transform the output value from A to B
//
// Returns:
//   - A Kleisli arrow that takes a State[S1, A] and returns a State[S2, B]
//
//go:inline
func IMap[A, S2, S1, B any](f iso.Iso[S2, S1], g func(A) B) Kleisli[S2, State[S1, A], B] {
	return F.Bind13of3(F.Flow3[func(s S2) S1, State[S1, A], func(pair.Pair[S1, A]) pair.Pair[S2, B]])(f.Get, pair.BiMap(f.ReverseGet, g))
}

// MapState is a contravariant-like operation for State that transforms the state type using an isomorphism.
// It applies an isomorphism f to convert between state types while preserving the value type.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Convert the input state from S2 to S1 using the isomorphism's Get function
//   - Run the State computation with S1
//   - Convert the output state back from S1 to S2 using the isomorphism's ReverseGet function
//   - Keep the value type A unchanged
//
// Type Parameters:
//   - A: The value type (unchanged)
//   - S2: The new state type
//   - S1: The original state type expected by the State
//
// Parameters:
//   - f: An isomorphism between S2 and S1
//
// Returns:
//   - A Kleisli arrow that takes a State[S1, A] and returns a State[S2, A]
//
//go:inline
func MapState[A, S2, S1 any](f iso.Iso[S2, S1]) Kleisli[S2, State[S1, A], A] {
	return F.Bind13of3(F.Flow3[func(S2) S1, State[S1, A], func(pair.Pair[S1, A]) pair.Pair[S2, A]])(f.Get, pair.MapHead[A](f.ReverseGet))
}
