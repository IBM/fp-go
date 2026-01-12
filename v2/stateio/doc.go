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

// Package stateio provides a functional programming abstraction that combines
// stateful computations with side effects.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - State monad: https://github.com/fantasyland/fantasy-land
//   - IO monad: https://github.com/fantasyland/fantasy-land
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//
// # StateIO
//
// StateIO[S, A] represents a computation that:
//   - Manages state of type S (State monad)
//   - Performs side effects (IO monad)
//   - Produces a value of type A
//
// The type is defined as:
//
//	StateIO[S, A] = Reader[S, IO[Pair[S, A]]]
//
// This is particularly useful for:
//   - Stateful computations with side effects
//   - Managing application state while performing IO operations
//   - Composing operations that need both state management and effectful computation
//
// # Core Operations
//
// Construction:
//   - Of: Create a computation with a pure value
//   - FromIO: Lift an IO computation into StateIO
//
// Transformation:
//   - Map: Transform the value within the computation
//   - Chain: Sequence dependent computations (monadic bind)
//
// Combination:
//   - Ap: Apply a function in a context to a value in a context
//
// Kleisli Arrows:
//   - FromIOK: Lift an IO-returning function to a Kleisli arrow
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
//	type AppState struct {
//	    RequestCount int
//	    LastError    error
//	}
//
//	// A computation that manages state and performs IO
//	func incrementCounter(data string) StateIO[AppState, string] {
//	    return func(state AppState) IO[Pair[AppState, string]] {
//	        return func() Pair[AppState, string] {
//	            // Update state.RequestCount
//	            // Perform IO operations
//	            newState := AppState{RequestCount: state.RequestCount + 1}
//	            result := "processed: " + data
//	            return pair.MakePair(newState, result)
//	        }
//	    }
//	}
//
//	// Compose operations using do-notation
//	type Result struct {
//	    result string
//	    count  int
//	}
//
//	computation := function.Pipe3(
//	    Do[AppState](Result{}),
//	    Bind(
//	        func(result string) func(Result) Result {
//	            return func(r Result) Result { return Result{result: result, count: r.count} }
//	        },
//	        func(r Result) StateIO[AppState, string] {
//	            return incrementCounter("data")
//	        },
//	    ),
//	    Map[AppState](func(r Result) string { return r.result }),
//	)
//
//	// Execute with initial state
//	initialState := AppState{RequestCount: 0}
//	outcome := computation(initialState)() // Returns Pair[AppState, string]
//
// # Monad Laws
//
// StateIO satisfies the monad laws:
//   - Left Identity: Of(a) >>= f ≡ f(a)
//   - Right Identity: m >>= Of ≡ m
//   - Associativity: (m >>= f) >>= g ≡ m >>= (x => f(x) >>= g)
//
// Where >>= represents the Chain operation (monadic bind).
//
// These laws ensure that StateIO computations compose predictably and that
// the order of composition doesn't affect the final result (beyond the order
// of effects and state updates).
package stateio
