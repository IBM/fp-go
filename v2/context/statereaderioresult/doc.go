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

// Package statereaderioresult provides a functional programming abstraction that combines
// four powerful concepts: State, Reader, IO, and Result monads, specialized for Go's context.Context.
//
// # StateReaderIOResult
//
// StateReaderIOResult[S, A] represents a computation that:
//   - Manages state of type S (State)
//   - Depends on a [context.Context] (Reader)
//   - Performs side effects (IO)
//   - Can fail with an [error] or succeed with a value of type A (Result)
//
// This is a specialization of StateReaderIOEither with:
//   - Context type fixed to [context.Context]
//   - Error type fixed to [error]
//
// This is particularly useful for:
//   - Stateful computations with dependency injection using Go contexts
//   - Error handling in effectful computations with state
//   - Composing operations that need access to context, manage state, and can fail
//   - Working with Go's standard context patterns (cancellation, deadlines, values)
//
// # Core Operations
//
// Construction:
//   - Of/Right: Create a successful computation with a value
//   - Left: Create a failed computation with an error
//   - FromState: Lift a State into StateReaderIOResult
//   - FromIO: Lift an IO into StateReaderIOResult
//   - FromResult: Lift a Result into StateReaderIOResult
//   - FromIOResult: Lift an IOResult into StateReaderIOResult
//   - FromReaderIOResult: Lift a ReaderIOResult into StateReaderIOResult
//
// Transformation:
//   - Map: Transform the success value
//   - Chain: Sequence dependent computations (monadic bind)
//   - Flatten: Flatten nested StateReaderIOResult
//
// Combination:
//   - Ap: Apply a function in a context to a value in a context
//
// Context Access:
//   - Asks: Get a value derived from the context
//   - Local: Run a computation with a modified context
//
// Kleisli Arrows:
//   - FromResultK: Lift a Result-returning function to a Kleisli arrow
//   - FromIOK: Lift an IO-returning function to a Kleisli arrow
//   - FromIOResultK: Lift an IOResult-returning function to a Kleisli arrow
//   - FromReaderIOResultK: Lift a ReaderIOResult-returning function to a Kleisli arrow
//   - ChainResultK: Chain with a Result-returning function
//   - ChainIOResultK: Chain with an IOResult-returning function
//   - ChainReaderIOResultK: Chain with a ReaderIOResult-returning function
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
//	// A computation that manages state, depends on context, performs IO, and can fail
//	func processRequest(data string) statereaderioresult.StateReaderIOResult[AppState, string] {
//	    return func(state AppState) readerioresult.ReaderIOResult[pair.Pair[AppState, string]] {
//	        return func(ctx context.Context) ioresult.IOResult[pair.Pair[AppState, string]] {
//	            return func() result.Result[pair.Pair[AppState, string]] {
//	                // Check context for cancellation
//	                if ctx.Err() != nil {
//	                    return result.Error[pair.Pair[AppState, string]](ctx.Err())
//	                }
//	                // Update state
//	                newState := AppState{RequestCount: state.RequestCount + 1}
//	                // Perform IO operations
//	                return result.Of(pair.MakePair(newState, "processed: " + data))
//	            }
//	        }
//	    }
//	}
//
//	// Compose operations using do-notation
//	result := function.Pipe3(
//	    statereaderioresult.Do[AppState](State{}),
//	    statereaderioresult.Bind(
//	        func(result string) func(State) State { return func(s State) State { return State{result: result} } },
//	        func(s State) statereaderioresult.StateReaderIOResult[AppState, string] {
//	            return processRequest(s.input)
//	        },
//	    ),
//	    statereaderioresult.Map[AppState](func(s State) string { return s.result }),
//	)
//
//	// Execute with initial state and context
//	initialState := AppState{RequestCount: 0}
//	ctx := t.Context()
//	outcome := result(initialState)(ctx)() // Returns result.Result[pair.Pair[AppState, string]]
//
// # Context Integration
//
// This package is designed to work seamlessly with Go's context.Context:
//
//	// Using context values
//	getUserID := statereaderioresult.Asks[AppState, string](func(ctx context.Context) statereaderioresult.StateReaderIOResult[AppState, string] {
//	    userID, ok := ctx.Value("userID").(string)
//	    if !ok {
//	        return statereaderioresult.Left[AppState, string](errors.New("missing userID"))
//	    }
//	    return statereaderioresult.Of[AppState](userID)
//	})
//
//	// Using context cancellation
//	withTimeout := statereaderioresult.Local[AppState, string](func(ctx context.Context) context.Context {
//	    ctx, _ = context.WithTimeout(ctx, 5*time.Second)
//	    return ctx
//	})
//
// # Monad Laws
//
// StateReaderIOResult satisfies the monad laws:
//   - Left Identity: Of(a) >>= f ≡ f(a)
//   - Right Identity: m >>= Of ≡ m
//   - Associativity: (m >>= f) >>= g ≡ m >>= (x => f(x) >>= g)
//
// These laws are verified in the testing subpackage.
package statereaderioresult
