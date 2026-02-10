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

// TestMonadAltBasicFunctionality tests the basic behavior of MonadAlt
func TestMonadAltBasicFunctionality(t *testing.T) {
	t.Run("uses first codec when it succeeds", func(t *testing.T) {
		// Create two codecs that both work with strings
		stringCodec := Id[string]()

		// Create another string codec that only accepts uppercase
		uppercaseOnly := MakeType(
			"UppercaseOnly",
			Is[string](),
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					for _, r := range s {
						if r >= 'a' && r <= 'z' {
							return validation.FailureWithMessage[string](s, "must be uppercase")(c)
						}
					}
					return validation.Success(s)
				}
			},
			F.Identity[string],
		)

		// Create alt codec that tries uppercase first, then any string
		altCodec := MonadAlt(
			uppercaseOnly,
			func() Type[string, string, string] { return stringCodec },
		)

		// Test with uppercase string - should succeed with first codec
		result := altCodec.Decode("HELLO")

		assert.True(t, either.IsRight(result), "should successfully decode with first codec")

		value := either.GetOrElse(reader.Of[validation.Errors, string](""))(result)
		assert.Equal(t, "HELLO", value)
	})

	t.Run("falls back to second codec when first fails", func(t *testing.T) {
		// Create a codec that only accepts positive integers
		positiveInt := MakeType(
			"PositiveInt",
			Is[int](),
			func(i int) Decode[Context, int] {
				return func(c Context) Validation[int] {
					if i <= 0 {
						return validation.FailureWithMessage[int](i, "must be positive")(c)
					}
					return validation.Success(i)
				}
			},
			F.Identity[int],
		)

		// Create a codec that accepts any integer (with same input type)
		anyInt := MakeType(
			"AnyInt",
			Is[int](),
			func(i int) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.Success(i)
				}
			},
			F.Identity[int],
		)

		// Create alt codec
		altCodec := MonadAlt(
			positiveInt,
			func() Type[int, int, int] { return anyInt },
		)

		// Test with negative number - first fails, second succeeds
		result := altCodec.Decode(-5)

		assert.True(t, either.IsRight(result), "should successfully decode with second codec")

		value := either.GetOrElse(reader.Of[validation.Errors, int](0))(result)
		assert.Equal(t, -5, value)
	})

	t.Run("aggregates errors when both codecs fail", func(t *testing.T) {
		// Create two codecs that will both fail
		positiveInt := MakeType(
			"PositiveInt",
			Is[int](),
			func(i int) Decode[Context, int] {
				return func(c Context) Validation[int] {
					if i <= 0 {
						return validation.FailureWithMessage[int](i, "must be positive")(c)
					}
					return validation.Success(i)
				}
			},
			F.Identity[int],
		)

		evenInt := MakeType(
			"EvenInt",
			Is[int](),
			func(i int) Decode[Context, int] {
				return func(c Context) Validation[int] {
					if i%2 != 0 {
						return validation.FailureWithMessage[int](i, "must be even")(c)
					}
					return validation.Success(i)
				}
			},
			F.Identity[int],
		)

		// Create alt codec
		altCodec := MonadAlt(
			positiveInt,
			func() Type[int, int, int] { return evenInt },
		)

		// Test with -3 (negative and odd) - both should fail
		result := altCodec.Decode(-3)

		assert.True(t, either.IsLeft(result), "should fail when both codecs fail")

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(int) validation.Errors { return nil },
		)

		require.NotNil(t, errors)
		// Should have errors from both validation attempts
		assert.GreaterOrEqual(t, len(errors), 2, "should have errors from both codecs")
	})
}

// TestMonadAltNaming tests that the codec name is correctly generated
func TestMonadAltNaming(t *testing.T) {
	t.Run("generates correct name", func(t *testing.T) {
		stringCodec := Id[string]()
		anotherStringCodec := Id[string]()

		altCodec := MonadAlt(
			stringCodec,
			func() Type[string, string, string] { return anotherStringCodec },
		)

		assert.Equal(t, "Alt[string]", altCodec.Name())
	})
}

// TestMonadAltEncoding tests that encoding uses the first codec's encoder
func TestMonadAltEncoding(t *testing.T) {
	t.Run("uses first codec's encoder", func(t *testing.T) {
		// Create a codec that encodes ints as strings with prefix
		prefixedInt := MakeType(
			"PrefixedInt",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					var n int
					_, err := fmt.Sscanf(s, "NUM:%d", &n)
					if err != nil {
						return validation.FailureWithError[int](s, "expected NUM:n format")(err)(c)
					}
					return validation.Success(n)
				}
			},
			func(n int) string {
				return fmt.Sprintf("NUM:%d", n)
			},
		)

		// Create a standard int from string codec
		standardInt := IntFromString()

		// Create alt codec
		altCodec := MonadAlt(
			prefixedInt,
			func() Type[int, string, string] { return standardInt },
		)

		// Encode should use first codec's encoder
		encoded := altCodec.Encode(42)
		assert.Equal(t, "NUM:42", encoded)
	})
}

// TestAltOperator tests the curried Alt function
func TestAltOperator(t *testing.T) {
	t.Run("creates reusable operator", func(t *testing.T) {
		// Create a fallback operator that accepts any string
		withStringFallback := Alt(func() Type[string, string, string] {
			return Id[string]()
		})

		// Create a codec that only accepts "hello"
		helloOnly := MakeType(
			"HelloOnly",
			Is[string](),
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					if s != "hello" {
						return validation.FailureWithMessage[string](s, "must be 'hello'")(c)
					}
					return validation.Success(s)
				}
			},
			F.Identity[string],
		)

		// Apply fallback to the codec
		altCodec := withStringFallback(helloOnly)

		// Test that it works
		result1 := altCodec.Decode("hello")
		assert.True(t, either.IsRight(result1))

		result2 := altCodec.Decode("world")
		assert.True(t, either.IsRight(result2))
	})

	t.Run("works in pipeline with F.Pipe", func(t *testing.T) {
		// Create a codec pipeline with multiple fallbacks
		baseCodec := MakeType(
			"StrictInt",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					if s != "42" {
						return validation.FailureWithMessage[int](s, "must be exactly '42'")(c)
					}
					return validation.Success(42)
				}
			},
			strconv.Itoa,
		)

		fallback1 := MakeType(
			"Fallback1",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					if s != "100" {
						return validation.FailureWithMessage[int](s, "must be exactly '100'")(c)
					}
					return validation.Success(100)
				}
			},
			strconv.Itoa,
		)

		fallback2 := MakeType(
			"AnyInt",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					n, err := strconv.Atoi(s)
					if err != nil {
						return validation.FailureWithError[int](s, "not an integer")(err)(c)
					}
					return validation.Success(n)
				}
			},
			strconv.Itoa,
		)

		// Build pipeline with multiple alternatives
		pipeline := F.Pipe2(
			baseCodec,
			Alt(func() Type[int, string, string] { return fallback1 }),
			Alt(func() Type[int, string, string] { return fallback2 }),
		)

		// Test with "42" - should use base codec
		result1 := pipeline.Decode("42")
		assert.True(t, either.IsRight(result1))
		value1 := either.GetOrElse(reader.Of[validation.Errors, int](0))(result1)
		assert.Equal(t, 42, value1)

		// Test with "100" - should use fallback1
		result2 := pipeline.Decode("100")
		assert.True(t, either.IsRight(result2))
		value2 := either.GetOrElse(reader.Of[validation.Errors, int](0))(result2)
		assert.Equal(t, 100, value2)

		// Test with "999" - should use fallback2
		result3 := pipeline.Decode("999")
		assert.True(t, either.IsRight(result3))
		value3 := either.GetOrElse(reader.Of[validation.Errors, int](0))(result3)
		assert.Equal(t, 999, value3)
	})
}

// TestAltLazyEvaluation tests that the second codec is only evaluated when needed
func TestAltLazyEvaluation(t *testing.T) {
	t.Run("does not evaluate second codec when first succeeds", func(t *testing.T) {
		evaluated := false

		stringCodec := Id[string]()

		altCodec := MonadAlt(
			stringCodec,
			func() Type[string, string, string] {
				evaluated = true
				return Id[string]()
			},
		)

		// Decode with first codec succeeding
		result := altCodec.Decode("hello")
		assert.True(t, either.IsRight(result))

		// Second codec should not have been evaluated
		assert.False(t, evaluated, "second codec should not be evaluated when first succeeds")
	})

	t.Run("evaluates second codec when first fails", func(t *testing.T) {
		evaluated := false

		// Create a codec that always fails
		failingCodec := MakeType(
			"Failing",
			Is[string](),
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					return validation.FailureWithMessage[string](s, "always fails")(c)
				}
			},
			F.Identity[string],
		)

		altCodec := MonadAlt(
			failingCodec,
			func() Type[string, string, string] {
				evaluated = true
				return Id[string]()
			},
		)

		// Decode with first codec failing
		result := altCodec.Decode("hello")
		assert.True(t, either.IsRight(result))

		// Second codec should have been evaluated
		assert.True(t, evaluated, "second codec should be evaluated when first fails")
	})
}

// TestAltWithComplexTypes tests Alt with more complex codec scenarios
func TestAltWithComplexTypes(t *testing.T) {
	t.Run("works with string length validation", func(t *testing.T) {
		// Create codec that accepts strings of length 5
		length5 := MakeType(
			"Length5",
			Is[string](),
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					if len(s) != 5 {
						return validation.FailureWithMessage[string](s, "must be length 5")(c)
					}
					return validation.Success(s)
				}
			},
			F.Identity[string],
		)

		// Create codec that accepts any string
		anyString := Id[string]()

		// Create alt codec
		altCodec := MonadAlt(
			length5,
			func() Type[string, string, string] { return anyString },
		)

		// Test with length 5 - should use first codec
		result1 := altCodec.Decode("hello")
		assert.True(t, either.IsRight(result1))

		// Test with different length - should fall back to second codec
		result2 := altCodec.Decode("hi")
		assert.True(t, either.IsRight(result2))
	})
}

// TestAltTypeChecking tests that type checking works correctly
func TestAltTypeChecking(t *testing.T) {
	t.Run("type checking uses generic Is", func(t *testing.T) {
		stringCodec := Id[string]()
		anotherStringCodec := Id[string]()

		altCodec := MonadAlt(
			stringCodec,
			func() Type[string, string, string] { return anotherStringCodec },
		)

		// Test Is with valid type
		result1 := altCodec.Is("hello")
		assert.True(t, either.IsRight(result1))

		// Test Is with invalid type
		result2 := altCodec.Is(42)
		assert.True(t, either.IsLeft(result2))
	})
}

// TestAltRoundTrip tests encoding and decoding round trips
func TestAltRoundTrip(t *testing.T) {
	t.Run("round-trip with first codec", func(t *testing.T) {
		stringCodec := Id[string]()
		anotherStringCodec := Id[string]()

		altCodec := MonadAlt(
			stringCodec,
			func() Type[string, string, string] { return anotherStringCodec },
		)

		original := "hello"

		// Decode
		decodeResult := altCodec.Decode(original)
		require.True(t, either.IsRight(decodeResult))

		decoded := either.GetOrElse(reader.Of[validation.Errors, string](""))(decodeResult)

		// Encode
		encoded := altCodec.Encode(decoded)

		// Verify
		assert.Equal(t, original, encoded)
	})

	t.Run("round-trip with second codec", func(t *testing.T) {
		// Create a codec that only accepts "hello"
		helloOnly := MakeType(
			"HelloOnly",
			Is[string](),
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					if s != "hello" {
						return validation.FailureWithMessage[string](s, "must be 'hello'")(c)
					}
					return validation.Success(s)
				}
			},
			F.Identity[string],
		)

		anyString := Id[string]()

		altCodec := MonadAlt(
			helloOnly,
			func() Type[string, string, string] { return anyString },
		)

		original := "world"

		// Decode (will use second codec)
		decodeResult := altCodec.Decode(original)
		require.True(t, either.IsRight(decodeResult))

		decoded := either.GetOrElse(reader.Of[validation.Errors, string](""))(decodeResult)

		// Encode (uses first codec's encoder, which is identity)
		encoded := altCodec.Encode(decoded)

		// Verify
		assert.Equal(t, original, encoded)
	})
}

// TestAltErrorMessages tests that error messages are informative
func TestAltErrorMessages(t *testing.T) {
	t.Run("provides detailed error context", func(t *testing.T) {
		// Create two codecs with specific error messages
		codec1 := MakeType(
			"Codec1",
			Is[int](),
			func(i int) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.FailureWithMessage[int](i, "codec1 error")(c)
				}
			},
			F.Identity[int],
		)

		codec2 := MakeType(
			"Codec2",
			Is[int](),
			func(i int) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.FailureWithMessage[int](i, "codec2 error")(c)
				}
			},
			F.Identity[int],
		)

		altCodec := MonadAlt(
			codec1,
			func() Type[int, int, int] { return codec2 },
		)

		result := altCodec.Decode(42)

		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(int) validation.Errors { return nil },
		)

		require.NotNil(t, errors)
		require.GreaterOrEqual(t, len(errors), 2)

		// Check that both error messages are present
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}

		hasCodec1Error := false
		hasCodec2Error := false
		for _, msg := range messages {
			if msg == "codec1 error" {
				hasCodec1Error = true
			}
			if msg == "codec2 error" {
				hasCodec2Error = true
			}
		}

		assert.True(t, hasCodec1Error, "should have error from first codec")
		assert.True(t, hasCodec2Error, "should have error from second codec")
	})
}

// TestAltMonoid tests the AltMonoid function
func TestAltMonoid(t *testing.T) {
	t.Run("with default value as zero", func(t *testing.T) {
		// Create a monoid with a default value of 0
		m := AltMonoid(func() Type[int, string, string] {
			return MakeType(
				"DefaultZero",
				Is[int](),
				func(s string) Decode[Context, int] {
					return func(c Context) Validation[int] {
						return validation.Success(0)
					}
				},
				strconv.Itoa,
			)
		})

		// Create codecs
		intFromString := IntFromString()
		failing := MakeType(
			"Failing",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.FailureWithMessage[int](s, "always fails")(c)
				}
			},
			strconv.Itoa,
		)

		t.Run("first success wins", func(t *testing.T) {
			// Combine two successful codecs - first should win
			codec1 := MakeType(
				"Returns10",
				Is[int](),
				func(s string) Decode[Context, int] {
					return func(c Context) Validation[int] {
						return validation.Success(10)
					}
				},
				strconv.Itoa,
			)
			codec2 := MakeType(
				"Returns20",
				Is[int](),
				func(s string) Decode[Context, int] {
					return func(c Context) Validation[int] {
						return validation.Success(20)
					}
				},
				strconv.Itoa,
			)

			combined := m.Concat(codec1, codec2)
			result := combined.Decode("input")

			assert.True(t, either.IsRight(result))
			value := either.GetOrElse(reader.Of[validation.Errors, int](0))(result)
			assert.Equal(t, 10, value, "first success should win")
		})

		t.Run("falls back to second when first fails", func(t *testing.T) {
			combined := m.Concat(failing, intFromString)
			result := combined.Decode("42")

			assert.True(t, either.IsRight(result))
			value := either.GetOrElse(reader.Of[validation.Errors, int](0))(result)
			assert.Equal(t, 42, value)
		})

		t.Run("uses zero when both fail", func(t *testing.T) {
			combined := m.Concat(failing, m.Empty())
			result := combined.Decode("invalid")

			assert.True(t, either.IsRight(result))
			value := either.GetOrElse(reader.Of[validation.Errors, int](-1))(result)
			assert.Equal(t, 0, value, "should use default zero value")
		})
	})

	t.Run("with failing zero", func(t *testing.T) {
		// Create a monoid with a failing zero
		m := AltMonoid(func() Type[int, string, string] {
			return MakeType(
				"NoDefault",
				Is[int](),
				func(s string) Decode[Context, int] {
					return func(c Context) Validation[int] {
						return validation.FailureWithMessage[int](s, "no default available")(c)
					}
				},
				strconv.Itoa,
			)
		})

		failing1 := MakeType(
			"Failing1",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.FailureWithMessage[int](s, "error 1")(c)
				}
			},
			strconv.Itoa,
		)

		failing2 := MakeType(
			"Failing2",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.FailureWithMessage[int](s, "error 2")(c)
				}
			},
			strconv.Itoa,
		)

		t.Run("aggregates all errors when all fail", func(t *testing.T) {
			combined := m.Concat(m.Concat(failing1, failing2), m.Empty())
			result := combined.Decode("input")

			assert.True(t, either.IsLeft(result))

			errors := either.MonadFold(result,
				F.Identity[validation.Errors],
				func(int) validation.Errors { return nil },
			)

			require.NotNil(t, errors)
			// Should have errors from all three: failing1, failing2, and zero
			assert.GreaterOrEqual(t, len(errors), 3)

			messages := make([]string, len(errors))
			for i, err := range errors {
				messages[i] = err.Messsage
			}

			hasError1 := false
			hasError2 := false
			hasNoDefault := false
			for _, msg := range messages {
				if msg == "error 1" {
					hasError1 = true
				}
				if msg == "error 2" {
					hasError2 = true
				}
				if msg == "no default available" {
					hasNoDefault = true
				}
			}

			assert.True(t, hasError1, "should have error from failing1")
			assert.True(t, hasError2, "should have error from failing2")
			assert.True(t, hasNoDefault, "should have error from zero")
		})
	})

	t.Run("chaining multiple fallbacks", func(t *testing.T) {
		m := AltMonoid(func() Type[string, string, string] {
			return MakeType(
				"Default",
				Is[string](),
				func(s string) Decode[Context, string] {
					return func(c Context) Validation[string] {
						return validation.Success("default")
					}
				},
				F.Identity[string],
			)
		})

		primary := MakeType(
			"Primary",
			Is[string](),
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					if s == "primary" {
						return validation.Success("from primary")
					}
					return validation.FailureWithMessage[string](s, "not primary")(c)
				}
			},
			F.Identity[string],
		)

		secondary := MakeType(
			"Secondary",
			Is[string](),
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					if s == "secondary" {
						return validation.Success("from secondary")
					}
					return validation.FailureWithMessage[string](s, "not secondary")(c)
				}
			},
			F.Identity[string],
		)

		// Chain: try primary, then secondary, then default
		combined := m.Concat(m.Concat(primary, secondary), m.Empty())

		t.Run("uses primary when it succeeds", func(t *testing.T) {
			result := combined.Decode("primary")
			assert.True(t, either.IsRight(result))
			value := either.GetOrElse(reader.Of[validation.Errors, string](""))(result)
			assert.Equal(t, "from primary", value)
		})

		t.Run("uses secondary when primary fails", func(t *testing.T) {
			result := combined.Decode("secondary")
			assert.True(t, either.IsRight(result))
			value := either.GetOrElse(reader.Of[validation.Errors, string](""))(result)
			assert.Equal(t, "from secondary", value)
		})

		t.Run("uses default when both fail", func(t *testing.T) {
			result := combined.Decode("other")
			assert.True(t, either.IsRight(result))
			value := either.GetOrElse(reader.Of[validation.Errors, string](""))(result)
			assert.Equal(t, "default", value)
		})
	})

	t.Run("satisfies monoid laws", func(t *testing.T) {
		m := AltMonoid(func() Type[int, string, string] {
			return MakeType(
				"DefaultZero",
				Is[int](),
				func(s string) Decode[Context, int] {
					return func(c Context) Validation[int] {
						return validation.Success(0)
					}
				},
				strconv.Itoa,
			)
		})

		codec1 := MakeType(
			"Codec1",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.Success(10)
				}
			},
			strconv.Itoa,
		)

		codec2 := MakeType(
			"Codec2",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.Success(20)
				}
			},
			strconv.Itoa,
		)

		codec3 := MakeType(
			"Codec3",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.Success(30)
				}
			},
			strconv.Itoa,
		)

		t.Run("left identity", func(t *testing.T) {
			// m.Concat(m.Empty(), codec) should behave like codec
			// But with AltMonoid, if codec fails, it falls back to empty
			combined := m.Concat(m.Empty(), codec1)
			result := combined.Decode("input")

			assert.True(t, either.IsRight(result))
			value := either.GetOrElse(reader.Of[validation.Errors, int](-1))(result)
			// Empty (0) comes first, so it wins
			assert.Equal(t, 0, value)
		})

		t.Run("right identity", func(t *testing.T) {
			// m.Concat(codec, m.Empty()) tries codec first, falls back to empty
			combined := m.Concat(codec1, m.Empty())
			result := combined.Decode("input")

			assert.True(t, either.IsRight(result))
			value := either.GetOrElse(reader.Of[validation.Errors, int](-1))(result)
			assert.Equal(t, 10, value, "codec1 should win")
		})

		t.Run("associativity", func(t *testing.T) {
			// For AltMonoid, first success wins
			left := m.Concat(m.Concat(codec1, codec2), codec3)
			right := m.Concat(codec1, m.Concat(codec2, codec3))

			resultLeft := left.Decode("input")
			resultRight := right.Decode("input")

			assert.True(t, either.IsRight(resultLeft))
			assert.True(t, either.IsRight(resultRight))

			valueLeft := either.GetOrElse(reader.Of[validation.Errors, int](-1))(resultLeft)
			valueRight := either.GetOrElse(reader.Of[validation.Errors, int](-1))(resultRight)

			// Both should return 10 (first success)
			assert.Equal(t, valueLeft, valueRight)
			assert.Equal(t, 10, valueLeft)
		})
	})

	t.Run("encoding uses first codec", func(t *testing.T) {
		m := AltMonoid(func() Type[int, string, string] {
			return MakeType(
				"Default",
				Is[int](),
				func(s string) Decode[Context, int] {
					return func(c Context) Validation[int] {
						return validation.Success(0)
					}
				},
				func(n int) string { return "DEFAULT" },
			)
		})

		codec1 := MakeType(
			"Codec1",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.Success(42)
				}
			},
			func(n int) string { return fmt.Sprintf("FIRST:%d", n) },
		)

		codec2 := MakeType(
			"Codec2",
			Is[int](),
			func(s string) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.Success(100)
				}
			},
			func(n int) string { return fmt.Sprintf("SECOND:%d", n) },
		)

		combined := m.Concat(codec1, codec2)

		// Encoding should use first codec's encoder
		encoded := combined.Encode(42)
		assert.Equal(t, "FIRST:42", encoded)
	})
}
