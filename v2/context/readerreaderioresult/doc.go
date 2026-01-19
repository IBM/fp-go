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

// Package readerreaderioresult provides a functional programming abstraction that combines
// four powerful concepts: Reader, Reader, IO, and Result (Either[error, A]) monads in a nested structure.
// This is a specialized version of readerreaderioeither where the error type is fixed to `error` and
// the inner context is fixed to `context.Context`.
//
// # Type Definition
//
// ReaderReaderIOResult[R, A] is defined as:
//
//	type ReaderReaderIOResult[R, A] = ReaderReaderIOEither[R, context.Context, error, A]
//
// Which expands to:
//
//	func(R) func(context.Context) func() Either[error, A]
//
// This represents a computation that:
//   - Takes an outer environment/context of type R
//   - Returns a function that takes a context.Context
//   - Returns an IO operation (a thunk/function with no parameters)
//   - Produces an Either[error, A] (Result[A]) when executed
//
// # Type Parameter Ordering Convention
//
// This package follows a consistent convention for ordering type parameters in function signatures.
// The general rule is: R -> C -> E -> T (outer context, inner context, error, type), where:
//   - R: The outer Reader context/environment type
//   - C: The inner Reader context/environment type (for the ReaderIOEither)
//   - E: The Either error type
//   - T: The value type(s) (A, B, etc.)
//
// However, when some type parameters can be automatically inferred by the Go compiler from
// function arguments, the convention is modified to minimize explicit type annotations:
//
// Rule: Undetectable types come first, followed by detectable types, while preserving
// the relative order within each group (R -> C -> E -> T).
//
// Examples:
//
//  1. All types detectable from first argument:
//     MonadMap[R, C, E, A, B](fa ReaderReaderIOEither[R, C, E, A], f func(A) B)
//     - R, C, E, A are detectable from fa
//     - B is detectable from f
//     - Order: R, C, E, A, B (standard order, all detectable)
//
//  2. Some types undetectable:
//     FromReader[C, E, R, A](ma Reader[R, A]) ReaderReaderIOEither[R, C, E, A]
//     - R, A are detectable from ma
//     - C, E are undetectable (not in any argument)
//     - Order: C, E, R, A (C, E first as undetectable, then R, A in standard order)
//
//  3. Multiple undetectable types:
//     Local[C, E, A, R1, R2](f func(R2) R1) func(ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A]
//     - C, E, A are undetectable
//     - R1, R2 are detectable from f
//     - Order: C, E, A, R1, R2 (undetectable first, then detectable)
//
//  4. Functions returning Kleisli arrows:
//     ChainReaderOptionK[R, C, A, B, E](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, C, E, A, B]
//     - Canonical order would be R, C, E, A, B
//     - E is detectable from onNone parameter
//     - R, C, A, B are not detectable (they're in the Kleisli argument type)
//     - Order: R, C, A, B, E (undetectable R, C, A, B first, then detectable E)
//
// This convention allows for more ergonomic function calls:
//
//	// Without convention - need to specify all types:
//	result := FromReader[OuterCtx, InnerCtx, error, User](readerFunc)
//
//	// With convention - only specify undetectable types:
//	result := FromReader[InnerCtx, error](readerFunc)  // R and A inferred from readerFunc
//
// The reasoning behind this approach is to reduce the number of explicit type parameters
// that developers need to specify when calling functions, improving code readability and
// reducing verbosity while maintaining type safety.
//
// Additional examples demonstrating the convention:
//
//  5. FromReaderOption[R, C, A, E](onNone Lazy[E]) Kleisli[R, C, E, ReaderOption[R, A], A]
//     - Canonical order would be R, C, E, A
//     - E is detectable from onNone parameter
//     - R, C, A are not detectable (they're in the return type's Kleisli)
//     - Order: R, C, A, E (undetectable R, C, A first, then detectable E)
//
//  6. MapLeft[R, C, A, E1, E2](f func(E1) E2) func(ReaderReaderIOEither[R, C, E1, A]) ReaderReaderIOEither[R, C, E2, A]
//     - Canonical order would be R, C, E1, E2, A
//     - E1, E2 are detectable from f parameter
//     - R, C, A are not detectable (they're in the return type)
//     - Order: R, C, A, E1, E2 (undetectable R, C, A first, then detectable E1, E2)
//
// Additional special cases:
//
//   - Ap[B, R, C, E, A]: B is undetectable (in function return type), so B comes first
//   - ChainOptionK[R, C, A, B, E]: R, C, A, B are undetectable, E is detectable from onNone
//   - FromReaderIO[C, E, R, A]: C, E are undetectable, R, A are detectable from ReaderIO[R, A]
//
// All functions in this package follow this convention consistently.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - Reader monad: https://github.com/fantasyland/fantasy-land
//   - Reader monad (nested): https://github.com/fantasyland/fantasy-land
//   - IO monad: https://github.com/fantasyland/fantasy-land
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
// # ReaderReaderIOEither
//
// ReaderReaderIOEither[R, C, E, A] represents a computation that:
//   - Depends on an outer context/environment of type R (outer Reader)
//   - Returns a computation that depends on an inner context/environment of type C (inner Reader)
//   - Performs side effects (IO)
//   - Can fail with an error of type E or succeed with a value of type A (Either)
//
// This is particularly useful for:
//   - Multi-level dependency injection patterns
//   - Layered architectures with different context requirements at each layer
//   - Composing operations that need access to multiple levels of configuration or context
//   - Building reusable components that can be configured at different stages
//
// # Core Operations
//
// Construction:
//   - Of/Right: Create a successful computation
//   - Left: Create a failed computation
//   - FromEither: Lift an Either into ReaderReaderIOEither
//   - FromIO: Lift an IO into ReaderReaderIOEither
//   - FromReader: Lift a Reader into ReaderReaderIOEither
//   - FromReaderIO: Lift a ReaderIO into ReaderReaderIOEither
//   - FromIOEither: Lift an IOEither into ReaderReaderIOEither
//   - FromReaderEither: Lift a ReaderEither into ReaderReaderIOEither
//   - FromReaderIOEither: Lift a ReaderIOEither into ReaderReaderIOEither
//   - FromReaderOption: Lift a ReaderOption into ReaderReaderIOEither
//
// Transformation:
//   - Map: Transform the success value
//   - MapLeft: Transform the error value
//   - Chain/Bind: Sequence dependent computations
//   - Flatten: Flatten nested ReaderReaderIOEither
//
// Combination:
//   - Ap: Apply a function in a context to a value in a context
//   - ApSeq: Sequential application
//   - ApPar: Parallel application
//
// Error Handling:
//   - Alt: Choose the first successful computation
//
// Context Access:
//   - Ask: Get the current outer context
//   - Asks: Get a value derived from the outer context
//   - Local: Run a computation with a modified outer context
//   - Read: Execute with a specific outer context
//
// Kleisli Composition:
//   - ChainEitherK: Chain with Either-returning functions
//   - ChainReaderK: Chain with Reader-returning functions
//   - ChainReaderIOK: Chain with ReaderIO-returning functions
//   - ChainReaderEitherK: Chain with ReaderEither-returning functions
//   - ChainReaderOptionK: Chain with ReaderOption-returning functions
//   - ChainIOEitherK: Chain with IOEither-returning functions
//   - ChainIOK: Chain with IO-returning functions
//   - ChainOptionK: Chain with Option-returning functions
//
// First/Tap Operations (execute for side effects, return original value):
//   - ChainFirst/Tap: Execute a computation but return the original value
//   - ChainFirstEitherK/TapEitherK: Tap with Either-returning functions
//   - ChainFirstReaderK/TapReaderK: Tap with Reader-returning functions
//   - ChainFirstReaderIOK/TapReaderIOK: Tap with ReaderIO-returning functions
//   - ChainFirstReaderEitherK/TapReaderEitherK: Tap with ReaderEither-returning functions
//   - ChainFirstReaderOptionK/TapReaderOptionK: Tap with ReaderOption-returning functions
//   - ChainFirstIOK/TapIOK: Tap with IO-returning functions
//
// # Example Usage
//
//	type AppConfig struct {
//	    DatabaseURL string
//	    LogLevel    string
//	}
//
//	// A computation that depends on AppConfig and context.Context
//	func fetchUser(id int) ReaderReaderIOResult[AppConfig, User] {
//	    return func(cfg AppConfig) readerioresult.ReaderIOResult[context.Context, User] {
//	        // Use cfg.DatabaseURL and cfg.LogLevel
//	        return func(ctx context.Context) ioresult.IOResult[User] {
//	            // Use ctx for cancellation/timeout
//	            return func() result.Result[User] {
//	                // Perform the actual IO operation
//	                // Return result.Of(user) or result.Error[User](err)
//	            }
//	        }
//	    }
//	}
//
//	// Compose operations
//	result := function.Pipe2(
//	    fetchUser(123),
//	    Map[AppConfig](func(u User) string { return u.Name }),
//	    Chain[AppConfig](func(name string) ReaderReaderIOResult[AppConfig, string] {
//	        return Of[AppConfig]("Hello, " + name)
//	    }),
//	)
//
//	// Execute with config and context
//	appConfig := AppConfig{DatabaseURL: "postgres://...", LogLevel: "info"}
//	ctx := t.Context()
//	outcome := result(appConfig)(ctx)() // Returns result.Result[string]
//
// # Use Cases
//
// This monad is particularly useful for:
//   - Applications with layered configuration (app config + request context)
//   - HTTP handlers that need both application config and request context
//   - Database operations with connection pool config and query context
//   - Retry logic with policy configuration and execution context
//   - Resource management with bracket pattern across multiple contexts
//
// # Dependency Injection with the Outer Context
//
// The outer Reader context (type parameter R) provides a powerful mechanism for dependency injection
// in functional programming. This pattern is explained in detail in Scott Wlaschin's talk:
// "Dependency Injection, The Functional Way" - https://www.youtube.com/watch?v=xPlsVVaMoB0
//
// ## Core Concept
//
// Instead of using traditional OOP dependency injection frameworks, the Reader monad allows you to:
//  1. Define functions that declare their dependencies as type parameters
//  2. Compose these functions without providing the dependencies
//  3. Supply all dependencies at the "end of the world" (program entry point)
//
// This approach provides:
//   - Compile-time safety: Missing dependencies cause compilation errors
//   - Explicit dependencies: Function signatures show exactly what they need
//   - Easy testing: Mock dependencies by providing different values
//   - Pure functions: Dependencies are passed as parameters, not global state
//
// ## Examples from the Video Adapted to fp-go
//
// ### Example 1: Basic Reader Pattern (Video: "Reader Monad Basics")
//
// In the video, Scott shows how to pass configuration through a chain of functions.
// In fp-go with ReaderReaderIOResult:
//
//	// Define your dependencies
//	type AppConfig struct {
//	    DatabaseURL string
//	    APIKey      string
//	    MaxRetries  int
//	}
//
//	// Functions declare their dependencies via the R type parameter
//	func getConnectionString() ReaderReaderIOResult[AppConfig, string] {
//	    return Asks[AppConfig](func(cfg AppConfig) string {
//	        return cfg.DatabaseURL
//	    })
//	}
//
//	func connectToDatabase() ReaderReaderIOResult[AppConfig, *sql.DB] {
//	    return MonadChain(
//	        getConnectionString(),
//	        func(connStr string) ReaderReaderIOResult[AppConfig, *sql.DB] {
//	            return FromIO[AppConfig](func() result.Result[*sql.DB] {
//	                db, err := sql.Open("postgres", connStr)
//	                return result.FromEither(either.FromError(db, err))
//	            })
//	        },
//	    )
//	}
//
// ### Example 2: Composing Dependencies (Video: "Composing Reader Functions")
//
// The video demonstrates how Reader functions compose naturally.
// In fp-go, you can compose operations that all share the same dependency:
//
//	func fetchUser(id int) ReaderReaderIOResult[AppConfig, User] {
//	    return MonadChain(
//	        connectToDatabase(),
//	        func(db *sql.DB) ReaderReaderIOResult[AppConfig, User] {
//	            return FromIO[AppConfig](func() result.Result[User] {
//	                // Query database using db and return user
//	                // The AppConfig is still available if needed
//	            })
//	        },
//	    )
//	}
//
//	func enrichUser(user User) ReaderReaderIOResult[AppConfig, EnrichedUser] {
//	    return Asks[AppConfig, EnrichedUser](func(cfg AppConfig) EnrichedUser {
//	        // Use cfg.APIKey to call external service
//	        return EnrichedUser{User: user, Extra: "data"}
//	    })
//	}
//
//	// Compose without providing dependencies
//	pipeline := function.Pipe2(
//	    fetchUser(123),
//	    Chain[AppConfig](enrichUser),
//	)
//
//	// Provide dependencies at the end
//	config := AppConfig{DatabaseURL: "...", APIKey: "...", MaxRetries: 3}
//	ctx := context.Background()
//	result := pipeline(config)(ctx)()
//
// ### Example 3: Local Context Modification (Video: "Local Environment")
//
// The video shows how to temporarily modify the environment for a sub-computation.
// In fp-go, use the Local function:
//
//	// Run a computation with modified configuration
//	func withRetries(retries int, action ReaderReaderIOResult[AppConfig, string]) ReaderReaderIOResult[AppConfig, string] {
//	    return Local[string](func(cfg AppConfig) AppConfig {
//	        // Create a modified config with different retry count
//	        return AppConfig{
//	            DatabaseURL: cfg.DatabaseURL,
//	            APIKey:      cfg.APIKey,
//	            MaxRetries:  retries,
//	        }
//	    })(action)
//	}
//
//	// Use it
//	result := withRetries(5, fetchUser(123))
//
// ### Example 4: Testing with Mock Dependencies (Video: "Testing with Reader")
//
// The video emphasizes how Reader makes testing easy by allowing mock dependencies.
// In fp-go:
//
//	func TestFetchUser(t *testing.T) {
//	    // Create a test configuration
//	    testConfig := AppConfig{
//	        DatabaseURL: "mock://test",
//	        APIKey:      "test-key",
//	        MaxRetries:  1,
//	    }
//
//	    // Run the computation with test config
//	    ctx := context.Background()
//	    result := fetchUser(123)(testConfig)(ctx)()
//
//	    // Assert on the result
//	    assert.True(t, either.IsRight(result))
//	}
//
// ### Example 5: Multi-Layer Dependencies (Video: "Nested Readers")
//
// The video discusses nested readers for multi-layer architectures.
// ReaderReaderIOResult provides exactly this with R (outer) and context.Context (inner):
//
//	type AppConfig struct {
//	    DatabaseURL string
//	}
//
//	// Outer context: Application-level configuration (AppConfig)
//	// Inner context: Request-level context (context.Context)
//	func handleRequest(userID int) ReaderReaderIOResult[AppConfig, Response] {
//	    return func(cfg AppConfig) readerioresult.ReaderIOResult[context.Context, Response] {
//	        // cfg is available here (outer context)
//	        return func(ctx context.Context) ioresult.IOResult[Response] {
//	            // ctx is available here (inner context)
//	            // Both cfg and ctx can be used
//	            return func() result.Result[Response] {
//	                // Perform operation using both contexts
//	                select {
//	                case <-ctx.Done():
//	                    return result.Error[Response](ctx.Err())
//	                default:
//	                    // Use cfg.DatabaseURL to connect
//	                    return result.Of(Response{})
//	                }
//	            }
//	        }
//	    }
//	}
//
// ### Example 6: Avoiding Global State (Video: "Problems with Global State")
//
// The video criticizes global state and shows how Reader solves this.
// In fp-go, instead of:
//
//	// BAD: Global state
//	var globalConfig AppConfig
//
//	func fetchUser(id int) result.Result[User] {
//	    // Uses globalConfig implicitly
//	    db := connectTo(globalConfig.DatabaseURL)
//	    // ...
//	}
//
// Use Reader to make dependencies explicit:
//
//	// GOOD: Explicit dependencies
//	func fetchUser(id int) ReaderReaderIOResult[AppConfig, User] {
//	    return MonadChain(
//	        Ask[AppConfig](), // Explicitly request the config
//	        func(cfg AppConfig) ReaderReaderIOResult[AppConfig, User] {
//	            // Use cfg explicitly
//	            return FromIO[AppConfig](func() result.Result[User] {
//	                db := connectTo(cfg.DatabaseURL)
//	                // ...
//	            })
//	        },
//	    )
//	}
//
// ## Benefits of This Approach
//
// 1. **Type Safety**: The compiler ensures all dependencies are provided
// 2. **Testability**: Easy to provide mock dependencies for testing
// 3. **Composability**: Functions compose naturally without dependency wiring
// 4. **Explicitness**: Function signatures document their dependencies
// 5. **Immutability**: Dependencies are immutable values, not mutable global state
// 6. **Flexibility**: Use Local to modify dependencies for sub-computations
// 7. **Separation of Concerns**: Business logic is separate from dependency resolution
//
// ## Comparison with Traditional DI
//
// Traditional OOP DI (e.g., Spring, Guice):
//   - Runtime dependency resolution
//   - Magic/reflection-based wiring
//   - Implicit dependencies (hidden in constructors)
//   - Mutable containers
//
// Reader-based DI (fp-go):
//   - Compile-time dependency resolution
//   - Explicit function composition
//   - Explicit dependencies (in type signatures)
//   - Immutable values
//
// ## When to Use Each Layer
//
// - **Outer Reader (R)**: Application-level dependencies that rarely change
//   - Database connection pools
//   - API keys and secrets
//   - Feature flags
//   - Application configuration
//
// - **Inner Reader (context.Context)**: Request-level dependencies that change per operation
//   - Request IDs and tracing
//   - Cancellation signals
//   - Deadlines and timeouts
//   - User authentication tokens
//
// This two-layer approach mirrors the video's discussion of nested readers and provides
// a clean separation between application-level and request-level concerns.
//
// # Relationship to Other Packages
//
//   - readerreaderioeither: The generic version with configurable error and context types
//   - readerioresult: Single reader with context.Context and error
//   - readerresult: Single reader with error (no IO)
//   - context/readerioresult: Alias for readerioresult with context.Context
package readerreaderioresult
