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

package chain

import (
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/functor"
)

// Chainable represents a type that supports sequential composition of computations,
// where each computation depends on the result of the previous one.
//
// Chainable extends Apply by adding the Chain method (also known as flatMap or bind),
// which allows for dependent sequencing of computations. Unlike Ap, which combines
// independent computations, Chain allows the structure of the second computation to
// depend on the value produced by the first.
//
// A Chainable must satisfy the following laws:
//
// Associativity:
//
//	Chain(f)(Chain(g)(m)) == Chain(x => Chain(f)(g(x)))(m)
//
// Type Parameters:
//   - A: The input value type
//   - B: The output value type after chaining
//   - HKTA: The higher-kinded type containing A
//   - HKTB: The higher-kinded type containing B
//   - HKTFAB: The higher-kinded type containing a function from A to B
//
// Example:
//
//	// Given a Chainable for Option
//	var c Chainable[int, string, Option[int], Option[string], Option[func(int) string]]
//	chainFn := c.Chain(func(x int) Option[string] {
//	  if x > 0 {
//	    return Some(strconv.Itoa(x))
//	  }
//	  return None[string]()
//	})
//	result := chainFn(Some(42)) // Returns Some("42")
type Chainable[A, B, HKTA, HKTB, HKTFAB any] interface {
	apply.Apply[A, B, HKTA, HKTB, HKTFAB]

	// Chain sequences computations where the second computation depends on the
	// value produced by the first.
	//
	// Takes a function that produces a new context-wrapped value based on the
	// unwrapped input, and returns a function that applies this to a context-wrapped
	// input, flattening the nested context.
	Chain(func(A) HKTB) func(HKTA) HKTB
}

// ToFunctor converts from [Chainable] to [functor.Functor]
func ToFunctor[A, B, HKTA, HKTB, HKTFAB any](ap Chainable[A, B, HKTA, HKTB, HKTFAB]) functor.Functor[A, B, HKTA, HKTB] {
	return ap
}

// ToApply converts from [Chainable] to [functor.Functor]
func ToApply[A, B, HKTA, HKTB, HKTFAB any](ap Chainable[A, B, HKTA, HKTB, HKTFAB]) apply.Apply[A, B, HKTA, HKTB, HKTFAB] {
	return ap
}

type (
	// Kleisli represents a Kleisli arrow - a function from A to a monadic value HKTB.
	// It's used for composing monadic computations where each step depends on the previous result.
	Kleisli[A, HKTB any] = func(A) HKTB

	// Operator represents a transformation from one monadic value to another.
	// It takes a value in context HKTA and produces a value in context HKTB.
	Operator[HKTA, HKTB any] = func(HKTA) HKTB

	ChainType[A, HKTA, HKTB any] = func(func(A) HKTB) func(HKTA) HKTB
)
