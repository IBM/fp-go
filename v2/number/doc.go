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

/*
Package number provides algebraic structures and utility functions for numeric types.

# Overview

This package defines common algebraic structures (Magma, Semigroup, Monoid) for
numeric types, along with utility functions for arithmetic operations. It works
with any type that satisfies the Number constraint (integers, floats, complex numbers).

The Number type constraint:

	type Number interface {
	    constraints.Integer | constraints.Float | constraints.Complex
	}

# Algebraic Structures

Magma - Binary operations (not necessarily associative):
  - MagmaSub[A Number]() - Subtraction magma
  - MagmaDiv[A Number]() - Division magma

Semigroup - Associative binary operations:
  - SemigroupSum[A Number]() - Addition semigroup
  - SemigroupProduct[A Number]() - Multiplication semigroup

Monoid - Semigroups with identity elements:
  - MonoidSum[A Number]() - Addition monoid (identity: 0)
  - MonoidProduct[A Number]() - Multiplication monoid (identity: 1)

# Basic Usage

Using monoids for numeric operations:

	// Sum monoid
	sumMonoid := number.MonoidSum[int]()
	result := sumMonoid.Concat(5, 3)  // 8
	empty := sumMonoid.Empty()         // 0

	// Product monoid
	prodMonoid := number.MonoidProduct[int]()
	result := prodMonoid.Concat(5, 3)  // 15
	empty := prodMonoid.Empty()         // 1

Using semigroups:

	// Addition semigroup
	addSemigroup := number.SemigroupSum[int]()
	result := addSemigroup.Concat(10, 20)  // 30

	// Multiplication semigroup
	mulSemigroup := number.SemigroupProduct[float64]()
	result := mulSemigroup.Concat(2.5, 4.0)  // 10.0

Using magmas (non-associative operations):

	// Subtraction magma
	subMagma := number.MagmaSub[int]()
	result := subMagma.Concat(10, 3)  // 7

	// Division magma
	divMagma := number.MagmaDiv[float64]()
	result := divMagma.Concat(10.0, 2.0)  // 5.0

# Curried Arithmetic Functions

The package provides curried versions of arithmetic operations for use in
functional composition:

	// Add - curried addition
	add5 := number.Add(5)
	result := add5(10)  // 15

	// Sub - curried subtraction
	sub3 := number.Sub(3)
	result := sub3(10)  // 7

	// Mul - curried multiplication
	double := number.Mul(2)
	result := double(5)  // 10

	// Div - curried division
	half := number.Div[float64](2)
	result := half(10.0)  // 5.0

Using with array operations:

	import (
	    A "github.com/IBM/fp-go/v2/array"
	    N "github.com/IBM/fp-go/v2/number"
	)

	numbers := []int{1, 2, 3, 4, 5}

	// Add 10 to each number
	result := A.Map(N.Add(10))(numbers)
	// result: [11, 12, 13, 14, 15]

	// Double each number
	doubled := A.Map(N.Mul(2))(numbers)
	// doubled: [2, 4, 6, 8, 10]

# Utility Functions

Inc - Increment a number:

	result := number.Inc(5)  // 6
	result := number.Inc(2.5)  // 3.5

Min - Get the minimum of two values:

	result := number.Min(5, 3)  // 3
	result := number.Min(2.5, 7.8)  // 2.5

Max - Get the maximum of two values:

	result := number.Max(5, 3)  // 5
	result := number.Max(2.5, 7.8)  // 7.8

# Working with Different Numeric Types

Integers:

	sumMonoid := number.MonoidSum[int]()
	result := sumMonoid.Concat(100, 200)  // 300

	add10 := number.Add(10)
	result := add10(5)  // 15

Floats:

	sumMonoid := number.MonoidSum[float64]()
	result := sumMonoid.Concat(3.14, 2.86)  // 6.0

	half := number.Div[float64](2.0)
	result := half(10.5)  // 5.25

Complex numbers:

	sumMonoid := number.MonoidSum[complex128]()
	c1 := complex(1, 2)
	c2 := complex(3, 4)
	result := sumMonoid.Concat(c1, c2)  // (4+6i)

# Combining with Monoid Operations

Using with monoid.ConcatAll:

	import (
	    M "github.com/IBM/fp-go/v2/monoid"
	    N "github.com/IBM/fp-go/v2/number"
	)

	// Sum all numbers
	sumMonoid := N.MonoidSum[int]()
	numbers := []int{1, 2, 3, 4, 5}
	total := M.ConcatAll(sumMonoid)(numbers)
	// total: 15

	// Product of all numbers
	prodMonoid := N.MonoidProduct[int]()
	product := M.ConcatAll(prodMonoid)(numbers)
	// product: 120

# Practical Examples

Calculate average:

	import (
	    M "github.com/IBM/fp-go/v2/monoid"
	    N "github.com/IBM/fp-go/v2/number"
	)

	numbers := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
	sumMonoid := N.MonoidSum[float64]()
	sum := M.ConcatAll(sumMonoid)(numbers)
	average := sum / float64(len(numbers))
	// average: 30.0

Factorial using product monoid:

	import (
	    M "github.com/IBM/fp-go/v2/monoid"
	    N "github.com/IBM/fp-go/v2/number"
	)

	factorial := func(n int) int {
	    if n <= 1 {
	        return 1
	    }
	    numbers := make([]int, n)
	    for i := range numbers {
	        numbers[i] = i + 1
	    }
	    prodMonoid := N.MonoidProduct[int]()
	    return M.ConcatAll(prodMonoid)(numbers)
	}

	result := factorial(5)  // 120

Transform and sum:

	import (
	    A "github.com/IBM/fp-go/v2/array"
	    F "github.com/IBM/fp-go/v2/function"
	    M "github.com/IBM/fp-go/v2/monoid"
	    N "github.com/IBM/fp-go/v2/number"
	)

	// Sum of squares
	numbers := []int{1, 2, 3, 4, 5}
	squares := A.Map(func(x int) int { return x * x })(numbers)
	sumMonoid := N.MonoidSum[int]()
	sumOfSquares := M.ConcatAll(sumMonoid)(squares)
	// sumOfSquares: 55 (1 + 4 + 9 + 16 + 25)

# Type Safety

All functions are type-safe and work with any numeric type:

	// Works with int
	intSum := number.MonoidSum[int]()

	// Works with float64
	floatSum := number.MonoidSum[float64]()

	// Works with complex128
	complexSum := number.MonoidSum[complex128]()

	// Compile-time error for non-numeric types
	// stringSum := number.MonoidSum[string]()  // Error!

# Functions

Algebraic structures:
  - MagmaSub[A Number]() - Subtraction magma
  - MagmaDiv[A Number]() - Division magma
  - SemigroupSum[A Number]() - Addition semigroup
  - SemigroupProduct[A Number]() - Multiplication semigroup
  - MonoidSum[A Number]() - Addition monoid (identity: 0)
  - MonoidProduct[A Number]() - Multiplication monoid (identity: 1)

Curried arithmetic:
  - Add[T Number](T) func(T) T - Curried addition
  - Sub[T Number](T) func(T) T - Curried subtraction
  - Mul[T Number](T) func(T) T - Curried multiplication
  - Div[T Number](T) func(T) T - Curried division

Utilities:
  - Inc[T Number](T) T - Increment by 1
  - Min[A Ordered](A, A) A - Minimum of two values
  - Max[A Ordered](A, A) A - Maximum of two values

# Related Packages

  - magma: Base algebraic structure (binary operation)
  - semigroup: Associative binary operation
  - monoid: Semigroup with identity element
  - constraints: Type constraints for generics
*/
package number
