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

		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(alwaysFailPred)

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
		optCodec := Optional(S.Monoid, alwaysFailInner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

		encoded := optCodec.Encode(option.Some("hello"))

		// S.Monoid concatenates: inner.Encode("hello") + pred.Encode(true) = "hello" + "true"
		assert.Equal(t, "hellotrue", encoded)
	})

	t.Run("encodes Some with empty inner value", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

		res := optCodec.Is(option.Some("hello"))
		assert.True(t, either.IsRight(res))
	})

	t.Run("Is succeeds for option.None", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional(S.Monoid, inner)(pred)

		res := optCodec.Is(option.None[string]())
		assert.True(t, either.IsRight(res))
	})

	t.Run("Is fails for a raw string (not an Option)", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional(S.Monoid, inner)(pred)

		res := optCodec.Is("not an option")
		assert.True(t, either.IsLeft(res))
	})

	t.Run("Is fails for an int (not an Option)", func(t *testing.T) {
		pred := boolFromOptionalTestInput()
		inner := stringFromOptionalTestInput()
		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

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
		optCodec := Optional(S.Monoid, inner)(pred)

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

// eitherTestInput carries a discriminant flag and one value slot for each
// branch, used to exercise EitherOf with a concrete input/output type.
type eitherTestInput struct {
	isRight  bool
	leftVal  string
	rightVal int
}

// boolFromEitherTestInput is a Type[bool, string, eitherTestInput] that reads
// the "isRight" flag and encodes it as "true"/"false".
func boolFromEitherTestInput() Type[bool, string, eitherTestInput] {
	return MakeType(
		"BoolFromEitherInput",
		Is[bool](),
		func(dep eitherTestInput) Decode[Context, bool] {
			return func(ctx Context) validation.Validation[bool] {
				return validation.Success(dep.isRight)
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

// stringFromEitherTestInput is a Type[string, string, eitherTestInput] that
// reads the "leftVal" field and encodes it as-is.
func stringFromEitherTestInput() Type[string, string, eitherTestInput] {
	return MakeType(
		"StringFromEitherInput",
		Is[string](),
		func(dep eitherTestInput) Decode[Context, string] {
			return func(ctx Context) validation.Validation[string] {
				return validation.Success(dep.leftVal)
			}
		},
		F.Identity[string],
	)
}

// intFromEitherTestInput is a Type[int, string, eitherTestInput] that reads
// the "rightVal" field and encodes it via strconv.Itoa.
func intFromEitherTestInput() Type[int, string, eitherTestInput] {
	return MakeType(
		"IntFromEitherInput",
		Is[int](),
		func(dep eitherTestInput) Decode[Context, int] {
			return func(ctx Context) validation.Validation[int] {
				return validation.Success(dep.rightVal)
			}
		},
		strconv.Itoa,
	)
}

// TestEitherOf_Naming verifies that the codec name contains the predicate and
// both branch codec names.
func TestEitherOf_Naming(t *testing.T) {
	t.Run("name includes predicate and both branch codec names", func(t *testing.T) {
		pred := boolFromEitherTestInput()
		onLeft := stringFromEitherTestInput()
		onRight := intFromEitherTestInput()

		c := EitherOf(S.Monoid, onLeft, onRight)(pred)

		name := c.Name()
		assert.Contains(t, name, "EitherOf[")
		assert.Contains(t, name, pred.Name())
		assert.Contains(t, name, onLeft.Name())
		assert.Contains(t, name, onRight.Name())
	})
}

// TestEitherOf_Decoding_Right verifies that when the predicate decodes to true
// the right codec is invoked and the result is wrapped in either.Right.
func TestEitherOf_Decoding_Right(t *testing.T) {
	t.Run("returns Right when predicate is true", func(t *testing.T) {
		c := EitherOf(S.Monoid, stringFromEitherTestInput(), intFromEitherTestInput())(boolFromEitherTestInput())

		input := eitherTestInput{isRight: true, leftVal: "ignored", rightVal: 42}
		res := c.Decode(input)

		assert.True(t, either.IsRight(res))
		decoded := either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)
		assert.Equal(t, either.Right[string](42), decoded)
	})

	t.Run("Right branch carries the decoded int value", func(t *testing.T) {
		c := EitherOf(S.Monoid, stringFromEitherTestInput(), intFromEitherTestInput())(boolFromEitherTestInput())

		for _, v := range []int{0, 1, -99, 1000} {
			input := eitherTestInput{isRight: true, rightVal: v}
			res := c.Decode(input)

			assert.True(t, either.IsRight(res), "expected Right for rightVal %d", v)
			decoded := either.MonadFold(res,
				func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
				F.Identity[either.Either[string, int]],
			)
			assert.Equal(t, either.Right[string](v), decoded, "rightVal %d", v)
		}
	})
}

// TestEitherOf_Decoding_Left verifies that when the predicate decodes to false
// the left codec is invoked and the result is wrapped in either.Left.
func TestEitherOf_Decoding_Left(t *testing.T) {
	t.Run("returns Left when predicate is false", func(t *testing.T) {
		c := EitherOf(S.Monoid, stringFromEitherTestInput(), intFromEitherTestInput())(boolFromEitherTestInput())

		input := eitherTestInput{isRight: false, leftVal: "hello", rightVal: 99}
		res := c.Decode(input)

		assert.True(t, either.IsRight(res))
		decoded := either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)
		assert.Equal(t, either.Left[int]("hello"), decoded)
	})

	t.Run("Left branch carries the decoded string value", func(t *testing.T) {
		c := EitherOf(S.Monoid, stringFromEitherTestInput(), intFromEitherTestInput())(boolFromEitherTestInput())

		for _, s := range []string{"", "foo", "bar baz"} {
			input := eitherTestInput{isRight: false, leftVal: s}
			res := c.Decode(input)

			assert.True(t, either.IsRight(res), "expected successful decode for leftVal %q", s)
			decoded := either.MonadFold(res,
				func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
				F.Identity[either.Either[string, int]],
			)
			assert.Equal(t, either.Left[int](s), decoded, "leftVal %q", s)
		}
	})
}

// TestEitherOf_Decoding_ValidationFailure verifies that a failed predicate or
// branch validation propagates as Left in the outer Validation result.
func TestEitherOf_Decoding_ValidationFailure(t *testing.T) {
	t.Run("returns Left when predicate validation fails", func(t *testing.T) {
		alwaysFailPred := MakeType(
			"AlwaysFailPred",
			Is[bool](),
			func(dep eitherTestInput) Decode[Context, bool] {
				return validation.FailureWithMessage[bool](dep, "pred always fails")
			},
			func(b bool) string { return "" },
		)
		c := EitherOf(S.Monoid, stringFromEitherTestInput(), intFromEitherTestInput())(alwaysFailPred)

		res := c.Decode(eitherTestInput{isRight: true})
		assert.True(t, either.IsLeft(res))
	})

	t.Run("returns Left when right branch validation fails", func(t *testing.T) {
		alwaysFailRight := MakeType(
			"AlwaysFailRight",
			Is[int](),
			func(dep eitherTestInput) Decode[Context, int] {
				return validation.FailureWithMessage[int](dep, "right always fails")
			},
			strconv.Itoa,
		)
		c := EitherOf(S.Monoid, stringFromEitherTestInput(), alwaysFailRight)(boolFromEitherTestInput())

		// pred=true routes to right, which fails
		res := c.Decode(eitherTestInput{isRight: true})
		assert.True(t, either.IsLeft(res))
	})

	t.Run("returns Left when left branch validation fails", func(t *testing.T) {
		alwaysFailLeft := MakeType(
			"AlwaysFailLeft",
			Is[string](),
			func(dep eitherTestInput) Decode[Context, string] {
				return validation.FailureWithMessage[string](dep, "left always fails")
			},
			F.Identity[string],
		)
		c := EitherOf(S.Monoid, alwaysFailLeft, intFromEitherTestInput())(boolFromEitherTestInput())

		// pred=false routes to left, which fails
		res := c.Decode(eitherTestInput{isRight: false})
		assert.True(t, either.IsLeft(res))
	})
}

// TestEitherOf_Encoding_Right verifies that encoding Right(r) calls the right
// branch encoder and pred encoder for true, combined by the monoid.
func TestEitherOf_Encoding_Right(t *testing.T) {
	t.Run("encodes Right by combining right value and predicate true outputs", func(t *testing.T) {
		c := EitherOf(S.Monoid, stringFromEitherTestInput(), intFromEitherTestInput())(boolFromEitherTestInput())

		// onRight.Encode(42) = "42", pred.Encode(true) = "true"  → "42true"
		encoded := c.Encode(either.Right[string](42))
		assert.Equal(t, "42true", encoded)
	})
}

// TestEitherOf_Encoding_Left verifies that encoding Left(l) calls the left
// branch encoder and pred encoder for false, combined by the monoid.
func TestEitherOf_Encoding_Left(t *testing.T) {
	t.Run("encodes Left by combining left value and predicate false outputs", func(t *testing.T) {
		c := EitherOf(S.Monoid, stringFromEitherTestInput(), intFromEitherTestInput())(boolFromEitherTestInput())

		// onLeft.Encode("hello") = "hello", pred.Encode(false) = "false" → "hellofalse"
		encoded := c.Encode(either.Left[int]("hello"))
		assert.Equal(t, "hellofalse", encoded)
	})
}

// TestEitherOf_TypeChecking verifies that the Is function correctly type-checks
// Either[L, R] values.
func TestEitherOf_TypeChecking(t *testing.T) {
	c := EitherOf(S.Monoid, stringFromEitherTestInput(), intFromEitherTestInput())(boolFromEitherTestInput())

	t.Run("Is succeeds for either.Right", func(t *testing.T) {
		assert.True(t, either.IsRight(c.Is(either.Right[string](1))))
	})

	t.Run("Is succeeds for either.Left", func(t *testing.T) {
		assert.True(t, either.IsRight(c.Is(either.Left[int]("e"))))
	})

	t.Run("Is fails for a plain string (not an Either)", func(t *testing.T) {
		assert.True(t, either.IsLeft(c.Is("not an either")))
	})

	t.Run("Is fails for an int (not an Either)", func(t *testing.T) {
		assert.True(t, either.IsLeft(c.Is(99)))
	})
}

// TestEitherOf_RoundTrip verifies encode-then-decode consistency.
func TestEitherOf_RoundTrip(t *testing.T) {
	c := EitherOf(S.Monoid, stringFromEitherTestInput(), intFromEitherTestInput())(boolFromEitherTestInput())

	t.Run("Right round-trip: encoding contains int and true marker", func(t *testing.T) {
		encoded := c.Encode(either.Right[string](7))
		assert.Contains(t, encoded, "7")
		assert.Contains(t, encoded, "true")
	})

	t.Run("Left round-trip: encoding contains string and false marker", func(t *testing.T) {
		encoded := c.Encode(either.Left[int]("err"))
		assert.Contains(t, encoded, "err")
		assert.Contains(t, encoded, "false")
	})

	t.Run("decode Right then encode preserves Right", func(t *testing.T) {
		input := eitherTestInput{isRight: true, rightVal: 5}
		res := c.Decode(input)
		assert.True(t, either.IsRight(res))

		decoded := either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)
		assert.Equal(t, either.Right[string](5), decoded)
		assert.Contains(t, c.Encode(decoded), "true")
	})

	t.Run("decode Left then encode preserves Left", func(t *testing.T) {
		input := eitherTestInput{isRight: false, leftVal: "msg"}
		res := c.Decode(input)
		assert.True(t, either.IsRight(res))

		decoded := either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)
		assert.Equal(t, either.Left[int]("msg"), decoded)
		assert.Contains(t, c.Encode(decoded), "false")
	})
}

// ExampleEitherOf demonstrates using EitherOf to build a codec that
// dispatches between a string (Left) and an int (Right) branch based on a
// boolean flag parsed from the input string.
//
// The input format is "<flag>:<value>" where flag is "true" or "false".
// "true"  routes to the right (int) branch.
// "false" routes to the left (string) branch.
func ExampleEitherOf() {
	// BoolFromString decodes "true"/"false" and encodes bool back to string.
	// Id[string]()    encodes string → string (identity).
	// IntFromString() decodes a string digit and encodes int → string.
	//
	// We wire them together so:
	//   pred=true  → Right branch (int)
	//   pred=false → Left  branch (string)
	c := F.Pipe1(
		BoolFromString(),
		EitherOf(S.Monoid, Id[string](), IntFromString()),
	)

	// Encode Right(42): IntFromString encodes 42 → "42", pred encodes true → "true"
	// S.Monoid concatenates: "42" + "true" = "42true"
	fmt.Println(c.Encode(either.Right[string](42)))

	// Encode Left("err"): Id encodes "err" → "err", pred encodes false → "false"
	// S.Monoid concatenates: "err" + "false" = "errfalse"
	fmt.Println(c.Encode(either.Left[int]("err")))

	// Output:
	// 42true
	// errfalse
}
