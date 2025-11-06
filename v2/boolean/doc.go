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

// Package boolean provides functional programming utilities for working with boolean values.
//
// This package offers algebraic structures (Monoid, Eq, Ord) for boolean values,
// enabling functional composition and reasoning about boolean operations.
//
// # Monoids
//
// The package provides two monoid instances for booleans:
//
//   - MonoidAny: Combines booleans using logical OR (disjunction), with false as identity
//   - MonoidAll: Combines booleans using logical AND (conjunction), with true as identity
//
// # MonoidAny - Logical OR
//
// MonoidAny implements the boolean monoid under disjunction (OR operation).
// The identity element is false, meaning false OR x = x for any boolean x.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/boolean"
//
//	// Combine multiple booleans with OR
//	result := boolean.MonoidAny.Concat(false, true)  // true
//	result2 := boolean.MonoidAny.Concat(false, false) // false
//
//	// Identity element
//	identity := boolean.MonoidAny.Empty() // false
//
//	// Check if any value in a collection is true
//	import "github.com/IBM/fp-go/v2/array"
//	values := []bool{false, false, true, false}
//	anyTrue := array.Fold(boolean.MonoidAny)(values) // true
//
// # MonoidAll - Logical AND
//
// MonoidAll implements the boolean monoid under conjunction (AND operation).
// The identity element is true, meaning true AND x = x for any boolean x.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/boolean"
//
//	// Combine multiple booleans with AND
//	result := boolean.MonoidAll.Concat(true, true)   // true
//	result2 := boolean.MonoidAll.Concat(true, false) // false
//
//	// Identity element
//	identity := boolean.MonoidAll.Empty() // true
//
//	// Check if all values in a collection are true
//	import "github.com/IBM/fp-go/v2/array"
//	values := []bool{true, true, true}
//	allTrue := array.Fold(boolean.MonoidAll)(values) // true
//
// # Equality
//
// The Eq instance provides structural equality for booleans:
//
//	import "github.com/IBM/fp-go/v2/boolean"
//
//	equal := boolean.Eq.Equals(true, true)   // true
//	equal2 := boolean.Eq.Equals(true, false) // false
//
// # Ordering
//
// The Ord instance provides a total ordering for booleans where false < true:
//
//	import "github.com/IBM/fp-go/v2/boolean"
//
//	cmp := boolean.Ord.Compare(false, true) // -1 (false < true)
//	cmp2 := boolean.Ord.Compare(true, false) // +1 (true > false)
//	cmp3 := boolean.Ord.Compare(true, true)  // 0 (equal)
//
// # Use Cases
//
// The boolean package is particularly useful for:
//
//   - Combining multiple boolean conditions functionally
//   - Implementing validation logic that accumulates results
//   - Working with predicates in a composable way
//   - Folding collections of boolean values
//
// Example - Validation:
//
//	import (
//	    "github.com/IBM/fp-go/v2/array"
//	    "github.com/IBM/fp-go/v2/boolean"
//	)
//
//	type User struct {
//	    Name  string
//	    Email string
//	    Age   int
//	}
//
//	// Define validation predicates
//	validations := []func(User) bool{
//	    func(u User) bool { return len(u.Name) > 0 },
//	    func(u User) bool { return len(u.Email) > 0 },
//	    func(u User) bool { return u.Age >= 18 },
//	}
//
//	// Check if user passes all validations
//	user := User{"Alice", "alice@example.com", 25}
//	results := array.Map(func(v func(User) bool) bool {
//	    return v(user)
//	})(validations)
//	allValid := array.Fold(boolean.MonoidAll)(results) // true
//
// Example - Any Match:
//
//	import (
//	    "github.com/IBM/fp-go/v2/array"
//	    "github.com/IBM/fp-go/v2/boolean"
//	)
//
//	// Check if any number is even
//	numbers := []int{1, 3, 5, 7, 8, 9}
//	checks := array.Map(func(n int) bool {
//	    return n%2 == 0
//	})(numbers)
//	hasEven := array.Fold(boolean.MonoidAny)(checks) // true
package boolean
