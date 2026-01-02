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

package io

// TailRec creates a tail-recursive computation in the IO monad.
// It enables writing recursive algorithms that don't overflow the call stack by using
// trampolining - a technique where recursive calls are converted into iterations.
//
// The function takes a step function that returns a Trampoline:
//   - Bounce(A): Continue recursion with a new value of type A
//   - Land(B): Terminate recursion with a final result of type B
//
// This is particularly useful for implementing recursive algorithms like:
//   - Iterative calculations (factorial, fibonacci, sum, etc.)
//   - State machines with multiple steps
//   - Loops over large data structures
//   - Processing collections with complex iteration logic
//
// The recursion is stack-safe because each step returns a value that indicates
// whether to continue (Bounce) or stop (Land), rather than making direct recursive calls.
// This allows processing arbitrarily large inputs without stack overflow.
//
// Type Parameters:
//   - A: The intermediate type used during recursion (loop state)
//   - B: The final result type when recursion terminates
//
// Parameters:
//   - f: A step function that takes the current state (A) and returns an IO
//     containing either Bounce(A) to continue with a new state, or Land(B) to
//     terminate with a final result
//
// Returns:
//   - A Kleisli arrow (function from A to IO[B]) that executes the
//     tail-recursive computation starting from the initial value
//
// Example - Computing factorial in a stack-safe way:
//
//	type FactState struct {
//	    n      int
//	    result int
//	}
//
//	factorial := io.TailRec(func(state FactState) io.IO[tailrec.Trampoline[FactState, int]] {
//	    if state.n <= 1 {
//	        // Terminate with final result
//	        return io.Of(tailrec.Land[FactState](state.result))
//	    }
//	    // Continue with next iteration
//	    return io.Of(tailrec.Bounce[int](FactState{
//	        n:      state.n - 1,
//	        result: state.result * state.n,
//	    }))
//	})
//
//	result := factorial(FactState{n: 5, result: 1})() // 120
//
// Example - Sum of numbers from 1 to N:
//
//	type SumState struct {
//	    current int
//	    limit   int
//	    sum     int
//	}
//
//	sumToN := io.TailRec(func(state SumState) io.IO[tailrec.Trampoline[SumState, int]] {
//	    if state.current > state.limit {
//	        return io.Of(tailrec.Land[SumState](state.sum))
//	    }
//	    return io.Of(tailrec.Bounce[int](SumState{
//	        current: state.current + 1,
//	        limit:   state.limit,
//	        sum:     state.sum + state.current,
//	    }))
//	})
//
//	result := sumToN(SumState{current: 1, limit: 100, sum: 0})() // 5050
//
// Example - Processing a list with accumulation:
//
//	type ListState struct {
//	    items []int
//	    acc   []int
//	}
//
//	doubleAll := io.TailRec(func(state ListState) io.IO[tailrec.Trampoline[ListState, []int]] {
//	    if len(state.items) == 0 {
//	        return io.Of(tailrec.Land[ListState](state.acc))
//	    }
//	    doubled := append(state.acc, state.items[0]*2)
//	    return io.Of(tailrec.Bounce[[]int](ListState{
//	        items: state.items[1:],
//	        acc:   doubled,
//	    }))
//	})
//
//	result := doubleAll(ListState{items: []int{1, 2, 3}, acc: []int{}})() // [2, 4, 6]
//
// Example - Fibonacci sequence:
//
//	type FibState struct {
//	    n    int
//	    prev int
//	    curr int
//	}
//
//	fibonacci := io.TailRec(func(state FibState) io.IO[tailrec.Trampoline[FibState, int]] {
//	    if state.n == 0 {
//	        return io.Of(tailrec.Land[FibState](state.curr))
//	    }
//	    return io.Of(tailrec.Bounce[int](FibState{
//	        n:    state.n - 1,
//	        prev: state.curr,
//	        curr: state.prev + state.curr,
//	    }))
//	})
//
//	result := fibonacci(FibState{n: 10, prev: 0, curr: 1})() // 55
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B] {
	return func(a A) IO[B] {
		initial := f(a)
		return func() B {
			current := initial()
			for {
				if current.Landed {
					return current.Land
				}
				current = f(current.Bounce)()
			}
		}
	}
}
