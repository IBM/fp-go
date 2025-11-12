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

package endomorphism

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
)

// MonadAp applies an endomorphism to a value in a monadic context.
//
// This function applies the endomorphism fab to the value fa, returning the result.
// It's the monadic application operation for endomorphisms.
//
// Parameters:
//   - fab: An endomorphism to apply
//   - fa: The value to apply the endomorphism to
//
// Returns:
//   - The result of applying fab to fa
//
// Example:
//
//	double := func(x int) int { return x * 2 }
//	result := endomorphism.MonadAp(double, 5) // Returns: 10
func MonadAp[A any](fab Endomorphism[A], fa A) A {
	return identity.MonadAp(fab, fa)
}

// Ap returns a function that applies a value to an endomorphism.
//
// This is the curried version of MonadAp. It takes a value and returns a function
// that applies that value to any endomorphism.
//
// Parameters:
//   - fa: The value to be applied
//
// Returns:
//   - A function that takes an endomorphism and applies fa to it
//
// Example:
//
//	applyFive := endomorphism.Ap(5)
//	double := func(x int) int { return x * 2 }
//	result := applyFive(double) // Returns: 10
func Ap[A any](fa A) func(Endomorphism[A]) A {
	return identity.Ap[A](fa)
}

// MonadCompose composes two endomorphisms, executing them from right to left.
//
// MonadCompose creates a new endomorphism that applies f2 first, then f1.
// This follows the mathematical notation of function composition: (f1 ∘ f2)(x) = f1(f2(x))
//
// IMPORTANT: The execution order is RIGHT-TO-LEFT:
//   - f2 is applied first to the input
//   - f1 is applied to the result of f2
//
// This is different from Chain/MonadChain which executes LEFT-TO-RIGHT.
//
// Parameters:
//   - f1: The second function to apply (outer function)
//   - f2: The first function to apply (inner function)
//
// Returns:
//   - A new endomorphism that applies f2, then f1
//
// Example:
//
//	double := func(x int) int { return x * 2 }
//	increment := func(x int) int { return x + 1 }
//
//	// MonadCompose executes RIGHT-TO-LEFT: increment first, then double
//	composed := endomorphism.MonadCompose(double, increment)
//	result := composed(5) // (5 + 1) * 2 = 12
//
//	// Compare with Chain which executes LEFT-TO-RIGHT:
//	chained := endomorphism.MonadChain(double, increment)
//	result2 := chained(5) // (5 * 2) + 1 = 11
func MonadCompose[A any](f, g Endomorphism[A]) Endomorphism[A] {
	return function.Flow2(g, f)
}

// Compose returns a function that composes an endomorphism with another, executing right to left.
//
// This is the curried version of MonadCompose. It takes an endomorphism g and returns
// a function that composes any endomorphism with g, applying g first (inner function),
// then the input endomorphism (outer function).
//
// IMPORTANT: Execution order is RIGHT-TO-LEFT (mathematical composition):
//   - g is applied first to the input
//   - The endomorphism passed to the returned function is applied to the result of g
//
// This follows the mathematical composition notation where Compose(g)(f) = f ∘ g
//
// Parameters:
//   - g: The first endomorphism to apply (inner function)
//
// Returns:
//   - A function that takes an endomorphism f and composes it with g (right-to-left)
//
// Example:
//
//	increment := func(x int) int { return x + 1 }
//	composeWithIncrement := endomorphism.Compose(increment)
//	double := func(x int) int { return x * 2 }
//
//	// Composes double with increment (RIGHT-TO-LEFT: increment first, then double)
//	composed := composeWithIncrement(double)
//	result := composed(5) // (5 + 1) * 2 = 12
//
//	// Compare with Chain which executes LEFT-TO-RIGHT:
//	chainWithIncrement := endomorphism.Chain(increment)
//	chained := chainWithIncrement(double)
//	result2 := chained(5) // (5 * 2) + 1 = 11
func Compose[A any](g Endomorphism[A]) Endomorphism[Endomorphism[A]] {
	return function.Bind2nd(MonadCompose, g)
}

// MonadChain chains two endomorphisms together, executing them from left to right.
//
// This is the monadic bind operation for endomorphisms. It composes two endomorphisms
// ma and f, returning a new endomorphism that applies ma first, then f.
//
// IMPORTANT: The execution order is LEFT-TO-RIGHT:
//   - f is applied first to the input
//   - g is applied to the result of ma
//
// This is different from Compose which executes RIGHT-TO-LEFT.
//
// Parameters:
//   - f: The first endomorphism to apply
//   - g: The second endomorphism to apply
//
// Returns:
//   - A new endomorphism that applies ma, then f
//
// Example:
//
//	double := func(x int) int { return x * 2 }
//	increment := func(x int) int { return x + 1 }
//
//	// MonadChain executes LEFT-TO-RIGHT: double first, then increment
//	chained := endomorphism.MonadChain(double, increment)
//	result := chained(5) // (5 * 2) + 1 = 11
//
//	// Compare with Compose which executes RIGHT-TO-LEFT:
//	composed := endomorphism.Compose(increment, double)
//	result2 := composed(5) // (5 * 2) + 1 = 11 (same result, different parameter order)
func MonadChain[A any](f Endomorphism[A], g Endomorphism[A]) Endomorphism[A] {
	return function.Flow2(f, g)
}

// Chain returns a function that chains an endomorphism with another, executing left to right.
//
// This is the curried version of MonadChain. It takes an endomorphism f and returns
// a function that chains any endomorphism with f, applying the input endomorphism first,
// then f.
//
// IMPORTANT: Execution order is LEFT-TO-RIGHT:
//   - The endomorphism passed to the returned function is applied first
//   - f is applied to the result
//
// Parameters:
//   - f: The second endomorphism to apply
//
// Returns:
//   - A function that takes an endomorphism and chains it with f (left-to-right)
//
// Example:
//
//	increment := func(x int) int { return x + 1 }
//	chainWithIncrement := endomorphism.Chain(increment)
//	double := func(x int) int { return x * 2 }
//
//	// Chains double (first) with increment (second)
//	chained := chainWithIncrement(double)
//	result := chained(5) // (5 * 2) + 1 = 11
func Chain[A any](f Endomorphism[A]) Endomorphism[Endomorphism[A]] {
	return function.Bind2nd(MonadChain, f)
}
