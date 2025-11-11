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

package bounded

import "github.com/IBM/fp-go/v2/ord"

type Bounded[T any] interface {
	ord.Ord[T]
	Top() T
	Bottom() T
}

type bounded[T any] struct {
	c func(x, y T) int
	e func(x, y T) bool
	t T
	b T
}

func (b bounded[T]) Equals(x, y T) bool {
	return b.e(x, y)
}

func (b bounded[T]) Compare(x, y T) int {
	return b.c(x, y)
}

func (b bounded[T]) Top() T {
	return b.t
}

func (b bounded[T]) Bottom() T {
	return b.b
}

// MakeBounded creates an instance of a bounded type
func MakeBounded[T any](o ord.Ord[T], t, b T) Bounded[T] {
	return bounded[T]{c: o.Compare, e: o.Equals, t: t, b: b}
}

// Clamp returns a function that clamps against the bounds defined in the bounded type
func Clamp[T any](b Bounded[T]) func(T) T {
	return ord.Clamp(b)(b.Bottom(), b.Top())
}

// Reverse reverses the ordering and swaps the bounds
func Reverse[T any](b Bounded[T]) Bounded[T] {
	return MakeBounded(ord.Reverse(b), b.Bottom(), b.Top())
}
