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

// Package readerresult implements a specialization of the Reader monad assuming a golang context as the context of the monad and a standard golang error
package readerresult

import (
	"context"

	"github.com/IBM/fp-go/v2/idiomatic/result"
)

// TailRec implements tail-recursive computation for ReaderResult with context cancellation support.
//
// TailRec takes a Kleisli function that returns Trampoline[A, B] and converts it into a stack-safe,
// tail-recursive computation. The function repeatedly applies the Kleisli until it produces a Land value.
//
// The implementation includes a short-circuit mechanism that checks for context cancellation on each
// iteration. If the context is canceled (ctx.Err() != nil), the computation immediately returns an
// error result containing the context's cause error, preventing unnecessary computation.
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
//   - If canceled, returns (zero value, context.Cause(ctx))
//   - If the step returns an error, propagates the error as (zero value, error)
//   - If the step returns Bounce(a), continues recursion with new value 'a'
//   - If the step returns Land(b), terminates with success value (b, nil)
//
// Example - Factorial computation with context:
//
//	type State struct {
//	    n   int
//	    acc int
//	}
//
//	factorialStep := func(state State) ReaderResult[tailrec.Trampoline[State, int]] {
//	    return func(ctx context.Context) (tailrec.Trampoline[State, int], error) {
//	        if state.n <= 0 {
//	            return tailrec.Land[State](state.acc), nil
//	        }
//	        return tailrec.Bounce[int](State{state.n - 1, state.acc * state.n}), nil
//	    }
//	}
//
//	factorial := TailRec(factorialStep)
//	result, err := factorial(State{5, 1})(ctx) // Returns (120, nil)
//
// Example - Context cancellation:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	cancel() // Cancel immediately
//
//	computation := TailRec(someStep)
//	result, err := computation(initialValue)(ctx)
//	// Returns (zero value, context.Cause(ctx)) without executing any steps
//
// Example - Error handling:
//
//	errorStep := func(n int) ReaderResult[tailrec.Trampoline[int, int]] {
//	    return func(ctx context.Context) (tailrec.Trampoline[int, int], error) {
//	        if n == 5 {
//	            return tailrec.Trampoline[int, int]{}, errors.New("computation error")
//	        }
//	        if n <= 0 {
//	            return tailrec.Land[int](n), nil
//	        }
//	        return tailrec.Bounce[int](n - 1), nil
//	    }
//	}
//
//	computation := TailRec(errorStep)
//	result, err := computation(10)(ctx) // Returns (0, errors.New("computation error"))
//
//go:inline
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B] {
	return func(a A) ReaderResult[B] {
		initialReader := f(a)
		return func(ctx context.Context) (B, error) {
			rdr := initialReader
			for {
				// short circuit
				if ctx.Err() != nil {
					return result.Left[B](context.Cause(ctx))
				}
				rec, e := rdr(ctx)
				if e != nil {
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
