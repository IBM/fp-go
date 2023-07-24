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

package bounded

import (
	O "github.com/IBM/fp-go/ord"
)

type Bounded[T any] interface {
	O.Ord[T]
	Top() T
	Bottom() T
}

type bounded[T any] struct {
	c func(x, y T) int
	e func(x, y T) bool
	t T
	b T
}

func (self bounded[T]) Equals(x, y T) bool {
	return self.e(x, y)
}

func (self bounded[T]) Compare(x, y T) int {
	return self.c(x, y)
}

func (self bounded[T]) Top() T {
	return self.t
}

func (self bounded[T]) Bottom() T {
	return self.b
}

// MakeBounded creates an instance of a bounded type
func MakeBounded[T any](o O.Ord[T], t, b T) Bounded[T] {
	return bounded[T]{c: o.Compare, e: o.Equals, t: t, b: b}
}

// Clamp returns a function that clamps against the bounds defined in the bounded type
func Clamp[T any](b Bounded[T]) func(T) T {
	return O.Clamp[T](b)(b.Bottom(), b.Top())
}

// Reverse reverses the ordering and swaps the bounds
func Reverse[T any](b Bounded[T]) Bounded[T] {
	return MakeBounded(O.Reverse[T](b), b.Bottom(), b.Top())
}
