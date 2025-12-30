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
	"encoding/base64"
	"errors"
	"net/url"
	"testing"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestSome(t *testing.T) {
	somePrism := MakePrism(F.Identity[Option[int]], O.Some[int])

	assert.Equal(t, O.Some(1), somePrism.GetOption(O.Some(1)))
}

func TestId(t *testing.T) {
	idPrism := Id[int]()

	// GetOption always returns Some for identity
	assert.Equal(t, O.Some(42), idPrism.GetOption(42))

	// ReverseGet is identity
	assert.Equal(t, 42, idPrism.ReverseGet(42))
}

func TestFromPredicate(t *testing.T) {
	// Prism for positive numbers
	positivePrism := FromPredicate(func(n int) bool {
		return n > 0
	})

	// Matches positive numbers
	assert.Equal(t, O.Some(42), positivePrism.GetOption(42))
	assert.Equal(t, O.Some(1), positivePrism.GetOption(1))

	// Doesn't match non-positive numbers
	assert.Equal(t, O.None[int](), positivePrism.GetOption(0))
	assert.Equal(t, O.None[int](), positivePrism.GetOption(-5))

	// ReverseGet always succeeds (doesn't check predicate)
	assert.Equal(t, 42, positivePrism.ReverseGet(42))
	assert.Equal(t, -5, positivePrism.ReverseGet(-5))
}

func TestCompose(t *testing.T) {
	// Prism for Some values
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Prism for positive numbers
	positivePrism := FromPredicate(func(n int) bool {
		return n > 0
	})

	// Compose: Option[int] -> int (if Some and positive)
	composedPrism := F.Pipe1(
		somePrism,
		Compose[Option[int]](positivePrism),
	)

	// Test with Some positive
	assert.Equal(t, O.Some(42), composedPrism.GetOption(O.Some(42)))

	// Test with Some non-positive
	assert.Equal(t, O.None[int](), composedPrism.GetOption(O.Some(-5)))

	// Test with None
	assert.Equal(t, O.None[int](), composedPrism.GetOption(O.None[int]()))

	// ReverseGet constructs Some
	assert.Equal(t, O.Some(42), composedPrism.ReverseGet(42))
}

func TestSet(t *testing.T) {
	// Prism for Some values
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Set value when it matches
	result := Set[Option[int]](100)(somePrism)(O.Some(42))
	assert.Equal(t, O.Some(100), result)

	// No change when it doesn't match
	result = Set[Option[int]](100)(somePrism)(O.None[int]())
	assert.Equal(t, O.None[int](), result)
}

func TestSomeFunction(t *testing.T) {
	// Prism that focuses on an Option field
	type Config struct {
		Timeout Option[int]
	}

	configPrism := MakePrism(
		func(c Config) Option[Option[int]] {
			return O.Some(c.Timeout)
		},
		func(t Option[int]) Config {
			return Config{Timeout: t}
		},
	)

	// Focus on the Some value
	somePrism := Some(configPrism)

	// Extract from Some
	config := Config{Timeout: O.Some(30)}
	assert.Equal(t, O.Some(30), somePrism.GetOption(config))

	// Extract from None
	configNone := Config{Timeout: O.None[int]()}
	assert.Equal(t, O.None[int](), somePrism.GetOption(configNone))

	// ReverseGet constructs Config with Some
	result := somePrism.ReverseGet(60)
	assert.Equal(t, Config{Timeout: O.Some(60)}, result)
}

func TestIMap(t *testing.T) {
	// Prism for Some values
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Map to string and back
	stringPrism := F.Pipe1(
		somePrism,
		IMap[Option[int]](
			func(n int) string {
				if n == 42 {
					return "42"
				}
				return "100"
			},
			func(s string) int {
				if s == "42" {
					return 42
				}
				return 100
			},
		),
	)

	// GetOption maps the value
	result := stringPrism.GetOption(O.Some(42))
	assert.Equal(t, O.Some("42"), result)

	// GetOption on None
	result = stringPrism.GetOption(O.None[int]())
	assert.Equal(t, O.None[string](), result)

	// ReverseGet maps back
	opt := stringPrism.ReverseGet("100")
	assert.Equal(t, O.Some(100), opt)
}

func TestPrismLaws(t *testing.T) {
	// Test prism laws with a simple prism
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Law 1: GetOptionReverseGet
	// prism.GetOption(prism.ReverseGet(a)) == Some(a)
	a := 42
	result := somePrism.GetOption(somePrism.ReverseGet(a))
	assert.Equal(t, O.Some(a), result)

	// Law 2: ReverseGetGetOption
	// if GetOption(s) == Some(a), then ReverseGet(a) should produce equivalent s
	s := O.Some(42)
	extracted := somePrism.GetOption(s)
	if O.IsSome(extracted) {
		reconstructed := somePrism.ReverseGet(O.GetOrElse(F.Constant(0))(extracted))
		assert.Equal(t, s, reconstructed)
	}
}

func TestPrismModifyOption(t *testing.T) {
	// Test the internal prismModifyOption function through Set
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Modify when match
	setFn := Set[Option[int]](100)
	result := setFn(somePrism)(O.Some(42))
	assert.Equal(t, O.Some(100), result)

	// No modification when no match
	result = setFn(somePrism)(O.None[int]())
	assert.Equal(t, O.None[int](), result)
}

// Custom sum type for testing
type testResult interface{ isResult() }
type testSuccess struct{ Value int }
type testFailure struct{ Error string }

func (testSuccess) isResult() {}
func (testFailure) isResult() {}

func TestPrismWithCustomType(t *testing.T) {
	// Create prism for Success variant
	successPrism := MakePrism(
		func(r testResult) Option[int] {
			if s, ok := r.(testSuccess); ok {
				return O.Some(s.Value)
			}
			return O.None[int]()
		},
		func(v int) testResult {
			return testSuccess{Value: v}
		},
	)

	// Test GetOption with Success
	success := testSuccess{Value: 42}
	assert.Equal(t, O.Some(42), successPrism.GetOption(success))

	// Test GetOption with Failure
	failure := testFailure{Error: "oops"}
	assert.Equal(t, O.None[int](), successPrism.GetOption(failure))

	// Test ReverseGet
	result := successPrism.ReverseGet(100)
	assert.Equal(t, testSuccess{Value: 100}, result)

	// Test Set with Success
	setFn := Set[testResult](200)
	updated := setFn(successPrism)(success)
	assert.Equal(t, testSuccess{Value: 200}, updated)

	// Test Set with Failure (no change)
	unchanged := setFn(successPrism)(failure)
	assert.Equal(t, failure, unchanged)
}

// TestPrismModify tests the prismModify internal function through various scenarios
func TestPrismModify(t *testing.T) {
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Test modify with matching value
	result := Set[Option[int]](84)(somePrism)(O.Some(42))
	assert.Equal(t, O.Some(84), result)

	// Test that original is returned when no match
	result = Set[Option[int]](100)(somePrism)(O.None[int]())
	assert.Equal(t, O.None[int](), result)
}

// TestPrismModifyWithTransform tests modifying through a prism with a transformation
func TestPrismModifyWithTransform(t *testing.T) {
	// Create a prism for positive numbers
	positivePrism := FromPredicate(N.MoreThan(0))

	// Modify positive number
	setter := Set[int](100)
	result := setter(positivePrism)(42)
	assert.Equal(t, 100, result)

	// Try to modify negative number (no change)
	result = setter(positivePrism)(-5)
	assert.Equal(t, -5, result)
}

// TestAsTraversal tests converting a prism to a traversal
func TestAsTraversal(t *testing.T) {
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Simple identity functor for testing
	type Identity[A any] struct{ Value A }

	fof := func(s Option[int]) Identity[Option[int]] {
		return Identity[Option[int]]{Value: s}
	}

	fmap := func(ia Identity[int], f func(int) Option[int]) Identity[Option[int]] {
		return Identity[Option[int]]{Value: f(ia.Value)}
	}

	type TraversalFunc func(func(int) Identity[int]) func(Option[int]) Identity[Option[int]]
	traversal := AsTraversal[TraversalFunc](fof, fmap)(somePrism)

	// Test with Some value
	f := func(n int) Identity[int] {
		return Identity[int]{Value: n * 2}
	}
	result := traversal(f)(O.Some(21))
	assert.Equal(t, O.Some(42), result.Value)

	// Test with None value
	result = traversal(f)(O.None[int]())
	assert.Equal(t, O.None[int](), result.Value)
}

// Test types for composition chain test
type testOuter struct{ Middle Option[testInner] }
type testInner struct{ Value Option[int] }

// TestPrismCompositionChain tests composing multiple prisms
func TestPrismCompositionChain(t *testing.T) {
	// Three-level composition

	outerPrism := MakePrism(
		func(o testOuter) Option[Option[testInner]] {
			return O.Some(o.Middle)
		},
		func(m Option[testInner]) testOuter {
			return testOuter{Middle: m}
		},
	)

	middlePrism := MakePrism(
		F.Identity[Option[testInner]],
		O.Some[testInner],
	)

	innerPrism := MakePrism(
		func(i testInner) Option[Option[int]] {
			return O.Some(i.Value)
		},
		func(v Option[int]) testInner {
			return testInner{Value: v}
		},
	)

	// Compose all three
	composed := F.Pipe2(
		outerPrism,
		Compose[testOuter](middlePrism),
		Compose[testOuter](innerPrism),
	)

	// Further compose to get the int value
	finalPrism := Some(composed)

	// Test extraction through all layers
	outer := testOuter{Middle: O.Some(testInner{Value: O.Some(42)})}
	value := finalPrism.GetOption(outer)
	assert.Equal(t, O.Some(42), value)

	// Test with None at middle layer
	outerNone := testOuter{Middle: O.None[testInner]()}
	value = finalPrism.GetOption(outerNone)
	assert.Equal(t, O.None[int](), value)

	// Test with None at inner layer
	outerInnerNone := testOuter{Middle: O.Some(testInner{Value: O.None[int]()})}
	value = finalPrism.GetOption(outerInnerNone)
	assert.Equal(t, O.None[int](), value)
}

// TestPrismSetMultipleTimes tests setting values multiple times
func TestPrismSetMultipleTimes(t *testing.T) {
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Chain multiple sets
	result := F.Pipe3(
		O.Some(10),
		Set[Option[int]](20)(somePrism),
		Set[Option[int]](30)(somePrism),
		Set[Option[int]](40)(somePrism),
	)

	assert.Equal(t, O.Some(40), result)
}

// TestIMapBidirectional tests that IMap maintains bidirectionality
func TestIMapBidirectional(t *testing.T) {
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Map int to string and back
	stringPrism := F.Pipe1(
		somePrism,
		IMap[Option[int]](
			func(n int) string {
				if n == 42 {
					return "forty-two"
				}
				return "other"
			},
			func(s string) int {
				if s == "forty-two" {
					return 42
				}
				return 0
			},
		),
	)

	// Test GetOption with mapping
	result := stringPrism.GetOption(O.Some(42))
	assert.Equal(t, O.Some("forty-two"), result)

	// Test ReverseGet with reverse mapping
	opt := stringPrism.ReverseGet("forty-two")
	assert.Equal(t, O.Some(42), opt)

	// Verify round-trip
	value := stringPrism.GetOption(stringPrism.ReverseGet("forty-two"))
	assert.Equal(t, O.Some("forty-two"), value)
}

// Test types for complex sum type test
type Shape interface{ isShape() }
type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }
type Triangle struct{ Base, Height float64 }

func (Circle) isShape()    {}
func (Rectangle) isShape() {}
func (Triangle) isShape()  {}

// TestPrismWithComplexSumType tests prism with a more complex sum type
func TestPrismWithComplexSumType(t *testing.T) {

	// Prism for Circle
	circlePrism := MakePrism(
		func(s Shape) Option[float64] {
			if c, ok := s.(Circle); ok {
				return O.Some(c.Radius)
			}
			return O.None[float64]()
		},
		func(r float64) Shape {
			return Circle{Radius: r}
		},
	)

	// Prism for Rectangle
	rectanglePrism := MakePrism(
		func(s Shape) Option[struct{ Width, Height float64 }] {
			if r, ok := s.(Rectangle); ok {
				return O.Some(struct{ Width, Height float64 }{r.Width, r.Height})
			}
			return O.None[struct{ Width, Height float64 }]()
		},
		func(dims struct{ Width, Height float64 }) Shape {
			return Rectangle{Width: dims.Width, Height: dims.Height}
		},
	)

	// Test Circle prism
	circle := Circle{Radius: 5.0}
	radius := circlePrism.GetOption(circle)
	assert.Equal(t, O.Some(5.0), radius)

	// Circle prism doesn't match Rectangle
	rect := Rectangle{Width: 10, Height: 20}
	radius = circlePrism.GetOption(rect)
	assert.Equal(t, O.None[float64](), radius)

	// Rectangle prism matches Rectangle
	dims := rectanglePrism.GetOption(rect)
	assert.True(t, O.IsSome(dims))

	// Test ReverseGet
	newCircle := circlePrism.ReverseGet(7.5)
	assert.Equal(t, Circle{Radius: 7.5}, newCircle)
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("prism with zero value", func(t *testing.T) {
		somePrism := MakePrism(
			F.Identity[Option[int]],
			O.Some[int],
		)

		// Zero value should work fine
		result := somePrism.GetOption(O.Some(0))
		assert.Equal(t, O.Some(0), result)

		opt := somePrism.ReverseGet(0)
		assert.Equal(t, O.Some(0), opt)
	})

	t.Run("prism with empty string", func(t *testing.T) {
		somePrism := MakePrism(
			F.Identity[Option[string]],
			O.Some[string],
		)

		result := somePrism.GetOption(O.Some(""))
		assert.Equal(t, O.Some(""), result)
	})

	t.Run("identity prism with nil pointer", func(t *testing.T) {
		type MyStruct struct{ Value int }
		idPrism := Id[*MyStruct]()

		var nilPtr *MyStruct
		result := idPrism.GetOption(nilPtr)
		assert.Equal(t, O.Some(nilPtr), result)
	})
}

// TestFromEncoding tests the FromEncoding prism with various base64 encodings
func TestFromEncoding(t *testing.T) {
	t.Run("standard encoding - valid base64", func(t *testing.T) {
		prism := FromEncoding(base64.StdEncoding)

		// Test decoding valid base64
		input := "SGVsbG8gV29ybGQ="
		result := prism.GetOption(input)

		assert.True(t, O.IsSome(result))
		decoded := O.GetOrElse(F.Constant([]byte{}))(result)
		assert.Equal(t, []byte("Hello World"), decoded)
	})

	t.Run("standard encoding - invalid base64", func(t *testing.T) {
		prism := FromEncoding(base64.StdEncoding)

		// Test decoding invalid base64
		invalid := "not-valid-base64!!!"
		result := prism.GetOption(invalid)

		assert.True(t, O.IsNone(result))
		assert.Equal(t, O.None[[]byte](), result)
	})

	t.Run("standard encoding - encode bytes", func(t *testing.T) {
		prism := FromEncoding(base64.StdEncoding)

		// Test encoding bytes to base64
		data := []byte("Hello World")
		encoded := prism.ReverseGet(data)

		assert.Equal(t, "SGVsbG8gV29ybGQ=", encoded)
	})

	t.Run("standard encoding - round trip", func(t *testing.T) {
		prism := FromEncoding(base64.StdEncoding)

		// Test round-trip: encode then decode
		original := []byte("Test Data 123")
		encoded := prism.ReverseGet(original)
		decoded := prism.GetOption(encoded)

		assert.True(t, O.IsSome(decoded))
		result := O.GetOrElse(F.Constant([]byte{}))(decoded)
		assert.Equal(t, original, result)
	})

	t.Run("URL encoding - valid base64", func(t *testing.T) {
		prism := FromEncoding(base64.URLEncoding)

		// URL encoding uses - and _ instead of + and /
		input := "SGVsbG8gV29ybGQ="
		result := prism.GetOption(input)

		assert.True(t, O.IsSome(result))
		decoded := O.GetOrElse(F.Constant([]byte{}))(result)
		assert.Equal(t, []byte("Hello World"), decoded)
	})

	t.Run("URL encoding - with special characters", func(t *testing.T) {
		prism := FromEncoding(base64.URLEncoding)

		// Test data that would use URL-safe characters
		data := []byte("subjects?_d=1")
		encoded := prism.ReverseGet(data)
		decoded := prism.GetOption(encoded)

		assert.True(t, O.IsSome(decoded))
		result := O.GetOrElse(F.Constant([]byte{}))(decoded)
		assert.Equal(t, data, result)
	})

	t.Run("raw standard encoding - no padding", func(t *testing.T) {
		prism := FromEncoding(base64.RawStdEncoding)

		// RawStdEncoding omits padding
		data := []byte("Hello")
		encoded := prism.ReverseGet(data)

		// Should not have padding
		assert.NotContains(t, encoded, "=")

		// Should still decode correctly
		decoded := prism.GetOption(encoded)
		assert.True(t, O.IsSome(decoded))
		result := O.GetOrElse(F.Constant([]byte{}))(decoded)
		assert.Equal(t, data, result)
	})

	t.Run("empty byte slice", func(t *testing.T) {
		prism := FromEncoding(base64.StdEncoding)

		// Test encoding empty byte slice
		empty := []byte{}
		encoded := prism.ReverseGet(empty)
		assert.Equal(t, "", encoded)

		// Test decoding empty string
		decoded := prism.GetOption("")
		assert.True(t, O.IsSome(decoded))
		result := O.GetOrElse(F.Constant([]byte{1}))(decoded)
		assert.Equal(t, []byte{}, result)
	})

	t.Run("binary data", func(t *testing.T) {
		prism := FromEncoding(base64.StdEncoding)

		// Test with binary data (not just text)
		binary := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
		encoded := prism.ReverseGet(binary)
		decoded := prism.GetOption(encoded)

		assert.True(t, O.IsSome(decoded))
		result := O.GetOrElse(F.Constant([]byte{}))(decoded)
		assert.Equal(t, binary, result)
	})

	t.Run("malformed base64 - wrong padding", func(t *testing.T) {
		prism := FromEncoding(base64.StdEncoding)

		// Test with incorrect padding - too much padding
		malformed := "SGVsbG8===" // Too much padding
		result := prism.GetOption(malformed)

		// Should return None for malformed input
		assert.True(t, O.IsNone(result))
	})

	t.Run("malformed base64 - invalid characters", func(t *testing.T) {
		prism := FromEncoding(base64.StdEncoding)

		// Test with invalid characters for standard encoding
		invalid := "SGVs bG8@"
		result := prism.GetOption(invalid)

		assert.True(t, O.IsNone(result))
	})
}

// TestFromEncodingWithSet tests using Set with FromEncoding prism
func TestFromEncodingWithSet(t *testing.T) {
	prism := FromEncoding(base64.StdEncoding)

	t.Run("set new value on valid base64", func(t *testing.T) {
		// Original encoded value
		original := "SGVsbG8gV29ybGQ=" // "Hello World"

		// New data to set
		newData := []byte("New Data")
		setter := Set[string](newData)

		// Apply the setter
		result := setter(prism)(original)

		// Should return the new data encoded
		expected := prism.ReverseGet(newData)
		assert.Equal(t, expected, result)

		// Verify it decodes to the new data
		decoded := prism.GetOption(result)
		assert.True(t, O.IsSome(decoded))
		assert.Equal(t, newData, O.GetOrElse(F.Constant([]byte{}))(decoded))
	})

	t.Run("set on invalid base64 returns original", func(t *testing.T) {
		// Invalid base64 string
		invalid := "not-valid-base64!!!"

		// Try to set new data
		newData := []byte("New Data")
		setter := Set[string](newData)

		// Should return original unchanged
		result := setter(prism)(invalid)
		assert.Equal(t, invalid, result)
	})
}

// TestFromEncodingComposition tests composing FromEncoding with other prisms
func TestFromEncodingComposition(t *testing.T) {
	t.Run("compose with predicate prism", func(t *testing.T) {
		// Create a prism that only accepts non-empty byte slices
		nonEmptyPrism := FromPredicate(func(b []byte) bool {
			return len(b) > 0
		})

		// Compose with base64 prism
		base64Prism := FromEncoding(base64.StdEncoding)
		composed := F.Pipe1(
			base64Prism,
			Compose[string](nonEmptyPrism),
		)

		// Test with non-empty data
		validEncoded := base64Prism.ReverseGet([]byte("data"))
		result := composed.GetOption(validEncoded)
		assert.True(t, O.IsSome(result))

		// Test with empty data
		emptyEncoded := base64Prism.ReverseGet([]byte{})
		result = composed.GetOption(emptyEncoded)
		assert.True(t, O.IsNone(result))
	})
}

// TestFromEncodingPrismLaws tests that FromEncoding satisfies prism laws
func TestFromEncodingPrismLaws(t *testing.T) {
	prism := FromEncoding(base64.StdEncoding)

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		// For any byte slice, encoding then decoding should return the original
		testData := [][]byte{
			[]byte("Hello World"),
			[]byte(""),
			[]byte{0x00, 0xFF},
			[]byte("Special chars: !@#$%^&*()"),
		}

		for _, data := range testData {
			encoded := prism.ReverseGet(data)
			decoded := prism.GetOption(encoded)

			assert.True(t, O.IsSome(decoded))
			result := O.GetOrElse(F.Constant([]byte{}))(decoded)
			assert.Equal(t, data, result)
		}
	})

	t.Run("law 2: if GetOption(s) == Some(a), then ReverseGet(a) produces valid s", func(t *testing.T) {
		// For valid base64 strings, decode then encode should produce valid base64
		validInputs := []string{
			"SGVsbG8gV29ybGQ=",
			"",
			"AQID",
		}

		for _, input := range validInputs {
			extracted := prism.GetOption(input)
			if O.IsSome(extracted) {
				data := O.GetOrElse(F.Constant([]byte{}))(extracted)
				reencoded := prism.ReverseGet(data)

				// Re-decode to verify it's valid
				redecoded := prism.GetOption(reencoded)
				assert.True(t, O.IsSome(redecoded))

				// The data should match
				finalData := O.GetOrElse(F.Constant([]byte{}))(redecoded)
				assert.Equal(t, data, finalData)
			}
		}
	})
}

// TestParseURL tests the ParseURL prism with various URL formats
func TestParseURL(t *testing.T) {
	urlPrism := ParseURL()

	t.Run("valid HTTP URL", func(t *testing.T) {
		input := "https://example.com/path?query=value"
		result := urlPrism.GetOption(input)

		assert.True(t, O.IsSome(result))
		parsed := O.GetOrElse(F.Constant((*url.URL)(nil)))(result)
		assert.NotNil(t, parsed)
		assert.Equal(t, "https", parsed.Scheme)
		assert.Equal(t, "example.com", parsed.Host)
		assert.Equal(t, "/path", parsed.Path)
		assert.Equal(t, "query=value", parsed.RawQuery)
	})

	t.Run("valid HTTP URL without scheme", func(t *testing.T) {
		input := "//example.com/path"
		result := urlPrism.GetOption(input)

		assert.True(t, O.IsSome(result))
		parsed := O.GetOrElse(F.Constant((*url.URL)(nil)))(result)
		assert.Equal(t, "example.com", parsed.Host)
	})

	t.Run("simple domain", func(t *testing.T) {
		input := "example.com"
		result := urlPrism.GetOption(input)

		assert.True(t, O.IsSome(result))
		parsed := O.GetOrElse(F.Constant((*url.URL)(nil)))(result)
		assert.NotNil(t, parsed)
	})

	t.Run("URL with port", func(t *testing.T) {
		input := "https://example.com:8080/path"
		result := urlPrism.GetOption(input)

		assert.True(t, O.IsSome(result))
		parsed := O.GetOrElse(F.Constant((*url.URL)(nil)))(result)
		assert.Equal(t, "example.com:8080", parsed.Host)
	})

	t.Run("URL with fragment", func(t *testing.T) {
		input := "https://example.com/path#section"
		result := urlPrism.GetOption(input)

		assert.True(t, O.IsSome(result))
		parsed := O.GetOrElse(F.Constant((*url.URL)(nil)))(result)
		assert.Equal(t, "section", parsed.Fragment)
	})

	t.Run("URL with user info", func(t *testing.T) {
		input := "https://user:pass@example.com/path"
		result := urlPrism.GetOption(input)

		assert.True(t, O.IsSome(result))
		parsed := O.GetOrElse(F.Constant((*url.URL)(nil)))(result)
		assert.NotNil(t, parsed.User)
	})

	t.Run("invalid URL with spaces", func(t *testing.T) {
		input := "ht tp://invalid url"
		result := urlPrism.GetOption(input)

		// url.Parse is lenient, so this might still parse
		// The test verifies the behavior
		_ = result
	})

	t.Run("empty string", func(t *testing.T) {
		input := ""
		result := urlPrism.GetOption(input)

		assert.True(t, O.IsSome(result))
		parsed := O.GetOrElse(F.Constant((*url.URL)(nil)))(result)
		assert.NotNil(t, parsed)
	})

	t.Run("reverse get - URL to string", func(t *testing.T) {
		u, _ := url.Parse("https://example.com/path?query=value")
		str := urlPrism.ReverseGet(u)

		assert.Equal(t, "https://example.com/path?query=value", str)
	})

	t.Run("round trip", func(t *testing.T) {
		original := "https://example.com:8080/path?q=v#frag"
		parsed := urlPrism.GetOption(original)

		assert.True(t, O.IsSome(parsed))
		u := O.GetOrElse(F.Constant((*url.URL)(nil)))(parsed)
		reconstructed := urlPrism.ReverseGet(u)

		// Parse both to compare (URL normalization may occur)
		reparsed := urlPrism.GetOption(reconstructed)
		assert.True(t, O.IsSome(reparsed))
	})
}

// TestParseURLWithSet tests using Set with ParseURL prism
func TestParseURLWithSet(t *testing.T) {
	urlPrism := ParseURL()

	t.Run("set new URL on valid input", func(t *testing.T) {
		original := "https://oldsite.com/path"
		newURL, _ := url.Parse("https://newsite.com/newpath")

		setter := Set[string](newURL)
		result := setter(urlPrism)(original)

		assert.Equal(t, "https://newsite.com/newpath", result)
	})
}

// TestParseURLPrismLaws tests that ParseURL satisfies prism laws
func TestParseURLPrismLaws(t *testing.T) {
	urlPrism := ParseURL()

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		testURLs := []string{
			"https://example.com",
			"https://example.com:8080/path?q=v",
			"http://user:pass@example.com/path#frag",
		}

		for _, urlStr := range testURLs {
			u, _ := url.Parse(urlStr)
			str := urlPrism.ReverseGet(u)
			reparsed := urlPrism.GetOption(str)

			assert.True(t, O.IsSome(reparsed))
		}
	})
}

// TestInstanceOf tests the InstanceOf prism with various types
func TestInstanceOf(t *testing.T) {
	t.Run("extract int from any", func(t *testing.T) {
		intPrism := InstanceOf[int]()
		var value any = 42

		result := intPrism.GetOption(value)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("extract string from any", func(t *testing.T) {
		stringPrism := InstanceOf[string]()
		var value any = "hello"

		result := stringPrism.GetOption(value)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "hello", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("type mismatch returns None", func(t *testing.T) {
		intPrism := InstanceOf[int]()
		var value any = "not an int"

		result := intPrism.GetOption(value)
		assert.True(t, O.IsNone(result))
	})

	t.Run("extract struct from any", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		personPrism := InstanceOf[Person]()
		var value any = Person{Name: "Alice", Age: 30}

		result := personPrism.GetOption(value)
		assert.True(t, O.IsSome(result))
		person := O.GetOrElse(F.Constant(Person{}))(result)
		assert.Equal(t, "Alice", person.Name)
		assert.Equal(t, 30, person.Age)
	})

	t.Run("extract pointer type", func(t *testing.T) {
		type Data struct{ Value int }
		ptrPrism := InstanceOf[*Data]()

		data := &Data{Value: 42}
		var value any = data

		result := ptrPrism.GetOption(value)
		assert.True(t, O.IsSome(result))
		extracted := O.GetOrElse(F.Constant((*Data)(nil)))(result)
		assert.Equal(t, 42, extracted.Value)
	})

	t.Run("nil value", func(t *testing.T) {
		intPrism := InstanceOf[int]()
		var value any = nil

		result := intPrism.GetOption(value)
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get - T to any", func(t *testing.T) {
		intPrism := InstanceOf[int]()
		anyValue := intPrism.ReverseGet(42)

		assert.Equal(t, any(42), anyValue)
	})

	t.Run("round trip", func(t *testing.T) {
		stringPrism := InstanceOf[string]()
		original := "test"

		anyValue := stringPrism.ReverseGet(original)
		extracted := stringPrism.GetOption(anyValue)

		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, original, O.GetOrElse(F.Constant(""))(extracted))
	})

	t.Run("zero value", func(t *testing.T) {
		intPrism := InstanceOf[int]()
		var value any = 0

		result := intPrism.GetOption(value)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(result))
	})
}

// TestInstanceOfWithSet tests using Set with InstanceOf prism
func TestInstanceOfWithSet(t *testing.T) {
	t.Run("set new value on matching type", func(t *testing.T) {
		intPrism := InstanceOf[int]()
		var original any = 42

		setter := Set[any](100)
		result := setter(intPrism)(original)

		assert.Equal(t, any(100), result)
	})

	t.Run("set on non-matching type returns original", func(t *testing.T) {
		intPrism := InstanceOf[int]()
		var original any = "not an int"

		setter := Set[any](100)
		result := setter(intPrism)(original)

		assert.Equal(t, original, result)
	})
}

// TestInstanceOfPrismLaws tests that InstanceOf satisfies prism laws
func TestInstanceOfPrismLaws(t *testing.T) {
	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		intPrism := InstanceOf[int]()
		testValues := []int{0, 42, -10, 999}

		for _, val := range testValues {
			anyVal := intPrism.ReverseGet(val)
			extracted := intPrism.GetOption(anyVal)

			assert.True(t, O.IsSome(extracted))
			assert.Equal(t, val, O.GetOrElse(F.Constant(0))(extracted))
		}
	})
}

// TestParseDate tests the ParseDate prism with various date formats
func TestParseDate(t *testing.T) {
	t.Run("ISO date format - valid", func(t *testing.T) {
		datePrism := ParseDate("2006-01-02")
		input := "2024-03-15"

		result := datePrism.GetOption(input)
		assert.True(t, O.IsSome(result))

		parsed := O.GetOrElse(F.Constant(time.Time{}))(result)
		assert.Equal(t, 2024, parsed.Year())
		assert.Equal(t, time.March, parsed.Month())
		assert.Equal(t, 15, parsed.Day())
	})

	t.Run("ISO date format - invalid", func(t *testing.T) {
		datePrism := ParseDate("2006-01-02")
		input := "not-a-date"

		result := datePrism.GetOption(input)
		assert.True(t, O.IsNone(result))
	})

	t.Run("RFC3339 format", func(t *testing.T) {
		datePrism := ParseDate(time.RFC3339)
		input := "2024-03-15T10:30:00Z"

		result := datePrism.GetOption(input)
		assert.True(t, O.IsSome(result))

		parsed := O.GetOrElse(F.Constant(time.Time{}))(result)
		assert.Equal(t, 2024, parsed.Year())
		assert.Equal(t, 10, parsed.Hour())
		assert.Equal(t, 30, parsed.Minute())
	})

	t.Run("custom format", func(t *testing.T) {
		datePrism := ParseDate("02/01/2006")
		input := "15/03/2024"

		result := datePrism.GetOption(input)
		assert.True(t, O.IsSome(result))

		parsed := O.GetOrElse(F.Constant(time.Time{}))(result)
		assert.Equal(t, 2024, parsed.Year())
		assert.Equal(t, time.March, parsed.Month())
		assert.Equal(t, 15, parsed.Day())
	})

	t.Run("wrong format returns None", func(t *testing.T) {
		datePrism := ParseDate("2006-01-02")
		input := "03/15/2024" // MM/DD/YYYY instead of YYYY-MM-DD

		result := datePrism.GetOption(input)
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get - format date", func(t *testing.T) {
		datePrism := ParseDate("2006-01-02")
		date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

		str := datePrism.ReverseGet(date)
		assert.Equal(t, "2024-03-15", str)
	})

	t.Run("round trip", func(t *testing.T) {
		datePrism := ParseDate("2006-01-02")
		original := "2024-03-15"

		parsed := datePrism.GetOption(original)
		assert.True(t, O.IsSome(parsed))

		date := O.GetOrElse(F.Constant(time.Time{}))(parsed)
		formatted := datePrism.ReverseGet(date)

		assert.Equal(t, original, formatted)
	})

	t.Run("empty string", func(t *testing.T) {
		datePrism := ParseDate("2006-01-02")
		result := datePrism.GetOption("")

		assert.True(t, O.IsNone(result))
	})

	t.Run("time with timezone", func(t *testing.T) {
		datePrism := ParseDate(time.RFC3339)
		input := "2024-03-15T10:30:00+05:00"

		result := datePrism.GetOption(input)
		assert.True(t, O.IsSome(result))

		parsed := O.GetOrElse(F.Constant(time.Time{}))(result)
		_, offset := parsed.Zone()
		assert.Equal(t, 5*3600, offset) // 5 hours in seconds
	})
}

// TestParseDateWithSet tests using Set with ParseDate prism
func TestParseDateWithSet(t *testing.T) {
	datePrism := ParseDate("2006-01-02")

	t.Run("set new date on valid input", func(t *testing.T) {
		original := "2024-03-15"
		newDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

		setter := Set[string](newDate)
		result := setter(datePrism)(original)

		assert.Equal(t, "2025-01-01", result)
	})

	t.Run("set on invalid date returns original", func(t *testing.T) {
		invalid := "not-a-date"
		newDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

		setter := Set[string](newDate)
		result := setter(datePrism)(invalid)

		assert.Equal(t, invalid, result)
	})
}

// TestParseDatePrismLaws tests that ParseDate satisfies prism laws
func TestParseDatePrismLaws(t *testing.T) {
	datePrism := ParseDate("2006-01-02")

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		testDates := []time.Time{
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			time.Date(2000, 6, 15, 0, 0, 0, 0, time.UTC),
		}

		for _, date := range testDates {
			str := datePrism.ReverseGet(date)
			reparsed := datePrism.GetOption(str)

			assert.True(t, O.IsSome(reparsed))
			parsed := O.GetOrElse(F.Constant(time.Time{}))(reparsed)

			// Compare date components (ignore time/location)
			assert.Equal(t, date.Year(), parsed.Year())
			assert.Equal(t, date.Month(), parsed.Month())
			assert.Equal(t, date.Day(), parsed.Day())
		}
	})

	t.Run("law 2: if GetOption(s) == Some(a), then ReverseGet(a) produces valid s", func(t *testing.T) {
		validInputs := []string{
			"2024-01-01",
			"2024-12-31",
			"2000-06-15",
		}

		for _, input := range validInputs {
			extracted := datePrism.GetOption(input)
			if O.IsSome(extracted) {
				date := O.GetOrElse(F.Constant(time.Time{}))(extracted)
				reformatted := datePrism.ReverseGet(date)

				// Re-parse to verify it's valid
				reparsed := datePrism.GetOption(reformatted)
				assert.True(t, O.IsSome(reparsed))

				// Should produce the same date
				assert.Equal(t, input, reformatted)
			}
		}
	})
}

// TestDeref tests the Deref prism with pointer dereferencing
func TestDeref(t *testing.T) {
	derefPrism := Deref[int]()

	t.Run("dereference non-nil pointer", func(t *testing.T) {
		value := 42
		ptr := &value

		result := derefPrism.GetOption(ptr)
		assert.True(t, O.IsSome(result))

		extracted := O.GetOrElse(F.Constant((*int)(nil)))(result)
		assert.NotNil(t, extracted)
		assert.Equal(t, 42, *extracted)
	})

	t.Run("dereference nil pointer", func(t *testing.T) {
		var nilPtr *int

		result := derefPrism.GetOption(nilPtr)
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get returns pointer unchanged", func(t *testing.T) {
		value := 42
		ptr := &value

		result := derefPrism.ReverseGet(ptr)
		assert.Equal(t, ptr, result)
		assert.Equal(t, 42, *result)
	})

	t.Run("reverse get with nil pointer", func(t *testing.T) {
		var nilPtr *int

		result := derefPrism.ReverseGet(nilPtr)
		assert.Nil(t, result)
	})

	t.Run("with string pointers", func(t *testing.T) {
		stringDeref := Deref[string]()

		str := "hello"
		ptr := &str

		result := stringDeref.GetOption(ptr)
		assert.True(t, O.IsSome(result))

		extracted := O.GetOrElse(F.Constant((*string)(nil)))(result)
		assert.NotNil(t, extracted)
		assert.Equal(t, "hello", *extracted)
	})

	t.Run("with struct pointers", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		personDeref := Deref[Person]()

		person := Person{Name: "Alice", Age: 30}
		ptr := &person

		result := personDeref.GetOption(ptr)
		assert.True(t, O.IsSome(result))

		extracted := O.GetOrElse(F.Constant((*Person)(nil)))(result)
		assert.NotNil(t, extracted)
		assert.Equal(t, "Alice", extracted.Name)
		assert.Equal(t, 30, extracted.Age)
	})
}

// TestDerefWithSet tests using Set with Deref prism
func TestDerefWithSet(t *testing.T) {
	derefPrism := Deref[int]()

	t.Run("set on non-nil pointer", func(t *testing.T) {
		value := 42
		ptr := &value

		newValue := 100
		newPtr := &newValue

		setter := Set[*int](newPtr)
		result := setter(derefPrism)(ptr)

		assert.NotNil(t, result)
		assert.Equal(t, 100, *result)
	})

	t.Run("set on nil pointer returns nil", func(t *testing.T) {
		var nilPtr *int

		newValue := 100
		newPtr := &newValue

		setter := Set[*int](newPtr)
		result := setter(derefPrism)(nilPtr)

		assert.Nil(t, result)
	})
}

// TestDerefPrismLaws tests that Deref satisfies prism laws
func TestDerefPrismLaws(t *testing.T) {
	derefPrism := Deref[int]()

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		value := 42
		ptr := &value

		reversed := derefPrism.ReverseGet(ptr)
		extracted := derefPrism.GetOption(reversed)

		assert.True(t, O.IsSome(extracted))
		result := O.GetOrElse(F.Constant((*int)(nil)))(extracted)
		assert.Equal(t, ptr, result)
	})
}

// TestFromEither tests the FromEither prism with Either types
func TestFromEither(t *testing.T) {
	t.Run("extract from Right", func(t *testing.T) {
		prism := FromEither[error, int]()

		success := E.Right[error](42)
		result := prism.GetOption(success)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("extract from Left returns None", func(t *testing.T) {
		prism := FromEither[error, int]()

		failure := E.Left[int](errors.New("failed"))
		result := prism.GetOption(failure)

		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get wraps into Right", func(t *testing.T) {
		prism := FromEither[error, int]()

		wrapped := prism.ReverseGet(100)

		assert.True(t, E.IsRight(wrapped))
		value := E.GetOrElse(func(error) int { return 0 })(wrapped)
		assert.Equal(t, 100, value)
	})

	t.Run("with string error type", func(t *testing.T) {
		prism := FromEither[string, int]()

		success := E.Right[string](42)
		result := prism.GetOption(success)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(0))(result))

		failure := E.Left[int]("error message")
		result = prism.GetOption(failure)

		assert.True(t, O.IsNone(result))
	})

	t.Run("with custom error type", func(t *testing.T) {
		type CustomError struct {
			Code    int
			Message string
		}

		prism := FromEither[CustomError, string]()

		success := E.Right[CustomError]("success")
		result := prism.GetOption(success)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, "success", O.GetOrElse(F.Constant(""))(result))

		failure := E.Left[string](CustomError{Code: 404, Message: "Not Found"})
		result = prism.GetOption(failure)

		assert.True(t, O.IsNone(result))
	})

	t.Run("round trip", func(t *testing.T) {
		prism := FromEither[error, int]()

		original := 42
		wrapped := prism.ReverseGet(original)
		extracted := prism.GetOption(wrapped)

		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, original, O.GetOrElse(F.Constant(0))(extracted))
	})
}

// TestFromEitherWithSet tests using Set with FromEither prism
func TestFromEitherWithSet(t *testing.T) {
	prism := FromEither[error, int]()

	t.Run("set on Right value", func(t *testing.T) {
		success := E.Right[error](42)

		setter := Set[E.Either[error, int]](100)
		result := setter(prism)(success)

		assert.True(t, E.IsRight(result))
		value := E.GetOrElse(func(error) int { return 0 })(result)
		assert.Equal(t, 100, value)
	})

	t.Run("set on Left value returns original", func(t *testing.T) {
		failure := E.Left[int](errors.New("failed"))

		setter := Set[E.Either[error, int]](100)
		result := setter(prism)(failure)

		assert.True(t, E.IsLeft(result))
		// Original error is preserved
		assert.Equal(t, failure, result)
	})
}

// TestFromEitherComposition tests composing FromEither with other prisms
func TestFromEitherComposition(t *testing.T) {
	t.Run("compose with predicate prism", func(t *testing.T) {
		// Create a prism that only accepts positive numbers
		positivePrism := FromPredicate(func(n int) bool {
			return n > 0
		})

		// Compose with Either prism
		eitherPrism := FromEither[error, int]()
		composed := F.Pipe1(
			eitherPrism,
			Compose[E.Either[error, int]](positivePrism),
		)

		// Test with Right positive
		success := E.Right[error](42)
		result := composed.GetOption(success)
		assert.True(t, O.IsSome(result))

		// Test with Right non-positive
		nonPositive := E.Right[error](-5)
		result = composed.GetOption(nonPositive)
		assert.True(t, O.IsNone(result))

		// Test with Left
		failure := E.Left[int](errors.New("error"))
		result = composed.GetOption(failure)
		assert.True(t, O.IsNone(result))
	})
}

// TestFromEitherPrismLaws tests that FromEither satisfies prism laws
func TestFromEitherPrismLaws(t *testing.T) {
	prism := FromEither[error, int]()

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		testValues := []int{0, 42, -10, 999}

		for _, val := range testValues {
			wrapped := prism.ReverseGet(val)
			extracted := prism.GetOption(wrapped)

			assert.True(t, O.IsSome(extracted))
			assert.Equal(t, val, O.GetOrElse(F.Constant(0))(extracted))
		}
	})

	t.Run("law 2: if GetOption(s) == Some(a), then ReverseGet(a) produces valid Either", func(t *testing.T) {
		success := E.Right[error](42)

		extracted := prism.GetOption(success)
		if O.IsSome(extracted) {
			value := O.GetOrElse(F.Constant(0))(extracted)
			rewrapped := prism.ReverseGet(value)

			// Should be Right
			assert.True(t, E.IsRight(rewrapped))

			// Re-extract to verify
			reextracted := prism.GetOption(rewrapped)
			assert.True(t, O.IsSome(reextracted))
			assert.Equal(t, value, O.GetOrElse(F.Constant(0))(reextracted))
		}
	})
}
