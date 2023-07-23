// Copyright (c) 2023 IBM Corp.
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
	F "github.com/IBM/fp-go/function"
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

// Semigroup implements a two level ordering
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

// Monoid implements a two level ordering such that
// - its `Concat(ord1, ord2)` operation will order first by `ord1`, and then by `ord2`
// - its `Empty` value is an `Ord` that always considers compared elements equal
func Monoid[A any]() M.Monoid[Ord[A]] {
	return M.MakeMonoid(Semigroup[A]().Concat, FromCompare(F.Constant2[A, A](0)))
}

// MaxSemigroup returns a semigroup where `concat` will return the maximum, based on the provided order.
func MaxSemigroup[A any](O Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(Max(O))
}

// MaxSemigroup returns a semigroup where `concat` will return the minimum, based on the provided order.
func MinSemigroup[A any](O Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(Min(O))
}
