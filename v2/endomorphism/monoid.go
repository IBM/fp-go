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

import (
	"github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Of converts any function to an Endomorphism.
//
// This function provides a way to explicitly convert a function with the signature
// func(A) A into an Endomorphism[A] type. Due to Go's type system, this is often
// not necessary as the types are compatible, but it can be useful for clarity.
//
// Parameters:
//   - f: A function from type A to type A
//
// Returns:
//   - The same function as an Endomorphism[A]
//
// Example:
//
//	myFunc := N.Mul(2)
//	endo := endomorphism.Of(myFunc)
func Of[F ~func(A) A, A any](f F) Endomorphism[A] {
	return f
}

// Wrap converts any function to an Endomorphism.
//
// Deprecated: This function is no longer needed due to Go's type compatibility.
// You can directly use functions where Endomorphism is expected.
func Wrap[F ~func(A) A, A any](f F) Endomorphism[A] {
	return f
}

// Unwrap converts any Endomorphism to a function.
//
// Deprecated: This function is no longer needed due to Go's type compatibility.
// Endomorphisms can be used directly as functions.
func Unwrap[F ~func(A) A, A any](f Endomorphism[A]) F {
	return f
}

// Identity returns the identity endomorphism.
//
// The identity endomorphism is a function that returns its input unchanged.
// It serves as the identity element for endomorphism composition, meaning:
//   - Compose(Identity(), f) = f
//   - Compose(f, Identity()) = f
//
// This is the empty element of the endomorphism monoid.
//
// Returns:
//   - An endomorphism that returns its input unchanged
//
// Example:
//
//	id := endomorphism.Identity[int]()
//	result := id(42) // Returns: 42
//
//	// Identity is neutral for composition
//	double := N.Mul(2)
//	composed := endomorphism.Compose(id, double)
//	// composed behaves exactly like double
func Identity[A any]() Endomorphism[A] {
	return function.Identity[A]
}

// Semigroup returns a Semigroup for endomorphisms where the concat operation is function composition.
//
// A semigroup is an algebraic structure with an associative binary operation.
// For endomorphisms, this operation is composition (Compose). This means:
//   - Concat(f, Concat(g, h)) = Concat(Concat(f, g), h)
//
// IMPORTANT: Concat uses Compose, which executes RIGHT-TO-LEFT:
//   - Concat(f, g) applies g first, then f
//   - This is equivalent to Compose(f, g)
//
// The returned semigroup can be used with semigroup operations to combine
// multiple endomorphisms.
//
// Returns:
//   - A Semigroup[Endomorphism[A]] where concat is composition (right-to-left)
//
// Example:
//
//	import S "github.com/IBM/fp-go/v2/semigroup"
//
//	sg := endomorphism.Semigroup[int]()
//	double := N.Mul(2)
//	increment := N.Add(1)
//
//	// Combine using the semigroup (RIGHT-TO-LEFT execution)
//	combined := sg.Concat(double, increment)
//	result := combined(5) // (5 + 1) * 2 = 12 (increment first, then double)
func Semigroup[A any]() S.Semigroup[Endomorphism[A]] {
	return S.MakeSemigroup(MonadCompose[A])
}

// Monoid returns a Monoid for endomorphisms where concat is composition and empty is identity.
//
// A monoid is a semigroup with an identity element. For endomorphisms:
//   - The binary operation is composition (Compose)
//   - The identity element is the identity function (Identity)
//
// IMPORTANT: Concat uses Compose, which executes RIGHT-TO-LEFT:
//   - Concat(f, g) applies g first, then f
//   - ConcatAll applies functions from right to left
//
// This satisfies the monoid laws:
//   - Right identity: Concat(x, Empty) = x
//   - Left identity: Concat(Empty, x) = x
//   - Associativity: Concat(x, Concat(y, z)) = Concat(Concat(x, y), z)
//
// The returned monoid can be used with monoid operations like ConcatAll to
// combine multiple endomorphisms.
//
// Returns:
//   - A Monoid[Endomorphism[A]] with composition (right-to-left) and identity
//
// Example:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	monoid := endomorphism.Monoid[int]()
//	double := N.Mul(2)
//	increment := N.Add(1)
//	square := func(x int) int { return x * x }
//
//	// Combine multiple endomorphisms (RIGHT-TO-LEFT execution)
//	combined := M.ConcatAll(monoid)(double, increment, square)
//	result := combined(5) // square(increment(double(5))) = square(increment(10)) = square(11) = 121
func Monoid[A any]() M.Monoid[Endomorphism[A]] {
	return M.MakeMonoid(MonadCompose[A], Identity[A]())
}
