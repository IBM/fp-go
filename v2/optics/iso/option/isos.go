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

// Package option provides isomorphisms for working with Option types.
// It offers utilities to convert between regular values and Option-wrapped values,
// particularly useful for handling zero values and optional data.
package option

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/iso"
	"github.com/IBM/fp-go/v2/option"
)

// FromZero creates an isomorphism between a comparable type T and Option[T].
// The isomorphism treats the zero value of T as None and non-zero values as Some.
//
// This is particularly useful for types where the zero value has special meaning
// (e.g., 0 for numbers, "" for strings, nil for pointers) and you want to represent
// the absence of a meaningful value using Option.
//
// Type Parameters:
//   - T: A comparable type (must support == and != operators)
//
// Returns:
//   - An Iso[T, Option[T]] where:
//   - Get: Converts T to Option[T] (zero value → None, non-zero → Some)
//   - ReverseGet: Converts Option[T] to T (None → zero value, Some → unwrapped value)
//
// Behavior:
//   - Get direction: If the value equals the zero value of T, returns None; otherwise returns Some(value)
//   - ReverseGet direction: If the Option is None, returns the zero value; otherwise returns the unwrapped value
//
// Example with integers:
//
//	isoInt := FromZero[int]()
//	opt := isoInt.Get(0)        // None (0 is the zero value)
//	opt = isoInt.Get(42)        // Some(42)
//	val := isoInt.ReverseGet(option.None[int]())  // 0
//	val = isoInt.ReverseGet(option.Some(42))      // 42
//
// Example with strings:
//
//	isoStr := FromZero[string]()
//	opt := isoStr.Get("")       // None ("" is the zero value)
//	opt = isoStr.Get("hello")   // Some("hello")
//	val := isoStr.ReverseGet(option.None[string]())  // ""
//	val = isoStr.ReverseGet(option.Some("world"))    // "world"
//
// Example with pointers:
//
//	isoPtr := FromZero[*int]()
//	opt := isoPtr.Get(nil)      // None (nil is the zero value)
//	num := 42
//	opt = isoPtr.Get(&num)      // Some(&num)
//
// Use cases:
//   - Converting between database nullable columns and Go types
//   - Handling optional configuration values with defaults
//   - Working with APIs that use zero values to indicate absence
//   - Simplifying validation logic for required vs optional fields
//
// Note: This isomorphism satisfies the round-trip laws:
//   - ReverseGet(Get(t)) == t for all t: T
//   - Get(ReverseGet(opt)) == opt for all opt: Option[T]
func FromZero[T comparable]() iso.Iso[T, option.Option[T]] {
	var zero T
	return iso.MakeIso(
		option.FromPredicate(func(t T) bool { return t != zero }),
		option.GetOrElse(F.Constant(zero)),
	)
}
