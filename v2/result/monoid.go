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

package result

import (
	"github.com/IBM/fp-go/v2/either"
	M "github.com/IBM/fp-go/v2/monoid"
)

// AlternativeMonoid creates a monoid for Either using applicative semantics.
// The empty value is Right with the monoid's empty value.
// Combines values using applicative operations.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/monoid"
//	intAdd := monoid.MakeMonoid(0, func(a, b int) int { return a + b })
//	m := either.AlternativeMonoid[error](intAdd)
//	result := m.Concat(either.Right[error](1), either.Right[error](2))
//	// result is Right(3)
//
//go:inline
func AlternativeMonoid[A any](m M.Monoid[A]) Monoid[A] {
	return either.AlternativeMonoid[error](m)
}

// AltMonoid creates a monoid for Either using the Alt operation.
// The empty value is provided as a lazy computation.
// When combining, returns the first Right value, or the second if the first is Left.
//
// Example:
//
//	zero := func() either.Result[int] { return either.Left[int](errors.New("empty")) }
//	m := either.AltMonoid[error, int](zero)
//	result := m.Concat(either.Left[int](errors.New("err1")), either.Right[error](42))
//	// result is Right(42)
//
//go:inline
func AltMonoid[A any](zero Lazy[Result[A]]) Monoid[A] {
	return either.AltMonoid(zero)
}
