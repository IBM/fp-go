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

package eq

import (
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

func Semigroup[A any]() S.Semigroup[Eq[A]] {
	return S.MakeSemigroup(func(x, y Eq[A]) Eq[A] {
		return FromEquals(func(a, b A) bool {
			return x.Equals(a, b) && y.Equals(a, b)
		})
	})
}

func Monoid[A any]() M.Monoid[Eq[A]] {
	return M.MakeMonoid(Semigroup[A]().Concat, Empty[A]())
}
