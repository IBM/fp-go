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
