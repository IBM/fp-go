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
	"github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// optionalTestInput is a simple struct that carries both a presence flag and a
// string value, used to exercise Optional with a concrete input/output type.
type optionalTestInput struct {
	present bool
	value   string
}

// boolFromOptionalTestInput is a Type[bool, string, optionalTestInput] that reads
// the "present" field of an optionalTestInput and encodes it as "true"/"false".
func boolFromOptionalTestInput() Type[bool, string, optionalTestInput] {
	return MakeType(
		"BoolFromInput",
		Is[bool](),
		func(dep optionalTestInput) Decode[Context, bool] {
			return func(ctx Context) validation.Validation[bool] {
				return validation.Success(dep.present)
			}
		},
		func(b bool) string {
			if b {
				return "true"
			}
			return "false"
		},
	)
}

// stringFromOptionalTestInput is a Type[string, string, optionalTestInput] that reads
// the "value" field of an optionalTestInput and encodes it as-is.
func stringFromOptionalTestInput() Type[string, string, optionalTestInput] {
	return MakeType(
		"StringFromInput",
		Is[string](),
		func(dep optionalTestInput) Decode[Context, string] {
			return func(ctx Context) validation.Validation[string] {
				return validation.Success(dep.value)
			}
		},
		F.Identity[string],
	)
}

// TestOptional_Naming verifies that the codec name contains both the predicate
// and the inner codec names.
func TestOptional_Naming(t *testing.T) {
	t.Run("name includes predicate and inner codec names", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()

		optCodec := Optional[string](S.Monoid, pred)(inner)

		name := optCodec.Name()
		assert.Contains(t, name, "Optional[")
		assert.Contains(t, name, pred.Name())
		assert.Contains(t, name, inner.Name())
	})
}

// TestOptional_Decoding_PresentTrue verifies that when the predicate decodes to
// true the inner codec is invoked and the result is wrapped in option.Some.
func TestOptional_Decoding_PresentTrue(t *testing.T) {
	t.Run("returns Some when predicate is true", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		input := optionalTestInput{present: true, value: "hello"}
		result := optCodec.Decode(input)

		assert.True(t, either.IsRight(result))
		decoded := either.MonadFold(result,
			func(validation.Errors) option.Option[string] { return option.None[string]() },
			F.Identity[option.Option[string]],
		)
		assert.Equal(t, option.Some("hello"), decoded)
	})

	t.Run("decoded Some contains the inner value", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		for _, value := range []string{"", "world", "foo bar"} {
			input := optionalTestInput{present: true, value: value}
			result := optCodec.Decode(input)

			assert.True(t, either.IsRight(result), "expected Right for value %q", value)
			decoded := either.MonadFold(result,
				func(validation.Errors) option.Option[string] { return option.None[string]() },
				F.Identity[option.Option[string]],
			)
			assert.Equal(t, option.Some(value), decoded, "value %q", value)
		}
	})
}

// TestOptional_Decoding_PresentFalse verifies that when the predicate decodes to
// false the result is option.None regardless of the inner value.
func TestOptional_Decoding_PresentFalse(t *testing.T) {
	t.Run("returns None when predicate is false", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		input := optionalTestInput{present: false, value: "ignored"}
		result := optCodec.Decode(input)

		assert.True(t, either.IsRight(result))
		decoded := either.MonadFold(result,
			func(validation.Errors) option.Option[string] { return option.None[string]() },
			F.Identity[option.Option[string]],
		)
		assert.Equal(t, option.None[string](), decoded)
	})

	t.Run("None is returned regardless of inner value when predicate is false", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		input := optionalTestInput{present: false, value: "some value that must be ignored"}
		result := optCodec.Decode(input)

		assert.True(t, either.IsRight(result))
		decoded := either.MonadFold(result,
			func(validation.Errors) option.Option[string] { return option.None[string]() },
			F.Identity[option.Option[string]],
		)
		assert.Equal(t, option.None[string](), decoded)
	})
}

// TestOptional_Decoding_ValidationFailure verifies that a failed predicate
// validation propagates as Left.
func TestOptional_Decoding_ValidationFailure(t *testing.T) {
	t.Run("returns Left when predicate validation fails", func(t *testing.T) {
		// A predicate codec that always fails decoding.
		alwaysFailPred := MakeType(
			"AlwaysFail",
			Is[bool](),
			func(dep optionalTestInput) Decode[Context, bool] {
				return validation.FailureWithMessage[bool](dep, "predicate always fails")
			},
			func(b bool) string { return "" },
		)
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, alwaysFailPred)(inner)

		input := optionalTestInput{present: true, value: "hello"}
		result := optCodec.Decode(input)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("returns Left when inner codec validation fails on Some path", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		// Inner codec that always fails decoding.
		alwaysFailInner := MakeType(
			"AlwaysFailInner",
			Is[string](),
			func(dep optionalTestInput) Decode[Context, string] {
				return validation.FailureWithMessage[string](dep, "inner always fails")
			},
			F.Identity[string],
		)
		optCodec := Optional[string](S.Monoid, pred)(alwaysFailInner)

		// Predicate succeeds (present=true), but inner codec fails.
		input := optionalTestInput{present: true, value: "hello"}
		result := optCodec.Decode(input)

		assert.True(t, either.IsLeft(result))
	})
}

// TestOptional_Encoding_Some verifies that encoding Some(a) calls the inner
// codec's encoder and the predicate's encoder, combining them with the monoid.
func TestOptional_Encoding_Some(t *testing.T) {
	t.Run("encodes Some by combining inner value and predicate true outputs", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		encoded := optCodec.Encode(option.Some("hello"))

		// S.Monoid concatenates: inner.Encode("hello") + pred.Encode(true) = "hello" + "true"
		assert.Equal(t, "hellotrue", encoded)
	})

	t.Run("encodes Some with empty inner value", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		encoded := optCodec.Encode(option.Some(""))
		// inner.Encode("") = "", pred.Encode(true) = "true"
		assert.Equal(t, "true", encoded)
	})
}

// TestOptional_Encoding_None verifies that encoding None uses the monoid empty
// for the inner value and the predicate's encoder for false.
func TestOptional_Encoding_None(t *testing.T) {
	t.Run("encodes None using monoid empty and predicate false output", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		encoded := optCodec.Encode(option.None[string]())

		// inner value absent => m.Empty() = "", pred.Encode(false) = "false"
		// S.Monoid.Concat("", "false") = "false"
		assert.Equal(t, "false", encoded)
	})
}

// TestOptional_TypeChecking verifies that the Is function correctly type-checks
// Option[A] values.
func TestOptional_TypeChecking(t *testing.T) {
	t.Run("Is succeeds for option.Some", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		res := optCodec.Is(option.Some("hello"))
		assert.True(t, either.IsRight(res))
	})

	t.Run("Is succeeds for option.None", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		res := optCodec.Is(option.None[string]())
		assert.True(t, either.IsRight(res))
	})

	t.Run("Is fails for a raw string (not an Option)", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		res := optCodec.Is("not an option")
		assert.True(t, either.IsLeft(res))
	})

	t.Run("Is fails for an int (not an Option)", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		res := optCodec.Is(42)
		assert.True(t, either.IsLeft(res))
	})
}

// TestOptional_RoundTrip verifies encode-then-decode consistency for the
// typical case where the output type equals the input type.
func TestOptional_RoundTrip(t *testing.T) {
	// This test exercises Optional with I = bool so pred can be Bool() and
	// onSome can be Id[bool](), demonstrating a self-referential round-trip.
	//
	// pred  : Bool()  – Type[bool, bool, bool] (identity)
	// inner : Id[bool]() – Type[bool, bool, bool] (identity)
	// m     : BoolMonoid (AND) – not available; use a custom OR-based monoid to
	//         keep things simple.
	//
	// Instead we test with a bespoke pair of codecs sharing I = string (pure
	// string-based inputs) to verify decode→encode gives back the original string.

	// Codec whose input is a raw string representing a JSON-like "present:value" pair
	// is too complex; keep the round-trip test focused on encoding idempotence.

	t.Run("encode(Some(v)).isSome mirrors decode.present=true", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		// Encoding Some("x") must contain "true"
		encoded := optCodec.Encode(option.Some("x"))
		assert.Contains(t, encoded, "true")

		// Encoding None must contain "false"
		encodedNone := optCodec.Encode(option.None[string]())
		assert.Contains(t, encodedNone, "false")
	})

	t.Run("decoding present=true followed by encoding preserves Some", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		input := optionalTestInput{present: true, value: "roundtrip"}
		decodeResult := optCodec.Decode(input)
		assert.True(t, either.IsRight(decodeResult))

		decoded := either.MonadFold(decodeResult,
			func(validation.Errors) option.Option[string] { return option.None[string]() },
			F.Identity[option.Option[string]],
		)
		assert.True(t, option.IsSome(decoded))

		encoded := optCodec.Encode(decoded)
		assert.Contains(t, encoded, "roundtrip")
		assert.Contains(t, encoded, "true")
	})

	t.Run("decoding present=false followed by encoding yields false marker", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional[string](S.Monoid, pred)(inner)

		input := optionalTestInput{present: false, value: "ignored"}
		decodeResult := optCodec.Decode(input)
		assert.True(t, either.IsRight(decodeResult))

		decoded := either.MonadFold(decodeResult,
			func(validation.Errors) option.Option[string] { return option.None[string]() },
			F.Identity[option.Option[string]],
		)
		assert.False(t, option.IsSome(decoded))

		encoded := optCodec.Encode(decoded)
		assert.Contains(t, encoded, "false")
	})
}
