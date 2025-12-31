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

package predicate

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/semigroup"
)

type (
	// Semigroup represents a semigroup instance for predicates, providing a way to combine
	// two predicates into one using an associative operation.
	Semigroup[A any] = semigroup.Semigroup[Predicate[A]]

	// Monoid represents a monoid instance for predicates, extending Semigroup with an
	// identity element (empty predicate).
	Monoid[A any] = monoid.Monoid[Predicate[A]]
)

// SemigroupAny creates a semigroup that combines predicates using logical OR (||).
//
// When two predicates are combined with this semigroup, the resulting predicate returns
// true if either of the original predicates returns true. This implements the associative
// operation for disjunction.
//
// Returns:
//   - A Semigroup[A] that combines predicates with OR logic
//
// Example:
//
//	s := SemigroupAny[int]()
//	isPositive := N.MoreThan(0)
//	isEven := func(n int) bool { return n%2 == 0 }
//	isPositiveOrEven := s.Concat(isPositive, isEven)
//	isPositiveOrEven(4)  // true (even)
//	isPositiveOrEven(3)  // true (positive)
//	isPositiveOrEven(-2) // true (even)
//	isPositiveOrEven(-3) // false (neither)
func SemigroupAny[A any]() Semigroup[A] {
	return semigroup.MakeSemigroup(func(first Predicate[A], second Predicate[A]) Predicate[A] {
		return F.Pipe1(
			first,
			Or(second),
		)
	})
}

// SemigroupAll creates a semigroup that combines predicates using logical AND (&&).
//
// When two predicates are combined with this semigroup, the resulting predicate returns
// true only if both of the original predicates return true. This implements the associative
// operation for conjunction.
//
// Returns:
//   - A Semigroup[A] that combines predicates with AND logic
//
// Example:
//
//	s := SemigroupAll[int]()
//	isPositive := N.MoreThan(0)
//	isEven := func(n int) bool { return n%2 == 0 }
//	isPositiveAndEven := s.Concat(isPositive, isEven)
//	isPositiveAndEven(4)  // true (both)
//	isPositiveAndEven(3)  // false (not even)
//	isPositiveAndEven(-2) // false (not positive)
//	isPositiveAndEven(-3) // false (neither)
func SemigroupAll[A any]() Semigroup[A] {
	return semigroup.MakeSemigroup(func(first Predicate[A], second Predicate[A]) Predicate[A] {
		return F.Pipe1(
			first,
			And(second),
		)
	})
}

// MonoidAny creates a monoid that combines predicates using logical OR (||).
//
// This extends SemigroupAny with an identity element: a predicate that always returns false.
// The identity element satisfies the property that combining it with any predicate p yields p.
// This is useful for folding/reducing a collection of predicates where an empty collection
// should result in a predicate that always returns false.
//
// Returns:
//   - A Monoid[A] that combines predicates with OR logic and has false as identity
//
// Example:
//
//	m := MonoidAny[int]()
//	predicates := []Predicate[int]{
//	    func(n int) bool { return n > 10 },
//	    func(n int) bool { return n < 0 },
//	}
//	combined := A.Reduce(m.Empty(), m.Concat)(predicates)
//	combined(15)  // true (> 10)
//	combined(-5)  // true (< 0)
//	combined(5)   // false (neither)
func MonoidAny[A any]() Monoid[A] {
	return monoid.MakeMonoid(
		SemigroupAny[A]().Concat,
		F.Constant1[A](false),
	)
}

// MonoidAll creates a monoid that combines predicates using logical AND (&&).
//
// This extends SemigroupAll with an identity element: a predicate that always returns true.
// The identity element satisfies the property that combining it with any predicate p yields p.
// This is useful for folding/reducing a collection of predicates where an empty collection
// should result in a predicate that always returns true.
//
// Returns:
//   - A Monoid[A] that combines predicates with AND logic and has true as identity
//
// Example:
//
//	m := MonoidAll[int]()
//	predicates := []Predicate[int]{
//	    N.MoreThan(0),
//	    func(n int) bool { return n < 100 },
//	}
//	combined := A.Reduce(m.Empty(), m.Concat)(predicates)
//	combined(50)  // true (both conditions)
//	combined(-5)  // false (not > 0)
//	combined(150) // false (not < 100)
func MonoidAll[A any]() Monoid[A] {
	return monoid.MakeMonoid(
		SemigroupAll[A]().Concat,
		F.Constant1[A](true),
	)
}
