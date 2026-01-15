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

package readerioeither

import (
	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/lens/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readeroption"
)

type (
	// Either represents a value of one of two possible types (a disjoint union).
	// An instance of Either is either Left (representing an error) or Right (representing a success).
	Either[E, A any] = either.Either[E, A]

	// Reader represents a computation that depends on some context/environment of type R
	// and produces a value of type A. It's useful for dependency injection patterns.
	Reader[R, A any] = reader.Reader[R, A]

	// IO represents a synchronous computation that cannot fail.
	IO[T any] = io.IO[T]

	// ReaderIO represents a computation that depends on some context R and performs
	// side effects to produce a value of type A.
	ReaderIO[R, A any] = readerio.ReaderIO[R, A]

	// IOEither represents a computation that performs side effects and can either
	// fail with an error of type E or succeed with a value of type A.
	IOEither[E, A any] = ioeither.IOEither[E, A]

	// ReaderEither represents a computation that depends on an environment R and can fail
	// with an error E or succeed with a value A (without side effects).
	ReaderEither[R, E, A any] = readereither.ReaderEither[R, E, A]

	// ReaderIOEither represents a computation that:
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
	//   func fetchUser(id int) ReaderIOEither[Config, error, User] {
	//       return func(cfg Config) IOEither[error, User] {
	//           return func() Either[error, User] {
	//               // Use cfg.BaseURL to fetch user
	//               // Return either.Right(user) or either.Left(err)
	//           }
	//       }
	//   }
	ReaderIOEither[R, E, A any] = Reader[R, IOEither[E, A]]

	// Kleisli represents a Kleisli arrow for the ReaderIOEither monad.
	// It's a function from A to ReaderIOEither[R, E, B], used for composing operations that
	// depend on an environment, perform side effects, and may fail.
	Kleisli[R, E, A, B any] = reader.Reader[A, ReaderIOEither[R, E, B]]

	// Operator represents a transformation from one ReaderIOEither to another.
	// It's a Reader that takes a ReaderIOEither[R, E, A] and produces a ReaderIOEither[R, E, B].
	// This type is commonly used for composing operations in a point-free style.
	//
	// Type parameters:
	//   - R: The context type
	//   - E: The error type
	//   - A: The input value type
	//   - B: The output value type
	//
	// Example:
	//   var doubleOp Operator[Config, error, int, int] = Map(N.Mul(2))
	Operator[R, E, A, B any] = Kleisli[R, E, ReaderIOEither[R, E, A], B]

	// ReaderOption represents a computation that depends on an environment R and may not produce a value.
	ReaderOption[R, A any] = readeroption.ReaderOption[R, A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Consumer represents a function that consumes a value of type A.
	Consumer[A any] = consumer.Consumer[A]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	Predicate[A any] = predicate.Predicate[A]

	Lazy[A any] = lazy.Lazy[A]
)
