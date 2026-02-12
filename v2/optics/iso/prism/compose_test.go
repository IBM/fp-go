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

package prism

import (
	"encoding/json"
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/optics/iso"
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestComposeWithEitherPrism tests composing a prism with an isomorphism using Either
func TestComposeWithEitherPrism(t *testing.T) {
	// Create a prism that extracts Right values from Either[error, string]
	rightPrism := P.FromEither[error, string]()

	// Create an isomorphism between []byte and Either[error, string]
	bytesEitherIso := I.MakeIso(
		func(b []byte) E.Either[error, string] {
			return E.Right[error](string(b))
		},
		func(e E.Either[error, string]) []byte {
			return []byte(E.GetOrElse(func(error) string { return "" })(e))
		},
	)

	// Compose them: Prism[Either, string] with Iso[[]byte, Either] -> Prism[[]byte, string]
	bytesPrism := Compose[[]byte](rightPrism)(bytesEitherIso)

	t.Run("GetOption extracts string from []byte", func(t *testing.T) {
		bytes := []byte("hello")
		result := bytesPrism.GetOption(bytes)

		assert.True(t, O.IsSome(result))
		str := O.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, "hello", str)
	})

	t.Run("ReverseGet constructs []byte from string", func(t *testing.T) {
		value := "world"
		result := bytesPrism.ReverseGet(value)

		assert.Equal(t, []byte("world"), result)
	})

	t.Run("Round-trip through GetOption and ReverseGet", func(t *testing.T) {
		original := "test"

		// ReverseGet to create []byte
		bytes := bytesPrism.ReverseGet(original)

		// GetOption to extract string back
		result := bytesPrism.GetOption(bytes)

		assert.True(t, O.IsSome(result))
		extracted := O.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, original, extracted)
	})
}

// TestComposeWithOptionPrism tests composing a prism with an isomorphism using Option
func TestComposeWithOptionPrism(t *testing.T) {
	// Create a prism that extracts Some values from Option[int]
	somePrism := P.FromOption[int]()

	// Create an isomorphism between string and Option[int]
	stringOptionIso := I.MakeIso(
		func(s string) O.Option[int] {
			i, err := strconv.Atoi(s)
			if err != nil {
				return O.None[int]()
			}
			return O.Some(i)
		},
		func(opt O.Option[int]) string {
			return strconv.Itoa(O.GetOrElse(F.Constant(0))(opt))
		},
	)

	// Compose them: Prism[Option, int] with Iso[string, Option] -> Prism[string, int]
	stringPrism := Compose[string](somePrism)(stringOptionIso)

	t.Run("GetOption extracts int from valid string", func(t *testing.T) {
		result := stringPrism.GetOption("42")

		assert.True(t, O.IsSome(result))
		num := O.GetOrElse(F.Constant(0))(result)
		assert.Equal(t, 42, num)
	})

	t.Run("GetOption returns None for invalid string", func(t *testing.T) {
		result := stringPrism.GetOption("invalid")

		assert.True(t, O.IsNone(result))
	})

	t.Run("ReverseGet constructs string from int", func(t *testing.T) {
		result := stringPrism.ReverseGet(100)

		assert.Equal(t, "100", result)
	})
}

// Custom types for testing
type Celsius float64
type Fahrenheit float64

type Temperature interface {
	isTemperature()
}

type CelsiusTemp struct {
	Value Celsius
}

func (c CelsiusTemp) isTemperature() {}

type FahrenheitTemp struct {
	Value Fahrenheit
}

func (f FahrenheitTemp) isTemperature() {}

// TestComposeWithCustomPrism tests composing with custom types
func TestComposeWithCustomPrism(t *testing.T) {
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

	// Isomorphism between Fahrenheit and Temperature
	fahrenheitTempIso := I.MakeIso(
		func(f Fahrenheit) Temperature {
			celsius := Celsius((f - 32) * 5 / 9)
			return CelsiusTemp{Value: celsius}
		},
		func(t Temperature) Fahrenheit {
			if ct, ok := t.(CelsiusTemp); ok {
				return Fahrenheit(ct.Value*9/5 + 32)
			}
			return 0
		},
	)

	// Compose: Prism[Temperature, Celsius] with Iso[Fahrenheit, Temperature] -> Prism[Fahrenheit, Celsius]
	fahrenheitPrism := Compose[Fahrenheit](celsiusPrism)(fahrenheitTempIso)

	t.Run("GetOption extracts Celsius from Fahrenheit", func(t *testing.T) {
		fahrenheit := Fahrenheit(68)
		result := fahrenheitPrism.GetOption(fahrenheit)

		assert.True(t, O.IsSome(result))
		celsius := O.GetOrElse(F.Constant(Celsius(0)))(result)
		assert.InDelta(t, 20.0, float64(celsius), 0.01)
	})

	t.Run("ReverseGet constructs Fahrenheit from Celsius", func(t *testing.T) {
		celsius := Celsius(20)
		result := fahrenheitPrism.ReverseGet(celsius)

		assert.InDelta(t, 68.0, float64(result), 0.01)
	})

	t.Run("Round-trip preserves value", func(t *testing.T) {
		original := Celsius(25)

		// ReverseGet to create Fahrenheit
		fahrenheit := fahrenheitPrism.ReverseGet(original)

		// GetOption to extract Celsius back
		result := fahrenheitPrism.GetOption(fahrenheit)

		assert.True(t, O.IsSome(result))
		extracted := O.GetOrElse(F.Constant(Celsius(0)))(result)
		assert.InDelta(t, float64(original), float64(extracted), 0.01)
	})
}

// TestComposeIdentityIso tests composing with an identity isomorphism
func TestComposeIdentityIso(t *testing.T) {
	// Prism that extracts Right values
	rightPrism := P.FromEither[error, string]()

	// Identity isomorphism on Either
	idIso := I.Id[E.Either[error, string]]()

	// Compose with identity should not change behavior
	composedPrism := Compose[E.Either[error, string]](rightPrism)(idIso)

	t.Run("Composed prism behaves like original prism", func(t *testing.T) {
		either := E.Right[error]("test")

		// Original prism
		originalResult := rightPrism.GetOption(either)

		// Composed prism
		composedResult := composedPrism.GetOption(either)

		assert.Equal(t, originalResult, composedResult)
	})

	t.Run("ReverseGet produces same result", func(t *testing.T) {
		value := "test"

		// Original prism
		originalResult := rightPrism.ReverseGet(value)

		// Composed prism
		composedResult := composedPrism.ReverseGet(value)

		assert.Equal(t, originalResult, composedResult)
	})
}

// TestComposeChaining tests chaining multiple compositions
func TestComposeChaining(t *testing.T) {
	// Prism: extracts Right values from Either[error, int]
	rightPrism := P.FromEither[error, int]()

	// Iso 1: string to Either[error, int]
	stringEitherIso := I.MakeIso(
		func(s string) E.Either[error, int] {
			i, err := strconv.Atoi(s)
			if err != nil {
				return E.Left[int](err)
			}
			return E.Right[error](i)
		},
		func(e E.Either[error, int]) string {
			return strconv.Itoa(E.GetOrElse(func(error) int { return 0 })(e))
		},
	)

	// Iso 2: []byte to string
	bytesStringIso := I.MakeIso(
		func(b []byte) string { return string(b) },
		func(s string) []byte { return []byte(s) },
	)

	// First composition: Prism[Either, int] with Iso[string, Either] -> Prism[string, int]
	step1 := Compose[string](rightPrism)(stringEitherIso)

	// Second composition: Prism[string, int] with Iso[[]byte, string] -> Prism[[]byte, int]
	step2 := Compose[[]byte](step1)(bytesStringIso)

	t.Run("Chained composition extracts correctly", func(t *testing.T) {
		bytes := []byte("42")
		result := step2.GetOption(bytes)

		assert.True(t, O.IsSome(result))
		num := O.GetOrElse(F.Constant(0))(result)
		assert.Equal(t, 42, num)
	})

	t.Run("Chained composition ReverseGet works correctly", func(t *testing.T) {
		num := 100
		result := step2.ReverseGet(num)

		assert.Equal(t, []byte("100"), result)
	})
}

// TestComposePrismLaws verifies that the composed prism satisfies prism laws
func TestComposePrismLaws(t *testing.T) {
	// Create a prism
	prism := P.FromEither[error, int]()

	// Create an isomorphism from string to Either[error, int]
	iso := I.MakeIso(
		func(s string) E.Either[error, int] {
			i, err := strconv.Atoi(s)
			if err != nil {
				return E.Left[int](err)
			}
			return E.Right[error](i)
		},
		func(e E.Either[error, int]) string {
			return strconv.Itoa(E.GetOrElse(func(error) int { return 0 })(e))
		},
	)

	// Compose them
	composed := Compose[string](prism)(iso)

	t.Run("Law 1: GetOption(ReverseGet(b)) == Some(b)", func(t *testing.T) {
		value := 42

		// ReverseGet then GetOption should return Some(value)
		source := composed.ReverseGet(value)
		result := composed.GetOption(source)

		assert.True(t, O.IsSome(result))
		extracted := O.GetOrElse(F.Constant(0))(result)
		assert.Equal(t, value, extracted)
	})

	t.Run("Law 2: If GetOption(s) == Some(a), then GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		source := "100"

		// First GetOption
		firstResult := composed.GetOption(source)
		assert.True(t, O.IsSome(firstResult))

		// Extract the value
		value := O.GetOrElse(F.Constant(0))(firstResult)

		// ReverseGet then GetOption again
		reconstructed := composed.ReverseGet(value)
		secondResult := composed.GetOption(reconstructed)

		assert.True(t, O.IsSome(secondResult))
		finalValue := O.GetOrElse(F.Constant(0))(secondResult)
		assert.Equal(t, value, finalValue)
	})
}

// TestComposeWithJSON tests a practical example with JSON parsing
func TestComposeWithJSON(t *testing.T) {
	type Config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	// Prism that extracts Config from []byte (via JSON parsing)
	configPrism := P.MakePrism(
		func(b []byte) O.Option[Config] {
			var cfg Config
			if err := json.Unmarshal(b, &cfg); err != nil {
				return O.None[Config]()
			}
			return O.Some(cfg)
		},
		func(cfg Config) []byte {
			b, _ := json.Marshal(cfg)
			return b
		},
	)

	// Isomorphism between string and []byte
	stringBytesIso := I.MakeIso(
		func(s string) []byte { return []byte(s) },
		func(b []byte) string { return string(b) },
	)

	// Compose: Prism[[]byte, Config] with Iso[string, []byte] -> Prism[string, Config]
	stringConfigPrism := Compose[string](configPrism)(stringBytesIso)

	t.Run("GetOption parses valid JSON string", func(t *testing.T) {
		jsonStr := `{"host":"localhost","port":8080}`
		result := stringConfigPrism.GetOption(jsonStr)

		assert.True(t, O.IsSome(result))
		cfg := O.GetOrElse(F.Constant(Config{}))(result)
		assert.Equal(t, "localhost", cfg.Host)
		assert.Equal(t, 8080, cfg.Port)
	})

	t.Run("GetOption returns None for invalid JSON", func(t *testing.T) {
		invalidJSON := `{invalid json}`
		result := stringConfigPrism.GetOption(invalidJSON)

		assert.True(t, O.IsNone(result))
	})

	t.Run("ReverseGet creates JSON string from Config", func(t *testing.T) {
		cfg := Config{Host: "example.com", Port: 443}
		result := stringConfigPrism.ReverseGet(cfg)

		// Parse it back to verify
		var parsed Config
		err := json.Unmarshal([]byte(result), &parsed)
		assert.NoError(t, err)
		assert.Equal(t, cfg, parsed)
	})
}

// TestComposeWithEmptyValues tests edge cases with empty/zero values
func TestComposeWithEmptyValues(t *testing.T) {
	// Prism that extracts Right values
	prism := P.FromEither[error, []byte]()

	// Isomorphism between string and Either[error, []byte]
	iso := I.MakeIso(
		func(s string) E.Either[error, []byte] {
			return E.Right[error]([]byte(s))
		},
		func(e E.Either[error, []byte]) string {
			return string(E.GetOrElse(func(error) []byte { return []byte{} })(e))
		},
	)

	composed := Compose[string](prism)(iso)

	t.Run("Empty string is handled correctly", func(t *testing.T) {
		result := composed.GetOption("")

		assert.True(t, O.IsSome(result))
		bytes := O.GetOrElse(F.Constant([]byte("default")))(result)
		assert.Equal(t, []byte{}, bytes)
	})

	t.Run("ReverseGet with empty bytes", func(t *testing.T) {
		bytes := []byte{}
		result := composed.ReverseGet(bytes)

		assert.Equal(t, "", result)
	})
}
