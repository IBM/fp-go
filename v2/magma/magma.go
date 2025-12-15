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

package magma

type Magma[A any] interface {
	Concat(x A, y A) A
}

type magma[A any] struct {
	c func(A, A) A
}

func (m magma[A]) Concat(x, y A) A {
	return m.c(x, y)
}

func MakeMagma[A any](c func(A, A) A) Magma[A] {
	return magma[A]{c: c}
}

func Reverse[A any](m Magma[A]) Magma[A] {
	return MakeMagma(func(x A, y A) A {
		return m.Concat(y, x)
	})
}

func filterFirst[A any](p func(A) bool, c func(A, A) A, x, y A) A {
	if p(x) {
		return c(x, y)
	}
	return y
}

func filterSecond[A any](p func(A) bool, c func(A, A) A, x, y A) A {
	if p(y) {
		return c(x, y)
	}
	return x
}

func FilterFirst[A any](p func(A) bool) func(Magma[A]) Magma[A] {
	return func(m Magma[A]) Magma[A] {
		c := m.Concat
		return MakeMagma(func(x A, y A) A {
			return filterFirst(p, c, x, y)
		})
	}
}

func FilterSecond[A any](p func(A) bool) func(Magma[A]) Magma[A] {
	return func(m Magma[A]) Magma[A] {
		c := m.Concat
		return MakeMagma(func(x, y A) A {
			return filterSecond(p, c, x, y)
		})
	}
}

func first[A any](x, _ A) A {
	return x
}

func second[A any](_, y A) A {
	return y
}

func First[A any]() Magma[A] {
	return MakeMagma(first[A])
}

func Second[A any]() Magma[A] {
	return MakeMagma(second[A])
}

func endo[A any](f func(A) A, c func(A, A) A, x, y A) A {
	return c(f(x), f(y))
}

func Endo[A any](f func(A) A) func(Magma[A]) Magma[A] {
	return func(m Magma[A]) Magma[A] {
		c := m.Concat
		return MakeMagma(func(x A, y A) A {
			return endo(f, c, x, y)
		})
	}
}
