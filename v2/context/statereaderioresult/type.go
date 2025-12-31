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

package statereaderioresult

import (
	RIORES "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/optics/iso/lens"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
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

	// Result represents a value that can be either an error or a success value.
	// This is specialized to use [error] as the error type.
	Result[A any] = result.Result[A]

	// IO represents a computation that performs side effects and produces a value of type A.
	IO[A any] = io.IO[A]

	// IOResult represents a computation that performs side effects and can fail with an error
	// or succeed with a value A.
	IOResult[A any] = ioresult.IOResult[A]

	// ReaderIOResult represents a computation that depends on a context.Context,
	// performs side effects, and can fail with an error or succeed with a value A.
	ReaderIOResult[A any] = RIORES.ReaderIOResult[A]

	// StateReaderIOResult represents a stateful computation that:
	//   - Takes an initial state S
	//   - Depends on a [context.Context]
	//   - Performs side effects (IO)
	//   - Can fail with an [error] or succeed with a value A
	//   - Returns a pair of the new state S and the result
	//
	// This is the main type of this package, combining State, Reader, IO, and Result monads.
	// It is a specialization of StateReaderIOEither with:
	//   - Context type fixed to [context.Context]
	//   - Error type fixed to [error]
	StateReaderIOResult[S, A any] = Reader[S, ReaderIOResult[Pair[S, A]]]

	// Kleisli represents a Kleisli arrow - a function that takes a value A and returns
	// a StateReaderIOResult computation producing B.
	// This is used for monadic composition via Chain.
	Kleisli[S, A, B any] = Reader[A, StateReaderIOResult[S, B]]

	// Operator represents a function that transforms one StateReaderIOResult into another.
	// This is commonly used for building composable operations via Map, Chain, etc.
	Operator[S, A, B any] = Reader[StateReaderIOResult[S, A], StateReaderIOResult[S, B]]

	Predicate[A any] = predicate.Predicate[A]
)
