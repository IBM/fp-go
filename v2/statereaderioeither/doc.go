// Copyright (c) 2024 - 2025 IBM Corp.
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

// Package statereaderioeither provides a functional programming abstraction that combines
// four powerful concepts: State, Reader, IO, and Either monads.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - State monad: https://github.com/fantasyland/fantasy-land
//   - Reader monad: https://github.com/fantasyland/fantasy-land
//   - IO monad: https://github.com/fantasyland/fantasy-land
//   - Either monad: https://github.com/fantasyland/fantasy-land#either
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//
// # StateReaderIOEither
//
// StateReaderIOEither[S, R, E, A] represents a computation that:
//   - Manages state of type S (State)
//   - Depends on some context/environment of type R (Reader)
//   - Performs side effects (IO)
//   - Can fail with an error of type E or succeed with a value of type A (Either)
//
// The type is defined as:
//
//	StateReaderIOEither[S, R, E, A] = Reader[S, ReaderIOEither[R, E, Pair[S, A]]]
//
// This is particularly useful for:
//   - Stateful computations with dependency injection
//   - Error handling in effectful computations with state
//   - Composing operations that need access to shared configuration, manage state, and can fail
//
// # Core Operations
//
// Construction:
//   - Of/Right: Create a successful computation with a value
//   - Left: Create a failed computation with an error
//   - FromState: Lift a State into StateReaderIOEither
//   - FromReader: Lift a Reader into StateReaderIOEither
//   - FromIO: Lift an IO into StateReaderIOEither
//   - FromEither: Lift an Either into StateReaderIOEither
//   - FromIOEither: Lift an IOEither into StateReaderIOEither
//   - FromReaderEither: Lift a ReaderEither into StateReaderIOEither
//   - FromReaderIOEither: Lift a ReaderIOEither into StateReaderIOEither
//
// Transformation:
//   - Map: Transform the success value
//   - Chain: Sequence dependent computations (monadic bind)
//   - Flatten: Flatten nested StateReaderIOEither
//
// Combination:
//   - Ap: Apply a function in a context to a value in a context
//
// Context Access:
//   - Asks: Get a value derived from the context
//   - Local: Run a computation with a modified context
//
// Kleisli Arrows:
//   - FromEitherK: Lift an Either-returning function to a Kleisli arrow
//   - FromIOK: Lift an IO-returning function to a Kleisli arrow
//   - FromIOEitherK: Lift an IOEither-returning function to a Kleisli arrow
//   - FromReaderIOEitherK: Lift a ReaderIOEither-returning function to a Kleisli arrow
//   - ChainEitherK: Chain with an Either-returning function
//   - ChainIOEitherK: Chain with an IOEither-returning function
//   - ChainReaderIOEitherK: Chain with a ReaderIOEither-returning function
//
// Do Notation (Monadic Composition):
//   - Do: Start a do-notation chain
//   - Bind: Bind a value from a computation
//   - BindTo: Bind a value to a simple constructor
//   - Let: Compute a derived value
//   - LetTo: Set a constant value
//   - ApS: Apply in sequence (for applicative composition)
//   - BindL/ApSL/LetL/LetToL: Lens-based variants for working with nested structures
//
// # Example Usage
//
//	type Config struct {
//	    DatabaseURL string
//	    MaxRetries  int
//	}
//
//	type AppState struct {
//	    RequestCount int
//	    LastError    error
//	}
//
//	// A computation that manages state, depends on config, performs IO, and can fail
//	func processRequest(data string) statereaderioeither.StateReaderIOEither[AppState, Config, error, string] {
//	    return func(state AppState) readerioeither.ReaderIOEither[Config, error, pair.Pair[AppState, string]] {
//	        return func(cfg Config) ioeither.IOEither[error, pair.Pair[AppState, string]] {
//	            return func() either.Either[error, pair.Pair[AppState, string]] {
//	                // Use cfg.DatabaseURL and cfg.MaxRetries
//	                // Update state.RequestCount
//	                // Perform IO operations
//	                // Return either.Right(pair.MakePair(newState, result)) or either.Left(err)
//	                newState := AppState{RequestCount: state.RequestCount + 1}
//	                return either.Right(pair.MakePair(newState, "processed: " + data))
//	            }
//	        }
//	    }
//	}
//
//	// Compose operations using do-notation
//	result := function.Pipe3(
//	    statereaderioeither.Do[AppState, Config, error](State{}),
//	    statereaderioeither.Bind(
//	        func(result string) func(State) State { return func(s State) State { return State{result: result} } },
//	        func(s State) statereaderioeither.StateReaderIOEither[AppState, Config, error, string] {
//	            return processRequest(s.input)
//	        },
//	    ),
//	    statereaderioeither.Map[AppState, Config, error](func(s State) string { return s.result }),
//	)
//
//	// Execute with initial state and config
//	initialState := AppState{RequestCount: 0}
//	config := Config{DatabaseURL: "postgres://localhost", MaxRetries: 3}
//	outcome := result(initialState)(config)() // Returns either.Either[error, pair.Pair[AppState, string]]
//
// # Monad Laws
//
// StateReaderIOEither satisfies the monad laws:
//   - Left Identity: Of(a) >>= f ≡ f(a)
//   - Right Identity: m >>= Of ≡ m
//   - Associativity: (m >>= f) >>= g ≡ m >>= (x => f(x) >>= g)
//
// These laws are verified in the testing subpackage.
package statereaderioeither
