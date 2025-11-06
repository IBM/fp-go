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
Package magma provides the Magma algebraic structure.

# Overview

A Magma is the most basic algebraic structure in abstract algebra. It consists
of a set equipped with a single binary operation (Concat) that combines two
elements to produce another element of the same type. Unlike more constrained
structures like Semigroup or Monoid, a Magma's operation doesn't need to be
associative or have an identity element.

The Magma interface:

	type Magma[A any] interface {
	    Concat(x A, y A) A
	}

This simple structure serves as the foundation for more complex algebraic
structures in the type class hierarchy:
  - Magma (no laws)
  - Semigroup (adds associativity)
  - Monoid (adds identity element)

# Basic Usage

Creating and using magmas:

	// Create a magma for integer addition
	addMagma := magma.MakeMagma(func(a, b int) int {
	    return a + b
	})

	result := addMagma.Concat(5, 3)
	// result is 8

	// Create a magma for string concatenation
	stringMagma := magma.MakeMagma(func(a, b string) string {
	    return a + b
	})

	result := stringMagma.Concat("Hello", " World")
	// result is "Hello World"

# Built-in Magmas

The package provides several pre-defined magmas:

First - Always returns the first argument:

	m := magma.First[int]()
	result := m.Concat(1, 2)
	// result is 1

Second - Always returns the second argument:

	m := magma.Second[int]()
	result := m.Concat(1, 2)
	// result is 2

# Transforming Magmas

Reverse - Swaps the order of arguments:

	addMagma := magma.MakeMagma(func(a, b int) int {
	    return a - b
	})

	reversedMagma := magma.Reverse(addMagma)

	result1 := addMagma.Concat(10, 3)      // 10 - 3 = 7
	result2 := reversedMagma.Concat(10, 3) // 3 - 10 = -7

FilterFirst - Only applies operation if first argument satisfies predicate:

	addMagma := magma.MakeMagma(func(a, b int) int {
	    return a + b
	})

	// Only add if first number is positive
	filteredMagma := magma.FilterFirst(func(n int) bool {
	    return n > 0
	})(addMagma)

	result1 := filteredMagma.Concat(5, 3)   // 5 + 3 = 8 (5 is positive)
	result2 := filteredMagma.Concat(-5, 3)  // 3 (5 is negative, return second)

FilterSecond - Only applies operation if second argument satisfies predicate:

	addMagma := magma.MakeMagma(func(a, b int) int {
	    return a + b
	})

	// Only add if second number is positive
	filteredMagma := magma.FilterSecond(func(n int) bool {
	    return n > 0
	})(addMagma)

	result1 := filteredMagma.Concat(5, 3)   // 5 + 3 = 8 (3 is positive)
	result2 := filteredMagma.Concat(5, -3)  // 5 (-3 is negative, return first)

Endo - Applies a function to both arguments before combining:

	addMagma := magma.MakeMagma(func(a, b int) int {
	    return a + b
	})

	// Double both numbers before adding
	doubledMagma := magma.Endo(func(n int) int {
	    return n * 2
	})(addMagma)

	result := doubledMagma.Concat(3, 4)
	// (3*2) + (4*2) = 6 + 8 = 14

# Array Operations

ConcatAll - Combines all elements in a slice using a magma:

	addMagma := magma.MakeMagma(func(a, b int) int {
	    return a + b
	})

	numbers := []int{1, 2, 3, 4, 5}
	result := magma.ConcatAll(addMagma)(0)(numbers)
	// 0 + 1 + 2 + 3 + 4 + 5 = 15

MonadConcatAll - Uncurried version:

	addMagma := magma.MakeMagma(func(a, b int) int {
	    return a + b
	})

	numbers := []int{1, 2, 3, 4, 5}
	result := magma.MonadConcatAll(addMagma)(numbers, 0)
	// 0 + 1 + 2 + 3 + 4 + 5 = 15

# Generic Array Operations

The package provides generic versions that work with custom slice types:

	type IntSlice []int

	addMagma := magma.MakeMagma(func(a, b int) int {
	    return a + b
	})

	numbers := IntSlice{1, 2, 3}
	result := magma.GenericConcatAll[IntSlice](addMagma)(0)(numbers)
	// result is 6

# Practical Examples

Building a max magma:

	maxMagma := magma.MakeMagma(func(a, b int) int {
	    if a > b {
	        return a
	    }
	    return b
	})

	numbers := []int{3, 7, 2, 9, 1}
	maximum := magma.ConcatAll(maxMagma)(0)(numbers)
	// maximum is 9

Building a min magma:

	minMagma := magma.MakeMagma(func(a, b int) int {
	    if a < b {
	        return a
	    }
	    return b
	})

	numbers := []int{3, 7, 2, 9, 1}
	minimum := magma.ConcatAll(minMagma)(10)(numbers)
	// minimum is 1

Combining strings with separator:

	joinMagma := magma.MakeMagma(func(a, b string) string {
	    if a == "" {
	        return b
	    }
	    if b == "" {
	        return a
	    }
	    return a + ", " + b
	})

	words := []string{"apple", "banana", "cherry"}
	result := magma.ConcatAll(joinMagma)("")(words)
	// result is "apple, banana, cherry"

# Relationship to Other Structures

Magma is the base of the algebraic hierarchy:

	Magma (no laws)
	  ↓
	Semigroup (associative)
	  ↓
	Monoid (identity element)

Any Semigroup or Monoid is also a Magma, but not vice versa.

# Functions

Core operations:
  - MakeMagma[A any](func(A, A) A) - Create a magma from a binary operation
  - First[A any]() - Magma that returns first argument
  - Second[A any]() - Magma that returns second argument

Transformations:
  - Reverse[A any](Magma[A]) - Swap argument order
  - FilterFirst[A any](func(A) bool) - Filter based on first argument
  - FilterSecond[A any](func(A) bool) - Filter based on second argument
  - Endo[A any](func(A) A) - Apply function before combining

Array operations:
  - ConcatAll[A any](Magma[A]) - Combine all elements (curried)
  - MonadConcatAll[A any](Magma[A]) - Combine all elements (uncurried)
  - GenericConcatAll[GA ~[]A, A any](Magma[A]) - Generic version (curried)
  - GenericMonadConcatAll[GA ~[]A, A any](Magma[A]) - Generic version (uncurried)

# Related Packages

  - semigroup: Adds associativity law to Magma
  - monoid: Adds identity element to Semigroup
*/
package magma
