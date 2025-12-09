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
	"context"

	RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	M "github.com/IBM/fp-go/v2/monoid"
)

// AlternativeMonoid creates a Monoid for ReaderResult using the Alternative semantics.
//
// The Alternative semantics means that the monoid operation tries the first computation,
// and if it fails, tries the second one. The empty element is a computation that always fails.
// The inner values are combined using the provided monoid when both computations succeed.
//
// Type Parameters:
//   - A: The value type
//
// Parameters:
//   - m: A Monoid[A] for combining successful values
//
// Returns:
//   - A Monoid[ReaderResult[A]] with Alternative semantics
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/monoid"
//
//	// Monoid for integers with addition
//	intMonoid := monoid.MonoidSum[int]()
//	rrMonoid := readerresult.AlternativeMonoid(intMonoid)
//
//	rr1 := readerresult.Right(10)
//	rr2 := readerresult.Right(20)
//	combined := rrMonoid.Concat(rr1, rr2)
//	value, err := combined(ctx)  // Returns (30, nil)
//
//go:inline
func AlternativeMonoid[A any](m M.Monoid[A]) Monoid[A] {
	return RR.AlternativeMonoid[context.Context](m)
}

// AltMonoid creates a Monoid for ReaderResult using Alt semantics with a custom zero.
//
// The Alt semantics means that the monoid operation tries the first computation,
// and if it fails, tries the second one. The provided zero is used as the empty element.
//
// Type Parameters:
//   - A: The value type
//
// Parameters:
//   - zero: A lazy ReaderResult[A] to use as the empty element
//
// Returns:
//   - A Monoid[ReaderResult[A]] with Alt semantics
//
// Example:
//
//	zero := func() readerresult.ReaderResult[int] {
//	    return readerresult.Left[int](errors.New("empty"))
//	}
//	rrMonoid := readerresult.AltMonoid(zero)
//
//	rr1 := readerresult.Left[int](errors.New("failed"))
//	rr2 := readerresult.Right(42)
//	combined := rrMonoid.Concat(rr1, rr2)
//	value, err := combined(ctx)  // Returns (42, nil) - uses second on first failure
//
//go:inline
func AltMonoid[A any](zero Lazy[ReaderResult[A]]) Monoid[A] {
	return RR.AltMonoid(zero)
}

// ApplicativeMonoid creates a Monoid for ReaderResult using Applicative semantics.
//
// The Applicative semantics means that both computations are executed independently,
// and their results are combined using the provided monoid. If either fails, the
// entire operation fails.
//
// Type Parameters:
//   - A: The value type
//
// Parameters:
//   - m: A Monoid[A] for combining successful values
//
// Returns:
//   - A Monoid[ReaderResult[A]] with Applicative semantics
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/monoid"
//
//	// Monoid for integers with addition
//	intMonoid := monoid.MonoidSum[int]()
//	rrMonoid := readerresult.ApplicativeMonoid(intMonoid)
//
//	rr1 := readerresult.Right(10)
//	rr2 := readerresult.Right(20)
//	combined := rrMonoid.Concat(rr1, rr2)
//	value, err := combined(ctx)  // Returns (30, nil)
//
//go:inline
func ApplicativeMonoid[A any](m M.Monoid[A]) Monoid[A] {
	return RR.ApplicativeMonoid[context.Context](m)
}
