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
//	user, err := getUser(42)(cfg)  // Returns (User, error)
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
// # Object-Oriented Patterns with Curry Functions
//
// The Curry functions enable an interesting pattern where you can treat the Reader context (R)
// as an object instance, effectively creating method-like functions that compose functionally.
//
// When you curry a function like func(R, T1, T2) (A, error), the context R becomes the last
// argument to be applied, even though it appears first in the original function signature.
// This is intentional and follows Go's context-first convention while enabling functional
// composition patterns.
//
// Why R is the last curried argument:
//
//   - In Go, context conventionally comes first: func(ctx Context, params...) (Result, error)
//   - In curried form: Curry2(f)(param1)(param2) returns ReaderResult[R, A]
//   - The ReaderResult is then applied to R: Curry2(f)(param1)(param2)(ctx)
//   - This allows partial application of business parameters before providing the context/object
//
// Object-Oriented Example:
//
//	// A service struct that acts as the Reader context
//	type UserService struct {
//	    db *sql.DB
//	    cache Cache
//	}
//
//	// A method-like function following Go conventions (context first)
//	func (s *UserService) GetUserByID(ctx context.Context, id int) (User, error) {
//	    // Use s.db and s.cache...
//	}
//
//	func (s *UserService) UpdateUser(ctx context.Context, id int, name string) (User, error) {
//	    // Use s.db and s.cache...
//	}
//
//	// Curry these into composable operations
//	getUser := readerresult.Curry1((*UserService).GetUserByID)
//	updateUser := readerresult.Curry2((*UserService).UpdateUser)
//
//	// Now compose operations that will be bound to a UserService instance
//	type Context struct {
//	    Svc *UserService
//	}
//
//	pipeline := F.Pipe2(
//	    getUser(42),  // ReaderResult[Context, User]
//	    readerresult.Chain(func(user User) readerresult.ReaderResult[Context, User] {
//	        newName := user.Name + " (updated)"
//	        return updateUser(user.ID)(newName)
//	    }),
//	)
//
//	// Execute by providing the service instance as context
//	svc := &UserService{db: db, cache: cache}
//	ctx := Context{Svc: svc}
//	updatedUser, err := pipeline(ctx)
//
// The key insight is that currying creates a chain where:
//  1. Business parameters are applied first: getUser(42)
//  2. This returns a ReaderResult that waits for the context
//  3. Multiple operations can be composed before providing the context
//  4. Finally, the context/object is provided to execute everything: pipeline(ctx)
//
// This pattern is particularly useful for:
//   - Creating reusable operation pipelines independent of service instances
//   - Testing with mock service instances
//   - Dependency injection in a functional style
//   - Composing operations that share the same service context
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
