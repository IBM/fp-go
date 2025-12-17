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

/*
Package iso provides isomorphisms - bidirectional transformations between types without loss of information.

# Overview

An Isomorphism (Iso) is an optic that converts elements of type S into elements of type A
and back again without any loss of information. Unlike lenses which focus on a part of a
structure, isomorphisms represent a complete, reversible transformation between two types.

Isomorphisms are useful for:
  - Converting between equivalent representations (e.g., Celsius ↔ Fahrenheit)
  - Wrapping and unwrapping newtypes
  - Encoding and decoding data formats
  - Normalizing data structures

# Mathematical Foundation

An Iso[S, A] consists of two functions:
  - Get: S → A (convert from S to A)
  - ReverseGet: A → S (convert from A back to S)

Isomorphisms must satisfy the round-trip laws:
 1. ReverseGet(Get(s)) == s (for all s: S)
 2. Get(ReverseGet(a)) == a (for all a: A)

These laws ensure that the transformation is truly reversible with no information loss.

# Basic Usage

Creating an isomorphism between Celsius and Fahrenheit:

	celsiusToFahrenheit := iso.MakeIso(
		func(c float64) float64 { return c*9/5 + 32 },
		func(f float64) float64 { return (f - 32) * 5 / 9 },
	)

	// Convert Celsius to Fahrenheit
	fahrenheit := celsiusToFahrenheit.Get(20.0) // 68.0

	// Convert back to Celsius
	celsius := celsiusToFahrenheit.ReverseGet(68.0) // 20.0

# Identity Isomorphism

The identity isomorphism represents no transformation:

	idIso := iso.Id[int]()

	value := idIso.Get(42)        // 42
	same := idIso.ReverseGet(42)  // 42

# Composing Isomorphisms

Isomorphisms can be composed to create more complex transformations:

	metersToKm := iso.MakeIso(
		func(m float64) float64 { return m / 1000 },
		func(km float64) float64 { return km * 1000 },
	)

	kmToMiles := iso.MakeIso(
		func(km float64) float64 { return km * 0.621371 },
		func(mi float64) float64 { return mi / 0.621371 },
	)

	// Compose: meters → kilometers → miles
	metersToMiles := F.Pipe1(
		metersToKm,
		iso.Compose[float64](kmToMiles),
	)

	miles := metersToMiles.Get(5000)        // ~3.11 miles
	meters := metersToMiles.ReverseGet(3.11) // ~5000 meters

# Reversing Isomorphisms

Any isomorphism can be reversed to swap the direction:

	fahrenheitToCelsius := iso.Reverse(celsiusToFahrenheit)

	celsius := fahrenheitToCelsius.Get(68.0)        // 20.0
	fahrenheit := fahrenheitToCelsius.ReverseGet(20.0) // 68.0

# Modifying Through Isomorphisms

Apply transformations in the target space and convert back:

	type Meters float64
	type Kilometers float64

	mToKm := iso.MakeIso(
		func(m Meters) Kilometers { return Kilometers(m / 1000) },
		func(km Kilometers) Meters { return Meters(km * 1000) },
	)

	// Double the distance in kilometers, result in meters
	doubled := iso.Modify[Meters](func(km Kilometers) Kilometers {
		return km * 2
	})(mToKm)(Meters(5000))
	// Result: 10000 meters

# Wrapping and Unwrapping

Convenient functions for working with newtypes:

	type UserId int
	type User struct {
		id UserId
	}

	userIdIso := iso.MakeIso(
		func(id UserId) int { return int(id) },
		func(i int) UserId { return UserId(i) },
	)

	// Unwrap (Get)
	rawId := iso.Unwrap[int](UserId(42))(userIdIso) // 42
	// Also available as: iso.To[int](UserId(42))(userIdIso)

	// Wrap (ReverseGet)
	userId := iso.Wrap[UserId](42)(userIdIso) // UserId(42)
	// Also available as: iso.From[UserId](42)(userIdIso)

# Bidirectional Mapping

Transform both directions of an isomorphism:

	type Celsius float64
	type Kelvin float64

	celsiusIso := iso.Id[Celsius]()

	// Create isomorphism to Kelvin
	celsiusToKelvin := F.Pipe1(
		celsiusIso,
		iso.IMap(
			func(c Celsius) Kelvin { return Kelvin(c + 273.15) },
			func(k Kelvin) Celsius { return Celsius(k - 273.15) },
		),
	)

	kelvin := celsiusToKelvin.Get(Celsius(20))      // 293.15 K
	celsius := celsiusToKelvin.ReverseGet(Kelvin(293.15)) // 20°C

# Real-World Example: Data Encoding

	type JSON string
	type User struct {
		Name string
		Age  int
	}

	userJsonIso := iso.MakeIso(
		func(u User) JSON {
			data, _ := json.Marshal(u)
			return JSON(data)
		},
		func(j JSON) User {
			var u User
			json.Unmarshal([]byte(j), &u)
			return u
		},
	)

	user := User{Name: "Alice", Age: 30}

	// Encode to JSON
	jsonData := userJsonIso.Get(user)
	// `{"Name":"Alice","Age":30}`

	// Decode from JSON
	decoded := userJsonIso.ReverseGet(jsonData)
	// User{Name: "Alice", Age: 30}

# Real-World Example: Unit Conversions

	type Distance float64
	type DistanceUnit int

	const (
		Meters DistanceUnit = iota
		Kilometers
		Miles
	)

	type MeasuredDistance struct {
		Value Distance
		Unit  DistanceUnit
	}

	// Normalize all distances to meters
	normalizeIso := iso.MakeIso(
		func(md MeasuredDistance) Distance {
			switch md.Unit {
			case Kilometers:
				return md.Value * 1000
			case Miles:
				return md.Value * 1609.34
			default:
				return md.Value
			}
		},
		func(d Distance) MeasuredDistance {
			return MeasuredDistance{Value: d, Unit: Meters}
		},
	)

	// Convert 5 km to meters
	meters := normalizeIso.Get(MeasuredDistance{
		Value: 5,
		Unit:  Kilometers,
	}) // 5000 meters

	// Convert back (always to meters)
	measured := normalizeIso.ReverseGet(5000)
	// MeasuredDistance{Value: 5000, Unit: Meters}

# Real-World Example: Newtype Pattern

	type Email string
	type ValidatedEmail struct {
		value Email
	}

	func validateEmail(s string) (ValidatedEmail, error) {
		if !strings.Contains(s, "@") {
			return ValidatedEmail{}, errors.New("invalid email")
		}
		return ValidatedEmail{value: Email(s)}, nil
	}

	// Note: This iso assumes validation has already occurred
	emailIso := iso.MakeIso(
		func(ve ValidatedEmail) Email { return ve.value },
		func(e Email) ValidatedEmail { return ValidatedEmail{value: e} },
	)

	validated := ValidatedEmail{value: "user@example.com"}

	// Extract raw email
	raw := emailIso.Get(validated) // "user@example.com"

	// Wrap back (assumes valid)
	wrapped := emailIso.ReverseGet(Email("admin@example.com"))

# Isomorphisms vs Lenses

While both are optics, they serve different purposes:

**Isomorphisms:**
  - Represent complete, reversible transformations
  - No information loss
  - Both directions are equally important
  - Example: Celsius ↔ Fahrenheit

**Lenses:**
  - Focus on a part of a larger structure
  - Information loss when setting (other fields unchanged)
  - Asymmetric (get vs set)
  - Example: Person → Name

# Performance Considerations

Isomorphisms are lightweight and have minimal overhead:
  - No allocations for the iso structure itself
  - Performance depends on the Get and ReverseGet functions
  - Composition creates new function closures but is still efficient
  - Consider caching results if transformations are expensive

# Type Safety

Isomorphisms are fully type-safe:
  - The compiler ensures Get and ReverseGet have compatible types
  - Composition maintains type relationships
  - No runtime type assertions needed

# Function Reference

Core Functions:
  - MakeIso: Create an isomorphism from two functions
  - Id: Create an identity isomorphism
  - Compose: Compose two isomorphisms
  - Reverse: Reverse the direction of an isomorphism

Transformation:
  - Modify: Apply a transformation in the target space
  - IMap: Bidirectionally map an isomorphism

Convenience Functions:
  - Unwrap/To: Extract the target value (Get)
  - Wrap/From: Wrap into the source value (ReverseGet)

# Useful Iso Implementations

The package provides several ready-to-use isomorphisms for common transformations:

**String and Byte Conversions:**
  - UTF8String: []byte ↔ string (UTF-8 encoding)
  - Lines: []string ↔ string (newline-separated text)

**Time Conversions:**
  - UnixMilli: int64 ↔ time.Time (Unix millisecond timestamps)

**Numeric Operations:**
  - Add[T]: T ↔ T (shift by constant addition)
  - Sub[T]: T ↔ T (shift by constant subtraction)

**Collection Operations:**
  - ReverseArray[A]: []A ↔ []A (reverse slice order, self-inverse)
  - Head[A]: A ↔ NonEmptyArray[A] (singleton array conversion)

**Pair and Either Operations:**
  - SwapPair[A, B]: Pair[A, B] ↔ Pair[B, A] (swap pair elements, self-inverse)
  - SwapEither[E, A]: Either[E, A] ↔ Either[A, E] (swap Either types, self-inverse)

**Option Conversions (optics/iso/option):**
  - FromZero[T]: T ↔ Option[T] (zero value ↔ None, non-zero ↔ Some)

**Lens Conversions (optics/iso/lens):**
  - IsoAsLens: Convert Iso[S, A] to Lens[S, A]
  - IsoAsLensRef: Convert Iso[*S, A] to Lens[*S, A]

Example usage of built-in isomorphisms:

	// String/byte conversion
	utf8 := UTF8String()
	str := utf8.Get([]byte("hello"))  // "hello"

	// Time conversion
	unixTime := UnixMilli()
	t := unixTime.Get(1609459200000)  // 2021-01-01 00:00:00 UTC

	// Numeric shift
	addTen := Add(10)
	result := addTen.Get(5)  // 15

	// Array reversal
	reverse := ReverseArray[int]()
	reversed := reverse.Get([]int{1, 2, 3})  // []int{3, 2, 1}

	// Pair swap
	swap := SwapPair[string, int]()
	swapped := swap.Get(pair.MakePair("a", 1))  // Pair[int, string](1, "a")

	// Option conversion
	optIso := option.FromZero[int]()
	opt := optIso.Get(0)   // None
	opt = optIso.Get(42)   // Some(42)

# Related Packages

  - github.com/IBM/fp-go/v2/optics/lens: Lenses for focusing on parts of structures
  - github.com/IBM/fp-go/v2/optics/prism: Prisms for sum types
  - github.com/IBM/fp-go/v2/optics/optional: Optional optics
  - github.com/IBM/fp-go/v2/function: Function composition utilities
  - github.com/IBM/fp-go/v2/endomorphism: Endomorphisms (A → A functions)
*/
package iso
