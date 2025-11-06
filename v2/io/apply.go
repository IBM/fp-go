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

package io

import (
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// ApplySemigroup lifts a Semigroup[A] into a Semigroup[IO[A]].
// This allows combining IO computations using the semigroup operation on their results.
//
// Example:
//
//	intAdd := semigroup.MakeSemigroup(func(a, b int) int { return a + b })
//	ioAdd := io.ApplySemigroup(intAdd)
//	result := ioAdd.Concat(io.Of(1), io.Of(2)) // IO[3]
func ApplySemigroup[A any](s S.Semigroup[A]) Semigroup[A] {
	return S.ApplySemigroup(MonadMap[A, func(A) A], MonadAp[A, A], s)
}

// ApplicativeMonoid lifts a Monoid[A] into a Monoid[IO[A]].
// This allows combining IO computations using the monoid operation on their results,
// including an empty/identity element.
//
// Example:
//
//	intAdd := monoid.MakeMonoid(func(a, b int) int { return a + b }, 0)
//	ioAdd := io.ApplicativeMonoid(intAdd)
//	result := ioAdd.Concat(io.Of(1), io.Of(2)) // IO[3]
//	empty := ioAdd.Empty() // IO[0]
func ApplicativeMonoid[A any](m M.Monoid[A]) Monoid[A] {
	return M.ApplicativeMonoid(Of[A], MonadMap[A, func(A) A], MonadAp[A, A], m)
}
