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
)

// MonadAp applies an endomorphism in a function to an endomorphism value.
//
// For endomorphisms, Ap composes two endomorphisms using RIGHT-TO-LEFT composition.
// This is the applicative functor operation for endomorphisms.
//
// IMPORTANT: Execution order is RIGHT-TO-LEFT (same as MonadCompose):
//   - fa is applied first to the input
//   - fab is applied to the result
//
// Parameters:
//   - fab: An endomorphism to apply (outer function)
//   - fa: An endomorphism to apply first (inner function)
//
// Returns:
//   - A new endomorphism that applies fa, then fab
//
// Example:
//
//	double := N.Mul(2)
//	increment := N.Add(1)
//	result := endomorphism.MonadAp(double, increment) // Composes: double ∘ increment
//	// result(5) = double(increment(5)) = double(6) = 12
func MonadAp[A any](fab, fa Endomorphism[A]) Endomorphism[A] {
	return MonadCompose(fab, fa)
}

// Ap returns a function that applies an endomorphism to another endomorphism.
//
// This is the curried version of MonadAp. It takes an endomorphism fa and returns
// a function that composes any endomorphism with fa using RIGHT-TO-LEFT composition.
//
// IMPORTANT: Execution order is RIGHT-TO-LEFT:
//   - fa is applied first to the input
//   - The endomorphism passed to the returned function is applied to the result
//
// Parameters:
//   - fa: The first endomorphism to apply (inner function)
//
// Returns:
//   - A function that takes an endomorphism and composes it with fa (right-to-left)
//
// Example:
//
//	increment := N.Add(1)
//	applyIncrement := endomorphism.Ap(increment)
//	double := N.Mul(2)
//	composed := applyIncrement(double) // double ∘ increment
//	// composed(5) = double(increment(5)) = double(6) = 12
func Ap[A any](fa Endomorphism[A]) Operator[A] {
	return Compose(fa)
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
//	double := N.Mul(2)
//	increment := N.Add(1)
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

// MonadMap maps an endomorphism over another endomorphism using function composition.
//
// For endomorphisms, Map is equivalent to Compose (RIGHT-TO-LEFT composition).
// This is the functor map operation for endomorphisms.
//
// IMPORTANT: Execution order is RIGHT-TO-LEFT:
//   - g is applied first to the input
//   - f is applied to the result
//
// Parameters:
//   - f: The function to map (outer function)
//   - g: The endomorphism to map over (inner function)
//
// Returns:
//   - A new endomorphism that applies g, then f
//
// Example:
//
//	double := N.Mul(2)
//	increment := N.Add(1)
//	mapped := endomorphism.MonadMap(double, increment)
//	// mapped(5) = double(increment(5)) = double(6) = 12
func MonadMap[A any](f, g Endomorphism[A]) Endomorphism[A] {
	return MonadCompose(f, g)
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
//	increment := N.Add(1)
//	composeWithIncrement := endomorphism.Compose(increment)
//	double := N.Mul(2)
//
//	// Composes double with increment (RIGHT-TO-LEFT: increment first, then double)
//	composed := composeWithIncrement(double)
//	result := composed(5) // (5 + 1) * 2 = 12
//
//	// Compare with Chain which executes LEFT-TO-RIGHT:
//	chainWithIncrement := endomorphism.Chain(increment)
//	chained := chainWithIncrement(double)
//	result2 := chained(5) // (5 * 2) + 1 = 11
func Compose[A any](g Endomorphism[A]) Operator[A] {
	return function.Bind2nd(MonadCompose, g)
}

// Map returns a function that maps an endomorphism over another endomorphism.
//
// This is the curried version of MonadMap. It takes an endomorphism f and returns
// a function that maps f over any endomorphism using RIGHT-TO-LEFT composition.
//
// IMPORTANT: Execution order is RIGHT-TO-LEFT (same as Compose):
//   - The endomorphism passed to the returned function is applied first
//   - f is applied to the result
//
// For endomorphisms, Map is equivalent to Compose.
//
// Parameters:
//   - f: The function to map (outer function)
//
// Returns:
//   - A function that takes an endomorphism and maps f over it (right-to-left)
//
// Example:
//
//	double := N.Mul(2)
//	mapDouble := endomorphism.Map(double)
//	increment := N.Add(1)
//	mapped := mapDouble(increment)
//	// mapped(5) = double(increment(5)) = double(6) = 12
func Map[A any](f Endomorphism[A]) Operator[A] {
	return Compose(f)
}

// MonadChain chains two endomorphisms together, executing them from left to right.
//
// This is the monadic bind operation for endomorphisms. For endomorphisms, bind is
// simply left-to-right function composition: ma is applied first, then f.
//
// IMPORTANT: The execution order is LEFT-TO-RIGHT:
//   - ma is applied first to the input
//   - f is applied to the result of ma
//
// This is different from MonadCompose which executes RIGHT-TO-LEFT.
//
// Parameters:
//   - ma: The first endomorphism to apply
//   - f: The second endomorphism to apply
//
// Returns:
//   - A new endomorphism that applies ma, then f
//
// Example:
//
//	double := N.Mul(2)
//	increment := N.Add(1)
//
//	// MonadChain executes LEFT-TO-RIGHT: double first, then increment
//	chained := endomorphism.MonadChain(double, increment)
//	result := chained(5) // (5 * 2) + 1 = 11
//
//	// Compare with MonadCompose which executes RIGHT-TO-LEFT:
//	composed := endomorphism.MonadCompose(increment, double)
//	result2 := composed(5) // (5 * 2) + 1 = 11 (same result, different parameter order)
func MonadChain[A any](ma, f Endomorphism[A]) Endomorphism[A] {
	return function.Flow2(ma, f)
}

// MonadChainFirst chains two endomorphisms but returns the result of the first.
//
// This applies ma first, then f, but discards the result of f and returns the result of ma.
// Useful for performing side-effects while preserving the original value.
//
// Parameters:
//   - ma: The endomorphism whose result to keep
//   - f: The endomorphism to apply for its effect
//
// Returns:
//   - A new endomorphism that applies both but returns ma's result
//
// Example:
//
//	double := N.Mul(2)
//	log := func(x int) int { fmt.Println(x); return x }
//	chained := endomorphism.MonadChainFirst(double, log)
//	result := chained(5) // Prints 10, returns 10
func MonadChainFirst[A any](ma, f Endomorphism[A]) Endomorphism[A] {
	return func(a A) A {
		result := ma(a)
		f(result)     // Apply f for its effect
		return result // But return ma's result
	}
}

// ChainFirst returns a function that chains for effect but preserves the original result.
//
// This is the curried version of MonadChainFirst.
//
// Parameters:
//   - f: The endomorphism to apply for its effect
//
// Returns:
//   - A function that takes an endomorphism and chains it with f, keeping the first result
//
// Example:
//
//	log := func(x int) int { fmt.Println(x); return x }
//	chainLog := endomorphism.ChainFirst(log)
//	double := N.Mul(2)
//	chained := chainLog(double)
//	result := chained(5) // Prints 10, returns 10
func ChainFirst[A any](f Endomorphism[A]) Operator[A] {
	return function.Bind2nd(MonadChainFirst, f)
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
//	increment := N.Add(1)
//	chainWithIncrement := endomorphism.Chain(increment)
//	double := N.Mul(2)
//
//	// Chains double (first) with increment (second)
//	chained := chainWithIncrement(double)
//	result := chained(5) // (5 * 2) + 1 = 11
func Chain[A any](f Endomorphism[A]) Operator[A] {
	return function.Bind2nd(MonadChain, f)
}

// Flatten collapses a nested endomorphism into a single endomorphism.
//
// Given an endomorphism that transforms endomorphisms (Endomorphism[Endomorphism[A]]),
// Flatten produces a simple endomorphism by applying the outer transformation to the
// identity function. This is the monadic join operation for the Endomorphism monad.
//
// The function applies the nested endomorphism to Identity[A] to extract the inner
// endomorphism, effectively "flattening" the two layers into one.
//
// Type Parameters:
//   - A: The type being transformed by the endomorphisms
//
// Parameters:
//   - mma: A nested endomorphism that transforms endomorphisms
//
// Returns:
//   - An endomorphism that applies the transformation directly to values of type A
//
// Example:
//
//	type Counter struct {
//	    Value int
//	}
//
//	// An endomorphism that wraps another endomorphism
//	addThenDouble := func(endo Endomorphism[Counter]) Endomorphism[Counter] {
//	    return func(c Counter) Counter {
//	        c = endo(c)        // Apply the input endomorphism
//	        c.Value = c.Value * 2  // Then double
//	        return c
//	    }
//	}
//
//	flattened := Flatten(addThenDouble)
//	result := flattened(Counter{Value: 5})  // Counter{Value: 10}
func Flatten[A any](mma Endomorphism[Endomorphism[A]]) Endomorphism[A] {
	return mma(function.Identity[A])
}

// Join performs self-application of a function that produces endomorphisms.
//
// Given a function that takes a value and returns an endomorphism of that same type,
// Join creates an endomorphism that applies the value to itself through the function.
// This operation is also known as the W combinator (warbler) in combinatory logic,
// or diagonal application.
//
// The resulting endomorphism evaluates f(a)(a), applying the same value a to both
// the function f and the resulting endomorphism.
//
// Type Parameters:
//   - A: The type being transformed
//
// Parameters:
//   - f: A function that takes a value and returns an endomorphism of that type
//
// Returns:
//   - An endomorphism that performs self-application: f(a)(a)
//
// Example:
//
//	type Point struct {
//	    X, Y int
//	}
//
//	// Create an endomorphism based on the input point
//	scaleBy := func(p Point) Endomorphism[Point] {
//	    return func(p2 Point) Point {
//	        return Point{
//	            X: p2.X * p.X,
//	            Y: p2.Y * p.Y,
//	        }
//	    }
//	}
//
//	selfScale := Join(scaleBy)
//	result := selfScale(Point{X: 3, Y: 4})  // Point{X: 9, Y: 16}
func Join[A any](f Kleisli[A]) Endomorphism[A] {
	return func(a A) A {
		return f(a)(a)
	}
}
