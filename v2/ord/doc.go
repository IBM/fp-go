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
Package ord provides the Ord type class for types that support total ordering.

# Overview

An Ord represents a total ordering on a type. It extends Eq (equality) with
a comparison function that returns -1, 0, or 1 to indicate less than, equal to,
or greater than relationships.

The Ord interface:

	type Ord[T any] interface {
	    Eq[T]                    // Provides Equals(x, y T) bool
	    Compare(x, y T) int      // Returns -1, 0, or 1
	}

Ord laws:
  - Reflexivity: Compare(x, x) = 0
  - Antisymmetry: if Compare(x, y) <= 0 and Compare(y, x) <= 0 then x = y
  - Transitivity: if Compare(x, y) <= 0 and Compare(y, z) <= 0 then Compare(x, z) <= 0
  - Totality: Compare(x, y) <= 0 or Compare(y, x) <= 0

# Basic Usage

Creating an Ord for integers:

	intOrd := ord.FromStrictCompare[int]()

	result := intOrd.Compare(5, 3)   // 1 (5 > 3)
	result := intOrd.Compare(3, 5)   // -1 (3 < 5)
	result := intOrd.Compare(5, 5)   // 0 (5 == 5)

	equal := intOrd.Equals(5, 5)     // true

Creating a custom Ord:

	type Person struct {
	    Name string
	    Age  int
	}

	// Order by age
	personOrd := ord.MakeOrd(
	    func(p1, p2 Person) int {
	        if p1.Age < p2.Age {
	            return -1
	        } else if p1.Age > p2.Age {
	            return 1
	        }
	        return 0
	    },
	    func(p1, p2 Person) bool {
	        return p1.Age == p2.Age
	    },
	)

Creating Ord from compare function only:

	stringOrd := ord.FromCompare(func(a, b string) int {
	    if a < b {
	        return -1
	    } else if a > b {
	        return 1
	    }
	    return 0
	})
	// Equals is automatically derived from Compare

# Comparison Functions

Lt - Less than:

	intOrd := ord.FromStrictCompare[int]()
	isLessThan5 := ord.Lt(intOrd)(5)

	result := isLessThan5(3)  // true
	result := isLessThan5(5)  // false
	result := isLessThan5(7)  // false

Leq - Less than or equal:

	isAtMost5 := ord.Leq(intOrd)(5)

	result := isAtMost5(3)  // true
	result := isAtMost5(5)  // true
	result := isAtMost5(7)  // false

Gt - Greater than:

	isGreaterThan5 := ord.Gt(intOrd)(5)

	result := isGreaterThan5(3)  // false
	result := isGreaterThan5(5)  // false
	result := isGreaterThan5(7)  // true

Geq - Greater than or equal:

	isAtLeast5 := ord.Geq(intOrd)(5)

	result := isAtLeast5(3)  // false
	result := isAtLeast5(5)  // true
	result := isAtLeast5(7)  // true

Between - Check if value is in range [low, high):

	intOrd := ord.FromStrictCompare[int]()
	isBetween3And7 := ord.Between(intOrd)(3, 7)

	result := isBetween3And7(2)  // false
	result := isBetween3And7(3)  // true
	result := isBetween3And7(5)  // true
	result := isBetween3And7(7)  // false
	result := isBetween3And7(8)  // false

# Min and Max

Min - Get the minimum of two values:

	intOrd := ord.FromStrictCompare[int]()
	min := ord.Min(intOrd)

	result := min(5, 3)  // 3
	result := min(3, 5)  // 3
	result := min(5, 5)  // 5 (first argument when equal)

Max - Get the maximum of two values:

	max := ord.Max(intOrd)

	result := max(5, 3)  // 5
	result := max(3, 5)  // 5
	result := max(5, 5)  // 5 (first argument when equal)

Clamp - Restrict a value to a range:

	intOrd := ord.FromStrictCompare[int]()
	clamp := ord.Clamp(intOrd)(0, 100)

	result := clamp(-10)  // 0 (clamped to minimum)
	result := clamp(50)   // 50 (within range)
	result := clamp(150)  // 100 (clamped to maximum)

# Transforming Ord

Reverse - Invert the ordering:

	intOrd := ord.FromStrictCompare[int]()
	reversedOrd := ord.Reverse(intOrd)

	result := intOrd.Compare(5, 3)         // 1 (5 > 3)
	result := reversedOrd.Compare(5, 3)    // -1 (3 > 5 in reversed order)

	min := ord.Min(reversedOrd)
	result := min(5, 3)  // 5 (max in original order)

Contramap - Transform the input before comparing:

	type Person struct {
	    Name string
	    Age  int
	}

	intOrd := ord.FromStrictCompare[int]()

	// Order persons by age
	personOrd := ord.Contramap(func(p Person) int {
	    return p.Age
	})(intOrd)

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}

	result := personOrd.Compare(p1, p2)  // 1 (30 > 25)

ToEq - Convert Ord to Eq:

	intOrd := ord.FromStrictCompare[int]()
	intEq := ord.ToEq(intOrd)

	result := intEq.Equals(5, 5)  // true
	result := intEq.Equals(5, 3)  // false

# Semigroup and Monoid

Semigroup - Combine orderings (try first, then second):

	type Person struct {
	    LastName  string
	    FirstName string
	}

	stringOrd := ord.FromStrictCompare[string]()

	// Order by last name
	byLastName := ord.Contramap(func(p Person) string {
	    return p.LastName
	})(stringOrd)

	// Order by first name
	byFirstName := ord.Contramap(func(p Person) string {
	    return p.FirstName
	})(stringOrd)

	// Combine: order by last name, then first name
	sg := ord.Semigroup[Person]()
	personOrd := sg.Concat(byLastName, byFirstName)

	p1 := Person{LastName: "Smith", FirstName: "Alice"}
	p2 := Person{LastName: "Smith", FirstName: "Bob"}

	result := personOrd.Compare(p1, p2)  // -1 (Alice < Bob)

Monoid - Semigroup with identity (always equal):

	m := ord.Monoid[int]()

	// Empty ordering considers everything equal
	emptyOrd := m.Empty()
	result := emptyOrd.Compare(5, 3)  // 0 (always equal)

	// Concat with empty returns the original
	intOrd := ord.FromStrictCompare[int]()
	result := m.Concat(intOrd, emptyOrd)  // same as intOrd

MaxSemigroup - Semigroup that returns maximum:

	intOrd := ord.FromStrictCompare[int]()
	maxSg := ord.MaxSemigroup(intOrd)

	result := maxSg.Concat(5, 3)  // 5
	result := maxSg.Concat(3, 5)  // 5

MinSemigroup - Semigroup that returns minimum:

	minSg := ord.MinSemigroup(intOrd)

	result := minSg.Concat(5, 3)  // 3
	result := minSg.Concat(3, 5)  // 3

# Practical Examples

Sorting with custom order:

	import (
	    "sort"
	    O "github.com/IBM/fp-go/v2/ord"
	)

	type Person struct {
	    Name string
	    Age  int
	}

	people := []Person{
	    {Name: "Alice", Age: 30},
	    {Name: "Bob", Age: 25},
	    {Name: "Charlie", Age: 35},
	}

	intOrd := O.FromStrictCompare[int]()
	personOrd := O.Contramap(func(p Person) int {
	    return p.Age
	})(intOrd)

	sort.Slice(people, func(i, j int) bool {
	    return personOrd.Compare(people[i], people[j]) < 0
	})
	// people is now sorted by age

Finding min/max in a slice:

	import (
	    A "github.com/IBM/fp-go/v2/array"
	    O "github.com/IBM/fp-go/v2/ord"
	)

	numbers := []int{5, 2, 8, 1, 9, 3}
	intOrd := O.FromStrictCompare[int]()

	min := O.Min(intOrd)
	max := O.Max(intOrd)

	// Find minimum
	minimum := A.Reduce(numbers, min, numbers[0])  // 1

	// Find maximum
	maximum := A.Reduce(numbers, max, numbers[0])  // 9

Multi-level sorting:

	type Employee struct {
	    Department string
	    Name       string
	    Salary     int
	}

	stringOrd := O.FromStrictCompare[string]()
	intOrd := O.FromStrictCompare[int]()

	// Order by department
	byDept := O.Contramap(func(e Employee) string {
	    return e.Department
	})(stringOrd)

	// Order by salary (descending)
	bySalary := O.Reverse(O.Contramap(func(e Employee) int {
	    return e.Salary
	})(intOrd))

	// Order by name
	byName := O.Contramap(func(e Employee) string {
	    return e.Name
	})(stringOrd)

	// Combine: dept, then salary (desc), then name
	sg := O.Semigroup[Employee]()
	employeeOrd := sg.Concat(sg.Concat(byDept, bySalary), byName)

Filtering with comparisons:

	import (
	    A "github.com/IBM/fp-go/v2/array"
	    O "github.com/IBM/fp-go/v2/ord"
	)

	numbers := []int{1, 5, 3, 8, 2, 9, 4}
	intOrd := O.FromStrictCompare[int]()

	// Filter numbers greater than 5
	gt5 := O.Gt(intOrd)(5)
	result := A.Filter(gt5)(numbers)  // [8, 9]

	// Filter numbers between 3 and 7
	between3And7 := O.Between(intOrd)(3, 7)
	result := A.Filter(between3And7)(numbers)  // [5, 3, 4]

# Functions

Core operations:
  - MakeOrd[T any](func(T, T) int, func(T, T) bool) Ord[T] - Create Ord from compare and equals
  - FromCompare[T any](func(T, T) int) Ord[T] - Create Ord from compare (derives equals)
  - FromStrictCompare[A Ordered]() Ord[A] - Create Ord for built-in ordered types
  - ToEq[T any](Ord[T]) Eq[T] - Convert Ord to Eq

Transformations:
  - Reverse[T any](Ord[T]) Ord[T] - Invert the ordering
  - Contramap[A, B any](func(B) A) func(Ord[A]) Ord[B] - Transform input before comparing

Comparisons:
  - Lt[A any](Ord[A]) func(A) func(A) bool - Less than
  - Leq[A any](Ord[A]) func(A) func(A) bool - Less than or equal
  - Gt[A any](Ord[A]) func(A) func(A) bool - Greater than
  - Geq[A any](Ord[A]) func(A) func(A) bool - Greater than or equal
  - Between[A any](Ord[A]) func(A, A) func(A) bool - Check if in range [low, high)

Min/Max/Clamp:
  - Min[A any](Ord[A]) func(A, A) A - Get minimum of two values
  - Max[A any](Ord[A]) func(A, A) A - Get maximum of two values
  - Clamp[A any](Ord[A]) func(A, A) func(A) A - Clamp value to range

Algebraic structures:
  - Semigroup[A any]() Semigroup[Ord[A]] - Combine orderings
  - Monoid[A any]() Monoid[Ord[A]] - Semigroup with identity (always equal)
  - MaxSemigroup[A any](Ord[A]) Semigroup[A] - Semigroup returning maximum
  - MinSemigroup[A any](Ord[A]) Semigroup[A] - Semigroup returning minimum

# Related Packages

  - eq: Equality type class (parent of Ord)
  - constraints: Type constraints for generics
  - semigroup: Associative binary operation
  - monoid: Semigroup with identity element
*/
package ord
