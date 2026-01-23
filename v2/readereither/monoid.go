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

package readereither

import (
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
)

// ApplicativeMonoid returns a [Monoid] that concatenates [ReaderEither] instances via their applicative.
// This combines two ReaderEither values by applying the underlying monoid's combine operation
// to their success values using applicative application.
//
// The applicative behavior means that if either computation fails (returns Left), the entire
// combination fails. Both computations must succeed (return Right) for the result to succeed.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderEither[R, E, A].
//
// Example:
//
//	intMonoid := number.MonoidSum[int]()
//	reMonoid := ApplicativeMonoid[Config, string](intMonoid)
//
//	re1 := Right[Config, string](5)
//	re2 := Right[Config, string](3)
//	combined := reMonoid.Concat(re1, re2)
//	// Result: Right(8)
//
//	re3 := Left[Config, int]("error")
//	failed := reMonoid.Concat(re1, re3)
//	// Result: Left("error")
//
//go:inline
func ApplicativeMonoid[R, E, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderEither[R, E, A]] {
	return monoid.ApplicativeMonoid(
		Of[R, E, A],
		MonadMap[R, E, A, func(A) A],
		MonadAp[A, R, E, A],
		m,
	)
}

// AlternativeMonoid is the alternative [Monoid] for [ReaderEither].
// This combines ReaderEither values using the alternative semantics,
// where the second value is only evaluated if the first fails.
//
// The alternative behavior provides fallback semantics: if the first computation
// succeeds (returns Right), its value is used. If it fails (returns Left), the
// second computation is tried. If both succeed, their values are combined using
// the underlying monoid.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderEither[R, E, A] with alternative semantics.
//
// Example:
//
//	intMonoid := number.MonoidSum[int]()
//	reMonoid := AlternativeMonoid[Config, string](intMonoid)
//
//	re1 := Left[Config, int]("error1")
//	re2 := Right[Config, string](42)
//	combined := reMonoid.Concat(re1, re2)
//	// Result: Right(42) - falls back to second
//
//	re3 := Right[Config, string](5)
//	re4 := Right[Config, string](3)
//	both := reMonoid.Concat(re3, re4)
//	// Result: Right(8) - combines both successes
//
//go:inline
func AlternativeMonoid[R, E, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderEither[R, E, A]] {
	return monoid.AlternativeMonoid(
		Of[R, E, A],
		MonadMap[R, E, A, func(A) A],
		MonadAp[A, R, E, A],
		MonadAlt[R, E, A],
		m,
	)
}

// AltMonoid is the alternative [Monoid] for a [ReaderEither].
// This creates a monoid where the empty value is provided lazily,
// and combination uses the Alt operation (try first, fallback to second on failure).
//
// Unlike AlternativeMonoid, this does not combine successful values using an underlying
// monoid. Instead, it simply returns the first successful value, or falls back to the
// second if the first fails.
//
// Parameters:
//   - zero: Lazy computation that provides the empty/identity value
//
// Returns a Monoid for ReaderEither[R, E, A] with Alt-based combination.
//
// Example:
//
//	zero := lazy.MakeLazy(func() ReaderEither[Config, string, int] {
//	    return Left[Config, int]("no value")
//	})
//	reMonoid := AltMonoid(zero)
//
//	re1 := Left[Config, int]("error1")
//	re2 := Right[Config, string](42)
//	combined := reMonoid.Concat(re1, re2)
//	// Result: Right(42) - uses first success
//
//	re3 := Right[Config, string](100)
//	re4 := Right[Config, string](200)
//	first := reMonoid.Concat(re3, re4)
//	// Result: Right(100) - uses first success, doesn't combine
//
//go:inline
func AltMonoid[R, E, A any](zero lazy.Lazy[ReaderEither[R, E, A]]) monoid.Monoid[ReaderEither[R, E, A]] {
	return monoid.AltMonoid(
		zero,
		MonadAlt[R, E, A],
	)
}
