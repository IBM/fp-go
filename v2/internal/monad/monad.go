// Copyright (c) 2024 - 2025 IBM Corp.
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

package monad

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

// Monad represents the full monadic interface, combining Applicative and Chainable.
//
// A Monad provides the complete set of operations for working with computational contexts:
// - Map (from Functor): transform values within a context
// - Of (from Pointed): lift pure values into a context
// - Ap (from Apply): apply wrapped functions to wrapped values
// - Chain (from Chainable): sequence dependent computations
//
// Monads must satisfy the monad laws:
//
// Left Identity:
//
//	Chain(f)(Of(a)) == f(a)
//
// Right Identity:
//
//	Chain(Of)(m) == m
//
// Associativity:
//
//	Chain(g)(Chain(f)(m)) == Chain(x => Chain(g)(f(x)))(m)
//
// Type Parameters:
//   - A: The input value type
//   - B: The output value type
//   - HKTA: The higher-kinded type containing A
//   - HKTB: The higher-kinded type containing B
//   - HKTFAB: The higher-kinded type containing a function from A to B
//
// Example:
//
//	// Given a Monad for Option
//	var m Monad[int, string, Option[int], Option[string], Option[func(int) string]]
//
//	// Use Of to create a value
//	value := m.Of(42) // Some(42)
//
//	// Use Chain for dependent operations
//	chainFn := m.Chain(func(x int) Option[string] {
//	  if x > 0 {
//	    return Some(strconv.Itoa(x))
//	  }
//	  return None[string]()
//	})
//	result := chainFn(value) // Some("42")
type Monad[A, B, HKTA, HKTB, HKTFAB any] interface {
	applicative.Applicative[A, B, HKTA, HKTB, HKTFAB]
	chain.Chainable[A, B, HKTA, HKTB, HKTFAB]
}

// ToFunctor converts from [Monad] to [functor.Functor]
func ToFunctor[A, B, HKTA, HKTB, HKTFAB any](ap Monad[A, B, HKTA, HKTB, HKTFAB]) functor.Functor[A, B, HKTA, HKTB] {
	return ap
}

// ToApply converts from [Monad] to [apply.Apply]
func ToApply[A, B, HKTA, HKTB, HKTFAB any](ap Monad[A, B, HKTA, HKTB, HKTFAB]) apply.Apply[A, B, HKTA, HKTB, HKTFAB] {
	return ap
}

// ToPointed converts from [Monad] to [pointed.Pointed]
func ToPointed[A, B, HKTA, HKTB, HKTFAB any](ap Monad[A, B, HKTA, HKTB, HKTFAB]) pointed.Pointed[A, HKTA] {
	return ap
}

// ToApplicative converts from [Monad] to [applicative.Applicative]
func ToApplicative[A, B, HKTA, HKTB, HKTFAB any](ap Monad[A, B, HKTA, HKTB, HKTFAB]) applicative.Applicative[A, B, HKTA, HKTB, HKTFAB] {
	return ap
}

// ToChainable converts from [Monad] to [chain.Chainable]
func ToChainable[A, B, HKTA, HKTB, HKTFAB any](ap Monad[A, B, HKTA, HKTB, HKTFAB]) chain.Chainable[A, B, HKTA, HKTB, HKTFAB] {
	return ap
}
