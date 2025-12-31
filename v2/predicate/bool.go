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

// Not negates a predicate, returning a new predicate that returns the opposite boolean value.
//
// Given a predicate that returns true for some input, Not returns a predicate that returns false
// for the same input, and vice versa.
//
// Example:
//
//	isPositive := N.MoreThan(0)
//	isNotPositive := Not(isPositive)
//	isNotPositive(5)  // false
//	isNotPositive(-3) // true
func Not[A any](predicate Predicate[A]) Predicate[A] {
	return func(a A) bool {
		return !predicate(a)
	}
}

// And creates an operator that combines two predicates using logical AND (&&).
//
// The resulting predicate returns true only if both the first and second predicates return true.
// This function is curried, taking the second predicate first and returning an operator that
// takes the first predicate.
//
// Example:
//
//	isPositive := N.MoreThan(0)
//	isEven := func(n int) bool { return n%2 == 0 }
//	isPositiveAndEven := F.Pipe1(isPositive, And(isEven))
//	isPositiveAndEven(4)  // true
//	isPositiveAndEven(-2) // false
//	isPositiveAndEven(3)  // false
func And[A any](second Predicate[A]) Operator[A, A] {
	return func(first Predicate[A]) Predicate[A] {
		return func(a A) bool {
			return first(a) && second(a)
		}
	}
}

// Or creates an operator that combines two predicates using logical OR (||).
//
// The resulting predicate returns true if either the first or second predicate returns true.
// This function is curried, taking the second predicate first and returning an operator that
// takes the first predicate.
//
// Example:
//
//	isPositive := N.MoreThan(0)
//	isEven := func(n int) bool { return n%2 == 0 }
//	isPositiveOrEven := F.Pipe1(isPositive, Or(isEven))
//	isPositiveOrEven(4)  // true
//	isPositiveOrEven(-2) // true
//	isPositiveOrEven(3)  // true
//	isPositiveOrEven(-3) // false
func Or[A any](second Predicate[A]) Operator[A, A] {
	return func(first Predicate[A]) Predicate[A] {
		return func(a A) bool {
			return first(a) || second(a)
		}
	}
}
