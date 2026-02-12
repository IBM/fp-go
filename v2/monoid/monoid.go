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

// Package monoid provides an implementation of the Monoid algebraic structure.
//
// # Monoid
//
// A Monoid is an algebraic structure that extends [Semigroup] by adding an identity element.
// It consists of:
//   - A type A
//   - An associative binary operation Concat: (A, A) → A
//   - An identity element Empty: () → A
//
// # Laws
//
// A Monoid must satisfy the following laws:
//
//  1. Associativity (from Semigroup):
//     Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
//
//  2. Left Identity:
//     Concat(Empty(), x) = x
//
//  3. Right Identity:
//     Concat(x, Empty()) = x
//
// # Common Examples
//
//   - Integer addition: Concat = (+), Empty = 0
//   - Integer multiplication: Concat = (*), Empty = 1
//   - String concatenation: Concat = (++), Empty = ""
//   - List concatenation: Concat = (++), Empty = []
//   - Boolean AND: Concat = (&&), Empty = true
//   - Boolean OR: Concat = (||), Empty = false
//   - Function composition: Concat = (∘), Empty = id
//
// # References
//
//   - Haskell Data.Monoid: https://hackage.haskell.org/package/base/docs/Data-Monoid.html
//   - Fantasy Land Monoid: https://github.com/fantasyland/fantasy-land#monoid
//   - Semigroup: https://github.com/IBM/fp-go/v2/semigroup
//
// # Related Concepts
//
//   - [Semigroup]: A Monoid without the identity element requirement
//   - Magma: A set with a binary operation (no laws required)
//   - Group: A Monoid where every element has an inverse
package monoid

import (
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Monoid represents an algebraic structure with an associative binary operation and an identity element.
//
// A Monoid extends [Semigroup] by adding an identity element (Empty) that satisfies:
//   - Left identity: Concat(Empty(), x) = x
//   - Right identity: Concat(x, Empty()) = x
//
// The Monoid must also satisfy the associativity law from Semigroup:
//   - Associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
//
// # Methods
//
//   - Concat(x, y A) A: Inherited from Semigroup, combines two values associatively
//   - Empty() A: Returns the identity element for the monoid
//
// # Common Examples
//
//   - Integer addition with 0 as identity
//   - Integer multiplication with 1 as identity
//   - String concatenation with "" as identity
//   - List concatenation with [] as identity
//   - Boolean AND with true as identity
//   - Boolean OR with false as identity
//
// # References
//
//   - Haskell Monoid typeclass: https://hackage.haskell.org/package/base/docs/Data-Monoid.html#t:Monoid
//   - Fantasy Land Monoid specification: https://github.com/fantasyland/fantasy-land#monoid
type Monoid[A any] interface {
	S.Semigroup[A]
	Empty() A
}

type monoid[A any] struct {
	c func(A, A) A
	e A
}

func (m monoid[A]) Concat(x, y A) A {
	return m.c(x, y)
}

func (m monoid[A]) Empty() A {
	return m.e
}

// MakeMonoid creates a monoid from a binary operation and an identity element.
//
// The provided concat function must be associative, and the empty element must
// satisfy the identity laws (left and right identity).
//
// This is the primary constructor for creating custom monoid instances. It's the
// equivalent of defining a Monoid instance in Haskell or implementing the Fantasy Land
// Monoid specification.
//
// # Parameters
//
//   - c: An associative binary operation func(A, A) A (equivalent to Haskell's mappend or <>)
//   - e: The identity element of type A (equivalent to Haskell's mempty)
//
// # Returns
//
//   - A [Monoid][A] instance
//
// # Example
//
//	// Integer addition monoid (Sum in Haskell)
//	addMonoid := MakeMonoid(
//	    func(a, b int) int { return a + b },
//	    0,  // identity element
//	)
//	result := addMonoid.Concat(5, 3)  // 8
//	empty := addMonoid.Empty()         // 0
//
//	// String concatenation monoid
//	stringMonoid := MakeMonoid(
//	    func(a, b string) string { return a + b },
//	    "",  // identity element
//	)
//	result := stringMonoid.Concat("Hello", " World")  // "Hello World"
//
// # References
//
//   - Haskell Monoid instance: https://hackage.haskell.org/package/base/docs/Data-Monoid.html#t:Monoid
//   - Fantasy Land Monoid.empty: https://github.com/fantasyland/fantasy-land#monoid
func MakeMonoid[A any](c func(A, A) A, e A) Monoid[A] {
	return monoid[A]{c: c, e: e}
}

// Reverse returns the dual of a Monoid by swapping the arguments of Concat.
//
// The reversed monoid has the same identity element but applies the binary
// operation in the opposite order. This is useful for operations that are
// not commutative.
//
// This corresponds to the Dual newtype wrapper in Haskell's Data.Monoid, which
// provides a Monoid instance with reversed operation order.
//
// # Parameters
//
//   - m: The monoid to reverse
//
// # Returns
//
//   - A new [Monoid][A] with reversed operation order
//
// # Example
//
//	// Subtraction monoid (not commutative)
//	subMonoid := MakeMonoid(
//	    func(a, b int) int { return a - b },
//	    0,
//	)
//	reversedMonoid := Reverse(subMonoid)
//
//	result1 := subMonoid.Concat(10, 3)      // 10 - 3 = 7
//	result2 := reversedMonoid.Concat(10, 3) // 3 - 10 = -7
//
//	// String concatenation
//	stringMonoid := MakeMonoid(
//	    func(a, b string) string { return a + b },
//	    "",
//	)
//	reversed := Reverse(stringMonoid)
//	result := reversed.Concat("Hello", "World")  // "WorldHello"
//
// # References
//
//   - Haskell Data.Monoid.Dual: https://hackage.haskell.org/package/base/docs/Data-Monoid.html#t:Dual
func Reverse[A any](m Monoid[A]) Monoid[A] {
	return MakeMonoid(S.Reverse(m).Concat, m.Empty())
}

// ToSemigroup converts a Monoid to a Semigroup by discarding the identity element.
//
// This is useful when you need to use a monoid in a context that only requires
// a semigroup (associative binary operation without identity).
//
// Since every Monoid is also a Semigroup (Monoid extends Semigroup), this conversion
// is always safe. This reflects the mathematical relationship where monoids form a
// subset of semigroups.
//
// # Parameters
//
//   - m: The monoid to convert
//
// # Returns
//
//   - A [Semigroup][A] that uses the same Concat operation
//
// # Example
//
//	addMonoid := MakeMonoid(
//	    func(a, b int) int { return a + b },
//	    0,
//	)
//	sg := ToSemigroup(addMonoid)
//	result := sg.Concat(5, 3)  // 8 (identity not available)
//
// # References
//
//   - Haskell Semigroup: https://hackage.haskell.org/package/base/docs/Data-Semigroup.html
//   - Fantasy Land Semigroup: https://github.com/fantasyland/fantasy-land#semigroup
func ToSemigroup[A any](m Monoid[A]) S.Semigroup[A] {
	return S.Semigroup[A](m)
}
