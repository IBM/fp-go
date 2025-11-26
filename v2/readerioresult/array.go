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

package readerioresult

import (
	"github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/monoid"
)

// MonadReduceArray reduces an array of ReaderIOResults to a single ReaderIOResult by applying a reduction function.
// This is the monadic version that takes the array of ReaderIOResults as the first parameter.
//
// Each ReaderIOResult is evaluated with the same environment R, and the results are accumulated using
// the provided reduce function starting from the initial value. If any ReaderIOResult fails, the entire
// operation fails with that error.
//
// Parameters:
//   - as: Array of ReaderIOResults to reduce
//   - reduce: Binary function that combines accumulated value with each ReaderIOResult's result
//   - initial: Starting value for the reduction
//
// Example:
//
//	type Config struct { Base int }
//	readers := []readerioresult.ReaderIOResult[Config, int]{
//	    readerioresult.Of[Config](func(c Config) int { return c.Base + 1 }),
//	    readerioresult.Of[Config](func(c Config) int { return c.Base + 2 }),
//	    readerioresult.Of[Config](func(c Config) int { return c.Base + 3 }),
//	}
//	sum := func(acc, val int) int { return acc + val }
//	r := readerioresult.MonadReduceArray(readers, sum, 0)
//	result := r(Config{Base: 10})() // result.Of(36) (11 + 12 + 13)
//
//go:inline
func MonadReduceArray[R, A, B any](as []ReaderIOResult[R, A], reduce func(B, A) B, initial B) ReaderIOResult[R, B] {
	return RA.MonadTraverseReduce(
		Of,
		Map,
		Ap,

		as,

		function.Identity[ReaderIOResult[R, A]],
		reduce,
		initial,
	)
}

// ReduceArray returns a curried function that reduces an array of ReaderIOResults to a single ReaderIOResult.
// This is the curried version where the reduction function and initial value are provided first,
// returning a function that takes the array of ReaderIOResults.
//
// Parameters:
//   - reduce: Binary function that combines accumulated value with each ReaderIOResult's result
//   - initial: Starting value for the reduction
//
// Returns:
//   - A function that takes an array of ReaderIOResults and returns a ReaderIOResult of the reduced result
//
// Example:
//
//	type Config struct { Multiplier int }
//	product := func(acc, val int) int { return acc * val }
//	reducer := readerioresult.ReduceArray[Config](product, 1)
//	readers := []readerioresult.ReaderIOResult[Config, int]{
//	    readerioresult.Of[Config](func(c Config) int { return c.Multiplier * 2 }),
//	    readerioresult.Of[Config](func(c Config) int { return c.Multiplier * 3 }),
//	}
//	r := reducer(readers)
//	result := r(Config{Multiplier: 5})() // result.Of(150) (10 * 15)
//
//go:inline
func ReduceArray[R, A, B any](reduce func(B, A) B, initial B) Kleisli[R, []ReaderIOResult[R, A], B] {
	return RA.TraverseReduce[[]ReaderIOResult[R, A]](
		Of,
		Map,
		Ap,

		function.Identity[ReaderIOResult[R, A]],
		reduce,
		initial,
	)
}

// MonadReduceArrayM reduces an array of ReaderIOResults using a Monoid to combine the results.
// This is the monadic version that takes the array of ReaderIOResults as the first parameter.
//
// The Monoid provides both the binary operation (Concat) and the identity element (Empty)
// for the reduction, making it convenient when working with monoidal types. If any ReaderIOResult
// fails, the entire operation fails with that error.
//
// Parameters:
//   - as: Array of ReaderIOResults to reduce
//   - m: Monoid that defines how to combine the ReaderIOResult results
//
// Example:
//
//	type Config struct { Factor int }
//	readers := []readerioresult.ReaderIOResult[Config, int]{
//	    readerioresult.Of[Config](func(c Config) int { return c.Factor }),
//	    readerioresult.Of[Config](func(c Config) int { return c.Factor * 2 }),
//	    readerioresult.Of[Config](func(c Config) int { return c.Factor * 3 }),
//	}
//	intAddMonoid := monoid.MakeMonoid(func(a, b int) int { return a + b }, 0)
//	r := readerioresult.MonadReduceArrayM(readers, intAddMonoid)
//	result := r(Config{Factor: 5})() // result.Of(30) (5 + 10 + 15)
//
//go:inline
func MonadReduceArrayM[R, A any](as []ReaderIOResult[R, A], m monoid.Monoid[A]) ReaderIOResult[R, A] {
	return MonadReduceArray(as, m.Concat, m.Empty())
}

// ReduceArrayM returns a curried function that reduces an array of ReaderIOResults using a Monoid.
// This is the curried version where the Monoid is provided first, returning a function
// that takes the array of ReaderIOResults.
//
// The Monoid provides both the binary operation (Concat) and the identity element (Empty)
// for the reduction.
//
// Parameters:
//   - m: Monoid that defines how to combine the ReaderIOResult results
//
// Returns:
//   - A function that takes an array of ReaderIOResults and returns a ReaderIOResult of the reduced result
//
// Example:
//
//	type Config struct { Scale int }
//	intMultMonoid := monoid.MakeMonoid(func(a, b int) int { return a * b }, 1)
//	reducer := readerioresult.ReduceArrayM[Config](intMultMonoid)
//	readers := []readerioresult.ReaderIOResult[Config, int]{
//	    readerioresult.Of[Config](func(c Config) int { return c.Scale }),
//	    readerioresult.Of[Config](func(c Config) int { return c.Scale * 2 }),
//	}
//	r := reducer(readers)
//	result := r(Config{Scale: 3})() // result.Of(18) (3 * 6)
//
//go:inline
func ReduceArrayM[R, A any](m monoid.Monoid[A]) Kleisli[R, []ReaderIOResult[R, A], A] {
	return ReduceArray[R](m.Concat, m.Empty())
}

// MonadTraverseReduceArray transforms and reduces an array in one operation.
// This is the monadic version that takes the array as the first parameter.
//
// First, each element is transformed using the provided Kleisli function into a ReaderIOResult.
// Then, the ReaderIOResult results are reduced using the provided reduction function.
// If any transformation fails, the entire operation fails with that error.
//
// This is more efficient than calling TraverseArray followed by a separate reduce operation,
// as it combines both operations into a single traversal.
//
// Parameters:
//   - as: Array of elements to transform and reduce
//   - trfrm: Function that transforms each element into a ReaderIOResult
//   - reduce: Binary function that combines accumulated value with each transformed result
//   - initial: Starting value for the reduction
//
// Example:
//
//	type Config struct { Multiplier int }
//	numbers := []int{1, 2, 3, 4}
//	multiply := func(n int) readerioresult.ReaderIOResult[Config, int] {
//	    return readerioresult.Of[Config](func(c Config) int { return n * c.Multiplier })
//	}
//	sum := func(acc, val int) int { return acc + val }
//	r := readerioresult.MonadTraverseReduceArray(numbers, multiply, sum, 0)
//	result := r(Config{Multiplier: 10})() // result.Of(100) (10 + 20 + 30 + 40)
//
//go:inline
func MonadTraverseReduceArray[R, A, B, C any](as []A, trfrm Kleisli[R, A, B], reduce func(C, B) C, initial C) ReaderIOResult[R, C] {
	return RA.MonadTraverseReduce(
		Of,
		Map,
		Ap,

		as,

		trfrm,
		reduce,
		initial,
	)
}

// TraverseReduceArray returns a curried function that transforms and reduces an array.
// This is the curried version where the transformation function, reduce function, and initial value
// are provided first, returning a function that takes the array.
//
// First, each element is transformed using the provided Kleisli function into a ReaderIOResult.
// Then, the ReaderIOResult results are reduced using the provided reduction function.
//
// Parameters:
//   - trfrm: Function that transforms each element into a ReaderIOResult
//   - reduce: Binary function that combines accumulated value with each transformed result
//   - initial: Starting value for the reduction
//
// Returns:
//   - A function that takes an array and returns a ReaderIOResult of the reduced result
//
// Example:
//
//	type Config struct { Base int }
//	addBase := func(n int) readerioresult.ReaderIOResult[Config, int] {
//	    return readerioresult.Of[Config](func(c Config) int { return n + c.Base })
//	}
//	product := func(acc, val int) int { return acc * val }
//	transformer := readerioresult.TraverseReduceArray(addBase, product, 1)
//	r := transformer([]int{2, 3, 4})
//	result := r(Config{Base: 10})() // result.Of(2184) (12 * 13 * 14)
//
//go:inline
func TraverseReduceArray[R, A, B, C any](trfrm Kleisli[R, A, B], reduce func(C, B) C, initial C) Kleisli[R, []A, C] {
	return RA.TraverseReduce[[]A](
		Of,
		Map,
		Ap,

		trfrm,
		reduce,
		initial,
	)
}

// MonadTraverseReduceArrayM transforms and reduces an array using a Monoid.
// This is the monadic version that takes the array as the first parameter.
//
// First, each element is transformed using the provided Kleisli function into a ReaderIOResult.
// Then, the ReaderIOResult results are reduced using the Monoid's binary operation and identity element.
// If any transformation fails, the entire operation fails with that error.
//
// This combines transformation and monoidal reduction in a single efficient operation.
//
// Parameters:
//   - as: Array of elements to transform and reduce
//   - trfrm: Function that transforms each element into a ReaderIOResult
//   - m: Monoid that defines how to combine the transformed results
//
// Example:
//
//	type Config struct { Offset int }
//	numbers := []int{1, 2, 3}
//	addOffset := func(n int) readerioresult.ReaderIOResult[Config, int] {
//	    return readerioresult.Of[Config](func(c Config) int { return n + c.Offset })
//	}
//	intSumMonoid := monoid.MakeMonoid(func(a, b int) int { return a + b }, 0)
//	r := readerioresult.MonadTraverseReduceArrayM(numbers, addOffset, intSumMonoid)
//	result := r(Config{Offset: 100})() // result.Of(306) (101 + 102 + 103)
//
//go:inline
func MonadTraverseReduceArrayM[R, A, B any](as []A, trfrm Kleisli[R, A, B], m monoid.Monoid[B]) ReaderIOResult[R, B] {
	return MonadTraverseReduceArray(as, trfrm, m.Concat, m.Empty())
}

// TraverseReduceArrayM returns a curried function that transforms and reduces an array using a Monoid.
// This is the curried version where the transformation function and Monoid are provided first,
// returning a function that takes the array.
//
// First, each element is transformed using the provided Kleisli function into a ReaderIOResult.
// Then, the ReaderIOResult results are reduced using the Monoid's binary operation and identity element.
//
// Parameters:
//   - trfrm: Function that transforms each element into a ReaderIOResult
//   - m: Monoid that defines how to combine the transformed results
//
// Returns:
//   - A function that takes an array and returns a ReaderIOResult of the reduced result
//
// Example:
//
//	type Config struct { Factor int }
//	scale := func(n int) readerioresult.ReaderIOResult[Config, int] {
//	    return readerioresult.Of[Config](func(c Config) int { return n * c.Factor })
//	}
//	intProdMonoid := monoid.MakeMonoid(func(a, b int) int { return a * b }, 1)
//	transformer := readerioresult.TraverseReduceArrayM(scale, intProdMonoid)
//	r := transformer([]int{2, 3, 4})
//	result := r(Config{Factor: 5})() // result.Of(3000) (10 * 15 * 20)
//
//go:inline
func TraverseReduceArrayM[R, A, B any](trfrm Kleisli[R, A, B], m monoid.Monoid[B]) Kleisli[R, []A, B] {
	return TraverseReduceArray(trfrm, m.Concat, m.Empty())
}
