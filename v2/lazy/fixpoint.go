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

package lazy

// Fixpoint computes the fixed point of a function by finding a value x such that f(x) = x.
//
// This function implements the Y combinator pattern for lazy evaluation, enabling
// recursive definitions without explicit recursion. It takes a function that accepts
// a lazy computation and returns a value, then finds the fixed point by creating
// a self-referential lazy computation.
//
// The fixpoint is computed by creating a lazy value that refers to itself. When the
// function f is applied, it receives a lazy computation that, when evaluated, returns
// the result of applying f to itself. This creates a recursive structure that can be
// used to define recursive computations in a purely functional way.
//
// Comparison with lambda.Y:
//
// Fixpoint and lambda.Y both implement fixed-point combinators but differ in their
// approach and use cases:
//
//   - Fixpoint uses lazy evaluation (Lazy[T]) to defer computation, making it simpler
//     to use and more intuitive for Go developers. The lazy parameter is only evaluated
//     when explicitly called, providing natural control over recursion termination.
//
//   - lambda.Y is a classical Y-combinator implementation using complex type constraints
//     (Endo, RecFct) that more closely follows the theoretical definition. It requires
//     more sophisticated type parameters but doesn't rely on lazy evaluation.
//
//   - Fixpoint is recommended for most practical use cases due to its simplicity and
//     clear lazy evaluation semantics. Use lambda.Y when you need a more traditional
//     Y-combinator implementation or want to avoid the lazy evaluation overhead.
//
// Type Parameters:
//   - T: The type of the value being computed
//
// Parameters:
//   - f: A function that takes a lazy computation of T and returns T. This function
//     defines the recursive structure by describing how to compute the result given
//     access to a lazy version of the result itself.
//
// Returns:
//   - T: The fixed point value where f(lazy(x)) = x
//
// Common Use Cases:
//   - Defining recursive data structures (e.g., infinite lists, trees)
//   - Implementing recursive algorithms without explicit recursion
//   - Creating self-referential computations
//   - Implementing the Y combinator pattern
//
// Note: The function f should be careful about when it evaluates the lazy parameter.
// Evaluating it immediately will cause infinite recursion. The lazy parameter should
// only be evaluated when needed, allowing the recursion to terminate based on some
// condition within f.
//
// See Also:
//   - lambda.Y: Classical Y-combinator implementation without lazy evaluation
//   - Memoize: For caching lazy computation results
//   - Defer: For creating lazy computations from generators
func Fixpoint[T any](f func(Lazy[T]) T) T {
	var x T
	x = f(func() T {
		return x
	})
	return x
}
