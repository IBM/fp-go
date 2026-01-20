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

package readerreaderioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readerioeither"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// Option represents an optional value that may or may not be present.
	// It's an alias for option.Option[A].
	Option[A any] = option.Option[A]

	// Lazy represents a lazily evaluated computation that produces a value of type A.
	// It's an alias for lazy.Lazy[A].
	Lazy[A any] = lazy.Lazy[A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	// It's an alias for reader.Reader[R, A].
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderOption represents a computation that depends on an environment R and produces
	// an optional value A. It's an alias for readeroption.ReaderOption[R, A].
	ReaderOption[R, A any] = readeroption.ReaderOption[R, A]

	// ReaderIO represents a computation that depends on an environment R and performs
	// side effects to produce a value A. It's an alias for readerio.ReaderIO[R, A].
	ReaderIO[R, A any] = readerio.ReaderIO[R, A]

	// ReaderIOEither represents a computation that depends on an environment R, performs
	// side effects, and produces either an error E or a value A.
	// It's an alias for readerioeither.ReaderIOEither[R, E, A].
	ReaderIOEither[R, E, A any] = readerioeither.ReaderIOEither[R, E, A]

	// ReaderEither represents a computation that depends on an environment R and produces
	// either an error E or a value A. It's an alias for readereither.ReaderEither[R, E, A].
	ReaderEither[R, E, A any] = readereither.ReaderEither[R, E, A]

	// Either represents a value that is either a Left (error) E or a Right (success) A.
	// It's an alias for either.Either[E, A].
	Either[E, A any] = either.Either[E, A]

	// IOEither represents a side-effecting computation that produces either an error E
	// or a value A. It's an alias for ioeither.IOEither[E, A].
	IOEither[E, A any] = ioeither.IOEither[E, A]

	// IO represents a side-effecting computation that produces a value A.
	// It's an alias for io.IO[A].
	IO[A any] = io.IO[A]

	// ReaderReaderIOEither represents a nested Reader monad transformer that:
	// 1. Takes an outer environment R
	// 2. Returns a ReaderIOEither that takes an inner environment C
	// 3. Performs side effects
	// 4. Produces either an error E or a value A
	//
	// This is the core type of this package, enabling computations with two levels
	// of environment dependencies, side effects, and error handling.
	//
	// Type Parameters:
	//   - R: The outer environment type (first Reader layer)
	//   - C: The inner environment type (ReaderIOEither layer)
	//   - E: The error type
	//   - A: The success value type
	//
	// Example:
	//
	//	type OuterConfig struct { DatabaseURL string }
	//	type InnerContext struct { RequestID string }
	//
	//	computation := func(outer OuterConfig) readerioeither.ReaderIOEither[InnerContext, error, string] {
	//	    return func(inner InnerContext) ioeither.IOEither[error, string] {
	//	        return func() either.Either[error, string] {
	//	            return either.Right[error](fmt.Sprintf("DB: %s, Request: %s",
	//	                outer.DatabaseURL, inner.RequestID))
	//	        }
	//	    }
	//	}
	ReaderReaderIOEither[R, C, E, A any] = Reader[R, ReaderIOEither[C, E, A]]

	// Kleisli represents a Kleisli arrow for ReaderReaderIOEither.
	// It's a function that takes a value A and returns a ReaderReaderIOEither[R, C, E, B].
	//
	// Kleisli arrows are used for monadic composition, allowing you to chain
	// computations that depend on two environments, perform side effects, and may fail.
	//
	// Type Parameters:
	//   - R: The outer environment type
	//   - C: The inner environment type
	//   - E: The error type
	//   - A: The input value type
	//   - B: The output value type
	//
	// Example:
	//
	//	// A Kleisli arrow that validates and processes a user ID
	//	validateUser := func(userID int) ReaderReaderIOEither[Config, Context, error, User] {
	//	    return func(cfg Config) readerioeither.ReaderIOEither[Context, error, User] {
	//	        return func(ctx Context) ioeither.IOEither[error, User] {
	//	            return func() either.Either[error, User] {
	//	                if userID <= 0 {
	//	                    return either.Left[User](errors.New("invalid user ID"))
	//	                }
	//	                return either.Right[error](User{ID: userID})
	//	            }
	//	        }
	//	    }
	//	}
	Kleisli[R, C, E, A, B any] = Reader[A, ReaderReaderIOEither[R, C, E, B]]

	// Operator represents an endomorphism in the Kleisli category for ReaderReaderIOEither.
	// It's a Kleisli arrow where the input is itself a ReaderReaderIOEither[R, C, E, A].
	//
	// Operators are useful for transforming or enhancing existing computations,
	// such as adding logging, retry logic, or resource management.
	//
	// Type Parameters:
	//   - R: The outer environment type
	//   - C: The inner environment type
	//   - E: The error type
	//   - A: The input computation's value type
	//   - B: The output value type
	//
	// Example:
	//
	//	// An operator that adds retry logic to a computation
	//	withRetry := func(computation ReaderReaderIOEither[Config, Context, error, int]) ReaderReaderIOEither[Config, Context, error, int] {
	//	    return func(cfg Config) readerioeither.ReaderIOEither[Context, error, int] {
	//	        return func(ctx Context) ioeither.IOEither[error, int] {
	//	            return func() either.Either[error, int] {
	//	                // Retry logic here
	//	                return computation(cfg)(ctx)()
	//	            }
	//	        }
	//	    }
	//	}
	Operator[R, C, E, A, B any] = Kleisli[R, C, E, ReaderReaderIOEither[R, C, E, A], B]

	// Lens represents an optic for focusing on a part of a data structure.
	// It's an alias for lens.Lens[S, T].
	Lens[S, T any] = lens.Lens[S, T]

	// Trampoline represents a tail-recursive computation that can be executed
	// without stack overflow. It's an alias for tailrec.Trampoline[L, B].
	Trampoline[L, B any] = tailrec.Trampoline[L, B]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's an alias for predicate.Predicate[A].
	Predicate[A any] = predicate.Predicate[A]
)
