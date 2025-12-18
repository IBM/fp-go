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

package readerioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/tailrec"
)

// TailRec implements stack-safe tail recursion for the ReaderIOEither monad.
//
// This function enables recursive computations that combine three powerful concepts:
//   - Environment dependency (Reader aspect): Access to configuration, context, or dependencies
//   - Side effects (IO aspect): Logging, file I/O, network calls, etc.
//   - Error handling (Either aspect): Computations that can fail with an error
//
// The function uses an iterative loop to execute the recursion, making it safe for deep
// or unbounded recursion without risking stack overflow.
//
// # How It Works
//
// TailRec takes a Kleisli arrow that returns IOEither[E, Trampoline[A, B]]:
//   - Left(E): Computation failed with error E - recursion terminates
//   - Right(Bounce(A)): Continue recursion with the new state A
//   - Right(Land(B)): Terminate recursion successfully and return the final result B
//
// The function iteratively applies the Kleisli arrow, passing the environment R to each
// iteration, until either an error (Left) or a final result (Right(Right(B))) is produced.
//
// # Type Parameters
//
//   - R: The environment type (Reader context) - e.g., Config, Logger, Database connection
//   - E: The error type that can occur during computation
//   - A: The state type that changes during recursion
//   - B: The final result type when recursion terminates successfully
//
// # Parameters
//
//   - f: A Kleisli arrow (A => ReaderIOEither[R, E, Either[A, B]]) that:
//   - Takes the current state A
//   - Returns a ReaderIOEither that depends on environment R
//   - Can fail with error E (Left)
//   - Produces Either[A, B] to control recursion flow (Right)
//
// # Returns
//
// A Kleisli arrow (A => ReaderIOEither[R, E, B]) that:
//   - Takes an initial state A
//   - Returns a ReaderIOEither that requires environment R
//   - Can fail with error E
//   - Produces the final result B after recursion completes
//
// # Comparison with Other Monads
//
// Compared to other tail recursion implementations:
//   - Like IOEither: Has error handling (Left for errors)
//   - Like ReaderIO: Has environment dependency (R passed to each iteration)
//   - Unlike IOOption: Uses Either for both errors and recursion control
//   - Most powerful: Combines all three aspects (Reader + IO + Either)
//
// # Use Cases
//
//  1. Environment-dependent recursive algorithms with error handling:
//     - Recursive computations that need configuration and can fail
//     - Algorithms that log progress and handle errors gracefully
//     - Recursive operations with retry logic based on environment settings
//
//  2. Stateful computations with context and error handling:
//     - Tree traversals that need environment context and can fail
//     - Graph algorithms with configuration-dependent behavior and error cases
//     - Recursive parsers with environment-based rules and error reporting
//
//  3. Recursive operations with side effects and error handling:
//     - File system traversals with logging and error handling
//     - Network operations with retry configuration and failure handling
//     - Database operations with connection pooling and error recovery
//
// # Example: Factorial with Logging and Error Handling
//
//	type Env struct {
//	    Logger func(string)
//	    MaxN   int
//	}
//
//	type State struct {
//	    n   int
//	    acc int
//	}
//
//	// Factorial that logs each step and validates input
//	factorialStep := func(state State) readerioeither.ReaderIOEither[Env, string, tailrec.Trampoline[State, int]] {
//	    return func(env Env) ioeither.IOEither[string, tailrec.Trampoline[State, int]] {
//	        return func() either.Either[string, tailrec.Trampoline[State, int]] {
//	            if state.n > env.MaxN {
//	                return either.Left[tailrec.Trampoline[State, int]](fmt.Sprintf("n too large: %d > %d", state.n, env.MaxN))
//	            }
//	            if state.n <= 0 {
//	                env.Logger(fmt.Sprintf("Factorial complete: %d", state.acc))
//	                return either.Right[string](tailrec.Land[State](state.acc))
//	            }
//	            env.Logger(fmt.Sprintf("Computing: %d * %d", state.n, state.acc))
//	            return either.Right[string](tailrec.Bounce[int](State{state.n - 1, state.acc * state.n}))
//	        }
//	    }
//	}
//
//	factorial := readerioeither.TailRec(factorialStep)
//	env := Env{Logger: func(msg string) { fmt.Println(msg) }, MaxN: 20}
//	result := factorial(State{5, 1})(env)() // Returns Right(120), logs each step
//	// If n > MaxN, returns Left("n too large: ...")
//
// # Example: File Processing with Error Handling
//
//	type Config struct {
//	    MaxRetries int
//	    Logger     func(string)
//	}
//
//	type ProcessState struct {
//	    files   []string
//	    results []string
//	    retries int
//	}
//
//	processFilesStep := func(state ProcessState) readerioeither.ReaderIOEither[Config, error, tailrec.Trampoline[ProcessState, []string]] {
//	    return func(cfg Config) ioeither.IOEither[error, tailrec.Trampoline[ProcessState, []string]] {
//	        return func() either.Either[error, tailrec.Trampoline[ProcessState, []string]] {
//	            if len(state.files) == 0 {
//	                cfg.Logger("All files processed")
//	                return either.Right[error](tailrec.Land[ProcessState](state.results))
//	            }
//	            file := state.files[0]
//	            cfg.Logger(fmt.Sprintf("Processing: %s", file))
//
//	            // Simulate file processing that might fail
//	            if err := processFile(file); err != nil {
//	                if state.retries >= cfg.MaxRetries {
//	                    return either.Left[tailrec.Trampoline[ProcessState, []string]](
//	                        fmt.Errorf("max retries exceeded for %s: %w", file, err))
//	                }
//	                cfg.Logger(fmt.Sprintf("Retry %d for %s", state.retries+1, file))
//	                return either.Right[error](tailrec.Bounce[[]string](ProcessState{
//	                    files:   state.files,
//	                    results: state.results,
//	                    retries: state.retries + 1,
//	                }))
//	            }
//
//	            return either.Right[error](tailrec.Bounce[[]string](ProcessState{
//	                files:   state.files[1:],
//	                results: append(state.results, file),
//	                retries: 0,
//	            }))
//	        }
//	    }
//	}
//
//	processFiles := readerioeither.TailRec(processFilesStep)
//	config := Config{MaxRetries: 3, Logger: func(msg string) { fmt.Println(msg) }}
//	result := processFiles(ProcessState{files: []string{"a.txt", "b.txt"}, results: []string{}, retries: 0})(config)()
//	// Returns Right([]string{"a.txt", "b.txt"}) on success
//	// Returns Left(error) if max retries exceeded
//
// # Stack Safety
//
// The iterative implementation ensures that even deeply recursive computations
// (thousands or millions of iterations) will not cause stack overflow:
//
//	// Safe for very large inputs
//	countdownStep := func(n int) readerioeither.ReaderIOEither[Env, error, tailrec.Trampoline[int, int]] {
//	    return func(env Env) ioeither.IOEither[error, tailrec.Trampoline[int, int]] {
//	        return func() either.Either[error, tailrec.Trampoline[int, int]] {
//	            if n <= 0 {
//	                return either.Right[error](tailrec.Land[int](0))
//	            }
//	            return either.Right[error](tailrec.Bounce[int](n - 1))
//	        }
//	    }
//	}
//	countdown := readerioeither.TailRec(countdownStep)
//	result := countdown(1000000)(env)() // Safe, no stack overflow
//
// # Error Handling Patterns
//
// The Either[E, Trampoline[A, B]] structure provides two levels of control:
//
//  1. Outer Either (Left(E)): Unrecoverable errors that terminate recursion
//     - Validation failures
//     - Resource exhaustion
//     - Fatal errors
//
//  2. Inner Trampoline (Right(Bounce(A)) or Right(Land(B))): Recursion control
//     - Bounce(A): Continue with new state
//     - Land(B): Terminate successfully
//
// This separation allows for:
//   - Early termination on errors
//   - Graceful error propagation
//   - Clear distinction between "continue" and "error" states
//
// # Performance Considerations
//
//   - Each iteration creates a new IOEither action by calling f(a)(r)()
//   - The environment R is passed to every iteration
//   - Error checking happens on each iteration (Left vs Right)
//   - For performance-critical code, consider if error handling is necessary at each step
//   - Memoization of environment-derived values may improve performance
//
// # See Also
//
//   - [ioeither.TailRec]: Tail recursion with error handling (no environment)
//   - [readerio.TailRec]: Tail recursion with environment (no error handling)
//   - [iooption.TailRec]: Tail recursion with optional results
//   - [Chain]: For sequencing ReaderIOEither computations
//   - [Ask]: For accessing the environment
//   - [Left]/[Right]: For creating error/success values
func TailRec[R, E, A, B any](f Kleisli[R, E, A, tailrec.Trampoline[A, B]]) Kleisli[R, E, A, B] {
	return func(a A) ReaderIOEither[R, E, B] {
		initialReader := f(a)
		return func(r R) IOEither[E, B] {
			initialB := initialReader(r)
			return func() either.Either[E, B] {
				current := initialB()
				for {
					rec, e := either.Unwrap(current)
					if either.IsLeft(current) {
						return either.Left[B](e)
					}
					if rec.Landed {
						return either.Right[E](rec.Land)
					}
					current = f(rec.Bounce)(r)()
				}
			}
		}
	}
}
