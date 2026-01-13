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

package monoid

import (
	F "github.com/IBM/fp-go/v2/function"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// FunctionMonoid creates a monoid for functions when the codomain (return type) has a monoid.
//
// Given a monoid for type B, this creates a monoid for functions of type func(A) B.
// The resulting monoid combines functions by combining their results using the provided
// monoid's Concat operation. The empty function always returns the monoid's Empty value.
//
// This allows you to compose functions that return monoidal values in a point-wise manner.
//
// Type Parameters:
//   - A: The domain (input type) of the functions
//   - B: The codomain (return type) of the functions, which must have a monoid
//
// Parameters:
//   - m: A monoid for the codomain type B
//
// Returns:
//   - A Monoid[func(A) B] that combines functions point-wise
//
// Example:
//
//	// Monoid for functions returning integers
//	intAddMonoid := MakeMonoid(
//	    func(a, b int) int { return a + b },
//	    0,
//	)
//
//	funcMonoid := FunctionMonoid[string, int](intAddMonoid)
//
//	// Define some functions
//	f1 := S.Size
//	f2 := func(s string) int { return len(s) * 2 }
//
//	// Combine functions: result(x) = f1(x) + f2(x)
//	combined := funcMonoid.Concat(f1, f2)
//	result := combined("hello")  // len("hello") + len("hello")*2 = 5 + 10 = 15
//
//	// Empty function always returns 0
//	emptyFunc := funcMonoid.Empty()
//	result := emptyFunc("anything")  // 0
//
//	// Verify identity laws
//	assert.Equal(t, f1("test"), funcMonoid.Concat(funcMonoid.Empty(), f1)("test"))
//	assert.Equal(t, f1("test"), funcMonoid.Concat(f1, funcMonoid.Empty())("test"))
func FunctionMonoid[A, B any](m Monoid[B]) Monoid[func(A) B] {
	return MakeMonoid(
		S.FunctionSemigroup[A](m).Concat,
		F.Constant1[A](m.Empty()),
	)
}
