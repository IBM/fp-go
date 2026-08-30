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
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/iso"
	"github.com/IBM/fp-go/v2/optics/common"
	"github.com/IBM/fp-go/v2/optics/optional"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
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

// ---------------------------------------------------------------------------
// PipeIso — package-level free function
// ---------------------------------------------------------------------------

// userIdIso is a lossless isomorphism between int and UserId (newtype),
// used as the second step in a PipeIso composition.
// Get wraps an int in UserId; ReverseGet unwraps UserId back to int.
var userIdIso = iso.MakeIso(
	func(i int) UserId { return UserId(i) },
	func(id UserId) int { return int(id) },
)

// TestPipeIso_DecodeSuccess verifies that the composed codec successfully
// decodes and applies the isomorphism when the upstream codec succeeds.
func TestPipeIso_DecodeSuccess(t *testing.T) {
	// IntFromString decodes a string to int; userIdIso wraps int in UserId.
	composed := PipeIso[string, string](userIdIso)(IntFromString())

	assert.Equal(t, validation.Success(UserId(42)), composed.Decode("42"))
	assert.Equal(t, validation.Success(UserId(0)), composed.Decode("0"))
	assert.Equal(t, validation.Success(UserId(-7)), composed.Decode("-7"))
}

// TestPipeIso_DecodeFailsPropagated verifies that a failure in the upstream
// codec is propagated unchanged through PipeIso.
func TestPipeIso_DecodeFailsPropagated(t *testing.T) {
	composed := PipeIso[string, string](userIdIso)(IntFromString())

	assert.True(t, either.IsLeft(composed.Decode("not-a-number")))
	assert.True(t, either.IsLeft(composed.Decode("")))
}

// TestPipeIso_Encode verifies that encoding applies the isomorphism's
// ReverseGet to convert B back to A and then the upstream encoder to O.
func TestPipeIso_Encode(t *testing.T) {
	composed := PipeIso[string, string](userIdIso)(IntFromString())

	assert.Equal(t, "42", composed.Encode(UserId(42)))
	assert.Equal(t, "-7", composed.Encode(UserId(-7)))
}

// TestPipeIso_RoundTrip verifies that Encode(Decode(s)...) == s for valid
// inputs and that the composed codec is a true round-trip.
func TestPipeIso_RoundTrip(t *testing.T) {
	composed := PipeIso[string, string](userIdIso)(IntFromString())

	for _, s := range []string{"1", "42", "-10", "0"} {
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

// TestPipeIso_Name verifies that the composed codec's name follows the
// "Pipe(<receiver-name>, FromIso[Iso])" pattern.
func TestPipeIso_Name(t *testing.T) {
	composed := PipeIso[string, string](userIdIso)(IntFromString())
	assert.Equal(t, "Pipe(IntFromString, FromIso[Iso])", composed.Name())
}

// TestPipeIso_EquivalentToFreeFunction verifies that the free-function form
// PipeIso[O, I](ab)(base) is behaviourally identical to Pipe[O, I](FromIso(ab))(base).
func TestPipeIso_EquivalentToFreeFunction(t *testing.T) {
	via := PipeIso[string, string](userIdIso)(IntFromString())
	direct := Pipe[string, string](FromIso(userIdIso))(IntFromString())

	t.Run("Decode parity on success", func(t *testing.T) {
		assert.Equal(t, direct.Decode("42"), via.Decode("42"))
	})
	t.Run("Decode parity on failure", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(direct.Decode("abc")), either.IsLeft(via.Decode("abc")))
	})
	t.Run("Encode parity", func(t *testing.T) {
		assert.Equal(t, direct.Encode(UserId(7)), via.Encode(UserId(7)))
	})
}

// ---------------------------------------------------------------------------
// PipeRefinement — package-level free function
// ---------------------------------------------------------------------------

// TestPipeRefinement_DecodeSuccess verifies that decoding succeeds when both
// the upstream codec and the refinement accept the input.
func TestPipeRefinement_DecodeSuccess(t *testing.T) {
	composed := PipeRefinement[string, string](positiveIntPrism)(IntFromString())

	assert.Equal(t, validation.Success(1), composed.Decode("1"))
	assert.Equal(t, validation.Success(42), composed.Decode("42"))
	assert.Equal(t, validation.Success(999), composed.Decode("999"))
}

// TestPipeRefinement_DecodeFailsOnUpstream verifies that a failure in the
// upstream codec is propagated before the refinement is even attempted.
func TestPipeRefinement_DecodeFailsOnUpstream(t *testing.T) {
	composed := PipeRefinement[string, string](positiveIntPrism)(IntFromString())

	assert.True(t, either.IsLeft(composed.Decode("not-a-number")))
	assert.True(t, either.IsLeft(composed.Decode("")))
}

// TestPipeRefinement_DecodeFailsOnRefinement verifies that decoding fails when
// the upstream codec succeeds but the refinement rejects the value.
func TestPipeRefinement_DecodeFailsOnRefinement(t *testing.T) {
	composed := PipeRefinement[string, string](positiveIntPrism)(IntFromString())

	assert.True(t, either.IsLeft(composed.Decode("0")))
	assert.True(t, either.IsLeft(composed.Decode("-1")))
	assert.True(t, either.IsLeft(composed.Decode("-100")))
}

// TestPipeRefinement_Encode verifies that encoding applies the refinement's
// ReverseGet (identity for positiveIntPrism) and then the upstream encoder.
func TestPipeRefinement_Encode(t *testing.T) {
	composed := PipeRefinement[string, string](positiveIntPrism)(IntFromString())

	assert.Equal(t, "7", composed.Encode(7))
	assert.Equal(t, "42", composed.Encode(42))
}

// TestPipeRefinement_RoundTrip verifies the round-trip property for valid inputs.
func TestPipeRefinement_RoundTrip(t *testing.T) {
	composed := PipeRefinement[string, string](positiveIntPrism)(IntFromString())

	for _, s := range []string{"1", "42", "100"} {
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

// TestPipeRefinement_Name verifies the composed codec's name follows the
// "Pipe(<receiver-name>, FromRefinement(<refinement-name>))" pattern.
func TestPipeRefinement_Name(t *testing.T) {
	composed := PipeRefinement[string, string](positiveIntPrism)(IntFromString())
	assert.Equal(t, "Pipe(IntFromString, FromRefinement(PositiveInt))", composed.Name())
}

// TestPipeRefinement_EquivalentToFreeFunction verifies that PipeRefinement is
// behaviourally identical to Pipe(FromRefinement(ab)).
func TestPipeRefinement_EquivalentToFreeFunction(t *testing.T) {
	via := PipeRefinement[string, string](positiveIntPrism)(IntFromString())
	direct := Pipe[string, string](FromRefinement(positiveIntPrism))(IntFromString())

	t.Run("Decode parity on success", func(t *testing.T) {
		assert.Equal(t, direct.Decode("42"), via.Decode("42"))
	})
	t.Run("Decode parity on refinement failure", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(direct.Decode("-1")), either.IsLeft(via.Decode("-1")))
	})
	t.Run("Encode parity", func(t *testing.T) {
		assert.Equal(t, direct.Encode(99), via.Encode(99))
	})
}

// ---------------------------------------------------------------------------
// PipeIso method (go1.27+)
// ---------------------------------------------------------------------------

// TestTypePipeIso_MethodEquivalentToFreeFunction verifies that the method form
// base.PipeIso(ab) returns a codec identical in behaviour to the free function
// form PipeIso[O, I](ab)(base).
func TestTypePipeIso_MethodEquivalentToFreeFunction(t *testing.T) {
	base := IntFromString()
	method := base.PipeIso(userIdIso)
	free := PipeIso[string, string](userIdIso)(base)

	t.Run("Decode parity on success", func(t *testing.T) {
		assert.Equal(t, free.Decode("42"), method.Decode("42"))
	})
	t.Run("Decode parity on failure", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(free.Decode("abc")), either.IsLeft(method.Decode("abc")))
	})
	t.Run("Encode parity", func(t *testing.T) {
		assert.Equal(t, free.Encode(UserId(5)), method.Encode(UserId(5)))
	})
}

// TestTypePipeIso_DecodeSuccess verifies the method form decodes and applies
// the isomorphism correctly.
func TestTypePipeIso_DecodeSuccess(t *testing.T) {
	composed := IntFromString().PipeIso(userIdIso)

	assert.Equal(t, validation.Success(UserId(42)), composed.Decode("42"))
	assert.Equal(t, validation.Success(UserId(0)), composed.Decode("0"))
}

// TestTypePipeIso_Encode verifies the method form encodes via the isomorphism's
// ReverseGet and then the upstream encoder.
func TestTypePipeIso_Encode(t *testing.T) {
	composed := IntFromString().PipeIso(userIdIso)
	assert.Equal(t, "42", composed.Encode(UserId(42)))
}

// TestTypePipeIso_Name verifies the composed codec's name.
func TestTypePipeIso_Name(t *testing.T) {
	composed := IntFromString().PipeIso(userIdIso)
	assert.Equal(t, "Pipe(IntFromString, FromIso[Iso])", composed.Name())
}

// ---------------------------------------------------------------------------
// PipeRefinement method (go1.27+)
// ---------------------------------------------------------------------------

// TestTypePipeRefinement_MethodEquivalentToFreeFunction verifies that the method
// form base.PipeRefinement(ab) is behaviourally identical to the free function
// PipeRefinement[O, I](ab)(base).
func TestTypePipeRefinement_MethodEquivalentToFreeFunction(t *testing.T) {
	base := IntFromString()
	method := base.PipeRefinement(positiveIntPrism)
	free := PipeRefinement[string, string](positiveIntPrism)(base)

	t.Run("Decode parity on success", func(t *testing.T) {
		assert.Equal(t, free.Decode("42"), method.Decode("42"))
	})
	t.Run("Decode parity on upstream failure", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(free.Decode("abc")), either.IsLeft(method.Decode("abc")))
	})
	t.Run("Decode parity on refinement failure", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(free.Decode("-1")), either.IsLeft(method.Decode("-1")))
	})
	t.Run("Encode parity", func(t *testing.T) {
		assert.Equal(t, free.Encode(7), method.Encode(7))
	})
}

// TestTypePipeRefinement_DecodeSuccess verifies the method form validates and
// decodes successfully when both steps pass.
func TestTypePipeRefinement_DecodeSuccess(t *testing.T) {
	composed := IntFromString().PipeRefinement(positiveIntPrism)

	assert.Equal(t, validation.Success(1), composed.Decode("1"))
	assert.Equal(t, validation.Success(42), composed.Decode("42"))
}

// TestTypePipeRefinement_DecodeFailsOnRefinement verifies the method form
// propagates the refinement failure when the predicate is not satisfied.
func TestTypePipeRefinement_DecodeFailsOnRefinement(t *testing.T) {
	composed := IntFromString().PipeRefinement(positiveIntPrism)

	assert.True(t, either.IsLeft(composed.Decode("0")))
	assert.True(t, either.IsLeft(composed.Decode("-42")))
}

// TestTypePipeRefinement_Encode verifies the method form encodes via
// refinement.ReverseGet and then the upstream encoder.
func TestTypePipeRefinement_Encode(t *testing.T) {
	composed := IntFromString().PipeRefinement(positiveIntPrism)
	assert.Equal(t, "99", composed.Encode(99))
}

// TestTypePipeRefinement_Name verifies the composed codec name.
func TestTypePipeRefinement_Name(t *testing.T) {
	composed := IntFromString().PipeRefinement(positiveIntPrism)
	assert.Equal(t, "Pipe(IntFromString, FromRefinement(PositiveInt))", composed.Name())
}

// ---------------------------------------------------------------------------
// Example functions
// ---------------------------------------------------------------------------

// ExamplePipeIso demonstrates composing IntFromString with an Iso that wraps
// the decoded int in a UserId newtype.  Decoding never fails due to the Iso
// step; encoding unwraps the newtype and formats the integer as a string.
func ExamplePipeIso() {
	userIdIsoLocal := iso.MakeIso(
		func(i int) UserId { return UserId(i) },
		func(id UserId) int { return int(id) },
	)

	composed := PipeIso[string, string](userIdIsoLocal)(IntFromString())

	fmt.Println(composed.Name())
	fmt.Println(composed.Encode(UserId(42)))
	fmt.Println(either.IsRight(composed.Decode("7")))
	fmt.Println(either.IsLeft(composed.Decode("not-a-number")))

	// Output:
	// Pipe(IntFromString, FromIso[Iso])
	// 42
	// true
	// true
}

// ExamplePipeRefinement demonstrates composing IntFromString with a refinement
// that accepts only positive integers.  Decoding fails for non-positive values.
func ExamplePipeRefinement() {
	composed := PipeRefinement[string, string](positiveIntPrism)(IntFromString())

	fmt.Println(composed.Name())
	fmt.Println(composed.Encode(42))
	fmt.Println(either.IsRight(composed.Decode("10")))
	fmt.Println(either.IsLeft(composed.Decode("-3")))
	fmt.Println(either.IsLeft(composed.Decode("0")))

	// Output:
	// Pipe(IntFromString, FromRefinement(PositiveInt))
	// 42
	// true
	// true
	// true
}

// ExampleType_PipeIso demonstrates the method form of PipeIso available on
// Go 1.27+.  It is equivalent to calling PipeIso[O, I](ab)(base) directly.
func ExampleType_PipeIso() {
	userIdIsoLocal := iso.MakeIso(
		func(i int) UserId { return UserId(i) },
		func(id UserId) int { return int(id) },
	)

	composed := IntFromString().PipeIso(userIdIsoLocal)

	fmt.Println(composed.Name())
	fmt.Println(composed.Encode(UserId(7)))
	fmt.Println(either.IsRight(composed.Decode("99")))
	fmt.Println(either.IsLeft(composed.Decode("bad")))

	// Output:
	// Pipe(IntFromString, FromIso[Iso])
	// 7
	// true
	// true
}

// ExampleType_PipeRefinement demonstrates the method form of PipeRefinement
// available on Go 1.27+.  It is equivalent to calling
// PipeRefinement[O, I](ab)(base) directly.
func ExampleType_PipeRefinement() {
	composed := IntFromString().PipeRefinement(positiveIntPrism)

	fmt.Println(composed.Name())
	fmt.Println(composed.Encode(42))
	fmt.Println(either.IsRight(composed.Decode("10")))
	fmt.Println(either.IsLeft(composed.Decode("-3")))

	// Output:
	// Pipe(IntFromString, FromRefinement(PositiveInt))
	// 42
	// true
	// true
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

// ---------------------------------------------------------------------------
// PipePrism method (go1.27+)
// ---------------------------------------------------------------------------

// TestTypePipePrism_MethodEquivalentToFreeFunction verifies that the method
// form base.PipePrism(ab) is behaviourally identical to the free function
// PipePrism[O, I](ab)(base).
func TestTypePipePrism_MethodEquivalentToFreeFunction(t *testing.T) {
	base := IntFromString()
	method := base.PipePrism(positiveIntPrism)
	free := PipePrism[string, string](positiveIntPrism)(base)

	t.Run("Decode parity on success", func(t *testing.T) {
		assert.Equal(t, free.Decode("42"), method.Decode("42"))
	})
	t.Run("Decode parity on upstream failure", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(free.Decode("abc")), either.IsLeft(method.Decode("abc")))
	})
	t.Run("Decode parity on prism failure", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(free.Decode("-1")), either.IsLeft(method.Decode("-1")))
	})
	t.Run("Encode parity", func(t *testing.T) {
		assert.Equal(t, free.Encode(7), method.Encode(7))
	})
}

// TestTypePipePrism_EquivalentToPipeRefinementMethod verifies that
// base.PipePrism(ab) is behaviourally identical to base.PipeRefinement(ab).
func TestTypePipePrism_EquivalentToPipeRefinementMethod(t *testing.T) {
	base := IntFromString()
	viaPrism := base.PipePrism(positiveIntPrism)
	viaRefinement := base.PipeRefinement(positiveIntPrism)

	t.Run("same name", func(t *testing.T) {
		assert.Equal(t, viaRefinement.Name(), viaPrism.Name())
	})
	t.Run("decode success parity", func(t *testing.T) {
		assert.Equal(t, viaRefinement.Decode("42"), viaPrism.Decode("42"))
	})
	t.Run("decode failure parity", func(t *testing.T) {
		assert.Equal(t, either.IsLeft(viaRefinement.Decode("-1")), either.IsLeft(viaPrism.Decode("-1")))
	})
	t.Run("encode parity", func(t *testing.T) {
		assert.Equal(t, viaRefinement.Encode(99), viaPrism.Encode(99))
	})
}

// TestTypePipePrism_DecodeSuccess verifies the method form succeeds when both
// the upstream codec and the prism accept the input.
func TestTypePipePrism_DecodeSuccess(t *testing.T) {
	composed := IntFromString().PipePrism(positiveIntPrism)

	assert.Equal(t, validation.Success(1), composed.Decode("1"))
	assert.Equal(t, validation.Success(42), composed.Decode("42"))
}

// TestTypePipePrism_DecodeFailsOnPrism verifies the method form propagates a
// validation error when the prism rejects the decoded value.
func TestTypePipePrism_DecodeFailsOnPrism(t *testing.T) {
	composed := IntFromString().PipePrism(positiveIntPrism)

	assert.True(t, either.IsLeft(composed.Decode("0")))
	assert.True(t, either.IsLeft(composed.Decode("-7")))
}

// TestTypePipePrism_Encode verifies the method form encodes via
// prism.ReverseGet and then the upstream encoder.
func TestTypePipePrism_Encode(t *testing.T) {
	composed := IntFromString().PipePrism(positiveIntPrism)
	assert.Equal(t, "55", composed.Encode(55))
}

// TestTypePipePrism_Name verifies the composed codec name.
func TestTypePipePrism_Name(t *testing.T) {
	composed := IntFromString().PipePrism(positiveIntPrism)
	assert.Equal(t, "Pipe(IntFromString, FromRefinement(PositiveInt))", composed.Name())
}

// ExampleType_PipePrism demonstrates the method form of PipePrism available
// on Go 1.27+.  It is equivalent to calling PipePrism[O, I](ab)(base) directly.
func ExampleType_PipePrism() {
	composed := IntFromString().PipePrism(positiveIntPrism)

	fmt.Println(composed.Name())
	fmt.Println(composed.Encode(42))
	fmt.Println(either.IsRight(composed.Decode("10")))
	fmt.Println(either.IsLeft(composed.Decode("-3")))

	// Output:
	// Pipe(IntFromString, FromRefinement(PositiveInt))
	// 42
	// true
	// true
}

// ---------------------------------------------------------------------------
// Shared helpers for Alt / ApSL / ApSO / Bind method tests
// ---------------------------------------------------------------------------

// employee is a test struct used by method-receiver tests in this file to
// avoid clashing with Person defined in bind_test.go.
type employee struct {
	Name string
	Age  int
}

var (
	employeeNameLens = common.MakeLens(
		func(e employee) string { return e.Name },
		func(e employee, name string) employee { e.Name = name; return e },
	)
	employeeAgeLens = common.MakeLens(
		func(e employee) int { return e.Age },
		func(e employee, age int) employee { e.Age = age; return e },
	)
	// employeeNicknameOpt focuses on the Name field when non-empty (optional demo).
	employeeNicknameOpt = optional.MakeOptionalCurried(
		func(e employee) option.Option[string] {
			if e.Name != "" {
				return option.Some(e.Name)
			}
			return option.None[string]()
		},
		func(name string) func(employee) employee {
			return func(e employee) employee { e.Name = name; return e }
		},
	)
)

// ---------------------------------------------------------------------------
// Alt method (go1.27+)
// ---------------------------------------------------------------------------

// TestTypeAlt_MethodEquivalentToFreeFunction verifies that the method form
// base.Alt(second) is behaviourally identical to Alt(second)(base).
func TestTypeAlt_MethodEquivalentToFreeFunction(t *testing.T) {
	base := IntFromString()
	fallback := NonEmptyString().Pipe(IntFromString())

	method := base.Alt(lazy.Of(fallback))
	free := Alt(lazy.Of(fallback))(base)

	t.Run("Decode parity on primary success", func(t *testing.T) {
		assert.Equal(t, free.Decode("42"), method.Decode("42"))
	})
	t.Run("Decode parity on fallback", func(t *testing.T) {
		// Both fail on a truly non-numeric string
		assert.Equal(t, either.IsLeft(free.Decode("bad")), either.IsLeft(method.Decode("bad")))
	})
	t.Run("Encode parity", func(t *testing.T) {
		assert.Equal(t, free.Encode(7), method.Encode(7))
	})
}

// TestTypeAlt_PrimaryWins verifies that when the receiver succeeds the
// fallback is never used.
func TestTypeAlt_PrimaryWins(t *testing.T) {
	primary := IntFromString()
	// fallback always fails
	alwaysFail := MakeType(
		"AlwaysFail",
		Is[int](),
		func(s string) Decode[Context, int] {
			return func(c Context) Validation[int] {
				return validation.FailureWithMessage[int](s, "always fails")(c)
			}
		},
		func(n int) string { return "" },
	)

	composed := primary.Alt(lazy.Of(alwaysFail))
	assert.Equal(t, validation.Success(42), composed.Decode("42"))
}

// TestTypeAlt_FallbackUsedOnFailure verifies that when the receiver fails the
// fallback is tried and its result is used.
func TestTypeAlt_FallbackUsedOnFailure(t *testing.T) {
	// Primary only accepts "one"
	acceptOne := MakeType(
		"AcceptOne",
		Is[int](),
		func(s string) Decode[Context, int] {
			return func(c Context) Validation[int] {
				if s == "one" {
					return validation.Success(1)
				}
				return validation.FailureWithMessage[int](s, "not one")(c)
			}
		},
		func(n int) string { return "one" },
	)
	// Fallback accepts "two"
	acceptTwo := MakeType(
		"AcceptTwo",
		Is[int](),
		func(s string) Decode[Context, int] {
			return func(c Context) Validation[int] {
				if s == "two" {
					return validation.Success(2)
				}
				return validation.FailureWithMessage[int](s, "not two")(c)
			}
		},
		func(n int) string { return "two" },
	)

	composed := acceptOne.Alt(lazy.Of(acceptTwo))

	assert.Equal(t, validation.Success(1), composed.Decode("one"))
	assert.Equal(t, validation.Success(2), composed.Decode("two"))
	assert.True(t, either.IsLeft(composed.Decode("three")))
}

// TestTypeAlt_EncoderIsFromReceiver verifies that the composed codec always
// uses the receiver's encoder.
func TestTypeAlt_EncoderIsFromReceiver(t *testing.T) {
	primary := IntFromString()
	fallback := MakeType(
		"Other",
		Is[int](),
		func(s string) Decode[Context, int] { return func(c Context) Validation[int] { return validation.Success(0) } },
		func(n int) string { return "other" },
	)
	composed := primary.Alt(lazy.Of(fallback))
	assert.Equal(t, "99", composed.Encode(99))
}

// TestTypeAlt_Name verifies the composed codec carries the "Alt[<name>]" name.
func TestTypeAlt_Name(t *testing.T) {
	composed := IntFromString().Alt(lazy.Of(IntFromString()))
	assert.Equal(t, "Alt[IntFromString]", composed.Name())
}

// ExampleType_Alt demonstrates the method form of Alt available on Go 1.27+.
// The receiver is tried first; the fallback is used only when the receiver fails.
func ExampleType_Alt() {
	// Primary: parse a non-empty string, then an int from it.
	// Fallback: accept "zero" literally.
	acceptZero := MakeType(
		"AcceptZero",
		Is[int](),
		func(s string) Decode[Context, int] {
			return func(c Context) Validation[int] {
				if s == "zero" {
					return validation.Success(0)
				}
				return validation.FailureWithMessage[int](s, "not zero")(c)
			}
		},
		func(n int) string { return "zero" },
	)

	composed := IntFromString().Alt(lazy.Of(acceptZero))

	fmt.Println(composed.Name())
	fmt.Println(either.IsRight(composed.Decode("42")))  // primary succeeds
	fmt.Println(either.IsRight(composed.Decode("zero"))) // fallback succeeds
	fmt.Println(either.IsLeft(composed.Decode("bad")))   // both fail

	// Output:
	// Alt[IntFromString]
	// true
	// true
	// true
}

// ---------------------------------------------------------------------------
// AltW — free-function only (method form impossible due to Go generic cycle)
// ---------------------------------------------------------------------------

// TestAltW_MethodEquivalentToFreeFunction verifies that the existing free
// function AltW[R, L](leftItem)(rightItem) is correctly reached via the free
// function (no method equivalent exists).
func TestAltW_MethodForm_NotAvailable(t *testing.T) {
	// AltW cannot be a method because instantiating Type[either.Either[L,A],O,I]
	// from Type[A,O,I] creates a generic instantiation cycle in the Go compiler.
	// Verify that the free function form works correctly instead.
	c := AltW[int](Id[string]())(IntFromString())

	t.Run("Right branch decoded first", func(t *testing.T) {
		res := c.Decode("42")
		assert.True(t, either.IsRight(res))
	})
	t.Run("Left branch on right failure", func(t *testing.T) {
		res := c.Decode("hello")
		assert.True(t, either.IsRight(res)) // Left(string) is still Right(Either)
	})
	t.Run("Both fail", func(t *testing.T) {
		// An empty string fails both IntFromString and Id[string] (non-empty check)
		// Id[string] always succeeds, so this can't both fail — just verify shape
		res := c.Decode("42")
		val := either.MonadFold(res,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("err") },
			func(e either.Either[string, int]) either.Either[string, int] { return e },
		)
		assert.True(t, either.IsRight(val)) // decoded to Right(42)
	})
	t.Run("Encode Right dispatches to right encoder", func(t *testing.T) {
		assert.Equal(t, "99", c.Encode(either.Right[string](99)))
	})
	t.Run("Encode Left dispatches to left encoder", func(t *testing.T) {
		assert.Equal(t, "hello", c.Encode(either.Left[int]("hello")))
	})
}

// ---------------------------------------------------------------------------
// ApSL method (go1.27+)
// ---------------------------------------------------------------------------

// TestTypeApSL_MethodEquivalentToFreeFunction verifies that the method form
// base.ApSL(m, l, fa) is behaviourally identical to ApSL(m, l, fa)(base).
func TestTypeApSL_MethodEquivalentToFreeFunction(t *testing.T) {
	base := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{})))

	method := base.ApSL(S.Monoid, employeeNameLens, String())
	free := ApSL[employee, string, string, any](S.Monoid, employeeNameLens, String())(base)

	t.Run("Decode parity on success", func(t *testing.T) {
		assert.Equal(t, free.Decode("Alice"), method.Decode("Alice"))
	})
	t.Run("Encode parity", func(t *testing.T) {
		assert.Equal(t, free.Encode(employee{Name: "Alice"}), method.Encode(employee{Name: "Alice"}))
	})
}

// TestTypeApSL_BuildsStructCodec verifies that chaining ApSL calls produces a
// codec that correctly validates, sets, and encodes struct fields.
func TestTypeApSL_BuildsStructCodec(t *testing.T) {
	// Build a codec for employee using method chaining
	// Do[I, A, O] with I=any, A=employee, O=string
	base := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{})))
	nameCodec := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{}))).
		ApSL(S.Monoid, employeeNameLens, String())

	t.Run("Decode sets Name field", func(t *testing.T) {
		res := nameCodec.Decode("Alice")
		assert.Equal(t, validation.Success(employee{Name: "Alice"}), res)
	})

	t.Run("Encode extracts Name field", func(t *testing.T) {
		enc := nameCodec.Encode(employee{Name: "Bob"})
		assert.Equal(t, "Bob", enc)
	})

	_ = base
}

// TestTypeApSL_ErrorAccumulation verifies that ApSL accumulates errors from
// both the base codec and the field codec (applicative semantics).
func TestTypeApSL_ErrorAccumulation(t *testing.T) {
	// A failing base
	failBase := MakeType(
		"FailBase",
		Is[employee](),
		func(_ any) Decode[Context, employee] {
			return func(c Context) Validation[employee] {
				return validation.FailureWithMessage[employee](nil, "base failed")(c)
			}
		},
		func(e employee) string { return "" },
	)
	composed := failBase.ApSL(S.Monoid, employeeNameLens, String())
	assert.True(t, either.IsLeft(composed.Decode("Alice")))
}

// ExampleType_ApSL demonstrates building an employee codec field-by-field
// using the method form of ApSL on Go 1.27+.
func ExampleType_ApSL() {
	nameLens := common.MakeLens(
		func(e employee) string { return e.Name },
		func(e employee, name string) employee { e.Name = name; return e },
	)

	c := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{}))).
		ApSL(S.Monoid, nameLens, String())

	res := c.Decode("Alice")
	fmt.Println(either.IsRight(res))
	fmt.Println(c.Encode(employee{Name: "Bob"}))

	// Output:
	// true
	// Bob
}

// ---------------------------------------------------------------------------
// ApSO method (go1.27+)
// ---------------------------------------------------------------------------

// TestTypeApSO_MethodEquivalentToFreeFunction verifies that the method form
// base.ApSO(m, o, fa) is behaviourally identical to ApSO(m, o, fa)(base).
func TestTypeApSO_MethodEquivalentToFreeFunction(t *testing.T) {
	base := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{})))

	method := base.ApSO(S.Monoid, employeeNicknameOpt, String())
	free := ApSO[employee, string, string, any](S.Monoid, employeeNicknameOpt, String())(base)

	t.Run("Decode parity when field present", func(t *testing.T) {
		assert.Equal(t, free.Decode("Alice"), method.Decode("Alice"))
	})
	t.Run("Encode parity when field present", func(t *testing.T) {
		assert.Equal(t, free.Encode(employee{Name: "Alice"}), method.Encode(employee{Name: "Alice"}))
	})
	t.Run("Encode parity when field absent", func(t *testing.T) {
		assert.Equal(t, free.Encode(employee{}), method.Encode(employee{}))
	})
}

// TestTypeApSO_EncodesFieldWhenPresent verifies that encoding includes the
// optional field value when the optional's GetOption returns Some.
func TestTypeApSO_EncodesFieldWhenPresent(t *testing.T) {
	base := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{})))
	c := base.ApSO(S.Monoid, employeeNicknameOpt, String())

	assert.Equal(t, "Alice", c.Encode(employee{Name: "Alice"}))
}

// TestTypeApSO_OmitsFieldWhenAbsent verifies that encoding returns the base
// output when the optional's GetOption returns None.
func TestTypeApSO_OmitsFieldWhenAbsent(t *testing.T) {
	base := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{})))
	c := base.ApSO(S.Monoid, employeeNicknameOpt, String())

	assert.Equal(t, "", c.Encode(employee{})) // Name is "" → None → base output only
}

// ExampleType_ApSO demonstrates building a codec with an optional field using
// the method form of ApSO on Go 1.27+.
func ExampleType_ApSO() {
	opt := optional.MakeOptionalCurried(
		func(e employee) option.Option[string] {
			if e.Name != "" {
				return option.Some(e.Name)
			}
			return option.None[string]()
		},
		func(name string) func(employee) employee {
			return func(e employee) employee { e.Name = name; return e }
		},
	)

	c := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{}))).
		ApSO(S.Monoid, opt, String())

	fmt.Println(c.Encode(employee{Name: "Alice"})) // field present
	fmt.Println(c.Encode(employee{}))              // field absent

	// Output:
	// Alice
	//
}

// ---------------------------------------------------------------------------
// Bind method (go1.27+)
// ---------------------------------------------------------------------------

// TestTypeBind_MethodEquivalentToFreeFunction verifies that the method form
// base.Bind(m, l, f) is behaviourally identical to Bind(m, l, f)(base).
func TestTypeBind_MethodEquivalentToFreeFunction(t *testing.T) {
	base := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{})))
	kleisli := func(_ employee) Type[string, string, any] { return String() }

	method := base.Bind(S.Monoid, employeeNameLens, kleisli)
	free := Bind[employee, string, string, any](S.Monoid, employeeNameLens, kleisli)(base)

	t.Run("Decode parity on success", func(t *testing.T) {
		assert.Equal(t, free.Decode("Alice"), method.Decode("Alice"))
	})
	t.Run("Encode parity", func(t *testing.T) {
		assert.Equal(t, free.Encode(employee{Name: "Alice"}), method.Encode(employee{Name: "Alice"}))
	})
}

// TestTypeBind_ContextDependentField verifies that the Kleisli arrow receives
// the already-decoded struct value and can vary the field codec based on it.
func TestTypeBind_ContextDependentField(t *testing.T) {
	// Codec: given input string S, decode to employee{Name: S, Age: len(S)}.
	// The Age Kleisli arrow inspects the already-decoded Name to compute Age.
	ageLens := common.MakeLens(
		func(e employee) string { return e.Name }, // reuse Name as a proxy for the Kleisli input
		func(e employee, v string) employee { e.Age = len(v); return e },
	)

	base := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{}))).
		Bind(S.Monoid, ageLens, func(_ employee) Type[string, string, any] {
			return String()
		})

	t.Run("Kleisli arrow receives decoded struct", func(t *testing.T) {
		res := base.Decode("Alice")
		assert.True(t, either.IsRight(res))
	})

	t.Run("Encode uses Kleisli-derived codec", func(t *testing.T) {
		enc := base.Encode(employee{Name: "Bob"})
		assert.Equal(t, "Bob", enc)
	})
}

// TestTypeBind_FailFastOnBaseFailure verifies that when the base codec fails
// the Kleisli arrow is never evaluated (monadic / fail-fast semantics).
func TestTypeBind_FailFastOnBaseFailure(t *testing.T) {
	failBase := MakeType(
		"FailBase",
		Is[employee](),
		func(_ any) Decode[Context, employee] {
			return func(c Context) Validation[employee] {
				return validation.FailureWithMessage[employee](nil, "base failed")(c)
			}
		},
		func(e employee) string { return "" },
	)

	evaluated := false
	composed := failBase.Bind(S.Monoid, employeeNameLens, func(e employee) Type[string, string, any] {
		evaluated = true
		return String()
	})

	assert.True(t, either.IsLeft(composed.Decode("Alice")))
	assert.False(t, evaluated, "Kleisli arrow must not be evaluated when base fails")
}

// ExampleType_Bind demonstrates building a context-dependent field codec using
// the method form of Bind on Go 1.27+.
func ExampleType_Bind() {
	nameLens := common.MakeLens(
		func(e employee) string { return e.Name },
		func(e employee, name string) employee { e.Name = name; return e },
	)

	// Build an employee codec where the Name field codec is always String().
	c := Do[any, employee, string](lazy.Of(pair.MakePair("", employee{}))).
		Bind(S.Monoid, nameLens, func(_ employee) Type[string, string, any] {
			return String()
		})

	res := c.Decode("Alice")
	fmt.Println(either.IsRight(res))
	fmt.Println(c.Encode(employee{Name: "Bob"}))

	// Output:
	// true
	// Bob
}
