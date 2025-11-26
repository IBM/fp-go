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

// Package readerioresult provides a ReaderIOResult monad that combines Reader, IO, and Result monads.
//
// A ReaderIOResult[R, A] represents a computation that:
//   - Depends on an environment of type R (Reader aspect)
//   - Performs IO operations (IO aspect)
//   - May fail with an error (Result aspect, which is Either[error, A])
//
// This is equivalent to Reader[R, IOResult[A]] or Reader[R, func() (A, error)].
//
// # Use Cases
//
// ReaderIOResult is particularly useful for:
//
//  1. Dependency injection with IO and error handling - pass configuration/services through
//     computations that perform side effects and may fail
//  2. Functional IO with context - compose IO operations that depend on environment and may error
//  3. Testing - easily mock dependencies and IO operations by changing the environment value
//  4. Resource management - manage resources that depend on configuration
//
// # Basic Example
//
//	type Config struct {
//	    DatabaseURL string
//	    Timeout     time.Duration
//	}
//
//	// Function that needs config, performs IO, and may fail
//	func fetchUser(id int) readerioresult.ReaderIOResult[Config, User] {
//	    return func(cfg Config) ioresult.IOResult[User] {
//	        return func() (User, error) {
//	            // Use cfg.DatabaseURL and cfg.Timeout to fetch user
//	            return queryDatabase(cfg.DatabaseURL, id, cfg.Timeout)
//	        }
//	    }
//	}
//
//	// Execute by providing the config
//	cfg := Config{DatabaseURL: "postgres://...", Timeout: 5 * time.Second}
//	ioResult := fetchUser(42)(cfg)  // Returns IOResult[User]
//	user, err := ioResult()         // Execute the IO operation
//
// # Composition
//
// ReaderIOResult provides several ways to compose computations:
//
//  1. Map - transform successful values
//  2. Chain (FlatMap) - sequence dependent IO operations
//  3. Ap - combine independent IO computations
//  4. ChainFirst - perform IO for side effects while keeping original value
//
// # Example with Composition
//
//	type AppContext struct {
//	    DB    *sql.DB
//	    Cache Cache
//	    Log   Logger
//	}
//
//	getUserWithCache := F.Pipe2(
//	    getFromCache(userID),
//	    readerioresult.Alt(func() readerioresult.ReaderIOResult[AppContext, User] {
//	        return F.Pipe2(
//	            getFromDB(userID),
//	            readerioresult.ChainFirst(saveToCache),
//	        )
//	    }),
//	)
//
//	ctx := AppContext{DB: db, Cache: cache, Log: logger}
//	user, err := getUserWithCache(ctx)()
//
// # Error Handling
//
// ReaderIOResult provides several functions for error handling:
//
//   - Left/Right - create failed/successful values
//   - GetOrElse - provide a default value for errors
//   - OrElse - recover from errors with an alternative computation
//   - Fold - handle both success and failure cases
//   - ChainLeft - transform error values into new computations
//
// # Relationship to Other Monads
//
// ReaderIOResult is related to several other monads in this library:
//
//   - Reader[R, A] - ReaderIOResult without IO or error handling
//   - IOResult[A] - ReaderIOResult without environment dependency
//   - ReaderResult[R, A] - ReaderIOResult without IO (pure computations)
//   - ReaderIO[R, A] - ReaderIOResult without error handling
//   - ReaderIOEither[R, E, A] - like ReaderIOResult but with custom error type E
//
// # Performance Note
//
// ReaderIOResult is a zero-cost abstraction - it compiles to a simple function type
// with no runtime overhead beyond the underlying computation. The IO operations are
// lazy and only executed when the final IOResult is called.
package readerioresult
