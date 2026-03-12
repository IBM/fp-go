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

package codec

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/decode"
)

// FromIso creates a Type codec from an Iso (isomorphism).
//
// An isomorphism represents a bidirectional transformation between types I and A
// without any loss of information. This function converts an Iso[I, A] into a
// Type[A, I, I] codec that can validate, decode, and encode values using the
// isomorphism's transformations.
//
// The resulting codec:
//   - Decode: Uses iso.Get to transform I → A, always succeeds (no validation)
//   - Encode: Uses iso.ReverseGet to transform A → I
//   - Validation: Always succeeds since isomorphisms are lossless transformations
//   - Type checking: Uses standard type checking for type A
//
// This is particularly useful for:
//   - Creating codecs for newtype patterns (wrapping/unwrapping types)
//   - Building codecs for types with lossless conversions
//   - Composing with other codecs using Pipe or other operators
//   - Implementing bidirectional transformations in codec pipelines
//
// # Type Parameters
//
//   - A: The target type (what we decode to and encode from)
//   - I: The input/output type (what we decode from and encode to)
//
// # Parameters
//
//   - iso: An Iso[I, A] that defines the bidirectional transformation:
//   - Get: I → A (converts input to target type)
//   - ReverseGet: A → I (converts target back to input type)
//
// # Returns
//
//   - A Type[A, I, I] codec where:
//   - Decode: I → Validation[A] - transforms using iso.Get, always succeeds
//   - Encode: A → I - transforms using iso.ReverseGet
//   - Is: Checks if a value is of type A
//   - Name: Returns "FromIso[iso_string_representation]"
//
// # Behavior
//
// Decoding:
//   - Applies iso.Get to transform the input value
//   - Wraps the result in decode.Of (always successful validation)
//   - No validation errors can occur since isomorphisms are lossless
//
// Encoding:
//   - Applies iso.ReverseGet to transform back to the input type
//   - Always succeeds as isomorphisms guarantee reversibility
//
// # Example Usage
//
// Creating a codec for a newtype pattern:
//
//	type UserId int
//
//	// Define an isomorphism between int and UserId
//	userIdIso := iso.MakeIso(
//	    func(id UserId) int { return int(id) },
//	    func(i int) UserId { return UserId(i) },
//	)
//
//	// Create a codec from the isomorphism
//	userIdCodec := codec.FromIso[int, UserId](userIdIso)
//
//	// Decode: UserId → int
//	result := userIdCodec.Decode(UserId(42))  // Success: Right(42)
//
//	// Encode: int → UserId
//	encoded := userIdCodec.Encode(42)         // Returns: UserId(42)
//
// Using with temperature conversions:
//
//	type Celsius float64
//	type Fahrenheit float64
//
//	celsiusToFahrenheit := iso.MakeIso(
//	    func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
//	    func(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) },
//	)
//
//	tempCodec := codec.FromIso[Fahrenheit, Celsius](celsiusToFahrenheit)
//
//	// Decode: Celsius → Fahrenheit
//	result := tempCodec.Decode(Celsius(20))   // Success: Right(68°F)
//
//	// Encode: Fahrenheit → Celsius
//	encoded := tempCodec.Encode(Fahrenheit(68)) // Returns: 20°C
//
// Composing with other codecs:
//
//	type Email string
//	type ValidatedEmail struct{ value Email }
//
//	emailIso := iso.MakeIso(
//	    func(ve ValidatedEmail) Email { return ve.value },
//	    func(e Email) ValidatedEmail { return ValidatedEmail{value: e} },
//	)
//
//	// Compose with string codec for validation
//	emailCodec := F.Pipe2(
//	    codec.String(),                           // Type[string, string, any]
//	    codec.Pipe(codec.FromIso[Email, string](  // Add string → Email iso
//	        iso.MakeIso(
//	            func(s string) Email { return Email(s) },
//	            func(e Email) string { return string(e) },
//	        ),
//	    )),
//	    codec.Pipe(codec.FromIso[ValidatedEmail, Email](emailIso)), // Add Email → ValidatedEmail iso
//	)
//
// # Use Cases
//
//   - Newtype patterns: Wrapping primitive types for type safety
//   - Unit conversions: Temperature, distance, time, etc.
//   - Format transformations: Between equivalent representations
//   - Type aliasing: Creating semantic types from base types
//   - Codec composition: Building complex codecs from simple isomorphisms
//
// # Notes
//
//   - Isomorphisms must satisfy the round-trip laws:
//   - iso.ReverseGet(iso.Get(i)) == i
//   - iso.Get(iso.ReverseGet(a)) == a
//   - Validation always succeeds since isomorphisms are lossless
//   - The codec name includes the isomorphism's string representation
//   - Type checking is performed using the standard Is[A]() function
//   - This codec is ideal for lossless transformations without validation logic
//
// # See Also
//
//   - iso.Iso: The isomorphism type used by this function
//   - iso.MakeIso: Constructor for creating isomorphisms
//   - Pipe: For composing this codec with other codecs
//   - MakeType: For creating codecs with custom validation logic
func FromIso[A, I any](iso Iso[I, A]) Type[A, I, I] {
	return MakeType(
		fmt.Sprintf("FromIso[%s]", iso),
		Is[A](),
		F.Flow2(
			iso.Get,
			decode.Of[Context],
		),
		iso.ReverseGet,
	)
}
