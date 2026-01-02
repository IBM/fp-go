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

package endomorphism

type (
	// Endomorphism represents a function from a type to itself.
	//
	// An endomorphism is a unary function that takes a value of type A and returns
	// a value of the same type A. Mathematically, it's a function A → A.
	//
	// Endomorphisms have several important properties:
	//   - They can be composed: if f and g are endomorphisms, then f ∘ g is also an endomorphism
	//   - The identity function is an endomorphism
	//   - They form a monoid under composition
	//
	// Example:
	//
	//	// Simple endomorphisms on integers
	//	double := N.Mul(2)
	//	increment := N.Add(1)
	//
	//	// Both are endomorphisms of type Endomorphism[int]
	//	var f endomorphism.Endomorphism[int] = double
	//	var g endomorphism.Endomorphism[int] = increment
	Endomorphism[A any] = func(A) A

	// Kleisli represents a Kleisli arrow for endomorphisms.
	// It's a function from A to Endomorphism[A], used for composing endomorphic operations.
	Kleisli[A any] = func(A) Endomorphism[A]

	// Operator represents a transformation from one endomorphism to another.
	//
	// An Operator takes an endomorphism on type A and produces an endomorphism on type B.
	// This is useful for lifting operations or transforming endomorphisms in a generic way.
	//
	// Example:
	//
	//	// An operator that converts an int endomorphism to a string endomorphism
	//	intToString := func(f endomorphism.Endomorphism[int]) endomorphism.Endomorphism[string] {
	//		return func(s string) string {
	//			n, _ := strconv.Atoi(s)
	//			result := f(n)
	//			return strconv.Itoa(result)
	//		}
	//	}
	Operator[A any] = Endomorphism[Endomorphism[A]]
)
