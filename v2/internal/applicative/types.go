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

package applicative

import (
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

// Applicative represents a type that combines the ability to lift pure values into
// a context (Pointed) with the ability to apply wrapped functions to wrapped values (Apply).
//
// Applicative functors allow for function application lifted over a computational context,
// enabling multiple independent effects to be combined. This is the composition of Apply
// and Pointed, providing both the ability to create wrapped values and to apply wrapped
// functions.
//
// An Applicative must satisfy the following laws:
//
// Identity:
//
//	Ap(Of(identity))(v) == v
//
// Homomorphism:
//
//	Ap(Of(f))(Of(x)) == Of(f(x))
//
// Interchange:
//
//	Ap(Of(f))(u) == Ap(Map(f => f(y))(u))(Of(y))
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
//	// Given an Applicative for Option
//	var app Applicative[int, string, Option[int], Option[string], Option[func(int) string]]
//	value := app.Of(42) // Returns Some(42)
//	fn := app.Of(strconv.Itoa)
//	result := app.Ap(value)(fn) // Returns Some("42")
type Applicative[A, B, HKTA, HKTB, HKTFAB any] interface {
	apply.Apply[A, B, HKTA, HKTB, HKTFAB]
	pointed.Pointed[A, HKTA]
}

// ToFunctor converts from [Applicative] to [functor.Functor]
func ToFunctor[A, B, HKTA, HKTB, HKTFAB any](ap Applicative[A, B, HKTA, HKTB, HKTFAB]) functor.Functor[A, B, HKTA, HKTB] {
	return ap
}

// ToApply converts from [Applicative] to [apply.Apply]
func ToApply[A, B, HKTA, HKTB, HKTFAB any](ap Applicative[A, B, HKTA, HKTB, HKTFAB]) apply.Apply[A, B, HKTA, HKTB, HKTFAB] {
	return ap
}

// ToPointed converts from [Applicative] to [pointed.Pointed]
func ToPointed[A, B, HKTA, HKTB, HKTFAB any](ap Applicative[A, B, HKTA, HKTB, HKTFAB]) pointed.Pointed[A, HKTA] {
	return ap
}
