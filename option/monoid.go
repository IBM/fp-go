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

package option

import (
	F "github.com/IBM/fp-go/function"
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
)

func Semigroup[A any]() func(S.Semigroup[A]) S.Semigroup[Option[A]] {
	return func(s S.Semigroup[A]) S.Semigroup[Option[A]] {
		concat := s.Concat
		return S.MakeSemigroup(
			func(x, y Option[A]) Option[A] {
				return MonadFold(x, F.Constant(y), func(left A) Option[A] {
					return MonadFold(y, F.Constant(x), func(right A) Option[A] {
						return Some(concat(left, right))
					})
				})
			},
		)
	}
}

// Monoid returning the left-most non-`None` value. If both operands are `Some`s then the inner values are
// concatenated using the provided `Semigroup`
//
// | x       | y       | concat(x, y)       |
// | ------- | ------- | ------------------ |
// | none    | none    | none               |
// | some(a) | none    | some(a)            |
// | none    | some(b) | some(b)            |
// | some(a) | some(b) | some(concat(a, b)) |
func Monoid[A any]() func(S.Semigroup[A]) M.Monoid[Option[A]] {
	sg := Semigroup[A]()
	return func(s S.Semigroup[A]) M.Monoid[Option[A]] {
		return M.MakeMonoid(sg(s).Concat, None[A]())
	}
}
