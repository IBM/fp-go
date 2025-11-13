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
Package lens provides conversions from isomorphisms to lenses.

# Overview

This package bridges the gap between isomorphisms (bidirectional transformations)
and lenses (focused accessors). Since every isomorphism can be viewed as a lens,
this package provides functions to perform that conversion.

An isomorphism Iso[S, A] represents a lossless bidirectional transformation between
types S and A. A lens Lens[S, A] provides focused access to a part A within a
structure S. Since an isomorphism can transform the entire structure S to A and back,
it naturally forms a lens that focuses on the "whole as a part".

# Mathematical Foundation

Given an Iso[S, A] with:
  - Get: S → A (forward transformation)
  - ReverseGet: A → S (reverse transformation)

We can construct a Lens[S, A] with:
  - Get: S → A (same as iso's Get)
  - Set: A → S → S (implemented as: a => s => ReverseGet(a))

The lens laws are automatically satisfied because the isomorphism laws guarantee:
 1. GetSet: Set(Get(s))(s) == s (from iso's round-trip law)
 2. SetGet: Get(Set(a)(s)) == a (from iso's inverse law)
 3. SetSet: Set(a2)(Set(a1)(s)) == Set(a2)(s) (trivially true)

# Basic Usage

Converting an isomorphism to a lens:

	type Celsius float64
	type Kelvin float64

	// Create an isomorphism between Celsius and Kelvin
	celsiusKelvinIso := iso.MakeIso(
		func(c Celsius) Kelvin { return Kelvin(c + 273.15) },
		func(k Kelvin) Celsius { return Celsius(k - 273.15) },
	)

	// Convert to a lens
	celsiusKelvinLens := lens.IsoAsLens(celsiusKelvinIso)

	// Use as a lens
	celsius := Celsius(20.0)
	kelvin := celsiusKelvinLens.Get(celsius)        // 293.15 K
	updated := celsiusKelvinLens.Set(Kelvin(300))(celsius) // 26.85°C

# Working with Pointers

For pointer-based structures, use IsoAsLensRef:

	type UserId int
	type User struct {
		id   UserId
		name string
	}

	// Isomorphism between User pointer and UserId
	userIdIso := iso.MakeIso(
		func(u *User) UserId { return u.id },
		func(id UserId) *User { return &User{id: id, name: "Unknown"} },
	)

	// Convert to a reference lens
	userIdLens := lens.IsoAsLensRef(userIdIso)

	user := &User{id: 42, name: "Alice"}
	id := userIdLens.Get(user)                    // 42
	updated := userIdLens.Set(UserId(100))(user)  // New user with id 100

# Use Cases

1. Type Wrappers: Convert between newtype wrappers and their underlying types

	type Email string
	type ValidatedEmail struct{ value Email }

	emailIso := iso.MakeIso(
		func(ve ValidatedEmail) Email { return ve.value },
		func(e Email) ValidatedEmail { return ValidatedEmail{value: e} },
	)

	emailLens := lens.IsoAsLens(emailIso)

2. Unit Conversions: Work with different units of measurement

	type Meters float64
	type Feet float64

	metersFeetIso := iso.MakeIso(
		func(m Meters) Feet { return Feet(m * 3.28084) },
		func(f Feet) Meters { return Meters(f / 3.28084) },
	)

	distanceLens := lens.IsoAsLens(metersFeetIso)

3. Encoding/Decoding: Transform between different representations

	type JSON string
	type Config struct {
		Host string
		Port int
	}

	// Assuming encode/decode functions exist
	configIso := iso.MakeIso(encode, decode)
	configLens := lens.IsoAsLens(configIso)

# Composition

Lenses created from isomorphisms can be composed with other lenses:

	type Temperature struct {
		celsius Celsius
	}

	// Lens to access celsius field
	celsiusFieldLens := L.MakeLens(
		func(t Temperature) Celsius { return t.celsius },
		func(t Temperature, c Celsius) Temperature {
			t.celsius = c
			return t
		},
	)

	// Compose with iso-based lens to work with Kelvin
	tempKelvinLens := F.Pipe1(
		celsiusFieldLens,
		L.Compose[Temperature](celsiusKelvinLens),
	)

	temp := Temperature{celsius: 20}
	kelvin := tempKelvinLens.Get(temp)              // 293.15 K
	updated := tempKelvinLens.Set(Kelvin(300))(temp) // 26.85°C

# Comparison with Direct Lenses

While you can create a lens directly, using an isomorphism provides benefits:

1. Reusability: The isomorphism can be used in multiple contexts
2. Bidirectionality: The inverse transformation is explicitly available
3. Type Safety: Isomorphism laws ensure correctness
4. Composability: Isomorphisms compose naturally

Direct lens approach requires defining both get and set operations separately,
while the isomorphism approach defines the bidirectional transformation once
and converts it to a lens when needed.

# Performance Considerations

Converting an isomorphism to a lens has minimal overhead. The resulting lens
simply delegates to the isomorphism's Get and ReverseGet functions. However,
keep in mind:

1. Each Set operation performs a full transformation via ReverseGet
2. For pointer types, use IsoAsLensRef to ensure proper copying
3. The lens ignores the original structure in Set, using only the new value

# Function Reference

Conversion Functions:
  - IsoAsLens: Convert Iso[S, A] to Lens[S, A] for value types
  - IsoAsLensRef: Convert Iso[*S, A] to Lens[*S, A] for pointer types

# Related Packages

  - github.com/IBM/fp-go/v2/optics/iso: Isomorphisms (bidirectional transformations)
  - github.com/IBM/fp-go/v2/optics/lens: Lenses (focused accessors)
  - github.com/IBM/fp-go/v2/optics/lens/iso: Convert lenses to isomorphisms (inverse operation)
  - github.com/IBM/fp-go/v2/endomorphism: Endomorphisms (A → A functions)
  - github.com/IBM/fp-go/v2/function: Function composition utilities

# Examples

Complete example with type wrappers:

	type UserId int
	type Username string

	type User struct {
		id   UserId
		name Username
	}

	// Isomorphism for UserId
	userIdIso := iso.MakeIso(
		func(u User) UserId { return u.id },
		func(id UserId) User { return User{id: id, name: "Unknown"} },
	)

	// Isomorphism for Username
	usernameIso := iso.MakeIso(
		func(u User) Username { return u.name },
		func(name Username) User { return User{id: 0, name: name} },
	)

	// Convert to lenses
	idLens := lens.IsoAsLens(userIdIso)
	nameLens := lens.IsoAsLens(usernameIso)

	user := User{id: 42, name: "Alice"}

	// Access and modify through lenses
	id := idLens.Get(user)                      // 42
	name := nameLens.Get(user)                  // "Alice"
	renamed := nameLens.Set("Bob")(user)        // User{id: 0, name: "Bob"}
	reidentified := idLens.Set(UserId(100))(user) // User{id: 100, name: "Unknown"}

Note: When using Set with iso-based lenses, the entire structure is replaced
via ReverseGet, so other fields may be reset to default values. For partial
updates, use regular lenses instead.
*/
package lens
