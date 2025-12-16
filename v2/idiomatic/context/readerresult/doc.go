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

// Package readerresult provides a ReaderResult monad specialized for context.Context.
//
// A ReaderResult[A] represents an effectful computation that:
//   - Takes a context.Context as input
//   - May fail with an error (Result aspect, which is Either[error, A])
//   - Returns a value of type A on success
//
// The type is defined as: ReaderResult[A any] = func(context.Context) (A, error)
//
// This is equivalent to Reader[context.Context, Result[A]] or Reader[context.Context, Either[error, A]],
// but specialized to always use context.Context as the environment type.
//
// # Effectful Computations with Context
//
// ReaderResult is particularly well-suited for representing effectful computations in Go. An effectful
// computation is one that:
//
//   - Performs side effects (I/O, network calls, database operations, etc.)
//   - May fail with an error
//   - Requires contextual information (cancellation, deadlines, request-scoped values)
//
// By using context.Context as the fixed environment type, ReaderResult[A] provides:
//
//  1. Cancellation propagation - operations can be cancelled via context
//  2. Deadline/timeout handling - operations respect context deadlines
//  3. Request-scoped values - access to request metadata, trace IDs, etc.
//  4. Functional composition - chain effectful operations while maintaining context
//  5. Error handling - explicit error propagation through the Result type
//
// This pattern is idiomatic in Go, where functions performing I/O conventionally accept
// context.Context as their first parameter: func(ctx context.Context, ...) (Result, error).
// ReaderResult preserves this convention while enabling functional composition.
//
// Example of an effectful computation:
//
//	// An effectful operation that queries a database
//	func fetchUser(ctx context.Context, id int) (User, error) {
//	    // ctx provides cancellation, deadlines, and request context
//	    row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)
//	    var user User
//	    err := row.Scan(&user.ID, &user.Name)
//	    return user, err
//	}
//
//	// Lift into ReaderResult for functional composition
//	getUser := readerresult.Curry1(fetchUser)
//
//	// Compose multiple effectful operations
//	pipeline := F.Pipe2(
//	    getUser(42),  // ReaderResult[User]
//	    readerresult.Chain(func(user User) readerresult.ReaderResult[[]Post] {
//	        return getPosts(user.ID)  // Another effectful operation
//	    }),
//	)
//
//	// Execute with a context (e.g., from an HTTP request)
//	ctx := r.Context()  // HTTP request context
//	posts, err := pipeline(ctx)
//
// # Use Cases
//
// ReaderResult is particularly useful for:
//
//  1. Effectful computations with context - operations that perform I/O and need cancellation/deadlines
//  2. Functional error handling - compose operations that depend on context and may error
//  3. Testing - easily mock context-dependent operations
//  4. HTTP handlers - chain request processing operations with proper context propagation
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
//	    readerresult.Do(State{}),
//	    readerresult.Bind(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) readerresult.ReaderResult[User] {
//	            return getUser(42)
//	        },
//	    ),
//	    readerresult.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) readerresult.ReaderResult[[]Post] {
//	            return getPosts(s.User.ID)
//	        },
//	    ),
//	)
//
// # Currying Functions with Context
//
// The Curry functions enable partial application of function parameters while deferring
// the context.Context parameter until execution time.
//
// When you curry a function like func(context.Context, T1, T2) (A, error), the context.Context
// becomes the last argument to be applied, even though it appears first in the original function
// signature. This is intentional and follows Go's context-first convention while enabling
// functional composition patterns.
//
// Why context.Context is the last curried argument:
//
//   - In Go, context conventionally comes first: func(ctx context.Context, params...) (Result, error)
//   - In curried form: Curry2(f)(param1)(param2) returns ReaderResult[A]
//   - The ReaderResult is then applied to ctx: Curry2(f)(param1)(param2)(ctx)
//   - This allows partial application of business parameters before providing the context
//
// Example with database operations:
//
//	// Database operations following Go conventions (context first)
//	func fetchUser(ctx context.Context, db *sql.DB, id int) (User, error) {
//	    row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)
//	    var user User
//	    err := row.Scan(&user.ID, &user.Name)
//	    return user, err
//	}
//
//	func updateUser(ctx context.Context, db *sql.DB, id int, name string) (User, error) {
//	    _, err := db.ExecContext(ctx, "UPDATE users SET name = ? WHERE id = ?", name, id)
//	    if err != nil {
//	        return User{}, err
//	    }
//	    return fetchUser(ctx, db, id)
//	}
//
//	// Curry these into composable operations
//	getUser := readerresult.Curry2(fetchUser)
//	updateUserName := readerresult.Curry3(updateUser)
//
//	// Compose operations with partial application
//	pipeline := F.Pipe2(
//	    getUser(db)(42),  // ReaderResult[User] - db and id applied, waiting for ctx
//	    readerresult.Chain(func(user User) readerresult.ReaderResult[User] {
//	        newName := user.Name + " (updated)"
//	        return updateUserName(db)(user.ID)(newName)  // Waiting for ctx
//	    }),
//	)
//
//	// Execute by providing the context
//	ctx := context.Background()
//	updatedUser, err := pipeline(ctx)
//
// The key insight is that currying creates a chain where:
//  1. Business parameters are applied first: getUser(db)(42)
//  2. This returns a ReaderResult[User] that waits for the context
//  3. Multiple operations can be composed before providing the context
//  4. Finally, the context is provided to execute everything: pipeline(ctx)
//
// This pattern is particularly useful for:
//   - Creating reusable operation pipelines independent of specific contexts
//   - Testing with different contexts (with timeouts, cancellation, etc.)
//   - Composing operations that share the same context
//   - Deferring context creation until execution time
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
//   - Reader[context.Context, A] - ReaderResult without error handling
//   - Result[A] (Either[error, A]) - error handling without context dependency
//   - IOResult[A] - similar to ReaderResult but without explicit context parameter
//   - ReaderIOResult[R, A] - generic version that allows custom environment type R
//
// # Performance Note
//
// ReaderResult is a zero-cost abstraction - it compiles to a simple function type
// with no runtime overhead beyond the underlying computation.
package readerresult
