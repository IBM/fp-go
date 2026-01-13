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
Package monoid provides the Monoid algebraic structure.

# Overview

A Monoid is a Semigroup with an identity element (called Empty). It extends
the Magma → Semigroup hierarchy by adding a neutral element that doesn't
change other elements when combined with them.

The Monoid interface:

	type Monoid[A any] interface {
	    Semigroup[A]  // Provides Concat(x, y A) A
	    Empty() A     // Identity element
	}

Monoid laws:
  - Associativity (from Semigroup): Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
  - Left identity: Concat(Empty(), x) = x
  - Right identity: Concat(x, Empty()) = x

# Algebraic Hierarchy

	Magma (binary operation)
	  ↓
	Semigroup (associative)
	  ↓
	Monoid (identity element)

# Basic Usage

Creating monoids for common types:

	// Integer addition monoid (identity: 0)
	addMonoid := monoid.MakeMonoid(
	    func(a, b int) int { return a + b },
	    0,
	)

	result := addMonoid.Concat(5, 3)  // 8
	empty := addMonoid.Empty()         // 0

	// Verify identity laws
	assert.Equal(t, 5, addMonoid.Concat(addMonoid.Empty(), 5))  // Left identity
	assert.Equal(t, 5, addMonoid.Concat(5, addMonoid.Empty()))  // Right identity

	// Integer multiplication monoid (identity: 1)
	mulMonoid := monoid.MakeMonoid(
	    func(a, b int) int { return a * b },
	    1,
	)

	result := mulMonoid.Concat(5, 3)  // 15
	empty := mulMonoid.Empty()         // 1

	// String concatenation monoid (identity: "")
	stringMonoid := monoid.MakeMonoid(
	    func(a, b string) string { return a + b },
	    "",
	)

	result := stringMonoid.Concat("Hello", " World")  // "Hello World"
	empty := stringMonoid.Empty()                      // ""

# Array Operations

ConcatAll - Combines all elements using the monoid's empty as starting value:

	addMonoid := monoid.MakeMonoid(
	    func(a, b int) int { return a + b },
	    0,
	)

	numbers := []int{1, 2, 3, 4, 5}
	sum := monoid.ConcatAll(addMonoid)(numbers)
	// sum is 15 (0 + 1 + 2 + 3 + 4 + 5)

	// Empty slice returns the identity
	empty := monoid.ConcatAll(addMonoid)([]int{})
	// empty is 0

Fold - Alias for ConcatAll:

	mulMonoid := monoid.MakeMonoid(
	    func(a, b int) int { return a * b },
	    1,
	)

	numbers := []int{2, 3, 4}
	product := monoid.Fold(mulMonoid)(numbers)
	// product is 24 (1 * 2 * 3 * 4)

GenericConcatAll - Works with custom slice types:

	type IntSlice []int

	addMonoid := monoid.MakeMonoid(
	    func(a, b int) int { return a + b },
	    0,
	)

	numbers := IntSlice{1, 2, 3}
	sum := monoid.GenericConcatAll[IntSlice](addMonoid)(numbers)
	// sum is 6

# Transforming Monoids

Reverse - Swaps the order of arguments:

	subMonoid := monoid.MakeMonoid(
	    func(a, b int) int { return a - b },
	    0,
	)

	reversedMonoid := monoid.Reverse(subMonoid)

	result1 := subMonoid.Concat(10, 3)         // 10 - 3 = 7
	result2 := reversedMonoid.Concat(10, 3)    // 3 - 10 = -7

ToSemigroup - Converts a Monoid to a Semigroup:

	m := monoid.MakeMonoid(
	    func(a, b int) int { return a + b },
	    0,
	)

	sg := monoid.ToSemigroup(m)
	result := sg.Concat(5, 3)  // 8

# Function Monoid

FunctionMonoid - Creates a monoid for functions when the codomain has a monoid:

	// Monoid for functions that return integers
	intAddMonoid := monoid.MakeMonoid(
	    func(a, b int) int { return a + b },
	    0,
	)

	funcMonoid := monoid.FunctionMonoid[string, int](intAddMonoid)

	f1 := S.Size
	f2 := func(s string) int { return len(s) * 2 }

	// Combine functions: result(x) = f1(x) + f2(x)
	combined := funcMonoid.Concat(f1, f2)

	result := combined("hello")  // len("hello") + len("hello")*2 = 5 + 10 = 15

	// Empty function always returns the monoid's empty
	emptyFunc := funcMonoid.Empty()
	result := emptyFunc("anything")  // 0

# Applicative Monoid

ApplicativeMonoid - Lifts a monoid into an applicative functor:

	// This is used internally for combining applicative effects
	// Example with Option-like types:

	type Option[A any] struct {
	    value *A
	}

	func Some[A any](a A) Option[A] {
	    return Option[A]{value: &a}
	}

	func None[A any]() Option[A] {
	    return Option[A]{value: nil}
	}

	// Define applicative operations for Option
	// Then use ApplicativeMonoid to combine Option values

# Alternative Monoid

AltMonoid - Creates a monoid from an Alt type class:

	// For types with alternative/fallback semantics
	// Used with Option, Either, etc. to provide fallback behavior

AlternativeMonoid - Combines applicative and alternative:

	// Advanced usage for types that are both Applicative and Alternative
	// Provides rich composition of effects with fallback

# Practical Examples

Boolean AND monoid:

	andMonoid := monoid.MakeMonoid(
	    func(a, b bool) bool { return a && b },
	    true,  // Identity: true AND x = x
	)

	values := []bool{true, true, true}
	result := monoid.ConcatAll(andMonoid)(values)  // true

	values2 := []bool{true, false, true}
	result2 := monoid.ConcatAll(andMonoid)(values2)  // false

Boolean OR monoid:

	orMonoid := monoid.MakeMonoid(
	    func(a, b bool) bool { return a || b },
	    false,  // Identity: false OR x = x
	)

	values := []bool{false, false, false}
	result := monoid.ConcatAll(orMonoid)(values)  // false

	values2 := []bool{false, true, false}
	result2 := monoid.ConcatAll(orMonoid)(values2)  // true

Max monoid (with bounded integers):

	maxMonoid := monoid.MakeMonoid(
	    func(a, b int) int {
	        if a > b {
	            return a
	        }
	        return b
	    },
	    math.MinInt,  // Identity: smallest possible int
	)

	numbers := []int{3, 7, 2, 9, 1}
	maximum := monoid.ConcatAll(maxMonoid)(numbers)  // 9

Min monoid (with bounded integers):

	minMonoid := monoid.MakeMonoid(
	    func(a, b int) int {
	        if a < b {
	            return a
	        }
	        return b
	    },
	    math.MaxInt,  // Identity: largest possible int
	)

	numbers := []int{3, 7, 2, 9, 1}
	minimum := monoid.ConcatAll(minMonoid)(numbers)  // 1

List concatenation monoid:

	listMonoid := monoid.MakeMonoid(
	    func(a, b []int) []int {
	        result := make([]int, len(a)+len(b))
	        copy(result, a)
	        copy(result[len(a):], b)
	        return result
	    },
	    []int{},  // Identity: empty slice
	)

	lists := [][]int{{1, 2}, {3, 4}, {5}}
	flattened := monoid.ConcatAll(listMonoid)(lists)
	// flattened is []int{1, 2, 3, 4, 5}

# Monoid Laws

All monoid instances must satisfy these laws:

 1. Associativity (from Semigroup):
    Concat(Concat(x, y), z) = Concat(x, Concat(y, z))

 2. Left Identity:
    Concat(Empty(), x) = x

 3. Right Identity:
    Concat(x, Empty()) = x

Example verification:

	m := monoid.MakeMonoid(
	    func(a, b int) int { return a + b },
	    0,
	)

	// Associativity
	assert.Equal(t,
	    m.Concat(m.Concat(1, 2), 3),
	    m.Concat(1, m.Concat(2, 3)),
	)  // Both equal 6

	// Left identity
	assert.Equal(t, 5, m.Concat(m.Empty(), 5))

	// Right identity
	assert.Equal(t, 5, m.Concat(5, m.Empty()))

# Common Monoids

Additive monoid (integers):
  - Concat: addition
  - Empty: 0

Multiplicative monoid (integers):
  - Concat: multiplication
  - Empty: 1

String monoid:
  - Concat: concatenation
  - Empty: ""

List monoid:
  - Concat: list concatenation
  - Empty: []

Boolean AND monoid:
  - Concat: logical AND
  - Empty: true

Boolean OR monoid:
  - Concat: logical OR
  - Empty: false

# Functions

Core operations:
  - MakeMonoid[A any](func(A, A) A, A) - Create a monoid
  - Reverse[A any](Monoid[A]) - Swap argument order
  - ToSemigroup[A any](Monoid[A]) - Convert to semigroup

Array operations:
  - ConcatAll[A any](Monoid[A]) - Combine all elements
  - Fold[A any](Monoid[A]) - Alias for ConcatAll
  - GenericConcatAll[GA ~[]A, A any](Monoid[A]) - Generic version

Higher-order:
  - FunctionMonoid[A, B any](Monoid[B]) - Monoid for functions
  - ApplicativeMonoid[A, HKTA, HKTFA any](...) - Lift into applicative
  - AltMonoid[HKTA any, LAZYHKTA ~func() HKTA](...) - From Alt type class
  - AlternativeMonoid[A, HKTA, HKTFA any, LAZYHKTA ~func() HKTA](...) - Applicative + Alternative

# Related Packages

  - semigroup: Parent structure (associative binary operation)
  - magma: Grandparent structure (binary operation)
*/
package monoid
