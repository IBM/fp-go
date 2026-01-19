// Copyright (c) 2024 IBM Corp.
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

package readerreaderioresult

import (
	"context"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/traversal/result"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/readerreaderioeither"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// Option represents an optional value that may or may not be present.
	// It's an alias for option.Option[A].
	Option[A any] = option.Option[A]

	// Lazy represents a lazily evaluated computation that produces a value of type A.
	// It's an alias for lazy.Lazy[A].
	Lazy[A any] = lazy.Lazy[A]

	// Reader represents a computation that depends on an environment of type R
	// and produces a value of type A.
	// It's an alias for reader.Reader[R, A].
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderOption represents a computation that depends on an environment of type R
	// and produces an optional value of type A.
	// It's an alias for readeroption.ReaderOption[R, A].
	ReaderOption[R, A any] = readeroption.ReaderOption[R, A]

	// ReaderIO represents a computation that depends on an environment of type R
	// and performs side effects to produce a value of type A.
	// It's an alias for readerio.ReaderIO[R, A].
	ReaderIO[R, A any] = readerio.ReaderIO[R, A]

	// ReaderIOResult represents a computation that depends on an environment of type R,
	// performs side effects, and may fail with an error.
	// It's an alias for readerioresult.ReaderIOResult[R, A].
	ReaderIOResult[R, A any] = readerioresult.ReaderIOResult[R, A]

	// Either represents a value that can be one of two types: Left (error) or Right (success).
	// It's an alias for either.Either[E, A].
	Either[E, A any] = either.Either[E, A]

	// Result is a specialized Either with error as the left type.
	// It's an alias for result.Result[A] which is Either[error, A].
	Result[A any] = result.Result[A]

	// IOEither represents a side-effecting computation that may fail with an error of type E
	// or succeed with a value of type A.
	// It's an alias for ioeither.IOEither[E, A].
	IOEither[E, A any] = ioeither.IOEither[E, A]

	// IOResult represents a side-effecting computation that may fail with an error
	// or succeed with a value of type A.
	// It's an alias for ioresult.IOResult[A] which is IOEither[error, A].
	IOResult[A any] = ioresult.IOResult[A]

	// IO represents a side-effecting computation that produces a value of type A.
	// It's an alias for io.IO[A].
	IO[A any] = io.IO[A]

	// ReaderReaderIOEither is the base monad transformer that combines:
	// - Reader[R, ...] for outer dependency injection
	// - Reader[C, ...] for inner dependency injection (typically context.Context)
	// - IO for side effects
	// - Either[E, A] for error handling
	// It's an alias for readerreaderioeither.ReaderReaderIOEither[R, C, E, A].
	ReaderReaderIOEither[R, C, E, A any] = readerreaderioeither.ReaderReaderIOEither[R, C, E, A]

	// ReaderReaderIOResult is the main type of this package, specializing ReaderReaderIOEither
	// with context.Context as the inner reader type and error as the error type.
	//
	// Type structure:
	//   ReaderReaderIOResult[R, A] = R -> context.Context -> IO[Either[error, A]]
	//
	// This represents a computation that:
	// 1. Depends on an outer environment of type R (e.g., application config)
	// 2. Depends on a context.Context for cancellation and request-scoped values
	// 3. Performs side effects (IO)
	// 4. May fail with an error or succeed with a value of type A
	//
	// This is the primary type used throughout the package for composing
	// context-aware, effectful computations with error handling.
	ReaderReaderIOResult[R, A any] = ReaderReaderIOEither[R, context.Context, error, A]

	// Kleisli represents a function from A to a monadic value ReaderReaderIOResult[R, B].
	// It's used for composing monadic functions using Kleisli composition.
	//
	// Type structure:
	//   Kleisli[R, A, B] = A -> ReaderReaderIOResult[R, B]
	//
	// Kleisli arrows can be composed using Chain operations to build complex
	// data transformation pipelines.
	Kleisli[R, A, B any] = Reader[A, ReaderReaderIOResult[R, B]]

	// Operator is a specialized Kleisli arrow that operates on monadic values.
	// It takes a ReaderReaderIOResult[R, A] and produces a ReaderReaderIOResult[R, B].
	//
	// Type structure:
	//   Operator[R, A, B] = ReaderReaderIOResult[R, A] -> ReaderReaderIOResult[R, B]
	//
	// Operators are useful for transforming monadic computations, such as
	// adding retry logic, logging, or error recovery.
	Operator[R, A, B any] = Kleisli[R, ReaderReaderIOResult[R, A], B]

	// Lens represents an optic for focusing on a part of a data structure.
	// It provides a way to get and set a field T within a structure S.
	// It's an alias for lens.Lens[S, T].
	Lens[S, T any] = lens.Lens[S, T]

	// Trampoline is used for stack-safe recursion through tail call optimization.
	// It's an alias for tailrec.Trampoline[L, B].
	Trampoline[L, B any] = tailrec.Trampoline[L, B]

	// Predicate represents a function that tests whether a value of type A
	// satisfies some condition.
	// It's an alias for predicate.Predicate[A].
	Predicate[A any] = predicate.Predicate[A]

	// Endmorphism represents a function from type A to type A.
	// It's an alias for endomorphism.Endomorphism[A].
	Endmorphism[A any] = endomorphism.Endomorphism[A]

	Void = function.Void
)
