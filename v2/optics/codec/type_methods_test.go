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
	"log/slog"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// AsEncode
// ---------------------------------------------------------------------------

// TestType_AsEncode verifies that AsEncode returns the Encode function field
// and that calling it is equivalent to calling Encode directly.
func TestType_AsEncode(t *testing.T) {
	t.Run("returns the encode function for IntFromString", func(t *testing.T) {
		c := IntFromString()
		encodeFn := c.AsEncode()

		assert.Equal(t, c.Encode(42), encodeFn(42))
		assert.Equal(t, c.Encode(-7), encodeFn(-7))
		assert.Equal(t, c.Encode(0), encodeFn(0))
	})

	t.Run("returned function is independent of the original codec value", func(t *testing.T) {
		encodeFn := IntFromString().AsEncode()
		assert.Equal(t, "99", encodeFn(99))
	})

	t.Run("works for identity codec", func(t *testing.T) {
		c := String()
		encodeFn := c.AsEncode()
		assert.Equal(t, "hello", encodeFn("hello"))
	})
}

// ---------------------------------------------------------------------------
// AsValidate
// ---------------------------------------------------------------------------

// TestType_AsValidate verifies that AsValidate returns the Validate function
// field and that calling it produces the same result as going through
// Decoder.Validate directly.
func TestType_AsValidate(t *testing.T) {
	t.Run("returns the validate function for IntFromString", func(t *testing.T) {
		c := IntFromString()
		validateFn := c.AsValidate()

		ctx := emptyContext
		assert.Equal(t, c.Validate("42")(ctx), validateFn("42")(ctx))
	})

	t.Run("validate function succeeds on valid input", func(t *testing.T) {
		validateFn := IntFromString().AsValidate()
		result := validateFn("100")(emptyContext)
		assert.Equal(t, validation.Success(100), result)
	})

	t.Run("validate function fails on invalid input", func(t *testing.T) {
		validateFn := IntFromString().AsValidate()
		result := validateFn("abc")(emptyContext)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("returned function is independent of the original codec value", func(t *testing.T) {
		validateFn := Int().AsValidate()
		result := validateFn(42)(emptyContext)
		assert.Equal(t, validation.Success(42), result)
	})
}

// ---------------------------------------------------------------------------
// AsDecode
// ---------------------------------------------------------------------------

// TestType_AsDecode verifies that AsDecode returns the Decode function field
// and that calling it is equivalent to calling Decode directly on the codec.
func TestType_AsDecode(t *testing.T) {
	t.Run("returns the decode function for IntFromString", func(t *testing.T) {
		c := IntFromString()
		decodeFn := c.AsDecode()

		assert.Equal(t, c.Decode("42"), decodeFn("42"))
		assert.Equal(t, c.Decode("abc"), decodeFn("abc"))
	})

	t.Run("decode function succeeds on valid input", func(t *testing.T) {
		decodeFn := IntFromString().AsDecode()
		assert.Equal(t, validation.Success(42), decodeFn("42"))
	})

	t.Run("decode function fails on invalid input", func(t *testing.T) {
		decodeFn := IntFromString().AsDecode()
		assert.True(t, either.IsLeft(decodeFn("not-a-number")))
	})

	t.Run("returned function is independent of the original codec value", func(t *testing.T) {
		decodeFn := String().AsDecode()
		assert.Equal(t, validation.Success("hello"), decodeFn("hello"))
	})

	t.Run("differs from AsValidate in that it provides no context argument", func(t *testing.T) {
		c := IntFromString()
		// AsDecode: Decode[I, A] = func(I) Validation[A]
		decodeFn := c.AsDecode()
		// AsValidate: Validate[I, A] = func(I) Reader[Context, Validation[A]]
		validateFn := c.AsValidate()

		// Both should agree on the successful result
		decodeResult := decodeFn("7")
		validateResult := validateFn("7")(emptyContext)
		assert.Equal(t, either.IsRight(decodeResult), either.IsRight(validateResult))
	})
}

// ---------------------------------------------------------------------------
// Example functions for all public Type methods
// ---------------------------------------------------------------------------

// ExampleType_Name shows that Name returns the descriptive name assigned at
// construction time and used in validation error messages.
func ExampleType_Name() {
	fmt.Println(String().Name())
	fmt.Println(Int().Name())
	fmt.Println(IntFromString().Name())
	fmt.Println(Array(Int()).Name())

	// Output:
	// string
	// int
	// IntFromString
	// Array[int]
}

// ExampleType_Decode shows that Decode validates the input and returns
// Either[Errors, A]: Right on success, Left on failure.
func ExampleType_Decode() {
	c := IntFromString()

	fmt.Println(either.IsRight(c.Decode("42")))
	fmt.Println(either.IsLeft(c.Decode("abc")))

	// Output:
	// true
	// true
}

// ExampleType_Encode shows that Encode converts A back to O without error.
func ExampleType_Encode() {
	fmt.Println(IntFromString().Encode(42))
	fmt.Println(String().Encode("hello"))
	fmt.Println(Bool().Encode(true))

	// Output:
	// 42
	// hello
	// true
}

// ExampleType_Validate shows the context-aware form of decoding.
// Unlike Decode, Validate returns a function that accepts a Context, allowing
// the caller to thread an existing validation path into the result.
func ExampleType_Validate() {
	c := IntFromString()

	validateFn := c.Validate("42")     // returns Reader[Context, Validation[int]]
	result := validateFn(emptyContext) // apply with an empty context

	fmt.Println(either.IsRight(result))

	// Output:
	// true
}

// ExampleType_Is shows that Is performs a runtime type assertion against A.
// It returns Right(a) when the value is of type A and Left(err) otherwise.
func ExampleType_Is() {
	c := Int()

	fmt.Println(either.IsRight(c.Is(42)))
	fmt.Println(either.IsLeft(c.Is("42")))

	// Output:
	// true
	// true
}

// ExampleType_AsDecoder shows that AsDecoder extracts the embedded Decoder
// struct, which can then be passed to functions expecting a Decoder[I, A].
func ExampleType_AsDecoder() {
	dec := IntFromString().AsDecoder()

	fmt.Println(either.IsRight(dec.Decode("10")))
	fmt.Println(either.IsLeft(dec.Decode("x")))

	// Output:
	// true
	// true
}

// ExampleType_AsEncoder shows that AsEncoder extracts the embedded Encoder
// struct, which can then be passed to functions expecting an Encoder[A, O].
func ExampleType_AsEncoder() {
	enc := IntFromString().AsEncoder()

	fmt.Println(enc.Encode(7))
	fmt.Println(enc.Encode(-3))

	// Output:
	// 7
	// -3
}

// ExampleType_AsEncode shows that AsEncode extracts the bare Encode[A, O]
// function from the codec, useful when a plain function value is needed.
func ExampleType_AsEncode() {
	encodeFn := IntFromString().AsEncode()

	fmt.Println(encodeFn(100))
	fmt.Println(encodeFn(0))

	// Output:
	// 100
	// 0
}

// ExampleType_AsValidate shows that AsValidate extracts the bare
// Validate[I, A] function — the context-aware form of decoding.
func ExampleType_AsValidate() {
	validateFn := IntFromString().AsValidate()

	result := validateFn("55")(emptyContext)
	fmt.Println(either.IsRight(result))

	// Output:
	// true
}

// ExampleType_AsDecode shows that AsDecode extracts the bare Decode[I, A]
// function — equivalent to calling Decode directly on the codec.
func ExampleType_AsDecode() {
	decodeFn := IntFromString().AsDecode()

	fmt.Println(either.IsRight(decodeFn("42")))
	fmt.Println(either.IsLeft(decodeFn("oops")))

	// Output:
	// true
	// true
}

// ExampleType_String shows that String() returns the codec name, satisfying
// the fmt.Stringer interface.
func ExampleType_String() {
	fmt.Println(IntFromString().String())
	fmt.Println(Array(Bool()).String())

	// Output:
	// IntFromString
	// Array[bool]
}

// ExampleType_Format shows the fmt.Formatter implementation: %s and %v use
// the codec name, while %q quotes it.
func ExampleType_Format() {
	c := IntFromString()

	fmt.Printf("%s\n", c)
	fmt.Printf("%v\n", c)
	fmt.Printf("%q\n", c)

	// Output:
	// IntFromString
	// IntFromString
	// "IntFromString"
}

// ExampleType_LogValue shows that LogValue returns a slog.StringValue
// containing the codec name, satisfying the slog.LogValuer interface.
func ExampleType_LogValue() {
	v := IntFromString().LogValue()

	fmt.Println(v.Kind() == slog.KindString)
	fmt.Println(v.String())

	// Output:
	// true
	// IntFromString
}

// ExampleMakeType shows how to construct a custom Type using MakeType.
func ExampleMakeType() {
	double := MakeType(
		"Double",
		Is[int](),
		func(i int) func([]validation.ContextEntry) validation.Validation[int] {
			return func(_ []validation.ContextEntry) validation.Validation[int] {
				return validation.Success(i * 2)
			}
		},
		func(n int) int { return n / 2 },
	)

	fmt.Println(double.Name())
	fmt.Println(either.IsRight(double.Decode(21)))
	fmt.Println(double.Encode(42))

	// Output:
	// Double
	// true
	// 21
}

// ExampleString shows that String() returns a Type that validates and encodes
// string values with identity semantics.
func ExampleString() {
	c := String()

	fmt.Println(c.Name())
	fmt.Println(c.Encode("hello"))
	fmt.Println(either.IsRight(c.Decode("world")))
	fmt.Println(either.IsLeft(c.Decode(42)))

	// Output:
	// string
	// hello
	// true
	// true
}

// ExampleInt shows that Int() returns a Type that validates and encodes int
// values with identity semantics.
func ExampleInt() {
	c := Int()

	fmt.Println(c.Name())
	fmt.Println(c.Encode(99))
	fmt.Println(either.IsRight(c.Decode(7)))
	fmt.Println(either.IsLeft(c.Decode("7")))

	// Output:
	// int
	// 99
	// true
	// true
}

// ExampleBool shows that Bool() returns a Type that validates and encodes bool
// values with identity semantics.
func ExampleBool() {
	c := Bool()

	fmt.Println(c.Name())
	fmt.Println(c.Encode(true))
	fmt.Println(either.IsRight(c.Decode(false)))
	fmt.Println(either.IsLeft(c.Decode(1)))

	// Output:
	// bool
	// true
	// true
	// true
}

// ExampleNil shows that Nil[A]() validates nil values and rejects non-nil ones.
func ExampleNil() {
	c := Nil[string]()

	fmt.Println(c.Name())
	fmt.Println(either.IsRight(c.Decode(nil)))
	fmt.Println(either.IsLeft(c.Decode("not nil")))

	// Output:
	// nil
	// true
	// true
}
