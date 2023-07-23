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

package monoid

import (
	S "github.com/IBM/fp-go/semigroup"
)

type Monoid[A any] interface {
	S.Semigroup[A]
	Empty() A
}

type monoid[A any] struct {
	c func(A, A) A
	e A
}

func (self monoid[A]) Concat(x A, y A) A {
	return self.c(x, y)
}

func (self monoid[A]) Empty() A {
	return self.e
}

// MakeMonoid creates a monoid given a concat function and an empty element
func MakeMonoid[A any](c func(A, A) A, e A) Monoid[A] {
	return monoid[A]{c: c, e: e}
}

// Reverse returns the dual of a `Monoid`, obtained by swapping the arguments of `Concat`.
func Reverse[A any](m Monoid[A]) Monoid[A] {
	return MakeMonoid(S.Reverse[A](m).Concat, m.Empty())
}

func ToSemigroup[A any](m Monoid[A]) S.Semigroup[A] {
	return S.Semigroup[A](m)
}
