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
Package semigroup provides implementations of the Semigroup algebraic structure.

# Semigroup

A Semigroup is an algebraic structure consisting of a set together with an associative
binary operation. It extends the Magma structure by adding the associativity law.

Mathematical Definition:

A semigroup is a pair (S, •) where:
  - S is a set
  - • is a binary operation: S × S → S
  - The operation must be associative: (a • b) • c = a • (b • c)

The key difference from Magma is the associativity requirement, which allows operations
to be chained without worrying about parentheses.

# Basic Usage

Creating and using a semigroup:

	import (
		"fmt"
		SG "github.com/IBM/fp-go/v2/semigroup"
	)

	// Create a semigroup for string concatenation
	stringConcat := SG.MakeSemigroup(func(a, b string) string {
		return a + b
	})

	result := stringConcat.Concat("Hello, ", "World!")
	fmt.Println(result) // Output: Hello, World!

	// Associativity holds
	s1 := stringConcat.Concat(stringConcat.Concat("a", "b"), "c")
	s2 := stringConcat.Concat("a", stringConcat.Concat("b", "c"))
	fmt.Println(s1 == s2) // Output: true

# Built-in Semigroups

The package provides several pre-defined semigroups:

First - Always returns the first argument:

	first := SG.First[int]()
	result := first.Concat(1, 2) // Returns: 1

Last - Always returns the last argument:

	last := SG.Last[int]()
	result := last.Concat(1, 2) // Returns: 2

# Semigroup Transformations

Reverse - Swaps the order of arguments:

	import N "github.com/IBM/fp-go/v2/number"

	sub := SG.MakeSemigroup(func(a, b int) int { return a - b })
	reversed := SG.Reverse(sub)

	result1 := sub.Concat(10, 3)      // 10 - 3 = 7
	result2 := reversed.Concat(10, 3) // 3 - 10 = -7

FunctionSemigroup - Lifts a semigroup to work with functions:

	import N "github.com/IBM/fp-go/v2/number"

	// Semigroup for integers
	intSum := N.SemigroupSum[int]()

	// Lift to functions that return integers
	funcSG := SG.FunctionSemigroup[string](intSum)

	f := S.Size
	g := func(s string) int { return len(s) * 2 }

	// Combine functions
	combined := funcSG.Concat(f, g)
	result := combined("hello") // len("hello") + len("hello")*2 = 5 + 10 = 15

# Array Operations

ConcatAll - Concatenates all elements in an array with a starting value:

	import N "github.com/IBM/fp-go/v2/number"

	sum := N.SemigroupSum[int]()
	concatAll := SG.ConcatAll(sum)

	result := concatAll(10)([]int{1, 2, 3, 4}) // 10 + 1 + 2 + 3 + 4 = 20

MonadConcatAll - Concatenates all elements with a starting value (uncurried):

	import N "github.com/IBM/fp-go/v2/number"

	sum := N.SemigroupSum[int]()
	result := SG.MonadConcatAll(sum)([]int{1, 2, 3, 4}, 10) // 20

GenericConcatAll - Generic version for custom slice types:

	type MyInts []int

	sum := N.SemigroupSum[int]()
	concatAll := SG.GenericConcatAll[MyInts](sum)

	result := concatAll(0)(MyInts{1, 2, 3}) // 6

# Higher-Kinded Type Semigroups

ApplySemigroup - Creates a semigroup for applicative functors:

	// For a type HKT<A> with map and ap operations
	applySG := SG.ApplySemigroup(
		fmap, // func(HKT<A>, func(A) func(A) A) HKT<func(A) A>
		fap,  // func(HKT<func(A) A>, HKT<A>) HKT<A>
		baseSemigroup,
	)

AltSemigroup - Creates a semigroup for alternative functors:

	// For a type HKT<A> with an alt operation
	altSG := SG.AltSemigroup(
		falt, // func(HKT<A>, func() HKT<A>) HKT<A>
	)

# Practical Examples

Example 1: Merging Configurations

	type Config struct {
		Timeout int
		Retries int
	}

	configSG := SG.MakeSemigroup(func(a, b Config) Config {
		return Config{
			Timeout: max(a.Timeout, b.Timeout),
			Retries: a.Retries + b.Retries,
		}
	})

	default := Config{Timeout: 30, Retries: 3}
	user := Config{Timeout: 60, Retries: 5}
	override := Config{Timeout: 45, Retries: 2}

	// Merge configurations (associative)
	final := configSG.Concat(configSG.Concat(default, user), override)
	// Result: Config{Timeout: 60, Retries: 10}

Example 2: Combining Validators

	import S "github.com/IBM/fp-go/v2/string"

	type Validator func(string) []string // Returns list of errors

	validatorSG := SG.MakeSemigroup(func(v1, v2 Validator) Validator {
		return func(s string) []string {
			errors1 := v1(s)
			errors2 := v2(s)
			return append(errors1, errors2...)
		}
	})

	notEmpty := func(s string) []string {
		if S.IsEmpty(s) {
			return []string{"must not be empty"}
		}
		return nil
	}

	minLength := func(s string) []string {
		if len(s) < 3 {
			return []string{"must be at least 3 characters"}
		}
		return nil
	}

	// Combine validators
	combined := validatorSG.Concat(notEmpty, minLength)
	errors := combined("ab") // ["must be at least 3 characters"]

Example 3: Aggregating Statistics

	type Stats struct {
		Count int
		Sum   float64
		Min   float64
		Max   float64
	}

	statsSG := SG.MakeSemigroup(func(a, b Stats) Stats {
		return Stats{
			Count: a.Count + b.Count,
			Sum:   a.Sum + b.Sum,
			Min:   min(a.Min, b.Min),
			Max:   max(a.Max, b.Max),
		}
	})

	s1 := Stats{Count: 3, Sum: 15.0, Min: 2.0, Max: 8.0}
	s2 := Stats{Count: 2, Sum: 12.0, Min: 5.0, Max: 7.0}
	s3 := Stats{Count: 4, Sum: 20.0, Min: 1.0, Max: 9.0}

	// Aggregate statistics (order doesn't matter due to associativity)
	total := statsSG.Concat(statsSG.Concat(s1, s2), s3)
	// Result: Stats{Count: 9, Sum: 47.0, Min: 1.0, Max: 9.0}

Example 4: Building Query Strings

	import S "github.com/IBM/fp-go/v2/string"

	querySG := SG.MakeSemigroup(func(a, b string) string {
		if S.IsEmpty(a) {
			return b
		}
		if S.IsEmpty(b) {
			return a
		}
		return a + "&" + b
	})

	base := "api/users"
	filter := "status=active"
	sort := "sort=name"
	page := "page=1"

	// Build query string
	query := querySG.Concat(querySG.Concat(base+"?"+filter, sort), page)
	// Result: "api/users?status=active&sort=name&page=1"

# Relationship to Other Structures

Semigroup extends Magma by adding the associativity law:
  - Magma: Has a binary operation
  - Semigroup: Has an associative binary operation
  - Monoid: Semigroup with an identity element

Converting between structures:

	// Semigroup to Magma
	magma := SG.ToMagma(semigroup)

	// Semigroup to Monoid (requires identity element)
	// See the monoid package

# Laws

A valid Semigroup must satisfy the associativity law:

	// Associativity: (a • b) • c = a • (b • c)
	s.Concat(s.Concat(a, b), c) == s.Concat(a, s.Concat(b, c))

This law ensures that the order of evaluation doesn't matter, allowing for
parallel computation and optimization.

# Function Reference

Core Functions:
  - MakeSemigroup[A](func(A, A) A) Semigroup[A] - Creates a semigroup from a binary operation
  - Reverse[A](Semigroup[A]) Semigroup[A] - Returns the dual semigroup with swapped arguments
  - ToMagma[A](Semigroup[A]) Magma[A] - Converts a semigroup to a magma

Built-in Semigroups:
  - First[A]() Semigroup[A] - Always returns the first argument
  - Last[A]() Semigroup[A] - Always returns the last argument

Higher-Order Functions:
  - FunctionSemigroup[A, B](Semigroup[B]) Semigroup[func(A) B] - Lifts a semigroup to functions

Array Operations:
  - ConcatAll[A](Semigroup[A]) func(A) func([]A) A - Concatenates array elements with initial value
  - MonadConcatAll[A](Semigroup[A]) func([]A, A) A - Uncurried version of ConcatAll
  - GenericConcatAll[GA ~[]A, A](Semigroup[A]) func(A) func(GA) A - Generic version for custom slices
  - GenericMonadConcatAll[GA ~[]A, A](Semigroup[A]) func(GA, A) A - Generic uncurried version

Higher-Kinded Type Operations:
  - ApplySemigroup[A, HKTA, HKTFA](fmap, fap, Semigroup[A]) Semigroup[HKTA] - Semigroup for applicatives
  - AltSemigroup[HKTA, LAZYHKTA](falt) Semigroup[HKTA] - Semigroup for alternatives

# Related Packages

  - github.com/IBM/fp-go/v2/magma - Base algebraic structure without associativity
  - github.com/IBM/fp-go/v2/monoid - Semigroup with identity element
  - github.com/IBM/fp-go/v2/number - Numeric semigroups (sum, product, min, max)
  - github.com/IBM/fp-go/v2/string - String semigroups
  - github.com/IBM/fp-go/v2/array - Array operations using semigroups
  - github.com/IBM/fp-go/v2/function - Function composition utilities
*/
package semigroup
