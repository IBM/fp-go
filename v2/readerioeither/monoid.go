//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package readerioeither

import (
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
)

type (
	Monoid[R, E, A any] = monoid.Monoid[ReaderIOEither[R, E, A]]
)

// ApplicativeMonoid returns a [Monoid] that concatenates [ReaderIOEither] instances via their applicative.
// This uses the default applicative behavior (parallel or sequential based on useParallel flag).
//
// The monoid combines two ReaderIOEither values by applying the underlying monoid's combine operation
// to their success values using applicative application.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderIOEither[R, E, A].
func ApplicativeMonoid[R, E, A any](m monoid.Monoid[A]) Monoid[R, E, A] {
	return monoid.ApplicativeMonoid(
		Of[R, E, A],
		MonadMap[R, E, A, func(A) A],
		MonadAp[R, E, A, A],
		m,
	)
}

// ApplicativeMonoidSeq returns a [Monoid] that concatenates [ReaderIOEither] instances via their applicative.
// This explicitly uses sequential execution for combining values.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderIOEither[R, E, A] with sequential execution.
func ApplicativeMonoidSeq[R, E, A any](m monoid.Monoid[A]) Monoid[R, E, A] {
	return monoid.ApplicativeMonoid(
		Of[R, E, A],
		MonadMap[R, E, A, func(A) A],
		MonadApSeq[R, E, A, A],
		m,
	)
}

// ApplicativeMonoidPar returns a [Monoid] that concatenates [ReaderIOEither] instances via their applicative.
// This explicitly uses parallel execution for combining values.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderIOEither[R, E, A] with parallel execution.
func ApplicativeMonoidPar[R, E, A any](m monoid.Monoid[A]) Monoid[R, E, A] {
	return monoid.ApplicativeMonoid(
		Of[R, E, A],
		MonadMap[R, E, A, func(A) A],
		MonadApPar[R, E, A, A],
		m,
	)
}

// AlternativeMonoid is the alternative [Monoid] for [ReaderIOEither].
// This combines ReaderIOEither values using the alternative semantics,
// where the second value is only evaluated if the first fails.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderIOEither[R, E, A] with alternative semantics.
func AlternativeMonoid[R, E, A any](m monoid.Monoid[A]) Monoid[R, E, A] {
	return monoid.AlternativeMonoid(
		Of[R, E, A],
		MonadMap[R, E, A, func(A) A],
		MonadAp[R, E, A, A],
		MonadAlt[R, E, A],
		m,
	)
}

// AltMonoid is the alternative [Monoid] for a [ReaderIOEither].
// This creates a monoid where the empty value is provided lazily,
// and combination uses the Alt operation (try first, fallback to second on failure).
//
// Parameters:
//   - zero: Lazy computation that provides the empty/identity value
//
// Returns a Monoid for ReaderIOEither[R, E, A] with Alt-based combination.
func AltMonoid[R, E, A any](zero lazy.Lazy[ReaderIOEither[R, E, A]]) Monoid[R, E, A] {
	return monoid.AltMonoid(
		zero,
		MonadAlt[R, E, A],
	)
}
