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

// package readeriooption provides a monad transformer that combines the Reader and Option monads.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - Reader monad: https://github.com/fantasyland/fantasy-land
//   - Maybe (Option) monad: https://github.com/fantasyland/fantasy-land#maybe
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//   - Alt: https://github.com/fantasyland/fantasy-land#alt
//
// ReaderIOOption[R, A] represents a computation that:
//   - Depends on a shared environment of type R (Reader monad)
//   - May fail to produce a value of type A (Option monad)
//
// This is useful for computations that need access to configuration, context, or dependencies
// while also being able to represent the absence of a value without using errors.
//
// The ReaderIOOption monad is defined as: Reader[R, Option[A]]
//
// Key operations:
//   - Of: Wraps a value in a ReaderIOOption
//   - None: Creates a ReaderIOOption representing no value
//   - Map: Transforms the value inside a ReaderIOOption
//   - Chain: Sequences ReaderIOOption computations
//   - Ask/Asks: Accesses the environment
//
// Example:
//
//	type Config struct {
//	    DatabaseURL string
//	    Timeout     int
//	}
//
//	// A computation that may or may not find a user
//	func findUser(id int) readeriooption.ReaderIOOption[Config, User] {
//	    return readeriooption.Asks(func(cfg Config) option.Option[User] {
//	        // Use cfg.DatabaseURL to query database
//	        // Return Some(user) if found, None() if not found
//	    })
//	}
//
//	// Chain multiple operations
//	result := F.Pipe2(
//	    findUser(123),
//	    readeriooption.Chain(func(user User) readeriooption.ReaderIOOption[Config, Profile] {
//	        return loadProfile(user.ProfileID)
//	    }),
//	    readeriooption.Map(func(profile Profile) string {
//	        return profile.DisplayName
//	    }),
//	)
//
//	// Execute with config
//	config := Config{DatabaseURL: "localhost:5432", Timeout: 30}
//	displayName := result(config) // Returns Option[string]
package readeriooption

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
)

type (
	// Lazy represents a deferred computation that produces a value of type A.
	Lazy[A any] = lazy.Lazy[A]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's commonly used for filtering and conditional operations.
	Predicate[A any] = predicate.Predicate[A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// IOOption represents an IO computation that may produce a value of type A.
	// It combines IO effects with the Option monad for optional values.
	IOOption[A any] = iooption.IOOption[A]

	// Either represents a value of one of two possible types (a disjoint union).
	// An instance of Either is either Left (representing an error) or Right (representing a success).
	Either[E, A any] = either.Either[E, A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderIO represents a computation that depends on an environment R and performs IO to produce a value A.
	// It combines the Reader monad (for dependency injection) with IO effects.
	ReaderIO[R, A any] = readerio.ReaderIO[R, A]

	// ReaderIOOption represents a computation that depends on an environment R and may produce a value A.
	// It combines the Reader monad (for dependency injection) with IO effects and the Option monad (for optional values).
	// This is the main type of this package, defined as Reader[R, IOOption[A]].
	ReaderIOOption[R, A any] = Reader[R, IOOption[A]]

	// Kleisli represents a function that takes a value A and returns a ReaderIOOption[R, B].
	// This is the type of functions used with Chain/Bind operations, enabling monadic composition.
	Kleisli[R, A, B any] = Reader[A, ReaderIOOption[R, B]]

	// Operator represents a function that transforms one ReaderIOOption into another.
	// It takes a ReaderIOOption[R, A] and produces a ReaderIOOption[R, B].
	// This is commonly used for lifting functions into the ReaderIOOption context.
	Operator[R, A, B any] = Reader[ReaderIOOption[R, A], ReaderIOOption[R, B]]
)
