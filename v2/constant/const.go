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

package constant

import (
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

type Const[E, A any] struct {
	value E
}

func Make[E, A any](e E) Const[E, A] {
	return Const[E, A]{value: e}
}

func Unwrap[E, A any](c Const[E, A]) E {
	return c.value
}

func Of[E, A any](m M.Monoid[E]) func(A) Const[E, A] {
	return F.Constant1[A](Make[E, A](m.Empty()))
}

func MonadMap[E, A, B any](fa Const[E, A], _ func(A) B) Const[E, B] {
	return Make[E, B](fa.value)
}

func MonadAp[E, A, B any](s S.Semigroup[E]) func(fab Const[E, func(A) B], fa Const[E, A]) Const[E, B] {
	return func(fab Const[E, func(A) B], fa Const[E, A]) Const[E, B] {
		return Make[E, B](s.Concat(fab.value, fa.value))
	}
}

func Map[E, A, B any](f func(A) B) func(fa Const[E, A]) Const[E, B] {
	return F.Bind2nd(MonadMap[E, A, B], f)
}

func Ap[E, A, B any](s S.Semigroup[E]) func(fa Const[E, A]) func(fab Const[E, func(A) B]) Const[E, B] {
	monadap := MonadAp[E, A, B](s)
	return func(fa Const[E, A]) func(fab Const[E, func(A) B]) Const[E, B] {
		return F.Bind2nd(monadap, fa)
	}
}
