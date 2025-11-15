// Copyright (c) 2025 IBM Corp.
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

package option

import (
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// ApplySemigroup lifts a Semigroup over a type A to a Semigroup over Option[A].
// The resulting semigroup combines two Options using the applicative functor pattern.
//
// Example:
//
//	intSemigroup := semigroup.MakeSemigroup(func(a, b int) int { return a + b })
//	optSemigroup := ApplySemigroup(intSemigroup)
//	result := optSemigroup.Concat(Some(2), Some(3)) // Some(5)
//	result := optSemigroup.Concat(Some(2), None[int]()) // None
func ApplySemigroup[A any](s S.Semigroup[A]) S.Semigroup[Option[A]] {
	return S.ApplySemigroup(MonadMap[A, func(A) A], MonadAp[A, A], s)
}

// ApplicativeMonoid returns a Monoid that concatenates Option instances via their applicative functor.
// This combines the monoid structure of the underlying type with the Option structure.
//
// Example:
//
//	intMonoid := monoid.MakeMonoid(func(a, b int) int { return a + b }, 0)
//	optMonoid := ApplicativeMonoid(intMonoid)
//	result := optMonoid.Concat(Some(2), Some(3)) // Some(5)
//	result := optMonoid.Empty() // Some(0)
//
//go:inline
func ApplicativeMonoid[A any](m M.Monoid[A]) M.Monoid[Option[A]] {
	return M.ApplicativeMonoid(Of[A], MonadMap[A, func(A) A], MonadAp[A, A], m)
}
