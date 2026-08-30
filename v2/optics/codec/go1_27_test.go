//go:build go1.27

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
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/stretchr/testify/assert"
)

// positiveIntPrism is a Prism that matches integers > 0 and reviews them unchanged.
var positiveIntPrism = prism.MakePrismWithName(
	func(n int) option.Option[int] {
		if n > 0 {
			return option.Some(n)
		}
		return option.None[int]()
	},
	func(n int) int { return n },
	"PositiveInt",
)

// positiveIntCodec is a codec that validates int values are positive.
var positiveIntCodec = FromRefinement(positiveIntPrism)

// TestTypePipe_MethodEquivalentToFreeFunction verifies that the method form
// base.Pipe(ab) returns a codec identical in behaviour to the free function
// Pipe[O, I](ab)(base) for both Decode and Encode.
func TestTypePipe_MethodEquivalentToFreeFunction(t *testing.T) {
	method := IntFromString().Pipe(positiveIntCodec)
	free := Pipe[string, string](positiveIntCodec)(IntFromString())

	t.Run("Decode parity on success", func(t *testing.T) {
		assert.Equal(t, free.Decode("42"), method.Decode("42"))
	})

	t.Run("Decode parity on first-codec failure", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(free.Decode("abc")), either.IsLeft(method.Decode("abc")))
	})

	t.Run("Decode parity on second-codec failure", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(free.Decode("-1")), either.IsLeft(method.Decode("-1")))
	})

	t.Run("Encode parity", func(t *testing.T) {
		assert.Equal(t, free.Encode(7), method.Encode(7))
	})
}

// TestTypePipe_Method_DecodeSuccess verifies that Decode succeeds when both
// codecs accept the input.
func TestTypePipe_Method_DecodeSuccess(t *testing.T) {
	composed := IntFromString().Pipe(positiveIntCodec)

	assert.Equal(t, validation.Success(42), composed.Decode("42"))
	assert.Equal(t, validation.Success(1), composed.Decode("1"))
	assert.Equal(t, validation.Success(100), composed.Decode("100"))
}

// TestTypePipe_Method_DecodeFailsOnFirstCodec verifies that Decode fails when
// the receiver codec rejects the input, before the second codec is applied.
func TestTypePipe_Method_DecodeFailsOnFirstCodec(t *testing.T) {
	composed := IntFromString().Pipe(positiveIntCodec)

	assert.True(t, either.IsLeft(composed.Decode("not-a-number")))
	assert.True(t, either.IsLeft(composed.Decode("")))
	assert.True(t, either.IsLeft(composed.Decode("1.5")))
}

// TestTypePipe_Method_DecodeFailsOnSecondCodec verifies that Decode fails when
// the receiver succeeds but the second codec rejects the intermediate value.
func TestTypePipe_Method_DecodeFailsOnSecondCodec(t *testing.T) {
	composed := IntFromString().Pipe(positiveIntCodec)

	assert.True(t, either.IsLeft(composed.Decode("0")))
	assert.True(t, either.IsLeft(composed.Decode("-1")))
	assert.True(t, either.IsLeft(composed.Decode("-42")))
}

// TestTypePipe_Method_Encode verifies that Encode applies both codecs in
// reverse order: the argument codec encodes B → A, then the receiver
// encodes A → O.
func TestTypePipe_Method_Encode(t *testing.T) {
	composed := IntFromString().Pipe(positiveIntCodec)

	assert.Equal(t, "42", composed.Encode(42))
	assert.Equal(t, "1", composed.Encode(1))
	assert.Equal(t, "100", composed.Encode(100))
}

// TestTypePipe_Method_RoundTrip verifies that Encode(Decode(s)) == s for valid
// inputs, i.e. the composed codec is a true round-trip codec.
func TestTypePipe_Method_RoundTrip(t *testing.T) {
	composed := IntFromString().Pipe(positiveIntCodec)

	for _, s := range []string{"1", "7", "42", "999"} {
		t.Run(s, func(t *testing.T) {
			res := composed.Decode(s)
			assert.True(t, either.IsRight(res))
			encoded := either.MonadFold(res,
				func(validation.Errors) string { return "" },
				composed.Encode,
			)
			assert.Equal(t, s, encoded)
		})
	}
}

// TestTypePipe_Method_Name verifies that the composed codec's name is
// "Pipe(<receiver-name>, <ab-name>)".
func TestTypePipe_Method_Name(t *testing.T) {
	composed := IntFromString().Pipe(positiveIntCodec)
	assert.Equal(t, "Pipe(IntFromString, FromRefinement(PositiveInt))", composed.Name())
}

// TestTypePipe_Method_Chained verifies that three .Pipe calls chain correctly
// through three decode steps.
func TestTypePipe_Method_Chained(t *testing.T) {
	// nonEmptyString → IntFromString → positiveInt
	// Decodes: string → non-empty string → int → positive int
	composed := NonEmptyString().
		Pipe(IntFromString()).
		Pipe(positiveIntCodec)

	t.Run("decodes a positive integer string", func(t *testing.T) {
		assert.Equal(t, validation.Success(7), composed.Decode("7"))
	})

	t.Run("fails on empty string (first codec)", func(t *testing.T) {
		assert.True(t, either.IsLeft(composed.Decode("")))
	})

	t.Run("fails on non-numeric string (second codec)", func(t *testing.T) {
		assert.True(t, either.IsLeft(composed.Decode("abc")))
	})

	t.Run("fails on zero (third codec)", func(t *testing.T) {
		assert.True(t, either.IsLeft(composed.Decode("0")))
	})

	t.Run("encodes back through all three levels", func(t *testing.T) {
		assert.Equal(t, "42", composed.Encode(42))
	})
}

// ExampleType_Pipe demonstrates chaining two codecs with the method form.
// IntFromString decodes a string to int; the positive-int refinement then
// validates that the integer is greater than zero.
func ExampleType_Pipe() {
	positiveFromString := IntFromString().Pipe(FromRefinement(positiveIntPrism))

	fmt.Println(positiveFromString.Name())
	fmt.Println(positiveFromString.Encode(42))
	fmt.Println(either.IsRight(positiveFromString.Decode("10")))
	fmt.Println(either.IsLeft(positiveFromString.Decode("-3")))

	// Output:
	// Pipe(IntFromString, FromRefinement(PositiveInt))
	// 42
	// true
	// true
}
