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

// Package readerresult provides a ReaderResult monad that combines the Reader and Result monads.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - Reader monad: https://github.com/fantasyland/fantasy-land
//   - Either monad: https://github.com/fantasyland/fantasy-land#either
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Bifunctor: https://github.com/fantasyland/fantasy-land#bifunctor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//   - Alt: https://github.com/fantasyland/fantasy-land#alt
//
// A ReaderResult[R, A] represents a computation that:
//   - Depends on an environment of type R (Reader aspect)
//   - May fail with an error (Result aspect, which is Either[error, A])
//
// This is equivalent to Reader[R, Result[A]] or Reader[R, Either[error, A]].
//
// # Use Cases
//
// ReaderResult is particularly useful for:
//
//  1. Dependency injection with error handling - pass configuration/services through
//     computations that may fail
//  2. Functional error handling - compose operations that depend on context and may error
//  3. Testing - easily mock dependencies by changing the environment value
//
// # Basic Example
//
//	type Config struct {
//	    DatabaseURL string
//	}
//
//	// Function that needs config and may fail
//	func getUser(id int) readerresult.ReaderResult[Config, User] {
//	    return readerresult.Asks(func(cfg Config) result.Result[User] {
//	        // Use cfg.DatabaseURL to fetch user
//	        return result.Of(user)
//	    })
//	}
//
//	// Execute by providing the config
//	cfg := Config{DatabaseURL: "postgres://..."}
//	res := getUser(42)(cfg)  // Returns result.Result[User]
//
// # Composition
//
// ReaderResult provides several ways to compose computations:
//
//  1. Map - transform successful values
//  2. Chain (FlatMap) - sequence dependent operations
//  3. Ap - combine independent computations
//  4. Do-notation - imperative-style composition with Bind
//
// # Do-Notation Example
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//
//	result := F.Pipe2(
//	    readerresult.Do[Config](State{}),
//	    readerresult.Bind(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) readerresult.ReaderResult[Config, User] {
//	            return getUser(42)
//	        },
//	    ),
//	    readerresult.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) readerresult.ReaderResult[Config, []Post] {
//	            return getPosts(s.User.ID)
//	        },
//	    ),
//	)
//
// # Error Handling
//
// ReaderResult provides several functions for error handling:
//
//   - Left/Right - create failed/successful values
//   - GetOrElse - provide a default value for errors
//   - OrElse - recover from errors with an alternative computation
//   - Fold - handle both success and failure cases
//   - ChainEitherK - lift result.Result computations into ReaderResult
//
// # Relationship to Other Monads
//
// ReaderResult is related to several other monads in this library:
//
//   - Reader[R, A] - ReaderResult without error handling
//   - Result[A] (Either[error, A]) - error handling without environment
//   - ReaderEither[R, E, A] - like ReaderResult but with custom error type E
//   - IOResult[A] - like ReaderResult but with no environment (IO with errors)
//
// # Performance Note
//
// ReaderResult is a zero-cost abstraction - it compiles to a simple function type
// with no runtime overhead beyond the underlying computation.
package readerresult
