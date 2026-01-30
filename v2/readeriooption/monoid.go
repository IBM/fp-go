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

package readeriooption

import "github.com/IBM/fp-go/v2/monoid"

// ApplicativeMonoid creates a Monoid for ReaderIOOption based on Applicative functor composition.
// The empty element is Of(m.Empty()), and concat combines two computations using the underlying monoid.
// Both computations must succeed (return Some) for the result to succeed.
//
// This is useful for accumulating results from multiple independent computations that all need
// to succeed. If any computation returns None, the entire result is None.
//
// The resulting monoid satisfies the monoid laws:
//   - Left identity: Concat(Empty(), x) = x
//   - Right identity: Concat(x, Empty()) = x
//   - Associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
//
// Parameters:
//   - m: The underlying monoid for combining success values of type A
//
// Returns:
//   - A Monoid[ReaderIOOption[R, A]] that combines ReaderIOOption computations
//
// Example:
//
//	import (
//	    N "github.com/IBM/fp-go/v2/number"
//	    RO "github.com/IBM/fp-go/v2/readeroption"
//	)
//
//	// Create a monoid for integer addition
//	intAdd := N.MonoidSum[int]()
//	roMonoid := RO.ApplicativeMonoid[Config](intAdd)
//
//	// Combine successful computations
//	ro1 := RO.Of[Config](5)
//	ro2 := RO.Of[Config](3)
//	combined := roMonoid.Concat(ro1, ro2)
//	// combined(cfg) returns option.Some(8)
//
//	// If either fails, the whole computation fails
//	ro3 := RO.None[Config, int]()
//	failed := roMonoid.Concat(ro1, ro3)
//	// failed(cfg) returns option.None[int]()
//
//	// Empty element is the identity
//	withEmpty := roMonoid.Concat(ro1, roMonoid.Empty())
//	// withEmpty(cfg) returns option.Some(5)
//
//go:inline
func ApplicativeMonoid[R, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderIOOption[R, A]] {
	return monoid.ApplicativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadAp[R, A, A],
		m,
	)
}

// AlternativeMonoid creates a Monoid for ReaderIOOption that combines both Alternative and Applicative behavior.
// It uses the provided monoid for the success values and falls back to alternative computations on failure.
//
// The empty element is Of(m.Empty()), and concat tries the first computation, falling back to the second
// if it fails (returns None), then combines successful values using the underlying monoid.
//
// This is particularly useful when you want to:
//   - Try multiple computations and accumulate their results
//   - Provide fallback behavior when computations fail
//   - Combine results from computations that may or may not succeed
//
// The behavior differs from ApplicativeMonoid in that it provides fallback semantics:
//   - If the first computation succeeds, use its value
//   - If the first fails but the second succeeds, use the second's value
//   - If both succeed, combine their values using the underlying monoid
//   - If both fail, the result is None
//
// The resulting monoid satisfies the monoid laws:
//   - Left identity: Concat(Empty(), x) = x
//   - Right identity: Concat(x, Empty()) = x
//   - Associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
//
// Parameters:
//   - m: The underlying monoid for combining success values of type A
//
// Returns:
//   - A Monoid[ReaderIOOption[R, A]] that combines ReaderIOOption computations with fallback
//
// Example:
//
//	import (
//	    N "github.com/IBM/fp-go/v2/number"
//	    RO "github.com/IBM/fp-go/v2/readeroption"
//	)
//
//	// Create a monoid for integer addition with alternative behavior
//	intAdd := N.MonoidSum[int]()
//	roMonoid := RO.AlternativeMonoid[Config](intAdd)
//
//	// Combine successful computations
//	ro1 := RO.Of[Config](5)
//	ro2 := RO.Of[Config](3)
//	combined := roMonoid.Concat(ro1, ro2)
//	// combined(cfg) returns option.Some(8)
//
//	// Fallback when first fails
//	ro3 := RO.None[Config, int]()
//	ro4 := RO.Of[Config](10)
//	withFallback := roMonoid.Concat(ro3, ro4)
//	// withFallback(cfg) returns option.Some(10)
//
//	// Use first success when available
//	withFirst := roMonoid.Concat(ro1, ro3)
//	// withFirst(cfg) returns option.Some(5)
//
//	// Accumulate multiple values with some failures
//	result := roMonoid.Concat(
//	    roMonoid.Concat(ro3, ro1),  // None + 5 = 5
//	    ro2,                         // 5 + 3 = 8
//	)
//	// result(cfg) returns option.Some(8)
//
//go:inline
func AlternativeMonoid[R, A any](m monoid.Monoid[A]) monoid.Monoid[ReaderIOOption[R, A]] {
	return monoid.AlternativeMonoid(
		Of[R, A],
		MonadMap[R, A, func(A) A],
		MonadAp[R, A, A],
		MonadAlt[R, A],
		m,
	)
}
