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

// Package lens provides utilities for composing prisms with lenses to create optionals.
//
// This package enables composition of a Prism (which focuses on a variant within a sum type)
// with a Lens (which focuses on a field within a product type) to create an Optional that
// combines both focusing operations.
package lens

import (
	F "github.com/IBM/fp-go/v2/function"
	OL "github.com/IBM/fp-go/v2/optics/optional/lens"
	OP "github.com/IBM/fp-go/v2/optics/optional/prism"
)

// Compose composes a Prism with a Lens to create an Optional.
//
// This composition allows you to first focus on a variant within a sum type (using a Prism),
// and then focus on a field within that variant (using a Lens). The result is an Optional
// because the Prism may not match the source value.
//
// The composition works by:
//  1. Converting the Prism to an Optional using AsOptional
//  2. Composing that Optional with the Lens (converted to Optional) using optional composition
//
// The resulting Optional satisfies the three optional laws:
//
//  1. GetSet Law (No-op on None):
//     If GetOption(s) returns None (prism doesn't match), then Set(b)(s) returns s unchanged.
//     This is satisfied because the prism's AsOptional conversion ensures Set is a no-op
//     when the prism doesn't match.
//
//     Formally: GetOption(s) = None => Set(b)(s) = s
//
//  2. SetGet Law (Get what you Set):
//     If GetOption(s) returns Some(_) (prism matches), then GetOption(Set(b)(s)) returns Some(b).
//     This is satisfied because both the prism-to-optional and lens-to-optional conversions
//     preserve this property, and optional composition maintains it.
//
//     Formally: GetOption(s) = Some(_) => GetOption(Set(b)(s)) = Some(b)
//
//  3. SetSet Law (Last Set Wins):
//     Set(c)(Set(b)(s)) equals Set(c)(s).
//     This is satisfied because optional composition preserves this property from both
//     the prism and lens components.
//
//     Formally: Set(c)(Set(b)(s)) = Set(c)(s)
//
// Type Parameters:
//   - S: The source type (sum type)
//   - A: The intermediate type (variant within the sum type)
//   - B: The target type (field within the variant)
//
// Parameters:
//   - l: A Lens[A, B] that focuses on field B within variant A
//
// Returns:
//   - A function that takes a Prism[S, A] and returns an Optional[S, B]
//
// Example:
//
//	type Result interface{ isResult() }
//	type Success struct{ Value int }
//	type Failure struct{ Error string }
//
//	func (Success) isResult() {}
//	func (Failure) isResult() {}
//
//	// Prism to focus on Success variant
//	successPrism := prism.MakePrism(
//	    func(r Result) option.Option[Success] {
//	        if s, ok := r.(Success); ok {
//	            return option.Some(s)
//	        }
//	        return option.None[Success]()
//	    },
//	    func(s Success) Result { return s },
//	)
//
//	// Lens to focus on Value field within Success
//	valueLens := lens.MakeLens(
//	    func(s Success) int { return s.Value },
//	    func(s Success, v int) Success { s.Value = v; return s },
//	)
//
//	// Compose to create Optional[Result, int]
//	resultValueOptional := Compose[Result, Success, int](valueLens)(successPrism)
//
//	// Use the optional
//	result := Success{Value: 42}
//	value := resultValueOptional.GetOption(result)  // Some(42)
//	updated := resultValueOptional.Set(100)(result) // Success{Value: 100}
//
//	// Set is no-op when prism doesn't match (Law 1)
//	failure := Failure{Error: "failed"}
//	unchanged := resultValueOptional.Set(100)(failure) // failure (unchanged)
//
// See Also:
//   - github.com/IBM/fp-go/v2/optics/optional/prism.AsOptional for prism-to-optional conversion
//   - github.com/IBM/fp-go/v2/optics/optional/lens.Compose for lens-optional composition
//   - github.com/IBM/fp-go/v2/optics/lens/prism for the inverse composition (lens then prism)
func Compose[S, A, B any](l Lens[A, B]) func(Prism[S, A]) Optional[S, B] {
	return F.Flow2(
		OP.AsOptional,
		OL.Compose[S](l),
	)
}
