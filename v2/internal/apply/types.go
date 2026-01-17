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

package apply

import (
	"github.com/IBM/fp-go/v2/internal/functor"
)

// Apply represents a Functor with the ability to apply a function wrapped in a context
// to a value wrapped in a context.
//
// Apply extends Functor by adding the Ap method, which allows for application of
// functions that are themselves wrapped in the same computational context. This enables
// independent computations to be combined.
//
// An Apply must satisfy the following laws:
//
// Composition:
//
//	Ap(Ap(Map(compose)(f))(g))(x) == Ap(f)(Ap(g)(x))
//
// Type Parameters:
//   - A: The input value type
//   - B: The output value type after function application
//   - HKTA: The higher-kinded type containing A
//   - HKTB: The higher-kinded type containing B
//   - HKTFAB: The higher-kinded type containing a function from A to B
//
// Example:
//
//	// Given an Apply for Option
//	var ap Apply[int, string, Option[int], Option[string], Option[func(int) string]]
//	fn := Some(strconv.Itoa)
//	applyFn := ap.Ap(Some(42))
//	result := applyFn(fn) // Returns Some("42")
type Apply[A, B, HKTA, HKTB, HKTFAB any] interface {
	functor.Functor[A, B, HKTA, HKTB]

	// Ap applies a function wrapped in a context to a value wrapped in a context.
	//
	// Takes a value in context (HKTA) and returns a function that takes a function
	// in context (HKTFAB) and produces a result in context (HKTB).
	Ap(HKTA) func(HKTFAB) HKTB
}

// ToFunctor converts from [Apply] to [functor.Functor]
func ToFunctor[A, B, HKTA, HKTB, HKTFAB any](ap Apply[A, B, HKTA, HKTB, HKTFAB]) functor.Functor[A, B, HKTA, HKTB] {
	return ap
}

type ApType[HKTA, HKTB, HKTFAB any] = func(HKTA) func(HKTFAB) HKTB
