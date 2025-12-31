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

package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/optics/iso/lens"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerioeither"
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

	// IOEither represents a computation that performs side effects and can fail with an error E
	// or succeed with a value A.
	IOEither[E, A any] = ioeither.IOEither[E, A]

	// ReaderIOEither represents a computation that depends on an environment R,
	// performs side effects, and can fail with an error E or succeed with a value A.
	ReaderIOEither[R, E, A any] = readerioeither.ReaderIOEither[R, E, A]

	// ReaderEither represents a computation that depends on an environment R and can fail
	// with an error E or succeed with a value A (without side effects).
	ReaderEither[R, E, A any] = readereither.ReaderEither[R, E, A]

	// StateReaderIOEither represents a stateful computation that:
	//   - Takes an initial state S
	//   - Depends on an environment/context R
	//   - Performs side effects (IO)
	//   - Can fail with an error E or succeed with a value A
	//   - Returns a pair of the new state S and the result
	//
	// This is the main type of this package, combining State, Reader, IO, and Either monads.
	StateReaderIOEither[S, R, E, A any] = Reader[S, ReaderIOEither[R, E, Pair[S, A]]]

	// Kleisli represents a Kleisli arrow - a function that takes a value A and returns
	// a StateReaderIOEither computation producing B.
	// This is used for monadic composition via Chain.
	Kleisli[S, R, E, A, B any] = Reader[A, StateReaderIOEither[S, R, E, B]]

	// Operator represents a function that transforms one StateReaderIOEither into another.
	// This is commonly used for building composable operations via Map, Chain, etc.
	Operator[S, R, E, A, B any] = Reader[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]]

	Predicate[A any] = predicate.Predicate[A]
)
