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

// Package tailrec provides a trampoline implementation for tail-call optimization in Go.
//
// # Overview
//
// Go does not support tail-call optimization (TCO) at the language level, which means
// deeply recursive functions can cause stack overflow errors. The trampoline pattern
// provides a way to convert recursive algorithms into iterative ones, avoiding stack
// overflow while maintaining the clarity of recursive code.
//
// A trampoline works by returning instructions about what to do next instead of
// directly making recursive calls. The trampoline executor then interprets these
// instructions in a loop, effectively converting recursion into iteration.
//
// # Core Concepts
//
// The package provides three main operations:
//
// **Bounce**: Indicates that the computation should continue with a new value.
// This represents a recursive call in the original algorithm.
//
// **Land**: Indicates that the computation is complete and returns a final result.
// This represents the base case in the original algorithm.
//
// **Unwrap**: Extracts the state from a Trampoline, allowing the executor to
// determine whether to continue (Bounce) or terminate (Land).
//
// # Type Parameters
//
// The Trampoline type has two type parameters:
//
//   - B: The "bounce" type - the intermediate state passed between recursive steps
//   - L: The "land" type - the final result type when computation completes
//
// # Basic Usage
//
// Converting a recursive factorial function to use trampolines:
//
//	// Traditional recursive factorial (can overflow stack)
//	func factorial(n int) int {
//	    if n <= 1 {
//	        return 1
//	    }
//	    return n * factorial(n-1)
//	}
//
//	// Trampoline-based factorial (stack-safe)
//	type State struct {
//	    n   int
//	    acc int
//	}
//
//	func factorialStep(state State) tailrec.Trampoline[State, int] {
//	    if state.n <= 1 {
//	        return tailrec.Land[State](state.acc)
//	    }
//	    return tailrec.Bounce[int](State{state.n - 1, state.acc * state.n})
//	}
//
//	// Execute the trampoline
//	func factorial(n int) int {
//	    current := tailrec.Bounce[int](State{n, 1})
//	    for {
//	        bounce, land, isLand := tailrec.Unwrap(current)
//	        if isLand {
//	            return land
//	        }
//	        current = factorialStep(bounce)
//	    }
//	}
//
// # Fibonacci Example
//
// Computing Fibonacci numbers with tail recursion:
//
//	type FibState struct {
//	    n    int
//	    curr int
//	    prev int
//	}
//
//	func fibStep(state FibState) tailrec.Trampoline[FibState, int] {
//	    if state.n <= 0 {
//	        return tailrec.Land[FibState](state.curr)
//	    }
//	    return tailrec.Bounce[int](FibState{
//	        n:    state.n - 1,
//	        curr: state.prev + state.curr,
//	        prev: state.curr,
//	    })
//	}
//
//	func fibonacci(n int) int {
//	    current := tailrec.Bounce[int](FibState{n, 1, 0})
//	    for {
//	        bounce, land, isLand := tailrec.Unwrap(current)
//	        if isLand {
//	            return land
//	        }
//	        current = fibStep(bounce)
//	    }
//	}
//
// # List Processing Example
//
// Summing a list with tail recursion:
//
//	type SumState struct {
//	    list []int
//	    sum  int
//	}
//
//	func sumStep(state SumState) tailrec.Trampoline[SumState, int] {
//	    if len(state.list) == 0 {
//	        return tailrec.Land[SumState](state.sum)
//	    }
//	    return tailrec.Bounce[int](SumState{
//	        list: state.list[1:],
//	        sum:  state.sum + state.list[0],
//	    })
//	}
//
//	func sumList(list []int) int {
//	    current := tailrec.Bounce[int](SumState{list, 0})
//	    for {
//	        bounce, land, isLand := tailrec.Unwrap(current)
//	        if isLand {
//	            return land
//	        }
//	        current = sumStep(bounce)
//	    }
//	}
//
// # Integration with Reader Monads
//
// The tailrec package is commonly used with Reader monads (readerio, context/readerio)
// to implement stack-safe recursive computations that also depend on an environment:
//
//	import (
//	    "github.com/IBM/fp-go/v2/readerio"
//	    "github.com/IBM/fp-go/v2/tailrec"
//	)
//
//	type Env struct {
//	    Multiplier int
//	}
//
//	func compute(n int) readerio.ReaderIO[Env, int] {
//	    return readerio.TailRec(
//	        n,
//	        func(n int) readerio.ReaderIO[Env, tailrec.Trampoline[int, int]] {
//	            return func(env Env) func() tailrec.Trampoline[int, int] {
//	                return func() tailrec.Trampoline[int, int] {
//	                    if n <= 0 {
//	                        return tailrec.Land[int](n * env.Multiplier)
//	                    }
//	                    return tailrec.Bounce[int](n - 1)
//	                }
//	            }
//	        },
//	    )
//	}
//
// # Benefits
//
//   - **Stack Safety**: Prevents stack overflow for deep recursion
//   - **Clarity**: Maintains the structure of recursive algorithms
//   - **Performance**: Converts recursion to iteration without manual rewriting
//   - **Composability**: Works well with functional programming patterns
//
// # When to Use
//
// Use trampolines when:
//   - You have a naturally recursive algorithm
//   - The recursion depth could be large (thousands of calls)
//   - You want to maintain the clarity of recursive code
//   - You're working with functional programming patterns
//
// # Performance Considerations
//
// While trampolines prevent stack overflow, they do have some overhead:
//   - Each step allocates a Trampoline struct
//   - The executor loop adds some indirection
//
// For shallow recursion (< 1000 calls), direct recursion may be faster.
// For deep recursion, trampolines are essential to avoid stack overflow.
//
// # Key Functions
//
// **Bounce**: Create a trampoline that continues computation with a new state
//
// **Land**: Create a trampoline that terminates with a final result
//
// **Unwrap**: Extract the state and determine if computation should continue
//
// # See Also
//
//   - readerio.TailRec: Tail-recursive Reader IO computations
//   - context/readerio.TailRec: Tail-recursive Reader IO with context
package tailrec
