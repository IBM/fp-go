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

package ioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/tailrec"
)

// TailRec creates a tail-recursive computation in the IOEither monad.
// It enables writing recursive algorithms that don't overflow the call stack by using
// trampolining - a technique where recursive calls are converted into iterations.
//
// The function takes a step function that returns a Trampoline:
//   - Bounce(A): Continue recursion with a new value of type A
//   - Land(B): Terminate recursion with a final result of type B
//
// This is particularly useful for implementing recursive algorithms like:
//   - Iterative calculations (factorial, fibonacci, etc.)
//   - State machines with multiple steps
//   - Loops that may fail at any iteration
//   - Processing collections with early termination
//
// The recursion is stack-safe because each step returns a value that indicates
// whether to continue (Bounce) or stop (Land), rather than making direct recursive calls.
//
// Type Parameters:
//   - E: The error type that may occur during computation
//   - A: The intermediate type used during recursion (loop state)
//   - B: The final result type when recursion terminates
//
// Parameters:
//   - f: A step function that takes the current state (A) and returns an IOEither
//     containing either Bounce(A) to continue with a new state, or Land(B) to
//     terminate with a final result
//
// Returns:
//   - A Kleisli arrow (function from A to IOEither[E, B]) that executes the
//     tail-recursive computation starting from the initial value
//
// Example - Computing factorial in a stack-safe way:
//
//	type FactState struct {
//	    n      int
//	    result int
//	}
//
//	factorial := TailRec(func(state FactState) IOEither[error, tailrec.Trampoline[FactState, int]] {
//	    if state.n <= 1 {
//	        // Terminate with final result
//	        return Of[error](tailrec.Land[FactState](state.result))
//	    }
//	    // Continue with next iteration
//	    return Of[error](tailrec.Bounce[int](FactState{
//	        n:      state.n - 1,
//	        result: state.result * state.n,
//	    }))
//	})
//
//	result := factorial(FactState{n: 5, result: 1})() // Right(120)
//
// Example - Processing a list with potential errors:
//
//	type ProcessState struct {
//	    items []string
//	    sum   int
//	}
//
//	processItems := TailRec(func(state ProcessState) IOEither[error, tailrec.Trampoline[ProcessState, int]] {
//	    if len(state.items) == 0 {
//	        return Of[error](tailrec.Land[ProcessState](state.sum))
//	    }
//	    val, err := strconv.Atoi(state.items[0])
//	    if err != nil {
//	        return Left[tailrec.Trampoline[ProcessState, int]](err)
//	    }
//	    return Of[error](tailrec.Bounce[int](ProcessState{
//	        items: state.items[1:],
//	        sum:   state.sum + val,
//	    }))
//	})
//
//	result := processItems(ProcessState{items: []string{"1", "2", "3"}, sum: 0})() // Right(6)
func TailRec[E, A, B any](f Kleisli[E, A, tailrec.Trampoline[A, B]]) Kleisli[E, A, B] {
	return func(a A) IOEither[E, B] {
		initial := f(a)
		return func() either.Either[E, B] {
			current := initial()
			for {
				r, e := either.Unwrap(current)
				if either.IsLeft(current) {
					return either.Left[B](e)
				}
				if r.Landed {
					return either.Right[E](r.Land)
				}
				current = f(r.Bounce)()
			}
		}
	}
}
