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

package either

import (
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
func AlternativeMonoid[E, A any](m M.Monoid[A]) Monoid[E, A] {
	return M.AlternativeMonoid(
		Of[E, A],
		MonadMap[E, A, func(A) A],
		MonadAp[A, E, A],
		MonadAlt[E, A],
		m,
	)
}

// AltMonoid creates a monoid for Either using the Alt operation.
// The empty value is provided as a lazy computation.
// When combining, returns the first Right value, or the second if the first is Left.
//
// Example:
//
//	zero := func() either.Either[error, int] { return either.Left[int](errors.New("empty")) }
//	m := either.AltMonoid[error, int](zero)
//	result := m.Concat(either.Left[int](errors.New("err1")), either.Right[error](42))
//	// result is Right(42)
func AltMonoid[E, A any](zero Lazy[Either[E, A]]) Monoid[E, A] {
	return M.AltMonoid(
		zero,
		MonadAlt[E, A],
	)
}

// takeFirst is a helper function that returns the first Right value, or the second if the first is Left.
func takeFirst[E, A any](l, r Either[E, A]) Either[E, A] {
	if IsRight(l) {
		return l
	}
	return r
}

// FirstMonoid creates a Monoid for Either[E, A] that returns the first Right value.
// This monoid prefers the left operand when it is Right, otherwise returns the right operand.
// The empty value is provided as a lazy computation.
//
// This is equivalent to AltMonoid but implemented more directly.
//
// Truth table:
//
//	| x         | y         | concat(x, y) |
//	| --------- | --------- | ------------ |
//	| left(e1)  | left(e2)  | left(e2)     |
//	| right(a)  | left(e)   | right(a)     |
//	| left(e)   | right(b)  | right(b)     |
//	| right(a)  | right(b)  | right(a)     |
//
// Example:
//
//	import "errors"
//	zero := func() either.Either[error, int] { return either.Left[int](errors.New("empty")) }
//	m := either.FirstMonoid[error, int](zero)
//	m.Concat(either.Right[error](2), either.Right[error](3)) // Right(2) - returns first Right
//	m.Concat(either.Left[int](errors.New("err")), either.Right[error](3)) // Right(3)
//	m.Concat(either.Right[error](2), either.Left[int](errors.New("err"))) // Right(2)
//	m.Empty() // Left(error("empty"))
//
//go:inline
func FirstMonoid[E, A any](zero Lazy[Either[E, A]]) M.Monoid[Either[E, A]] {
	return M.MakeMonoid(takeFirst[E, A], zero())
}

// takeLast is a helper function that returns the last Right value, or the first if the last is Left.
func takeLast[E, A any](l, r Either[E, A]) Either[E, A] {
	if IsRight(r) {
		return r
	}
	return l
}

// LastMonoid creates a Monoid for Either[E, A] that returns the last Right value.
// This monoid prefers the right operand when it is Right, otherwise returns the left operand.
// The empty value is provided as a lazy computation.
//
// Truth table:
//
//	| x         | y         | concat(x, y) |
//	| --------- | --------- | ------------ |
//	| left(e1)  | left(e2)  | left(e1)     |
//	| right(a)  | left(e)   | right(a)     |
//	| left(e)   | right(b)  | right(b)     |
//	| right(a)  | right(b)  | right(b)     |
//
// Example:
//
//	import "errors"
//	zero := func() either.Either[error, int] { return either.Left[int](errors.New("empty")) }
//	m := either.LastMonoid[error, int](zero)
//	m.Concat(either.Right[error](2), either.Right[error](3)) // Right(3) - returns last Right
//	m.Concat(either.Left[int](errors.New("err")), either.Right[error](3)) // Right(3)
//	m.Concat(either.Right[error](2), either.Left[int](errors.New("err"))) // Right(2)
//	m.Empty() // Left(error("empty"))
//
//go:inline
func LastMonoid[E, A any](zero Lazy[Either[E, A]]) M.Monoid[Either[E, A]] {
	return M.MakeMonoid(takeLast[E, A], zero())
}
