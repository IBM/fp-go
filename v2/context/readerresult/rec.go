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

// Package readerresult implements a specialization of the Reader monad assuming a golang context as the context of the monad and a standard golang error
package readerresult

import (
	"context"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/result"
)

// TailRec implements tail-recursive computation for ReaderResult with context cancellation support.
//
// TailRec takes a Kleisli function that returns Trampoline[A, B] and converts it into a stack-safe,
// tail-recursive computation. The function repeatedly applies the Kleisli until it produces a Land value.
//
// The implementation includes a short-circuit mechanism that checks for context cancellation on each
// iteration. If the context is canceled (ctx.Err() != nil), the computation immediately returns a
// Left result containing the context's cause error, preventing unnecessary computation.
//
// Type Parameters:
//   - A: The input type for the recursive step
//   - B: The final result type
//
// Parameters:
//   - f: A Kleisli function that takes an A and returns a ReaderResult containing Trampoline[A, B].
//     When the result is Bounce(a), recursion continues with the new value 'a'.
//     When the result is Land(b), recursion terminates with the final value 'b'.
//
// Returns:
//   - A Kleisli function that performs the tail-recursive computation in a stack-safe manner.
//
// Behavior:
//   - On each iteration, checks if the context has been canceled (short circuit)
//   - If canceled, returns result.Left[B](context.Cause(ctx))
//   - If the step returns Left[B](error), propagates the error
//   - If the step returns Right[A](Bounce(a)), continues recursion with new value 'a'
//   - If the step returns Right[A](Land(b)), terminates with success value 'b'
//
// Example - Factorial computation with context:
//
//	type State struct {
//	    n   int
//	    acc int
//	}
//
//	factorialStep := func(state State) ReaderResult[tailrec.Trampoline[State, int]] {
//	    return func(ctx context.Context) result.Result[tailrec.Trampoline[State, int]] {
//	        if state.n <= 0 {
//	            return result.Of(tailrec.Land[State](state.acc))
//	        }
//	        return result.Of(tailrec.Bounce[int](State{state.n - 1, state.acc * state.n}))
//	    }
//	}
//
//	factorial := TailRec(factorialStep)
//	result := factorial(State{5, 1})(ctx) // Returns result.Of(120)
//
// Example - Context cancellation:
//
//	ctx, cancel := context.WithCancel(t.Context())
//	cancel() // Cancel immediately
//
//	computation := TailRec(someStep)
//	result := computation(initialValue)(ctx)
//	// Returns result.Left[B](context.Cause(ctx)) without executing any steps
//
//go:inline
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B] {
	return func(a A) ReaderResult[B] {
		initialReader := f(a)
		return func(ctx context.Context) result.Result[B] {
			rdr := initialReader
			for {
				// short circuit
				if ctx.Err() != nil {
					return result.Left[B](context.Cause(ctx))
				}
				current := rdr(ctx)
				rec, e := either.Unwrap(current)
				if either.IsLeft(current) {
					return result.Left[B](e)
				}
				if rec.Landed {
					return result.Of(rec.Land)
				}
				rdr = f(rec.Bounce)
			}
		}
	}
}
