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

package effect

import (
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/monoid"
)

// ApplicativeMonoid creates a monoid for effects using applicative semantics.
// This combines effects by running both and combining their results using the provided monoid.
// If either effect fails, the combined effect fails.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - A: The value type that has a monoid instance
//
// # Parameters
//
//   - m: The monoid instance for combining values of type A
//
// # Returns
//
//   - Monoid[Effect[C, A]]: A monoid for combining effects
//
// # Example
//
//	stringMonoid := monoid.MakeMonoid(
//		func(a, b string) string { return a + b },
//		"",
//	)
//	effectMonoid := effect.ApplicativeMonoid[MyContext](stringMonoid)
//	eff1 := effect.Of[MyContext]("Hello")
//	eff2 := effect.Of[MyContext](" World")
//	combined := effectMonoid.Concat(eff1, eff2)
//	// combined produces "Hello World"
func ApplicativeMonoid[C, A any](m monoid.Monoid[A]) Monoid[Effect[C, A]] {
	return readerreaderioresult.ApplicativeMonoid[C](m)
}

// AlternativeMonoid creates a monoid for effects using alternative semantics.
// This tries the first effect, and if it fails, tries the second effect.
// If both succeed, their results are combined using the provided monoid.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - A: The value type that has a monoid instance
//
// # Parameters
//
//   - m: The monoid instance for combining values of type A
//
// # Returns
//
//   - Monoid[Effect[C, A]]: A monoid for combining effects with fallback behavior
//
// # Example
//
//	stringMonoid := monoid.MakeMonoid(
//		func(a, b string) string { return a + b },
//		"",
//	)
//	effectMonoid := effect.AlternativeMonoid[MyContext](stringMonoid)
//	eff1 := effect.Fail[MyContext, string](errors.New("failed"))
//	eff2 := effect.Of[MyContext]("fallback")
//	combined := effectMonoid.Concat(eff1, eff2)
//	// combined produces "fallback" (first failed, so second is used)
func AlternativeMonoid[C, A any](m monoid.Monoid[A]) Monoid[Effect[C, A]] {
	return readerreaderioresult.AlternativeMonoid[C](m)
}

// AltMonoid creates a monoid for effects using alternative semantics with a custom zero element.
// This tries the first effect, and if it fails, tries the second effect.
// The zero element is provided as a lazy computation.
//
// Type Parameters:
//   - R: The environment type required by the effects
//   - A: The value type produced by the effects
//
// Parameters:
//   - zero: A lazy computation that produces the zero/empty effect
//
// Returns:
//   - Monoid[Effect[R, A]]: A monoid for combining effects with custom zero
//
// Example:
//
//	import (
//	    "errors"
//	    L "github.com/IBM/fp-go/lazy"
//	)
//
//	// Create a monoid with a custom zero that returns a default value
//	zero := L.Of(Of[string]("default"))
//	effectMonoid := AltMonoid(zero)
//
//	// Empty returns the zero effect
//	empty := effectMonoid.Empty()
//	result := empty("env") // Returns Result containing "default"
//
//	// Concat tries first effect, falls back to second if first fails
//	eff1 := Fail[string, string](errors.New("failed"))
//	eff2 := Of[string]("fallback")
//	combined := effectMonoid.Concat(eff1, eff2)
//	// combined produces Result containing "fallback"
//
// See Also:
//   - ApplicativeMonoid: Combines effects by running both and combining results
//   - AlternativeMonoid: Alternative semantics with standard monoid for values
func AltMonoid[R, A any](zero Lazy[Effect[R, A]]) Monoid[Effect[R, A]] {
	return readerreaderioresult.AltMonoid(zero)
}
