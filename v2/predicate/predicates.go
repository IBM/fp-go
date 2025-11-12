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

package predicate

import (
	F "github.com/IBM/fp-go/v2/function"

	"github.com/IBM/fp-go/v2/eq"
)

// IsEqual creates a Kleisli arrow that tests if two values are equal using a custom equality function.
//
// This function takes an Eq instance (which defines how to compare values of type A) and returns
// a curried function that can be used to create predicates for equality testing.
//
// Parameters:
//   - pred: An Eq[A] instance that defines equality for type A
//
// Returns:
//   - A Kleisli[A, A] that takes a value and returns a predicate testing equality with that value
//
// Example:
//
//	type Person struct { Name string; Age int }
//	personEq := eq.MakeEq(func(a, b Person) bool {
//	    return a.Name == b.Name && a.Age == b.Age
//	})
//	isEqualToPerson := IsEqual(personEq)
//	alice := Person{Name: "Alice", Age: 30}
//	isAlice := isEqualToPerson(alice)
//	isAlice(Person{Name: "Alice", Age: 30}) // true
//	isAlice(Person{Name: "Bob", Age: 30})   // false
func IsEqual[A any](pred eq.Eq[A]) Kleisli[A, A] {
	return F.Curry2(pred.Equals)
}

// IsStrictEqual creates a Kleisli arrow that tests if two values are equal using Go's == operator.
//
// This is a convenience function for comparable types that uses strict equality (==) for comparison.
// It's equivalent to IsEqual with an Eq instance based on ==.
//
// Returns:
//   - A Kleisli[A, A] that takes a value and returns a predicate testing strict equality
//
// Example:
//
//	isEqualTo5 := IsStrictEqual[int]()(5)
//	isEqualTo5(5)  // true
//	isEqualTo5(10) // false
//
//	isEqualToHello := IsStrictEqual[string]()("hello")
//	isEqualToHello("hello") // true
//	isEqualToHello("world") // false
func IsStrictEqual[A comparable]() Kleisli[A, A] {
	return IsEqual(eq.FromStrictEquals[A]())
}

// IsZero creates a predicate that tests if a value equals the zero value for its type.
//
// The zero value is the default value for a type in Go (e.g., 0 for int, "" for string,
// false for bool, nil for pointers, etc.).
//
// Returns:
//   - A Predicate[A] that returns true if the value is the zero value for type A
//
// Example:
//
//	isZeroInt := IsZero[int]()
//	isZeroInt(0)  // true
//	isZeroInt(5)  // false
//
//	isZeroString := IsZero[string]()
//	isZeroString("")      // true
//	isZeroString("hello") // false
//
//	isZeroBool := IsZero[bool]()
//	isZeroBool(false) // true
//	isZeroBool(true)  // false
func IsZero[A comparable]() Predicate[A] {
	var zero A
	return IsStrictEqual[A]()(zero)
}

// IsNonZero creates a predicate that tests if a value is not equal to the zero value for its type.
//
// This is the negation of IsZero, returning true for any non-zero value.
//
// Returns:
//   - A Predicate[A] that returns true if the value is not the zero value for type A
//
// Example:
//
//	isNonZeroInt := IsNonZero[int]()
//	isNonZeroInt(0)  // false
//	isNonZeroInt(5)  // true
//	isNonZeroInt(-3) // true
//
//	isNonZeroString := IsNonZero[string]()
//	isNonZeroString("")      // false
//	isNonZeroString("hello") // true
//
//	isNonZeroPtr := IsNonZero[*int]()
//	isNonZeroPtr(nil)      // false
//	isNonZeroPtr(new(int)) // true
func IsNonZero[A comparable]() Predicate[A] {
	return Not(IsZero[A]())
}
