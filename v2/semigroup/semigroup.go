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

package semigroup

import (
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/magma"
)

type Semigroup[A any] interface {
	M.Magma[A]
}

type semigroup[A any] struct {
	c func(A, A) A
}

func (self semigroup[A]) Concat(x A, y A) A {
	return self.c(x, y)
}

func MakeSemigroup[A any](c func(A, A) A) Semigroup[A] {
	return semigroup[A]{c: c}
}

// Reverse returns The dual of a `Semigroup`, obtained by swapping the arguments of `concat`.
func Reverse[A any](m Semigroup[A]) Semigroup[A] {
	return MakeSemigroup(M.Reverse[A](m).Concat)
}

// FunctionSemigroup forms a semigroup as long as you can provide a semigroup for the codomain.
func FunctionSemigroup[A, B any](s Semigroup[B]) Semigroup[func(A) B] {
	return MakeSemigroup(func(f func(A) B, g func(A) B) func(A) B {
		return func(a A) B {
			return s.Concat(f(a), g(a))
		}
	})
}

// First always returns the first argument.
func First[A any]() Semigroup[A] {
	return MakeSemigroup(F.First[A, A])
}

// Last always returns the last argument.
func Last[A any]() Semigroup[A] {
	return MakeSemigroup(F.Second[A, A])
}

// ToMagma converts a semigroup to a magma
func ToMagma[A any](s Semigroup[A]) M.Magma[A] {
	return s
}
