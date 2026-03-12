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
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/iso"
	"github.com/stretchr/testify/assert"
)

// Test types for newtype pattern
type UserId int
type Email string
type Celsius float64
type Fahrenheit float64

func TestFromIso_Success(t *testing.T) {
	t.Run("decodes using iso.Get", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		result := codec.Decode(UserId(42))

		// Assert
		assert.Equal(t, validation.Success(42), result)
	})

	t.Run("encodes using iso.ReverseGet", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		encoded := codec.Encode(42)

		// Assert
		assert.Equal(t, UserId(42), encoded)
	})

	t.Run("round-trip preserves value", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)
		original := UserId(123)

		// Act
		decoded := codec.Decode(original)

		// Assert
		assert.True(t, either.IsRight(decoded))
		roundTrip := either.Fold[validation.Errors, int, UserId](
			func(validation.Errors) UserId { return UserId(0) },
			codec.Encode,
		)(decoded)
		assert.Equal(t, original, roundTrip)
	})
}

func TestFromIso_StringTypes(t *testing.T) {
	t.Run("handles string newtype", func(t *testing.T) {
		// Arrange
		emailIso := iso.MakeIso(
			func(e Email) string { return string(e) },
			func(s string) Email { return Email(s) },
		)
		codec := FromIso[string, Email](emailIso)

		// Act
		result := codec.Decode(Email("user@example.com"))

		// Assert
		assert.Equal(t, validation.Success("user@example.com"), result)
	})

	t.Run("encodes string newtype", func(t *testing.T) {
		// Arrange
		emailIso := iso.MakeIso(
			func(e Email) string { return string(e) },
			func(s string) Email { return Email(s) },
		)
		codec := FromIso[string, Email](emailIso)

		// Act
		encoded := codec.Encode("admin@example.com")

		// Assert
		assert.Equal(t, Email("admin@example.com"), encoded)
	})

	t.Run("handles empty string", func(t *testing.T) {
		// Arrange
		emailIso := iso.MakeIso(
			func(e Email) string { return string(e) },
			func(s string) Email { return Email(s) },
		)
		codec := FromIso[string, Email](emailIso)

		// Act
		result := codec.Decode(Email(""))

		// Assert
		assert.Equal(t, validation.Success(""), result)
	})
}

func TestFromIso_NumericConversions(t *testing.T) {
	t.Run("converts Celsius to Fahrenheit", func(t *testing.T) {
		// Arrange
		tempIso := iso.MakeIso(
			func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
			func(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) },
		)
		codec := FromIso[Fahrenheit, Celsius](tempIso)

		// Act
		result := codec.Decode(Celsius(0))

		// Assert
		assert.Equal(t, validation.Success(Fahrenheit(32)), result)
	})

	t.Run("converts Fahrenheit to Celsius", func(t *testing.T) {
		// Arrange
		tempIso := iso.MakeIso(
			func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
			func(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) },
		)
		codec := FromIso[Fahrenheit, Celsius](tempIso)

		// Act
		encoded := codec.Encode(Fahrenheit(68))

		// Assert
		assert.Equal(t, Celsius(20), encoded)
	})

	t.Run("handles negative temperatures", func(t *testing.T) {
		// Arrange
		tempIso := iso.MakeIso(
			func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
			func(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) },
		)
		codec := FromIso[Fahrenheit, Celsius](tempIso)

		// Act
		result := codec.Decode(Celsius(-40))

		// Assert
		assert.Equal(t, validation.Success(Fahrenheit(-40)), result)
	})

	t.Run("temperature round-trip", func(t *testing.T) {
		// Arrange
		tempIso := iso.MakeIso(
			func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
			func(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) },
		)
		codec := FromIso[Fahrenheit, Celsius](tempIso)
		original := Celsius(25)

		// Act
		decoded := codec.Decode(original)

		// Assert
		assert.True(t, either.IsRight(decoded))
		roundTrip := either.Fold[validation.Errors, Fahrenheit, Celsius](
			func(validation.Errors) Celsius { return Celsius(0) },
			codec.Encode,
		)(decoded)
		// Allow small floating point error
		diff := float64(original - roundTrip)
		if diff < 0 {
			diff = -diff
		}
		assert.True(t, diff < 0.0001)
	})
}

func TestFromIso_EdgeCases(t *testing.T) {
	t.Run("handles zero values", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		result := codec.Decode(UserId(0))

		// Assert
		assert.Equal(t, validation.Success(0), result)
	})

	t.Run("handles negative values", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		result := codec.Decode(UserId(-1))

		// Assert
		assert.Equal(t, validation.Success(-1), result)
	})

	t.Run("handles large values", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		result := codec.Decode(UserId(999999999))

		// Assert
		assert.Equal(t, validation.Success(999999999), result)
	})
}

func TestFromIso_TypeChecking(t *testing.T) {
	t.Run("Is checks target type", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		isResult := codec.Is(42)

		// Assert
		assert.True(t, either.IsRight(isResult))
	})

	t.Run("Is rejects wrong type", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		isResult := codec.Is("not an int")

		// Assert
		assert.True(t, either.IsLeft(isResult))
	})
}

func TestFromIso_Name(t *testing.T) {
	t.Run("includes iso in name", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		name := codec.Name()

		// Assert
		assert.True(t, len(name) > 0)
		assert.True(t, name[:7] == "FromIso")
	})
}

func TestFromIso_Composition(t *testing.T) {
	t.Run("composes with Pipe", func(t *testing.T) {
		// Arrange
		type PositiveInt int

		// First iso: UserId -> int
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)

		// Second iso: int -> PositiveInt (no validation, just type conversion)
		positiveIso := iso.MakeIso(
			func(i int) PositiveInt { return PositiveInt(i) },
			func(p PositiveInt) int { return int(p) },
		)

		// Compose codecs
		codec := F.Pipe1(
			FromIso[int, UserId](userIdIso),
			Pipe[UserId, UserId](FromIso[PositiveInt, int](positiveIso)),
		)

		// Act
		result := codec.Decode(UserId(42))

		// Assert
		assert.Equal(t, validation.Of(PositiveInt(42)), result)
	})

	t.Run("composed codec encodes correctly", func(t *testing.T) {
		// Arrange
		type PositiveInt int

		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)

		positiveIso := iso.MakeIso(
			func(i int) PositiveInt { return PositiveInt(i) },
			func(p PositiveInt) int { return int(p) },
		)

		codec := F.Pipe1(
			FromIso[int, UserId](userIdIso),
			Pipe[UserId, UserId](FromIso[PositiveInt, int](positiveIso)),
		)

		// Act
		encoded := codec.Encode(PositiveInt(42))

		// Assert
		assert.Equal(t, UserId(42), encoded)
	})
}

func TestFromIso_Integration(t *testing.T) {
	t.Run("works with Array codec", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		userIdCodec := FromIso[int, UserId](userIdIso)
		arrayCodec := TranscodeArray(userIdCodec)

		// Act
		result := arrayCodec.Decode([]UserId{UserId(1), UserId(2), UserId(3)})

		// Assert
		assert.Equal(t, validation.Success([]int{1, 2, 3}), result)
	})

	t.Run("encodes array correctly", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		userIdCodec := FromIso[int, UserId](userIdIso)
		arrayCodec := TranscodeArray(userIdCodec)

		// Act
		encoded := arrayCodec.Encode([]int{1, 2, 3})

		// Assert
		assert.Equal(t, []UserId{UserId(1), UserId(2), UserId(3)}, encoded)
	})

	t.Run("handles empty array", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		userIdCodec := FromIso[int, UserId](userIdIso)
		arrayCodec := TranscodeArray(userIdCodec)

		// Act
		result := arrayCodec.Decode([]UserId{})

		// Assert
		assert.True(t, either.IsRight(result))
		decoded := either.Fold[validation.Errors, []int, []int](
			func(validation.Errors) []int { return nil },
			func(arr []int) []int { return arr },
		)(result)
		assert.Equal(t, 0, len(decoded))
	})
}

func TestFromIso_ComplexTypes(t *testing.T) {
	t.Run("handles struct wrapping", func(t *testing.T) {
		// Arrange
		type Wrapper struct{ Value int }

		wrapperIso := iso.MakeIso(
			func(w Wrapper) int { return w.Value },
			func(i int) Wrapper { return Wrapper{Value: i} },
		)
		codec := FromIso[int, Wrapper](wrapperIso)

		// Act
		result := codec.Decode(Wrapper{Value: 42})

		// Assert
		assert.Equal(t, validation.Success(42), result)
	})

	t.Run("encodes struct wrapping", func(t *testing.T) {
		// Arrange
		type Wrapper struct{ Value int }

		wrapperIso := iso.MakeIso(
			func(w Wrapper) int { return w.Value },
			func(i int) Wrapper { return Wrapper{Value: i} },
		)
		codec := FromIso[int, Wrapper](wrapperIso)

		// Act
		encoded := codec.Encode(42)

		// Assert
		assert.Equal(t, Wrapper{Value: 42}, encoded)
	})
}

func TestFromIso_AsDecoder(t *testing.T) {
	t.Run("returns decoder interface", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		decoder := codec.AsDecoder()

		// Assert
		result := decoder.Decode(UserId(42))
		assert.Equal(t, validation.Success(42), result)
	})
}

func TestFromIso_AsEncoder(t *testing.T) {
	t.Run("returns encoder interface", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		encoder := codec.AsEncoder()

		// Assert
		encoded := encoder.Encode(42)
		assert.Equal(t, UserId(42), encoded)
	})
}

func TestFromIso_Validate(t *testing.T) {
	t.Run("validate method works correctly", func(t *testing.T) {
		// Arrange
		userIdIso := iso.MakeIso(
			func(id UserId) int { return int(id) },
			func(i int) UserId { return UserId(i) },
		)
		codec := FromIso[int, UserId](userIdIso)

		// Act
		validateFn := codec.Validate(UserId(42))
		result := validateFn([]validation.ContextEntry{})

		// Assert
		assert.Equal(t, validation.Success(42), result)
	})
}

// Made with Bob
