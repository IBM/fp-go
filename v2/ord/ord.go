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
	C "github.com/IBM/fp-go/v2/constraints"
	E "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/predicate"
)

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

// ToEq converts an [Ord] to [E.Eq]
func ToEq[T any](o Ord[T]) E.Eq[T] {
	return o
}

// MakeOrd creates an instance of an Ord
func MakeOrd[T any](c func(x, y T) int, e func(x, y T) bool) Ord[T] {
	return ord[T]{c: c, e: e}
}

// MakeOrd creates an instance of an Ord from a compare function
func FromCompare[T any](compare func(T, T) int) Ord[T] {
	return MakeOrd(compare, func(x, y T) bool {
		return compare(x, y) == 0
	})
}

// Reverse creates an inverted ordering
func Reverse[T any](o Ord[T]) Ord[T] {
	return MakeOrd(func(y, x T) int {
		return o.Compare(x, y)
	}, o.Equals)
}

// Contramap creates an ordering under a transformation function
func Contramap[A, B any](f func(B) A) func(Ord[A]) Ord[B] {
	return func(o Ord[A]) Ord[B] {
		return MakeOrd(func(x, y B) int {
			return o.Compare(f(x), f(y))
		}, func(x, y B) bool {
			return o.Equals(f(x), f(y))
		})
	}
}

// Min takes the minimum of two values. If they are considered equal, the first argument is chosen
func Min[A any](o Ord[A]) func(A, A) A {
	return func(a, b A) A {
		if o.Compare(a, b) < 1 {
			return a
		}
		return b
	}
}

// Max takes the maximum of two values. If they are considered equal, the first argument is chosen
func Max[A any](o Ord[A]) func(A, A) A {
	return func(a, b A) A {
		if o.Compare(a, b) >= 0 {
			return a
		}
		return b
	}
}

// Clamp clamps a value between a minimum and a maximum
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

func strictCompare[A C.Ordered](a, b A) int {
	if a < b {
		return -1
	} else if a > b {
		return +1
	} else {
		return 0
	}
}

func strictEq[A comparable](a, b A) bool {
	return a == b
}

// FromStrictCompare implements the ordering based on the built in native order
func FromStrictCompare[A C.Ordered]() Ord[A] {
	return MakeOrd(strictCompare[A], strictEq[A])
}

// Lt tests whether one value is strictly less than another
func Lt[A any](o Ord[A]) func(A) func(A) bool {
	return func(second A) func(A) bool {
		return func(first A) bool {
			return o.Compare(first, second) < 0
		}
	}
}

// Leq Tests whether one value is less or equal than another
func Leq[A any](O Ord[A]) func(A) func(A) bool {
	return func(second A) func(A) bool {
		return func(first A) bool {
			return O.Compare(first, second) <= 0
		}
	}
}

/**
 * Test whether one value is strictly greater than another
 */
func Gt[A any](O Ord[A]) func(A) func(A) bool {
	return func(second A) func(A) bool {
		return func(first A) bool {
			return O.Compare(first, second) > 0
		}
	}
}

// Geq tests whether one value is greater or equal than another
func Geq[A any](O Ord[A]) func(A) func(A) bool {
	return func(second A) func(A) bool {
		return func(first A) bool {
			return O.Compare(first, second) >= 0
		}
	}
}

// Between tests whether a value is between a minimum (inclusive) and a maximum (exclusive)
func Between[A any](O Ord[A]) func(A, A) func(A) bool {
	lt := Lt(O)
	geq := Geq(O)
	return func(lo, hi A) func(A) bool {
		// returns the predicate
		return P.And(lt(hi))(geq(lo))
	}
}
