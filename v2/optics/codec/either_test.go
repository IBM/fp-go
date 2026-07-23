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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAltW_Naming verifies that the codec name contains both branch codec names.
func TestAltW_Naming(t *testing.T) {
	t.Run("name follows AltW[left, right] format", func(t *testing.T) {
		c := AltW[int](Id[string]())(IntFromString())
		assert.Equal(t, "AltW[string, IntFromString]", c.Name())
	})
}

// TestAltW_Encode verifies that encoding dispatches to the correct branch encoder.
func TestAltW_Encode(t *testing.T) {
	c := AltW[int](Id[string]())(IntFromString())

	t.Run("encodes Left value using left codec", func(t *testing.T) {
		assert.Equal(t, "hello", c.Encode(either.Left[int]("hello")))
	})

	t.Run("encodes Right value using right codec", func(t *testing.T) {
		assert.Equal(t, "42", c.Encode(either.Right[string](42)))
	})
}

// TestAltW_Decode verifies that decoding routes to the correct branch.
func TestAltW_Decode(t *testing.T) {
	// Id[string] accepts any string; IntFromString accepts digit strings only.
	// Right branch (int) is tried first, so "42" → Right(42), "hello" → Left("hello").
	c := AltW[int](Id[string]())(IntFromString())

	extractInner := func(res Validation[either.Either[string, int]]) either.Either[string, int] {
		t.Helper()
		require.True(t, either.IsRight(res), "outer validation must succeed")
		return either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)
	}

	t.Run("decodes digit string as Right(int)", func(t *testing.T) {
		inner := extractInner(c.Decode("42"))
		assert.Equal(t, either.Right[string](42), inner)
	})

	t.Run("decodes non-digit string as Left(string)", func(t *testing.T) {
		inner := extractInner(c.Decode("hello"))
		assert.Equal(t, either.Left[int]("hello"), inner)
	})
}

// TestAltW_Decode_Prioritization verifies that the Right branch is tried before Left.
func TestAltW_Decode_Prioritization(t *testing.T) {
	c := AltW[int](Id[string]())(IntFromString())

	t.Run("prioritizes Right over Left when both could succeed", func(t *testing.T) {
		// "42" is a valid string AND a valid int; Right must win.
		res := c.Decode("42")
		require.True(t, either.IsRight(res))

		inner := either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)
		assert.True(t, either.IsRight(inner), "Right branch must be preferred")
		assert.Equal(t, either.Right[string](42), inner)
	})

	t.Run("falls back to Left when Right fails", func(t *testing.T) {
		// "hello" cannot be parsed as int, so Left must be used.
		res := c.Decode("hello")
		require.True(t, either.IsRight(res))

		inner := either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)
		assert.True(t, either.IsLeft(inner), "Left branch must be used when Right fails")
		assert.Equal(t, either.Left[int]("hello"), inner)
	})
}

// TestAltW_Decode_BothFail verifies that errors from both branches accumulate
// when neither branch can decode the input.
func TestAltW_Decode_BothFail(t *testing.T) {
	// NonEmptyString rejects empty strings; PositiveInt rejects non-positive values.
	nonEmptyString := MakeType(
		"NonEmptyString",
		Is[string](),
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

	positiveInt := MakeType(
		"PositiveInt",
		Is[int](),
		func(s string) Decode[Context, int] {
			return func(c Context) Validation[int] {
				n, err := strconv.Atoi(s)
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

	c := AltW[int](nonEmptyString)(positiveInt)

	t.Run("outer result is Left when both branches fail", func(t *testing.T) {
		// Empty string fails NonEmptyString (Left) and cannot be parsed as int (Right).
		assert.True(t, either.IsLeft(c.Decode("")))
	})

	t.Run("errors from both branches are accumulated", func(t *testing.T) {
		res := c.Decode("")
		require.True(t, either.IsLeft(res))

		errs := either.MonadFold(res,
			F.Identity[validation.Errors],
			func(either.Either[string, int]) validation.Errors { return nil },
		)
		require.NotNil(t, errs)
		assert.GreaterOrEqual(t, len(errs), 2,
			"should carry at least one error from each branch")

		var msgs []string
		for _, e := range errs {
			msgs = append(msgs, e.Messsage)
		}
		hasRightErr := false
		hasLeftErr := false
		for _, m := range msgs {
			if m == "expected integer string" || m == "must be positive" {
				hasRightErr = true
			}
			if m == "must not be empty" {
				hasLeftErr = true
			}
		}
		assert.True(t, hasRightErr, "error from Right branch must be present")
		assert.True(t, hasLeftErr, "error from Left branch must be present")
	})

	t.Run("non-positive integer fails both branches", func(t *testing.T) {
		// "-1": NonEmptyString succeeds (non-empty), so the overall result
		// should still be Right (Left branch succeeds as fallback for non-positive int).
		// Demonstrates that Left succeeds when Right fails.
		res := c.Decode("-1")
		require.True(t, either.IsRight(res))
		inner := either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)
		// "-1" fails PositiveInt but succeeds NonEmptyString → Left("-1")
		assert.Equal(t, either.Left[int]("-1"), inner)
	})
}

// TestAltW_Decode_CustomCodecs exercises additional validation rules.
func TestAltW_Decode_CustomCodecs(t *testing.T) {
	nonEmptyString := MakeType(
		"NonEmptyString",
		Is[string](),
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

	positiveInt := MakeType(
		"PositiveInt",
		Is[int](),
		func(s string) Decode[Context, int] {
			return func(c Context) Validation[int] {
				n, err := strconv.Atoi(s)
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

	c := AltW[int](nonEmptyString)(positiveInt)

	t.Run("non-empty string succeeds as Left when right fails", func(t *testing.T) {
		assert.True(t, either.IsRight(c.Decode("hello")))
	})

	t.Run("positive integer string succeeds as Right", func(t *testing.T) {
		assert.True(t, either.IsRight(c.Decode("42")))
	})

	t.Run("zero fails Right, succeeds as Left (non-empty string)", func(t *testing.T) {
		res := c.Decode("0")
		require.True(t, either.IsRight(res))
		inner := either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)
		// "0" fails PositiveInt, but is non-empty → Left("0")
		assert.Equal(t, either.Left[int]("0"), inner)
	})
}

// TestAltW_TypeChecking verifies that Is correctly accepts Either values and
// rejects plain values.
func TestAltW_TypeChecking(t *testing.T) {
	c := AltW[int](Id[string]())(IntFromString())

	t.Run("Is succeeds for either.Right", func(t *testing.T) {
		assert.True(t, either.IsRight(c.Is(either.Right[string](1))))
	})

	t.Run("Is succeeds for either.Left", func(t *testing.T) {
		assert.True(t, either.IsRight(c.Is(either.Left[int]("e"))))
	})

	t.Run("Is fails for a plain string", func(t *testing.T) {
		assert.True(t, either.IsLeft(c.Is("not an either")))
	})

	t.Run("Is fails for a plain int", func(t *testing.T) {
		assert.True(t, either.IsLeft(c.Is(42)))
	})
}

// TestAltW_RoundTrip verifies encode-then-decode consistency.
func TestAltW_RoundTrip(t *testing.T) {
	c := AltW[int](Id[string]())(IntFromString())

	extractInner := func(res Validation[either.Either[string, int]]) either.Either[string, int] {
		t.Helper()
		require.True(t, either.IsRight(res))
		return either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)
	}

	t.Run("Left round-trip: decode then encode returns original string", func(t *testing.T) {
		inner := extractInner(c.Decode("hello"))
		assert.Equal(t, "hello", c.Encode(inner))
	})

	t.Run("Right round-trip: decode then encode returns original string", func(t *testing.T) {
		inner := extractInner(c.Decode("42"))
		assert.Equal(t, "42", c.Encode(inner))
	})
}

// ExampleAltW demonstrates using AltW to build a codec that tries to decode
// the input as an int (Right branch) and falls back to a raw string (Left
// branch) when the int parse fails.
func ExampleAltW() {
	// AltW(leftCodec)(rightCodec): right branch tried first.
	// Id[string]() accepts any string and encodes it as-is (Left branch).
	// IntFromString() parses digit strings and encodes int back to string (Right branch).
	c := AltW[int](Id[string]())(IntFromString())

	// "42" parses successfully as int → Right(42) → encoded back as "42"
	fmt.Println(c.Encode(either.Right[string](42)))

	// "hello" fails int parsing → Left("hello") → encoded back as "hello"
	fmt.Println(c.Encode(either.Left[int]("hello")))

	// Output:
	// 42
	// hello
}
