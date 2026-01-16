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

package either

import (
	"github.com/IBM/fp-go/v2/tailrec"
)

// TailRec converts a tail-recursive Kleisli arrow into a stack-safe iterative computation.
//
// This function enables writing recursive algorithms in a functional style while avoiding
// stack overflow errors. It takes a Kleisli arrow that returns a Trampoline wrapped in Either,
// and converts it into a regular Kleisli arrow that executes the recursion iteratively.
//
// The function handles both success and failure cases:
//   - If any step returns Left[E], the recursion stops and returns that error
//   - If a step returns Right with Landed=true, the final result is returned
//   - If a step returns Right with Landed=false, recursion continues with the bounced value
//
// Type Parameters:
//   - E: The error type (Left case)
//   - A: The input type for each recursive step
//   - B: The final result type (Right case)
//
// Parameters:
//   - f: A Kleisli arrow that takes an input of type A and returns Either[E, Trampoline[A, B]]
//     The Trampoline indicates whether to continue (Bounce) or terminate (Land)
//
// Returns:
//   - A Kleisli arrow that executes the tail recursion iteratively and returns Either[E, B]
//
// Example - Factorial with error handling:
//
//	type State struct { n, acc int }
//
//	factorialStep := func(state State) either.Either[string, tailrec.Trampoline[State, int]] {
//	    if state.n < 0 {
//	        return either.Left[tailrec.Trampoline[State, int]]("negative input")
//	    }
//	    if state.n <= 1 {
//	        return either.Right[string](tailrec.Land[State](state.acc))
//	    }
//	    return either.Right[string](tailrec.Bounce[int](State{state.n - 1, state.acc * state.n}))
//	}
//
//	factorial := either.TailRec(factorialStep)
//	result := factorial(State{5, 1}) // Right(120)
//	error := factorial(State{-1, 1}) // Left("negative input")
//
// Example - Countdown with validation:
//
//	countdown := either.TailRec(func(n int) either.Either[string, tailrec.Trampoline[int, int]] {
//	    if n < 0 {
//	        return either.Left[tailrec.Trampoline[int, int]]("already negative")
//	    }
//	    if n == 0 {
//	        return either.Right[string](tailrec.Land[int](0))
//	    }
//	    return either.Right[string](tailrec.Bounce[int](n - 1))
//	})
//
//	result := countdown(5) // Right(0)
//
// The function is stack-safe and can handle arbitrarily deep recursion without
// causing stack overflow, as it uses iteration internally rather than actual recursion.
//
//go:inline
func TailRec[E, A, B any](f Kleisli[E, A, tailrec.Trampoline[A, B]]) Kleisli[E, A, B] {
	return func(a A) Either[E, B] {
		current := f(a)
		for {
			rec, e := Unwrap(current)
			if IsLeft(current) {
				return Left[B](e)
			}
			if rec.Landed {
				return Right[E](rec.Land)
			}
			current = f(rec.Bounce)
		}
	}
}
