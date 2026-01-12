// Copyright (c) 2024 - 2025 IBM Corp.
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
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	// Endomorphism represents a function from a type to itself (A -> A).
	// It's an alias for endomorphism.Endomorphism[A] and is commonly used for
	// state transformations and updates.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Lens represents a functional reference to a part of a data structure.
	// It's an alias for lens.Lens[S, A] where S is the whole structure and A is the part.
	// Lenses provide composable getters and setters for immutable data structures.
	Lens[S, A any] = lens.Lens[S, A]

	// Reader represents a computation that depends on an environment of type R and produces a value of type A.
	// It's an alias for reader.Reader[R, A] and is used for dependency injection patterns.
	Reader[R, A any] = reader.Reader[R, A]

	// Pair represents a tuple of two values of types L and R.
	// It's an alias for pair.Pair[L, R] where L is the head (left) and R is the tail (right).
	Pair[L, R any] = pair.Pair[L, R]

	// State represents a stateful computation that takes an initial state of type S,
	// performs some operation, and returns both a new state and a value of type A.
	// It's defined as Reader[S, Pair[S, A]], meaning it's a function that:
	// 1. Takes an initial state S as input
	// 2. Returns a Pair where the head is the new state S and the tail is the computed value A
	// The new state is in the head position and the value in the tail position because
	// the pair monad operates on the tail, allowing monadic operations to transform
	// the computed value while threading the state through the computation.
	State[S, A any] = Reader[S, pair.Pair[S, A]]

	// Kleisli represents a Kleisli arrow for the State monad.
	// It's a function from A to State[S, B], which allows composition of
	// stateful computations. This is the fundamental building block for chaining
	// operations that both depend on and modify state.
	Kleisli[S, A, B any] = Reader[A, State[S, B]]

	// Operator is a specialized Kleisli arrow that operates on State values.
	// It transforms a State[S, A] into a State[S, B], making it useful for
	// building pipelines of stateful transformations while maintaining the state type S.
	Operator[S, A, B any] = Kleisli[S, State[S, A], B]

	Void = function.Void
)
