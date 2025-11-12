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

package eq

import (
	F "github.com/IBM/fp-go/v2/function"
)

// Eq represents an equality type class for type T.
// It provides a way to define custom equality semantics for any type,
// not just those that are comparable with Go's == operator.
//
// Type Parameters:
//   - T: The type for which equality is defined
//
// Methods:
//   - Equals(x, y T) bool: Returns true if x and y are considered equal
//
// Laws:
// An Eq instance must satisfy the equivalence relation laws:
//  1. Reflexivity: Equals(x, x) = true for all x
//  2. Symmetry: Equals(x, y) = Equals(y, x) for all x, y
//  3. Transitivity: If Equals(x, y) and Equals(y, z), then Equals(x, z)
//
// Example:
//
//	// Create an equality predicate for integers
//	intEq := eq.FromStrictEquals[int]()
//	assert.True(t, intEq.Equals(42, 42))
//	assert.False(t, intEq.Equals(42, 43))
//
//	// Create a custom equality predicate
//	caseInsensitiveEq := eq.FromEquals(func(a, b string) bool {
//	    return strings.EqualFold(a, b)
//	})
//	assert.True(t, caseInsensitiveEq.Equals("Hello", "HELLO"))
type Eq[T any] interface {
	// Equals returns true if x and y are considered equal according to this equality predicate.
	//
	// Parameters:
	//   - x: The first value to compare
	//   - y: The second value to compare
	//
	// Returns:
	//   - true if x and y are equal, false otherwise
	Equals(x, y T) bool
}

// eq is the internal implementation of the Eq interface.
// It wraps a comparison function to provide the Eq interface.
type eq[T any] struct {
	c func(x, y T) bool
}

// Equals implements the Eq interface by delegating to the wrapped comparison function.
func (e eq[T]) Equals(x, y T) bool {
	return e.c(x, y)
}

// strictEq is a helper function that uses Go's built-in == operator for comparison.
// It can only be used with comparable types.
func strictEq[A comparable](a, b A) bool {
	return a == b
}

// FromStrictEquals constructs an Eq instance using Go's built-in == operator.
// This is the most common way to create an Eq for types that support ==.
//
// Type Parameters:
//   - T: Must be a comparable type (supports ==)
//
// Returns:
//   - An Eq[T] that uses == for equality comparison
//
// Example:
//
//	intEq := eq.FromStrictEquals[int]()
//	assert.True(t, intEq.Equals(42, 42))
//	assert.False(t, intEq.Equals(42, 43))
//
//	stringEq := eq.FromStrictEquals[string]()
//	assert.True(t, stringEq.Equals("hello", "hello"))
//	assert.False(t, stringEq.Equals("hello", "world"))
//
// Note: For types that are not comparable or require custom equality logic,
// use FromEquals instead.
func FromStrictEquals[T comparable]() Eq[T] {
	return FromEquals(strictEq[T])
}

// FromEquals constructs an Eq instance from a custom comparison function.
// This allows defining equality for any type, including non-comparable types
// or types that need custom equality semantics.
//
// Type Parameters:
//   - T: The type for which equality is being defined (can be any type)
//
// Parameters:
//   - c: A function that takes two values of type T and returns true if they are equal
//
// Returns:
//   - An Eq[T] that uses the provided function for equality comparison
//
// Example:
//
//	// Case-insensitive string equality
//	caseInsensitiveEq := eq.FromEquals(func(a, b string) bool {
//	    return strings.EqualFold(a, b)
//	})
//	assert.True(t, caseInsensitiveEq.Equals("Hello", "HELLO"))
//
//	// Approximate float equality
//	approxEq := eq.FromEquals(func(a, b float64) bool {
//	    return math.Abs(a-b) < 0.0001
//	})
//	assert.True(t, approxEq.Equals(1.0, 1.00009))
//
//	// Custom struct equality (compare by specific fields)
//	type Person struct { ID int; Name string }
//	personEq := eq.FromEquals(func(a, b Person) bool {
//	    return a.ID == b.ID  // Compare only by ID
//	})
//
// Note: The provided function should satisfy the equivalence relation laws
// (reflexivity, symmetry, transitivity) for correct behavior.
func FromEquals[T any](c func(x, y T) bool) Eq[T] {
	return eq[T]{c: c}
}

// Empty returns an Eq instance that always returns true for any comparison.
// This is the identity element for the Eq Monoid and is useful when you need
// an equality predicate that accepts everything.
//
// Type Parameters:
//   - T: The type for which the always-true equality is defined
//
// Returns:
//   - An Eq[T] where Equals(x, y) always returns true
//
// Example:
//
//	alwaysTrue := eq.Empty[int]()
//	assert.True(t, alwaysTrue.Equals(1, 2))
//	assert.True(t, alwaysTrue.Equals(42, 100))
//
//	// Useful as identity in monoid operations
//	monoid := eq.Monoid[string]()
//	combined := monoid.Concat(eq.FromStrictEquals[string](), monoid.Empty())
//	// combined behaves the same as FromStrictEquals
//
// Use cases:
//   - As the identity element in Monoid operations
//   - When you need a placeholder equality that accepts everything
//   - In generic code that requires an Eq but doesn't need actual comparison
func Empty[T any]() Eq[T] {
	return FromEquals(F.Constant2[T, T](true))
}

// Equals returns a curried equality checking function.
// This is useful for partial application and functional composition.
//
// Type Parameters:
//   - T: The type being compared
//
// Parameters:
//   - eq: The Eq instance to use for comparison
//
// Returns:
//   - A function that takes a value and returns another function that checks equality with that value
//
// Example:
//
//	intEq := eq.FromStrictEquals[int]()
//	equals42 := eq.Equals(intEq)(42)
//
//	assert.True(t, equals42(42))
//	assert.False(t, equals42(43))
//
//	// Use in higher-order functions
//	numbers := []int{40, 41, 42, 43, 44}
//	filtered := array.Filter(equals42)(numbers)
//	// filtered = [42]
//
//	// Partial application
//	equalsFunc := eq.Equals(intEq)
//	equals10 := equalsFunc(10)
//	equals20 := equalsFunc(20)
//
// This is particularly useful when working with functional programming patterns
// like map, filter, and other higher-order functions.
func Equals[T any](eq Eq[T]) func(T) func(T) bool {
	return func(other T) func(T) bool {
		return F.Bind2nd(eq.Equals, other)
	}
}
