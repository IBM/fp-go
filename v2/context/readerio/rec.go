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

package readerio

import (
	"github.com/IBM/fp-go/v2/readerio"
)

// TailRec implements stack-safe tail recursion for the ReaderIO monad.
//
// This function enables recursive computations that depend on a [context.Context] and
// perform side effects, without risking stack overflow. It uses an iterative loop to
// execute the recursion, making it safe for deep or unbounded recursion.
//
// The function takes a Kleisli arrow that returns Trampoline[A, B]:
//   - Bounce(A): Continue recursion with the new state A
//   - Land(B): Terminate recursion and return the final result B
//
// Type Parameters:
//   - A: The state type that changes during recursion
//   - B: The final result type when recursion terminates
//
// Parameters:
//   - f: A Kleisli arrow (A => ReaderIO[Trampoline[A, B]]) that controls recursion flow
//
// Returns:
//   - A Kleisli arrow (A => ReaderIO[B]) that executes the recursion safely
//
// Example - Countdown:
//
//	countdownStep := func(n int) ReaderIO[tailrec.Trampoline[int, string]] {
//	    return func(ctx context.Context) IO[tailrec.Trampoline[int, string]] {
//	        return func() tailrec.Trampoline[int, string] {
//	            if n <= 0 {
//	                return tailrec.Land[int]("Done!")
//	            }
//	            return tailrec.Bounce[string](n - 1)
//	        }
//	    }
//	}
//
//	countdown := TailRec(countdownStep)
//	result := countdown(10)(t.Context())() // Returns "Done!"
//
// Example - Sum with context:
//
//	type SumState struct {
//	    numbers []int
//	    total   int
//	}
//
//	sumStep := func(state SumState) ReaderIO[tailrec.Trampoline[SumState, int]] {
//	    return func(ctx context.Context) IO[tailrec.Trampoline[SumState, int]] {
//	        return func() tailrec.Trampoline[SumState, int] {
//	            if len(state.numbers) == 0 {
//	                return tailrec.Land[SumState](state.total)
//	            }
//	            return tailrec.Bounce[int](SumState{
//	                numbers: state.numbers[1:],
//	                total:   state.total + state.numbers[0],
//	            })
//	        }
//	    }
//	}
//
//	sum := TailRec(sumStep)
//	result := sum(SumState{numbers: []int{1, 2, 3, 4, 5}})(t.Context())()
//	// Returns 15, safe even for very large slices
//
//go:inline
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B] {
	return readerio.TailRec(f)
}
