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

package ord

import (
	"cmp"
	"time"

	C "github.com/IBM/fp-go/v2/constraints"
	E "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/predicate"
)

// Ord represents a total ordering type class for type T.
// It extends Eq with a comparison function that establishes a total order.
//
// Type Parameters:
//   - T: The type for which ordering is defined
//
// Methods:
//   - Equals(x, y T) bool: Inherited from Eq, returns true if x and y are equal
//   - Compare(x, y T) int: Returns -1 if x < y, 0 if x == y, 1 if x > y
//
// Laws:
// An Ord instance must satisfy the total order laws:
//  1. Reflexivity: Compare(x, x) = 0 for all x
//  2. Antisymmetry: if Compare(x, y) <= 0 and Compare(y, x) <= 0 then x = y
//  3. Transitivity: if Compare(x, y) <= 0 and Compare(y, z) <= 0 then Compare(x, z) <= 0
//  4. Totality: Compare(x, y) <= 0 or Compare(y, x) <= 0 for all x, y
//
// The Compare function must be consistent with Equals:
//   - Compare(x, y) = 0 if and only if Equals(x, y) = true
//
// Example:
//
//	// Using built-in ordering for integers
//	intOrd := ord.FromStrictCompare[int]()
//	result := intOrd.Compare(5, 3)   // 1 (5 > 3)
//	result := intOrd.Compare(3, 5)   // -1 (3 < 5)
//	result := intOrd.Compare(5, 5)   // 0 (5 == 5)
//
//	// Creating custom ordering
//	type Person struct { Name string; Age int }
//	personOrd := ord.MakeOrd(
//	    func(p1, p2 Person) int {
//	        if p1.Age < p2.Age { return -1 }
//	        if p1.Age > p2.Age { return 1 }
//	        return 0
//	    },
//	    func(p1, p2 Person) bool { return p1.Age == p2.Age },
//	)
type Ord[T any] interface {
	E.Eq[T]
	Compare(x, y T) int
}

type ord[T any] struct {
	c func(x, y T) int
	e func(x, y T) bool
}

func (self ord[T]) Equals(x, y T) bool {
	return self.e(x, y)
}

func (self ord[T]) Compare(x, y T) int {
	return self.c(x, y)
}

// ToEq converts an [Ord] to [E.Eq].
// This allows using an Ord instance where only equality checking is needed.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	intEq := ord.ToEq(intOrd)
//	result := intEq.Equals(5, 5)  // true
//
//go:inline
func ToEq[T any](o Ord[T]) E.Eq[T] {
	return o
}

// MakeOrd creates an instance of an Ord from a compare function and an equals function.
//
// Parameters:
//   - c: A comparison function that returns -1 if x < y, 0 if x == y, 1 if x > y
//   - e: An equality function that returns true if x and y are equal
//
// The compare and equals functions must be consistent: c(x, y) == 0 iff e(x, y) == true
//
// Example:
//
//	intOrd := ord.MakeOrd(
//	    func(a, b int) int {
//	        if a < b { return -1 }
//	        if a > b { return 1 }
//	        return 0
//	    },
//	    func(a, b int) bool { return a == b },
//	)
//
//go:inline
func MakeOrd[T any](c func(x, y T) int, e func(x, y T) bool) Ord[T] {
	return ord[T]{c: c, e: e}
}

// FromCompare creates an instance of an Ord from a compare function.
// The equals function is automatically derived from the compare function.
//
// Parameters:
//   - compare: A comparison function that returns -1 if x < y, 0 if x == y, 1 if x > y
//
// Example:
//
//	stringOrd := ord.FromCompare(func(a, b string) int {
//	    if a < b { return -1 }
//	    if a > b { return 1 }
//	    return 0
//	})
func FromCompare[T any](compare func(T, T) int) Ord[T] {
	return MakeOrd(compare, func(x, y T) bool {
		return compare(x, y) == 0
	})
}

// Reverse creates an inverted ordering where the comparison results are reversed.
// If the original ordering has x < y, the reversed ordering will have x > y.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	reversedOrd := ord.Reverse(intOrd)
//	result := reversedOrd.Compare(5, 3)  // -1 (reversed from 1)
func Reverse[T any](o Ord[T]) Ord[T] {
	return MakeOrd(func(y, x T) int {
		return o.Compare(x, y)
	}, o.Equals)
}

// Contramap creates an ordering under a transformation function.
// This allows ordering values of type B by first transforming them to type A
// and then using the ordering for type A.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Parameters:
//   - f: A transformation function from B to A
//
// Returns a function that takes an Ord[A] and returns an Ord[B]
//
// Example:
//
//	type Person struct { Name string; Age int }
//	intOrd := ord.FromStrictCompare[int]()
//	personOrd := ord.Contramap(func(p Person) int {
//	    return p.Age
//	})(intOrd)
//	// Now persons are ordered by age
func Contramap[A, B any](f func(B) A) Operator[A, B] {
	return func(o Ord[A]) Ord[B] {
		return MakeOrd(func(x, y B) int {
			return o.Compare(f(x), f(y))
		}, func(x, y B) bool {
			return o.Equals(f(x), f(y))
		})
	}
}

// Min takes the minimum of two values according to the given ordering.
// If the values are considered equal, the first argument is chosen.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	min := ord.Min(intOrd)
//	result := min(5, 3)  // 3
//	result := min(5, 5)  // 5 (first argument)
func Min[A any](o Ord[A]) func(A, A) A {
	return func(a, b A) A {
		if o.Compare(a, b) < 1 {
			return a
		}
		return b
	}
}

// Max takes the maximum of two values according to the given ordering.
// If the values are considered equal, the first argument is chosen.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	max := ord.Max(intOrd)
//	result := max(5, 3)  // 5
//	result := max(5, 5)  // 5 (first argument)
func Max[A any](o Ord[A]) func(A, A) A {
	return func(a, b A) A {
		if o.Compare(a, b) >= 0 {
			return a
		}
		return b
	}
}

// Clamp restricts a value to be within a specified range [low, hi].
// If the value is less than low, low is returned.
// If the value is greater than hi, hi is returned.
// Otherwise, the value itself is returned.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	clamp := ord.Clamp(intOrd)(0, 100)
//	result := clamp(-10)  // 0
//	result := clamp(50)   // 50
//	result := clamp(150)  // 100
func Clamp[A any](o Ord[A]) func(A, A) func(A) A {
	return func(low, hi A) func(A) A {
		clow := F.Bind2nd(o.Compare, low)
		chi := F.Bind2nd(o.Compare, hi)
		return func(a A) A {
			if clow(a) <= 0 {
				return low
			}
			if chi(a) >= 0 {
				return hi
			}
			return a
		}
	}
}

//go:inline
func strictCompare[A C.Ordered](a, b A) int {
	return cmp.Compare(a, b)
}

//go:inline
func strictEq[A comparable](a, b A) bool {
	return a == b
}

// FromStrictCompare implements the ordering based on the built-in native order
// for types that satisfy the Ordered constraint (integers, floats, strings).
//
// This is the most common way to create an Ord for built-in types.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	result := intOrd.Compare(5, 3)  // 1
//
//	stringOrd := ord.FromStrictCompare[string]()
//	result := stringOrd.Compare("apple", "banana")  // -1
//
//go:inline
func FromStrictCompare[A C.Ordered]() Ord[A] {
	return MakeOrd(strictCompare[A], strictEq[A])
}

// Lt tests whether one value is strictly less than another.
// Returns a curried function that first takes the comparison value,
// then takes the value to test.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	isLessThan5 := ord.Lt(intOrd)(5)
//	result := isLessThan5(3)  // true
//	result := isLessThan5(5)  // false
//	result := isLessThan5(7)  // false
func Lt[A any](o Ord[A]) func(A) func(A) bool {
	return func(second A) func(A) bool {
		return func(first A) bool {
			return o.Compare(first, second) < 0
		}
	}
}

// Leq tests whether one value is less than or equal to another.
// Returns a curried function that first takes the comparison value,
// then takes the value to test.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	isAtMost5 := ord.Leq(intOrd)(5)
//	result := isAtMost5(3)  // true
//	result := isAtMost5(5)  // true
//	result := isAtMost5(7)  // false
func Leq[A any](o Ord[A]) func(A) func(A) bool {
	return func(second A) func(A) bool {
		return func(first A) bool {
			return o.Compare(first, second) <= 0
		}
	}
}

// Gt tests whether one value is strictly greater than another.
// Returns a curried function that first takes the comparison value,
// then takes the value to test.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	isGreaterThan5 := ord.Gt(intOrd)(5)
//	result := isGreaterThan5(3)  // false
//	result := isGreaterThan5(5)  // false
//	result := isGreaterThan5(7)  // true
func Gt[A any](o Ord[A]) func(A) func(A) bool {
	return func(second A) func(A) bool {
		return func(first A) bool {
			return o.Compare(first, second) > 0
		}
	}
}

// Geq tests whether one value is greater than or equal to another.
// Returns a curried function that first takes the comparison value,
// then takes the value to test.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	isAtLeast5 := ord.Geq(intOrd)(5)
//	result := isAtLeast5(3)  // false
//	result := isAtLeast5(5)  // true
//	result := isAtLeast5(7)  // true
func Geq[A any](o Ord[A]) func(A) func(A) bool {
	return func(second A) func(A) bool {
		return func(first A) bool {
			return o.Compare(first, second) >= 0
		}
	}
}

// Between tests whether a value is between a minimum (inclusive) and a maximum (exclusive).
// Returns a curried function that first takes the range bounds,
// then takes the value to test.
//
// The range is [lo, hi), meaning lo is included but hi is excluded.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	isBetween3And7 := ord.Between(intOrd)(3, 7)
//	result := isBetween3And7(2)  // false (below range)
//	result := isBetween3And7(3)  // true (at lower bound)
//	result := isBetween3And7(5)  // true (within range)
//	result := isBetween3And7(7)  // false (at upper bound, excluded)
//	result := isBetween3And7(8)  // false (above range)
func Between[A any](o Ord[A]) func(A, A) func(A) bool {
	lt := Lt(o)
	geq := Geq(o)
	return func(lo, hi A) func(A) bool {
		// returns the predicate
		return P.And(lt(hi))(geq(lo))
	}
}

// compareTime is a helper function that compares two time.Time values.
// Returns -1 if a is before b, 1 if a is after b, and 0 if they are equal.
func compareTime(a, b time.Time) int {
	if a.Before(b) {
		return -1
	} else if a.After(b) {
		return 1
	}
	return 0
}

// OrdTime returns an Ord instance for time.Time values.
// Times are ordered chronologically using the Before and After methods.
//
// Example:
//
//	timeOrd := ord.OrdTime()
//	t1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
//	t2 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//	result := timeOrd.Compare(t1, t2)  // -1 (t1 is before t2)
func OrdTime() Ord[time.Time] {
	return MakeOrd(compareTime, time.Time.Equal)
}
