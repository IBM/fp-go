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

// Package readerresult implements a specialization of the Reader monad assuming a golang context as the context of the monad and a standard golang error.
//
// # Pure vs Effectful Functions
//
// This package distinguishes between pure (side-effect free) and effectful (side-effectful) functions:
//
// EFFECTFUL FUNCTIONS (depend on context.Context):
//   - ReaderResult[A]: func(context.Context) (A, error) - Effectful computation that needs context
//   - These functions are effectful because context.Context is effectful (can be cancelled, has deadlines, carries values)
//   - Use for: operations that need cancellation, timeouts, context values, or any context-dependent behavior
//   - Examples: database queries, HTTP requests, operations that respect cancellation
//
// PURE FUNCTIONS (side-effect free):
//   - func(State) (Value, error) - Pure computation that only depends on state, not context
//   - func(State) Value - Pure transformation without errors
//   - These functions are pure because they only read from their input state and don't depend on external context
//   - Use for: parsing, validation, calculations, data transformations that don't need context
//   - Examples: JSON parsing, input validation, mathematical computations
//
// The package provides different bind operations for each:
//   - Bind: For effectful ReaderResult computations (State -> ReaderResult[Value])
//   - BindResultK: For pure functions with errors (State -> (Value, error))
//   - Let: For pure functions without errors (State -> Value)
//   - BindReaderK: For context-dependent pure functions (State -> Reader[Context, Value])
//   - BindEitherK: For pure Result/Either values (State -> Result[Value])
package readerresult

import (
	"context"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Option[A any]    = option.Option[A]
	Either[A any]    = either.Either[error, A]
	Result[A any]    = result.Result[A]
	Reader[R, A any] = reader.Reader[R, A]
	// ReaderResult is a specialization of the Reader monad for the typical golang scenario
	ReaderResult[A any] = readereither.ReaderEither[context.Context, error, A]

	Kleisli[A, B any]   = reader.Reader[A, ReaderResult[B]]
	Operator[A, B any]  = Kleisli[ReaderResult[A], B]
	Endomorphism[A any] = endomorphism.Endomorphism[A]
	Prism[S, T any]     = prism.Prism[S, T]
	Lens[S, T any]      = lens.Lens[S, T]
)
