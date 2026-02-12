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

package iso

import (
	"errors"
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/optics/iso"
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestComposeWithEitherPrism tests composing an isomorphism with an Either prism
func TestComposeWithEitherPrism(t *testing.T) {
	// Create an isomorphism between string and []byte
	stringBytesIso := I.MakeIso(
		func(s string) []byte { return []byte(s) },
		func(b []byte) string { return string(b) },
	)

	// Create a prism that extracts Right values from Either[error, string]
	rightPrism := P.FromEither[error, string]()

	// Compose them
	bytesPrism := Compose[E.Either[error, string]](stringBytesIso)(rightPrism)

	t.Run("GetOption extracts and transforms Right value", func(t *testing.T) {
		success := E.Right[error]("hello")
		result := bytesPrism.GetOption(success)

		assert.True(t, O.IsSome(result))
		bytes := O.GetOrElse(F.Constant([]byte{}))(result)
		assert.Equal(t, []byte("hello"), bytes)
	})

	t.Run("GetOption returns None for Left value", func(t *testing.T) {
		failure := E.Left[string](errors.New("error"))
		result := bytesPrism.GetOption(failure)

		assert.True(t, O.IsNone(result))
	})

	t.Run("ReverseGet constructs Either from transformed value", func(t *testing.T) {
		bytes := []byte("world")
		result := bytesPrism.ReverseGet(bytes)

		assert.True(t, E.IsRight(result))
		str := E.GetOrElse(func(error) string { return "" })(result)
		assert.Equal(t, "world", str)
	})

	t.Run("Round-trip through GetOption and ReverseGet", func(t *testing.T) {
		// Start with bytes
		original := []byte("test")

		// ReverseGet to create Either
		either := bytesPrism.ReverseGet(original)

		// GetOption to extract bytes back
		result := bytesPrism.GetOption(either)

		assert.True(t, O.IsSome(result))
		extracted := O.GetOrElse(F.Constant([]byte{}))(result)
		assert.Equal(t, original, extracted)
	})
}

// TestComposeWithOptionPrism tests composing an isomorphism with an Option prism
func TestComposeWithOptionPrism(t *testing.T) {
	// Create an isomorphism between int and string
	intStringIso := I.MakeIso(
		func(i int) string { return strconv.Itoa(i) },
		func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		},
	)

	// Create a prism that extracts Some values from Option[int]
	somePrism := P.FromOption[int]()

	// Compose them
	stringPrism := Compose[O.Option[int]](intStringIso)(somePrism)

	t.Run("GetOption extracts and transforms Some value", func(t *testing.T) {
		some := O.Some(42)
		result := stringPrism.GetOption(some)

		assert.True(t, O.IsSome(result))
		str := O.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, "42", str)
	})

	t.Run("GetOption returns None for None value", func(t *testing.T) {
		none := O.None[int]()
		result := stringPrism.GetOption(none)

		assert.True(t, O.IsNone(result))
	})

	t.Run("ReverseGet constructs Option from transformed value", func(t *testing.T) {
		str := "100"
		result := stringPrism.ReverseGet(str)

		assert.True(t, O.IsSome(result))
		num := O.GetOrElse(F.Constant(0))(result)
		assert.Equal(t, 100, num)
	})
}

// TestComposeWithCustomPrism tests composing with a custom prism
// Custom types for TestComposeWithCustomPrism
type Celsius float64
type Fahrenheit float64

type Temperature interface {
	isTemperature()
}

type CelsiusTemp struct {
	Value Celsius
}

func (c CelsiusTemp) isTemperature() {}

type KelvinTemp struct {
	Value float64
}

func (k KelvinTemp) isTemperature() {}

func TestComposeWithCustomPrism(t *testing.T) {
	// Isomorphism between Celsius and Fahrenheit
	tempIso := I.MakeIso(
		func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
		func(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) },
	)

	// Prism that extracts Celsius from Temperature
	celsiusPrism := P.MakePrism(
		func(t Temperature) O.Option[Celsius] {
			if ct, ok := t.(CelsiusTemp); ok {
				return O.Some(ct.Value)
			}
			return O.None[Celsius]()
		},
		func(c Celsius) Temperature {
			return CelsiusTemp{Value: c}
		},
	)

	// Compose to work with Fahrenheit
	fahrenheitPrism := Compose[Temperature](tempIso)(celsiusPrism)

	t.Run("GetOption extracts and converts Celsius to Fahrenheit", func(t *testing.T) {
		temp := CelsiusTemp{Value: 0}
		result := fahrenheitPrism.GetOption(temp)

		assert.True(t, O.IsSome(result))
		fahrenheit := O.GetOrElse(F.Constant(Fahrenheit(0)))(result)
		assert.InDelta(t, 32.0, float64(fahrenheit), 0.01)
	})

	t.Run("GetOption returns None for non-Celsius temperature", func(t *testing.T) {
		temp := KelvinTemp{Value: 273.15}
		result := fahrenheitPrism.GetOption(temp)

		assert.True(t, O.IsNone(result))
	})

	t.Run("ReverseGet constructs Temperature from Fahrenheit", func(t *testing.T) {
		fahrenheit := Fahrenheit(68)
		result := fahrenheitPrism.ReverseGet(fahrenheit)

		celsiusTemp, ok := result.(CelsiusTemp)
		assert.True(t, ok)
		assert.InDelta(t, 20.0, float64(celsiusTemp.Value), 0.01)
	})

	t.Run("Round-trip preserves value", func(t *testing.T) {
		original := Fahrenheit(100)

		// ReverseGet to create Temperature
		temp := fahrenheitPrism.ReverseGet(original)

		// GetOption to extract Fahrenheit back
		result := fahrenheitPrism.GetOption(temp)

		assert.True(t, O.IsSome(result))
		extracted := O.GetOrElse(F.Constant(Fahrenheit(0)))(result)
		assert.InDelta(t, float64(original), float64(extracted), 0.01)
	})
}

// TestComposeIdentityIso tests composing with an identity isomorphism
func TestComposeIdentityIso(t *testing.T) {
	// Identity isomorphism (no transformation)
	idIso := I.Id[string]()

	// Prism that extracts Right values
	rightPrism := P.FromEither[error, string]()

	// Compose with identity should not change behavior
	composedPrism := Compose[E.Either[error, string]](idIso)(rightPrism)

	t.Run("Composed prism behaves like original prism", func(t *testing.T) {
		success := E.Right[error]("test")

		// Original prism
		originalResult := rightPrism.GetOption(success)

		// Composed prism
		composedResult := composedPrism.GetOption(success)

		assert.Equal(t, originalResult, composedResult)
	})

	t.Run("ReverseGet produces same result", func(t *testing.T) {
		value := "test"

		// Original prism
		originalEither := rightPrism.ReverseGet(value)

		// Composed prism
		composedEither := composedPrism.ReverseGet(value)

		assert.Equal(t, originalEither, composedEither)
	})
}

// TestComposeChaining tests chaining multiple compositions
func TestComposeChaining(t *testing.T) {
	// First isomorphism: int to string
	intStringIso := I.MakeIso(
		func(i int) string { return strconv.Itoa(i) },
		func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		},
	)

	// Second isomorphism: string to []byte
	stringBytesIso := I.MakeIso(
		func(s string) []byte { return []byte(s) },
		func(b []byte) string { return string(b) },
	)

	// Prism that extracts Right values
	rightPrism := P.FromEither[error, int]()

	// Chain compositions: Either[error, int] -> int -> string -> []byte
	step1 := Compose[E.Either[error, int]](intStringIso)(rightPrism) // Prism[Either[error, int], string]
	step2 := Compose[E.Either[error, int]](stringBytesIso)(step1)    // Prism[Either[error, int], []byte]

	t.Run("Chained composition extracts and transforms correctly", func(t *testing.T) {
		either := E.Right[error](42)
		result := step2.GetOption(either)

		assert.True(t, O.IsSome(result))
		bytes := O.GetOrElse(F.Constant([]byte{}))(result)
		assert.Equal(t, []byte("42"), bytes)
	})

	t.Run("Chained composition ReverseGet works correctly", func(t *testing.T) {
		bytes := []byte("100")
		result := step2.ReverseGet(bytes)

		assert.True(t, E.IsRight(result))
		num := E.GetOrElse(func(error) int { return 0 })(result)
		assert.Equal(t, 100, num)
	})
}

// TestComposePrismLaws verifies that the composed prism satisfies prism laws
func TestComposePrismLaws(t *testing.T) {
	// Create an isomorphism
	iso := I.MakeIso(
		func(i int) string { return strconv.Itoa(i) },
		func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		},
	)

	// Create a prism
	prism := P.FromEither[error, int]()

	// Compose them
	composed := Compose[E.Either[error, int]](iso)(prism)

	t.Run("Law 1: GetOption(ReverseGet(b)) == Some(b)", func(t *testing.T) {
		value := "42"

		// ReverseGet then GetOption should return Some(value)
		either := composed.ReverseGet(value)
		result := composed.GetOption(either)

		assert.True(t, O.IsSome(result))
		extracted := O.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, value, extracted)
	})

	t.Run("Law 2: If GetOption(s) == Some(a), then GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		either := E.Right[error](100)

		// First GetOption
		firstResult := composed.GetOption(either)
		assert.True(t, O.IsSome(firstResult))

		// Extract the value
		value := O.GetOrElse(F.Constant(""))(firstResult)

		// ReverseGet then GetOption again
		reconstructed := composed.ReverseGet(value)
		secondResult := composed.GetOption(reconstructed)

		assert.True(t, O.IsSome(secondResult))
		finalValue := O.GetOrElse(F.Constant(""))(secondResult)
		assert.Equal(t, value, finalValue)
	})
}

// TestComposeWithEmptyValues tests edge cases with empty/zero values
func TestComposeWithEmptyValues(t *testing.T) {
	// Isomorphism that handles empty strings
	iso := I.MakeIso(
		func(s string) []byte { return []byte(s) },
		func(b []byte) string { return string(b) },
	)

	prism := P.FromEither[error, string]()
	composed := Compose[E.Either[error, string]](iso)(prism)

	t.Run("Empty string is handled correctly", func(t *testing.T) {
		either := E.Right[error]("")
		result := composed.GetOption(either)

		assert.True(t, O.IsSome(result))
		bytes := O.GetOrElse(F.Constant([]byte("default")))(result)
		assert.Equal(t, []byte{}, bytes)
	})

	t.Run("ReverseGet with empty bytes", func(t *testing.T) {
		bytes := []byte{}
		result := composed.ReverseGet(bytes)

		assert.True(t, E.IsRight(result))
		str := E.GetOrElse(func(error) string { return "default" })(result)
		assert.Equal(t, "", str)
	})
}
