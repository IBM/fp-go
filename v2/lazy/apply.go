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

package lazy

import (
	IO "github.com/IBM/fp-go/v2/io"
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// ApplySemigroup lifts a Semigroup[A] to a Semigroup[Lazy[A]].
// This allows you to combine lazy computations using the semigroup operation
// on their underlying values.
//
// The resulting semigroup's Concat operation will evaluate both lazy computations
// and combine their results using the original semigroup's operation.
//
// Parameters:
//   - s: A semigroup for values of type A
//
// Returns:
//   - A semigroup for lazy computations of type A
//
// Example:
//
//	import (
//	    M "github.com/IBM/fp-go/v2/monoid"
//	    "github.com/IBM/fp-go/v2/lazy"
//	)
//
//	// Create a semigroup for lazy integers using addition
//	intAddSemigroup := lazy.ApplySemigroup(M.MonoidSum[int]())
//
//	lazy1 := lazy.Of(5)
//	lazy2 := lazy.Of(10)
//
//	// Combine the lazy computations
//	result := intAddSemigroup.Concat(lazy1, lazy2)() // 15
func ApplySemigroup[A any](s S.Semigroup[A]) S.Semigroup[Lazy[A]] {
	return IO.ApplySemigroup(s)
}

// ApplicativeMonoid lifts a Monoid[A] to a Monoid[Lazy[A]].
// This allows you to combine lazy computations using the monoid operation
// on their underlying values, with an identity element.
//
// The resulting monoid's Concat operation will evaluate both lazy computations
// and combine their results using the original monoid's operation. The Empty
// operation returns a lazy computation that produces the monoid's identity element.
//
// Parameters:
//   - m: A monoid for values of type A
//
// Returns:
//   - A monoid for lazy computations of type A
//
// Example:
//
//	import (
//	    M "github.com/IBM/fp-go/v2/monoid"
//	    "github.com/IBM/fp-go/v2/lazy"
//	)
//
//	// Create a monoid for lazy integers using addition
//	intAddMonoid := lazy.ApplicativeMonoid(M.MonoidSum[int]())
//
//	// Get the identity element (0 wrapped in lazy)
//	empty := intAddMonoid.Empty()() // 0
//
//	lazy1 := lazy.Of(5)
//	lazy2 := lazy.Of(10)
//
//	// Combine the lazy computations
//	result := intAddMonoid.Concat(lazy1, lazy2)() // 15
//
//	// Identity laws hold:
//	// Concat(Empty(), x) == x
//	// Concat(x, Empty()) == x
func ApplicativeMonoid[A any](m M.Monoid[A]) M.Monoid[Lazy[A]] {
	return IO.ApplicativeMonoid(m)
}
