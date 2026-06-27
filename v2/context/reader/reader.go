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

// Package reader provides a specialization of the Reader monad for [context.Context].
//
// This package offers a context-aware Reader monad that simplifies working with
// Go's [context.Context] in a functional programming style. It eliminates the need
// to explicitly thread context through function calls while maintaining type safety
// and composability.
//
// # Core Concept
//
// The Reader monad represents computations that depend on a shared environment.
// In this package, that environment is fixed to [context.Context], making it
// particularly useful for:
//
//   - Request-scoped data propagation
//   - Cancellation and timeout handling
//   - Dependency injection via context values
//   - Avoiding explicit context parameter threading
//
// # Type Definitions
//
//   - Reader[A]: A computation that depends on context.Context and produces A
//   - Kleisli[A, B]: A function from A to Reader[B] for composing computations
//   - Operator[A, B]: A transformation from Reader[A] to Reader[B]
//
// # Usage Pattern
//
// Instead of passing context explicitly through every function:
//
//	func processUser(ctx context.Context, userID string) (User, error) {
//	    user := fetchUser(ctx, userID)
//	    profile := fetchProfile(ctx, user.ProfileID)
//	    return enrichUser(ctx, user, profile), nil
//	}
//
// You can use Reader to compose context-dependent operations:
//
//	fetchUser := func(userID string) Reader[User] {
//	    return func(ctx context.Context) User {
//	        // Use ctx for database access, cancellation, etc.
//	        return queryDatabase(ctx, userID)
//	    }
//	}
//
//	processUser := func(userID string) Reader[User] {
//	    return F.Pipe2(
//	        fetchUser(userID),
//	        reader.Chain(func(user User) Reader[Profile] {
//	            return fetchProfile(user.ProfileID)
//	        }),
//	        reader.Map(func(profile Profile) User {
//	            return enrichUser(user, profile)
//	        }),
//	    )
//	}
//
//	// Execute with context
//	ctx := context.Background()
//	user := processUser("user123")(ctx)
//
// # Integration with Standard Library
//
// This package works seamlessly with Go's standard [context] package:
//
//   - Context cancellation and deadlines are preserved
//   - Context values can be accessed within Reader computations
//   - Readers can be composed with context-aware libraries
//
// # Relationship to Other Packages
//
// This package is a specialization of [github.com/IBM/fp-go/v2/reader] where
// the environment type R is fixed to [context.Context]. For more general
// Reader operations, see the base reader package.
//
// For combining Reader with other monads:
//   - [github.com/IBM/fp-go/v2/context/readerio]: Reader + IO effects
//   - [github.com/IBM/fp-go/v2/readeroption]: Reader + Option
//   - [github.com/IBM/fp-go/v2/readerresult]: Reader + Result (Either)
//
// # Example: HTTP Request Handler
//
//	type RequestContext struct {
//	    UserID    string
//	    RequestID string
//	}
//
//	// Extract request context from context.Context
//	getRequestContext := func(ctx context.Context) RequestContext {
//	    return RequestContext{
//	        UserID:    ctx.Value("userID").(string),
//	        RequestID: ctx.Value("requestID").(string),
//	    }
//	}
//
//	// A Reader that logs with request context
//	logInfo := func(message string) Reader[function.Void] {
//	    return func(ctx context.Context) function.Void {
//	        reqCtx := getRequestContext(ctx)
//	        log.Printf("[%s] User %s: %s", reqCtx.RequestID, reqCtx.UserID, message)
//	        return function.VOID
//	    }
//	}
//
//	// Compose operations
//	handleRequest := func(data string) Reader[Response] {
//	    return F.Pipe2(
//	        logInfo("Processing request"),
//	        reader.Chain(func(_ function.Void) Reader[Result] {
//	            return processData(data)
//	        }),
//	        reader.Map(func(result Result) Response {
//	            return Response{Data: result}
//	        }),
//	    )
//	}
package reader

import (
	"context"

	IC "github.com/IBM/fp-go/v2/internal/context"
)

// WithValue creates a Kleisli arrow that adds a value to the context.
//
// This function provides a functional way to add values to a context.Context,
// returning a new context with the key-value pair added. It's particularly useful
// for building context-dependent computations that need to propagate values through
// a chain of operations.
//
// Type Parameters:
//   - A: The type of the value to store in the context
//   - K: The type of the key (typically string or a custom type)
//
// Parameters:
//   - key: The key to associate with the value in the context
//
// Returns:
//   - Kleisli[A, context.Context]: A function that takes a value and returns a Reader
//     that produces a new context with the key-value pair added
//
// Example:
//
//	import (
//	    "context"
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/context/reader"
//	)
//
//	// Create a Kleisli arrow for adding a user ID to context
//	setUserID := reader.WithValue[string, string]("userID")
//
//	// Use it to create a context with a user ID
//	ctx := context.Background()
//	newCtx := setUserID("user123")(ctx)
//	userID := newCtx.Value("userID").(string) // "user123"
//
// Example: Chaining multiple context values
//
//	import (
//	    R "github.com/IBM/fp-go/v2/reader"
//	)
//
//	// Chain multiple WithValue operations
//	enrichContext := F.Pipe2(
//	    reader.WithValue[string, string]("userID")("user123"),
//	    R.Chain(reader.WithValue[string, string]("requestID")("req456")),
//	    R.Chain(reader.WithValue[int, string]("timeout")(30)),
//	)
//
//	ctx := context.Background()
//	enrichedCtx := enrichContext(ctx)
//	// enrichedCtx now contains userID, requestID, and timeout
//
// Example: Using with custom key types
//
//	type contextKey string
//
//	const (
//	    userKey    contextKey = "user"
//	    sessionKey contextKey = "session"
//	)
//
//	type User struct {
//	    ID   string
//	    Name string
//	}
//
//	// Type-safe context value setting
//	setUser := reader.WithValue[User, contextKey](userKey)
//	setSession := reader.WithValue[string, contextKey](sessionKey)
//
//	user := User{ID: "123", Name: "Alice"}
//	ctx := F.Pipe1(
//	    setUser(user),
//	    R.Chain(setSession("session-token")),
//	)(context.Background())
//
// Notes:
//   - The returned context is a new context; the original is not modified
//   - Keys should be comparable types (typically string or custom types)
//   - For type safety, consider using custom key types instead of strings
//   - Values stored in context should be request-scoped and immutable
//
// See Also:
//   - context.WithValue: The underlying standard library function
//   - Reader: For composing context-dependent computations
//   - Kleisli: For understanding Kleisli arrow composition
func WithValue[A, K any](key K) Kleisli[A, context.Context] {
	return IC.WithValue[A](key)
}

// NopCancel wraps a context in a ContextCancel whose cancel function is a no-op.
//
// The returned ContextCancel pairs the given context with a cancel function that
// does nothing when called. This is useful when an API requires a ContextCancel
// but no actual cancellation is needed — for example, when adapting a plain
// context.Context to a function that expects a cancellable context pair.
//
// The name is intentionally analogous to io.NopCloser, which wraps an io.Reader
// in an io.ReadCloser whose Close method is also a no-op.
//
// Parameters:
//   - ctx: The context to wrap. It is returned unchanged as the second element
//     of the pair.
//
// Returns:
//   - ContextCancel: A pair whose first element is a no-op context.CancelFunc
//     and whose second element is ctx.
//
// See Also:
//   - io.NopCloser: The standard library analogue for io.ReadCloser
//   - ContextCancel: The pair type returned by this function
func NopCancel(ctx context.Context) ContextCancel {
	return IC.NopCancel(ctx)
}
