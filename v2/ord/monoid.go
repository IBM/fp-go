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
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Semigroup implements a two-level ordering that combines two Ord instances.
// The resulting Ord will first compare using the first ordering, and only if
// the values are equal according to the first ordering, it will use the second ordering.
//
// This is useful for implementing multi-level sorting (e.g., sort by last name, then by first name).
//
// Example:
//
//	type Person struct { LastName, FirstName string }
//	stringOrd := ord.FromStrictCompare[string]()
//	byLastName := ord.Contramap(func(p Person) string { return p.LastName })(stringOrd)
//	byFirstName := ord.Contramap(func(p Person) string { return p.FirstName })(stringOrd)
//	sg := ord.Semigroup[Person]()
//	personOrd := sg.Concat(byLastName, byFirstName)
//	// Now persons are ordered by last name, then by first name
func Semigroup[A any]() S.Semigroup[Ord[A]] {
	return S.MakeSemigroup(func(first, second Ord[A]) Ord[A] {
		return FromCompare(func(a, b A) int {
			ox := first.Compare(a, b)
			if ox != 0 {
				return ox
			}
			return second.Compare(a, b)
		})
	})
}

// Monoid implements a two-level ordering with an identity element.
//
// Properties:
//   - Concat(ord1, ord2) will order first by ord1, and then by ord2
//   - Empty() returns an Ord that always considers compared elements equal
//
// The Empty ordering acts as an identity: Concat(ord, Empty()) == ord
//
// Example:
//
//	m := ord.Monoid[int]()
//	emptyOrd := m.Empty()
//	result := emptyOrd.Compare(5, 3)  // 0 (always equal)
//
//	intOrd := ord.FromStrictCompare[int]()
//	combined := m.Concat(intOrd, emptyOrd)  // same as intOrd
func Monoid[A any]() M.Monoid[Ord[A]] {
	return M.MakeMonoid(Semigroup[A]().Concat, FromCompare(F.Constant2[A, A](0)))
}

// MaxSemigroup returns a semigroup where Concat will return the maximum value
// according to the provided ordering.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	maxSg := ord.MaxSemigroup(intOrd)
//	result := maxSg.Concat(5, 3)  // 5
//	result := maxSg.Concat(3, 5)  // 5
func MaxSemigroup[A any](o Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(Max(o))
}

// MinSemigroup returns a semigroup where Concat will return the minimum value
// according to the provided ordering.
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	minSg := ord.MinSemigroup(intOrd)
//	result := minSg.Concat(5, 3)  // 3
//	result := minSg.Concat(3, 5)  // 3
func MinSemigroup[A any](o Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(Min(o))
}
