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

package readerioresult

import (
	"context"

	"github.com/IBM/fp-go/v2/monoid"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

type (
	Monoid[A any] = monoid.Monoid[ReaderIOResult[A]]
)

// ApplicativeMonoid returns a [Monoid] that concatenates [ReaderIOResult] instances via their applicative.
// This uses the default applicative behavior (parallel or sequential based on useParallel flag).
//
// The monoid combines two ReaderIOResult values by applying the underlying monoid's combine operation
// to their success values using applicative application.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderIOResult[A].
func ApplicativeMonoid[A any](m monoid.Monoid[A]) Monoid[A] {
	return RIOR.ApplicativeMonoid[context.Context](m)
}

// ApplicativeMonoidSeq returns a [Monoid] that concatenates [ReaderIOResult] instances via their applicative.
// This explicitly uses sequential execution for combining values.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderIOResult[A] with sequential execution.
func ApplicativeMonoidSeq[A any](m monoid.Monoid[A]) Monoid[A] {
	return RIOR.ApplicativeMonoidSeq[context.Context](m)
}

// ApplicativeMonoidPar returns a [Monoid] that concatenates [ReaderIOResult] instances via their applicative.
// This explicitly uses parallel execution for combining values.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderIOResult[A] with parallel execution.
func ApplicativeMonoidPar[A any](m monoid.Monoid[A]) Monoid[A] {
	return RIOR.ApplicativeMonoidPar[context.Context](m)
}

// AlternativeMonoid is the alternative [Monoid] for [ReaderIOResult].
// This combines ReaderIOResult values using the alternative semantics,
// where the second value is only evaluated if the first fails.
//
// Parameters:
//   - m: The underlying monoid for type A
//
// Returns a Monoid for ReaderIOResult[A] with alternative semantics.
func AlternativeMonoid[A any](m monoid.Monoid[A]) Monoid[A] {
	return RIOR.AlternativeMonoid[context.Context](m)
}

// AltMonoid is the alternative [Monoid] for a [ReaderIOResult].
// This creates a monoid where the empty value is provided lazily,
// and combination uses the Alt operation (try first, fallback to second on failure).
//
// Parameters:
//   - zero: Lazy computation that provides the empty/identity value
//
// Returns a Monoid for ReaderIOResult[A] with Alt-based combination.
func AltMonoid[A any](zero Lazy[ReaderIOResult[A]]) Monoid[A] {
	return RIOR.AltMonoid(zero)
}
