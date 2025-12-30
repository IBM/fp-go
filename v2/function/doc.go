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

// Package function provides functional programming utilities for function composition,
// transformation, and manipulation.
//
// This package offers a comprehensive set of tools for working with functions in a
// functional programming style, including:
//
//   - Function composition (Pipe, Flow)
//   - Currying and uncurrying
//   - Partial application (Bind)
//   - Function transformation (Flip, Swap)
//   - Utility functions (Identity, Constant, etc.)
//
// # Core Concepts
//
// Function Composition:
//
// Pipe and Flow are the primary composition primitives. They differ in the order
// of function application:
//
//   - Pipe: Applies functions left-to-right (data flows through the pipeline)
//   - Flow: Applies functions right-to-left (mathematical function composition)
//
// Example:
//
//	// Pipe: f(g(h(x)))
//	result := Pipe3(x, h, g, f)
//
//	// Flow: f(g(h(x)))
//	composed := Flow3(f, g, h)
//	result := composed(x)
//
// Currying:
//
// Currying transforms a function with multiple parameters into a sequence of
// functions each taking a single parameter.
//
// Example:
//
//	add := func(a, b int) int { return a + b }
//	curriedAdd := Curry2(add)
//	add5 := curriedAdd(5)
//	result := add5(3)  // 8
//
// Partial Application:
//
// Bind functions allow you to fix some arguments of a function, creating a new
// function with fewer parameters.
//
// Example:
//
//	multiply := func(a, b int) int { return a * b }
//	double := Bind1st(multiply, 2)
//	result := double(5)  // 10
//
// # Common Functions
//
// Identity and Constants:
//
//	Identity[A any](a A) A                    // Returns its argument unchanged
//	Constant[A any](a A) func() A             // Creates a nullary constant function
//	Constant1[B, A any](a A) func(B) A        // Creates a unary constant function
//	Constant2[B, C, A any](a A) func(B, C) A  // Creates a binary constant function
//
// Function Transformation:
//
//	Flip[T1, T2, R any](func(T1) func(T2) R) func(T2) func(T1) R  // Reverses curried function parameters
//	Swap[T1, T2, R any](func(T1, T2) R) func(T2, T1) R            // Swaps binary function parameters
//
// Pointer Utilities:
//
//	Ref[A any](a A) *A      // Creates a pointer to a value
//	Deref[A any](a *A) A    // Dereferences a pointer
//	IsNil[A any](a *A) bool // Checks if pointer is nil
//
// Selection:
//
//	First[T1, T2 any](t1 T1, t2 T2) T1   // Returns first argument
//	Second[T1, T2 any](t1 T1, t2 T2) T2  // Returns second argument
//
// Conditional:
//
//	Ternary[A, B any](pred func(A) bool, onTrue func(A) B, onFalse func(A) B) func(A) B
//	Switch[K comparable, T, R any](kf func(T) K, cases map[K]func(T) R, default func(T) R) func(T) R
//
// # Generated Functions
//
// This package includes generated functions for various arities (0-20 parameters):
//
//   - PipeN: Left-to-right composition
//   - FlowN: Right-to-left composition
//   - CurryN: Currying for N-ary functions
//   - UncurryN: Uncurrying for N-ary functions
//   - BindXofN: Partial application binding specific parameters
//   - IgnoreXofN: Partial application ignoring specific parameters
//
// # Usage Examples
//
// Basic composition:
//
//	// Transform a string: trim, lowercase, add prefix
//	process := Flow3(
//	    func(s string) string { return "processed: " + s },
//	    strings.ToLower,
//	    strings.TrimSpace,
//	)
//	result := process("  HELLO  ")  // "processed: hello"
//
// Currying and partial application:
//
//	// Create specialized functions from general ones
//	divide := func(a, b float64) float64 { return a / b }
//	divideBy2 := Bind2nd(divide, 2.0)
//	half := divideBy2(10.0)  // 5.0
//
// Working with predicates:
//
//	isPositive := N.MoreThan(0)
//	isEven := func(n int) bool { return n%2 == 0 }
//
//	classify := Ternary(
//	    isPositive,
//	    Constant1[int]("positive"),
//	    Constant1[int]("non-positive"),
//	)
//	result := classify(5)   // "positive"
//	result2 := classify(-3) // "non-positive"
//
// Memoization:
//
//	expensive := func(n int) int {
//	    time.Sleep(time.Second)
//	    return n * n
//	}
//	memoized := Memoize(expensive)
//	result1 := memoized(5)  // Takes 1 second
//	result2 := memoized(5)  // Instant (cached)
package function

//go:generate go run .. pipe --count 20 --filename gen.go

//go:generate go run .. bind --count 5 --filename binds.go
