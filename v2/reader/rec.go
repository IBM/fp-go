// Copyright (c) 2025 IBM Corp.
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

package reader

// TailRec converts a tail-recursive Kleisli arrow into a stack-safe iterative computation
// using the trampoline pattern.
//
// This function enables writing recursive algorithms in a functional style while avoiding
// stack overflow errors. It takes a Kleisli arrow that returns a Trampoline and converts
// it into a regular Kleisli arrow that executes the computation iteratively.
//
// The trampoline pattern works by:
//  1. Starting with an initial state A
//  2. Applying function f to get a Trampoline[A, B]
//  3. If the trampoline has "landed" (Landed == true), return the final result
//  4. If the trampoline should "bounce" (Landed == false), continue with the new state
//  5. Repeat steps 2-4 until landing
//
// Type Parameters:
//   - R: The environment/context type shared across all Reader computations
//   - A: The intermediate state type (bounce type) passed between recursive steps
//   - B: The final result type (land type) when the computation completes
//
// Parameters:
//   - f: A Kleisli arrow that takes a state A and returns a Reader producing a Trampoline[A, B].
//     This function represents one step of the recursive computation.
//
// Returns:
//   - A Kleisli arrow that takes an initial state A and returns a Reader producing the final result B.
//     The returned computation is stack-safe and executes iteratively.
//
// Example - Factorial with Reader environment:
//
//	type State struct {
//	    n   int
//	    acc int
//	}
//
//	type Config struct {
//	    maxIterations int
//	}
//
//	// One step of factorial computation
//	factorialStep := func(state State) reader.Reader[Config, tailrec.Trampoline[State, int]] {
//	    return func(cfg Config) tailrec.Trampoline[State, int] {
//	        if state.n <= 1 {
//	            return tailrec.Land[State](state.acc)  // Base case
//	        }
//	        if state.n > cfg.maxIterations {
//	            return tailrec.Land[State](-1)  // Error case using config
//	        }
//	        // Recursive case
//	        return tailrec.Bounce[int](State{n: state.n - 1, acc: state.acc * state.n})
//	    }
//	}
//
//	// Create stack-safe factorial
//	factorial := reader.TailRec(factorialStep)
//
//	// Execute with environment
//	config := Config{maxIterations: 1000}
//	result := factorial(State{n: 5, acc: 1})(config)  // Returns 120
//
// Example - Countdown:
//
//	type Env struct{ verbose bool }
//
//	countdown := func(n int) reader.Reader[Env, tailrec.Trampoline[int, string]] {
//	    return func(env Env) tailrec.Trampoline[int, string] {
//	        if n <= 0 {
//	            return tailrec.Land[int]("Done!")
//	        }
//	        if env.verbose {
//	            fmt.Printf("Counting: %d\n", n)
//	        }
//	        return tailrec.Bounce[string](n - 1)
//	    }
//	}
//
//	safeCountdown := reader.TailRec(countdown)
//	result := safeCountdown(10)(Env{verbose: true})  // "Done!"
//
// The key benefit is that even with very large recursion depths (e.g., factorial(10000)),
// the computation will not overflow the stack because it executes iteratively.
func TailRec[R, A, B any](f Kleisli[R, A, Trampoline[A, B]]) Kleisli[R, A, B] {
	return func(a A) Reader[R, B] {
		initialReader := f(a)
		return func(r R) B {
			current := initialReader(r)
			for {
				if current.Landed {
					return current.Land
				}
				current = f(current.Bounce)(r)
			}
		}
	}
}
