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

package readerreaderioresult

import (
	"github.com/IBM/fp-go/v2/monoid"
)

type (
	// Monoid represents a monoid structure for ReaderReaderIOResult[R, A].
	// A monoid provides an identity element (empty) and an associative binary operation (concat).
	Monoid[R, A any] = monoid.Monoid[ReaderReaderIOResult[R, A]]
)

// ApplicativeMonoid creates a monoid for ReaderReaderIOResult using applicative composition.
// It combines values using the provided monoid m and the applicative Ap operation.
// This allows combining multiple ReaderReaderIOResult values in parallel while merging their results.
//
// The resulting monoid satisfies:
//   - Identity: concat(empty, x) = concat(x, empty) = x
//   - Associativity: concat(concat(x, y), z) = concat(x, concat(y, z))
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/monoid"
//	import "github.com/IBM/fp-go/v2/number"
//
//	// Create a monoid for combining integers with addition
//	intMonoid := ApplicativeMonoid[Config](number.MonoidSum)
//
//	// Combine multiple computations
//	result := intMonoid.Concat(
//	    Of[Config](10),
//	    intMonoid.Concat(Of[Config](20), Of[Config](30)),
//	) // Results in 60
func ApplicativeMonoid[R, A any](m monoid.Monoid[A]) Monoid[R, A] {
	return monoid.ApplicativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadAp[R, A, A],
		m,
	)
}

// ApplicativeMonoidSeq creates a monoid for ReaderReaderIOResult using sequential applicative composition.
// Similar to ApplicativeMonoid but evaluates effects sequentially rather than in parallel.
//
// Use this when:
//   - Effects must be executed in a specific order
//   - Side effects depend on sequential execution
//   - You want to avoid concurrent execution
func ApplicativeMonoidSeq[R, A any](m monoid.Monoid[A]) Monoid[R, A] {
	return monoid.ApplicativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadApSeq[R, A, A],
		m,
	)
}

// ApplicativeMonoidPar creates a monoid for ReaderReaderIOResult using parallel applicative composition.
// Similar to ApplicativeMonoid but explicitly evaluates effects in parallel.
//
// Use this when:
//   - Effects are independent and can run concurrently
//   - You want to maximize performance through parallelism
//   - Order of execution doesn't matter
func ApplicativeMonoidPar[R, A any](m monoid.Monoid[A]) Monoid[R, A] {
	return monoid.ApplicativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadApPar[R, A, A],
		m,
	)
}

// AlternativeMonoid creates a monoid that combines ReaderReaderIOResult values using both
// applicative composition and alternative (Alt) semantics.
//
// This monoid:
//   - Uses Ap for combining successful values
//   - Uses Alt for handling failures (tries alternatives on failure)
//   - Provides a way to combine multiple computations with fallback behavior
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/monoid"
//	import "github.com/IBM/fp-go/v2/number"
//
//	intMonoid := AlternativeMonoid[Config](number.MonoidSum)
//
//	// If first computation fails, tries the second
//	result := intMonoid.Concat(
//	    Left[Config, int](errors.New("failed")),
//	    Of[Config](42),
//	) // Results in Right(42)
func AlternativeMonoid[R, A any](m monoid.Monoid[A]) Monoid[R, A] {
	return monoid.AlternativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadAp[R, A, A],
		MonadAlt[R, A],
		m,
	)
}

// AltMonoid creates a monoid based solely on the Alt operation.
// It provides a way to chain computations with fallback behavior.
//
// The monoid:
//   - Uses the provided zero as the identity element
//   - Uses Alt for concatenation (tries first, falls back to second on failure)
//   - Implements a "first success" strategy
//
// Example:
//
//	zero := func() ReaderReaderIOResult[Config, int] {
//	    return Left[Config, int](errors.New("no value"))
//	}
//	altMonoid := AltMonoid[Config, int](zero)
//
//	// Tries computations in order until one succeeds
//	result := altMonoid.Concat(
//	    Left[Config, int](errors.New("first failed")),
//	    altMonoid.Concat(
//	        Left[Config, int](errors.New("second failed")),
//	        Of[Config](42),
//	    ),
//	) // Results in Right(42)
func AltMonoid[R, A any](zero Lazy[ReaderReaderIOResult[R, A]]) Monoid[R, A] {
	return monoid.AltMonoid(
		zero,
		MonadAlt[R, A],
	)
}
