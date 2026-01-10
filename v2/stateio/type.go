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

package stateio

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/optics/iso/lens"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/state"
)

type (
	// Endomorphism represents a function from A to A.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Lens is an optic that focuses on a field of type A within a structure of type S.
	Lens[S, A any] = lens.Lens[S, A]

	// State represents a stateful computation that takes an initial state S and returns
	// a pair of the new state S and a value A.
	State[S, A any] = state.State[S, A]

	// Pair represents a tuple of two values.
	Pair[L, R any] = pair.Pair[L, R]

	// Reader represents a computation that depends on an environment/context of type R
	// and produces a value of type A.
	Reader[R, A any] = reader.Reader[R, A]

	// Either represents a value that can be either a Left (error) or Right (success).
	Either[E, A any] = either.Either[E, A]

	// IO represents a computation that performs side effects and produces a value of type A.
	IO[A any] = io.IO[A]

	StateIO[S, A any] = Reader[S, IO[Pair[S, A]]]

	Kleisli[S, A, B any] = Reader[A, StateIO[S, B]]

	Operator[S, A, B any] = Reader[StateIO[S, A], StateIO[S, B]]

	Predicate[A any] = predicate.Predicate[A]
)
