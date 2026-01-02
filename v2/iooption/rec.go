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

package iooption

import (
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/tailrec"
)

// TailRec creates a tail-recursive computation in the IOOption monad.
// It enables writing recursive algorithms that don't overflow the call stack by using
// an iterative loop - a technique where recursive calls are converted into iterations.
//
// The function takes a step function that returns an IOOption containing a Trampoline:
//   - None: Terminate recursion with no result
//   - Some(Bounce(A)): Continue recursion with a new value of type A
//   - Some(Land(B)): Terminate recursion with a final result of type B
//
// This is particularly useful for implementing recursive algorithms that may fail at any step:
//   - Iterative calculations that may not produce a result
//   - State machines with multiple steps that can fail
//   - Loops that may terminate early with None
//   - Processing collections with optional results
//
// Unlike the IOEither version which uses lazy recursion, this implementation uses
// an explicit iterative loop for better performance and simpler control flow.
//
// Type Parameters:
//   - E: Unused type parameter (kept for consistency with IOEither)
//   - A: The intermediate type used during recursion (loop state)
//   - B: The final result type when recursion terminates successfully
//
// Parameters:
//   - f: A step function that takes the current state (A) and returns an IOOption
//     containing either None (failure), Some(Bounce(A)) to continue with a new state,
//     or Some(Land(B)) to terminate with a final result
//
// Returns:
//   - A Kleisli arrow (function from A to IOOption[B]) that executes the
//     tail-recursive computation starting from the initial value
//
// Example - Computing factorial with optional result:
//
//	type FactState struct {
//	    n      int
//	    result int
//	}
//
//	factorial := TailRec[any](func(state FactState) IOOption[tailrec.Trampoline[FactState, int]] {
//	    if state.n < 0 {
//	        // Negative numbers have no factorial
//	        return None[tailrec.Trampoline[FactState, int]]()
//	    }
//	    if state.n <= 1 {
//	        // Terminate with final result
//	        return Of(tailrec.Land[FactState](state.result))
//	    }
//	    // Continue with next iteration
//	    return Of(tailrec.Bounce[int](FactState{
//	        n:      state.n - 1,
//	        result: state.result * state.n,
//	    }))
//	})
//
//	result := factorial(FactState{n: 5, result: 1})() // Some(120)
//	result := factorial(FactState{n: -1, result: 1})() // None
//
// Example - Safe division with early termination:
//
//	type DivState struct {
//	    numerator   int
//	    denominator int
//	    steps       int
//	}
//
//	safeDivide := TailRec[any](func(state DivState) IOOption[tailrec.Trampoline[DivState, int]] {
//	    if state.denominator == 0 {
//	        return None[tailrec.Trampoline[DivState, int]]() // Division by zero
//	    }
//	    if state.numerator < state.denominator {
//	        return Of(tailrec.Land[DivState](state.steps))
//	    }
//	    return Of(tailrec.Bounce[int](DivState{
//	        numerator:   state.numerator - state.denominator,
//	        denominator: state.denominator,
//	        steps:       state.steps + 1,
//	    }))
//	})
//
//	result := safeDivide(DivState{numerator: 10, denominator: 3, steps: 0})() // Some(3)
//	result := safeDivide(DivState{numerator: 10, denominator: 0, steps: 0})() // None
func TailRec[A, B any](f Kleisli[A, tailrec.Trampoline[A, B]]) Kleisli[A, B] {
	return func(a A) IOOption[B] {
		initial := f(a)
		return func() option.Option[B] {
			current := initial()
			for {
				r, ok := option.Unwrap(current)
				if !ok {
					return option.None[B]()
				}
				if r.Landed {
					return option.Some(r.Land)
				}
				current = f(r.Bounce)()
			}
		}
	}
}
