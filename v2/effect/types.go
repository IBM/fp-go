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

package effect

import (
	"github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Either represents a value that can be either a Left (error) or Right (success).
	Either[E, A any] = either.Either[E, A]

	// Reader represents a computation that depends on a context R and produces a value A.
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderIO represents a computation that depends on a context R and produces an IO action returning A.
	ReaderIO[R, A any] = readerio.ReaderIO[R, A]

	// IO represents a synchronous side effect that produces a value A.
	IO[A any] = io.IO[A]

	// IOEither represents a synchronous side effect that can fail with error E or succeed with value A.
	IOEither[E, A any] = ioeither.IOEither[E, A]

	// Lazy represents a lazily evaluated computation that produces a value A.
	Lazy[A any] = lazy.Lazy[A]

	// IOResult represents a synchronous side effect that can fail with an error or succeed with value A.
	IOResult[A any] = ioresult.IOResult[A]

	// ReaderIOResult represents a computation that depends on context and performs IO with error handling.
	ReaderIOResult[A any] = readerioresult.ReaderIOResult[A]

	// Monoid represents an algebraic structure with an associative binary operation and an identity element.
	Monoid[A any] = monoid.Monoid[A]

	// Effect represents an effectful computation that:
	//   - Requires a context of type C
	//   - Can perform I/O operations
	//   - Can fail with an error
	//   - Produces a value of type A on success
	//
	// This is the core type of the effect package, providing a complete effect system
	// for managing dependencies, errors, and side effects in a composable way.
	Effect[C, A any] = readerreaderioresult.ReaderReaderIOResult[C, A]

	// Thunk represents a computation that performs IO with error handling but doesn't require context.
	// It's equivalent to ReaderIOResult and is used as an intermediate step when providing context to an Effect.
	Thunk[A any] = ReaderIOResult[A]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	Predicate[A any] = predicate.Predicate[A]

	// Result represents a computation result that can be either an error (Left) or a success value (Right).
	Result[A any] = result.Result[A]

	// Lens represents an optic for focusing on a field T within a structure S.
	Lens[S, T any] = lens.Lens[S, T]

	// Kleisli represents a function from A to Effect[C, B], enabling monadic composition.
	// It's the fundamental building block for chaining effectful computations.
	Kleisli[C, A, B any] = readerreaderioresult.Kleisli[C, A, B]

	// Operator represents a function that transforms Effect[C, A] to Effect[C, B].
	// It's used for lifting operations over effects.
	Operator[C, A, B any] = readerreaderioresult.Operator[C, A, B]
)
