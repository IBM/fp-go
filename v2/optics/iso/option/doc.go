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
Package option provides isomorphisms for working with Option types.

# Overview

This package offers utilities to convert between regular values and Option-wrapped values,
particularly useful for handling zero values and optional data. It provides isomorphisms
that treat certain values (like zero values) as representing absence, mapping them to None,
while other values map to Some.

# Core Functionality

The main function in this package is FromZero, which creates an isomorphism between a
comparable type T and Option[T], treating the zero value as None.

# FromZero Isomorphism

FromZero creates a bidirectional transformation where:
  - Forward (Get): T → Option[T]
  - Zero value → None
  - Non-zero value → Some(value)
  - Reverse (ReverseGet): Option[T] → T
  - None → Zero value
  - Some(value) → value

# Basic Usage

Working with integers:

	import (
		"github.com/IBM/fp-go/v2/optics/iso/option"
		O "github.com/IBM/fp-go/v2/option"
	)

	isoInt := option.FromZero[int]()

	// Convert zero to None
	opt := isoInt.Get(0) // None[int]

	// Convert non-zero to Some
	opt = isoInt.Get(42) // Some(42)

	// Convert None to zero
	val := isoInt.ReverseGet(O.None[int]()) // 0

	// Convert Some to value
	val = isoInt.ReverseGet(O.Some(42)) // 42

# Use Cases

## Database Nullable Columns

Convert between database NULL and Go zero values:

	type User struct {
		ID       int
		Name     string
		Age      *int  // NULL in database
		Email    *string
	}

	ageIso := option.FromZero[*int]()

	// Reading from database
	var dbAge *int = nil
	optAge := ageIso.Get(dbAge) // None[*int]

	// Writing to database
	userAge := 25
	dbAge = ageIso.ReverseGet(O.Some(&userAge)) // &25

## Configuration with Defaults

Handle optional configuration values:

	type Config struct {
		Port    int
		Timeout int
		MaxConn int
	}

	portIso := option.FromZero[int]()

	// Use zero as "not configured"
	config := Config{Port: 0, Timeout: 30, MaxConn: 100}
	portOpt := portIso.Get(config.Port) // None[int] (use default)

	// Set explicit value
	config.Port = portIso.ReverseGet(O.Some(8080)) // 8080

## API Response Handling

Work with APIs that use zero values to indicate absence:

	type APIResponse struct {
		UserID   int     // 0 means not set
		Score    float64 // 0.0 means not available
		Message  string  // "" means no message
	}

	userIDIso := option.FromZero[int]()
	scoreIso := option.FromZero[float64]()
	messageIso := option.FromZero[string]()

	response := APIResponse{UserID: 0, Score: 0.0, Message: ""}

	userID := userIDIso.Get(response.UserID)     // None[int]
	score := scoreIso.Get(response.Score)        // None[float64]
	message := messageIso.Get(response.Message)  // None[string]

## Validation Logic

Simplify required vs optional field validation:

	import S "github.com/IBM/fp-go/v2/string"

	type FormData struct {
		Name     string // Required
		Email    string // Required
		Phone    string // Optional (empty = not provided)
		Comments string // Optional
	}

	phoneIso := option.FromZero[string]()
	commentsIso := option.FromZero[string]()

	form := FormData{
		Name:     "Alice",
		Email:    "alice@example.com",
		Phone:    "",
		Comments: "",
	}

	// Check optional fields
	phone := phoneIso.Get(form.Phone)       // None[string]
	comments := commentsIso.Get(form.Comments) // None[string]

	// Validate: required fields must be non-empty
	if S.IsEmpty(form.Name) || S.IsEmpty(form.Email) {
		// Validation error
	}

# Working with Different Types

## Strings

	strIso := option.FromZero[string]()

	opt := strIso.Get("")        // None[string]
	opt = strIso.Get("hello")    // Some("hello")

	val := strIso.ReverseGet(O.None[string]())      // ""
	val = strIso.ReverseGet(O.Some("world"))        // "world"

## Pointers

	ptrIso := option.FromZero[*int]()

	opt := ptrIso.Get(nil)       // None[*int]
	num := 42
	opt = ptrIso.Get(&num)       // Some(&num)

	val := ptrIso.ReverseGet(O.None[*int]())  // nil
	val = ptrIso.ReverseGet(O.Some(&num))     // &num

## Floating Point Numbers

	floatIso := option.FromZero[float64]()

	opt := floatIso.Get(0.0)     // None[float64]
	opt = floatIso.Get(3.14)     // Some(3.14)

	val := floatIso.ReverseGet(O.None[float64]())  // 0.0
	val = floatIso.ReverseGet(O.Some(2.71))        // 2.71

## Booleans

	boolIso := option.FromZero[bool]()

	opt := boolIso.Get(false)    // None[bool]
	opt = boolIso.Get(true)      // Some(true)

	val := boolIso.ReverseGet(O.None[bool]())   // false
	val = boolIso.ReverseGet(O.Some(true))      // true

# Composition with Other Optics

Combine with lenses for nested structures:

	import (
		L "github.com/IBM/fp-go/v2/optics/lens"
		I "github.com/IBM/fp-go/v2/optics/iso"
	)

	type Settings struct {
		Volume int // 0 means muted
	}

	volumeLens := L.MakeLens(
		func(s Settings) int { return s.Volume },
		func(s Settings, v int) Settings {
			s.Volume = v
			return s
		},
	)

	volumeIso := option.FromZero[int]()

	// Compose lens with iso
	volumeOptLens := F.Pipe1(
		volumeLens,
		L.IMap[Settings](volumeIso.Get, volumeIso.ReverseGet),
	)

	settings := Settings{Volume: 0}
	vol := volumeOptLens.Get(settings) // None[int] (muted)

	// Set volume
	updated := volumeOptLens.Set(O.Some(75))(settings)
	// updated.Volume == 75

# Isomorphism Laws

FromZero satisfies the isomorphism round-trip laws:

1. **ReverseGet(Get(t)) == t** for all t: T

	isoInt := option.FromZero[int]()
	value := 42
	result := isoInt.ReverseGet(isoInt.Get(value))
	// result == 42

2. **Get(ReverseGet(opt)) == opt** for all opt: Option[T]

	isoInt := option.FromZero[int]()
	opt := O.Some(42)
	result := isoInt.Get(isoInt.ReverseGet(opt))
	// result == Some(42)

These laws ensure that the transformation is truly reversible with no information loss.

# Performance Considerations

The FromZero isomorphism is very efficient:
  - No allocations for the iso structure itself
  - Simple equality comparison for zero check
  - Direct value unwrapping for ReverseGet
  - No reflection or runtime type assertions

# Type Safety

The isomorphism is fully type-safe:
  - Compile-time type checking ensures T is comparable
  - Generic type parameters prevent type mismatches
  - No runtime type assertions needed
  - The compiler enforces correct usage

# Limitations

The FromZero isomorphism has some limitations to be aware of:

1. **Zero Value Ambiguity**: Cannot distinguish between "intentionally zero" and "absent"
  - For int: 0 always maps to None, even if 0 is a valid value
  - For string: "" always maps to None, even if empty string is valid
  - Solution: Use a different representation (e.g., pointers) if zero is meaningful

2. **Comparable Constraint**: Only works with comparable types
  - Cannot use with slices, maps, or functions
  - Cannot use with structs containing non-comparable fields
  - Solution: Use pointers to such types, or custom isomorphisms

3. **Boolean Limitation**: false always maps to None
  - Cannot represent "explicitly false" vs "not set"
  - Solution: Use *bool or a custom type if this distinction matters

# Related Packages

  - github.com/IBM/fp-go/v2/optics/iso: Core isomorphism functionality
  - github.com/IBM/fp-go/v2/option: Option type and operations
  - github.com/IBM/fp-go/v2/optics/lens: Lenses for focused access
  - github.com/IBM/fp-go/v2/optics/lens/option: Lenses for optional values

# See Also

For more information on isomorphisms and optics:
  - optics/iso package documentation
  - optics package overview
  - option package documentation
*/
package option
