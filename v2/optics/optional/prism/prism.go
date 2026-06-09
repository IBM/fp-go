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

// Package prism provides utilities for converting prisms to optionals and working with Option types.
//
// This package bridges the gap between prisms (which focus on sum types) and optionals
// (which focus on values that may not exist). The key functions allow you to:
//   - Convert any prism into an optional using AsOptional
//   - Focus on the Some variant of Option types using Some
//
// These conversions maintain the optional laws, ensuring that the resulting optionals
// behave correctly with respect to GetOption and Set operations.
package prism

import (
	F "github.com/IBM/fp-go/v2/function"
	OPT "github.com/IBM/fp-go/v2/optics/optional"
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
)

// AsOptional converts a Prism into an Optional.
//
// A Prism[S, A] focuses on a specific variant within a sum type S, providing:
//   - GetOption: Attempts to extract A from S (returns Option[A])
//   - ReverseGet: Constructs S from A (always succeeds)
//
// An Optional[S, A] focuses on a value that may not exist within S, providing:
//   - GetOption: Attempts to extract A from S (returns Option[A])
//   - Set: Updates A within S if it exists (no-op if it doesn't)
//
// The conversion works by:
//   - Using the prism's GetOption directly as the optional's GetOption
//   - Implementing Set using the prism's Set operation, which internally uses
//     GetOption to check if the value exists before updating
//
// The resulting Optional satisfies the three optional laws:
//
//  1. GetSet Law (No-op on None):
//     If GetOption(s) returns None, then Set(a)(s) returns s unchanged.
//     This is satisfied because the prism's Set operation checks GetOption
//     and only updates when it returns Some.
//
//     Formally: GetOption(s) = None => Set(a)(s) = s
//
//  2. SetGet Law (Get what you Set):
//     If GetOption(s) returns Some(_), then GetOption(Set(a)(s)) returns Some(a).
//     This is satisfied because the prism's Set operation replaces the focused
//     value with the new value when GetOption returns Some.
//
//     Formally: GetOption(s) = Some(_) => GetOption(Set(a)(s)) = Some(a)
//
//  3. SetSet Law (Last Set Wins):
//     Set(b)(Set(a)(s)) equals Set(b)(s).
//     This is satisfied because both operations check GetOption and only update
//     when it returns Some, with the prism's Set operation ensuring the last set wins.
//
//     Formally: Set(b)(Set(a)(s)) = Set(b)(s)
//
// Type Parameters:
//   - S: The source type (sum type)
//   - A: The focus type (variant within the sum type)
//
// Parameters:
//   - sa: A prism focusing on variant A within sum type S
//
// Returns:
//   - An Optional[S, A] that focuses on the same variant
//
// Example:
//
//	type Result interface{ isResult() }
//	type Success struct{ Value int }
//	type Failure struct{ Error string }
//
//	// Create a prism for the Success variant
//	successPrism := prism.MakePrism(
//	    func(r Result) Option[int] {
//	        if s, ok := r.(Success); ok {
//	            return Some(s.Value)
//	        }
//	        return None[int]()
//	    },
//	    func(v int) Result { return Success{Value: v} },
//	)
//
//	// Convert to optional
//	successOptional := AsOptional(successPrism)
//
//	// Use the optional
//	result := Success{Value: 42}
//	value := successOptional.GetOption(result)  // Some(42)
//	updated := successOptional.Set(100)(result) // Success{Value: 100}
//
//	// Set is no-op when GetOption returns None (Law 1)
//	failure := Failure{Error: "failed"}
//	unchanged := successOptional.Set(100)(failure) // failure (unchanged)
//
// See Also:
//   - Some: Focuses on the Some variant of Option types
//   - github.com/IBM/fp-go/v2/optics/prism for prism operations
//   - github.com/IBM/fp-go/v2/optics/optional for optional operations
func AsOptional[S, A any](sa P.Prism[S, A]) OPT.Optional[S, A] {
	return OPT.MakeOptional(
		sa.GetOption,
		func(s S, a A) S {
			return P.Set[S](a)(sa)(s)
		},
	)
}

// PrismSome creates a prism that focuses on the Some variant of an Option type.
//
// This prism provides:
//   - GetOption: Returns the Option itself (identity function)
//   - ReverseGet: Wraps a value in Some
//
// This is a building block used by the Some function to create optionals that
// focus on values within Option types.
//
// Type Parameters:
//   - A: The type of value within the Option
//
// Returns:
//   - A Prism[Option[A], A] that focuses on Some values
//
// Example:
//
//	prism := PrismSome[int]()
//
//	// GetOption returns the Option itself
//	opt := Some(42)
//	result := prism.GetOption(opt)  // Some(42)
//
//	// ReverseGet wraps in Some
//	wrapped := prism.ReverseGet(42)  // Some(42)
//
// See Also:
//   - Some: Uses this prism to create optionals for Option types
//   - github.com/IBM/fp-go/v2/prism.FromOption for the standard prism version
func PrismSome[A any]() P.Prism[O.Option[A], A] {
	return P.MakePrismWithName(F.Identity[O.Option[A]], O.Some[A], "PrismSome")
}

// Some creates an Optional that focuses on the Some variant of an Option within a structure.
//
// Given an Optional[S, Option[A]] that focuses on an Option field, this function
// returns an Optional[S, A] that focuses directly on the value within Some.
//
// This is useful when you have a structure containing an Option field and want to
// work with the value inside Some without manually unwrapping the Option.
//
// The conversion works by composing the provided optional with a prism that
// extracts values from Some. The resulting optional:
//   - Returns Some(a) from GetOption only when both the outer optional matches
//     and the inner Option is Some
//   - Performs Set only when both conditions are met (no-op otherwise)
//
// The resulting Optional satisfies the three optional laws:
//
//  1. GetSet Law (No-op on None):
//     If GetOption(s) returns None (either because the outer optional doesn't match
//     or the inner Option is None), then Set(a)(s) returns s unchanged.
//
//     Formally: GetOption(s) = None => Set(a)(s) = s
//
//  2. SetGet Law (Get what you Set):
//     If GetOption(s) returns Some(_), then GetOption(Set(a)(s)) returns Some(a).
//
//     Formally: GetOption(s) = Some(_) => GetOption(Set(a)(s)) = Some(a)
//
//  3. SetSet Law (Last Set Wins):
//     Set(b)(Set(a)(s)) equals Set(b)(s).
//
//     Formally: Set(b)(Set(a)(s)) = Set(b)(s)
//
// Type Parameters:
//   - S: The structure type
//   - A: The type of value within the Option
//
// Parameters:
//   - soa: An optional focusing on an Option[A] field within S
//
// Returns:
//   - An Optional[S, A] that focuses directly on values within Some
//
// Example:
//
//	type Config struct {
//	    Timeout Option[int]
//	}
//
//	// Create an optional for the Timeout field
//	timeoutOptional := optional.MakeOptional(
//	    func(c Config) Option[Option[int]] {
//	        return Some(c.Timeout)
//	    },
//	    func(c Config, opt Option[int]) Config {
//	        c.Timeout = opt
//	        return c
//	    },
//	)
//
//	// Focus on the value within Some
//	valueOptional := Some(timeoutOptional)
//
//	// Use the optional
//	config := Config{Timeout: Some(30)}
//	value := valueOptional.GetOption(config)  // Some(30)
//	updated := valueOptional.Set(60)(config)  // Config{Timeout: Some(60)}
//
//	// Set is no-op when inner Option is None (Law 1)
//	emptyConfig := Config{Timeout: None[int]()}
//	unchanged := valueOptional.Set(60)(emptyConfig)  // emptyConfig (unchanged)
//
// See Also:
//   - AsOptional: Converts prisms to optionals
//   - PrismSome: The underlying prism for Option types
//   - github.com/IBM/fp-go/v2/optics/optional.Compose for composing optionals
func Some[S, A any](soa OPT.Optional[S, O.Option[A]]) OPT.Optional[S, A] {
	return OPT.Compose[S](AsOptional(PrismSome[A]()))(soa)
}
