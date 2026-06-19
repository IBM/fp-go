// Copyright (c) 2025 IBM Corp.
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

// Package predicate provides functional programming utilities for working with predicates.
//
// A predicate is a function that takes a value and returns a boolean, commonly used
// for filtering, validation, and conditional logic. This package offers combinators
// for composing predicates using logical operations (And, Or, Not), transforming
// predicates via ContraMap, and combining multiple predicates using Semigroup and
// Monoid abstractions.
//
// Key features:
//   - Boolean combinators: And, Or, Not
//   - ContraMap for transforming predicates
//   - Semigroup and Monoid instances for combining predicates
//
// Example usage:
//
//	import P "github.com/IBM/fp-go/v2/predicate"
//
//	// Create predicates
//	isPositive := N.MoreThan(0)
//	isEven := func(n int) bool { return n%2 == 0 }
//
//	// Combine predicates
//	isPositiveAndEven := F.Pipe1(isPositive, P.And(isEven))
//	isPositiveOrEven := F.Pipe1(isPositive, P.Or(isEven))
//	isNotPositive := P.Not(isPositive)
package predicate

type (
	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It is commonly used for filtering, validation, and conditional logic.
	Predicate[A any] = func(A) bool

	// Kleisli represents a function that takes a value of type A and returns a Predicate[B].
	// This is a Kleisli arrow in the context of predicates, allowing for the creation of
	// parameterized predicates. It's particularly useful for building predicates that depend
	// on some input value, such as equality testing or comparison operations.
	//
	// Type Parameters:
	//   - A: The input type for the Kleisli arrow
	//   - B: The type that the resulting predicate will test
	//
	// Common uses:
	//   - IsEqual: Takes a value and returns a predicate testing equality with that value
	//   - IsStrictEqual: Takes a value and returns a predicate testing strict equality
	//   - Custom parameterized predicates that depend on configuration or context
	//
	// See Also:
	//   - IsEqual: Returns a Kleisli[A, A] for custom equality testing
	//   - IsStrictEqual: Returns a Kleisli[A, A] for strict equality testing
	//   - Operator: A specialized Kleisli for transforming predicates
	Kleisli[A, B any] = func(A) Predicate[B]

	// Operator represents a function that transforms a Predicate[A] into a Predicate[B].
	// This is useful for composing and transforming predicates.
	Operator[A, B any] = Kleisli[Predicate[A], B]
)
