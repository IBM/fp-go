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

package reader

import (
	"github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/monoid"
	G "github.com/IBM/fp-go/v2/reader/generic"
)

// MonadTraverseArray transforms each element of an array using a function that returns a Reader,
// then collects the results into a single Reader containing an array.
// This is the monadic version that takes the array as the first parameter.
//
// All Readers share the same environment R and are evaluated with it.
//
// Example:
//
//	type Config struct { Prefix string }
//	numbers := []int{1, 2, 3}
//	addPrefix := func(n int) reader.Reader[Config, string] {
//	    return reader.Asks(func(c Config) string {
//	        return fmt.Sprintf("%s%d", c.Prefix, n)
//	    })
//	}
//	r := reader.MonadTraverseArray(numbers, addPrefix)
//	result := r(Config{Prefix: "num"}) // ["num1", "num2", "num3"]
func MonadTraverseArray[R, A, B any](ma []A, f Kleisli[R, A, B]) Reader[R, []B] {
	return G.MonadTraverseArray[Reader[R, B], Reader[R, []B]](ma, f)
}

// TraverseArray transforms each element of an array using a function that returns a Reader,
// then collects the results into a single Reader containing an array.
//
// This is useful for performing a Reader computation on each element of an array
// where all computations share the same environment.
//
// Example:
//
//	type Config struct { Multiplier int }
//	multiply := func(n int) reader.Reader[Config, int] {
//	    return reader.Asks(func(c Config) int { return n * c.Multiplier })
//	}
//	transform := reader.TraverseArray(multiply)
//	r := transform([]int{1, 2, 3})
//	result := r(Config{Multiplier: 10}) // [10, 20, 30]
func TraverseArray[R, A, B any](f Kleisli[R, A, B]) Kleisli[R, []A, []B] {
	return G.TraverseArray[Reader[R, B], Reader[R, []B], []A](f)
}

// TraverseArrayWithIndex transforms each element of an array using a function that takes
// both the index and the element, returning a Reader. The results are collected into
// a single Reader containing an array.
//
// This is useful when the transformation needs to know the position of each element.
//
// Example:
//
//	type Config struct { Prefix string }
//	addIndexPrefix := func(i int, s string) reader.Reader[Config, string] {
//	    return reader.Asks(func(c Config) string {
//	        return fmt.Sprintf("%s[%d]:%s", c.Prefix, i, s)
//	    })
//	}
//	transform := reader.TraverseArrayWithIndex(addIndexPrefix)
//	r := transform([]string{"a", "b", "c"})
//	result := r(Config{Prefix: "item"}) // ["item[0]:a", "item[1]:b", "item[2]:c"]
func TraverseArrayWithIndex[R, A, B any](f func(int, A) Reader[R, B]) Kleisli[R, []A, []B] {
	return G.TraverseArrayWithIndex[Reader[R, B], Reader[R, []B], []A](f)
}

// SequenceArray converts an array of Readers into a single Reader containing an array.
// All Readers in the input array share the same environment and are evaluated with it.
//
// This is useful when you have multiple independent Reader computations and want to
// collect all their results.
//
// Example:
//
//	type Config struct { X, Y, Z int }
//	readers := []reader.Reader[Config, int]{
//	    reader.Asks(func(c Config) int { return c.X }),
//	    reader.Asks(func(c Config) int { return c.Y }),
//	    reader.Asks(func(c Config) int { return c.Z }),
//	}
//	r := reader.SequenceArray(readers)
//	result := r(Config{X: 1, Y: 2, Z: 3}) // [1, 2, 3]
func SequenceArray[R, A any](ma []Reader[R, A]) Reader[R, []A] {
	return MonadTraverseArray(ma, function.Identity[Reader[R, A]])
}

// MonadReduceArray reduces an array of Readers to a single Reader by applying a reduction function.
// This is the monadic version that takes the array of Readers as the first parameter.
//
// Each Reader is evaluated with the same environment R, and the results are accumulated using
// the provided reduce function starting from the initial value.
//
// Parameters:
//   - as: Array of Readers to reduce
//   - reduce: Binary function that combines accumulated value with each Reader's result
//   - initial: Starting value for the reduction
//
// Example:
//
//	type Config struct { Base int }
//	readers := []reader.Reader[Config, int]{
//	    reader.Asks(func(c Config) int { return c.Base + 1 }),
//	    reader.Asks(func(c Config) int { return c.Base + 2 }),
//	    reader.Asks(func(c Config) int { return c.Base + 3 }),
//	}
//	sum := func(acc, val int) int { return acc + val }
//	r := reader.MonadReduceArray(readers, sum, 0)
//	result := r(Config{Base: 10}) // 36 (11 + 12 + 13)
//
//go:inline
func MonadReduceArray[R, A, B any](as []Reader[R, A], reduce func(B, A) B, initial B) Reader[R, B] {
	return RA.MonadTraverseReduce(
		Of,
		Map,
		Ap,

		as,

		function.Identity[Reader[R, A]],
		reduce,
		initial,
	)
}

// ReduceArray returns a curried function that reduces an array of Readers to a single Reader.
// This is the curried version where the reduction function and initial value are provided first,
// returning a function that takes the array of Readers.
//
// Parameters:
//   - reduce: Binary function that combines accumulated value with each Reader's result
//   - initial: Starting value for the reduction
//
// Returns:
//   - A function that takes an array of Readers and returns a Reader of the reduced result
//
// Example:
//
//	type Config struct { Multiplier int }
//	product := func(acc, val int) int { return acc * val }
//	reducer := reader.ReduceArray[Config](product, 1)
//	readers := []reader.Reader[Config, int]{
//	    reader.Asks(func(c Config) int { return c.Multiplier * 2 }),
//	    reader.Asks(func(c Config) int { return c.Multiplier * 3 }),
//	}
//	r := reducer(readers)
//	result := r(Config{Multiplier: 5}) // 150 (10 * 15)
//
//go:inline
func ReduceArray[R, A, B any](reduce func(B, A) B, initial B) Kleisli[R, []Reader[R, A], B] {
	return RA.TraverseReduce[[]Reader[R, A]](
		Of,
		Map,
		Ap,

		function.Identity[Reader[R, A]],
		reduce,
		initial,
	)
}

// MonadReduceArrayM reduces an array of Readers using a Monoid to combine the results.
// This is the monadic version that takes the array of Readers as the first parameter.
//
// The Monoid provides both the binary operation (Concat) and the identity element (Empty)
// for the reduction, making it convenient when working with monoidal types.
//
// Parameters:
//   - as: Array of Readers to reduce
//   - m: Monoid that defines how to combine the Reader results
//
// Example:
//
//	type Config struct { Factor int }
//	readers := []reader.Reader[Config, int]{
//	    reader.Asks(func(c Config) int { return c.Factor }),
//	    reader.Asks(func(c Config) int { return c.Factor * 2 }),
//	    reader.Asks(func(c Config) int { return c.Factor * 3 }),
//	}
//	intAddMonoid := monoid.MakeMonoid(func(a, b int) int { return a + b }, 0)
//	r := reader.MonadReduceArrayM(readers, intAddMonoid)
//	result := r(Config{Factor: 5}) // 30 (5 + 10 + 15)
//
//go:inline
func MonadReduceArrayM[R, A any](as []Reader[R, A], m monoid.Monoid[A]) Reader[R, A] {
	return MonadReduceArray(as, m.Concat, m.Empty())
}

// ReduceArrayM returns a curried function that reduces an array of Readers using a Monoid.
// This is the curried version where the Monoid is provided first, returning a function
// that takes the array of Readers.
//
// The Monoid provides both the binary operation (Concat) and the identity element (Empty)
// for the reduction.
//
// Parameters:
//   - m: Monoid that defines how to combine the Reader results
//
// Returns:
//   - A function that takes an array of Readers and returns a Reader of the reduced result
//
// Example:
//
//	type Config struct { Scale int }
//	intMultMonoid := monoid.MakeMonoid(func(a, b int) int { return a * b }, 1)
//	reducer := reader.ReduceArrayM[Config](intMultMonoid)
//	readers := []reader.Reader[Config, int]{
//	    reader.Asks(func(c Config) int { return c.Scale }),
//	    reader.Asks(func(c Config) int { return c.Scale * 2 }),
//	}
//	r := reducer(readers)
//	result := r(Config{Scale: 3}) // 18 (3 * 6)
//
//go:inline
func ReduceArrayM[R, A any](m monoid.Monoid[A]) Kleisli[R, []Reader[R, A], A] {
	return ReduceArray[R](m.Concat, m.Empty())
}

// MonadTraverseReduceArray transforms and reduces an array in one operation.
// This is the monadic version that takes the array as the first parameter.
//
// First, each element is transformed using the provided Kleisli function into a Reader.
// Then, the Reader results are reduced using the provided reduction function.
//
// This is more efficient than calling TraverseArray followed by a separate reduce operation,
// as it combines both operations into a single traversal.
//
// Parameters:
//   - as: Array of elements to transform and reduce
//   - trfrm: Function that transforms each element into a Reader
//   - reduce: Binary function that combines accumulated value with each transformed result
//   - initial: Starting value for the reduction
//
// Example:
//
//	type Config struct { Multiplier int }
//	numbers := []int{1, 2, 3, 4}
//	multiply := func(n int) reader.Reader[Config, int] {
//	    return reader.Asks(func(c Config) int { return n * c.Multiplier })
//	}
//	sum := func(acc, val int) int { return acc + val }
//	r := reader.MonadTraverseReduceArray(numbers, multiply, sum, 0)
//	result := r(Config{Multiplier: 10}) // 100 (10 + 20 + 30 + 40)
//
//go:inline
func MonadTraverseReduceArray[R, A, B, C any](as []A, trfrm Kleisli[R, A, B], reduce func(C, B) C, initial C) Reader[R, C] {
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
// First, each element is transformed using the provided Kleisli function into a Reader.
// Then, the Reader results are reduced using the provided reduction function.
//
// Parameters:
//   - trfrm: Function that transforms each element into a Reader
//   - reduce: Binary function that combines accumulated value with each transformed result
//   - initial: Starting value for the reduction
//
// Returns:
//   - A function that takes an array and returns a Reader of the reduced result
//
// Example:
//
//	type Config struct { Base int }
//	addBase := func(n int) reader.Reader[Config, int] {
//	    return reader.Asks(func(c Config) int { return n + c.Base })
//	}
//	product := func(acc, val int) int { return acc * val }
//	transformer := reader.TraverseReduceArray(addBase, product, 1)
//	r := transformer([]int{2, 3, 4})
//	result := r(Config{Base: 10}) // 2184 (12 * 13 * 14)
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
// First, each element is transformed using the provided Kleisli function into a Reader.
// Then, the Reader results are reduced using the Monoid's binary operation and identity element.
//
// This combines transformation and monoidal reduction in a single efficient operation.
//
// Parameters:
//   - as: Array of elements to transform and reduce
//   - trfrm: Function that transforms each element into a Reader
//   - m: Monoid that defines how to combine the transformed results
//
// Example:
//
//	type Config struct { Offset int }
//	numbers := []int{1, 2, 3}
//	addOffset := func(n int) reader.Reader[Config, int] {
//	    return reader.Asks(func(c Config) int { return n + c.Offset })
//	}
//	intSumMonoid := monoid.MakeMonoid(func(a, b int) int { return a + b }, 0)
//	r := reader.MonadTraverseReduceArrayM(numbers, addOffset, intSumMonoid)
//	result := r(Config{Offset: 100}) // 306 (101 + 102 + 103)
//
//go:inline
func MonadTraverseReduceArrayM[R, A, B any](as []A, trfrm Kleisli[R, A, B], m monoid.Monoid[B]) Reader[R, B] {
	return MonadTraverseReduceArray(as, trfrm, m.Concat, m.Empty())
}

// TraverseReduceArrayM returns a curried function that transforms and reduces an array using a Monoid.
// This is the curried version where the transformation function and Monoid are provided first,
// returning a function that takes the array.
//
// First, each element is transformed using the provided Kleisli function into a Reader.
// Then, the Reader results are reduced using the Monoid's binary operation and identity element.
//
// Parameters:
//   - trfrm: Function that transforms each element into a Reader
//   - m: Monoid that defines how to combine the transformed results
//
// Returns:
//   - A function that takes an array and returns a Reader of the reduced result
//
// Example:
//
//	type Config struct { Factor int }
//	scale := func(n int) reader.Reader[Config, int] {
//	    return reader.Asks(func(c Config) int { return n * c.Factor })
//	}
//	intProdMonoid := monoid.MakeMonoid(func(a, b int) int { return a * b }, 1)
//	transformer := reader.TraverseReduceArrayM(scale, intProdMonoid)
//	r := transformer([]int{2, 3, 4})
//	result := r(Config{Factor: 5}) // 3000 (10 * 15 * 20)
//
//go:inline
func TraverseReduceArrayM[R, A, B any](trfrm Kleisli[R, A, B], m monoid.Monoid[B]) Kleisli[R, []A, B] {
	return TraverseReduceArray(trfrm, m.Concat, m.Empty())
}
