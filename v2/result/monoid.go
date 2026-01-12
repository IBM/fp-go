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

// FirstMonoid creates a Monoid for Result[A] that returns the first Ok (Right) value.
// This monoid prefers the left operand when it is Ok, otherwise returns the right operand.
// The empty value is provided as a lazy computation.
//
// This is equivalent to AltMonoid but implemented more directly.
//
// Truth table:
//
//	| x         | y         | concat(x, y) |
//	| --------- | --------- | ------------ |
//	| err(e1)   | err(e2)   | err(e2)      |
//	| ok(a)     | err(e)    | ok(a)        |
//	| err(e)    | ok(b)     | ok(b)        |
//	| ok(a)     | ok(b)     | ok(a)        |
//
// Example:
//
//	import "errors"
//	zero := func() result.Result[int] { return result.Error[int](errors.New("empty")) }
//	m := result.FirstMonoid[int](zero)
//	m.Concat(result.Of(2), result.Of(3)) // Ok(2) - returns first Ok
//	m.Concat(result.Error[int](errors.New("err")), result.Of(3)) // Ok(3)
//	m.Concat(result.Of(2), result.Error[int](errors.New("err"))) // Ok(2)
//	m.Empty() // Error(error("empty"))
//
//go:inline
func FirstMonoid[A any](zero Lazy[Result[A]]) M.Monoid[Result[A]] {
	return either.FirstMonoid(zero)
}

// LastMonoid creates a Monoid for Result[A] that returns the last Ok (Right) value.
// This monoid prefers the right operand when it is Ok, otherwise returns the left operand.
// The empty value is provided as a lazy computation.
//
// Truth table:
//
//	| x         | y         | concat(x, y) |
//	| --------- | --------- | ------------ |
//	| err(e1)   | err(e2)   | err(e1)      |
//	| ok(a)     | err(e)    | ok(a)        |
//	| err(e)    | ok(b)     | ok(b)        |
//	| ok(a)     | ok(b)     | ok(b)        |
//
// Example:
//
//	import "errors"
//	zero := func() result.Result[int] { return result.Error[int](errors.New("empty")) }
//	m := result.LastMonoid[int](zero)
//	m.Concat(result.Of(2), result.Of(3)) // Ok(3) - returns last Ok
//	m.Concat(result.Error[int](errors.New("err")), result.Of(3)) // Ok(3)
//	m.Concat(result.Of(2), result.Error[int](errors.New("err"))) // Ok(2)
//	m.Empty() // Error(error("empty"))
//
//go:inline
func LastMonoid[A any](zero Lazy[Result[A]]) M.Monoid[Result[A]] {
	return either.LastMonoid(zero)
}
