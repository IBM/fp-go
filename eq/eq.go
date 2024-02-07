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
	F "github.com/IBM/fp-go/function"
)

type Eq[T any] interface {
	Equals(x, y T) bool
}

type eq[T any] struct {
	c func(x, y T) bool
}

func (e eq[T]) Equals(x, y T) bool {
	return e.c(x, y)
}

func strictEq[A comparable](a, b A) bool {
	return a == b
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[T comparable]() Eq[T] {
	return FromEquals(strictEq[T])
}

// FromEquals constructs an [EQ.Eq] from the comparison function
func FromEquals[T any](c func(x, y T) bool) Eq[T] {
	return eq[T]{c: c}
}

// Empty returns the equals predicate that is always true
func Empty[T any]() Eq[T] {
	return FromEquals(F.Constant2[T, T](true))
}

// Equals returns a predicate to test if one value equals the other under an equals predicate
func Equals[T any](eq Eq[T]) func(T) func(T) bool {
	return func(other T) func(T) bool {
		return F.Bind2nd(eq.Equals, other)
	}
}
