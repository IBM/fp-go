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

package readerioresult

import (
	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/readerresult"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Either represents a value of one of two possible types (a disjoint union).
	Either[E, A any] = either.Either[E, A]

	// Result represents a computation that may fail with an error.
	Result[A any] = result.Result[A]

	// Reader represents a computation that depends on some context/environment of type R
	// and produces a value of type A. It's useful for dependency injection patterns.
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderResult represents a computation that depends on an environment R and may fail with an error.
	ReaderResult[R, A any] = readerresult.ReaderResult[R, A]

	// ReaderIO represents a computation that depends on some context R and performs
	// side effects to produce a value of type A.
	ReaderIO[R, A any] = readerio.ReaderIO[R, A]

	// ReaderOption represents a computation that depends on an environment R and may not produce a value.
	ReaderOption[R, A any] = readeroption.ReaderOption[R, A]

	// IOEither represents a computation that performs side effects and can either
	// fail with an error of type E or succeed with a value of type A.
	IOEither[E, A any] = ioeither.IOEither[E, A]

	// IOResult represents a synchronous computation that may fail with an error.
	IOResult[A any] = ioresult.IOResult[A]

	// IO represents a synchronous computation that cannot fail.
	IO[A any] = io.IO[A]

	// Lazy represents a deferred computation that produces a value of type A.
	Lazy[A any] = lazy.Lazy[A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Endomorphism represents a function from a type to itself (A -> A).
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// ReaderIOResult represents a computation that:
	//   - Depends on some context/environment of type R (Reader)
	//   - Performs side effects (IO)
	//   - Can fail with an error of type E or succeed with a value of type A (Either)
	//
	// It combines three powerful functional programming concepts:
	//   1. Reader monad for dependency injection
	//   2. IO monad for side effects
	//   3. Either monad for error handling
	//
	// Type parameters:
	//   - R: The type of the context/environment (e.g., configuration, dependencies)
	//   - E: The type of errors that can occur
	//   - A: The type of the success value
	//
	// Example:
	//   type Config struct { BaseURL string }
	//   func fetchUser(id int) ReaderIOResult[Config, error, User] {
	//       return func(cfg Config) IOEither[error, User] {
	//           return func() Either[error, User] {
	//               // Use cfg.BaseURL to fetch user
	//               // Return either.Right(user) or either.Left(err)
	//           }
	//       }
	//   }
	ReaderIOResult[R, A any] = Reader[R, IOResult[A]]

	// Kleisli represents a Kleisli arrow for the ReaderIOResult monad.
	// It's a function from A to ReaderIOResult[R, B], used for composing operations that
	// depend on an environment, perform side effects, and may fail.
	Kleisli[R, A, B any] = reader.Reader[A, ReaderIOResult[R, B]]

	// Operator represents a transformation from one ReaderIOResult to another.
	// It's a Reader that takes a ReaderIOResult[R, A] and produces a ReaderIOResult[R, B].
	// This type is commonly used for composing operations in a point-free style.
	//
	// Type parameters:
	//   - R: The context type
	//   - A: The input value type
	//   - B: The output value type
	//
	// Example:
	//   var doubleOp Operator[Config, error, int, int] = Map(N.Mul(2))
	Operator[R, A, B any] = Kleisli[R, ReaderIOResult[R, A], B]

	// Consumer represents a function that consumes a value of type A.
	Consumer[A any] = consumer.Consumer[A]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	Predicate[A any] = predicate.Predicate[A]

	Void = function.Void
)
