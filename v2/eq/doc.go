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
Package eq provides type-safe equality comparisons for any type in Go.

# Overview

The eq package implements the Eq type class from functional programming, which
represents types that support equality comparison. Unlike Go's built-in == operator
which only works with comparable types, Eq allows defining custom equality semantics
for any type, including complex structures, functions, and non-comparable types.

# Core Concepts

The Eq[T] interface represents an equality predicate for type T:

	type Eq[T any] interface {
	    Equals(x, y T) bool
	}

This abstraction enables:
  - Custom equality semantics for any type
  - Composition of equality predicates
  - Contravariant mapping to transform equality predicates
  - Monoid structure for combining multiple equality checks

# Basic Usage

Creating equality predicates for comparable types:

	// For built-in comparable types
	intEq := eq.FromStrictEquals[int]()
	assert.True(t, intEq.Equals(42, 42))
	assert.False(t, intEq.Equals(42, 43))

	stringEq := eq.FromStrictEquals[string]()
	assert.True(t, stringEq.Equals("hello", "hello"))

Creating custom equality predicates:

	// Case-insensitive string equality
	caseInsensitiveEq := eq.FromEquals(func(a, b string) bool {
	    return strings.EqualFold(a, b)
	})
	assert.True(t, caseInsensitiveEq.Equals("Hello", "HELLO"))

	// Approximate float equality
	approxEq := eq.FromEquals(func(a, b float64) bool {
	    return math.Abs(a-b) < 0.0001
	})
	assert.True(t, approxEq.Equals(1.0, 1.00009))

# Contramap - Transforming Equality

Contramap allows you to create an equality predicate for type B from an equality
predicate for type A, given a function from B to A. This is useful for comparing
complex types by extracting comparable fields:

	type Person struct {
	    ID   int
	    Name string
	    Age  int
	}

	// Compare persons by ID only
	personEqByID := eq.Contramap(func(p Person) int {
	    return p.ID
	})(eq.FromStrictEquals[int]())

	p1 := Person{ID: 1, Name: "Alice", Age: 30}
	p2 := Person{ID: 1, Name: "Bob", Age: 25}
	assert.True(t, personEqByID.Equals(p1, p2)) // Same ID

	// Compare persons by name
	personEqByName := eq.Contramap(func(p Person) string {
	    return p.Name
	})(eq.FromStrictEquals[string]())

	assert.False(t, personEqByName.Equals(p1, p2)) // Different names

# Semigroup and Monoid

The eq package provides Semigroup and Monoid instances for Eq[A], allowing you to
combine multiple equality predicates using logical AND:

	type User struct {
	    Username string
	    Email    string
	}

	// Compare by username
	usernameEq := eq.Contramap(func(u User) string {
	    return u.Username
	})(eq.FromStrictEquals[string]())

	// Compare by email
	emailEq := eq.Contramap(func(u User) string {
	    return u.Email
	})(eq.FromStrictEquals[string]())

	// Combine: users are equal if BOTH username AND email match
	userEq := eq.Semigroup[User]().Concat(usernameEq, emailEq)

	u1 := User{Username: "alice", Email: "alice@example.com"}
	u2 := User{Username: "alice", Email: "alice@example.com"}
	u3 := User{Username: "alice", Email: "different@example.com"}

	assert.True(t, userEq.Equals(u1, u2))   // Both match
	assert.False(t, userEq.Equals(u1, u3))  // Email differs

The Monoid provides an identity element (Empty) that always returns true:

	monoid := eq.Monoid[int]()
	alwaysTrue := monoid.Empty()
	assert.True(t, alwaysTrue.Equals(1, 2)) // Always true

# Curried Equality

The Equals function provides a curried version of equality checking, useful for
partial application and functional composition:

	intEq := eq.FromStrictEquals[int]()
	equals42 := eq.Equals(intEq)(42)

	assert.True(t, equals42(42))
	assert.False(t, equals42(43))

	// Use in higher-order functions
	numbers := []int{40, 41, 42, 43, 44}
	filtered := array.Filter(equals42)(numbers)
	// filtered = [42]

# Advanced Examples

Comparing slices element-wise:

	sliceEq := eq.FromEquals(func(a, b []int) bool {
	    if len(a) != len(b) {
	        return false
	    }
	    intEq := eq.FromStrictEquals[int]()
	    for i := range a {
	        if !intEq.Equals(a[i], b[i]) {
	            return false
	        }
	    }
	    return true
	})

	assert.True(t, sliceEq.Equals([]int{1, 2, 3}, []int{1, 2, 3}))
	assert.False(t, sliceEq.Equals([]int{1, 2, 3}, []int{1, 2, 4}))

Comparing maps:

	mapEq := eq.FromEquals(func(a, b map[string]int) bool {
	    if len(a) != len(b) {
	        return false
	    }
	    for k, v := range a {
	        if bv, ok := b[k]; !ok || v != bv {
	            return false
	        }
	    }
	    return true
	})

# Type Class Laws

Eq instances should satisfy the following laws:

1. Reflexivity: For all x, Equals(x, x) = true
2. Symmetry: For all x, y, Equals(x, y) = Equals(y, x)
3. Transitivity: If Equals(x, y) and Equals(y, z), then Equals(x, z)

These laws ensure that Eq behaves as a proper equivalence relation.

# Functions

  - FromStrictEquals[T comparable]() - Create Eq from Go's == operator
  - FromEquals[T any](func(x, y T) bool) - Create Eq from custom comparison
  - Empty[T any]() - Create Eq that always returns true
  - Equals[T any](Eq[T]) - Curried equality checking
  - Contramap[A, B any](func(B) A) - Transform Eq by mapping input type
  - Semigroup[A any]() - Combine Eq instances with logical AND
  - Monoid[A any]() - Semigroup with identity element

# Related Packages

  - ord: Provides ordering comparisons (less than, greater than)
  - semigroup: Provides the Semigroup abstraction used by Eq
  - monoid: Provides the Monoid abstraction used by Eq
*/
package eq
