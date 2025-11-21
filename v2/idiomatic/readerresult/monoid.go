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

package readerresult

import (
	M "github.com/IBM/fp-go/v2/monoid"
)

// AlternativeMonoid creates a Monoid for ReaderResult that combines both Alternative and Applicative behavior.
// It uses the provided monoid for the success values and falls back to alternative computations on failure.
//
// The empty element is Of(m.Empty()), and concat tries the first computation, falling back to the second
// if it fails, then combines successful values using the underlying monoid.
//
// Example:
//
//	intAdd := monoid.MakeMonoid(0, func(a, b int) int { return a + b })
//	rrMonoid := readerresult.AlternativeMonoid[Config](intAdd)
//
//	rr1 := readerresult.Of[Config](5)
//	rr2 := readerresult.Of[Config](3)
//	combined := rrMonoid.Concat(rr1, rr2)
//	// combined(cfg) returns (8, nil)
//
//go:inline
func AlternativeMonoid[R, A any](m M.Monoid[A]) Monoid[R, A] {
	return M.AlternativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadAp[A, R, A],
		MonadAlt[R, A],
		m,
	)
}

// AltMonoid creates a Monoid for ReaderResult based on the Alternative pattern.
// The empty element is the provided zero computation, and concat tries the first computation,
// falling back to the second if it fails.
//
// This is useful for combining computations where you want to try alternatives until one succeeds.
//
// Example:
//
//	zero := func() readerresult.ReaderResult[Config, User] {
//	    return readerresult.Left[Config, User](errors.New("no user"))
//	}
//	userMonoid := readerresult.AltMonoid[Config](zero)
//
//	primary := getPrimaryUser()
//	backup := getBackupUser()
//	combined := userMonoid.Concat(primary, backup)
//	// Tries primary, falls back to backup if primary fails
//
//go:inline
func AltMonoid[R, A any](zero Lazy[ReaderResult[R, A]]) Monoid[R, A] {
	return M.AltMonoid(
		zero,
		MonadAlt[R, A],
	)
}

// ApplicativeMonoid creates a Monoid for ReaderResult based on Applicative functor composition.
// The empty element is Of(m.Empty()), and concat combines two computations using the underlying monoid.
// Both computations must succeed for the result to succeed.
//
// This is useful for accumulating results from multiple independent computations.
//
// Example:
//
//	intAdd := monoid.MakeMonoid(0, func(a, b int) int { return a + b })
//	rrMonoid := readerresult.ApplicativeMonoid[Config](intAdd)
//
//	rr1 := readerresult.Of[Config](5)
//	rr2 := readerresult.Of[Config](3)
//	combined := rrMonoid.Concat(rr1, rr2)
//	// combined(cfg) returns (8, nil)
//
//	// If either fails, the whole computation fails
//	rr3 := readerresult.Left[Config, int](errors.New("error"))
//	failed := rrMonoid.Concat(rr1, rr3)
//	// failed(cfg) returns (nil, error)
//
//go:inline
func ApplicativeMonoid[R, A any](m M.Monoid[A]) Monoid[R, A] {
	return M.ApplicativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadAp[A, R, A],
		m,
	)
}
