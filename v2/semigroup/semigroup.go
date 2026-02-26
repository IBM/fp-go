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

package semigroup

import (
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/magma"
)

// Semigroup represents an algebraic structure with an associative binary operation.
// It extends the Magma interface by requiring that the Concat operation be associative,
// meaning (a • b) • c = a • (b • c) for all values a, b, c.
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//	sum := N.SemigroupSum[int]()
//	result := sum.Concat(sum.Concat(1, 2), 3)  // Same as sum.Concat(1, sum.Concat(2, 3))
type Semigroup[A any] interface {
	M.Magma[A]
}

type semigroup[A any] struct {
	c func(A, A) A
}

func (self semigroup[A]) Concat(x, y A) A {
	return self.c(x, y)
}

// MakeSemigroup creates a Semigroup from a binary operation function.
// The provided function must be associative to form a valid semigroup.
//
// Example:
//
//	// Create a string concatenation semigroup
//	strConcat := semigroup.MakeSemigroup(func(a, b string) string {
//	    return a + b
//	})
//	result := strConcat.Concat("Hello, ", "World!")  // "Hello, World!"
func MakeSemigroup[A any](c func(A, A) A) Semigroup[A] {
	return semigroup[A]{c: c}
}

// Reverse returns the dual of a Semigroup, obtained by swapping the arguments of Concat.
// If the original operation is a • b, the reversed operation is b • a.
//
// Example:
//
//	sub := semigroup.MakeSemigroup(func(a, b int) int { return a - b })
//	reversed := semigroup.Reverse(sub)
//	result1 := sub.Concat(10, 3)      // 10 - 3 = 7
//	result2 := reversed.Concat(10, 3) // 3 - 10 = -7
func Reverse[A any](m Semigroup[A]) Semigroup[A] {
	return MakeSemigroup(M.Reverse(m).Concat)
}

// FunctionSemigroup lifts a Semigroup to work with functions that return values in that semigroup.
// Given a Semigroup[B], it creates a Semigroup[func(A) B] where functions are combined by
// applying both functions to the same input and combining the results using the original semigroup.
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//	intSum := N.SemigroupSum[int]()
//	funcSG := semigroup.FunctionSemigroup[string](intSum)
//
//	f := S.Size
//	g := func(s string) int { return len(s) * 2 }
//	combined := funcSG.Concat(f, g)
//	result := combined("hello")  // 5 + 10 = 15
func FunctionSemigroup[A, B any](s Semigroup[B]) Semigroup[func(A) B] {
	return MakeSemigroup(func(f func(A) B, g func(A) B) func(A) B {
		return func(a A) B {
			return s.Concat(f(a), g(a))
		}
	})
}

// First creates a Semigroup that always returns the first argument.
// This is useful when you want a semigroup where earlier values take precedence.
//
// Example:
//
//	first := semigroup.First[int]()
//	result := first.Concat(1, 2)  // Returns: 1
func First[A any]() Semigroup[A] {
	return MakeSemigroup(F.First[A, A])
}

// Last creates a Semigroup that always returns the last argument.
// This is useful when you want a semigroup where later values take precedence.
//
// Example:
//
//	last := semigroup.Last[int]()
//	result := last.Concat(1, 2)  // Returns: 2
func Last[A any]() Semigroup[A] {
	return MakeSemigroup(F.Second[A, A])
}

// ToMagma converts a Semigroup to a Magma.
// Since Semigroup extends Magma, this is simply an identity conversion that
// changes the type perspective without modifying the underlying structure.
//
// Example:
//
//	sg := semigroup.First[int]()
//	magma := semigroup.ToMagma(sg)
func ToMagma[A any](s Semigroup[A]) M.Magma[A] {
	return s
}

// ConcatWith creates a curried version of the Concat operation with the left argument fixed first.
// It returns a function that takes the left operand and returns another function that takes
// the right operand and performs the concatenation.
//
// This is useful for partial application and function composition patterns.
//
// # Type Parameters
//
//   - A: The type of elements in the semigroup
//
// # Parameters
//
//   - s: The semigroup to use for concatenation
//
// # Returns
//
//   - func(A) func(A) A: A curried function that takes left then right operand
//
// # Example Usage
//
//	import N "github.com/IBM/fp-go/v2/number"
//	sum := N.SemigroupSum[int]()
//	concatWith := ConcatWith(sum)
//	add5 := concatWith(5)
//	result := add5(3)  // 5 + 3 = 8
//
// # See Also
//
//   - AppendTo: Similar but fixes the right argument first
func ConcatWith[A any](s Semigroup[A]) func(A) func(A) A {
	return func(l A) func(A) A {
		return func(r A) A {
			return s.Concat(l, r)
		}
	}
}

// AppendTo creates a curried version of the Concat operation with the right argument fixed first.
// It returns a function that takes the right operand and returns another function that takes
// the left operand and performs the concatenation.
//
// This is useful for partial application where you want to fix the second argument first,
// which is common in append-style operations.
//
// # Type Parameters
//
//   - A: The type of elements in the semigroup
//
// # Parameters
//
//   - s: The semigroup to use for concatenation
//
// # Returns
//
//   - func(A) func(A) A: A curried function that takes right then left operand
//
// # Example Usage
//
//	import S "github.com/IBM/fp-go/v2/string"
//	strConcat := S.Semigroup
//	appendTo := AppendTo(strConcat)
//	addSuffix := appendTo("!")
//	result := addSuffix("Hello")  // "Hello" + "!" = "Hello!"
//
// # See Also
//
//   - ConcatWith: Similar but fixes the left argument first
func AppendTo[A any](s Semigroup[A]) func(A) func(A) A {
	return func(r A) func(A) A {
		return func(l A) A {
			return s.Concat(l, r)
		}
	}
}
