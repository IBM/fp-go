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

// package readerioresult provides a functional programming abstraction that combines
// three powerful concepts: Reader, IO, and Either monads.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - Reader monad: https://github.com/fantasyland/fantasy-land
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
// # ReaderIOResult
//
// ReaderIOResult[R, A] represents a computation that:
//   - Depends on some context/environment of type R (Reader)
//   - Performs side effects (IO)
//   - Can fail with an error of type E or succeed with a value of type A (Either)
//
// This is particularly useful for:
//   - Dependency injection patterns
//   - Error handling in effectful computations
//   - Composing operations that need access to shared configuration or context
//
// # Core Operations
//
// Construction:
//   - Of/Right: Create a successful computation
//   - Left/ThrowError: Create a failed computation
//   - FromEither: Lift an Either into ReaderIOResult
//   - FromIO: Lift an IO into ReaderIOResult
//   - FromReader: Lift a Reader into ReaderIOResult
//   - FromIOEither: Lift an IOEither into ReaderIOResult
//   - TryCatch: Wrap error-returning functions
//
// Transformation:
//   - Map: Transform the success value
//   - MapLeft: Transform the error value
//   - BiMap: Transform both success and error values
//   - Chain/Bind: Sequence dependent computations
//   - Flatten: Flatten nested ReaderIOResult
//
// Combination:
//   - Ap: Apply a function in a context to a value in a context
//   - SequenceArray: Convert array of ReaderIOResult to ReaderIOResult of array
//   - TraverseArray: Map and sequence in one operation
//
// Error Handling:
//   - Fold: Handle both success and error cases
//   - GetOrElse: Provide a default value on error
//   - OrElse: Try an alternative computation on error
//   - Alt: Choose the first successful computation
//
// Context Access:
//   - Ask: Get the current context
//   - Asks: Get a value derived from the context
//   - Local: Run a computation with a modified context
//
// Resource Management:
//   - Bracket: Ensure resource cleanup
//   - WithResource: Manage resource lifecycle
//
// # Example Usage
//
//	type Config struct {
//	    BaseURL string
//	    Timeout time.Duration
//	}
//
//	// A computation that depends on Config, performs IO, and can fail
//	func fetchUser(id int) readerioeither.ReaderIOResult[Config, error, User] {
//	    return func(cfg Config) ioeither.IOEither[error, User] {
//	        return func() either.Either[error, User] {
//	            // Use cfg.BaseURL and cfg.Timeout to fetch user
//	            // Return either.Right(user) or either.Left(err)
//	        }
//	    }
//	}
//
//	// Compose operations
//	result := function.Pipe2(
//	    fetchUser(123),
//	    readerioeither.Map[Config, error](func(u User) string { return u.Name }),
//	    readerioeither.Chain[Config, error](func(name string) readerioeither.ReaderIOResult[Config, error, string] {
//	        return readerioeither.Of[Config, error]("Hello, " + name)
//	    }),
//	)
//
//	// Execute with config
//	config := Config{BaseURL: "https://api.example.com", Timeout: 30 * time.Second}
//	outcome := result(config)() // Returns either.Either[error, string]
package readerioresult

//go:generate go run .. readerioeither --count 10 --filename gen.go
