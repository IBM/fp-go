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
	S "github.com/IBM/fp-go/v2/semigroup"
)

// FromSemigroup converts a semigroup into a Kleisli arrow for endomorphisms.
//
// This function takes a semigroup and returns a Kleisli arrow that, when given
// a value of type A, produces an endomorphism that concatenates that value with
// other values using the semigroup's Concat operation.
//
// The resulting Kleisli arrow has the signature: func(A) Endomorphism[A]
// When called with a value 'x', it returns an endomorphism that concatenates
// 'x' with its input using the semigroup's binary operation.
//
// # Data Last Principle
//
// FromSemigroup follows the "data last" principle by using function.Bind2of2,
// which binds the second parameter of the semigroup's Concat operation.
// This means that for a semigroup with Concat(a, b), calling FromSemigroup(s)(x)
// creates an endomorphism that computes Concat(input, x), where the input data
// comes first and the bound value 'x' comes last.
//
// For example, with string concatenation:
//   - Semigroup.Concat("Hello", "World") = "HelloWorld"
//   - FromSemigroup(semigroup)("World") creates: func(input) = Concat(input, "World")
//   - Applying it: endomorphism("Hello") = Concat("Hello", "World") = "HelloWorld"
//
// This is particularly useful for creating endomorphisms from associative operations
// like string concatenation, number addition, list concatenation, etc.
//
// Parameters:
//   - s: A semigroup providing the Concat operation for type A
//
// Returns:
//   - A Kleisli arrow that converts values of type A into endomorphisms
//
// Example:
//
//	import (
//		"github.com/IBM/fp-go/v2/endomorphism"
//		"github.com/IBM/fp-go/v2/semigroup"
//	)
//
//	// Create a semigroup for integer addition
//	addSemigroup := semigroup.MakeSemigroup(func(a, b int) int {
//		return a + b
//	})
//
//	// Convert it to a Kleisli arrow
//	addKleisli := endomorphism.FromSemigroup(addSemigroup)
//
//	// Use the Kleisli arrow to create an endomorphism that adds 5
//	// This follows "data last": the input data comes first, 5 comes last
//	addFive := addKleisli(5)
//
//	// Apply the endomorphism: Concat(10, 5) = 10 + 5 = 15
//	result := addFive(10) // result is 15
//
// The function uses function.Bind2of2 to partially apply the semigroup's Concat
// operation, effectively currying it to create the desired Kleisli arrow while
// maintaining the "data last" principle.
func FromSemigroup[A any](s S.Semigroup[A]) Kleisli[A] {
	return function.Bind2of2(s.Concat)
}
