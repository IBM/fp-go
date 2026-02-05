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
	"strconv"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEitherWithIdentityCodecs tests the Either function with identity codecs
// where both branches have the same output and input types
func TestEitherWithIdentityCodecs(t *testing.T) {
	t.Run("creates codec with correct name", func(t *testing.T) {
		// The Either function is designed for cases where both branches encode to the same type
		// For example, both encode to string or both encode to JSON

		// Create codecs that both encode to string
		stringToString := Id[string]()
		intToString := IntFromString()

		eitherCodec := Either(stringToString, intToString)

		assert.Equal(t, "Either[string, IntFromString]", eitherCodec.Name())
	})
}

// TestEitherEncode tests encoding of Either values
func TestEitherEncode(t *testing.T) {
	// Create codecs that both encode to string
	stringToString := Id[string]()
	intToString := IntFromString()

	eitherCodec := Either(stringToString, intToString)

	t.Run("encodes Left value", func(t *testing.T) {
		leftValue := either.Left[int]("hello")
		encoded := eitherCodec.Encode(leftValue)

		assert.Equal(t, "hello", encoded)
	})

	t.Run("encodes Right value", func(t *testing.T) {
		rightValue := either.Right[string](42)
		encoded := eitherCodec.Encode(rightValue)

		assert.Equal(t, "42", encoded)
	})
}

// TestEitherDecode tests decoding/validation of Either values
func TestEitherDecode(t *testing.T) {
	getOrElseNull := either.GetOrElse(reader.Of[validation.Errors, either.Either[string, int]](either.Left[int]("")))

	// Create codecs that both work with string input
	stringCodec := Id[string]()
	intFromString := IntFromString()

	eitherCodec := Either(stringCodec, intFromString)

	t.Run("decodes integer string as Right", func(t *testing.T) {
		result := eitherCodec.Decode("42")

		assert.True(t, either.IsRight(result), "should successfully decode integer string")

		value := getOrElseNull(result)
		assert.True(t, either.IsRight(value), "should be Right")

		rightValue := either.MonadFold(value,
			func(string) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, rightValue)
	})

	t.Run("decodes non-integer string as Left", func(t *testing.T) {
		result := eitherCodec.Decode("hello")

		assert.True(t, either.IsRight(result), "should successfully decode string")

		value := getOrElseNull(result)
		assert.True(t, either.IsLeft(value), "should be Left")

		leftValue := either.MonadFold(value,
			F.Identity[string],
			func(int) string { return "" },
		)
		assert.Equal(t, "hello", leftValue)
	})
}

// TestEitherValidation tests validation behavior
func TestEitherValidation(t *testing.T) {
	t.Run("validates with custom codecs", func(t *testing.T) {
		// Create a codec that only accepts non-empty strings
		nonEmptyString := MakeType(
			"NonEmptyString",
			func(u any) either.Either[error, string] {
				s, ok := u.(string)
				if !ok || len(s) == 0 {
					return either.Left[string](fmt.Errorf("not a non-empty string"))
				}
				return either.Of[error](s)
			},
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					if len(s) == 0 {
						return validation.FailureWithMessage[string](s, "must not be empty")(c)
					}
					return validation.Success(s)
				}
			},
			F.Identity[string],
		)

		// Create a codec that only accepts positive integers from strings
		positiveIntFromString := MakeType(
			"PositiveInt",
			func(u any) either.Either[error, int] {
				i, ok := u.(int)
				if !ok || i <= 0 {
					return either.Left[int](fmt.Errorf("not a positive integer"))
				}
				return either.Of[error](i)
			},
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					var n int
					_, err := fmt.Sscanf(s, "%d", &n)
					if err != nil {
						return validation.FailureWithError[int](s, "expected integer string")(err)(c)
					}
					if n <= 0 {
						return validation.FailureWithMessage[int](n, "must be positive")(c)
					}
					return validation.Success(n)
				}
			},
			func(n int) string {
				return fmt.Sprintf("%d", n)
			},
		)

		eitherCodec := Either(nonEmptyString, positiveIntFromString)

		// Valid non-empty string
		validLeft := eitherCodec.Decode("hello")
		assert.True(t, either.IsRight(validLeft))

		// Valid positive integer
		validRight := eitherCodec.Decode("42")
		assert.True(t, either.IsRight(validRight))

		// Invalid empty string - should fail both validations
		invalidEmpty := eitherCodec.Decode("")
		assert.True(t, either.IsLeft(invalidEmpty))

		// Invalid zero - should fail Right validation, succeed as Left
		zeroResult := eitherCodec.Decode("0")
		// "0" is a valid non-empty string, so it should succeed as Left
		assert.True(t, either.IsRight(zeroResult))
	})
}

// TestEitherRoundTrip tests encoding and decoding round trips
func TestEitherRoundTrip(t *testing.T) {
	stringCodec := Id[string]()
	intFromString := IntFromString()

	eitherCodec := Either(stringCodec, intFromString)

	t.Run("round-trip Left value", func(t *testing.T) {
		original := "hello"

		// Decode
		decodeResult := eitherCodec.Decode(original)
		require.True(t, either.IsRight(decodeResult))

		decoded := either.MonadFold(decodeResult,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)

		// Encode
		encoded := eitherCodec.Encode(decoded)

		// Verify
		assert.Equal(t, original, encoded)
	})

	t.Run("round-trip Right value", func(t *testing.T) {
		original := "42"

		// Decode
		decodeResult := eitherCodec.Decode(original)
		require.True(t, either.IsRight(decodeResult))

		decoded := either.MonadFold(decodeResult,
			func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)

		// Encode
		encoded := eitherCodec.Encode(decoded)

		// Verify
		assert.Equal(t, original, encoded)
	})
}

// TestEitherPrioritization tests that Right validation is prioritized over Left
func TestEitherPrioritization(t *testing.T) {
	stringCodec := Id[string]()
	intFromString := IntFromString()

	eitherCodec := Either(stringCodec, intFromString)

	t.Run("prioritizes Right over Left when both could succeed", func(t *testing.T) {
		// "42" can be validated as both string (Left) and int (Right)
		// The codec should prioritize Right
		result := eitherCodec.Decode("42")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)

		// Should be Right because int validation succeeds and is prioritized
		assert.True(t, either.IsRight(value))

		rightValue := either.MonadFold(value,
			func(string) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, rightValue)
	})

	t.Run("falls back to Left when Right fails", func(t *testing.T) {
		// "hello" can only be validated as string (Left), not as int (Right)
		result := eitherCodec.Decode("hello")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)

		// Should be Left because int validation failed
		assert.True(t, either.IsLeft(value))

		leftValue := either.MonadFold(value,
			F.Identity[string],
			func(int) string { return "" },
		)
		assert.Equal(t, "hello", leftValue)
	})
}

// TestEitherErrorAccumulation tests that errors from both branches are accumulated
func TestEitherErrorAccumulation(t *testing.T) {
	// Create codecs with specific validation rules that will both fail
	nonEmptyString := MakeType(
		"NonEmptyString",
		func(u any) either.Either[error, string] {
			s, ok := u.(string)
			if !ok || len(s) == 0 {
				return either.Left[string](fmt.Errorf("not a non-empty string"))
			}
			return either.Of[error](s)
		},
		func(s string) Decode[Context, string] {
			return func(c Context) Validation[string] {
				if len(s) == 0 {
					return validation.FailureWithMessage[string](s, "must not be empty")(c)
				}
				return validation.Success(s)
			}
		},
		F.Identity[string],
	)

	positiveIntFromString := MakeType(
		"PositiveInt",
		func(u any) either.Either[error, int] {
			i, ok := u.(int)
			if !ok || i <= 0 {
				return either.Left[int](fmt.Errorf("not a positive integer"))
			}
			return either.Of[error](i)
		},
		func(s string) Decode[Context, int] {
			return func(c Context) Validation[int] {
				var n int
				_, err := fmt.Sscanf(s, "%d", &n)
				if err != nil {
					return validation.FailureWithError[int](s, "expected integer string")(err)(c)
				}
				if n <= 0 {
					return validation.FailureWithMessage[int](n, "must be positive")(c)
				}
				return validation.Success(n)
			}
		},
		strconv.Itoa,
	)

	eitherCodec := Either(nonEmptyString, positiveIntFromString)

	t.Run("accumulates errors from both branches when both fail", func(t *testing.T) {
		// Empty string will fail both validations
		result := eitherCodec.Decode("")

		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(either.Either[string, int]) validation.Errors { return nil },
		)

		require.NotNil(t, errors)
		// Should have errors from both string and int validation attempts
		assert.GreaterOrEqual(t, len(errors), 2, "Should have at least 2 errors (one from Right validation, one from Left validation)")

		// Verify we have errors from both validation attempts
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}

		// Check that we have errors related to both validations
		hasIntError := false
		hasStringError := false
		for _, msg := range messages {
			if msg == "expected integer string" || msg == "must be positive" {
				hasIntError = true
			}
			if msg == "must not be empty" {
				hasStringError = true
			}
		}

		assert.True(t, hasIntError, "Should have error from integer validation (Right branch)")
		assert.True(t, hasStringError, "Should have error from string validation (Left branch)")
	})
}
